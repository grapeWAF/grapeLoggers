package remaps

////////////////////////////////////////////////////////////
//  add by koangel
//	date: 2017/11/02
//  同步remap.conf
////////////////////////////////////////////////////////////

import (
	"fmt"
	guard "grapeGuard"
	"net/url"
	"sort"
	"strings"
	"time"

	log "grapeLoggers/clientv1"
)

const (
	pluginAdd   = `@plugin=balancer.so`
	pluginParam = `@pparam=%v`
	pluginKey   = `4acb08ed296368a2967f12d1812487d9`
)

var (
	hostFairVals = []int{}
)

type addAction = func(rtype, realhost string, target *TargetUriXml) *LineNode

var (
	actions map[int]addAction = map[int]addAction{
		MainTypeNormal:     addInNormal,
		MainTypeHash:       addInHashOrRound,
		MainTypeRoundRobin: addInHashOrRound,
	}

	remapFlag = make(chan int, 36)
)

func init() {
	go procWriteRemap()
}

func procWriteRemap() {
	for {
		select {
		case <-remapFlag:
			tryCount := 0
			for {
				err := SyncAtsMap() // 同步ats文件
				if err != nil {
					tryCount++
					if tryCount > 5 {
						log.Error("同步Ats Remap错误:", err, "，持续失败，放弃。")
						break
					}
					log.Error("同步Ats Remap错误:", err)
					time.Sleep(time.Second)
				} else {
					break
				}
			}
		}
	}
}

func addInNormal(rtype, realhost string, target *TargetUriXml) *LineNode {
	singelHost := strings.TrimPrefix(target.Uris[0].Uri, "http://")
	// 添加新的
	return &LineNode{
		Data:      "",
		Key:       rtype,
		Values:    []string{fixUrl(realhost), fixUrl(guard.CmpStr(target.Https, "https://"+singelHost, "http://"+singelHost))},
		IsEmpty:   false,
		IsComment: false,
	}
}

func addInHashOrRound(rtype, realhost string, target *TargetUriXml) *LineNode {
	ecode, err := cryptParam(target)
	if err != nil {
		log.Error(err)
		return nil
	}

	return &LineNode{
		Data:      "",
		Key:       rtype,
		IsEmpty:   false,
		IsComment: false,
		Values:    []string{fixUrl(realhost), fixUrl(realhost), pluginAdd, fmt.Sprintf(pluginParam, ecode)},
	}
}

func fixUrl(s string) string {
	realUrl, err := url.Parse(s)
	if err != nil {
		return guard.CmpStr(strings.HasSuffix(s, "/"), s, s+"/")
	}

	if realUrl.Path == "" {
		realUrl.Path = "/"
	}

	return realUrl.String()
}

func AtsAddOrUpdate(remap *TrafficFile, retype, host string, target *TargetUriXml) bool {

	if target == nil {
		return false
	}

	if target.Type == Target_Fair {
		pingAddr := []string{}
		fairValues := []int{}
		maxValue := len(target.Uris) * 6
		for i, v := range target.Uris {
			addUrl := v.Uri
			if strings.Contains(v.Uri, ":") {
				spUrl := strings.Split(addUrl, ":")
				addUrl = spUrl[0]
			}
			pingAddr = append(pingAddr, addUrl)
			fairValues = append(fairValues, maxValue-((maxValue/len(target.Uris))*i)-(i*1))
		}

		// ping
		replay := PingArray(pingAddr...)

		sort.SliceStable(replay, func(i, j int) bool {
			return replay[i].Time < replay[j].Time
		})

		// 完成权重数值分配
		for i, rv := range replay {
			for ui, uv := range target.Uris {
				if strings.HasPrefix(uv.Uri, rv.Src) {
					target.Uris[ui].Weight = fairValues[i]
					break
				}
			}
		}
	}

	realHost := guard.CmpStr(strings.HasPrefix(host, "https://"), "http://"+strings.TrimPrefix(host, "https://"), host)

	//log.Debug("Remap目标数据：", retype, ",host:", host, ",data:", *target)
	remap.ReplaceOrAdd(retype, realHost, actions[target.MType](retype, realHost, target))
	return true
}

func SyncAtsMap() error {

	traffic := NewConfig(TypeRemapFile, guard.Conf.RemapConf)

	if err := traffic.LoadFile(); err != nil {
		return err
	}

	traffic.ClearVaild() // 删除非空行

	foreachOnce := func(value interface{}) {
		var item *RemapItem = value.(*RemapItem)

		// 跳转GUARD层实现，不再写入任何配置
		if item.Type == type_redirect || item.Type == type_redirect_regex || item.IsUseAts == false {
			return
		}

		vtype := "map"
		switch item.Type {
		case type_regex, type_Rewrite:
			vtype = "regex_map"
		}

		fromUrl := item.FromSrc
		if item.FromUrl != nil {
			fromUrl = item.FromUrl.String()
			if strings.Contains(item.FromUrl.String(), "*.") && type_maps == item.Type {
				vtype = "regex_map"
				fromUrl = strings.Replace(fromUrl, "*.", "(.*)", -1)
			}
		}

		if item.ToUris == nil {
			log.Error("写出Remap异常项目:", item.FromSrc, ",To:", item.ToUrlSrc)
			return
		}

		AtsAddOrUpdate(traffic, vtype, fromUrl, item.ToUris)
	}

	earchOf := func(key, value interface{}) bool {
		foreachOnce(value)
		return true
	}

	remaps[type_maps].SortRange(earchOf)
	remaps[type_Rewrite].SortRange(earchOf)
	remaps[type_regex].SortRange(earchOf)

	if err := traffic.SaveFile(guard.Conf.RemapConf); err != nil {
		return err
	}

	// 运行重新载入remap
	return traffic.Sync()
}

func BeginSyncAts() {
	if guard.Conf.SyncRemap {
		log.Info("开始同步Ats Remap文件...")
		remapFlag <- 1
	}
}
