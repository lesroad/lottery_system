package datasource

import (
	"fmt"
	"iris项目/my_lottery/conf"
	"log"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
)

var rdsLock sync.Mutex
var cacheInstance *RedisConn

type RedisConn struct {
	pool      *redis.Pool
	showDebug bool
}

func (rds *RedisConn) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	conn := rds.pool.Get()
	defer conn.Close()

	t1 := time.Now().UnixNano() //取出当前时间ns数
	reply, err = conn.Do(commandName, args...)
	if err != nil {
		e := conn.Err()
		if e != nil {
			log.Println("rdshelper.Do", err, e)
		}
	}
	t2 := time.Now().UnixNano()
	if rds.showDebug {
		fmt.Printf("[redis] [info] [%dus] cmd=%s, err=%s, args=%v, reply=%s\n",
			(t2-t1)/1000, commandName, err, args, reply)
	}
	return reply, err
}

func (rds *RedisConn) ShowDebug(b bool) {
	rds.showDebug = b
}

//单例化
func InstanceCache() *RedisConn {
	if cacheInstance != nil {
		return cacheInstance
	}
	rdsLock.Lock()
	defer rdsLock.Unlock()
	// 再次检查
	if cacheInstance != nil {
		return cacheInstance
	}
	return NewCache()
}

func NewCache() *RedisConn {
	pool := redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d",
				conf.RdsCache.Host, conf.RdsCache.Port))
			if err != nil {
				log.Fatal("rdshelper.NewCache Dial error=", err)
				return nil, err
			}
			return c, nil
		},
		//是否在从池中取出连接前进行检验,如果检验失败,则从池中去除连接并尝试取出另一个
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
		MaxIdle:         10000, //最多保持空闲连接数
		MaxActive:       10000, //最大活跃数
		IdleTimeout:     0,
		Wait:            false,
		MaxConnLifetime: 0,
	}
	instance := &RedisConn{
		pool: &pool,
	}
	cacheInstance = instance
	// instance.ShowDebug(true)
	instance.ShowDebug(false)
	return instance
}
