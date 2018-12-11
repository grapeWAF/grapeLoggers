package record

import (
	"encoding/json"
	"grapeLoggers/appConf"
	"time"

	proto "grapeLoggers/protos"

	log "github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	guard "grapeGuard"
	col "grapeGuard/collection"
)

////////////////////////////////////////////////////////////
//  add by koangel
//	date: 2017/09/24
//  日志系统，Mysql记录数据库行为
////////////////////////////////////////////////////////////

const (
	hCollect   = "hCollect"
	SData      = "singleData"
	QPSCollect = "QPSCollect"
	RemapColl  = "remapColl"
)

type MongoDBRecord struct {
	mog *mgo.Session
}

type mgoEntry = map[string]interface{}

type LogResult struct {
	Level     string `bson:"level"`
	Type      string `bson:"type"`
	LogTime   string `bson:"logTime"`
	TimeStamp int64  `bson:"timeStamp"`
	RemoteIP  string `bson:"remoteIp"`
	Host      string `bson:"hostName"`
	Msg       string `bson:"message"`
	Data      string `bson:"data"`
}

func (r *MongoDBRecord) SubmitQps(req *proto.QPSDataReq) {

	session := r.mog.Copy()
	defer session.Close()
	col := session.DB(config.C.MgoDBName).C(config.C.MgoColName)
	err := col.Insert(&bson.M{
		"qps":       req.Qps,
		"pv":        req.Pv,
		"logTime":   time.Now().Format("2006-01-02 15:04:05"),
		"timestamp": time.Now().Unix(),
	})

	if err != nil {
		log.Error("mongo add log:", err)
		return
	}

	return
}

func (r *MongoDBRecord) SearchRecord(search *LogSearch) []*LogEntry {

	var dm bson.M = bson.M{}
	if len(search.Type) > 0 {
		dm["type"] = search.Type
	}

	if len(search.Level) > 0 {
		dm["level"] = search.Level
	}

	if len(search.Host) > 0 {
		dm["hostName"] = search.Host
	}

	if len(search.BeginDate) > 0 && len(search.EndDate) > 0 {

		beginTime, _ := time.Parse(search.BeginDate, "2006-01-02 15:04:05")
		endTime, _ := time.Parse(search.EndDate, "2006-01-02 15:04:05")

		dm["timeStamp"] = bson.M{
			"$gte": beginTime.Unix(),
			"$lte": endTime.Unix(),
		}
	}

	if len(search.SearchKey) > 0 {
		dm["message"] = "/" + search.SearchKey + "/"
	}

	skip := 0
	if search.PageNum > -1 {
		skip = int(search.PageNum) * 60
	}

	var resultv []LogResult = []LogResult{}
	var respData []*LogEntry = []*LogEntry{}

	session := r.mog.Copy()
	defer session.Close()
	err := session.DB(config.C.MgoDBName).C(config.C.MgoColName).Find(dm).Sort("timeStamp").Skip(skip).Limit(60).All(&resultv)
	if err != nil {
		return []*LogEntry{}
	}

	for _, v := range resultv {
		entry := &LogEntry{
			Host:       v.Host,
			RemoteAddr: v.RemoteIP,
			Time:       v.LogTime,
			TimeStamp:  v.TimeStamp,
			Level:      v.Level,
			Type:       v.Type,
			Msg:        v.Msg,
			Data:       map[string]interface{}{},
		}

		json.Unmarshal([]byte(v.Data), &entry.Data)
		respData = append(respData, entry)
	}

	return respData
}

// 记录数据到数据库或某个RECORD
func (r *MongoDBRecord) Record(entrys *LogEntry) bool {

	session := r.mog.Copy()
	defer session.Close()
	col := session.DB(config.C.MgoDBName).C(config.C.MgoColName)
	err := col.Insert(&bson.M{
		"level":     entrys.Level,
		"type":      entrys.Type,
		"logTime":   entrys.Time,
		"timeStamp": entrys.TimeStamp,
		"remoteIp":  entrys.RemoteAddr,
		"hostName":  entrys.Host,
		"message":   entrys.Msg,
		"caller":    entrys.Caller,
		"data":      entrys.ExtraJson(),
	})

	if err != nil {
		log.Error("mongo add log:", err)
		return false
	}

	return true
}

func (r *MongoDBRecord) SubmitSingle(req *proto.SingleHostDataReq) {

	var docs []interface{}
	var reqData []col.HostCollection = []col.HostCollection{}

	err := guard.UnmarshalJsonGzip(req.JsonBody, &reqData)
	if err != nil {
		log.Error("解析数据错误：", err)
		return
	}

	var missPath []interface{}

	nowUnix := time.Now().Unix()

	for _, v := range reqData {
		docs = append(docs, &bson.M{
			"timestamp": nowUnix,
			"hostName":  req.HostName,
			"uuid":      req.MachineID,

			"hostAddr":  v.HostAddr,
			"RQ":        v.RQTotal,
			"Mobile":    v.Mobile,
			"PCDevice":  v.PCDevice,
			"PadDevice": v.PadDevice,

			// 防御数量
			"GuardCount": v.GuardCount,

			// 命中数据
			"HitCount":  v.HitCount,
			"MissCount": v.MissCount,
			"NoCache":   v.NoCache,
		})

		if len(v.MissRaw) > 0 {
			for _, miss := range v.MissRaw {
				missPath = append(missPath, bson.M{
					"Url":        miss.Url,
					"Path":       miss.Path,
					"Via":        miss.AtsVia,
					"StatusCode": miss.StatusCode,
					"timestamp":  miss.Timestamp,
				})
			}

		}
	}

	session := r.mog.Copy()
	defer session.Close()
	derr := session.DB(config.C.MgoDBName).C(SData).Insert(docs...)
	if derr != nil {
		log.Error("mongo add collect:", derr)
	}

	if len(missPath) == 0 {
		return
	}

	derr = session.DB(config.C.MgoDBName).C("cacheMiss").Insert(missPath...)
	if derr != nil {
		log.Error("mongo add collect:", derr)
	}

}

func (r *MongoDBRecord) SubmitHost(req *proto.HostInfoDataReq) {

	var mreq bson.M = bson.M{}
	jbody, err := json.Marshal(req)
	if err != nil {
		log.Error("mongo add collect:", err)
		return
	}

	json.Unmarshal(jbody, &mreq)

	session := r.mog.Copy()
	defer session.Close()
	err = session.DB(config.C.MgoDBName).C(hCollect).Insert(mreq)
	if err != nil {
		log.Error("mongo add collect:", err)
	}
}

// 类型
func (r *MongoDBRecord) Type() string {
	return "mongodb"
}

func (r *MongoDBRecord) RemapCommit(req *proto.RemapCommitReq) {
	session := r.mog.Copy()
	defer session.Close()
	col := session.DB(config.C.MgoDBName).C(RemapColl)

	var remapData []interface{}
	for _, v := range req.Data {
		remapData = append(remapData, &bson.M{
			"fromSrc":   v.FromSrc,
			"sendBytes": v.SendBytes,
			"recvBytes": v.RecvBytes,
			"reqTotal":  v.Rqtotal,
			"logTime":   time.Now().Format("2006-01-02 15:04:05"),
			"timestamp": v.Time,
		})
	}

	err := col.Insert(remapData...)
	if err != nil {
		log.Error("mongo add log:", err)
		return
	}
}

func (r *MongoDBRecord) ProcRemove() {
	tickTime := time.NewTicker(30 * time.Minute)
	for {
		select {
		case <-tickTime.C:
			chkHour := time.Now().Hour()
			if chkHour >= 5 && chkHour <= 6 {

				// 日志只保留7天在远端
				delTime := time.Now().AddDate(0, 0, -7).Unix()
				qpsDelTime := time.Now().Add(-20 * time.Hour)

				session := r.mog.Copy()
				defer session.Close()

				session.DB(config.C.MgoDBName).C(config.C.MgoColName).Remove(bson.M{
					"timeStamp": bson.M{
						"$lt": delTime,
					},
				})

				// 采集数据只保留7天
				session.DB(config.C.MgoDBName).C(hCollect).Remove(bson.M{
					"timestamp": bson.M{
						"$lt": delTime,
					},
				})

				// 采集数据只保留7天
				session.DB(config.C.MgoDBName).C(SData).Remove(bson.M{
					"timestamp": bson.M{
						"$lt": delTime,
					},
				})

				session.DB(config.C.MgoDBName).C(QPSCollect).Remove(bson.M{
					"timestamp": bson.M{
						"$lt": qpsDelTime,
					},
				})

				session.DB(config.C.MgoDBName).C(RemapColl).Remove(bson.M{
					"timestamp": bson.M{
						"$lt": qpsDelTime,
					},
				})
			}
		}
	}
}

// 初始化这个记录集
func (r *MongoDBRecord) InitRecord() bool {

	session, err := mgo.Dial(config.C.MongoDB)
	if err != nil {
		log.Error("MongoDB Dial Error:", err)
		return false
	}

	session.SetPoolLimit(300)
	r.mog = session

	//r.col = session.DB(config.C.MgoDBName).C(config.C.MgoColName)
	//r.col.EnsureIndexKey("timeStamp")
	//r.col.EnsureIndexKey("hostName")
	//r.col.EnsureIndexKey("type")

	//hcol := session.DB(config.C.MgoDBName).C("hCollect")
	//hcol.EnsureIndexKey("timestamp")
	//hcol.EnsureIndexKey("hostName")

	go r.ProcRemove() //

	return true
}
