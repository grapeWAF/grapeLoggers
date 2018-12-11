package clientv1

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/sirupsen/logrus"
)

func Test_Fire(t *testing.T) {
	logrus.AddHook(NewHook(LOGAPI_Config{
		Urls:      []string{"localhost:2996"},
		TryCount:  3,
		PoolSize:  4,
		AppSecret: "77439d6adc102edde14f1dbab93205d50ec05263",
		AppId:     "26483",
		Caller:    true,
	}))

	logrus.Info("Testabasdasd")
}

/*func Benchmark_Fire(b *testing.B) {
	logrus.AddHook(NewHook(LOGAPI_Config{
		Urls:      []string{"localhost:2996"},
		TryCount:  3,
		PoolSize:  4,
		AppSecret: "77439d6adc102edde14f1dbab93205d50ec05263",
		AppId:     "26483",
		Caller:    false,
	}))

	for i := 0; i < b.N; i++ {
		logrus.Info("测试日志INFO型...")
		//.Error("测试日志ERROR")
	}
}*/

func Benchmark_PFire(b *testing.B) {
	logrus.AddHook(NewHook(LOGAPI_Config{
		Urls:      []string{"localhost:2996"},
		TryCount:  3,
		PoolSize:  4,
		AppSecret: "77439d6adc102edde14f1dbab93205d50ec05263",
		AppId:     "26483",
		Caller:    true,
	}))

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logrus.Info("测试日志INFO型ssssssssssssssssssssssssssssss77439d6adc102edde14f1dbab93205d50ec0526377439d6adc102edde14f1dbab93205d50ec05263...")
			logrus.Warn("测试日志INFO型ssssssssssssssssssssssssssssss77439d6adc102edde14f1dbab93205d50ec0526377439d6adc102edde14f1dbab93205d50ec05263...")
			logrus.Debug("测试日志INFO型ssssssssssssssssssssssssssssss77439d6adc102edde14f1dbab93205d50ec0526377439d6adc102edde14f1dbab93205d50ec05263...")
			logrus.Error("测试日志ERROR")
		}
	})
}

func Benchmark_AllFire(b *testing.B) {
	logrus.AddHook(NewHook(LOGAPI_Config{
		Urls:      []string{"localhost:2996"},
		TryCount:  3,
		PoolSize:  4,
		AppSecret: "77439d6adc102edde14f1dbab93205d50ec05263",
		AppId:     "26483",
		Caller:    false,
	}))

	for i := 0; i < b.N; i++ {
		logrus.Info("测试日志INFO型ssssssssssssssssssssssssssssss77439d6adc102edde14f1dbab93205d50ec0526377439d6adc102edde14f1dbab93205d50ec05263...")
		logrus.Warn("测试日志INFO型ssssssssssssssssssssssssssssss77439d6adc102edde14f1dbab93205d50ec0526377439d6adc102edde14f1dbab93205d50ec05263...")
		logrus.Debug("测试日志INFO型ssssssssssssssssssssssssssssss77439d6adc102edde14f1dbab93205d50ec0526377439d6adc102edde14f1dbab93205d50ec05263...")
		logrus.Error("测试日志ERROR77439d6adc102edde14f1dbab93205d50ec0526377439d6adc102edde14f1dbab93205d50ec0526377439d6adc102edde14f1dbab93205d50ec0526377439d6adc102edde14f1dbab93205d50ec05263")
	}
}

func Test_JsonMsg(t *testing.T) {
	newMsg := &LogMsg{
		Version: "1.0",
		Host:    "test123123",
		Level:   1,
		Type:    "LOGS",
		Msg:     "TEST LOGS",
		Time:    time.Now().Unix(),
		Caller:  "uncall.go:123",
		Extra: map[string]interface{}{
			"caller": "uncall.go:123",
		},
	}

	msgJson := newMsg.Json()
	fmt.Print(msgJson)

	var JsonMap map[string]interface{} = map[string]interface{}{}
	json.Unmarshal([]byte(msgJson), &JsonMap)

	dmpMsg := &LogMsg{}
	dmpMsg.ParserMap(JsonMap)

	fmt.Print(dmpMsg)
}
