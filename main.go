package main

////////////////////////////////////////////////////////////
//  add by koangel
//	date: 2017/09/19
//  日志记录协议可以远端查询
////////////////////////////////////////////////////////////

import (
	"fmt"
	"grapeLoggers/record"
	"net"
	"os"
	"runtime"
	"time"

	"github.com/koangel/grapeTimer"

	"google.golang.org/grpc"

	"grapeLoggers/appConf"

	"github.com/fatih/color"
	log "github.com/sirupsen/logrus"

	proto "grapeLoggers/protos"
	"grapeLoggers/routers"

	util "github.com/koangel/grapeNet/Utils"
)

var (
	version = "0.4.5 beta"
)

func runLoop(signal chan os.Signal) string {
	runtime.GOMAXPROCS((runtime.NumCPU() * 4) + 2)

	color.Blue(" __                                              ")
	color.Blue("/\\ \\                                             ")
	color.Blue("\\ \\ \\        ___      __      __      __   _ __  ")
	color.Blue(" \\ \\ \\  __  / __`\\  /'_ `\\  /'_ `\\  /'__`\\/\\`'__\\")
	color.Blue("  \\ \\ \\L\\ \\/\\ \\L\\ \\/\\ \\L\\ \\/\\ \\L\\ \\/\\  __/\\ \\ \\/ ")
	color.Blue("   \\ \\____/\\ \\____/\\ \\____ \\ \\____ \\ \\____\\\\ \\_\\ ")
	color.Blue("    \\/___/  \\/___/  \\/___L\\ \\/___L\\ \\/____/ \\/_/ ")
	color.Blue("                      /\\____/ /\\____/            ")
	color.Blue("                      \\_/__/  \\_/__/             ")
	color.White(`-----------------------------------------------------`)
	color.White("*	Grape Loggers Server %v", version)
	color.White("*	Support Mysql & InfluxDB & MongoDB & Postgres")
	color.White("*	Log Collection, Powered by Grape Soft")
	color.White(`-----------------------------------------------------`)

	grapeTimer.CDebugMode = false
	grapeTimer.UseAsyncExec = true
	grapeTimer.InitGrapeScheduler(time.Second, false)

	lerr := config.LoadConf()
	if lerr != nil {
		log.Error(lerr)
		return "load config error..."
	}

	config.BuildLogger()

	log.Info("Grape Logger Server ", version)

	lis, err := net.Listen("tcp", config.C.HttpAddr)
	if err != nil {
		log.Errorf("failed to listen: %v", err)
		return "error"
	}

	grpcServer := grpc.NewServer()
	proto.RegisterLogServiceServer(grpcServer, &routers.RpcLoggersV1{})

	log.Info("初始化Records...")
	if record.InitRecords() == false {
		return "error"
	}

	log.Info("Listen:", config.C.HttpAddr)
	go grpcServer.Serve(lis)

	select {
	case killsign := <-signal:
		log.Info("logger stoping...")
		grpcServer.Stop()
		if killsign == os.Interrupt {
			return "Daemon was interruped by system signal"
		}

		return "logger stoped.."
	}

	return "finished"
}

func main() {
	fmt.Println(util.RunDaemon("grapeLogger", "Logs System", "", runLoop))
}
