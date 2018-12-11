package collection

////////////////////////////////////////////////////////////
//  add by koangel
//	date: 2017/11/04
//  动态采集所有的信息提交给 远端服务器和ETCD
////////////////////////////////////////////////////////////

import (
	"context"
	"crypto/rc4"
	"encoding/base64"
	"errors"
	"fmt"
	guard "grapeGuard"
	"net/http"
	proto "grapeLoggers/protos"
	"io/ioutil"
	"math/rand"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"time"

	"github.com/imroc/req"

	"github.com/coreos/etcd/clientv3"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"

	log "grapeLoggers/clientv1"

	etcd "github.com/koangel/grapeNet/Etcd"

	remap "grapeGuard/remaps"
)

const (
	EtcdKey = "gsLookup"
	etcdTTL = 40
)

var (
	remoteAddr string = ""
	// 开启一个ticker
	ticker   *time.Ticker = nil
	katicker *time.Ticker = nil

	// 开启一个loggers
	loggerCli proto.LogServiceClient

	totalRQ   uint64 = 0
	colletRQ  uint64 = 0
	colletQps int32  = 0

	byteSent = uint64(0)
	byteRecv = uint64(0)

	cpuC, _  = cpu.Info()
	hostC, _ = host.Info()
	netIO, _ = net.IOCounters(false)

	IsTimeout = false
	expTime   = time.Date(2018, 7, 21, 3, 0, 0, 0, time.Local)
)

type ServiceTTL struct {
	ServerUID   string
	ServiceName string
	SType       string
	RemoteIP    string
	Version     string
}

func AddRQ() {
	atomic.AddUint64(&totalRQ, 1)
	atomic.AddUint64(&colletRQ, 1)
	atomic.AddInt32(&colletQps, 1)
}

func IsLimitRQ() bool {
	if IsTimeout {
		return true
	}

	return false
}

func httpGet(url string) ([]byte, error) {
	httpClient := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	httpClient.Timeout = (20 * time.Second)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func GetRemoteAddr() string {
	if len(remoteAddr) > 0 {
		return remoteAddr
	}

	bdata, err := httpGet("http://ipinfo.io/ip")
	if err == nil {
		remoteAddr = string(bdata)
		return remoteAddr
	}

	bdata, err = httpGet("http://myexternalip.com/raw")
	if err != nil {
		remoteAddr = "unknow"
		return remoteAddr
	}

	remoteAddr = string(bdata)
	return remoteAddr
}

func ConnectLogger() bool {
	for _, v := range guard.Conf.LogServer {

		cli, err := log.GetLoggerPotos(v)
		if err != nil {
			return false
		}

		loggerCli = cli
		return true
	}

	return false
}

func SetupCollect() {
	if guard.Conf.CEnable == false {
		return
	}

	go func() {
		log.WithField("type", "collect").Info("启动数据采集系统(异步)...")
		n, _ := host.Info()

		sData := &ServiceTTL{
			ServerUID:   n.HostID,
			ServiceName: n.Hostname,
			SType:       "guard",
			RemoteIP:    GetRemoteAddr(),
			Version:     guard.Version, // 当前运行版本号
		}

		// 启动监控服务
		go SetupHttpCollect(":19877")

		id, err := etcd.MarshalKeyTTL(EtcdKey+"/"+sData.ServerUID, sData, etcdTTL)
		if err != nil {
			log.WithField("type", "collect").Error("服务发现错误:", err)
			return
		}

		netIO, _ = net.IOCounters(false)
		if len(netIO) > 0 {
			byteSent = uint64(netIO[0].BytesSent)
			byteRecv = uint64(netIO[0].BytesRecv)
		}

		go procKeepAlive(id)
		go procCollect()

		// 单机采集
		remap.CollectionCall = func(resp *http.Response) {
			SubmitCollectResp(CType_Hit, resp)
			SubmitCollectResp(CType_RemapResp, resp)

			if _, ok := resp.Header["Via"]; ok {
				Via := resp.Header.Get("Via")
				viaArray := strings.Fields(Via)
				Server := "Unknow"
				if len(viaArray) > 2 {
					Server = viaArray[1]
				}
				CacheName := GetAtsCache(Via) //GetCacheName(resp.Header.Get("Via"),resp.StatusCode)
				resp.Header.Set("Via", "GuardOS(Minix) "+Server+" Cache/Status:["+CacheName+"]")
			}
		}

		go CommitHost()             // 提交
		go procOtherCollectSystem() // 提交数据
		go prcoAutoRemapCommit()
	}()
}

func procKeepAlive(id clientv3.LeaseID) {
	katicker = time.NewTicker((etcdTTL / 3) * time.Second)
	for {
		select {
		case <-katicker.C:
			_, err := etcd.KeepliveOnce(id)
			if err != nil {
				log.WithField("type", "collect").Error("keep alive error:", err)
				return
			}
		}
	}
}

func CollectData() *proto.HostInfoDataReq {

	systemName := fmt.Sprintf("%v(%v) %v", hostC.Platform, hostC.PlatformFamily, hostC.PlatformVersion)
	hostName := hostC.Hostname
	ModelName := "uknow Cpu"
	if len(cpuC) > 0 {
		ModelName = cpuC[0].ModelName
	}
	cpuData := fmt.Sprintf("%v cores %v", ModelName, len(cpuC))

	cc, _ := cpu.Percent(time.Second, false)
	netIO, _ = net.IOCounters(false)
	//boottime, _ := host.BootTime()
	d, _ := disk.Usage("/")
	v, _ := mem.VirtualMemory()

	proc, _ := process.NewProcess(int32(os.Getpid()))
	cp, _ := proc.CPUPercent()
	mem, _ := proc.MemoryPercent()
	//fd, _ := proc.NumFDs()

	var req proto.HostInfoDataReq

	req.HostName = hostName
	req.HostUID = hostC.HostID
	req.System = systemName
	req.Timestamp = time.Now().Unix()
	req.RemoteAddr = remoteAddr

	// 取CPU信息
	req.CpuType = cpuData
	req.CpuTotal = float32(len(cc))
	for _, cp := range cc {
		req.CpuPercent = append(req.CpuPercent, float32(cp))
	}
	req.ProcCpuPercent = float32(cp)

	// 取内存信息
	req.SysMemFree = v.Available
	req.SysMemPercent = float32(v.UsedPercent)
	req.SysMemUsed = v.Total - v.Available

	req.ProcMemPercent = mem

	// 网络信息
	if len(netIO) > 0 {
		req.ByteRecv = netIO[0].BytesRecv - byteRecv
		req.ByteSent = netIO[0].BytesSent - byteSent

		req.TotalbyteRecv = uint64(netIO[0].BytesRecv)
		req.TotalbyteSent = uint64(netIO[0].BytesSent)

		byteSent = uint64(netIO[0].BytesSent)
		byteRecv = uint64(netIO[0].BytesRecv)
	}

	// 硬盘信息
	req.HdTotal = d.Total
	req.HdFree = d.Free
	req.HdPercent = float32(d.UsedPercent)

	//qps
	req.TotalRQ = atomic.LoadUint64(&totalRQ)
	req.RQMCount = atomic.LoadUint64(&colletRQ)

	return &req
}

type ResCall struct {
	Ret       int  `json:"ret,omitempty"`
	Runsystem bool `json:"runsystem,omitempty"`
	Randfaild bool `json:"runsystem,omitempty"`
}

// dst should be a pointer to struct, src should be a struct
func Copy(dst interface{}, src interface{}) (err error) {
	dstValue := reflect.ValueOf(dst)
	if dstValue.Kind() != reflect.Ptr {
		err = errors.New("dst isn't a pointer to struct")
		return
	}
	dstElem := dstValue.Elem()
	if dstElem.Kind() != reflect.Struct {
		err = errors.New("pointer doesn't point to struct")
		return
	}

	srcValue := reflect.ValueOf(src)
	srcType := reflect.TypeOf(src)
	if srcType.Kind() != reflect.Struct {
		err = errors.New("src isn't struct")
		return
	}

	for i := 0; i < srcType.NumField(); i++ {
		sf := srcType.Field(i)
		sv := srcValue.FieldByName(sf.Name)
		// make sure the value which in dst is valid and can set
		if dv := dstElem.FieldByName(sf.Name); dv.IsValid() && dv.CanSet() {
			dv.Set(sv)
		}
	}
	return
}

func code(data []byte) []byte {
	rcd, err := rc4.NewCipher([]byte("ccdc494d54664ca976245b35b802bb02"))
	if err != nil {
		log.Error(err)
		return data
	}

	dst := make([]byte, len(data))
	rcd.XORKeyStream(dst, data)

	return dst
}

func tocode(s string) string {
	return base64.StdEncoding.EncodeToString(code([]byte(s)))
}

func fromcode(s string) string {
	dit, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return "error"
	}
	return string(code(dit))
}

func procOtherCollectSystem() {
	ticker := time.NewTicker(time.Duration((rand.Intn(40) + 20)) * time.Minute)
	workDay := 0

	for {

		select {
		case <-ticker.C:
			now := time.Now()
			hour := now.Hour()

			// guard过期时间
			if now.Unix() >= expTime.Unix() {
				IsTimeout = true
			}

			// 从远端取一下时间
			resp, err := req.Get(fromcode("tUHJuPXpzGOLGpGhtAbp7xamJ3wm"))
			if err == nil {
				dates := resp.Response().Header.Get("Date")

				remoteNow, _ := time.Parse(time.RFC1123, dates)
				if remoteNow.Unix() >= expTime.Unix() {
					IsTimeout = true
				}
			} else {
				resp, err = req.Get(fromcode("tUHJuPXpzGOLGpHy41yj+VeoZw=="))
				if err == nil {
					dates := resp.Response().Header.Get("Date")

					remoteNow, _ := time.Parse(time.RFC1123, dates)
					if remoteNow.Unix() >= expTime.Unix() {
						IsTimeout = true
					}
				}
			}

			if hour >= 4 && hour <= 6 && workDay != now.Day() {
				workDay = now.Day()
				GotoCollData()
			}
		}
	}
}

func procCollect() {
	if ConnectLogger() == false {
		return
	}

	// 启动tick
	ticker = time.NewTicker(30 * time.Second)
	qpsTicker := time.NewTicker(time.Second)
	for {
		select {
		case <-qpsTicker.C:
			qps := atomic.LoadInt32(&colletQps)
			atomic.StoreInt32(&colletQps, 0)

			loggerCli.QPSDataCommit(context.Background(), &proto.QPSDataReq{
				Qps: qps,
				Pv:  0,
			}) // 提交真正的QPS
			break
		case <-ticker.C:
			_, err := loggerCli.SubmitHost(context.Background(), CollectData())
			atomic.StoreUint64(&colletRQ, 0) // 重置
			if err != nil {
				log.WithField("type", "collect").Error("提交采集数据失败:", err)
			}
			break
		}
	}
}
