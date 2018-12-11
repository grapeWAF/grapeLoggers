package config

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/go-ini/ini"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
)

type Server struct {
	HttpAddr string
	Verion   []string `delim:","`
}

type Log struct {
	LogPath  string
	LogLevel int
}

type Database struct {
	DBtype  string
	TimeLoc string
}

type MysqlData struct {
	MysqlDB string
}

type PgsqlData struct {
	PgsqlAddr string
}

type InfluxData struct {
	InfluxDB       string
	InfluxName     string
	InfluxUser     string
	InfluxPassword string
}

type MongoData struct {
	MongoDB    string
	MgoDBName  string
	MgoColName string
}

type LogApi struct {
	AppScret string
	AppId    string
}

type LoggerConf struct {
	Name string `ini:"name"`
	Server
	Log
	Database
	MysqlData
	PgsqlData
	InfluxData
	MongoData
	LogApi
}

var C *LoggerConf = new(LoggerConf)
var Loc *time.Location = nil

func LoadConf() error {
	cf, err := ini.Load("conf/app.conf")
	if err != nil {
		return err
	}

	cf.MapTo(C)

	Loc, _ = time.LoadLocation(C.TimeLoc)

	return nil
}

func BuildLogger() {
	log.SetFormatter(&prefixed.TextFormatter{
		ForceFormatting:  true,
		ForceColors:      true,
		FullTimestamp:    true,
		DisableTimestamp: false,
		TimestampFormat:  "2006-01-02 15:04:05",
	})

	AddHookDefault()

	logpath := C.LogPath
	os.MkdirAll(path.Dir(logpath), os.ModePerm)
	writer, err := rotatelogs.New(
		logpath+".%Y%m%d%H%M",
		rotatelogs.WithLinkName(logpath),
		rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(604800)*time.Second),
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	debugWriter, derr := rotatelogs.New(
		logpath+".debug.%Y%m%d%H%M",
		rotatelogs.WithLinkName(logpath),
		rotatelogs.WithMaxAge(time.Duration(86400)*time.Second),
		rotatelogs.WithRotationTime(time.Duration(604800)*time.Second),
	)

	if derr != nil {
		fmt.Println(derr)
		return
	}

	log.AddHook(lfshook.NewHook(lfshook.WriterMap{
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.DebugLevel: debugWriter,
	}, &log.TextFormatter{}))

	log.SetLevel(log.Level(C.LogLevel))
}
