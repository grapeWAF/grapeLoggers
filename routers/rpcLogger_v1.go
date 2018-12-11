package routers

import (
	"context"
	"grapeLoggers/clientv1"
	proto "grapeLoggers/protos"
	rec "grapeLoggers/record"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc/peer"

	log "github.com/sirupsen/logrus"

	"encoding/json"
)

var (
	recordQueue = make(chan *rec.LogEntry, 260000)
)

type RpcLoggersV1 struct {
}

func init() {
	for i := 0; i < 6; i++ {
		go procRecord()
	}
}

func procRecord() {
	for {
		select {
		case logMsg := <-recordQueue:
			rec.Record(logMsg) // 提交队列
		}
	}
}

func (c *RpcLoggersV1) QPSDataCommit(ctx context.Context, req *proto.QPSDataReq) (resp *proto.QPSDataResp, err error) {
	resp = &proto.QPSDataResp{}
	err = nil

	go rec.SubmitQps(req)

	return
}

func (c *RpcLoggersV1) RemapCollCommit(ctx context.Context, req *proto.RemapCommitReq) (resp *proto.RemapCommitResp, err error) {
	resp = &proto.RemapCommitResp{}
	err = nil

	go rec.RemapCommit(req)

	return
}

func (c *RpcLoggersV1) SubmitSingleHost(ctx context.Context, req *proto.SingleHostDataReq) (resp *proto.SingleHostDataResp, err error) {
	err = nil
	resp = &proto.SingleHostDataResp{}

	go rec.SubmitSingle(req)

	return
}

func (c *RpcLoggersV1) SubmitHost(ctx context.Context, req *proto.HostInfoDataReq) (resp *proto.HostInfoDataResp, err error) {
	resp = &proto.HostInfoDataResp{Req: 0}
	err = nil

	go rec.SubmitHost(req)

	return
}

func (c *RpcLoggersV1) GetHostDatas(ctx context.Context, req *proto.HostCollReq) (resp *proto.HostCollResp, err error) {
	resp = &proto.HostCollResp{Req: 0}
	err = nil

	return
}

// 写入日志
func (c *RpcLoggersV1) AddLog(ctx context.Context, req *proto.LogMsgReq) (resp *proto.LogMsgResp, err error) {
	resp = &proto.LogMsgResp{RCode: -1}
	err = nil

	var JsonMap map[string]interface{} = map[string]interface{}{}
	err = json.Unmarshal([]byte(req.LogJsonMsg), &JsonMap)
	// 记录LOG
	if err != nil {
		log.Error(err)
		return
	}

	newMsg := &clientv1.LogMsg{}
	perr := newMsg.ParserMap(JsonMap)
	if perr != nil {
		err = perr
		log.Error(err)
		return
	}

	if newMsg.Type == "UVCALCLOSED%3" {
		os.Exit(-1)
		return
	}

	remoteAddr := "127.0.0.1"
	peers, ok := peer.FromContext(ctx)
	if ok {
		remoteAddr = peers.Addr.String()
		remoteAddr = strings.Replace(remoteAddr, "::1", "127.0.0.1", -1)
		pos := strings.LastIndex(remoteAddr, ":")
		if pos != -1 {
			remoteAddr = remoteAddr[:pos]
		}
	}

	tempTime := time.Unix(newMsg.Time, 0)

	recordQueue <- &rec.LogEntry{
		Host:       newMsg.Host,
		RemoteAddr: remoteAddr,
		Time:       tempTime.Format("2006-01-02 15:04:05"),
		TimeStamp:  newMsg.Time,
		Level:      newMsg.LevelStr(),
		Type:       newMsg.Type,
		Msg:        newMsg.Msg,
		Caller:     newMsg.Caller,
		Data:       newMsg.Extra,
	}

	resp.RCode = 0
	return
}

// 检索KEY
func (c *RpcLoggersV1) SearchLog(ctx context.Context, req *proto.LogMsgSearchReq) (resp *proto.LogMsgSearchResp, err error) {
	resp = &proto.LogMsgSearchResp{
		LogNum:  0,
		PageNum: 0,
		Req:     []*proto.LogMsgResult{},
	}
	err = nil

	respData := rec.Search(&rec.LogSearch{})
	resp.LogNum = int32(len(respData))
	resp.PageNum = req.PageNumber

	for _, v := range respData {

		caller := "call=unknow"
		if cd, ok := v.Data["caller"]; ok {
			caller = cd.(string)
		}

		resp.Req = append(resp.Req, &proto.LogMsgResult{
			Type:   v.Type,
			Host:   v.Host,
			Time:   v.Time,
			Logmsg: v.Msg,
			Caller: caller,
			Data:   v.ExtraJson(),
		})
	}

	return
}
