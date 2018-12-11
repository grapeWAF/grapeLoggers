////////////////////////////////////////////////////////////
//  add by koangel
//	date: 2017/09/19
//  本地配置文件记录
////////////////////////////////////////////////////////////
package grapeGuard

import (
	"fmt"
	"grapeLoggers/clientv1"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sync"
	"time"

	"github.com/koangel/grapeTimer"

	"github.com/x-cray/logrus-prefixed-formatter"

	"crypto/rc4"

	"github.com/go-ini/ini"
	"github.com/lestrrat/go-file-rotatelogs"
	"github.com/rifflock/lfshook"
	log "github.com/sirupsen/logrus"

	etcd "github.com/koangel/grapeNet/Etcd"

	"github.com/koangel/go-filemutex"
)

type COLLECTION struct {
	CEnable  bool `ini:"enable"`
	CTick    int  `ini:"collectionTick"`
	QPSTick  int  `ini:"qpsTick"`
	QpsLimit int  `ini:"limitQps"`
}

type MSGQUEUE struct {
	Address []string `ini:"Dials" delim:","`
}

type ATSCONF struct {
	SyncRemap   bool   `ini:"AutoSyncRemap"`
	RemapConf   string `ini:"RemapConf"`
	RemapReload string `ini:"RemapReload"`
}

type SERVER struct {
	BindAddr string `ini:"httpAddr"`
	TlsAddr  string `ini:"httpsAddr"`

	ReadHeaderTimeout int `ini:"ReadHeaderTimeout"`
	ReadTimeout       int `ini:"ReadTimeout"`
	WriteTimeout      int `ini:"WriteTimeout"`

	UseHttp2 bool `ini:"UseHttp2"`
}

type TLS struct {
	HostName string `ini:"hostName"`
	CertPath string `ini:"certPath"`
	KeyPath  string `ini:"keyPath"`
}

type LOG struct {
	Type  string `ini:"TYPE"`
	Level int    `ini:"LEVEL"` //DEBUG
	// 日志路径 默认为本地的logs 仅对FILE有效
	LogPath string `ini:"LogPath"` //=logs
	// 远端日志服务器 可以为多个 防止出现故障节点无法捕获日志 仅对Server有效
	LogServer []string `ini:"LogServer" delim:","` //http://localhost:3301,http://localhost:3301
	LogAppKey string   `ini:"LogAppKey"`
	LogAppId  string   `ini:"LogAppId"`
}

type ETCD struct {
	EAddrs    []string `ini:"etcd_address" delim:","`
	IsAuth    bool     `ini:"etcd_auth"`
	EUserName string   `ini:"etcd_user"`
	EPassword string   `ini:"etcd_pass"`
}

type MainConf struct {
	Name string `ini:"name"`
	Mode string `ini:"mode"`
	SERVER
	TLS
	LOG
	ETCD
	ATSCONF
	COLLECTION
}

type VerConf struct {
	Ver string `ini:"version"`
}

type TempConf struct {
	UseJSGuard bool `ini:"UseJSGuard"`
	UseIPSet   bool `ini:"UseIPSet"`
}

type NetConf struct {
	COLLECTION
	MSGQUEUE
}

var (
	NetC    *NetConf             = new(NetConf)
	Conf    *MainConf            = new(MainConf)
	V       *VerConf             = new(VerConf)
	TempC   *TempConf            = new(TempConf)
	filemux *filemutex.FileMutex = nil

	WorkPath       string = ""
	updaConf       string = ""
	lastUpdateTime int64  = 0
	uplocker       sync.Mutex
)

const (
	cryptKeys = "grapesoft_guard_@)!&)!@#"

	Version = "0.5.3.8 beta"
)

func MuxAndLock() {
	var merr error
	filemux, merr = filemutex.New("./guard.lock")
	if merr != nil {
		log.Error("无法创建进程锁:", merr)
		return
	}

	LockWait()
}

func LockWait() {
	if err := filemux.Lock(); err != nil {
		log.Info("进程互斥无法启动:", err)
		fmt.Print(err)
		return
	}
}

func Unlock() {
	filemux.Unlock()
}

func CryptRC4(data []byte) []byte {
	rcd, err := rc4.NewCipher([]byte(cryptKeys))
	if err != nil {
		log.Error(err)
		return data
	}

	dst := make([]byte, len(data))
	rcd.XORKeyStream(dst, data)

	return dst
}

func IsUseGuard() bool {
	uplocker.Lock()
	defer uplocker.Unlock()

	return TempC.UseJSGuard
}

func IsIPSet() bool {
	uplocker.Lock()
	defer uplocker.Unlock()

	return TempC.UseIPSet
}

func onUpdateConf() {
	if finfo, err := os.Stat(updaConf); err == nil {
		if finfo.ModTime().Unix() != lastUpdateTime {
			cfg, err := ini.Load(updaConf)
			if err != nil {
				log.Error(err)
				return
			}
			finfo, err := os.Stat(updaConf)
			if err == nil {
				lastUpdateTime = finfo.ModTime().Unix()
			}

			uplocker.Lock()
			cfg.MapTo(TempC)
			uplocker.Unlock()
		}
	}
}

func Load() error {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	WorkPath = dir

	cfg, err := ini.Load(WorkPath + "/conf/app.conf")
	if err != nil {
		return err
	}

	cfg.MapTo(Conf)

	BuildLogger()
	loadTemplates()
	LoadVer()
	LoadTempConf()
	return nil
}

func LoadEtcdConf() error {
	encode, err := etcd.Read("apps")
	if err != nil {
		return err
	}

	cfg, err := ini.Load(CryptRC4(encode))
	if err != nil {
		return err
	}

	cfg.MapTo(NetC)

	return nil
}

func SyncNetToEtcd() error {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	WorkPath = dir

	data, err := ioutil.ReadFile(WorkPath + "/apps_conf/app.conf")
	if err != nil {
		return err
	}

	etcd.Write("apps", CryptRC4(data))
	return nil
}

func LoadVer() error {
	cfg, err := ini.Load(WorkPath + "/conf/ver.conf")
	if err != nil {
		return err
	}

	cfg.MapTo(V)
	return nil
}

func LoadTempConf() error {
	updaConf = WorkPath + "/conf/temp.conf"
	cfg, err := ini.Load(updaConf)
	if err != nil {
		TempC.UseJSGuard = true
		return err
	}
	finfo, err := os.Stat(updaConf)
	if err == nil {
		lastUpdateTime = finfo.ModTime().Unix()
	}

	grapeTimer.NewTickerLoop(10*1000, -1, onUpdateConf)

	uplocker.Lock()
	cfg.MapTo(TempC)
	uplocker.Unlock()

	return nil
}

func SaveVer() error {
	cfg := ini.Empty()
	ini.ReflectFrom(cfg, V)
	return cfg.SaveTo(WorkPath + "/conf/ver.conf")
}

func IsDebug() bool {
	return (Conf.Mode == "Dev")
}

func BuildLogger() {

	if Conf.Type == "server" {
		/*log.AddHook(clientv1.NewHook(clientv1.LOGAPI_Config{
			Urls:      Conf.LogServer,
			TryCount:  3,
			PoolSize:  8,
			AppSecret: Conf.LogAppKey,
			AppId:     Conf.LogAppId,
			Caller:    true,
		})) // 提交到远端*/

		clientv1.SetupHook(clientv1.LOGAPI_Config{
			Urls:      Conf.LogServer,
			TryCount:  3,
			PoolSize:  6,
			AppSecret: Conf.LogAppKey,
			AppId:     Conf.LogAppId,
			Caller:    true,
		})
	}

	log.SetFormatter(&prefixed.TextFormatter{
		ForceFormatting:  true,
		ForceColors:      false,
		FullTimestamp:    true,
		DisableTimestamp: false,
		TimestampFormat:  "2006-01-02 15:04:05",
	})

	if Conf.Type == "memory" {
		log.SetOutput(os.Stdout)
	} else {
		os.MkdirAll(path.Dir(Conf.LogPath), os.ModePerm)
		writer, err := rotatelogs.New(
			Conf.LogPath+".%Y%m%d%H%M",
			rotatelogs.WithLinkName(Conf.LogPath),
			rotatelogs.WithMaxAge(time.Duration(60)*time.Minute),
			rotatelogs.WithRotationTime(6*time.Hour),
		)
		if err != nil {
			fmt.Println(err)
			return
		}

		debugWriter, derr := rotatelogs.New(
			Conf.LogPath+".debug.%Y%m%d%H%M",
			rotatelogs.WithLinkName(Conf.LogPath),
			rotatelogs.WithMaxAge(time.Duration(60)*time.Minute),
			rotatelogs.WithRotationTime(6*time.Hour),
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
	}

	log.SetLevel(log.Level(Conf.Level))
}

func CmpVal(isVal bool, tv interface{}, fv interface{}) interface{} {
	if isVal {
		return tv
	}

	return fv
}

func CmpStr(isVal bool, tv string, fv string) string {
	return CmpVal(isVal, tv, fv).(string)
}

func CmpBool(isVal bool, tv bool, fv bool) bool {
	return CmpVal(isVal, tv, fv).(bool)
}
