package model

import "time"

type TbRemapdata struct {
	Logidx    int64 `xorm:"not null pk autoincr BIGINT(20)"`
	Formsrc   string
	Recvbytes int32
	Sendbytes int32
	Reqtotal  int32
	Timestamp int64
	Time      time.Time `xorm:"not null DATETIME"`
}
