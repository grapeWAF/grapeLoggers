// Code generated by protoc-gen-go. DO NOT EDIT.
// source: logger.proto

/*
Package grape_logs is a generated protocol buffer package.

It is generated from these files:
	logger.proto

It has these top-level messages:
	LogMsgReq
	LogMsgResp
	LogMsgSearchReq
	LogMsgResult
	LogMsgSearchResp
	HostInfoDataReq
	HostInfoDataResp
	HostCollReq
	HostCollResp
	SingleHostDataReq
	SingleHostDataResp
	QPSDataReq
	QPSDataResp
	RemapItem
	RemapCommitReq
	RemapCommitResp
*/
package grape_logs

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type LogMsgReq struct {
	LogJsonMsg string `protobuf:"bytes,1,opt,name=LogJsonMsg" json:"LogJsonMsg,omitempty"`
}

func (m *LogMsgReq) Reset()                    { *m = LogMsgReq{} }
func (m *LogMsgReq) String() string            { return proto.CompactTextString(m) }
func (*LogMsgReq) ProtoMessage()               {}
func (*LogMsgReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *LogMsgReq) GetLogJsonMsg() string {
	if m != nil {
		return m.LogJsonMsg
	}
	return ""
}

type LogMsgResp struct {
	RCode int32 `protobuf:"varint,1,opt,name=rCode" json:"rCode,omitempty"`
}

func (m *LogMsgResp) Reset()                    { *m = LogMsgResp{} }
func (m *LogMsgResp) String() string            { return proto.CompactTextString(m) }
func (*LogMsgResp) ProtoMessage()               {}
func (*LogMsgResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *LogMsgResp) GetRCode() int32 {
	if m != nil {
		return m.RCode
	}
	return 0
}

type LogMsgSearchReq struct {
	Type       string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
	Host       string `protobuf:"bytes,2,opt,name=host" json:"host,omitempty"`
	Key        string `protobuf:"bytes,3,opt,name=key" json:"key,omitempty"`
	BeginDate  string `protobuf:"bytes,4,opt,name=beginDate" json:"beginDate,omitempty"`
	EndDate    string `protobuf:"bytes,5,opt,name=endDate" json:"endDate,omitempty"`
	PageNumber int32  `protobuf:"varint,6,opt,name=page_number,json=pageNumber" json:"page_number,omitempty"`
}

func (m *LogMsgSearchReq) Reset()                    { *m = LogMsgSearchReq{} }
func (m *LogMsgSearchReq) String() string            { return proto.CompactTextString(m) }
func (*LogMsgSearchReq) ProtoMessage()               {}
func (*LogMsgSearchReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *LogMsgSearchReq) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *LogMsgSearchReq) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func (m *LogMsgSearchReq) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *LogMsgSearchReq) GetBeginDate() string {
	if m != nil {
		return m.BeginDate
	}
	return ""
}

func (m *LogMsgSearchReq) GetEndDate() string {
	if m != nil {
		return m.EndDate
	}
	return ""
}

func (m *LogMsgSearchReq) GetPageNumber() int32 {
	if m != nil {
		return m.PageNumber
	}
	return 0
}

type LogMsgResult struct {
	Type   string `protobuf:"bytes,1,opt,name=type" json:"type,omitempty"`
	Host   string `protobuf:"bytes,2,opt,name=host" json:"host,omitempty"`
	Time   string `protobuf:"bytes,3,opt,name=time" json:"time,omitempty"`
	Logmsg string `protobuf:"bytes,4,opt,name=logmsg" json:"logmsg,omitempty"`
	Caller string `protobuf:"bytes,5,opt,name=caller" json:"caller,omitempty"`
	Data   string `protobuf:"bytes,6,opt,name=data" json:"data,omitempty"`
}

func (m *LogMsgResult) Reset()                    { *m = LogMsgResult{} }
func (m *LogMsgResult) String() string            { return proto.CompactTextString(m) }
func (*LogMsgResult) ProtoMessage()               {}
func (*LogMsgResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *LogMsgResult) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

func (m *LogMsgResult) GetHost() string {
	if m != nil {
		return m.Host
	}
	return ""
}

func (m *LogMsgResult) GetTime() string {
	if m != nil {
		return m.Time
	}
	return ""
}

func (m *LogMsgResult) GetLogmsg() string {
	if m != nil {
		return m.Logmsg
	}
	return ""
}

func (m *LogMsgResult) GetCaller() string {
	if m != nil {
		return m.Caller
	}
	return ""
}

func (m *LogMsgResult) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

type LogMsgSearchResp struct {
	LogNum  int32           `protobuf:"varint,1,opt,name=logNum" json:"logNum,omitempty"`
	PageNum int32           `protobuf:"varint,2,opt,name=pageNum" json:"pageNum,omitempty"`
	Req     []*LogMsgResult `protobuf:"bytes,3,rep,name=req" json:"req,omitempty"`
}

func (m *LogMsgSearchResp) Reset()                    { *m = LogMsgSearchResp{} }
func (m *LogMsgSearchResp) String() string            { return proto.CompactTextString(m) }
func (*LogMsgSearchResp) ProtoMessage()               {}
func (*LogMsgSearchResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *LogMsgSearchResp) GetLogNum() int32 {
	if m != nil {
		return m.LogNum
	}
	return 0
}

func (m *LogMsgSearchResp) GetPageNum() int32 {
	if m != nil {
		return m.PageNum
	}
	return 0
}

func (m *LogMsgSearchResp) GetReq() []*LogMsgResult {
	if m != nil {
		return m.Req
	}
	return nil
}

type HostInfoDataReq struct {
	Timestamp      int64     `protobuf:"varint,1,opt,name=timestamp" json:"timestamp,omitempty"`
	HostName       string    `protobuf:"bytes,2,opt,name=hostName" json:"hostName,omitempty"`
	RemoteAddr     string    `protobuf:"bytes,3,opt,name=remoteAddr" json:"remoteAddr,omitempty"`
	System         string    `protobuf:"bytes,4,opt,name=system" json:"system,omitempty"`
	ByteSent       uint64    `protobuf:"varint,5,opt,name=byteSent" json:"byteSent,omitempty"`
	ByteRecv       uint64    `protobuf:"varint,6,opt,name=byteRecv" json:"byteRecv,omitempty"`
	SysMemPercent  float32   `protobuf:"fixed32,7,opt,name=sysMemPercent" json:"sysMemPercent,omitempty"`
	SysMemFree     uint64    `protobuf:"varint,8,opt,name=sysMemFree" json:"sysMemFree,omitempty"`
	SysMemUsed     uint64    `protobuf:"varint,9,opt,name=sysMemUsed" json:"sysMemUsed,omitempty"`
	ProcMemPercent float32   `protobuf:"fixed32,10,opt,name=procMemPercent" json:"procMemPercent,omitempty"`
	CpuType        string    `protobuf:"bytes,11,opt,name=cpuType" json:"cpuType,omitempty"`
	CpuPercent     []float32 `protobuf:"fixed32,12,rep,packed,name=cpuPercent" json:"cpuPercent,omitempty"`
	CpuTotal       float32   `protobuf:"fixed32,13,opt,name=cpuTotal" json:"cpuTotal,omitempty"`
	ProcCpuPercent float32   `protobuf:"fixed32,14,opt,name=procCpuPercent" json:"procCpuPercent,omitempty"`
	ProcNumFD      int32     `protobuf:"varint,15,opt,name=ProcNumFD" json:"ProcNumFD,omitempty"`
	// 硬盘信息
	HdTotal   uint64  `protobuf:"varint,16,opt,name=hdTotal" json:"hdTotal,omitempty"`
	HdFree    uint64  `protobuf:"varint,17,opt,name=hdFree" json:"hdFree,omitempty"`
	HdPercent float32 `protobuf:"fixed32,18,opt,name=hdPercent" json:"hdPercent,omitempty"`
	HostUID   string  `protobuf:"bytes,19,opt,name=hostUID" json:"hostUID,omitempty"`
	// 访问频率
	RQMCount uint64 `protobuf:"varint,20,opt,name=RQMCount" json:"RQMCount,omitempty"`
	// 每次启动开始计算 重启则归零
	TotalRQ uint64 `protobuf:"varint,21,opt,name=totalRQ" json:"totalRQ,omitempty"`
	// 总带宽流量
	TotalbyteSent uint64 `protobuf:"varint,22,opt,name=totalbyteSent" json:"totalbyteSent,omitempty"`
	TotalbyteRecv uint64 `protobuf:"varint,23,opt,name=totalbyteRecv" json:"totalbyteRecv,omitempty"`
}

func (m *HostInfoDataReq) Reset()                    { *m = HostInfoDataReq{} }
func (m *HostInfoDataReq) String() string            { return proto.CompactTextString(m) }
func (*HostInfoDataReq) ProtoMessage()               {}
func (*HostInfoDataReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *HostInfoDataReq) GetTimestamp() int64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *HostInfoDataReq) GetHostName() string {
	if m != nil {
		return m.HostName
	}
	return ""
}

func (m *HostInfoDataReq) GetRemoteAddr() string {
	if m != nil {
		return m.RemoteAddr
	}
	return ""
}

func (m *HostInfoDataReq) GetSystem() string {
	if m != nil {
		return m.System
	}
	return ""
}

func (m *HostInfoDataReq) GetByteSent() uint64 {
	if m != nil {
		return m.ByteSent
	}
	return 0
}

func (m *HostInfoDataReq) GetByteRecv() uint64 {
	if m != nil {
		return m.ByteRecv
	}
	return 0
}

func (m *HostInfoDataReq) GetSysMemPercent() float32 {
	if m != nil {
		return m.SysMemPercent
	}
	return 0
}

func (m *HostInfoDataReq) GetSysMemFree() uint64 {
	if m != nil {
		return m.SysMemFree
	}
	return 0
}

func (m *HostInfoDataReq) GetSysMemUsed() uint64 {
	if m != nil {
		return m.SysMemUsed
	}
	return 0
}

func (m *HostInfoDataReq) GetProcMemPercent() float32 {
	if m != nil {
		return m.ProcMemPercent
	}
	return 0
}

func (m *HostInfoDataReq) GetCpuType() string {
	if m != nil {
		return m.CpuType
	}
	return ""
}

func (m *HostInfoDataReq) GetCpuPercent() []float32 {
	if m != nil {
		return m.CpuPercent
	}
	return nil
}

func (m *HostInfoDataReq) GetCpuTotal() float32 {
	if m != nil {
		return m.CpuTotal
	}
	return 0
}

func (m *HostInfoDataReq) GetProcCpuPercent() float32 {
	if m != nil {
		return m.ProcCpuPercent
	}
	return 0
}

func (m *HostInfoDataReq) GetProcNumFD() int32 {
	if m != nil {
		return m.ProcNumFD
	}
	return 0
}

func (m *HostInfoDataReq) GetHdTotal() uint64 {
	if m != nil {
		return m.HdTotal
	}
	return 0
}

func (m *HostInfoDataReq) GetHdFree() uint64 {
	if m != nil {
		return m.HdFree
	}
	return 0
}

func (m *HostInfoDataReq) GetHdPercent() float32 {
	if m != nil {
		return m.HdPercent
	}
	return 0
}

func (m *HostInfoDataReq) GetHostUID() string {
	if m != nil {
		return m.HostUID
	}
	return ""
}

func (m *HostInfoDataReq) GetRQMCount() uint64 {
	if m != nil {
		return m.RQMCount
	}
	return 0
}

func (m *HostInfoDataReq) GetTotalRQ() uint64 {
	if m != nil {
		return m.TotalRQ
	}
	return 0
}

func (m *HostInfoDataReq) GetTotalbyteSent() uint64 {
	if m != nil {
		return m.TotalbyteSent
	}
	return 0
}

func (m *HostInfoDataReq) GetTotalbyteRecv() uint64 {
	if m != nil {
		return m.TotalbyteRecv
	}
	return 0
}

type HostInfoDataResp struct {
	Req int32 `protobuf:"varint,1,opt,name=req" json:"req,omitempty"`
}

func (m *HostInfoDataResp) Reset()                    { *m = HostInfoDataResp{} }
func (m *HostInfoDataResp) String() string            { return proto.CompactTextString(m) }
func (*HostInfoDataResp) ProtoMessage()               {}
func (*HostInfoDataResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *HostInfoDataResp) GetReq() int32 {
	if m != nil {
		return m.Req
	}
	return 0
}

type HostCollReq struct {
	HostName string `protobuf:"bytes,1,opt,name=hostName" json:"hostName,omitempty"`
	HostUUID string `protobuf:"bytes,2,opt,name=hostUUID" json:"hostUUID,omitempty"`
	Limit    int32  `protobuf:"varint,3,opt,name=limit" json:"limit,omitempty"`
}

func (m *HostCollReq) Reset()                    { *m = HostCollReq{} }
func (m *HostCollReq) String() string            { return proto.CompactTextString(m) }
func (*HostCollReq) ProtoMessage()               {}
func (*HostCollReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *HostCollReq) GetHostName() string {
	if m != nil {
		return m.HostName
	}
	return ""
}

func (m *HostCollReq) GetHostUUID() string {
	if m != nil {
		return m.HostUUID
	}
	return ""
}

func (m *HostCollReq) GetLimit() int32 {
	if m != nil {
		return m.Limit
	}
	return 0
}

type HostCollResp struct {
	Req   int32    `protobuf:"varint,1,opt,name=req" json:"req,omitempty"`
	Datas []string `protobuf:"bytes,2,rep,name=datas" json:"datas,omitempty"`
}

func (m *HostCollResp) Reset()                    { *m = HostCollResp{} }
func (m *HostCollResp) String() string            { return proto.CompactTextString(m) }
func (*HostCollResp) ProtoMessage()               {}
func (*HostCollResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

func (m *HostCollResp) GetReq() int32 {
	if m != nil {
		return m.Req
	}
	return 0
}

func (m *HostCollResp) GetDatas() []string {
	if m != nil {
		return m.Datas
	}
	return nil
}

type SingleHostDataReq struct {
	HostName  string `protobuf:"bytes,1,opt,name=hostName" json:"hostName,omitempty"`
	MachineID string `protobuf:"bytes,2,opt,name=machineID" json:"machineID,omitempty"`
	JsonBody  string `protobuf:"bytes,3,opt,name=jsonBody" json:"jsonBody,omitempty"`
}

func (m *SingleHostDataReq) Reset()                    { *m = SingleHostDataReq{} }
func (m *SingleHostDataReq) String() string            { return proto.CompactTextString(m) }
func (*SingleHostDataReq) ProtoMessage()               {}
func (*SingleHostDataReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

func (m *SingleHostDataReq) GetHostName() string {
	if m != nil {
		return m.HostName
	}
	return ""
}

func (m *SingleHostDataReq) GetMachineID() string {
	if m != nil {
		return m.MachineID
	}
	return ""
}

func (m *SingleHostDataReq) GetJsonBody() string {
	if m != nil {
		return m.JsonBody
	}
	return ""
}

type SingleHostDataResp struct {
	Req int32  `protobuf:"varint,1,opt,name=req" json:"req,omitempty"`
	Msg string `protobuf:"bytes,2,opt,name=msg" json:"msg,omitempty"`
}

func (m *SingleHostDataResp) Reset()                    { *m = SingleHostDataResp{} }
func (m *SingleHostDataResp) String() string            { return proto.CompactTextString(m) }
func (*SingleHostDataResp) ProtoMessage()               {}
func (*SingleHostDataResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *SingleHostDataResp) GetReq() int32 {
	if m != nil {
		return m.Req
	}
	return 0
}

func (m *SingleHostDataResp) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type QPSDataReq struct {
	Qps int32 `protobuf:"varint,1,opt,name=qps" json:"qps,omitempty"`
	Pv  int32 `protobuf:"varint,2,opt,name=pv" json:"pv,omitempty"`
}

func (m *QPSDataReq) Reset()                    { *m = QPSDataReq{} }
func (m *QPSDataReq) String() string            { return proto.CompactTextString(m) }
func (*QPSDataReq) ProtoMessage()               {}
func (*QPSDataReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func (m *QPSDataReq) GetQps() int32 {
	if m != nil {
		return m.Qps
	}
	return 0
}

func (m *QPSDataReq) GetPv() int32 {
	if m != nil {
		return m.Pv
	}
	return 0
}

type QPSDataResp struct {
	Req int32 `protobuf:"varint,1,opt,name=req" json:"req,omitempty"`
}

func (m *QPSDataResp) Reset()                    { *m = QPSDataResp{} }
func (m *QPSDataResp) String() string            { return proto.CompactTextString(m) }
func (*QPSDataResp) ProtoMessage()               {}
func (*QPSDataResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

func (m *QPSDataResp) GetReq() int32 {
	if m != nil {
		return m.Req
	}
	return 0
}

type RemapItem struct {
	FromSrc   string `protobuf:"bytes,1,opt,name=FromSrc" json:"FromSrc,omitempty"`
	RecvBytes int64  `protobuf:"varint,2,opt,name=RecvBytes" json:"RecvBytes,omitempty"`
	SendBytes int64  `protobuf:"varint,3,opt,name=SendBytes" json:"SendBytes,omitempty"`
	Rqtotal   int32  `protobuf:"varint,4,opt,name=Rqtotal" json:"Rqtotal,omitempty"`
	Time      int64  `protobuf:"varint,5,opt,name=time" json:"time,omitempty"`
}

func (m *RemapItem) Reset()                    { *m = RemapItem{} }
func (m *RemapItem) String() string            { return proto.CompactTextString(m) }
func (*RemapItem) ProtoMessage()               {}
func (*RemapItem) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

func (m *RemapItem) GetFromSrc() string {
	if m != nil {
		return m.FromSrc
	}
	return ""
}

func (m *RemapItem) GetRecvBytes() int64 {
	if m != nil {
		return m.RecvBytes
	}
	return 0
}

func (m *RemapItem) GetSendBytes() int64 {
	if m != nil {
		return m.SendBytes
	}
	return 0
}

func (m *RemapItem) GetRqtotal() int32 {
	if m != nil {
		return m.Rqtotal
	}
	return 0
}

func (m *RemapItem) GetTime() int64 {
	if m != nil {
		return m.Time
	}
	return 0
}

type RemapCommitReq struct {
	Data []*RemapItem `protobuf:"bytes,1,rep,name=data" json:"data,omitempty"`
}

func (m *RemapCommitReq) Reset()                    { *m = RemapCommitReq{} }
func (m *RemapCommitReq) String() string            { return proto.CompactTextString(m) }
func (*RemapCommitReq) ProtoMessage()               {}
func (*RemapCommitReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

func (m *RemapCommitReq) GetData() []*RemapItem {
	if m != nil {
		return m.Data
	}
	return nil
}

type RemapCommitResp struct {
	Req int32 `protobuf:"varint,1,opt,name=req" json:"req,omitempty"`
}

func (m *RemapCommitResp) Reset()                    { *m = RemapCommitResp{} }
func (m *RemapCommitResp) String() string            { return proto.CompactTextString(m) }
func (*RemapCommitResp) ProtoMessage()               {}
func (*RemapCommitResp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

func (m *RemapCommitResp) GetReq() int32 {
	if m != nil {
		return m.Req
	}
	return 0
}

func init() {
	proto.RegisterType((*LogMsgReq)(nil), "grape.logs.LogMsgReq")
	proto.RegisterType((*LogMsgResp)(nil), "grape.logs.LogMsgResp")
	proto.RegisterType((*LogMsgSearchReq)(nil), "grape.logs.LogMsgSearchReq")
	proto.RegisterType((*LogMsgResult)(nil), "grape.logs.LogMsgResult")
	proto.RegisterType((*LogMsgSearchResp)(nil), "grape.logs.LogMsgSearchResp")
	proto.RegisterType((*HostInfoDataReq)(nil), "grape.logs.HostInfoDataReq")
	proto.RegisterType((*HostInfoDataResp)(nil), "grape.logs.HostInfoDataResp")
	proto.RegisterType((*HostCollReq)(nil), "grape.logs.HostCollReq")
	proto.RegisterType((*HostCollResp)(nil), "grape.logs.HostCollResp")
	proto.RegisterType((*SingleHostDataReq)(nil), "grape.logs.SingleHostDataReq")
	proto.RegisterType((*SingleHostDataResp)(nil), "grape.logs.SingleHostDataResp")
	proto.RegisterType((*QPSDataReq)(nil), "grape.logs.QPSDataReq")
	proto.RegisterType((*QPSDataResp)(nil), "grape.logs.QPSDataResp")
	proto.RegisterType((*RemapItem)(nil), "grape.logs.RemapItem")
	proto.RegisterType((*RemapCommitReq)(nil), "grape.logs.RemapCommitReq")
	proto.RegisterType((*RemapCommitResp)(nil), "grape.logs.RemapCommitResp")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for LogService service

type LogServiceClient interface {
	AddLog(ctx context.Context, in *LogMsgReq, opts ...grpc.CallOption) (*LogMsgResp, error)
	SearchLog(ctx context.Context, in *LogMsgSearchReq, opts ...grpc.CallOption) (*LogMsgSearchResp, error)
	// 采集信息获取
	SubmitHost(ctx context.Context, in *HostInfoDataReq, opts ...grpc.CallOption) (*HostInfoDataResp, error)
	GetHostDatas(ctx context.Context, in *HostCollReq, opts ...grpc.CallOption) (*HostCollResp, error)
	SubmitSingleHost(ctx context.Context, in *SingleHostDataReq, opts ...grpc.CallOption) (*SingleHostDataResp, error)
	QPSDataCommit(ctx context.Context, in *QPSDataReq, opts ...grpc.CallOption) (*QPSDataResp, error)
	RemapCollCommit(ctx context.Context, in *RemapCommitReq, opts ...grpc.CallOption) (*RemapCommitResp, error)
}

type logServiceClient struct {
	cc *grpc.ClientConn
}

func NewLogServiceClient(cc *grpc.ClientConn) LogServiceClient {
	return &logServiceClient{cc}
}

func (c *logServiceClient) AddLog(ctx context.Context, in *LogMsgReq, opts ...grpc.CallOption) (*LogMsgResp, error) {
	out := new(LogMsgResp)
	err := grpc.Invoke(ctx, "/grape.logs.LogService/AddLog", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logServiceClient) SearchLog(ctx context.Context, in *LogMsgSearchReq, opts ...grpc.CallOption) (*LogMsgSearchResp, error) {
	out := new(LogMsgSearchResp)
	err := grpc.Invoke(ctx, "/grape.logs.LogService/SearchLog", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logServiceClient) SubmitHost(ctx context.Context, in *HostInfoDataReq, opts ...grpc.CallOption) (*HostInfoDataResp, error) {
	out := new(HostInfoDataResp)
	err := grpc.Invoke(ctx, "/grape.logs.LogService/SubmitHost", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logServiceClient) GetHostDatas(ctx context.Context, in *HostCollReq, opts ...grpc.CallOption) (*HostCollResp, error) {
	out := new(HostCollResp)
	err := grpc.Invoke(ctx, "/grape.logs.LogService/GetHostDatas", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logServiceClient) SubmitSingleHost(ctx context.Context, in *SingleHostDataReq, opts ...grpc.CallOption) (*SingleHostDataResp, error) {
	out := new(SingleHostDataResp)
	err := grpc.Invoke(ctx, "/grape.logs.LogService/SubmitSingleHost", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logServiceClient) QPSDataCommit(ctx context.Context, in *QPSDataReq, opts ...grpc.CallOption) (*QPSDataResp, error) {
	out := new(QPSDataResp)
	err := grpc.Invoke(ctx, "/grape.logs.LogService/QPSDataCommit", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *logServiceClient) RemapCollCommit(ctx context.Context, in *RemapCommitReq, opts ...grpc.CallOption) (*RemapCommitResp, error) {
	out := new(RemapCommitResp)
	err := grpc.Invoke(ctx, "/grape.logs.LogService/RemapCollCommit", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for LogService service

type LogServiceServer interface {
	AddLog(context.Context, *LogMsgReq) (*LogMsgResp, error)
	SearchLog(context.Context, *LogMsgSearchReq) (*LogMsgSearchResp, error)
	// 采集信息获取
	SubmitHost(context.Context, *HostInfoDataReq) (*HostInfoDataResp, error)
	GetHostDatas(context.Context, *HostCollReq) (*HostCollResp, error)
	SubmitSingleHost(context.Context, *SingleHostDataReq) (*SingleHostDataResp, error)
	QPSDataCommit(context.Context, *QPSDataReq) (*QPSDataResp, error)
	RemapCollCommit(context.Context, *RemapCommitReq) (*RemapCommitResp, error)
}

func RegisterLogServiceServer(s *grpc.Server, srv LogServiceServer) {
	s.RegisterService(&_LogService_serviceDesc, srv)
}

func _LogService_AddLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogMsgReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogServiceServer).AddLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grape.logs.LogService/AddLog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogServiceServer).AddLog(ctx, req.(*LogMsgReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogService_SearchLog_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogMsgSearchReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogServiceServer).SearchLog(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grape.logs.LogService/SearchLog",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogServiceServer).SearchLog(ctx, req.(*LogMsgSearchReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogService_SubmitHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HostInfoDataReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogServiceServer).SubmitHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grape.logs.LogService/SubmitHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogServiceServer).SubmitHost(ctx, req.(*HostInfoDataReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogService_GetHostDatas_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HostCollReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogServiceServer).GetHostDatas(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grape.logs.LogService/GetHostDatas",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogServiceServer).GetHostDatas(ctx, req.(*HostCollReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogService_SubmitSingleHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SingleHostDataReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogServiceServer).SubmitSingleHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grape.logs.LogService/SubmitSingleHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogServiceServer).SubmitSingleHost(ctx, req.(*SingleHostDataReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogService_QPSDataCommit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QPSDataReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogServiceServer).QPSDataCommit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grape.logs.LogService/QPSDataCommit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogServiceServer).QPSDataCommit(ctx, req.(*QPSDataReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _LogService_RemapCollCommit_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemapCommitReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogServiceServer).RemapCollCommit(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/grape.logs.LogService/RemapCollCommit",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogServiceServer).RemapCollCommit(ctx, req.(*RemapCommitReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _LogService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "grape.logs.LogService",
	HandlerType: (*LogServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddLog",
			Handler:    _LogService_AddLog_Handler,
		},
		{
			MethodName: "SearchLog",
			Handler:    _LogService_SearchLog_Handler,
		},
		{
			MethodName: "SubmitHost",
			Handler:    _LogService_SubmitHost_Handler,
		},
		{
			MethodName: "GetHostDatas",
			Handler:    _LogService_GetHostDatas_Handler,
		},
		{
			MethodName: "SubmitSingleHost",
			Handler:    _LogService_SubmitSingleHost_Handler,
		},
		{
			MethodName: "QPSDataCommit",
			Handler:    _LogService_QPSDataCommit_Handler,
		},
		{
			MethodName: "RemapCollCommit",
			Handler:    _LogService_RemapCollCommit_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "logger.proto",
}

func init() { proto.RegisterFile("logger.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 1002 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x56, 0xcd, 0x6e, 0xdb, 0xc6,
	0x13, 0x87, 0x44, 0xcb, 0x8e, 0xc6, 0xb2, 0xad, 0xec, 0x3f, 0xb6, 0x09, 0xc5, 0xff, 0x44, 0x60,
	0x83, 0x42, 0x6d, 0x01, 0x1d, 0x52, 0xa0, 0x2d, 0xd0, 0x4b, 0x1d, 0x19, 0x4e, 0x1d, 0xd8, 0x86,
	0xbd, 0xaa, 0x4f, 0x3d, 0x14, 0x14, 0xb9, 0xa5, 0xd8, 0x92, 0xdc, 0x15, 0x77, 0x65, 0x40, 0x2f,
	0xd1, 0xbe, 0x44, 0x5f, 0xa8, 0xaf, 0xd3, 0x53, 0x31, 0xc3, 0xe5, 0x87, 0x64, 0x25, 0xe8, 0x6d,
	0x7f, 0xbf, 0xf9, 0xde, 0x99, 0x1d, 0x12, 0x7a, 0x89, 0x8c, 0x22, 0x91, 0x8f, 0x55, 0x2e, 0x8d,
	0x64, 0x10, 0xe5, 0xbe, 0x12, 0xe3, 0x44, 0x46, 0xda, 0xfb, 0x0a, 0xba, 0xd7, 0x32, 0xba, 0xd1,
	0x11, 0x17, 0x0b, 0xf6, 0x0a, 0xe0, 0x5a, 0x46, 0x1f, 0xb4, 0xcc, 0x6e, 0x74, 0xe4, 0xb6, 0x86,
	0xad, 0x51, 0x97, 0x37, 0x18, 0xcf, 0x23, 0x39, 0x29, 0x6b, 0xc5, 0x5e, 0x40, 0x27, 0x9f, 0xc8,
	0x50, 0x90, 0x62, 0x87, 0x17, 0xc0, 0xfb, 0xab, 0x05, 0x47, 0x85, 0xd2, 0x54, 0xf8, 0x79, 0x30,
	0x47, 0xbf, 0x0c, 0x76, 0xcc, 0x4a, 0x09, 0xeb, 0x91, 0xce, 0xc8, 0xcd, 0xa5, 0x36, 0x6e, 0xbb,
	0xe0, 0xf0, 0xcc, 0xfa, 0xe0, 0xfc, 0x2e, 0x56, 0xae, 0x43, 0x14, 0x1e, 0xd9, 0x19, 0x74, 0x67,
	0x22, 0x8a, 0xb3, 0x0b, 0xdf, 0x08, 0x77, 0x87, 0xf8, 0x9a, 0x60, 0x2e, 0xec, 0x89, 0x2c, 0x24,
	0x59, 0x87, 0x64, 0x25, 0x64, 0xaf, 0x61, 0x5f, 0xf9, 0x91, 0xf8, 0x25, 0x5b, 0xa6, 0x33, 0x91,
	0xbb, 0xbb, 0x94, 0x21, 0x20, 0x75, 0x4b, 0x8c, 0xf7, 0x47, 0x0b, 0x7a, 0x55, 0x2d, 0xcb, 0xc4,
	0xfc, 0xe7, 0x1c, 0x51, 0x2f, 0x4e, 0x85, 0x4d, 0x92, 0xce, 0xec, 0x04, 0x76, 0x13, 0x19, 0xa5,
	0x3a, 0xb2, 0x29, 0x5a, 0x84, 0x7c, 0xe0, 0x27, 0x89, 0xc8, 0x6d, 0x7a, 0x16, 0xa1, 0x8f, 0xd0,
	0x37, 0x3e, 0xa5, 0xd5, 0xe5, 0x74, 0xf6, 0x14, 0xf4, 0xd7, 0xaf, 0x4d, 0x2b, 0xeb, 0xf7, 0x76,
	0x99, 0xda, 0x2b, 0xb6, 0x08, 0xeb, 0xb6, 0xa5, 0x50, 0x6a, 0x1d, 0x5e, 0x42, 0xf6, 0x25, 0x38,
	0xb9, 0x58, 0xb8, 0xce, 0xd0, 0x19, 0xed, 0xbf, 0x75, 0xc7, 0x75, 0xa3, 0xc7, 0xcd, 0x62, 0x39,
	0x2a, 0x79, 0x7f, 0x77, 0xe0, 0xe8, 0x47, 0xa9, 0xcd, 0x55, 0xf6, 0xab, 0xbc, 0xf0, 0x8d, 0x8f,
	0x9d, 0x3a, 0x83, 0x2e, 0x56, 0xa4, 0x8d, 0x9f, 0x2a, 0x0a, 0xea, 0xf0, 0x9a, 0x60, 0x03, 0x78,
	0x86, 0x77, 0x70, 0xeb, 0xa7, 0xc2, 0xde, 0x49, 0x85, 0x71, 0x76, 0x72, 0x91, 0x4a, 0x23, 0xce,
	0xc3, 0x30, 0xb7, 0xb7, 0xd3, 0x60, 0xb0, 0x16, 0xbd, 0xd2, 0x46, 0xa4, 0xe5, 0x1d, 0x15, 0x08,
	0x7d, 0xce, 0x56, 0x46, 0x4c, 0x45, 0x66, 0xe8, 0x96, 0x76, 0x78, 0x85, 0x4b, 0x19, 0x17, 0xc1,
	0x23, 0xdd, 0x95, 0x95, 0x21, 0x66, 0x6f, 0xe0, 0x40, 0xaf, 0xf4, 0x8d, 0x48, 0xef, 0x44, 0x1e,
	0xa0, 0xf1, 0xde, 0xb0, 0x35, 0x6a, 0xf3, 0x75, 0x12, 0xb3, 0x2a, 0x88, 0xcb, 0x5c, 0x08, 0xf7,
	0x19, 0xf9, 0x68, 0x30, 0xb5, 0xfc, 0x41, 0x8b, 0xd0, 0xed, 0x36, 0xe5, 0xc8, 0xb0, 0xcf, 0xe1,
	0x50, 0xe5, 0x32, 0x68, 0x84, 0x01, 0x0a, 0xb3, 0xc1, 0x62, 0x47, 0x02, 0xb5, 0xfc, 0x09, 0x07,
	0x68, 0xbf, 0x98, 0x44, 0x0b, 0x31, 0x42, 0xa0, 0x96, 0xa5, 0x75, 0x6f, 0xe8, 0x8c, 0xda, 0xbc,
	0xc1, 0x60, 0x8d, 0xa8, 0x2a, 0x8d, 0x9f, 0xb8, 0x07, 0xe4, 0xbb, 0xc2, 0x65, 0xf4, 0x49, 0x6d,
	0x7f, 0x58, 0x47, 0xaf, 0x59, 0xec, 0xda, 0x5d, 0x2e, 0x83, 0xdb, 0x65, 0x7a, 0x79, 0xe1, 0x1e,
	0xd1, 0x44, 0xd4, 0x04, 0xe6, 0x36, 0x0f, 0x8b, 0x00, 0x7d, 0x2a, 0xb0, 0x84, 0xd8, 0x93, 0x79,
	0x48, 0x37, 0xf3, 0x9c, 0x04, 0x16, 0xa1, 0xbf, 0x79, 0x58, 0x86, 0x64, 0x14, 0xb2, 0x26, 0xc8,
	0x9f, 0xd4, 0xe6, 0xe1, 0xea, 0xc2, 0xfd, 0x5f, 0x51, 0xab, 0x85, 0x58, 0x0b, 0xbf, 0xbf, 0x99,
	0xc8, 0x65, 0x66, 0xdc, 0x17, 0x45, 0xbf, 0x4a, 0x8c, 0x56, 0x06, 0x83, 0xf2, 0x7b, 0xf7, 0xb8,
	0xc8, 0xc2, 0x42, 0xec, 0x24, 0x1d, 0xab, 0x31, 0x38, 0x21, 0xf9, 0x3a, 0xb9, 0xa6, 0x45, 0x03,
	0x71, 0xba, 0xa1, 0x85, 0xa4, 0xf7, 0x06, 0xfa, 0xeb, 0x23, 0xad, 0x15, 0x6e, 0x15, 0x7c, 0x13,
	0xc5, 0x13, 0xa2, 0xc9, 0xff, 0x19, 0xf6, 0x51, 0x6b, 0x22, 0x93, 0x04, 0x87, 0xbe, 0x39, 0xd6,
	0xad, 0x8d, 0xb1, 0xb6, 0xb2, 0x07, 0xac, 0xb6, 0x31, 0xf2, 0x88, 0x71, 0x01, 0x26, 0x71, 0x1a,
	0x1b, 0x9a, 0xf6, 0x0e, 0x2f, 0x80, 0xf7, 0x0d, 0xf4, 0x6a, 0xe7, 0xdb, 0xc2, 0xa3, 0x1d, 0x3e,
	0x79, 0xed, 0xb6, 0x87, 0xce, 0xa8, 0xcb, 0x0b, 0xe0, 0xc5, 0xf0, 0x7c, 0x1a, 0x67, 0x51, 0x22,
	0xd0, 0xba, 0x7c, 0x8f, 0x9f, 0x4a, 0xed, 0x0c, 0xba, 0xa9, 0x1f, 0xcc, 0xe3, 0x4c, 0x54, 0xb9,
	0xd5, 0x04, 0x5a, 0xfe, 0xa6, 0x65, 0xf6, 0x4e, 0x86, 0xe5, 0x42, 0xad, 0xb0, 0xf7, 0x1d, 0xb0,
	0xcd, 0x50, 0x5b, 0x13, 0xed, 0x83, 0x83, 0x4b, 0xad, 0xf0, 0x8d, 0x47, 0x6f, 0x0c, 0x70, 0x7f,
	0x37, 0x2d, 0xb3, 0xeb, 0x83, 0xb3, 0x50, 0xba, 0xb4, 0x58, 0x28, 0xcd, 0x0e, 0xa1, 0xad, 0x1e,
	0xed, 0x52, 0x6a, 0xab, 0x47, 0xef, 0x35, 0xec, 0x57, 0xfa, 0x5b, 0x5b, 0xf1, 0x67, 0x0b, 0xba,
	0x5c, 0xa4, 0xbe, 0xba, 0xc2, 0x65, 0xe0, 0xc2, 0xde, 0x65, 0x2e, 0xd3, 0x69, 0x1e, 0xd8, 0x6a,
	0x4b, 0x88, 0xc5, 0x62, 0x83, 0xdf, 0xad, 0x8c, 0xd0, 0xe4, 0xdf, 0xe1, 0x35, 0x81, 0xd2, 0xa9,
	0xc8, 0xc2, 0x42, 0xea, 0x14, 0xd2, 0x8a, 0x40, 0xaf, 0x7c, 0x41, 0x73, 0x42, 0xbb, 0xa7, 0xc3,
	0x4b, 0x58, 0x2d, 0xf3, 0x0e, 0x99, 0xd0, 0xd9, 0xfb, 0x1e, 0x0e, 0x29, 0xa1, 0x89, 0x4c, 0xd3,
	0xd8, 0x60, 0x99, 0x5f, 0xd8, 0x75, 0xdd, 0xa2, 0xad, 0x7a, 0xdc, 0xdc, 0xaa, 0x55, 0xea, 0x76,
	0x8b, 0x7f, 0x06, 0x47, 0x6b, 0xc6, 0xdb, 0x6a, 0x7e, 0xfb, 0x8f, 0x43, 0xdf, 0xd1, 0xa9, 0xc8,
	0x1f, 0xe3, 0x40, 0xb0, 0x6f, 0x61, 0xf7, 0x3c, 0x0c, 0xaf, 0x65, 0xc4, 0x8e, 0xb7, 0x2d, 0xec,
	0xc5, 0xe0, 0x64, 0xeb, 0x1e, 0x57, 0xec, 0x12, 0xab, 0xc6, 0x8f, 0x05, 0xda, 0xbe, 0x7c, 0xaa,
	0x54, 0x7d, 0x80, 0x07, 0x67, 0x1f, 0x17, 0x6a, 0xc5, 0xde, 0x03, 0x4c, 0x97, 0xb3, 0x34, 0x36,
	0x38, 0x0e, 0xeb, 0x8e, 0x36, 0xbe, 0x0f, 0xeb, 0x8e, 0x9e, 0xbc, 0xb4, 0x73, 0xe8, 0xbd, 0x17,
	0xa6, 0x1c, 0x2a, 0xcd, 0x4e, 0x37, 0xb5, 0xed, 0x8b, 0x1b, 0xb8, 0xdb, 0x05, 0x5a, 0xb1, 0x7b,
	0xe8, 0x17, 0xb9, 0xd4, 0x03, 0xca, 0xfe, 0xdf, 0xd4, 0x7e, 0xf2, 0x46, 0x06, 0xaf, 0x3e, 0x25,
	0xd6, 0x8a, 0xfd, 0x00, 0x07, 0x76, 0x06, 0x8b, 0xae, 0xb0, 0xb5, 0xfb, 0xac, 0xc7, 0x79, 0x70,
	0xba, 0x95, 0xd7, 0x8a, 0x7d, 0xa8, 0xba, 0x9a, 0x24, 0xd6, 0xc7, 0xe0, 0xc9, 0x14, 0x54, 0xf3,
	0x32, 0x78, 0xf9, 0x51, 0x99, 0x56, 0xb3, 0x5d, 0xfa, 0x07, 0xfb, 0xfa, 0xdf, 0x00, 0x00, 0x00,
	0xff, 0xff, 0x92, 0x54, 0x1c, 0xe6, 0x93, 0x09, 0x00, 0x00,
}
