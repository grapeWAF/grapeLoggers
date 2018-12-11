package rules

import (
	"net/http"
	"regexp"

	log "grapeLoggers/clientv1"
)

////////////////////////////////////////////////////////////
//  add by koangel
//	date: 2017/09/30
//  rules host集
////////////////////////////////////////////////////////////
type GuardHosts struct {
	guardHost  string
	guardHosts []string
	hostRegex  *regexp.Regexp

	guards []GuardInterface
}

func (g *GuardHosts) Response(resp *http.Response) {
	for _, v := range g.guards {
		v.Response(resp)
	}
}

func (g *GuardHosts) MatchHost(req *http.Request) bool {
	if g.guardHost == "*" {
		return true
	}

	return g.hostRegex.MatchString(req.Host)
}

func (g *GuardHosts) Init(host string, conf *GuardRule) {
	g.hostRegex, _ = regexp.Compile(host)
	g.guardHost = conf.Host

	g.buildGuards(conf)
}

func (g *GuardHosts) appendGuard(conf *GuardRule, guardc GuardInterface) {
	guardc.Init(conf)

	log.Debugf("Host:%v，建立防御类型：%v", g.guardHost, guardc.Name())

	g.guards = append(g.guards, guardc)
}

func (g *GuardHosts) buildGuards(conf *GuardRule) {
	g.guards = []GuardInterface{}

	g.appendGuard(conf, &AVaildGuard{})
	g.appendGuard(conf, &AccessGuard{})
	g.appendGuard(conf, &ProxyGuard{})
	g.appendGuard(conf, &UserAgentGuard{})
	g.appendGuard(conf, &CheckUAGuard{})
	g.appendGuard(conf, &ReqPathGuard{})
	g.appendGuard(conf, &PickoutGuard{})
}

func (g *GuardHosts) GuardHttp(w http.ResponseWriter, req *http.Request, remoteIP string) (int, error) {

	for _, v := range g.guards {
		status, err := v.GuardHttp(w, req, remoteIP)
		if err != nil {
			return status, err
		}
	}

	return 200, nil // 成功
}
