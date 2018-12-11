package actions

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"

	"net/url"

	remap "grapeGuard/remaps"

	"github.com/imroc/req"
	log "grapeLoggers/clientv1"
)

type ClearCacheNode struct {
	SignKey string   `xm:"-" json:"signkey"`
	UrlHost []string `xml:"url_host,cdata" json:"urlhost"`
	ResPath []string `xml:"res_path,cdata" json:"respath"`
}

type clearCacheAction struct {
	doOne sync.Once
	queue chan *ClearCacheNode
}

var (
	ClearAction *clearCacheAction = new(clearCacheAction)
)

func SetupAction() {
	log.Info("安装所有Action...")
	ClearAction.Init()
}

func (c *ClearCacheNode) BuildAllUrls() []*url.URL {
	var urls []*url.URL = []*url.URL{}

	for _, hosts := range c.UrlHost {
		for _, resRaw := range c.ResPath {
			rawUrl := &url.URL{}
			rawUrl.Scheme = "http"
			if strings.HasPrefix(hosts, "http://") {
				rawUrl.Host = strings.TrimPrefix(hosts, "http://")
			} else {
				rawUrl.Host = hosts
			}

			rawUrl.Path = resRaw
			rawUrl.RawPath = resRaw

			urls = append(urls, rawUrl)
		}
	}

	return urls
}

func (c *clearCacheAction) Init() {
	c.doOne.Do(func() {
		c.queue = make(chan *ClearCacheNode, 20000)
		go c.procDoClearCache()
	})
}

func (c *clearCacheAction) Push(action *ClearCacheNode) {
	c.Init()
	c.queue <- action
}

func (c *clearCacheAction) procDoClearCache() {
	for {
		select {
		case msg := <-c.queue:
			{
				urls := msg.BuildAllUrls()
				for _, v := range urls {

					reqUrl := &url.URL{}

					reqUrl.Scheme = v.Scheme
					reqUrl.Host = remap.AtsTarget
					reqUrl.Path = v.Path
					reqUrl.RawQuery = v.RawQuery

					param := req.QueryParam{
						"v": fmt.Sprint(rand.Float32()),
					}

					headers := req.Header{
						"Host": v.Host,
					}

					log.Debug("执行清理缓存任务:", *v)

					_, err := req.Get(reqUrl.String(), headers, param)
					if err != nil {
						log.Errorf("Reflush Ats %v Cache,Error:%v", v.String(), err)
						continue
					}
				}
			}
		}
	}
}
