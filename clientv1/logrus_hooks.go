package clientv1

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	proto "grapeLoggers/protos"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

////////////////////////////////////////////////////////////
//  add by koangel
//	date: 2017/09/20
//  日志记录，在提交的过程中将日志提交给远端
////////////////////////////////////////////////////////////
const skipFrameCnt = 3

type LOGAPI_Config struct {
	Urls      []string
	TryCount  int
	AppSecret string
	AppId     string
	Caller    bool
	Gzip      bool
	PoolSize  int
}

type LOGAPI_Hook struct {
	conf   LOGAPI_Config
	urlPos int

	HostName string

	asyncChan chan *LogMsg

	ctx     context.Context
	cancel  context.CancelFunc
	mu      *sync.Mutex
	once    sync.Once
	grpcCli []proto.LogServiceClient
}

func NewHook(conf LOGAPI_Config) *LOGAPI_Hook {
	hostName, _ := os.Hostname()
	hook := &LOGAPI_Hook{
		conf:      conf,
		HostName:  hostName,
		mu:        &sync.Mutex{},
		asyncChan: make(chan *LogMsg, 5120*1024),
	}

	hook.ctx, hook.cancel = context.WithCancel(context.Background())
	idx := 0
	for i := 0; i < hook.conf.PoolSize; i++ {

		// 连接一个池
		if idx >= len(hook.conf.Urls) {
			idx = 0
		}

		cli, err := grpc.Dial(hook.conf.Urls[idx], grpc.WithInsecure())
		if err != nil {
			fmt.Print("Dial Log Faild:", err, "\n")
			continue
		}

		hook.grpcCli = append(hook.grpcCli, proto.NewLogServiceClient(cli))

		go hook.procCommit(i)

		idx++
	}

	return hook
}

func SetupHook(conf LOGAPI_Config) {
	logrus.AddHook(NewHook(conf))
}

func (hook *LOGAPI_Hook) procCommit(index int) {
	for {
		select {
		case <-hook.ctx.Done():
			return
		case msg := <-hook.asyncChan:
			hook.post(hook.formatter(msg), index) // 异步提交
		}
	}
}

func (hook *LOGAPI_Hook) getUrl(path string) string {
	pos := hook.urlPos
	if pos > len(hook.conf.Urls) {
		pos = 0
	}

	return fmt.Sprintf("%s%s", hook.conf.Urls[pos], path)
}

func (hook *LOGAPI_Hook) formatterMsg(entry *logrus.Entry) *LogMsg {
	TypeName := "LOGS"
	if tyn, ok := entry.Data["type"]; ok {
		TypeName = tyn.(string)
	}

	caller := "unknow";
	if clc,ok := entry.Data["caller"]; ok {
		caller = clc.(string)
	}

	return &LogMsg{
		Version: "1.0",
		Host:    hook.HostName,
		Level:   int32(entry.Level),
		Type:    TypeName,
		Msg:     entry.Message,
		Time:    time.Now().Unix(),
		Caller:  caller,
		Extra:   entry.Data,
	}

}

func (hook *LOGAPI_Hook) formatter(entry *LogMsg) string {

	if hook.conf.Gzip {
		return entry.GZip()
	}

	return entry.Json()
}

func (hook *LOGAPI_Hook) faild() bool {
	hook.urlPos++
	if hook.urlPos >= len(hook.conf.Urls) {
		hook.urlPos = 0
		return false
	}

	return true
}

func (hook *LOGAPI_Hook) asyncPost(entry *LogMsg) {
	hook.asyncChan <- entry
}

func (hook *LOGAPI_Hook) post(json string, index int) error {
	_, err := hook.grpcCli[index].AddLog(
		context.Background(),
		&proto.LogMsgReq{
			LogJsonMsg: json,
		}, grpc.FailFast(true))

	if err != nil {
		fmt.Print("post Log Faild:", err, "\n")
		return err
	}

	return nil
}

func (hook *LOGAPI_Hook) Fire(entry *logrus.Entry) error {

	hook.asyncPost(hook.formatterMsg(entry))
	return nil
}

func (hook *LOGAPI_Hook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
		logrus.DebugLevel,
	}
}
