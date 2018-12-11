package blacklist

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"

	guard "grapeGuard"

	log "grapeLoggers/clientv1"

	etcd "github.com/koangel/grapeNet/Etcd"
	utils "github.com/koangel/grapeNet/Utils"
)

////////////////////////////////////////////////////////////
//  add by koangel
//	date: 2017/09/27
//  黑名单体系，通过etcd的集群同步
////////////////////////////////////////////////////////////

var (
	bl sync.Map

	blqueue = utils.NewSQueue()
	blChans = make(chan *BlackNode, 320000)
)

const (
	EtcdKey = "blacklist"

	trimPrefix = "blacklist/"
)

type BlackListItem struct {
	IPAddr string
	Scheme string
	Host   string
	Msg    string
	Timout int64
	Ttl    int64
}

type BlackNode struct {
	TTL    int64  `json:"ttl"`
	Msg    string `json:"msg"`
	Scheme string `json:"scheme"`
	Host   string `json:"host"` // 攻击的谁
}

func init() {
	for i := 0; i < 8; i++ {
		go procBlackListPush()
	}
}

func procBlackListPush() {
	for {
		items := blqueue.Pop().(*BlackListItem)
		if items == nil {
			continue
		}

		if useIPSet != guard.IsIPSet() && guard.IsIPSet() {
			if err := checkIPSet(); err == nil {
				createIPSet()
			}
			useIPSet = guard.IsIPSet()
		}

		itemHost := &BlackNode{
			TTL:    items.Timout,
			Scheme: items.Scheme,
			Host:   items.Host,
			Msg:    items.Msg,
		}
		bl.Store(items.IPAddr, itemHost)

		if etcd.EtcdCli == nil {
			continue
		}

		if items.Ttl < 1000 {
			items.Ttl = 1000
		}

		log.WithField("type", "BlackList").Info("新增黑名单IP:", items.IPAddr, ",封停原因:", items.Msg)

		data, verr := json.Marshal(itemHost)

		if verr != nil {
			log.WithField("type", "BlackList").Info("推入黑名单错误：%v", verr)
			continue
		}

		_, err := etcd.WriteTTL(fmt.Sprint(EtcdKey, "/", items.IPAddr), data, items.Ttl)
		if err != nil {
			log.Error("添加TTL错误:", err)
			continue
		}
	}
}

func GetIPItem(conn net.Conn) (*BlackNode, string) {
	remote := GetRealIp(conn.RemoteAddr().String())
	item, ok := bl.Load(remote)
	if ok {
		return item.(*BlackNode), remote
	}

	return nil, remote
}

// 监听黑名单变化
func BlackListWatcher(ev *clientv3.Event) {

	keystr := strings.TrimPrefix(string(ev.Kv.Key), trimPrefix)
	//log.WithField("type", "BlackList").Info("黑名单数据修改,类型:", vtype, ",Key:", keystr)
	switch ev.Type.String() {
	case "PUT":
		if IsBlackList(keystr) {
			return
		}

		ttls, err := etcd.TimeToLive(clientv3.LeaseID(ev.Kv.Lease))

		if err != nil {
			log.WithField("type", "BlackList").Info("读取黑名单错误：", err)
			return
		}

		if ttls.TTL == -1 {
			log.WithField("type", "BlackList").Info("到期黑名单无需添加。")
			return
		}

		item := &BlackNode{}
		err = json.Unmarshal(ev.Kv.Value, item)
		if err != nil {
			bl.Store(keystr, &BlackNode{
				TTL:    time.Now().Add(time.Duration(ttls.TTL) * time.Second).Unix(),
				Host:   "",
				Scheme: "",
				Msg:    "unknow",
			})
		} else {
			bl.Store(keystr, &BlackNode{
				TTL:    time.Now().Add(time.Duration(ttls.TTL) * time.Second).Unix(),
				Host:   item.Host,
				Scheme: item.Scheme,
				Msg:    item.Msg,
			})
		}

		addIpset(keystr, int(ttls.TTL))

		//log.WithField("type", "BlackList").Info("新增黑名单IP:", keystr, ",封停second:", ttls.TTL, ",原因：", item.Msg)

	case "DELETE":

		unixTime, has := bl.Load(keystr)
		if !has {
			return
		}

		if time.Now().Unix() <= unixTime.(*BlackNode).TTL {
			log.WithField("type", "BlackList").Info("黑名单IP移除:", keystr, "但是未到时间。")
			return
		}

		removeIpset(keystr)
		//log.WithField("type", "BlackList").Info("黑名单IP移除:", keystr)
		bl.Delete(keystr)
	}
}

func IsBlackList(key string) bool {
	_, find := bl.Load(key)
	return find
}

func DeleteBlackIP(key string) {
	_, has := bl.Load(key)
	if !has {
		return
	}

	log.WithField("type", "BlackList").Info("手动黑名单IP移除:", key)
	bl.Delete(key)
}

func Push2EtchR(req *http.Request, errmsg string, ttl int64) {
	reqIP := GetIP(req)
	if IsBlackList(reqIP) {
		return
	}

	//log.WithField("type", "BlackList").Info("新增黑名单IP:", reqIP, ",封停原因:", errmsg)
	blqueue.Push(&BlackListItem{
		IPAddr: reqIP,
		Scheme: utils.Ifs(req.TLS == nil, "http", "https"),
		Host:   req.Host,
		Msg:    errmsg,
		Timout: time.Now().Add(time.Duration(ttl) * time.Second).Unix(),
		Ttl:    ttl,
	})
}

func Load4Etcd() {
	useIPSet = guard.IsIPSet()
	if useIPSet {
		log.Info("启动IPSet模式...")
		err := checkIPSet()
		if err != nil {
			log.Error("IPSet未安装或异常，无法启动。")
		}

		if useIPSet {
			log.Info("初始化Ipset表...")
			createIPSet()
		}
	}

	log.Info("载入集群黑名单数据...")
	resp, err := etcd.ReadAll(EtcdKey)
	if err != nil {
		log.Error("载入黑名单数据错误:", err)
		return
	}

	for _, v := range resp.Kvs {
		key := strings.TrimPrefix(string(v.Key), trimPrefix)

		ttls, err := etcd.TimeToLive(clientv3.LeaseID(v.Lease))

		if err != nil {
			log.WithField("type", "BlackList").Info("读取黑名单错误：", err)
			return
		}

		if ttls.TTL == -1 {
			log.WithField("type", "BlackList").Info("到期黑名单无需添加。")
			return
		}

		item := &BlackNode{}
		err = json.Unmarshal(v.Value, item)
		if err != nil {
			bl.Store(key, &BlackNode{
				TTL:    time.Now().Add(time.Duration(ttls.TTL) * time.Second).Unix(),
				Host:   "",
				Scheme: "",
				Msg:    "unknow",
			})
		} else {
			bl.Store(key, &BlackNode{
				TTL:    time.Now().Add(time.Duration(ttls.TTL) * time.Second).Unix(),
				Host:   item.Host,
				Scheme: item.Scheme,
				Msg:    item.Msg,
			})
		}
	}

	log.WithField("type", "BlackList").Info("载入黑名单成功。")
}

func GetRealIp(src string) string {
	if src == "::1" {
		return "127.0.0.1"
	}

	pos := strings.LastIndex(src, ":")
	if pos != -1 {
		return src[0:pos]
	}

	return src
}

func GetIP(req *http.Request) string {
	return GetRealIp(req.RemoteAddr)
}

func BlackListConn(conn net.Conn) error {
	remote := GetRealIp(conn.RemoteAddr().String())
	_, ok := bl.Load(remote)
	if ok {
		return fmt.Errorf("黑名单IP：%v", remote)
	}

	return nil
}

func BlackListRequest(w http.ResponseWriter, remoteIP string) (int, error) {

	//if remoteIP == "127.0.0.1" {
	//	return 200, nil
	//}

	ttl, ok := bl.Load(remoteIP)
	if ok {
		if time.Now().Unix() >= ttl.(*BlackNode).TTL {
			log.WithField("type", "BlackList").Info("黑名单IP移除:", remoteIP, ",超时。")
			bl.Delete(remoteIP)
		}

		return 503, fmt.Errorf("黑名单IP：%v", remoteIP)
	}

	return 200, nil
}
