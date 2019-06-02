package models

type LtGift struct {
	Displayorder int    `xorm:"not null default 0 comment('位置序号，小的排在前面') INT(10)"`
	Gdata        string `xorm:"not null default '' comment('扩展数据，如：虚拟币数量') VARCHAR(255)"`
	Gtype        int    `xorm:"not null default 0 comment('奖品类型，0 虚拟币，1 虚拟券，2 实物-小奖，3 实物-大奖') INT(10)"`
	Id           int    `xorm:"not null pk autoincr INT(10)"`
	Img          string `xorm:"not null default '' comment('奖品图片') VARCHAR(255)"`
	LeftNum      int    `xorm:"not null default 0 comment('剩余数量') INT(11)"`
	PrizeBegin   int    `xorm:"not null default 0 comment('发奖计划周期的开始') INT(11)"`
	PrizeCode    string `xorm:"not null default '' comment('0-9999表示100%，0-0表示万分之一的中奖概率') VARCHAR(50)"`
	PrizeData    string `xorm:"comment('发奖计划，[[时间1,数量1],[时间2,数量2]]') MEDIUMTEXT"`
	PrizeEnd     int    `xorm:"not null default 0 comment('发奖计划周期的结束') INT(11)"`
	PrizeNum     int    `xorm:"not null default -1 comment('奖品数量，0 无限量，>0限量，<0无奖品') INT(11)"`
	PrizeTime    int    `xorm:"not null default 0 comment('发奖周期，D天') INT(10)"`
	SysCreated   int    `xorm:"not null default 0 comment('创建时间') INT(10)"`
	SysIp        string `xorm:"not null default '' comment('操作人IP') VARCHAR(50)"`
	SysStatus    int    `xorm:"not null default 0 comment('状态，0 正常，1 删除') SMALLINT(5)"`
	SysUpdated   int    `xorm:"not null default 0 comment('修改时间') INT(10)"`
	TimeBegin    int    `xorm:"not null default 0 comment('开始时间') INT(11)"`
	TimeEnd      int    `xorm:"not null default 0 comment('结束时间') INT(11)"`
	Title        string `xorm:"not null default '' comment('奖品名称') VARCHAR(255)"`
}
