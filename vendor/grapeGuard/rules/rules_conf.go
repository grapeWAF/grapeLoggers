package rules

////////////////////////////////////////////////////////////
//  add by koangel
//	date: 2017/09/25
//  rules 规则基础库,配置文件解析
////////////////////////////////////////////////////////////

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"sync"

	guard "grapeGuard"

	log "grapeLoggers/clientv1"

	etcd "github.com/koangel/grapeNet/Etcd"
)

var (
	packupExt []string = []string{}
	oncePack  sync.Once
)

//go:generate easyjson -all rules_conf.go
type commonRule struct {
	Open      string `xml:"open,attr,omitempty" json:"-"`
	Bopen     bool   `xml:"bopen,attr" json:"open"`
	ErrorCode int    `xml:"errorCode,attr,omitempty" json:"errorCode"`
}

func (c *commonRule) IsOpen() bool {
	return (c.Open == "true")
}

func (c *commonRule) Exist() bool {
	return len(c.Open) > 0
}

type singleRule struct {
	XMLName string `xml:"single" json:"-"`
	commonRule
	CheckTime    int `xml:"checkTime,attr,omitempty" json:"checkTime"`
	Access       int `xml:"access,attr,omitempty" json:"access"`
	NoFoundChk   int `xml:"nfcheckTime,attr,omitempty" json:"nfcheckTime"`
	NoFoundCount int `xml:"nfcount,attr,omitempty" json:"nfcount"`
	BlockedTime  int `xml:"blockedTime,attr,omitempty" json:"blockedTime"`
}

type proxyRule struct {
	XMLName string `xml:"proxy" json:"-"`
	commonRule
	CheckTime int `xml:"checkTime,attr,omitempty" json:"checkTime,omitempty"`
	Ipcount   int `xml:"ipcount,attr,omitempty" json:"ipcount,omitempty"`
}

type captchaRule struct {
	XMLName string `xml:"captcha" json:"-"`
	commonRule
	Count     int    `xml:"count,attr,omitempty" json:"count"`
	TickTime  int    `xml:"tickTime,attr,omitempty" json:"tickTime"`
	Cookies   string `xml:"cookies,attr,omitempty" json:"cookies"`
	EcryptKey string `xml:"ecryptKey,attr,omitempty" json:"ecryptKey"`
}

type userAgentRule struct {
	XMLName string `xml:"userAgent" json:"-"`
	commonRule
	CheckTime    int      `xml:"checkTime,attr,omitempty" json:"checkTime"`
	AnitEmptyUA  bool     `xml:"AnitEmptyUA,attr,omitempty" json:"AnitEmptyUA"`
	ErrorUACount int      `xml:"errorUACount,attr,omitempty" json:"-"`
	BlockedTime  int      `xml:"blockedTime,attr,omitempty" json:"-"`
	SingelBUA    string   `xml:"singelBlocked,attr" json:"blocked"`
	BlockedUA    []string `xml:"blocked,omitempty" json:"-"`
	ShowUA       bool     `xml:"ShowUA,attr,omitempty" json:"ShowUA"`
}

// 每个path在同一个时间内的请求次数限制
type reqGuardRule struct {
	commonRule
	CheckTime   int    `xml:"checkTime,attr,omitempty" json:"checkTime"`
	Access      int    `xml:"access,attr,omitempty" json:"access"`
	BlockedTime int    `xml:"blockedTime,attr,omitempty" json:"blockedTime"`
	SomeCount   int    `xml:"someCount,attr,omitempty" json:"someCount"`
	Include     string `xml:"Include,attr,omitempty" json:"Include"`
}

// 是否检测错误异常数据内的UA数据
type checkBadUARule struct {
	commonRule
}

type AndroidCheckRule struct {
	commonRule
	// 是否自动开启JS盾 在OPEN为false下有效
	AutoOpen bool `xml:"autoOpen,attr,omitempty" json:"autoOpen"`
	// 在ATKTIME时间内被攻击ATCKCOUNT次数后，开启JS盾
	AttackTime  int `xml:"timeIn,attr,omitempty" json:"timeIn"`
	AttackCount int `xml:"count,attr,omitempty" json:"count"`
	// 当此时间内未达到攻击标值时，自动关闭JS盾，进入静止状态
	Timeout     int `xml:"timeout,attr,omitempty" json:"timeout"`
	Level       int `xml:"Level,attr,omitempty" json:"Level"`
	KeyTTL      int `xml:"KeyTTL,attr,omitempty" json:"KeyTTL"`
	TryCount    int `xml:"TryCount,attr,omitempty" json:"TryCount"`
	TryTTL      int `xml:"TryTTL,attr,omitempty" json:"TryTTL"`
	WhiteTTL    int `xml:"WhiteTTL,attr,omitempty" json:"WhiteTTL"`
	BlockedTime int `xml:"blockedTime,attr,omitempty" json:"blockedTime"`
}

type PickoutQueryRule struct {
	commonRule
	RemoveFP bool   `xml:"RemoveFP,attr,omitempty" json:"RemoveFP"`
	Exts     string `xml:"Exts,attr,omitempty" json:"Exts"`
}

type GuardRule struct {
	XMLName       string `xml:"guards" json:"xml_name,omitempty"`
	Host          string `xml:"host,attr,omitempty" json:"host"`
	ShowErrorPage bool   `xml:"showErrorPage,attr,omitempty" json:"show_error_page"`

	// 其他规则
	Single   *singleRule       `xml:"single,omitempty" json:"single"`
	Proxy    *proxyRule        `xml:"proxy,omitempty" json:"proxy"`
	Ctp      *captchaRule      `xml:"captcha,omitempty" json:"-"`
	UA       *userAgentRule    `xml:"userAgent,omitempty" json:"userAgent"`
	CBU      *checkBadUARule   `xml:"checkUA,omitempty" json:"checkUA"`
	ReqGuard *reqGuardRule     `xml:"reqGuard,omitempty" json:"reqGuard"`
	Pickout  *PickoutQueryRule `xml:"pickout,omitempty" json:"pickout"`
	ACheck   *AndroidCheckRule `xml:"acheck,omitempty" json:"acheck"`
}

func (r *GuardRule) Empty() {
	r.Host = "*"
	r.ShowErrorPage = true
	r.Single = &singleRule{}
	r.Proxy = &proxyRule{}
	r.Ctp = &captchaRule{}
	r.UA = &userAgentRule{}
	r.CBU = &checkBadUARule{}
	r.ReqGuard = &reqGuardRule{}
}

func (r *GuardRule) SetRule(src *GuardRule) {
	r.Empty()
	if src.Single != nil {
		*r.Single = *src.Single
	}

	if src.Proxy != nil {
		*r.Proxy = *src.Proxy
	}

	if src.Ctp != nil {
		*r.Ctp = *src.Ctp
	}

	if src.UA != nil {
		*r.UA = *src.UA
	}

	if src.CBU != nil {
		*r.CBU = *src.CBU
	}

	if src.ReqGuard != nil {
		*r.ReqGuard = *src.ReqGuard
	}
}

type ExtFilter struct {
	Include string `xml:"Include,attr,omitempty"`
}

type PackoutRule struct {
	Filters []ExtFilter `xml:"Filter,omitempty"`
	Open    bool        `xml:"open,attr,omitempty"`
	Type    string      `xml:"type,attr,omitempty"`
}

func (ex *PackoutRule) String() string {
	strInclude := []string{}

	for _, v := range ex.Filters {
		strInclude = append(strInclude, v.Include)
	}

	return strings.Join(strInclude, "\n")
}

func (ex *PackoutRule) Include(req *http.Request) bool {

	oncePack.Do(func() {
		for _, v := range ex.Filters {
			packupExt = append(packupExt, strings.Split(v.Include, ",")...)
		}

		for i, tl := range packupExt {
			packupExt[i] = strings.ToLower(tl)
		}
	})

	if len(req.URL.RawQuery) <= 0 {
		return false
	}

	if strings.Contains(req.URL.Path, ".") == false {
		return false
	}

	lowPath := strings.ToLower(req.URL.Path)
	for _, ext := range packupExt {
		if strings.HasSuffix(lowPath, ext) {
			return true
		}
	}

	return false

	/*stExt := strings.ToLower(filepath.Ext(req.URL.Path))
	if len(stExt) <= 1 {
		return false
	}

	stExt = stExt[1:] // 取出真的后缀

	for _, v := range ex.Filters {
		if strings.Contains(v.Include, stExt) {
			return true
		}
	}

	return false*/
}

type GuardXml struct {
	XMLName string      `xml:"rules"`
	Rules   []GuardRule `xml:"guards,omitempty"`
	Pickout PackoutRule `xml:"Pickout,omitempty"`
}

const (
	EtcdKey        = "rules"
	EtcdPickoutKey = EtcdKey + "/pickout"
	EtcdRulesKey   = EtcdKey + "/rs"
)

func Load4File(conf string) *GuardXml {
	body, err := ioutil.ReadFile(conf)
	if err != nil {
		log.Error("load File Remap Error:", err)
		return nil
	}
	var rmaps *GuardXml = new(GuardXml)
	err = xml.Unmarshal(body, rmaps)
	if err != nil {
		log.Error("load File Remap Error:", err)
		return nil
	}

	return rmaps
}

func Load4Etcd() *GuardXml {
	var rmaps *GuardXml = new(GuardXml)

	err := etcd.UnmarshalKey(EtcdPickoutKey, &rmaps.Pickout)
	if err != nil {
		log.Error("load Etcd rules Error:", err)
		return nil
	}

	var defRules GuardRule
	err = etcd.UnmarshalKey(EtcdRulesKey+"/*", &defRules)
	if err != nil {
		log.Error("载入默认配置错误，无法解析:", err)
		return nil
	}

	body, aerr := etcd.ReadAll(EtcdRulesKey)
	if aerr != nil {
		log.Error("load Etcd rules Error:", aerr)
		return nil
	}

	for _, v := range body.Kvs {
		var rules GuardRule
		err = guard.UnmarshalGzip(string(v.Value), &rules)
		if err != nil {
			log.Error("load Etcd rules:", err)
			continue
		}

		rmaps.Rules = append(rmaps.Rules, rules)
	}

	return rmaps
}

func Write2File(filename string) {
	guard := Load4Etcd()

	xmlBody, err := xml.Marshal(guard)
	if err != nil {
		return
	}

	ioutil.WriteFile(filename, xmlBody, 666)
}

func Write2MapXmlTo(cmpMaps, maps []GuardRule) {

	// 相同那么不需要
	if reflect.DeepEqual(cmpMaps, maps) {
		return // 相同，然后返回
	}

	// 建立一个临时map
	tempMap := map[string]GuardRule{}
	fromMap := map[string]GuardRule{}
	for _, v := range cmpMaps {
		tempMap[v.Host] = v
	}
	// 建立一个来源的临时MAP
	for _, fv := range maps {
		fromMap[fv.Host] = fv
	}

	// 循环检测并保存整个集群（保存以及修改）
	for _, v := range maps {

		val, ok := tempMap[v.Host]
		if ok {
			if reflect.DeepEqual(val, v) {
				continue
			}
		}

		// 需要增加或需要修改
		etcd.MarshalKey(EtcdRulesKey+"/"+v.Host, &v) // 保存MAPXML
	}

	// 反向查询并删除
	for _, v := range cmpMaps {
		_, ok := fromMap[v.Host]
		if ok {
			continue
		}

		// 删除KEY
		etcd.Delete(EtcdRulesKey+"/"+v.Host, false) // 删除MAPXML
	}
}

func Write2Etcd(conf *GuardXml) {

	cmpdConf := Load4Etcd()
	if cmpdConf == nil {
		cmpdConf = &GuardXml{}
	}

	if reflect.DeepEqual(&cmpdConf.Pickout, &conf.Pickout) == false {
		log.Info("开始同步Pickout规则代码...")
		err := etcd.MarshalKey(EtcdPickoutKey, &conf.Pickout)
		if err != nil {
			log.Error("write Etcd Remap Error:", err)
		}
	}

	// 构建规则表
	log.Info("构建RuleData,规则码...")
	Write2MapXmlTo(cmpdConf.Rules, conf.Rules)

	log.Info("write etcd remap,success...")
}
