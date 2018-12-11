package remaps

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"sync"

	"github.com/avct/uasurfer"
)

type CSPData struct {
	Name       uasurfer.BrowserName
	CSPVersion int
	CSPLevel   int
}

const (
	// 0级 使用替换
	CSP_Level_Zero = iota
	// 1级 阻止低级CSP策略
	CSP_Level_Once
	// 2级 使用高级CSP策略
	CSP_Level_Two
)

var (
	wordRegex = regexp.MustCompile(`^[a-zA-Z0-9.:-]{1,}$`)
	wUnrgex   = regexp.MustCompile(`^[a-zA-Z0-9.*:-]{1,}$`)
	hostcmap  sync.Map

	schemes = []string{"http://", "https://"}

	CSPTable []CSPData = []CSPData{
		{uasurfer.BrowserChrome, 43, CSP_Level_Two},
		{uasurfer.BrowserChrome, 25, CSP_Level_Once},

		{uasurfer.BrowserFirefox, 42, CSP_Level_Two},
		{uasurfer.BrowserFirefox, 33, CSP_Level_Once},

		{uasurfer.BrowserOpera, 30, CSP_Level_Two},
		{uasurfer.BrowserOpera, 25, CSP_Level_Once},

		{uasurfer.BrowserAndroid, 43, CSP_Level_Two},
		{uasurfer.BrowserAndroid, 25, CSP_Level_Once},

		{uasurfer.BrowserSafari, 7, CSP_Level_Once},
		{uasurfer.BrowserIE, 12, CSP_Level_Once},
	}
)

func HashScheme(src string) bool {
	for _, sc := range schemes {
		if strings.HasPrefix(src, sc) {
			return true
		}
	}
	return false
}

func RemoveMoreScheme(src string, keep int) string {
	res := src
	for _, sc := range schemes {
		count := strings.Count(res, sc)
		if count > keep {
			for i := 0; i < count-keep; i++ {
				res = strings.TrimPrefix(res, sc)
			}
		}
	}

	return res
}

func ConvertHost(host string) []string {

	// MAP中存在就不重新创建了
	val, ok := hostcmap.Load(host)
	if ok {
		return val.([]string)
	}

	temp := []string{}
	name := strings.ToLower(host)
	for len(name) > 0 && name[len(name)-1] == '.' {
		name = name[:len(name)-1]
	}

	temp = append(temp, name)

	labels := strings.Split(name, ".")
	for i := range labels {
		labels[i] = "*"
		temp = append(temp, strings.Join(labels, "."))
	}

	hostcmap.Store(host, temp)

	return temp
}

func GetScheme(s string) string {
	pos := strings.Index(s, "://")
	if pos != -1 {
		return s[:pos]
	}
	return "http"
}

func GetHostOnly(s string) string {
	pos := strings.Index(s, "://")
	if pos != -1 {
		end := s[pos+3:]
		// 搜索下一个
		pos = strings.Index(end, "/")
		if pos != -1 {
			return strings.ToLower(end[:pos])
		} else {
			return strings.ToLower(end)
		}
	}

	return strings.ToLower(s)
}

func GetItemKey(s string) string {
	pos := strings.Index(s, "://")
	if pos != -1 {
		scheme := s[:pos+3]
		end := s[pos+3:]
		// 搜索下一个
		pos = strings.Index(end, "/")
		if pos != -1 {
			return strings.ToLower(scheme + end[:pos])
		}
	}

	return strings.ToLower(s)
}

func GetHosts(scheme, host string, keeps int) []string {
	hosts := []string{scheme + host}
	labels := strings.Split(host, ".")

	if keeps <= 0 {
		keeps = len(labels)
	}

	if len(labels) <= 2 {
		labels = append([]string{"*"}, labels...)
	}

	for i := range labels {
		labels[i] = "*"
		if (i + keeps) >= len(labels) {
			break
		}

		hosts = append(hosts, scheme+strings.Join(labels[i:], "."))
	}

	return hosts
}

func convertType(etcdKey string) (searchKey string, vt int) {
	searchKey = etcdKey
	vt = type_maps
	for i := range TypeKey {
		vt = i
		if strings.HasPrefix(searchKey, TypeKey[i]+"/") {
			searchKey = strings.TrimPrefix(searchKey, TypeKey[i]+"/")
			break
		}
	}

	return
}

func TrimKeyName(s string) string {
	scheme := GetScheme(s)
	keyName := s
	if strings.HasPrefix(keyName, "https://") {
		keyName = strings.TrimPrefix(keyName, "https://")
	} else if strings.HasPrefix(keyName, "http://") {
		keyName = strings.TrimPrefix(keyName, "http://")
	}

	return scheme + "#" + keyName
}

func IsGzip(resp *http.Response) bool {
	return resp.Header.Get("Content-Encoding") == "gzip"
}

func UnGzip(resp *http.Response) (body []byte, isGzip bool, err error) {
	body = nil
	err = nil
	isGzip = false

	var buf bytes.Buffer
	var reader io.ReadCloser = resp.Body
	if IsGzip(resp) {
		isGzip = true
		gzipReader, gerr := gzip.NewReader(resp.Body)
		if gerr != nil {
			err = gerr
			return
		}

		reader = gzipReader
	}

	if isGzip {
		defer reader.Close()
	}

	if _, berr := buf.ReadFrom(reader); berr != nil {
		err = berr
		return
	}

	body = buf.Bytes()

	return
}

func HasExt(resp *http.Response, ext []string) bool {
	if len(ext) == 0 {
		return false
	}

	if len(ext) == 1 && ext[0] == "*" {
		return true
	}

	for _, ev := range ext {
		if strings.HasSuffix(resp.Request.URL.RawPath, ev) {
			return true
		}
	}

	return false
}

func DoGzip(resp *http.Response, ext []string) {

	if HasExt(resp, ext) == false {
		return // 不必理会
	}

	if IsGzip(resp) {
		return //已GZIP不必理会
	}

	if resp.ContentLength <= 0 {
		return // 长度不足 不必理会
	}

	uncorData, _ := ioutil.ReadAll(resp.Body)
	err := Gzip(resp, uncorData, true)
	if err != nil {
		return // 错误 不必理会
	}

	resp.Header.Set("Content-Encoding", "gzip") // 开启GZIP压缩
}

func Gzip(resp *http.Response, body []byte, isGzip bool) error {

	writeByte := body
	if isGzip {
		var b bytes.Buffer
		w, werr := gzip.NewWriterLevel(&b, gzip.BestSpeed)
		if werr != nil {
			return werr
		}
		_, err := w.Write(body)
		if err != nil {
			return err
		}
		w.Close()

		writeByte = b.Bytes()
	}

	resp.Body = ioutil.NopCloser(bytes.NewReader(writeByte))
	resp.ContentLength = int64(len(writeByte))
	resp.Header.Set("Content-Length", fmt.Sprint(len(writeByte)))

	return nil
}

func GetCSPLevel(sua string) int {
	UA := uasurfer.Parse(sua)
	for _, v := range CSPTable {
		if v.Name == UA.Browser.Name && UA.Browser.Version.Major >= v.CSPVersion {
			return v.CSPLevel
		}
	}

	return CSP_Level_Zero // 0级代表必须替换
}

func IsInclude(s string, src []string) bool {
	sl := strings.ToLower(s)

	for _, vu := range src {
		if strings.Contains(sl, strings.ToLower(vu)) {
			return true
		}
	}

	return false
}

func Fields(s string) []string {
	res := []string{}

	ext := strings.Fields(s)
	for _, vs := range ext {
		res = append(res, strings.Split(strings.TrimSpace(vs), ",")...)
	}

	return res
}

func ChkAndFixUrl(scheme, src string) string {
	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		return src
	}
	return scheme + "://" + src
}

func ChkAndFixUrlParser(scheme, src string) *url.URL {
	vurl, _ := url.Parse(ChkAndFixUrl(scheme, src))
	return vurl
}

func SplitHost(host string) []string {
	temp := host
	ap := []string{}

	for i := 0; i < 2; i++ {
		// 找出后面几个
		pos := strings.LastIndex(temp, ".")
		if pos == -1 {
			break
		}

		ap = append(ap, temp[pos+1:])
		temp = temp[:pos]
	}

	ap = append(ap, temp)

	for i := len(ap)/2 - 1; i >= 0; i-- {
		opp := len(ap) - 1 - i
		ap[i], ap[opp] = ap[opp], ap[i]
	}

	return ap
}

func VaildUrlRgex(src string, checkRegex bool) (ok bool, err error) {
	ok = false
	err = nil

	purl, perr := url.Parse(src)
	if perr != nil {
		err = perr
		return
	}

	tokens := SplitHost(purl.Host)
	if len(tokens) < 2 {
		err = fmt.Errorf("验证域名错误，请最少输入2个有效的数据，例如abcd.com")
		return
	}

	if len(tokens) > 2 && checkRegex {
		backhost := strings.Join(tokens[1:], ".")
		if wordRegex.MatchString(backhost) == false {
			err = fmt.Errorf("不匹配的域名类型，域名的关键区域不可使用正则，%v", backhost)
			return
		}

		_, rerr := regexp.Compile(tokens[0])
		if rerr != nil {
			err = fmt.Errorf("域名首编译错误：%v", rerr)
			return
		}
	} else {
		if wUnrgex.MatchString(purl.Host) == false {
			err = fmt.Errorf("不匹配的域名类型，域名区域不可使用正则，%v", purl.Host)
			return
		}
	}

	_, rerr := regexp.Compile(purl.RequestURI())
	if rerr != nil {
		err = fmt.Errorf("路径参数编译错误：%v", rerr)
		return
	}

	ok = true
	return
}
