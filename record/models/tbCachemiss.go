package model

type TbCachemiss struct {
	Idx        int64  `xorm:"not null pk autoincr BIGINT(255)"`
	Host       string `xorm:"default 'NULL' VARCHAR(255)"`
	Uuid       string `xorm:"default 'NULL' VARCHAR(255)"`
	Url        string `xorm:"not null VARCHAR(255)"`
	Path       string `xorm:"default 'NULL' VARCHAR(255)"`
	Via        string `xorm:"default 'NULL' VARCHAR(16)"`
	Statuscode int    `xorm:"default NULL INT(255)"`
	Timestamp  int64  `xorm:"default NULL index BIGINT(20)"`
}
