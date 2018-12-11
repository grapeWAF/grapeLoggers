package remaps

////////////////////////////////////////////////////////////
//  add by koangel
//	date: 2017/09/25
//  反向代理的基础remap库,配置文件
////////////////////////////////////////////////////////////

import (
	"encoding/xml"
	guard "grapeGuard"
	log "grapeLoggers/clientv1"
	"io/ioutil"
	"reflect"

	etcd "github.com/koangel/grapeNet/Etcd"

	util "github.com/koangel/grapeNet/Utils"
)

const (
	Target_Single  = iota
	Target_Round   // 轮询
	Target_Weight  // 权重
	Target_Fair    // 根据PING调整权重
	Target_IPHash  // 客户来源IPHASH
	Target_UrlHash // 根据URL进行HASH多源
)

const (
	MainTypeHash = iota
	MainTypeRoundRobin
	MainTypeNormal
)

const (
	tagOptName = "optData"
	tagName    = "remaps"
)

type SingleUriXml struct {
	Uri      string `xml:"Uri,attr" json:"uri"`
	Https    bool   `xml:"Https,attr" json:"https"`
	Backup   bool   `xml:"Backup,attr" json:"backup"`
	Weight   int    `xml:"Weight,attr" json:"weight"`
	TryCount int    `xml:"TryCount,attr" json:"trycount"`
	Timeout  int    `xml:"FaildTimeout,attr" json:"timeout"`
}

// 目标类型系统
type TargetUriXml struct {
	// 主类型 HASH或轮询
	MType int `xml:"MType,attr" json:"mtype"`
	// 子类型
	Type int `xml:"Type,attr" json:"type"`
	// 开启复杂轮询
	Open bool `xml:"Open,attr" json:"open"`

	// 底层开启
	Https       bool `xml:"Https,attr" json:"https"`
	HealthCheck bool `xml:"HealthCheck,attr" json:"health"`
	Follow      bool `xml:"Follow,attr" json:"follow"`

	Path string `xml:"Path,attr" json:"path"`

	// 目标数据点
	Uris []SingleUriXml `xml:"TargetUri" json:"uris"`
}

type SSLTypeXml struct {
	MinVersion uint16 `xml:"MinVersion,attr" json:"MinVersion,omitempty"`
	MaxVersion uint16 `xml:"MaxVersion,attr" json:"MaxVersion,omitempty"`

	CiphersVer []string `xml:"CiphersVer,attr" json:"CiphersVer,omitempty"`
}

// 主机缓存配置（集群策略或单策略）
type CacheTypeXml struct {
	UrlRegex string `xml:"urlRegex,attr" json:"urlRegex"`
	// 为空择跟原站数据走
	Scheme     string `xml:"scheme,attr" json:"scheme"`
	Suffix     string `xml:"suffix,attr" json:"suffix"`
	TTLTime    string `xml:"ttl,attr" json:"ttl"`
	NeverCache bool   `xml:"neverCache,attr" json:"neverCache"`
	// 是否强制
	Forced bool `xml:"forced,attr" json:"forced"`
}

type MapsXml struct {
	From         string        `xml:"Url,attr" json:"from,omitempty"`
	Target       string        `xml:"Target,attr" json:"target,omitempty"`
	RedirectCode int           `xml:"redirectCode,attr" json:"redirectCode,omitempty"`
	TargetUri    *TargetUriXml `xml:"TargetUri" json:"targeturis,omitempty"`
	ReplaceHost  string        `xml:"ReplaceHost,attr" json:"IsReplaceHost,omitempty"`
	IsClosed     string        `xml:"Closed,attr" json:"IsClosed,omitempty"`
	LockScheme   string        `xml:"LockScheme,attr" json:"IsLockScheme,omitempty"`
	IsAutoCert   string        `xml:"AutoGetCert,attr" json:"IsAutoGetCrt,omitempty"`
	IsUseAts     string        `xml:"httpLockATS,attr" json:"IsUseCache,omitempty"`
	ReplaceResp  string        `xml:"ReplaceResp,attr" json:"IsReplaceHttps,omitempty"`
	ForceCSP     string        `xml:"ForceCSP,attr" json:"IsForceCSP,omitempty"`

	// 开启GZip模式
	GZip    bool   `xml:"GZip,attr" json:"gzip,omitempty"`
	GZipExt string `xml:"GZipExt" json:"gzipExt,omitempty"`

	// ssl自定义
	SSLConf *SSLTypeXml `xml:"SSLConf" json:"SSLConf,omitempty"`

	// 缓存策略组
	Cache []CacheTypeXml `xml:"cache" json:"cache,omitempty"`

	// 检测是哪个Path
	Action    string `xml:"Action,attr" json:"action,omitempty"`
	Exclusion string `xml:"Exclusion,attr" json:"Exclusion,omitempty"`
	Include   string `xml:"Include,attr" json:"Includes,omitempty"`
}

type RemapOptions struct {
	ReplaceHost bool   `xml:"ReplaceHost,attr" json:"IsReplaceHost"`
	LockScheme  bool   `xml:"LockScheme,attr" json:"IsLockScheme"`
	HttpLockATS bool   `xml:"httpLockATS,attr" json:"IsUseCache"`
	ATSHost     string `xml:"ATSHost,attr" json:"atsHost"`
	ReplaceResp bool   `xml:"ReplaceResp,attr" json:"IsReplaceHttps"`
	ForceCSP    bool   `xml:"ForceCSP,attr" json:"IsForceCSP"`
	IsAutoCert  bool   `xml:"AutoGetCert,attr" json:"IsAutoGetCrt"`

	Exclusion string `xml:"Exclusion,attr" json:"Exclusion"`
	Include   string `xml:"Include,attr" json:"Includes"`

	// 全局策略
	// 开启GZip模式
	GZip    bool   `xml:"GZip,attr" json:"gzip"`
	GZipExt string `xml:"GZipExt" json:"gzipExt"`

	// 缓存策略组
	Cache []CacheTypeXml `xml:"cache" json:"cache"`
}

type RemapXml struct {
	XMLName      string       `xml:"Remaps"`
	Options      RemapOptions `xml:"Option"`
	Maps         []MapsXml    `xml:"Map"`
	Regexs       []MapsXml    `xml:"Map_Regex"`
	Rediects     []MapsXml    `xml:"Redirect"`
	RediectRegex []MapsXml    `xml:"Map_RedirectRegex"`
	Rewrites     []MapsXml    `xml:"Map_Rewrite"`
}

const (
	EtcdKey       = "remaps"
	EtcdOptionKey = EtcdKey + "/options"
	EtcdHostKey   = EtcdKey + "/hostv3"
)

func Load4File(conf string) *RemapXml {
	body, err := ioutil.ReadFile(conf)
	if err != nil {
		log.Error("load File Remap Error:", err)
		return nil
	}
	var rmaps *RemapXml = new(RemapXml)
	err = xml.Unmarshal(body, rmaps)
	if err != nil {
		log.Error("load File Remap Error:", err)
		return nil
	}

	return rmaps
}

func Load4Etcd() *RemapXml {
	var rmaps *RemapXml = new(RemapXml)

	err := etcd.UnmarshalKey(EtcdOptionKey, &rmaps.Options)
	if err != nil {
		log.Error("load Etcd Remap Option Error:", err)
		return nil
	}

	gzipExt = Fields(rmaps.Options.GZipExt)
	log.Info("全局GZIP后缀名：", len(gzipExt))

	if len(rmaps.Options.ATSHost) <= 0 {
		rmaps.Options.ATSHost = "127.0.0.1:81"
	}

	var jobs util.SyncJob
	for _, tv := range TypeKey {
		jobs.Append(func(typev string) {
			var mapData []MapsXml
			body, aerr := etcd.ReadAll(EtcdHostKey + "/" + typev + "/")
			if aerr != nil {
				log.Error("load Etcd Remap Host Error:", aerr)
				return
			}

			for _, v := range body.Kvs {
				var mapVal MapsXml
				ce := guard.UnmarshalGzip(string(v.Value), &mapVal)
				if ce != nil {
					log.Error("load Etcd Remap Host Error:", ce)
					continue
				}

				mapData = append(mapData, mapVal)
			}

			switch typev {
			case TypeKey[type_maps]:
				rmaps.Maps = mapData
			case TypeKey[type_regex]:
				rmaps.Regexs = mapData
			case TypeKey[type_redirect]:
				rmaps.Rediects = mapData
			case TypeKey[type_redirect_regex]:
				rmaps.RediectRegex = mapData
			case TypeKey[type_Rewrite]:
				rmaps.Rewrites = mapData
			}
		}, tv)
	}

	jobs.StartWait()

	return rmaps
}

func Write2MapXmlTo(typev string, cmpMaps, maps []MapsXml) {

	// 相同那么不需要
	if reflect.DeepEqual(cmpMaps, maps) {
		return // 相同，然后返回
	}

	// 建立一个临时map
	tempMap := map[string]MapsXml{}
	fromMap := map[string]MapsXml{}
	for _, v := range cmpMaps {
		tempMap[v.From] = v
	}
	// 建立一个来源的临时MAP
	for _, fv := range maps {
		fromMap[fv.From] = fv
	}

	// 循环检测并保存整个集群（保存以及修改）
	for _, v := range maps {

		val, ok := tempMap[v.From]
		if ok {
			if reflect.DeepEqual(val, v) {
				continue
			}
		}

		// 需要增加或需要修改
		etcd.MarshalKey(EtcdHostKey+"/"+typev+"/"+TrimKeyName(v.From), &v) // 保存MAPXML
	}

	// 反向查询并删除
	for _, v := range cmpMaps {
		_, ok := fromMap[v.From]
		if ok {
			continue
		}

		// 删除KEY
		etcd.Delete(EtcdHostKey+"/"+typev+"/"+TrimKeyName(v.From), false) // 删除MAPXML
	}
}

func Write2Etcd(conf *RemapXml) {

	remapXml := Load4Etcd()

	if remapXml == nil {
		remapXml = new(RemapXml)
	}

	// 检测是否需要保存
	if reflect.DeepEqual(remapXml.Options, conf.Options) == false {
		log.Info("Options有差异，开始同步数据...")
		// 分开保存到etcd
		err := etcd.MarshalKey(EtcdOptionKey, &conf.Options) // 单独保存选项
		if err != nil {
			log.Error("write Etcd Remap Error:", err)
		}
	} else {
		log.Info("Option相同无需同步...")
	}

	// 写入ETCD配置每个单独写入
	Write2MapXmlTo(TypeKey[type_maps], remapXml.Maps, conf.Maps)
	Write2MapXmlTo(TypeKey[type_regex], remapXml.Regexs, conf.Regexs)
	Write2MapXmlTo(TypeKey[type_redirect], remapXml.Rediects, conf.Rediects)
	Write2MapXmlTo(TypeKey[type_redirect_regex], remapXml.RediectRegex, conf.RediectRegex)

	log.Info("write etcd remap,success...")
}
