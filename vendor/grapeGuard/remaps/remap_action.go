package remaps

import (
	"grapeGuard/rules"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/koangel/grapeNet/Utils"

	log "grapeLoggers/clientv1"

	caddy "grapeGuard/proxy"
)

const (
	replaceHttpKey = `http://`
	targetHttpKey  = `https://`
)

var (
	portRegx, _              = regexp.Compile(`\:[0-9]{2,5}\/`)
	searchTab   []parseToken = []parseToken{
		{"", "script", "src"},
		{"", "iframe", "src"},
		{"", "audio", "src"},
		{"", "img", "src"},
		{"", "video", "src"},
		{"", "font", "src"},
		{"", "style", "src"},
		{"", "link", "href"},
	}

	tagTab []parseToken = []parseToken{
		{"head", "[src]", "src"},
		{"head", "[href]", "href"},
		{"", "[src]", "src"},
		{"", "[href]", "href"},
		{"", "[action]", "action"},
	}
)

type parseToken struct {
	First     string
	ItemToken string
	Attr      string
}

var (
	emptyBody                                                       = ioutil.NopCloser(strings.NewReader(""))
	CollectionCall func(resp *http.Response)                        = func(resp *http.Response) {}
	AtsSyncCall    func(resp *http.Response, options *RemapOptions) = func(resp *http.Response, options *RemapOptions) {}
)

/// 简易版替换代码
func replaceRespSimple(resp *http.Response) error {
	// 不处理任何和html无关的页面
	textType := strings.ToLower(resp.Header.Get("Content-Type"))
	if strings.Contains(textType, "text/html") == false {
		return nil
	}

	// resp没内容不处理
	if resp.ContentLength <= 0 {
		return nil
	}

	if resp.Body == nil {
		return nil
	}

	unBytes, isGzip, zerr := UnGzip(resp)
	if zerr != nil {
		return zerr
	}

	sbody := strings.Replace(string(unBytes), replaceHttpKey, targetHttpKey, -1)
	sbody = portRegx.ReplaceAllString(sbody, "/")

	return Gzip(resp, []byte(sbody), isGzip)
}

func responseData(remap *RemapItem, w http.ResponseWriter, req *http.Request, resp *http.Response) {

	// 防御剔除
	rules.GuardResponse(resp)

	CollectionCall(resp)

	// 指定资源以及指定ATS类型才使用
	if remap.IsUseAts {
		AtsSyncCall(resp, &options)
	}

	//if resp.StatusCode >= 400 {
	//	log.Debug("远端服务器或ATS返回错误，错误代码 - ", resp.StatusCode, resp.Status, resp.Request.URL.String())
	//}

	// smart csp response
	if strings.Contains(strings.ToLower(resp.Header.Get("Content-Type")), "text/html") &&
		IsInclude(req.URL.RequestURI(), remap.Exclusion) == false {

		if (remap.IsReplaceResp || IsInclude(req.URL.RequestURI(), remap.Include)) && req.TLS != nil {

			//log.Info("执行https数据替换，简易版...")
			// 此处无法使用goquery，因为goquery导出的html有严重问题，注意注意！！！
			if err := replaceRespSimple(resp); err != nil {
				log.Error("replace resp error:", err)
			}
		}

		// 注意 只对近代浏览器有效，早期浏览器无效
		if (remap.IsForceCSP) && req.TLS != nil {
			//log.Info("执行强制https数据替换...")
			// 强制https内容升级
			resp.Header.Set("Content-Security-Policy", "upgrade-insecure-requests")
		}
	}

	if remap.IsLockScheme {
		resp.Header.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	}

	// gzip启动 增加Gzip开启压缩算法
	if remap.IsGzip {
		DoGzip(resp, remap.GzipType) // 启动
	}

	// 启动全局的GZIP压缩
	if options.GZip {
		DoGzip(resp, gzipExt)
	}
}

func remap_action(remap *RemapItem, w http.ResponseWriter, req *http.Request) (status int, err error) {
	status = 200
	err = nil

	rsErr := remap.Proxy.ServeHTTP(w, req, func(resp *http.Response) {
		responseData(remap, w, req, resp)
	}) // 反向代理

	if rsErr != nil {
		log.Error("反向代理错误:", rsErr)
		status = 503
		err = rsErr
	}

	return
}

func remap_direct(remap *RemapItem, w http.ResponseWriter, req *http.Request) (status int, err error) {
	status = 200
	err = nil

	redirectCode := remap.RedirectCode
	targetUrl := ""
	if remap.Type == type_redirect {
		if remap.PathRegex == nil {
			ToUrl := &url.URL{
				Scheme:   remap.ToUrl.Scheme,
				Host:     remap.ToUrl.Host,
				Path:     remap.ToUrl.Path,
				RawQuery: req.URL.RawQuery,
			}

			// 优先默认目标的TO参数
			if len(remap.ToUrl.RawQuery) > 0 {
				ToUrl.RawQuery = remap.ToUrl.RawQuery
			}

			targetUrl = ToUrl.String()
		} else {
			ToUrl := &url.URL{
				Scheme:   Utils.Ifs(req.TLS != nil, "https", "http"),
				Host:     req.Host,
				Path:     req.URL.Path,
				RawQuery: req.URL.RawQuery,
			}
			dst := []byte{}
			src := ToUrl.RequestURI()
			match := remap.PathRegex.FindStringSubmatchIndex(src)
			dst = remap.FromRegex.ExpandString(dst, remap.ToUrlSrc, src, match)
			targetUrl = string(dst)
		}

	} else {
		ToUrl := &url.URL{
			Scheme:   Utils.Ifs(req.TLS != nil, "https", "http"),
			Host:     req.Host,
			Path:     req.URL.Path,
			RawQuery: req.URL.RawQuery,
		}
		dst := []byte{}
		targUrls := ToUrl.String()
		match := remap.FromRegex.FindStringSubmatchIndex(targUrls)
		dst = remap.FromRegex.ExpandString(dst, remap.ToUrlSrc, targUrls, match)
		targetUrl = string(dst)
	}

	if redirectCode != 301 && redirectCode != 302 && redirectCode != 307 {
		redirectCode = 302
	}

	http.Redirect(w, req, targetUrl, redirectCode) // 跳转
	status = http.StatusFound
	return
}

func remap_Rewrite(remap *RemapItem, w http.ResponseWriter, req *http.Request) (status int, err error) {
	status = 200
	err = nil
	ToUrl := &url.URL{
		Scheme:   Utils.Ifs(req.TLS != nil, "https", "http"),
		Host:     req.Host,
		Path:     req.URL.Path,
		RawQuery: req.URL.RawQuery,
	}
	dst := []byte{}
	targUrls := ToUrl.String()
	match := remap.FromRegex.FindStringSubmatchIndex(targUrls)
	dst = remap.FromRegex.ExpandString(dst, remap.ToUrlSrc, ToUrl.String(), match)
	turl, _ := url.Parse(string(dst))

	if turl == nil {
		return
	}

	proxy := caddy.NewSingleHostReverseProxy(turl, "", 10, 45*time.Second)
	proxy.UseInsecureTransport()
	rsErr := proxy.ServeHTTP(w, req, func(resp *http.Response) {
		responseData(remap, w, req, resp)
	}) // 反向代理

	if rsErr != nil {
		log.Error("反向代理错误:", rsErr)
		status = 503
		err = rsErr
	}

	return
}
