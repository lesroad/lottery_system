package datasource

import (
	"fmt"
	"iris项目/my_lottery/conf"
	"log"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

var masterInstance *xorm.Engine
var dbLock sync.Mutex

// 创建单例
func InstanceDbMaster() *xorm.Engine {
	if masterInstance != nil {
		return masterInstance
	}
	dbLock.Lock()
	defer dbLock.Unlock()
	// 再次检查
	if masterInstance != nil {
		return masterInstance
	}
	return NewDbMaster()
}

//实例化xorm引擎
func NewDbMaster() *xorm.Engine {
	sourcename := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8",
		conf.DbMaster.User,
		conf.DbMaster.Pwd,
		conf.DbMaster.Host,
		conf.DbMaster.Port,
		conf.DbMaster.Database)
	instance, err := xorm.NewEngine(conf.DriverName, sourcename)
	if err != nil {
		log.Fatal("dbhelper.NewDbMaster: NewEngine error", err)
	}
	//展示每一条sql语句，调试用
	// instance.ShowSQL(true)
	instance.ShowSQL(false)
	masterInstance = instance
	return instance
}
