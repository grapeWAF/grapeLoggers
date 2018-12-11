package model

import (
	"time"
)

type TbLogs struct {
	Logidx    int64     `xorm:"not null pk autoincr BIGINT(20)" json:"-"`
	Message   string    `xorm:"not null TEXT" json:"message"`
	Time      time.Time `xorm:"not null DATETIME"  json:"logTime"`
	Timestmap int64     `xorm:"not null index BIGINT(20)" json:"-"`
	Level     string    `xorm:"not null VARCHAR(255)" json:"level"`
	Type      string    `xorm:"not null VARCHAR(255)" json:"type"`
	Host      string    `xorm:"default 'NULL' index VARCHAR(255)" json:"hostName"`
	Serverip  string    `xorm:"default 'NULL' VARCHAR(255)" json:"remoteIp"`
	Datajson  string    `xorm:"default 'NULL' TEXT" json:"-"`
	Caller    string    `xorm:"default 'NULL' VARCHAR(255)" json:"caller"`
}
