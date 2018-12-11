package clientv1

import (
	"fmt"
	"path"
	"runtime"
	"strings"

	log "github.com/sirupsen/logrus"

	util "github.com/koangel/grapeNet/Utils"
)

type WlogField = log.Fields

type LogNode struct {
	LogType   int
	LogMsg    string
	LogArgs   []interface{}
	LogFields WlogField
	LogCaller string
}

const (
	logError = iota
	logErrorf

	logInfo
	logInfof

	logWarn
	logWarnf

	logDebug
	logDebugf
)

var (
	logsQueue = util.NewSQueue()
)

func init() {
	for i := 0; i < 12; i++ {
		go procCommitLogs()
	}
}

func procCommitLogs() {
	for {

		logMsg := logsQueue.Pop().(*LogNode)
		entry := log.WithFields(logMsg.LogFields)

		switch logMsg.LogType {
		case logError:
			entry.Error(logMsg.LogArgs)
		case logErrorf:
			entry.Errorf(logMsg.LogMsg)

		case logInfo:
			entry.Info(logMsg.LogArgs)
		case logInfof:
			entry.Infof(logMsg.LogMsg)

		case logWarn:
			entry.Warn(logMsg.LogArgs)
		case logWarnf:
			entry.Warnf(logMsg.LogMsg)

		case logDebug:
			entry.Debug(logMsg.LogArgs)
		case logDebugf:
			entry.Debugf(logMsg.LogMsg)
		}
	}
}

func (r *LogNode) getCaller() {
	caller := "unknow"
	pc := make([]uintptr, 64)
	cnt := runtime.Callers(skipFrameCnt, pc)

	for i := 0; i < cnt; i++ {
		fu := runtime.FuncForPC(pc[i])
		name := fu.Name()
		if !strings.Contains(name, "grapeLoggers/clientv1") {
			file, line := fu.FileLine(pc[i] - 1)
			caller = fmt.Sprintf("%v:%v", path.Base(file), line)
			break
		}
	}

	caller = strings.Replace(caller, ".go", ".cc", -1)
	r.LogCaller = caller

	if r.LogFields == nil {
		r.LogFields = WlogField{
			"caller": caller,
		}
	} else {
		r.LogFields["caller"] = caller
	}

}

func (r *LogNode) appendMsg(vtype int, msg string) {
	r.getCaller()
	r.LogMsg = msg
	r.LogType = vtype
	logsQueue.Push(r)
}

func (r *LogNode) appendArgs(vtype int, args ...interface{}) {
	r.getCaller()
	r.LogArgs = []interface{}(args)
	r.LogType = vtype
	logsQueue.Push(r)
}

func WithField(key string, value interface{}) *LogNode {
	return &LogNode{
		LogFields: WlogField{
			key: value,
		},
	}
}

func WithFields(fields WlogField) *LogNode {
	rs := WlogField{}
	for k, v := range fields {
		rs[k] = v
	}

	return &LogNode{
		LogFields: rs,
	}
}

func (r *LogNode) Error(args ...interface{}) {
	r.appendArgs(logInfo, args...)
}

func (r *LogNode) Errorf(msg string, args ...interface{}) {
	r.appendMsg(logErrorf, fmt.Sprintf(msg, args...))
}

func (r *LogNode) Info(args ...interface{}) {
	r.appendArgs(logInfo, args...)
}

func (r *LogNode) Infof(msg string, args ...interface{}) {
	r.appendMsg(logInfof, fmt.Sprintf(msg, args...))
}

func (r *LogNode) Warn(args ...interface{}) {
	r.appendArgs(logWarn, args...)
}

func (r *LogNode) Warnf(msg string, args ...interface{}) {
	r.appendMsg(logWarnf, fmt.Sprintf(msg, args...))
}

func (r *LogNode) Debug(args ...interface{}) {
	r.appendArgs(logDebug, args...)
}

func (r *LogNode) Debugf(msg string, args ...interface{}) {
	r.appendMsg(logDebugf, fmt.Sprintf(msg, args...))
}

func Error(args ...interface{}) {
	r := &LogNode{}
	r.appendArgs(logError, args...)
}

func Errorf(msg string, args ...interface{}) {
	r := &LogNode{}
	r.appendMsg(logErrorf, fmt.Sprintf(msg, args...))
}

func Info(args ...interface{}) {
	r := &LogNode{}
	r.appendArgs(logInfo, args...)
}

func Infof(msg string, args ...interface{}) {
	r := &LogNode{}
	r.appendMsg(logInfof, fmt.Sprintf(msg, args...))
}

func Warn(args ...interface{}) {
	r := &LogNode{}
	r.appendArgs(logWarn, args...)
}

func Warnf(msg string, args ...interface{}) {
	r := &LogNode{}
	r.appendMsg(logWarnf, fmt.Sprintf(msg, args...))
}

func Debug(args ...interface{}) {
	r := &LogNode{}
	r.appendArgs(logDebug, args...)
}

func Debugf(msg string, args ...interface{}) {
	r := &LogNode{}
	r.appendMsg(logDebugf, fmt.Sprintf(msg, args...))
}
