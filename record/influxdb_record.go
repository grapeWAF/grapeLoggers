package record

////////////////////////////////////////////////////////////
//  add by koangel
//	date: 2017/09/20
//  日志系统，InfluxDB的记录
////////////////////////////////////////////////////////////

import (
	"grapeLoggers/appConf"
	"time"

	proto "grapeLoggers/protos"

	"github.com/influxdata/influxdb/client/v2"
	log "github.com/sirupsen/logrus"
)

type recoredType = map[string]interface{}
type recoredTag = map[string]string

type InfluxDBRecord struct {
	influxC client.Client
}

func (r *InfluxDBRecord) RemapCommit(req *proto.RemapCommitReq) {

}

func (r *InfluxDBRecord) SubmitHost(req *proto.HostInfoDataReq) {
	// 让时序数据库支持 两个数据提交
}

func (r *InfluxDBRecord) SubmitQps(req *proto.QPSDataReq) {

}

func (r *InfluxDBRecord) SubmitSingle(req *proto.SingleHostDataReq) {
	// 让时序数据库支持 两个数据提交
}

func (r *InfluxDBRecord) SearchRecord(search *LogSearch) []*LogEntry {
	return []*LogEntry{}
}

// 记录数据到数据库或某个RECORD
func (r *InfluxDBRecord) Record(entrys *LogEntry) bool {

	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  config.C.InfluxName,
		Precision: "s",
	})

	if err != nil {
		log.Error(err)
		return false
	}

	//logTime,_ := time.ParseInLocation("2006-01-02 15:04:05",entrys.Time,config.Loc)

	timeTP := time.Unix(entrys.TimeStamp, 0)

	pb, perr := client.NewPoint("records",
		recoredTag{
			"tagLevel": entrys.Level,
			"tagType":  entrys.Type,
			"tagIP":    entrys.RemoteAddr,
			"tagHost":  entrys.Host,
		},
		recoredType{
			"level":    entrys.Level,
			"type":     entrys.Type,
			"logTime":  entrys.Time,
			"remoteIp": entrys.RemoteAddr,
			"hostName": entrys.Host,
			"message":  entrys.Msg,
			"Data":     entrys.Data,
		}, timeTP)

	if perr != nil {
		log.Error(perr)
		return false
	}

	bp.AddPoint(pb)

	r.influxC.Write(bp)
	return true
}

// 类型
func (r *InfluxDBRecord) Type() string {
	return "InfluxDB"
}

// 初始化这个记录集
func (r *InfluxDBRecord) InitRecord() bool {

	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     config.C.InfluxDB,
		Username: config.C.InfluxUser,
		Password: config.C.InfluxPassword,
	})

	if err != nil {
		log.Error(err)
		return false
	}

	r.influxC = c

	return true
}
