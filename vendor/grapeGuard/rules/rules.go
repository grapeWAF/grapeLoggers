package rules

import (
	"net/http"
	"strings"

	cta "grapeGuard/containers"

	guard "grapeGuard"

	log "grapeLoggers/clientv1"
)

////////////////////////////////////////////////////////////
//  add by koangel
//	date: 2017/09/25
//  rules 规则基础库
////////////////////////////////////////////////////////////

var (
	rulesConf *GuardXml = nil

	ruleMaps                  = &cta.SUrlMaps{}
	rulePickout  *PackoutRule = nil
	defaultConf  *GuardRule   = nil
	defaultRules *GuardHosts  = nil
)

func RuleWatcher(vtype string, key, val []byte) {
	log.Info("防御规则库更新配置信息:", vtype)

	switch vtype {
	case "PUT":
		var rule GuardRule
		err := guard.UnmarshalGzip(string(val), &rule)
		if err != nil {
			log.Error("rules err:", err)
			return
		}

		AppendRules(&rule, true)

	case "DELETE":

		prekey := strings.TrimPrefix(string(key), EtcdRulesKey+"/")

		if prekey == "*" {
			log.Error("默认规则不可删除，删除将出现永久异常!!!")
			return
		}

		log.Info("删除规则:", prekey)
		spHost := strings.Split(prekey, ",")
		for _, v := range spHost {
			if val, ok := ruleMaps.Lookup(v); ok {
				if val.(*GuardHosts).guardHost == prekey {
					log.Info("移除规则针对HOST:", v)
					ruleMaps.Delete(v) // 删除
				}
			}
		}
	}
}

func RulePickoutWatcher(vtype string, key, val []byte) {
	log.Info("同步Pickout数据...")
	if vtype == "PUT" {
		var pv *PackoutRule = &PackoutRule{}
		err := guard.UnmarshalGzip(string(val), pv)
		if err != nil {
			log.Error("pickout err:", err)
			return
		}

		rulePickout = pv
	}
}

func UpdateEtcd() {
	log.Info("开始解析防御规则...")
	BuildEtcdData()
	log.Info("解析防御规则完成！")
}

func AppendRules(rule *GuardRule, isUpdate bool) {

	// 根据默认配置构建真正的配置
	hSHsp := strings.Split(rule.Host, ",")
	log.Info("建立防御规则 于Host:", guard.CmpStr(rule.Host == "*", "any", rule.Host))
	vHosts := &GuardHosts{}
	vHosts.Init(rule.Host, rule)

	if rule.Host == "*" {
		defaultRules = vHosts
		if isUpdate == false {
			return
		}

		log.Info("由于更新默认配置，重建全部防御规则！")

		// 构建所有的非默认防御
		for _, v := range rulesConf.Rules {
			if v.Host == "*" {
				continue
			}

			AppendRules(&v, true)
		}
		return
	}

	// 建立SYNC
	for _, ul := range hSHsp {
		ruleMaps.Map(ul, vHosts)
	}
}

func BuildEtcdData() {
	rulesConf = Load4Etcd()
	if rulesConf == nil {
		log.Error("载入防御规则异常...")
		return
	}

	rulePickout = &rulesConf.Pickout
	// 先把默认配置拿出来
	for _, v := range rulesConf.Rules {
		if v.Host == "*" {
			defaultConf = &v
			break
		}
	}

	// 建立底层的配置信息
	for _, v := range rulesConf.Rules {
		AppendRules(&v, false)
	}
}

func GuardResponse(resp *http.Response) {
	if _, ok := resp.Header["X-Server-By"]; ok {
		resp.Header["X-Server-By"] = []string{"XGuard"}
	}

	if _, ok := resp.Header["Server"]; ok {
		resp.Header["Server"] = []string{"XGurd/" + guard.Version + " GuardOS(Minix)"}
	}

	if _, ok := resp.Header["X-Powered-By"]; ok {
		resp.Header["X-Powered-By"] = []string{"XGuardSystem/" + guard.Version}
	}

	// 隐藏.NET版本，防止针对性攻击
	if _, ok := resp.Header["x-aspnet-version"]; ok {
		resp.Header["x-aspnet-version"] = []string{".net core/" + guard.Version}
	}

	// 检测404次数
	guards := searchRules(resp.Request)
	if guards == nil {
		return
	}

	guards.Response(resp)
}

func searchRules(req *http.Request) *GuardHosts {
	guard, ok := ruleMaps.LookupS("", req.Host)
	if ok {
		return guard.(*GuardHosts)
	}

	return defaultRules
}

func Packout(req *http.Request) {
	if rulesConf.Pickout.Open {
		if rulesConf.Pickout.Include(req) {
			// 剔除所有参数
			req.URL.RawQuery = ""
			pos := strings.Index(req.RequestURI, "?")
			if pos != -1 {
				req.RequestURI = req.RequestURI[0:pos]
			}
		}
	}
}

func GuardLayer(w http.ResponseWriter, req *http.Request, remoteIP string) (int, error) {

	Packout(req) // 存在则直接处理,防止回源攻击

	guards := searchRules(req)
	if guards == nil {
		log.Error("无法查找到防御规则:", req.Host, ",远程IP:", remoteIP)
		return 200, nil // 成功
	}
	return guards.GuardHttp(w, req, remoteIP)
}
