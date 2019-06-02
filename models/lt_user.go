package models

type LtUser struct {
	Address    string `xorm:"not null default '' comment('联系地址') VARCHAR(255)"`
	Blacktime  int    `xorm:"not null default 0 comment('黑名单限制到期时间') INT(10)"`
	Id         int    `xorm:"not null pk autoincr INT(10)"`
	Mobile     string `xorm:"not null default '' comment('手机号') VARCHAR(50)"`
	Realname   string `xorm:"not null default '' comment('联系人') VARCHAR(50)"`
	SysCreated int    `xorm:"not null default 0 comment('创建时间') INT(10)"`
	SysIp      string `xorm:"not null default '' comment('IP地址') VARCHAR(50)"`
	SysUpdated int    `xorm:"not null default 0 comment('修改时间') INT(10)"`
	Username   string `xorm:"not null default '' comment('用户名') VARCHAR(50)"`
}
