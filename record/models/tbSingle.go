package model

type TbSingle struct {
	Idx        int64  `xorm:"not null pk autoincr BIGINT(255)"`
	Mobile     int    `xorm:"default NULL INT(255)"`
	Pcdevice   int    `xorm:"default NULL INT(255)"`
	Hostname   string `xorm:"default 'NULL' index(host) VARCHAR(255)"`
	Uuid       string `xorm:"default 'NULL' index(host) VARCHAR(255)"`
	Hostaddr   string `xorm:"default 'NULL' VARCHAR(255)"`
	Guardcount int    `xorm:"default NULL INT(255)"`
	Hitcount   int    `xorm:"default NULL INT(255)"`
	Misscount  int    `xorm:"default NULL INT(255)"`
	Nocache    int    `xorm:"default NULL INT(255)"`
	Rq         int    `xorm:"default NULL INT(255)"`
	Paddevice  int    `xorm:"default NULL INT(255)"`
	Timestamp  int64  `xorm:"default NULL index BIGINT(20)"`
}
