package remaps

import (
	"fmt"
	"runtime"
	"strings"
	"time"

	util "github.com/koangel/grapeNet/Utils"

	log "grapeLoggers/clientv1"
)

const (
	domain    = `dest_domain=%v`
	url_regex = `url_regex=%v`
	scheme    = `scheme=%v`

	forceTtl  = `ttl-in-cache=%v`
	normalTtl = `revalidate=%v`

	suffix = `suffix=%v`

	neverCache   = `action=never-cache`
	igonreClient = `action=ignore-client-no-cache`

	portCtrl = `port=%v`

	cacheFile = "/opt/ats/etc/trafficserver/cache.config"
)

var (
	cacheFlag = make(chan int, 36)
)

func init() {
	go procWriteCacheFiles()
}

func procWriteCacheFiles() {
	for {
		select {
		case <-cacheFlag:
			tryCount := 0
			for {
				err := SyncCache() // 同步ats文件
				if err != nil {
					tryCount++
					if tryCount > 5 {
						log.Error("同步Ats Cache错误:", err, "，持续失败，放弃。")
						break
					}
					log.Error("同步Ats Cache错误:", err, "，重试...")
					time.Sleep(time.Second)
				} else {
					break
				}
			}
		}
	}
}

// 把正则表达式的URL拆分和修复
func RegexFixUrl(url string) (scheme, host string) {
	scheme = util.Ifs(strings.HasPrefix(url, "http://"), "http", "https")
	host = strings.TrimPrefix(url, scheme+"://")
	return
}

func ToCacheLine(rtype int, hscheme, host string, xml CacheTypeXml) []string {
	lines := []string{}

	_, realHost := RegexFixUrl(host)
	//port := ""
	if strings.Contains(realHost, ":") {
		pos := strings.Index(realHost, ":")
		realHost = realHost[:pos]
		//port = host[pos+1:]
	}

	if rtype == type_maps {
		if strings.HasPrefix(realHost, "*.") {
			realHost = strings.Replace(realHost, "*.", ".*", -1)
			lines = append(lines, fmt.Sprintf(url_regex, realHost))
		} else {
			if len(xml.UrlRegex) > 0 {
				lines = append(lines, fmt.Sprintf(url_regex, realHost+xml.UrlRegex))
			} else {
				if len(host) > 0 {
					lines = append(lines, fmt.Sprintf(domain, realHost))
				}
			}
		}
	} else {
		lines = append(lines, fmt.Sprintf(url_regex, realHost))
	}

	//if len(port) > 0 {
	//	lines = append(lines,fmt.Sprintf(portCtrl,port))
	//}

	if len(xml.Suffix) > 0 {
		lines = append(lines, fmt.Sprintf(suffix, xml.Suffix)) // 默认强制缓存8小时
	}

	/*if len(xml.Scheme) > 0 {
		lines = append(lines, fmt.Sprintf(scheme, xml.Scheme))
	}else if len(hscheme) > 0{
		lines = append(lines, fmt.Sprintf(scheme, hscheme))
	}*/

	if xml.NeverCache {
		lines = append(lines, neverCache)
	} else {
		if len(xml.TTLTime) > 0 {
			if xml.Forced {
				lines = append(lines, fmt.Sprintf(forceTtl, xml.TTLTime))
			} else {
				lines = append(lines, fmt.Sprintf(normalTtl, xml.TTLTime))
			}
		} else {
			if xml.Forced {
				lines = append(lines, fmt.Sprintf(forceTtl, "8h")) // 默认强制缓存8小时
			}
		}
	}

	return lines
}

func SyncCache() error {

	cacheFiles := cacheFile
	if runtime.GOOS == "windows" {
		cacheFiles = "cache.config"
	}
	traffic := NewConfig(TypeCacheFile, cacheFiles)
	if err := traffic.LoadFile(); err != nil {
		return err
	}

	traffic.ClearVaild() // 删除部分数据

	// 插入全局的设置
	traffic.Append(&LineNode{
		IsComment: false,
		IsEmpty:   false,
		Key:       "",
		Values:    []string{"url_regex=.*", "action=ignore-client-no-cache"},
		Data:      "",
	})

	cacheOnce := func(value interface{}) {
		var item *RemapItem = value.(*RemapItem)

		// 跳转GUARD层实现，不再写入任何配置
		if item.Type == type_redirect || item.Type == type_redirect_regex || item.IsUseAts == false {
			return
		}

		Host := item.FromSrc
		if item.Type == type_maps && item.FromUrl != nil {
			Host = item.FromUrl.Host
		}

		for _, cache := range item.CacheData {
			traffic.Append(&LineNode{
				IsComment: false,
				IsEmpty:   false,
				Key:       "",
				Values:    ToCacheLine(item.Type, item.FromUrl.Scheme, Host, cache),
				Data:      "",
			})
		}
	}

	cacheEcho := func(key, value interface{}) bool {
		cacheOnce(value)
		return true
	}

	// 全局的CACHE设置
	for _, cache := range options.Cache {

		if len(cache.UrlRegex) > 0 {
			traffic.Append(&LineNode{
				IsComment: false,
				IsEmpty:   false,
				Key:       "",
				Values:    ToCacheLine(type_regex, "", cache.UrlRegex, cache),
				Data:      "",
			})
		} else {
			traffic.Append(&LineNode{
				IsComment: false,
				IsEmpty:   false,
				Key:       "",
				Values:    ToCacheLine(type_regex, "", ".*", cache),
				Data:      "",
			})
		}
	}

	remaps[type_maps].SortRange(cacheEcho)
	remaps[type_regex].SortRange(cacheEcho)

	traffic.SaveFile(cacheFiles)

	return traffic.Sync()
}

func BeginSyncCache() {
	log.Info("开始同步Ats Cache文件...")
	cacheFlag <- 1 // 触发写
}
