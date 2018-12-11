package model

import (
	"time"
)

type TbQps struct {
	Qpsidx    int64     `xorm:"not null pk autoincr BIGINT(255)"`
	Rqcount   int       `xorm:"not null INT(255)"`
	Pv        int       `xorm:"not null INT(255)"`
	Time      time.Time `xorm:"not null DATETIME"`
	Timestamp int64     `xorm:"default NULL index BIGINT(20)"`
}
