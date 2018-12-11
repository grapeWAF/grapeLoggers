package record

////////////////////////////////////////////////////////////
//  add by koangel
//	date: 2018/04/08
//  日志系统，PGSQL支持
////////////////////////////////////////////////////////////

import (
	"grapeLoggers/appConf"
	proto "grapeLoggers/protos"
	"strings"

	"github.com/koangel/grapeTimer"

	"github.com/go-xorm/xorm"

	log "github.com/sirupsen/logrus"

	"grapeLoggers/record/models"
	"time"

	_ "github.com/lib/pq"

	col "grapeGuard/collection"

	guard "grapeGuard"

	"encoding/json"
)

type PgsqlDBRecord struct {
	orm *xorm.Engine
}

func (r *PgsqlDBRecord) SubmitHost(req *proto.HostInfoDataReq) {

	addHost := &model.TbHcollect{
		Idx:      0,
		System:   req.System,
		Byterecv: int(req.ByteRecv),
		Bytesend: int(req.ByteSent),

		Procmemper: req.ProcMemPercent,
		Proccpuper: req.ProcCpuPercent,

		Cputotal: req.CpuTotal,
		Cpuper:   req.CpuPercent[0],
		Cputype:  req.CpuType,

		Hdfree:  int64(req.HdFree),
		Hdtotal: int64(req.HdTotal),
		Hdper:   req.HdPercent,

		Hostname: req.HostName,
		Hostuuid: req.HostUID,
		Remote:   req.RemoteAddr,

		Sysmemfree: int64(req.SysMemFree),
		Sysmemper:  req.SysMemPercent,
		Sysmemused: int64(req.SysMemUsed),

		Totalbyterecv: int64(req.TotalbyteRecv),
		Totalbytesend: int64(req.TotalbyteSent),
		Totalrq:       int(req.TotalRQ),

		Timestamp: req.Timestamp,
	}

	_, err := r.orm.InsertOne(addHost)
	if err != nil {
		log.Error("insert host error:", err)
		return
	}

}

func (r *PgsqlDBRecord) SubmitQps(req *proto.QPSDataReq) {
	addQps := model.TbQps{
		Qpsidx:    0,
		Rqcount:   int(req.Qps),
		Pv:        int(req.Pv),
		Time:      time.Now(),
		Timestamp: time.Now().Unix(),
	}

	_, err := r.orm.InsertOne(&addQps)
	if err != nil {
		log.Error("insert qps error:", err)
		return
	}
}

func (r *PgsqlDBRecord) SubmitSingle(req *proto.SingleHostDataReq) {
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
		docs = append(docs, &model.TbSingle{
			Idx:       0,
			Timestamp: nowUnix,
			Hostname:  req.HostName,
			Uuid:      req.MachineID,

			Hostaddr:  v.HostAddr,
			Rq:        int(v.RQTotal),
			Mobile:    int(v.Mobile),
			Pcdevice:  int(v.PCDevice),
			Paddevice: int(v.PadDevice),

			Guardcount: int(v.GuardCount),

			Hitcount:  int(v.HitCount),
			Misscount: int(v.MissCount),
			Nocache:   int(v.NoCache),
		})

		if len(v.MissRaw) > 0 {
			for _, miss := range v.MissRaw {
				missPath = append(missPath,
					&model.TbCachemiss{
						Idx:        0,
						Host:       req.HostName,
						Uuid:       req.MachineID,
						Url:        miss.Url,
						Path:       miss.Path,
						Via:        miss.AtsVia,
						Statuscode: miss.StatusCode,
						Timestamp:  miss.Timestamp,
					})
			}

		}
	}

	if len(missPath) > 0 {
		_, err = r.orm.Insert(missPath...)
		if err != nil {
			log.Error("insert missPath error:", err)
			return
		}
	}

	_, err = r.orm.Insert(docs...)
	if err != nil {
		log.Error("insert single error:", err)
		return
	}
}

func (r *PgsqlDBRecord) SearchRecord(search *LogSearch) []*LogEntry {
	var whereData []string = []string{}
	var argData []interface{} = []interface{}{}

	if len(search.Type) > 0 {
		whereData = append(whereData, "type = ?")
		argData = append(argData, search.Type)
	}

	if len(search.Level) > 0 {
		whereData = append(whereData, "level = ?")
		argData = append(argData, search.Level)
	}

	if len(search.Host) > 0 {
		whereData = append(whereData, "host = ?")
		argData = append(argData, search.Host)
	}

	if len(search.BeginDate) > 0 && len(search.EndDate) > 0 {
		whereData = append(whereData, "time >= ? and time <= ?")
		argData = append(argData, search.BeginDate)
		argData = append(argData, search.EndDate)
	}

	if len(search.SearchKey) > 0 {
		whereData = append(whereData, "time like ?")
		argData = append(argData, "%"+search.SearchKey+"%")
	}

	var data []model.TbLogs = []model.TbLogs{}
	whereVal := strings.Join(whereData, " and ")
	err := r.orm.Where(whereVal, argData).Desc("time").Find(&data)
	if err != nil {
		log.Error("search Log:", err)
		return []*LogEntry{}
	}

	var respData []*LogEntry = []*LogEntry{}
	for _, v := range data {
		entry := &LogEntry{
			Host:       v.Host,
			RemoteAddr: v.Serverip,
			Time:       v.Time.String(),
			TimeStamp:  v.Timestmap,
			Level:      v.Level,
			Type:       v.Type,
			Msg:        v.Message,
			Data:       map[string]interface{}{},
		}

		json.Unmarshal([]byte(v.Datajson), &entry.Data)
		respData = append(respData, entry)
	}

	return respData
}

// 记录数据到数据库或某个RECORD
func (r *PgsqlDBRecord) Record(entrys *LogEntry) bool {

	timesp := time.Unix(entrys.TimeStamp, 0)

	addLog := &model.TbLogs{
		Logidx:    0,
		Message:   entrys.Msg,
		Time:      timesp,
		Timestmap: entrys.TimeStamp,
		Level:     entrys.Level,
		Type:      entrys.Type,
		Host:      entrys.Host,
		Serverip:  entrys.RemoteAddr,
		Caller:    entrys.Caller,
		Datajson:  entrys.ExtraJson(),
	}

	_, err := r.orm.InsertOne(addLog)
	if err != nil {
		log.Error("insert log error:", err)
		return false
	}

	return true
}

func (r *PgsqlDBRecord) RemapCommit(req *proto.RemapCommitReq) {

	var remapData []interface{} = []interface{}{}

	for _, v := range req.Data {
		addRemap := &model.TbRemapdata{
			Formsrc:   v.FromSrc,
			Recvbytes: int32(v.RecvBytes),
			Sendbytes: int32(v.SendBytes),
			Reqtotal:  v.Rqtotal,
			Timestamp: v.Time,
			Time:      time.Now(),
		}

		remapData = append(remapData, addRemap)
	}

	_, err := r.orm.Insert(remapData...)
	if err != nil {
		log.Error("insert log error:", err)
		return
	}
}

// 类型
func (r *PgsqlDBRecord) Type() string {
	return "Postgres"
}

func (r *PgsqlDBRecord) procTimerRemove() {
	chkHour := time.Now().Hour()
	if chkHour >= 5 && chkHour <= 6 {
		// 日志只保留7天在远端
		delTime := time.Now().AddDate(0, 0, -7).Unix()
		qpsDelTime := time.Now().Add(-20 * time.Hour)

		r.orm.Where("timestmap <= ?", delTime).Delete(new(model.TbLogs))
		r.orm.Where("timestamp <= ?", delTime).Delete(new(model.TbRemapdata))
		r.orm.Where("timestamp <= ?", delTime).Delete(new(model.TbSingle))
		r.orm.Where("timestamp <= ?", delTime).Delete(new(model.TbHcollect))
		r.orm.Where("timestamp <= ?", qpsDelTime).Delete(new(model.TbQps))
	}
}

// 初始化这个记录集
func (r *PgsqlDBRecord) InitRecord() bool {
	engine, err := xorm.NewEngine("postgres", config.C.PgsqlAddr)
	if err != nil {
		log.Error("postgres init faild:", err)
		return false
	}

	r.orm = engine

	engine.Sync2(new(model.TbRemapdata))
	engine.Sync2(new(model.TbHcollect))
	engine.Sync2(new(model.TbHcollect))
	engine.Sync2(new(model.TbCachemiss))
	engine.Sync2(new(model.TbLogs))
	engine.Sync2(new(model.TbQps))
	engine.Sync2(new(model.TbSingle))

	grapeTimer.NewTickerLoop(30*60*1000, -1, r.procTimerRemove)

	return true
}
