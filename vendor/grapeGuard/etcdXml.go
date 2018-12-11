package grapeGuard

////////////////////////////////////////////////////////////
//  add by koangel
//	date: 2017/09/27
//  只有基本函数，将配置文件（大量文本，压缩后转为BASE64存取）
////////////////////////////////////////////////////////////

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"io/ioutil"

	"go.uber.org/zap/buffer"

	log "grapeLoggers/clientv1"
)

type XmlGzipFormatter struct {
}

func (f *XmlGzipFormatter) Marshal(v interface{}) ([]byte, error) {
	return Marshal(v)
}

func (f *XmlGzipFormatter) Unmarshal(data []byte, v interface{}) error {
	return Unmarshal(data, v)
}

func (f *XmlGzipFormatter) ToString(convert []byte) string {
	return base64.StdEncoding.EncodeToString(convert)
}

func (f *XmlGzipFormatter) FromString(convert string) []byte {
	b, _ := base64.StdEncoding.DecodeString(convert)
	return b
}

func Marshal(v interface{}) (body []byte, err error) {
	err = nil
	body = nil
	// step 1 XML序列化
	b, xerr := xml.Marshal(v)
	if xerr != nil {
		err = xerr
		log.Error("Etcd Xml Error:", xerr)
		return
	}

	// step 2压缩Gzip
	var zipBuf buffer.Buffer
	zap := gzip.NewWriter(&zipBuf)

	zap.Write(b)
	zap.Close()

	body = zipBuf.Bytes()
	return
}

func MarshalJson(v interface{}) (body []byte, err error) {
	err = nil
	body = nil
	// step 1 json序列化
	b, xerr := json.Marshal(v)
	if xerr != nil {
		err = xerr
		log.Error("Etcd Json Error:", xerr)
		return
	}

	// step 2压缩Gzip
	var zipBuf buffer.Buffer
	zap := gzip.NewWriter(&zipBuf)

	zap.Write(b)
	zap.Close()

	body = zipBuf.Bytes()
	return
}

func MarshalGzip(v interface{}) (body string, err error) {
	err = nil
	body = ""

	b, cerr := Marshal(v)
	if cerr != nil {
		err = cerr
		return
	}

	// step 3 BASE64
	body = base64.StdEncoding.EncodeToString(b)
	return
}

func MarshalJsonGzip(v interface{}) (body string, err error) {
	err = nil
	body = ""

	b, cerr := MarshalJson(v)
	if cerr != nil {
		err = cerr
		return
	}

	// step 3 BASE64
	body = base64.StdEncoding.EncodeToString(b)
	return
}

func Unmarshal(data []byte, v interface{}) error {
	// 解压缩
	zr, gerr := gzip.NewReader(bytes.NewReader(data))
	if gerr != nil {
		return gerr
	}

	unzipByte, uerr := ioutil.ReadAll(zr)
	zr.Close()
	if uerr != nil {
		return uerr
	}

	// xml反向序列化
	return xml.Unmarshal(unzipByte, v)
}

func UnmarshalJson(data []byte, v interface{}) error {
	// 解压缩
	zr, gerr := gzip.NewReader(bytes.NewReader(data))
	if gerr != nil {
		return gerr
	}

	unzipByte, uerr := ioutil.ReadAll(zr)
	zr.Close()
	if uerr != nil {
		return uerr
	}

	// json反向序列化
	return json.Unmarshal(unzipByte, v)
}

func UnmarshalGzip(data string, v interface{}) error {
	// 反向base64
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Error("Etcd Xml Error:", err)
		return err
	}

	return Unmarshal(b, v)
}

func UnmarshalJsonGzip(data string, v interface{}) error {
	// 反向base64
	b, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		log.Error("Etcd Xml Error:", err)
		return err
	}

	return UnmarshalJson(b, v)
}
