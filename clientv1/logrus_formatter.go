package clientv1

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

////////////////////////////////////////////////////////////
//  add by koangel
//	date: 2017/09/20
//  日志记录，格式方式代码，让格式符合本少爷的模式
////////////////////////////////////////////////////////////

const (
	nocolor = 0
	red     = 31
	green   = 32
	yellow  = 33
	blue    = 36
	gray    = 37
)

var (
	baseTimestamp time.Time
)

func init() {
	baseTimestamp = time.Now()
}

type RTextFormatter struct {
}

func (f *RTextFormatter) Format(entry *log.Entry) ([]byte, error) {
	var b *bytes.Buffer
	keys := make([]string, 0, len(entry.Data))
	for k := range entry.Data {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	levelText := strings.ToUpper(entry.Level.String())[0:4]

	fmt.Fprintf(b, "[%s][%s]:[%s]", levelText, entry.Time.Format("2006-01-02 15:04:05"), entry.Message)

	if len(keys) > 0 {
		fmt.Fprintf(b, " data:[")
		isFirst := true
		for _, k := range keys {
			v := entry.Data[k]
			if isFirst == false {
				b.WriteByte(' ')
			}
			isFirst = false
			fmt.Fprintf(b, "%s=", k)
			f.appendValue(b, v)
		}
		fmt.Fprintf(b, "]")
	}

	b.WriteByte('\n')
	return b.Bytes(), nil
}

func (f *RTextFormatter) appendKeyValue(b *bytes.Buffer, key string, value interface{}) {
	if b.Len() > 0 {
		b.WriteByte(' ')
	}
	b.WriteString(key)
	b.WriteByte('=')
	f.appendValue(b, value)
}

func (f *RTextFormatter) needsQuoting(text string) bool {
	for _, ch := range text {
		if !((ch >= 'a' && ch <= 'z') ||
			(ch >= 'A' && ch <= 'Z') ||
			(ch >= '0' && ch <= '9') ||
			ch == '-' || ch == '.' || ch == '_' || ch == '/' || ch == '@' || ch == '^' || ch == '+') {
			return true
		}
	}
	return false
}

func (f *RTextFormatter) appendValue(b *bytes.Buffer, value interface{}) {
	stringVal, ok := value.(string)
	if !ok {
		stringVal = fmt.Sprint(value)
	}

	if !f.needsQuoting(stringVal) {
		b.WriteString(stringVal)
	} else {
		b.WriteString(fmt.Sprintf("%q", stringVal))
	}
}
