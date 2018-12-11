package record

import (
	"encoding/json"
	"strings"

	"grapeLoggers/appConf"

	proto "grapeLoggers/protos"

	log "github.com/sirupsen/logrus"
)

////////////////////////////////////////////////////////////
//  add by koangel
//	date: 2017/09/19
//  日志系统，记录体系 基础库
////////////////////////////////////////////////////////////

type LogEntry struct {
	Host       string                 `json:"host"`
	RemoteAddr string                 `json:"remoteIp"`
	Time       string                 `josn:"time"`
	TimeStamp  int64                  `json:"ts"`
	Level      string                 `json:"loglevel"`
	Type       string                 `json:"logType"`
	Msg        string                 `json:"text"`
	Caller     string                 `json:"caller"`
	Data       map[string]interface{} `json:"data"`
}

type LogSearch struct {
	Host      string `json:"host"`
	Type      string `json:"type"`
	Level     string `json:"level"`
	SearchKey string `json:"key"`
	BeginDate string `json:"begin"`
	EndDate   string `json:"end"`
	PageNum   int32  `json:"pn"`
}

func (e *LogEntry) Data2Json() string {
	b, err := json.Marshal(e)
	if err != nil {
		log.Error("data to json:", err)
		return ""
	}

	return string(b)
}

func (e *LogEntry) ExtraJson() string {
	if e.Data == nil {
		return ""
	}

	b, err := json.Marshal(&e.Data)
	if err != nil {
		log.Error("data to json:", err)
		return ""
	}

	return string(b)
}

type LogRecord interface {
	// 记录数据到数据库或某个RECORD
	Record(entrys *LogEntry) bool
	// 类型
	Type() string
	// 初始化这个记录集
	InitRecord() bool
	// 检索数据
	SearchRecord(search *LogSearch) []*LogEntry

	// Host实时信息提交
	SubmitHost(req *proto.HostInfoDataReq)

	//
	SubmitSingle(req *proto.SingleHostDataReq)

	SubmitQps(req *proto.QPSDataReq)

	RemapCommit(req *proto.RemapCommitReq)
}

var records []LogRecord = []LogRecord{}

func InitRecords() bool {
	dbtype := strings.ToLower(config.C.DBtype)

	switch dbtype {
	case "influxdb":
		return Register(&InfluxDBRecord{})
	case "mysql":
		return Register(&MysqlRecord{})
	case "mongodb":
		return Register(&MongoDBRecord{})
	case "pgsql":
		return Register(&PgsqlDBRecord{})
	default:
		log.Error("unsupported db type...")
		return false
	}
}

// 附加一个记录集
func Register(rv LogRecord) bool {
	log.Info("register record:", rv.Type())
	if rv.InitRecord() == false {
		log.Error("register Record:", rv.Type(), ",faild!")
		return false
	}

	log.Info("register Record:", rv.Type(), ",success...")
	records = append(records, rv)

	return true
}

//记录数据到节点
func Record(entry *LogEntry) {
	for _, v := range records {
		v.Record(entry)
	}
}

// 读取
func Search(entry *LogSearch) []*LogEntry {
	for _, v := range records {
		return v.SearchRecord(entry)
	}

	return []*LogEntry{}
}

func SubmitHost(req *proto.HostInfoDataReq) {
	for _, v := range records {
		v.SubmitHost(req)
	}
}

func SubmitSingle(req *proto.SingleHostDataReq) {
	for _, v := range records {
		v.SubmitSingle(req)
	}
}

func SubmitQps(req *proto.QPSDataReq) {
	for _, v := range records {
		v.SubmitQps(req)
	}
}

func RemapCommit(req *proto.RemapCommitReq) {
	for _, v := range records {
		v.RemapCommit(req)
	}
}
