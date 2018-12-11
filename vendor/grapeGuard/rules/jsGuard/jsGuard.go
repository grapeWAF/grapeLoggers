package jsGuard

import (
	"fmt"
	"html/template"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/koangel/grapeTimer"

	"github.com/koangel/grapeNet/Utils"

	log "grapeLoggers/clientv1"
)

type vaildItem struct {
	Hash     string
	Host     string
	RemoteIP string
	UA       string
	TTL      int64
	Url      *url.URL
}

type ipData struct {
	RemoteIP string
	TryCount int
	TTL      int64
}

const (
	headerKey   = "__gKey"
	cookiesName = "gxToken"
)

var (
	passExts = []string{
		"ico", "css", "jpg", "png", "js", "jpge", "gif", "woff", "ttf", "woff2", "scss", "tga", "txt", "xml", "json",
	}

	SyncCB      = func(key string, TTL int64) error { return nil }
	vaildKeys   sync.Map
	whiteData   sync.Map
	ipVaild     sync.Map
	hostMap     sync.Map
	cookisNames sync.Map

	DJSGuard = NewJSGuard()

	WStatusCode   = 401
	WRedirectCode = 1001

	DTryCount     = 6
	DTryIconCount = 12
	DTryTTL       = 60 * time.Second

	DKeyTTL = 45 * time.Second

	DwhiteTTL = 4 * time.Hour

	guardTpl = template.New("ver2Tpl")
)

type WKeys struct {
	Key  string
	Path string
}

type JSGuard struct {
	TryCount     int
	TryTTL       time.Duration
	KeyTTL       time.Duration
	WhiteTTL     time.Duration
	Level        int
	UseRecaptcha bool
}

func init() {
	Tpl, err := guardTpl.Parse(ver2TplData)
	if err != nil {
		panic(err)
	}
	guardTpl = Tpl
	grapeTimer.CDebugMode = false
	grapeTimer.UseAsyncExec = true
	grapeTimer.SkipWaitTask = true
	grapeTimer.InitGrapeScheduler(time.Second, false)
	grapeTimer.NewTickerLoop(15*1000, -1, procDeleteKey)
}

func DeleteTTL(vmap *sync.Map, fn func(val interface{}) bool) {
	vKeys := []string{}
	vmap.Range(func(key, val interface{}) bool {
		if fn(val) {
			vKeys = append(vKeys, key.(string))
		}
		return true
	})

	for _, kv := range vKeys {
		vmap.Delete(kv)
	}
}

// 定时垃圾回收
func procDeleteKey() {
	now := time.Now().Unix()
	DeleteTTL(&vaildKeys, func(val interface{}) bool {
		return now >= val.(*vaildItem).TTL
	})

	DeleteTTL(&whiteData, func(val interface{}) bool {
		return now >= val.(int64)
	})

	DeleteTTL(&ipVaild, func(val interface{}) bool {
		return now >= val.(*ipData).TTL
	})
}

func NewJSGuard() *JSGuard {
	return &JSGuard{
		KeyTTL:       DKeyTTL,
		TryCount:     DTryCount,
		TryTTL:       DTryTTL,
		WhiteTTL:     DwhiteTTL,
		UseRecaptcha: false,
	}
}

func SyncJGuard(key string, ttl int64) {
	whiteData.Store(key, time.Now().Add(time.Duration(ttl)).Unix())
}

func GetHost(host string) string {

	val, has := hostMap.Load(host)
	if has {
		return val.(string)
	}

	shost := strings.Split(host, ".")
	if len(shost) < 2 {
		return host
	}

	if len(shost) == 2 {
		vhost := strings.Join(append([]string{"*"}, shost...), ".")
		hostMap.Store(host, vhost)
		return vhost
	}

	shost[0] = "*"
	vhost := strings.Join(shost, ".")
	hostMap.Store(host, vhost)
	return vhost
}

func (js *JSGuard) WriteRecaptchaTpl(w http.ResponseWriter, req *http.Request, remoteIP string) (int, error) {
	w.WriteHeader(WStatusCode)
	w.Write([]byte(ver3TplData))

	return WStatusCode, fmt.Errorf("Need Write Html...")
}

func (js *JSGuard) WriteTpl(w http.ResponseWriter, req *http.Request, remoteIP string) (int, error) {

	if js.UseRecaptcha {
		return js.WriteRecaptchaTpl(w, req, remoteIP)
	}

	// 上一次的删除掉
	randNameKey := fmt.Sprintf("%s_%d", cookiesName, rand.Int31n(60000)+1)
	cookies, err := req.Cookie(randNameKey)
	if err == nil {
		cookies.Expires = time.Now().AddDate(-1, 0, 0)
		http.SetCookie(w, cookies)
	}

	// 记录下来真实跳转
	query := RemoveKey(req.URL.RawQuery, headerKey)
	url := &url.URL{
		Scheme:   Utils.Ifs(req.TLS == nil, "http", "https"),
		Host:     req.Host,
		Path:     req.URL.Path,
		RawQuery: query,
	}

	/*if len(url.RawQuery) > 0 {
		url.RawQuery = url.RawQuery + "&" + headerKey + "="
	} else {
		url.RawQuery = headerKey + "="
	}*/

	randKey := byte(rand.Int31n(254) + 1)

	Keys, EnKeys := CreateHash(req.Host, remoteIP, randKey)
	w.WriteHeader(WStatusCode)

	tempHtml := strings.Replace(ver2TplData, "{{.Key}}", EnKeys, -1)
	tempHtml = strings.Replace(tempHtml, "{{.CName}}", randNameKey, -1)
	tempHtml = strings.Replace(tempHtml, "{{.Path}}", url.RequestURI(), -1)
	// {{.xorKey}}
	tempHtml = strings.Replace(tempHtml, "{{.xorKey}}", fmt.Sprint(randKey), -1)

	w.Write([]byte(tempHtml))

	cookisNames.Store(GetHost(req.Host)+remoteIP, randNameKey)
	vaildKeys.Store(string(Keys), &vaildItem{
		Hash:     string(Keys),
		Host:     req.Host,
		RemoteIP: remoteIP,
		UA:       req.UserAgent(),
		TTL:      time.Now().Add(js.KeyTTL).Unix(),
		Url:      url,
	})

	return WStatusCode, fmt.Errorf("Need Write Html...")
}

func (js *JSGuard) IsVaildKey(key, host, addr, ua string) (ret bool) {
	ret = false
	val, ok := vaildKeys.Load(key)
	if !ok {
		return
	}

	item := val.(*vaildItem)
	if time.Now().Unix() >= item.TTL {
		vaildKeys.Delete(key)
		return
	}

	if strings.ToLower(host) != strings.ToLower(item.Host) {
		return
	}

	if addr != item.RemoteIP {
		return
	}

	if ua != item.UA {
		return
	}

	ret = true
	vaildKeys.Delete(key)
	return
}

func (js *JSGuard) CheckKeys(w http.ResponseWriter, req *http.Request, remoteIP string) (ret bool) {
	ret = false

	keys := GetHost(req.Host) + remoteIP
	cName, has := cookisNames.Load(keys)
	if !has {
		return
	}

	cookies, err := req.Cookie(cName.(string))
	if err != nil {
		return
	}

	key, brower, platform := SplitBrower(cookies.Value)
	if js.IsVaildKey(key, req.Host, remoteIP, req.UserAgent()) {
		cookisNames.Delete(keys)
		cookies.Expires = time.Now().AddDate(-1, 0, 0)
		http.SetCookie(w, cookies)
		log.Debugf("验证通过，进入白名单，参数：%s - %s - %s - %s - %s - %s", cookies.Value, remoteIP, req.Host, req.UserAgent(), brower, platform)
		return true
	}

	if len(req.URL.RawQuery) <= 0 {
		return
	}

	token := req.URL.Query().Get(headerKey)
	if len(token) >= 20 && js.IsVaildKey(token, req.Host, remoteIP, req.UserAgent()) {
		log.Error("不在需求KEY，但是提供了KEY，可能是机器人,", token, req.Host, remoteIP)
		return
	}

	return
}

func (js *JSGuard) passSuffix(src string) bool {
	for _, ext := range passExts {
		if strings.HasSuffix(src, ext) {
			return true
		}
	}

	return false
}

func (js *JSGuard) CheckVaild(w http.ResponseWriter, req *http.Request, remoteIP string) (int, error) {
	unix := time.Now().Unix()

	key := GetHost(req.Host) + remoteIP
	ttl, ok := whiteData.Load(key)
	if ok {
		if unix >= ttl.(int64) {
			if js.passSuffix(req.URL.Path) {
				return 200, nil
			}

			whiteData.Delete(key) // whiteData
			return js.WriteTpl(w, req, remoteIP)
		}

		return 200, nil
	}

	if strings.HasSuffix(req.URL.Path, ".ico") {
		ipTry, has := ipVaild.LoadOrStore(remoteIP, &ipData{
			RemoteIP: remoteIP,
			TryCount: 0,
			TTL:      time.Now().Add(js.TryTTL).Unix(),
		})

		if has {
			item := ipTry.(*ipData)
			if item.TTL >= unix {
				if item.TryCount > DTryIconCount {
					return 503, fmt.Errorf("由于重试次数太多，写入黑名单:%v,拉取ICO", remoteIP)
				}
				item.TryCount++
			}
		}

		return 404, nil
	}

	if vaild := js.CheckKeys(w, req, remoteIP); vaild {
		ipVaild.Delete(remoteIP)

		// 放入白名单
		whiteData.Store(key, time.Now().Add(js.WhiteTTL).Unix())
		// 同步
		SyncCB(key, int64(js.WhiteTTL))
		//http.Redirect(w, req, url, 302)

		return 200, nil // 验证通过
	}

	// 超过3次没有成功校验
	ipTry, has := ipVaild.LoadOrStore(remoteIP, &ipData{
		RemoteIP: remoteIP,
		TryCount: 0,
		TTL:      time.Now().Add(js.TryTTL).Unix(),
	})
	if has {
		item := ipTry.(*ipData)
		if item.TTL >= unix {
			if item.TryCount > js.TryCount {
				return 503, fmt.Errorf("由于重试次数太多，写入黑名单:%v", remoteIP)
			}
			item.TryCount++
		}
	}

	return js.WriteTpl(w, req, remoteIP)
}
