package remaps

import (
	"container/list"
	"crypto/rc4"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
	"sync"
	"encoding/json"

	guard "grapeGuard"
)


const (
	firstLine = "# Auto Make By Grape Guard X & Don't Modify.(本文件由grapeGuardX自动生成，请勿手动修改。)"
)

const (
	TypeRemapFile = iota
	TypeCacheFile
)

func cryptParam(val interface{}) (string, error) {
	jsonBody, err := json.Marshal(val)
	if err != nil {
		return "", err
	}

	rc4crypt, _ := rc4.NewCipher([]byte(pluginKey))
	newData := make([]byte, len(jsonBody))

	rc4crypt.XORKeyStream(newData, jsonBody)

	return base64.StdEncoding.EncodeToString(newData), nil
}

func decryptParam(src string, val interface{}) error {
	enBody, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		return err
	}

	rc4crypt, _ := rc4.NewCipher([]byte(pluginKey))
	newData := make([]byte, len(enBody))
	rc4crypt.XORKeyStream(newData, enBody)

	return json.Unmarshal(newData, val)
}

type LineNode struct {
	// 数据
	Data string
	Key  string
	// 解析的数据
	Values []string
	// 是否是注释
	IsComment bool // 是否是注释
	// 是否为空行
	IsEmpty bool
}

func (l *LineNode) Append(s string) {
	l.Values = append(l.Values, s)
}

func (l *LineNode) Set(sv string) {
	l.Values = strings.Split(sv, " ")
}

func (l *LineNode) String() string {
	if l.IsEmpty {
		return ""
	}

	if l.IsComment {
		return l.Data
	}

	return  strings.TrimSpace(strings.Join(append([]string{l.Key}, l.Values...), " "))
}

type TrafficFile struct {
	Filename string
	ForceFile bool

	FileType int

	configLines *list.List
	locker sync.RWMutex
}

func (f *TrafficFile) LoadFile() error {
	f.locker.Lock()
	defer f.locker.Unlock()

	cDat, err := ioutil.ReadFile(f.Filename)
	if err != nil {
		return err
	}

	f.configLines = list.New() // 创建文件队列
	lines := strings.Split(string(cDat), "\n")

	for _, v := range lines {
		bData := []byte(strings.TrimSpace(v))
		item := &LineNode{
			Data:      v,
			Values:    []string{},
			IsEmpty:   false,
			IsComment: false,
		}

		if len(v) == 0 {
			item.IsEmpty = true
			f.configLines.PushBack(item)
			continue
		}

		if len(bData) > 0 && bData[0] == '#' {
			item.IsComment = true
			item.Key = "comment"
		}

		if item.IsComment == false && f.FileType == TypeRemapFile {
			if f.FileType == TypeRemapFile {
				kv := strings.Fields(v)
				if len(kv) > 1 {
					item.Key = kv[0]
					item.Values = kv[1:]
				}
			}else{
				item.Values = strings.Fields(v)
			}
		}

		f.configLines.PushBack(item)
	}

	return nil
}

func (f *TrafficFile) SaveFile(name string) error {

	f.locker.RLock()
	defer f.locker.RUnlock()

	var tempDoc []string = []string{}

	fe := f.configLines.Front()
	if fe != nil {
		fitem := fe.Value.(*LineNode)
		if fitem.IsComment {
			fitem.Data = firstLine
		} else {
			f.configLines.PushFront(&LineNode{
				Data:      firstLine,
				Key:       "comment",
				Values:    []string{},
				IsEmpty:   false,
				IsComment: true,
			})
		}
	}

	for e := f.configLines.Front(); e != nil; e = e.Next() {
		item := e.Value.(*LineNode)
		tempDoc = append(tempDoc, item.String())
	}

	return ioutil.WriteFile(name, []byte(strings.Join(tempDoc, "\n")), 0666)
}

// 删除所有有效行
func (f *TrafficFile) ClearVaild() {
	f.locker.Lock()
	defer f.locker.Unlock()

	if f.ForceFile {
		for e := f.configLines.Front(); e != nil; {
			item := e.Value.(*LineNode)
			if item.IsComment || item.IsEmpty {
				e = e.Next()
				continue
			}

			remove := e
			e = e.Next()
			f.configLines.Remove(remove)
		}
	}
}

func (f *TrafficFile) ReplaceOrAdd(rtype,key string,newItem *LineNode) {
	f.locker.Lock()
	defer f.locker.Unlock()

	for e := f.configLines.Front(); e != nil; e = e.Next() {
		item := e.Value.(*LineNode)
		if item.IsComment || item.IsEmpty {
			continue
		}

		if len(item.Values) > 0 && strings.Contains(item.Values[0], key) && rtype == item.Key {
			// 删除这个节点
			f.configLines.Remove(e)
			break
		}
	}

	f.configLines.PushBack(newItem)
}

func (f *TrafficFile) Append(newItem *LineNode) {
	f.locker.Lock()
	defer f.locker.Unlock()

	f.configLines.PushBack(newItem)
}

func (f *TrafficFile) Sync() error {
	// 运行重新载入remap
	cmd := exec.Command(guard.Conf.RemapReload, strings.Fields("config reload")...)
	if cmd != nil {
		return cmd.Run()
	}

	return nil
}

func (f *TrafficFile) Print() {
	f.locker.RLock()
	defer f.locker.RUnlock()

	for e := f.configLines.Front(); e != nil; e = e.Next() {
		item := e.Value.(*LineNode)
		fmt.Println(*item)
	}
}

func NewConfig(ftype int,filename string) *TrafficFile {
	return &TrafficFile{
		Filename:filename,
		ForceFile:true,
		FileType:ftype,
	}
}

