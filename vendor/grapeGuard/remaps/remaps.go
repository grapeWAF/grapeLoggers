package remaps

////////////////////////////////////////////////////////////
//  add by koangel
//	date: 2017/09/25
//  反向代理的基础remap库
////////////////////////////////////////////////////////////

import (
	"errors"
	"fmt"
	"grapeGuard/blacklist"
	cta "grapeGuard/containers"
	"net/http"
	"regexp"

	"github.com/koangel/grapeNet/Utils"

	"net"
	"time"

	"context"
	"net/url"
	"strings"

	caddy "grapeGuard/proxy"

	guard "grapeGuard"

	log "grapeLoggers/clientv1"

	"crypto/tls"
)

const (
	type_maps = iota
	type_regex
	type_redirect
	type_redirect_regex
	type_Rewrite // 永久重写
	max_type
)

var (
	TypeKey = []string{
		"remap",
		"regex",
		"redirect",
		"redirect_regex",
		"Rewrite",
	}

	wVRegex     = regexp.MustCompile(`\$\d`)
	wVReplaceFn = func(val string) string {
		return fmt.Sprintf("${%v}", val[1:])
	}
)

type RemapAction struct {
	// 执行特殊行为
	Type   string
	Values []string
}

func (c *RemapAction) Parser(s string) *RemapAction {
	keys := strings.Split(s, ":")
	if len(keys) < 1 {
		c.Type = "empty"
		return c
	}

	c.Type = keys[0]
	if len(keys) > 1 {
		c.Values = keys[1:]
	}

	return c
}

// 转化对应关系 MapXml -> RemapItem
// 通过参数转化和构建RemapItem数据，而RemapItem是运行阶段的实际数据
// 负责映射和部分参数的选择
type RemapItem struct {
	IsActive      bool // 激活状态
	Type          int
	FromSrc       string
	Scheme        string
	ItemKey       string
	RedirectCode  int
	FromUrl       *url.URL
	FromRegex     *regexp.Regexp
	PathRegex     *regexp.Regexp
	PathTemplate  string
	ToScheme      string
	ToUrlSrc      string
	ToUrl         *url.URL
	ToUris        *TargetUriXml
	Proxy         *caddy.ReverseProxy
	ProxyDirector func(*http.Request)
	IsCovertHost  bool
	IsLockScheme  bool
	IsUseAutoCert bool
	IsUseAts      bool
	IsReplaceResp bool
	IsForceCSP    bool // 强制开启HTTPS的CSP验证

	IsNeedUnpack bool // 需要获得BODY

	// 启动GZIP 并且标记类型
	IsGzip   bool
	GzipType []string

	// Cache
	CacheData []CacheTypeXml

	// SSL的参数(如果为空则采用默认的配置)
	TlsConfig *tls.Config

	Actions []*RemapAction

	Exclusion []string
	Include   []string

	RedirectUrl string

	reverseFunc func(remap *RemapItem, w http.ResponseWriter, req *http.Request) (int, error)
}

type RemapContiner = []*RemapItem

var (
	options RemapOptions

	/// 新版本 用MAP去命中常规数据，效率增加
	AtsTarget = "http://127.0.0.1:81"
	// 恢复老版本的REMAPS
	remaps     = cta.NewSMArray(max_type)
	rregexp    = cta.NewRegVecArray(max_type, []int{type_redirect, type_maps})
	remapNames = &cta.SUrlMaps{}

	urlRegex = regexp.MustCompile(`[A-Za-z0-9.*]+`)
	//regexConter *SyncContiner = New()
	//rewrites    *SyncContiner = New()

	//remapContiners []*SyncContiner = NewArray(max_type)
	remapPriority  []int = []int{type_redirect, type_maps}
	remapRegexType []int = []int{type_Rewrite, type_redirect_regex, type_regex}

	gzipExt []string = []string{}
)

func IsRegexType(vt int) bool {
	for _, t := range remapRegexType {
		if t == vt {
			return true
		}
	}

	return false
}

func RemapWatcher(vtype string, key, val []byte) {
	log.Info("Remap更新配置信息:", vtype)
	keystr := strings.TrimPrefix(string(key), EtcdHostKey+"/")

	switch vtype {
	case "PUT":
		log.Info("新改RemapData:", keystr)
		UpdateEtcd(keystr, val)

	case "DELETE":
		log.Info("移除RemapData:", keystr)
		keys, vt := convertType(keystr)

		//
		RemoveHost(vt, strings.Join(strings.Split(keys, "#"), "://"))
	}

	go BeginSyncAts()
	go BeginSyncCache() // 刷新缓存机制
}

func RemapOptWatcher(vtype string, key, val []byte) {
	log.Info("Remap更新选项信息:", vtype)

	if vtype == "PUT" {
		var opt RemapOptions
		err := guard.UnmarshalGzip(string(val), &opt)
		if err != nil {
			log.Error("Option选项更新错误:", err)
			return
		}

		options = opt
		go BeginSyncCache() // 刷新缓存机制
	}
}

func UpdateEtcd(keystr string, val []byte) {
	log.Info("开始建立Remap配置:", keystr)
	_, vtype := convertType(keystr)

	var xml MapsXml
	err := guard.UnmarshalGzip(string(val), &xml)
	if err != nil {
		log.Error("构建失败:", err)
		return
	}

	RemoveHost(vtype, xml.From)
	createOnceRemap(vtype, xml)
	log.Info("建立配置完成...")
}

func GetHttpsRemapItem(host string) *RemapItem {
	if strings.HasPrefix(host, "*.*") {
		host = strings.TrimPrefix(host, "*.")
	}

	for _, vtype := range remapPriority {
		item, ok := remaps[vtype].LookupS("https://", host)
		if ok {
			return item.(*RemapItem)
		}
	}

	return nil
}

func IsAutoGetCert(host string) bool {
	item := GetHttpsRemapItem(host)
	if item == nil {
		return false
	}

	return item.IsUseAutoCert
}

func RemoveHost(vt int, From string) {
	itemKey := GetItemKey(From)

	switch vt {
	case type_redirect_regex, type_regex, type_Rewrite:
		item, vhas := remaps[vt].Lookup(itemKey)
		if !vhas {
			return
		}

		fromKey := item.(*RemapItem).FromSrc
		redirectItem := item.(*RemapItem).RedirectUrl
		if len(redirectItem) > 0 {
			ridx, _, rhas := rregexp[type_redirect_regex].Prefix(redirectItem)
			if rhas {
				log.Info("Regex删除自动构建跳转：", redirectItem)
				rregexp[type_redirect_regex].Remove(ridx)
			}
		}

		idx, _, has := rregexp[vt].Prefix(fromKey)
		if has {
			rregexp[vt].Remove(idx)
		}

	case type_maps, type_redirect:
		if vt == type_maps || vt == type_regex {
			item, ok := remaps[vt].Lookup(itemKey)
			if !ok {
				return
			}

			redirectItem := item.(*RemapItem).RedirectUrl
			if len(redirectItem) > 0 {
				idx, _, has := rregexp[type_redirect_regex].Prefix(redirectItem)
				if has {
					log.Info("Hit删除自动构建跳转：", redirectItem)
					rregexp[type_redirect_regex].Remove(idx)
				}
			}
		}

		remaps[vt].Delete(itemKey)
	}

}

func createReverseProxy(item *RemapItem, mapv *MapsXml) {
	if item.Type == type_redirect || item.Type == type_redirect_regex {
		return
	}

	item.Proxy = caddy.NewSingleHostReverseProxy(item.ToUrl, "", 256, 90*time.Second)
	item.ProxyDirector = item.Proxy.Director
	item.Proxy.Director = func(req *http.Request) {
		item.ProxyDirector(req)

		if strings.HasSuffix(req.Host, ":443") {
			req.Host = strings.TrimSuffix(req.Host, ":443")
		}

		if item.IsCovertHost && options.HttpLockATS == false {
			req.Host = req.URL.Host
		}
	}

	// 先锁定HTTP
	item.Proxy.Transport = &http.Transport{
		MaxIdleConns:          0,
		MaxIdleConnsPerHost:   256,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   15 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		Proxy: http.ProxyFromEnvironment,
		Dial: func(network, addr string) (net.Conn, error) {

			netDialer := &net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				DualStack: true,
			}

			toAddr := options.ATSHost
			if item.IsUseAts == false {
				toAddr = addr
			}

			// 非单个
			if len(toAddr) == 0 || (item.ToUris.Type != Target_Single && item.Type == type_maps) {
				toAddr = options.ATSHost // 锁死81
			}

			log.Debug("dial ats address:", toAddr)

			conn, err := netDialer.DialContext(context.Background(), network, toAddr)
			if err != nil {
				return nil, err
			}

			return conn, nil
		},
	}
}

func createOnceRemap(vtype int, xml MapsXml) *RemapItem {
	item := &RemapItem{
		Type:      vtype,
		FromRegex: nil,
	}

	log.Debug("构建Remap:", xml)

	// 兼容版本 修复一个兼容BUG
	if len(xml.Target) >= 0 && xml.TargetUri == nil {
		item.ToUrl, _ = url.Parse(xml.Target)
		item.ToScheme = item.ToUrl.Scheme
		TargetAddr := RemoveMoreScheme(xml.Target, 1)
		if xml.TargetUri == nil {
			item.ToUris = &TargetUriXml{
				MType:       MainTypeNormal,
				Type:        Target_Single,
				Open:        false,
				Https:       false,
				HealthCheck: false,
				Follow:      false,
				Uris: []SingleUriXml{
					{
						Uri:      TargetAddr,
						Backup:   false,
						Weight:   5,
						TryCount: 1,
						Timeout:  30,
					},
				},
			}
		}
	} else if xml.TargetUri != nil {
		item.ToUris = xml.TargetUri
		if vtype != type_redirect_regex && vtype != type_redirect {
			item.ToUrl = ChkAndFixUrlParser("http", AtsTarget)
		} else {
			scheme := Utils.Ifs(xml.TargetUri.Https, "https://", "http://")
			item.ToUrl, _ = url.Parse(scheme + xml.TargetUri.Uris[0].Uri)
		}

		item.ToScheme = Utils.Ifs(xml.TargetUri.Https, "https", "http")
	}

	if xml.TargetUri == nil {
		log.Error("构建Remap异常(TargetUri = nil)：", xml.From, " 	To:", xml.Target)
		return nil
	}

	item.RedirectCode = xml.RedirectCode
	item.TlsConfig = CreateFromXml(&xml) // 创建CONFIG
	item.CacheData = xml.Cache           // 缓存数据
	item.IsGzip = xml.GZip               // gzip打开
	item.GzipType = Fields(xml.GZipExt)

	FormUrl := RemoveMoreScheme(xml.From, 1)
	item.FromUrl, _ = url.Parse(FormUrl)
	if item.FromUrl != nil {
		item.FromUrl.Host = strings.ToLower(item.FromUrl.Host)
	}
	item.FromSrc = FormUrl
	item.Scheme = GetScheme(FormUrl)
	item.ItemKey = GetItemKey(FormUrl)
	item.RedirectUrl = ""

	remapNames.Map(GetHostOnly(item.ItemKey), item.FromSrc)

	item.Actions = []*RemapAction{}
	if len(xml.Action) > 0 {
		actions := strings.Split(xml.Action, "|")
		for _, av := range actions {
			ac := &RemapAction{}
			item.Actions = append(item.Actions, ac.Parser(av))
		}
	}

	item.Exclusion = []string{}
	if len(xml.Exclusion) > 0 {
		item.Exclusion = strings.Split(xml.Exclusion, "|") // 排除升级规则
	}
	item.Include = []string{}
	if len(xml.Include) > 0 {
		item.Include = strings.Split(xml.Include, "|") // 排除升级规则
	}

	item.IsCovertHost = guard.CmpBool(len(xml.ReplaceHost) > 0, (xml.ReplaceHost == "true"), options.ReplaceHost)
	item.IsLockScheme = guard.CmpBool(len(xml.LockScheme) > 0, (xml.LockScheme == "true"), options.LockScheme)
	item.IsUseAutoCert = guard.CmpBool(len(xml.IsAutoCert) > 0, (xml.IsAutoCert == "true"), options.IsAutoCert)
	item.IsUseAts = guard.CmpBool(len(xml.IsUseAts) > 0, (xml.IsUseAts == "true"), options.HttpLockATS)
	item.IsReplaceResp = guard.CmpBool(len(xml.ReplaceResp) > 0, (xml.ReplaceResp == "true"), options.ReplaceResp)
	item.IsForceCSP = guard.CmpBool(len(xml.ForceCSP) > 0, (xml.ForceCSP == "true"), options.ForceCSP)

	if item.IsReplaceResp || len(item.Actions) > 0 {
		item.IsNeedUnpack = true
	}

	// 构建反向代理模型
	createReverseProxy(item, &xml)

	switch vtype {
	case type_maps, type_redirect:
		if item.Scheme == "http" && vtype == type_maps {
			_, has := remaps[vtype].Lookup(item.ItemKey)
			if has {
				log.Debug("冲突的Remap:%v", item.FromSrc)
				return nil
			}
		}

		if item.Scheme == "https" && item.IsLockScheme && vtype == type_maps {
			// 针对锁定https的行为，自动重建跳转
			redirect := createRedirectRegex(item.FromSrc)
			if redirect != nil {
				item.RedirectUrl = redirect.FromSrc
				_, has := remaps[type_maps].Lookup(item.ItemKey)
				if has {
					remaps[type_maps].Delete(item.ItemKey)
				}

				rregexp[type_redirect_regex].AddRegexS(redirect.FromSrc, redirect)
			}
		}

		item.PathRegex = nil

		if item.FromUrl != nil && vtype == type_redirect {
			reqPath := item.FromUrl.RequestURI()
			if len(reqPath) > 0 {
				pathRgex, err := regexp.Compile(reqPath)
				if err == nil {
					item.PathRegex = pathRgex

					if xml.TargetUri != nil {
						item.ToUrlSrc = ChkAndFixUrl(item.ToScheme, xml.TargetUri.Uris[0].Uri)
						item.ToUrlSrc = wVRegex.ReplaceAllStringFunc(item.ToUrlSrc, wVReplaceFn)
					}
				} else {
					log.Errorf("构建路径失败，但是整体可用:%v", err)
				}
			}
		}

		// 新增重写支持
		switch vtype {
		case type_maps:
			item.reverseFunc = remap_action
			break
		case type_redirect:
			item.reverseFunc = remap_direct
			break
		}

		remaps[vtype].Map(item.ItemKey, item)

	case type_regex, type_redirect_regex, type_Rewrite:

		_, err := VaildUrlRgex(xml.From, true)
		if err != nil {
			log.Errorf("构建Remap失败，不是有效的Url:%v - %v", xml.From, err)
			return nil
		}

		urls, err := url.Parse(xml.From)
		if err != nil {
			log.Errorf("构建Remap失败，不是有效的Url:%v - %v", xml.From, err)
			return nil
		}

		q := urlRegex.FindAllString(urls.Host, -1)
		if len(q) > 1 {
			log.Errorf("构建Remap失败，有限正则支持:%v - %v - %v", xml.From, err, q)
			return nil
		}

		toRegex, terr := regexp.Compile(xml.From)
		if terr != nil {
			log.Error("构建Remap Regex错误:", terr)
			return nil
		}

		/*pos := strings.Index(xml.From, "/")
		if pos != -1 {
			if pos == 0 || (xml.From[pos-1] != '\\') {
				log.Error("构建Remap Regex错误:不兼容的内容 </> ")
				return nil
			}
		}*/

		if xml.TargetUri != nil && item.Type == type_redirect_regex {
			item.ToUrlSrc = ChkAndFixUrl(item.ToScheme, xml.TargetUri.Uris[0].Uri)
			item.ToUrlSrc = wVRegex.ReplaceAllStringFunc(item.ToUrlSrc, wVReplaceFn)
			// 正则数据
		}

		if item.Scheme == "https" && item.IsLockScheme && vtype == type_regex {
			// 针对锁定https的行为，自动重建跳转
			regexD := createRedirectRegex(urls.Scheme + "://" + urls.Host)
			item.RedirectUrl = regexD.FromSrc
			rregexp[type_redirect_regex].AddRegexS(regexD.FromSrc, regexD)
		}

		switch vtype {
		case type_regex:
			item.reverseFunc = remap_action
		case type_Rewrite:
			item.reverseFunc = remap_Rewrite
		case type_redirect_regex:
			item.reverseFunc = remap_direct
		}

		item.FromSrc = xml.From
		item.FromRegex = toRegex
		rregexp[vtype].AddRegexS(item.FromSrc, item)
		remaps[vtype].Map(item.ItemKey, item)
	}

	return item
}

func createRedirectRegex(fromRegex string) *RemapItem {
	toRegex := strings.Replace(fromRegex, "https://", "http://", -1) + "(.*)"
	redirectItem := &RemapItem{
		Type:    type_redirect_regex,
		FromUrl: nil,
	}

	redirectItem.ItemKey = GetItemKey(strings.Replace(fromRegex, "https://", "http://", -1))

	params := "${1}"
	if strings.Contains(fromRegex, "*.") {
		toRegex = strings.Replace(toRegex, "*.", "(.*)", -1)
		params = "${2}"
	}

	redirectItem.RedirectCode = 307
	redirectItem.ToScheme = "https"
	redirectItem.ToUrlSrc = strings.Replace(fromRegex, "*.", "${1}", -1) + params
	redirectItem.FromSrc = toRegex
	redirectItem.FromRegex, _ = regexp.Compile(toRegex)
	redirectItem.PathRegex, _ = regexp.Compile("/(.*)")
	redirectItem.ToUrl = &url.URL{Scheme: "https", Host: "anyto"}
	redirectItem.IsLockScheme = false
	redirectItem.IsCovertHost = false
	redirectItem.IsUseAts = false // 临时数据不要锁定
	redirectItem.reverseFunc = remap_direct

	log.Debug("自动构建跳转Regex成功:", *redirectItem)

	return redirectItem
}

func createRemap(vtype int, vmap []MapsXml) {
	for _, v := range vmap {
		createOnceRemap(vtype, v)
	}
}

func makeMapFromConf(mapData *RemapXml) {
	if mapData == nil {
		log.Error("建立配置信息失败！")
		return
	}

	options = mapData.Options // 得到默认选项
	AtsTarget = options.ATSHost
	if len(AtsTarget) == 0 {
		AtsTarget = "127.0.0.1:81"
	} else if strings.HasPrefix(options.ATSHost, "http://") {
		options.ATSHost = strings.TrimPrefix(options.ATSHost, "http://")
	}
	// 建立临时的map
	createRemap(type_maps, mapData.Maps)
	createRemap(type_regex, mapData.Regexs)
	createRemap(type_redirect, mapData.Rediects)
	createRemap(type_redirect_regex, mapData.RediectRegex)
	createRemap(type_Rewrite, mapData.Rewrites)

	log.Info("构建数量全部Remap完成。")
}

// 建立合理的remap规则体系
func BuildRemapRules() {
	makeMapFromConf(Load4Etcd())

	go BeginSyncAts()
	go BeginSyncCache() // 同步缓存机制
}

func BuildRemapRulesFromFile(filename string) {
	makeMapFromConf(Load4File(filename))
}

func searchMap(host, rawPath, query string, ssl bool) *RemapItem {

	// 构建标准Host
	ToUrl := &url.URL{
		Scheme:   Utils.Ifs(ssl, "https", "http"),
		Host:     host,
		Path:     rawPath,
		RawQuery: query,
	}

	// 先命中remap或redirect
	scheme := Utils.Ifs(ssl, "https://", "http://")

	for _, vt := range remapPriority {
		val, ok := remaps[vt].LookupS(scheme, host)
		if ok {
			remapItem := val.(*RemapItem)
			if remapItem.PathRegex == nil {
				return remapItem
			}

			if remapItem.PathRegex != nil && remapItem.PathRegex.MatchString(ToUrl.RequestURI()) {
				return remapItem
			}
		}
	}

	for _, rvt := range remapRegexType {
		rregexp[rvt].BuildJobs(25)
		_, item, has := rregexp[rvt].LookupP(ToUrl.String())
		if has {
			return item.(*RemapItem)
		}
	}

	return nil
}

func SearchRemap(req *http.Request) *RemapItem {
	return searchMap(req.Host, req.URL.Path, req.URL.RawQuery, req.TLS != nil)
}

func SearchRemapNonReq(Host, Path string, ssl bool) *RemapItem {
	return searchMap(Host, Path, "", ssl)
}

func SearchRemapHost(scheme, host string) (string, bool) {
	value, ok := remapNames.LookupS(scheme, host)
	if ok {
		return value.(string), true
	}
	return "", false
}

// 负责处理反向代理核心模块，线程安全
func ReverseProxyHandler(w http.ResponseWriter, req *http.Request) (status int, err error) {
	status = 200
	err = nil

	if guard.IsDebug() {
		log.Info("连接请求进入:", req.Host, ",Uri:", req.RequestURI)
	}

	// 插入用户真实IP到数据头
	req.Header.Set("X-Forwarded-For", blacklist.GetIP(req)) // 把自己IP加进去

	remap := SearchRemap(req)
	if remap == nil {
		status = 404
		err = errors.New("未知的主机")
		return
	}

	status, err = remap.reverseFunc(remap, w, req)
	return
}
