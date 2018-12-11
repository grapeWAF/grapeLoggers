package rules

import (
	"fmt"
	guard "grapeGuard"
	bl "grapeGuard/blacklist"
	"grapeGuard/containers"
	"grapeGuard/rules/jsGuard"
	"net/http"
	"reflect"
	"regexp"
	"strings"
	"time"

	log "grapeLoggers/clientv1"
)

////////////////////////////////////////////////////////////
//  add by koangel
//	date: 2017/09/30
//  rules 防御体系的抽象类
////////////////////////////////////////////////////////////

type GuardInterface interface {
	GuardHttp(w http.ResponseWriter, req *http.Request, remoteIP string) (int, error)
	Name() string
	Open() bool
	Init(conf *GuardRule)
	Response(resp *http.Response)
}

func getConf(chk, def interface{}) interface{} {
	if !reflect.ValueOf(chk).IsNil() {
		return chk
	}

	return def
}

///////////////////////////////////////////////////////////////////////
// 通用数据
type CommonGuard struct {
	isopen     bool
	typename   string
	httpStatus int
}

func (c *CommonGuard) Response(resp *http.Response) {

}

func (c *CommonGuard) GuardHttp(w http.ResponseWriter, req *http.Request, remoteIP string) (int, error) {
	// 初级检测，如果已存在黑名单直接不处理
	return 200, nil
}

func (c *CommonGuard) Name() string {
	return "common"
}

func (c *CommonGuard) Open() bool {
	return c.isopen
}

func (c *CommonGuard) Init(conf *GuardRule) {

}

//////////////////////////////////////////////////////////////
// 单点IP高频率防御使用
type AccessGuard struct {
	CommonGuard
	BlockedTime int
	timeTick    *guard.TimeGroup
	secTick     *guard.TimeGroup
	noFoundTick *guard.TimeGroup
}

func (c *AccessGuard) Response(resp *http.Response) {
	if resp.StatusCode >= 400 {
		remoteIp := bl.GetIP(resp.Request)
		if c.noFoundTick.AddCount(remoteIp) == false {
			// 把ip加入黑名单，并且封停一段时间
			bl.Push2EtchR(resp.Request, "错误请求页次数过多，疑似攻击", int64(c.BlockedTime)) // 加入黑名单
		}
	}
}

func (c *AccessGuard) GuardHttp(w http.ResponseWriter, req *http.Request, remoteIP string) (int, error) {

	if c.isopen == false {
		return 200, nil
	}

	ipData := GetIP4A(remoteIP)
	if ipData == nil {
		return 200, nil
	}

	if ipData.IsSomePath(2) || ipData.IsLinePath(4) || ipData.IsLinePrefix(4) {
		if c.timeTick.AddCount(remoteIP) == false {
			containers.AddGuardCount(req.Host)
			// 把ip加入黑名单，并且封停一段时间
			bl.Push2EtchR(req, "访问频率过高的CC攻击", int64(c.BlockedTime)) // 加入黑名单
			// 超过指定数量，视为攻击
			return c.httpStatus, fmt.Errorf("检测到大量CC请求攻击行为，来自IP:%v", remoteIP)
		}
	}

	if c.secTick.AddCount(remoteIP) == false {
		containers.AddGuardCount(req.Host)
		// 把ip加入黑名单，并且封停一段时间
		bl.Push2EtchR(req, "访问频率过高的CC攻击", int64(c.BlockedTime)) // 加入黑名单
		// 超过指定数量，视为攻击
		return c.httpStatus, fmt.Errorf("检测到大量CC请求攻击行为，来自IP:%v", remoteIP)
	}

	return 200, nil
}

func (c *AccessGuard) Name() string {
	return "singleGuard"
}

func (c *AccessGuard) Init(conf *GuardRule) {
	c.CommonGuard.Init(conf)

	sconf := getConf(conf.Single, defaultConf.Single).(*singleRule)

	// 初始化自己的配置
	c.httpStatus = sconf.ErrorCode
	c.typename = sconf.XMLName
	c.isopen = (sconf.Open == "true") || sconf.Bopen
	if !c.isopen {
		log.Debug("AccessGuard Unopen...")
	}
	c.BlockedTime = sconf.BlockedTime
	c.timeTick = guard.NewTimeGroup(time.Duration(sconf.CheckTime)*time.Second, sconf.Access)
	c.secTick = guard.NewTimeGroup(time.Duration(sconf.CheckTime)*time.Second, sconf.Access+25)
	c.noFoundTick = guard.NewTimeGroup(time.Duration(sconf.NoFoundChk)*time.Second, sconf.NoFoundCount)
}

//////////////////////////////////////////////////////////////
// 代理在短时间的数量
type ProxyGuard struct {
	CommonGuard
	timeTick *guard.TimeGroup
}

func (c *ProxyGuard) IsProxy(req *http.Request) bool {

	_, ok := req.Header["HTTP_X_FORWARDED_FOR"]
	if ok {
		return true
	}

	_, ok = req.Header["HTTP_VIA"]
	if ok {
		return true
	}

	return false
}

func (c *ProxyGuard) GuardHttp(w http.ResponseWriter, req *http.Request, remoteIP string) (int, error) {
	if c.isopen == false {
		return 200, nil
	}

	if c.IsProxy(req) && c.timeTick.AddLimitMap(bl.GetIP(req)) == false {
		containers.AddGuardCount(req.Host)
		// 超过指定数量，视为攻击
		return c.httpStatus, fmt.Errorf("检测到大量代理CC攻击行为...")
	}

	return 200, nil
}

func (c *ProxyGuard) Name() string {
	return "proxyGuard"
}

func (c *ProxyGuard) Init(conf *GuardRule) {
	c.CommonGuard.Init(conf)

	sconf := getConf(conf.Proxy, defaultConf.Proxy).(*proxyRule)
	// 初始化自己的配置
	c.typename = sconf.XMLName
	c.httpStatus = sconf.ErrorCode
	c.isopen = (sconf.Open == "true") || sconf.Bopen
	if !c.isopen {
		log.Debug("ProxyGuard Unopen...")
	}
	c.timeTick = guard.NewTimeGroup(time.Duration(sconf.CheckTime)*time.Second, sconf.Ipcount)
}

//////////////////////////////////////////////////////////////
// UA过滤
type UserAgentGuard struct {
	CommonGuard
	emptyUA      bool
	showTrace    bool
	blocked      []string
	blockedTime  int
	blockedRegex *regexp.Regexp
	timeTick     *guard.TimeGroup
}

func (c *UserAgentGuard) AnitEmptyUA(req *http.Request) bool {
	if c.emptyUA == false {
		return false
	}

	ua := req.UserAgent()
	if len(ua) <= 5 {
		log.Error("User-Agent错误:", ua)
		return true
	}

	return false
}

func (c *UserAgentGuard) AddError(req *http.Request, remoteIp string) {
	if c.timeTick.AddCount(remoteIp) == false {
		// 把ip加入黑名单，并且封停一段时间
		bl.Push2EtchR(req, "USER AGENT异常黑名单", int64(c.blockedTime)) // 加入黑名单
	}
}

func (c *UserAgentGuard) GuardHttp(w http.ResponseWriter, req *http.Request, remoteIP string) (int, error) {
	if c.isopen == false {
		return 200, nil
	}

	if c.AnitEmptyUA(req) {
		c.AddError(req, remoteIP)
		return c.httpStatus, fmt.Errorf("空USER-AGENT数据，过滤:%v", remoteIP)
	}

	ua := req.UserAgent()
	if len(ua) < 5 {
		return 200, nil
	}

	if c.showTrace {
		log.WithField("type", "UATrace").Infof("From:%v,UserAgent:%v,referfer:%v", remoteIP, ua, req.Referer())
	}

	if c.blockedRegex != nil && c.blockedRegex.MatchString(ua) {
		c.AddError(req, remoteIP)
		return c.httpStatus, fmt.Errorf("遇到被阻止的USER-AGENT:%v,关键字:%v", remoteIP, ua)
	}

	// 检测阻止
	//lowUA := strings.ToLower(ua)
	//for _, v := range c.blocked {
	//	if strings.Contains(lowUA, v) {
	//		c.AddError(req, remoteIP)
	//		return c.httpStatus, fmt.Errorf("遇到被阻止的USER-AGENT:%v,关键字:%v", remoteIP, v)
	//	}
	//}

	return 200, nil
}

func (c *UserAgentGuard) Name() string {
	return "UserAgentGuard"
}

func (c *UserAgentGuard) Init(conf *GuardRule) {
	c.CommonGuard.Init(conf)

	sconf := getConf(conf.UA, defaultConf.UA).(*userAgentRule)

	// 初始化自己的配置
	c.typename = sconf.XMLName
	c.httpStatus = sconf.ErrorCode
	c.isopen = (sconf.Open == "true") || sconf.Bopen
	if !c.isopen {
		log.Debug("UserAgentGuard Unopen...")
	}
	c.emptyUA = sconf.AnitEmptyUA
	c.showTrace = sconf.ShowUA
	c.blockedTime = sconf.BlockedTime

	blua := []string{}

	// 有自己的配置，但是也要加上默认配置数据
	if conf.Host != "*" && conf.UA != nil {
		//默认配置的UA先加上
		for _, v := range defaultConf.UA.BlockedUA {
			blua = append(blua, strings.Split(v, "|")...)
		}
	}

	for _, v := range sconf.BlockedUA {
		blua = append(blua, strings.Split(v, "|")...)
	}

	c.blocked = RemoveDuplicatesAndEmpty(blua) // 去重
	regex, err := regexp.Compile(fmt.Sprintf("(?i)(%v)", strings.Join(c.blocked, "|")))
	if err != nil {
		c.blockedRegex = nil
		//
		log.Debugf("编译UserAgent错误:%v...", err)
	}

	c.blockedRegex = regex

	c.timeTick = guard.NewTimeGroup(time.Duration(sconf.CheckTime)*time.Second, sconf.ErrorUACount)
}

//// 此处增加一种针对单个区域请求的防御规则
//// 例如当请求某个指定连接时，在某个时间内不得超过多少次，限定指定的类型，例如png,jpg等。
type CheckUAGuard struct {
	CommonGuard
}

func (c *CheckUAGuard) GuardHttp(w http.ResponseWriter, req *http.Request, remoteIP string) (int, error) {
	if c.isopen == false {
		return 200, nil
	}

	ipData := GetIP4A(remoteIP)
	if ipData == nil {
		return 200, nil
	}

	if ipData.IsBadBrowser() {
		return c.httpStatus, fmt.Errorf("错误的浏览器请求，明确的错误浏览器非真实用户，屏蔽！")
	}

	return 200, nil
}

func (c *CheckUAGuard) Name() string {
	return "CheckUserAgent"
}

func (c *CheckUAGuard) Init(conf *GuardRule) {
	c.CommonGuard.Init(conf)

	sconf := getConf(conf.CBU, defaultConf.CBU).(*checkBadUARule)

	c.isopen = (sconf.Open == "true") || sconf.Bopen
	if !c.isopen {
		log.Debug("CheckUAGuard Unopen...")
	}
	c.httpStatus = sconf.ErrorCode
}

// 代表后台慢速度攻击
type ReqPathGuard struct {
	CommonGuard
	BlockedTime int
	SomeCount   int
	timeTick    *guard.TimeGroup
	Include     []string
}

func (c *ReqPathGuard) GuardHttp(w http.ResponseWriter, req *http.Request, remoteIP string) (int, error) {
	if c.isopen == false {
		return 200, nil
	}

	// 先检测后缀是否合法
	path := strings.ToLower(req.URL.Path)
	for _, ext := range c.Include {
		if !strings.HasSuffix(path, ext) {
			return 200, nil
		}
	}

	// 处理连续请求
	ipData := GetIP4A(remoteIP)
	if ipData == nil {
		return 200, nil
	}

	if ipData.IsSomePath(c.SomeCount) || ipData.IsLinePath(c.SomeCount) || ipData.IsLinePrefix(c.SomeCount) {
		if c.timeTick.AddCount(remoteIP) == false {
			containers.AddGuardCount(req.Host)
			// 把ip加入黑名单，并且封停一段时间
			bl.Push2EtchR(req, "频繁请求相同连接被封停！！！", int64(c.BlockedTime)) // 加入黑名单
			// 超过指定数量，视为攻击
			return c.httpStatus, fmt.Errorf("检测到频繁请求相同地址:%v IP:%v", req.URL.Path, remoteIP)
		}
	}

	return 200, nil
}

func (c *ReqPathGuard) Name() string {
	return "ReqPathGuard"
}

func (c *ReqPathGuard) Init(conf *GuardRule) {
	c.CommonGuard.Init(conf)

	sconf := getConf(conf.ReqGuard, defaultConf.ReqGuard).(*reqGuardRule)

	c.isopen = (sconf.Open == "true") || sconf.Bopen
	if !c.isopen {
		log.Debug("ReqPathGuard Unopen...")
	}

	c.httpStatus = sconf.ErrorCode

	c.Include = strings.Split(sconf.Include, ",")
	c.BlockedTime = sconf.BlockedTime
	c.SomeCount = sconf.SomeCount
	c.timeTick = guard.NewTimeGroup(time.Duration(sconf.CheckTime)*time.Second, sconf.Access)
}

///////// 增加针对首页和特殊参数的剔除
type PickoutGuard struct {
	CommonGuard
	RemoveFP bool
	Ext      []string
}

func (c *PickoutGuard) RemoveQuery(req *http.Request) {

	if len(req.URL.RawQuery) <= 0 {
		return
	}

	req.URL.RawQuery = ""
	pos := strings.Index(req.RequestURI, "?")
	if pos != -1 {
		req.RequestURI = req.RequestURI[0:pos]
	}
}

func (c *PickoutGuard) GuardHttp(w http.ResponseWriter, req *http.Request, remoteIP string) (int, error) {
	if c.isopen == false {
		return 200, nil
	}

	// 需要剔除特殊参数
	if len(req.URL.RawQuery) <= 0 {
		return 200, nil
	}

	// 剔除掉HOME的参数
	if c.RemoveFP && req.URL.Path == "/" {
		c.RemoveQuery(req)
		return 200, nil
	}

	// 剔除制定参数
	lowPath := strings.ToLower(req.URL.Path)
	for _, vext := range c.Ext {
		if strings.HasSuffix(lowPath, vext) {
			c.RemoveQuery(req)
			return 200, nil
		}
	}

	return 200, nil
}

func (c *PickoutGuard) Name() string {
	return "PickoutQuerys"
}

func (c *PickoutGuard) Init(conf *GuardRule) {
	c.CommonGuard.Init(conf)

	sconf := getConf(conf.Pickout, defaultConf.Pickout).(*PickoutQueryRule)

	if sconf == nil {
		c.isopen = false
		return
	}

	c.isopen = (sconf.Open == "true") || sconf.Bopen
	if !c.isopen {
		log.Debug("Pickout Data Unopen...")
	}
	c.httpStatus = sconf.ErrorCode
	c.Ext = strings.Split(sconf.Exts, ",")
	for i, to := range c.Ext {
		c.Ext[i] = strings.ToLower(to)
	}
	c.RemoveFP = sconf.RemoveFP
}

// 人机检测系统

type AVaildGuard struct {
	CommonGuard
	hostMap   containers.SUrlMaps
	autoOpen  bool
	autoLimit int32
	autoTime  time.Duration
	closeTime time.Duration
	needClose bool
	jGuard    *jsGuard.JSGuard
	Blocktime int
}

func (c *AVaildGuard) GuardHttp(w http.ResponseWriter, req *http.Request, remoteIP string) (int, error) {

	_, open := c.hostMap.Lookup(req.Host)
	if c.isopen == false && open == false {
		// 自动开启人机检测
		if c.autoOpen && containers.IsOverflow(req.Host, c.autoLimit, c.autoTime) {
			c.hostMap.Map(req.Host, true)
		}
		return 200, nil
	}

	// 自动关闭人机检测
	if c.needClose && containers.IsTimeout(req.Host, c.autoLimit, c.closeTime) {
		c.hostMap.Delete(req.Host)
		return 200, nil
	}

	// 直接防御
	status, err := c.jGuard.CheckVaild(w, req, remoteIP)
	if err != nil && status != jsGuard.WStatusCode && status != jsGuard.WRedirectCode {
		// 把ip加入黑名单，并且封停一段时间
		bl.Push2EtchR(req, "判断为机器人非法请求!", int64(c.Blocktime)) // 加入黑名单
		// 超过指定数量，视为攻击
		return c.httpStatus, fmt.Errorf("判断为机器人，封IP：%v", remoteIP)
	}

	if status == jsGuard.WStatusCode || status == jsGuard.WRedirectCode {
		return status, fmt.Errorf("已处理的代码")
	}

	return 200, nil
}

func (c *AVaildGuard) Name() string {
	return "AndroidGuard"
}

func (c *AVaildGuard) Init(conf *GuardRule) {
	c.CommonGuard.Init(conf)

	c.jGuard = jsGuard.NewJSGuard()
	sconf := getConf(conf.ACheck, defaultConf.ACheck).(*AndroidCheckRule)
	if sconf == nil {
		c.isopen = guard.IsUseGuard()
		c.Blocktime = 1800
		c.httpStatus = 503
	} else {
		c.isopen = (sconf.Open == "true") || sconf.Bopen
		if !c.isopen {
			log.Debug("AVaildGuard Unopen...")
		}

		c.needClose = !c.isopen
		c.httpStatus = sconf.ErrorCode
		c.jGuard.KeyTTL = time.Duration(sconf.KeyTTL) * time.Second
		c.jGuard.WhiteTTL = time.Duration(sconf.WhiteTTL) * time.Second
		c.jGuard.TryTTL = time.Duration(sconf.TryTTL) * time.Second
		c.jGuard.TryCount = sconf.TryCount
		c.Blocktime = sconf.BlockedTime

		// 自动开启人机检测
		c.autoOpen = sconf.AutoOpen
		c.autoLimit = int32(sconf.AttackCount)
		c.autoTime = time.Duration(sconf.AttackTime) * time.Second
		c.closeTime = time.Duration(sconf.Timeout) * time.Second
	}
}
