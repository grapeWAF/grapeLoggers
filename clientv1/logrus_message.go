package clientv1

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"go.uber.org/zap/buffer"

	proto "grapeLoggers/protos"
)

////////////////////////////////////////////////////////////
//  add by koangel
//	date: 2017/09/20
//  日志记录，改为标准JSON协议，并增加GZIP等行为
////////////////////////////////////////////////////////////

type LogMsg struct {
	Version string                 `form:"ver" json:"ver"`
	Host    string                 `form:"host" json:"host"`
	Msg     string                 `form:"msg" json:"msg"`
	Time    int64                  `form:"timestamp" json:"timestamp"`
	Level   int32                  `form:"level" json:"level"`
	Type    string                 `form:"type" json:"type"`
	Caller  string                 `form:"caller" json:"caller"`
	Extra   map[string]interface{} `form:"-" json:"-"`
}

func (m *LogMsg) LevelStr() string {
	return log.Level(m.Level).String()
}

func (m *LogMsg) Json() string {
	jv, _ := m.I2JSON()
	return fmt.Sprintf(`{"msg":"%s", "g":false }`, jv)
}

func (m *LogMsg) I2JSON() (bv []byte, err error) {
	var b, eb []byte
	err = nil
	bv = nil

	extra := m.Extra
	b, err = json.Marshal(m)
	m.Extra = extra
	if err != nil {
		return
	}

	if len(extra) == 0 {
		bv = []byte(base64.StdEncoding.EncodeToString(b))
		return
	}

	if eb, err = json.Marshal(extra); err != nil {
		return
	}

	// merge serialized message + serialized extra map
	b[len(b)-1] = ','
	bv = []byte(base64.StdEncoding.EncodeToString(append(b, eb[1:len(eb)]...)))
	return
}

func (m *LogMsg) GZip() string {
	zipB64, _ := m.I2Gzip()
	return fmt.Sprintf(`{"msg":"%s", "g":true }`, zipB64)
}

func (m *LogMsg) I2Gzip() (bz string, err error) {
	err = nil
	bz = ""
	JsonByte, jerr := m.I2JSON()
	if jerr != nil {
		err = jerr
		return
	}

	var ziper buffer.Buffer
	zw := gzip.NewWriter(&ziper)

	zw.Write(JsonByte)
	zw.Close()

	bz = base64.StdEncoding.EncodeToString(ziper.Bytes())
	return
}

func (m *LogMsg) ParserMap(msg map[string]interface{}) error {
	msgData, ok := msg["msg"]
	if !ok {
		return errors.New("msg data lost...")
	}

	isGzip, ok := msg["g"]
	if !ok {
		return errors.New("msg data lost[g]...")
	}

	return m.ParserMsg(msgData.(string), isGzip.(bool))
}

func (m *LogMsg) ParserMsg(msg string, isGizp bool) error {

	decodeByte, err := base64.StdEncoding.DecodeString(msg)
	if err != nil {
		return err
	}

	if isGizp {
		zr, gerr := gzip.NewReader(bytes.NewReader(decodeByte))
		if gerr != nil {
			return gerr
		}

		unzipByte, uerr := ioutil.ReadAll(zr)
		zr.Close()
		if uerr != nil {
			return uerr
		}

		return m.JSON2I(unzipByte)
	}

	return m.JSON2I(decodeByte)
}

func (m *LogMsg) JSON2I(data []byte) error {
	i := make(map[string]interface{}, 16)
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	for k, v := range i {
		if k[0] == '_' {
			if m.Extra == nil {
				m.Extra = make(map[string]interface{}, 1)
			}
			m.Extra[k] = v
			continue
		}
		switch k {
		case "ver":
			m.Version = v.(string)
		case "host":
			m.Host = v.(string)
		case "msg":
			m.Msg = v.(string)
		case "timestamp":
			m.Time = int64(v.(float64))
		case "type":
			m.Type = v.(string)
		case "level":
			m.Level = int32(v.(float64))
		case "caller":
			m.Caller = v.(string)
		}
	}
	return nil
}

func GetLoggerPotos(url string) (cli proto.LogServiceClient,err error) {
	cli = nil
	err = nil

	ccli, cerr := grpc.Dial(url, grpc.WithInsecure())
	if cerr != nil {
		WithField("type", "collect").Error("连接采集服务器失败:", cerr)
		err = cerr
		return
	}

	cli = proto.NewLogServiceClient(ccli)
	return
}