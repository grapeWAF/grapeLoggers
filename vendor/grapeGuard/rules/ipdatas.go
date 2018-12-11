package rules

import (
	"container/list"
	"net/http"
	"sync"
	"time"

	bl "grapeGuard/blacklist"

	"github.com/avct/uasurfer"

	utils "github.com/koangel/grapeNet/Utils"
)

////////////////////////////////////////////////////////////
//  add by koangel
//	date: 2017/12/05
//  ip数据统计，根据IP数据查看其浏览器和相关的内容
////////////////////////////////////////////////////////////
// 对每个来源IP的数据分析
type IPData struct {
	// 最后的15次请求内容是什么
	ReqPaths *StringLimit
	// 最后一次请求的数据是什么
	LastPath string
	// 最后一次请求的UA是什么
	UserAgent string
	// 请求设备是什么
	Device     string
	DeviceType uasurfer.DeviceType
	UA         uasurfer.UserAgent
	// 总IP请求数量
	TotalRQ int64
	// 锁
	mux sync.RWMutex

	lastUpdate time.Time
}

type IPQueueItem struct {
	remoteIP  string
	userAgent string
	path      string
}

var (
	IPMap sync.Map

	ipqueue = utils.NewSQueue()
)

func init() {
	for i := 0; i < 3; i++ {
		go procIPRequest()
	}

	go procIPDataTicker()
}

func SubmitRequest(req *http.Request, remoteIP string) {

	iQueue := &IPQueueItem{
		remoteIP:  remoteIP,
		userAgent: req.UserAgent(),
		path:      req.URL.Path,
	}

	ipqueue.Push(iQueue)
}

func procIPDataTicker() {
	checkTicker := time.NewTicker(time.Minute)
	for {
		select {
		case <-checkTicker.C:

			var listKey *list.List = list.New()
			nowTime := time.Now()

			IPMap.Range(func(key, value interface{}) bool {
				valData := value.(*IPData)
				if nowTime.Unix() >= valData.lastUpdate.Unix() {
					listKey.PushBack(key)
				}

				return true
			})

			for e := listKey.Front(); e != nil; e = e.Next() {
				IPMap.Delete(e.Value) // 删除所有ip
			}
			break
		}
	}
}

func procIPRequest() {
	for {
		item := ipqueue.Pop()
		ParserRequest(item.(*IPQueueItem)) // 进行一次记录
	}
}

func ParserRequest(item *IPQueueItem) *IPData {
	ua := uasurfer.Parse(item.userAgent)
	ipdata, has := IPMap.Load(item.remoteIP)
	if !has {
		newData := &IPData{
			ReqPaths:   NewSL(45),
			LastPath:   item.path,
			UserAgent:  item.userAgent,
			Device:     ua.DeviceType.String(),
			DeviceType: ua.DeviceType,
			lastUpdate: time.Now().Add(25 * time.Minute),
			UA:         *ua,
			TotalRQ:    1,
		}

		IPMap.Store(item.remoteIP, newData)
		return newData
	}

	ipOld := ipdata.(*IPData)

	ipOld.mux.Lock()
	ipOld.Device = ua.DeviceType.String()
	ipOld.DeviceType = ua.DeviceType
	ipOld.LastPath = item.path
	ipOld.UserAgent = item.userAgent
	ipOld.lastUpdate = time.Now().Add(25 * time.Minute)
	ipOld.ReqPaths.Push(ipOld.LastPath)
	ipOld.UA = *ua
	ipOld.TotalRQ++
	ipOld.mux.Unlock()

	return ipOld
}

func GetIP4A(addr string) *IPData {
	ipdata, has := IPMap.Load(addr)
	if !has {
		return nil
	}
	return ipdata.(*IPData)
}

func GetIPData(req *http.Request) *IPData {
	return GetIP4A(bl.GetIP(req))
}

func (ip *IPData) IsMobile() bool {
	defer ip.mux.RUnlock()
	ip.mux.RLock()
	return (ip.DeviceType == uasurfer.DevicePhone || ip.DeviceType == uasurfer.DeviceTablet)
}

func (ip *IPData) IsPC() bool {
	defer ip.mux.RUnlock()
	ip.mux.RLock()
	return (ip.DeviceType == uasurfer.DeviceComputer || ip.DeviceType == uasurfer.DeviceUnknown)
}

func (ip *IPData) IsSomePath(limit int) bool {
	defer ip.mux.RUnlock()
	ip.mux.RLock()
	return ip.ReqPaths.Match(ip.LastPath, limit)
}

func (ip *IPData) IsLinePath(limit int) bool {
	defer ip.mux.RUnlock()
	ip.mux.RLock()

	return ip.ReqPaths.LineMatch(limit)
}

func (ip *IPData) IsLinePrefix(limit int) bool {
	defer ip.mux.RUnlock()
	ip.mux.RLock()

	return ip.ReqPaths.MatchPrefix(ip.LastPath, limit)
}

type UAType map[uasurfer.BrowserName]uasurfer.Version

var (
	badUserType UAType = UAType{
		uasurfer.BrowserChrome: uasurfer.Version{32, 0, 0},
		uasurfer.BrowserSafari: uasurfer.Version{4, 0, 0},
		//uasurfer.BrowserIE:      uasurfer.Version{9, 0, 0},
		uasurfer.BrowserFirefox: uasurfer.Version{23, 0, 0},
	}
)

func (ip *IPData) IsBadBrowser() bool {
	defer ip.mux.RUnlock()
	ip.mux.RLock()

	ver, has := badUserType[ip.UA.Browser.Name]
	if !has {
		return false
	}

	if ip.UA.Browser.Version.Less(ver) {
		return true
	}

	return false
}

func (ip *IPData) GetTotalRQ() int64 {
	defer ip.mux.RUnlock()
	ip.mux.RLock()
	return ip.TotalRQ
}
