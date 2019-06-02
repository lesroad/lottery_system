/**
 * 同一个User抽奖，每天的操作限制，本地或者redis缓存
 */
package utils

import (
	"fmt"
	"log"
	"math"
	"time"

	"iris项目/my_lottery/comm"
	"iris项目/my_lottery/datasource"
)

const userFrameSize = 2 //将hash散列为多段数据，让每个hash小点

func init() {
	// User当天的统计数，整点归零，设置定时器
	//duration := comm.NextDayDuration()
	//time.AfterFunc(duration, resetGroupUserList)

	// TODO: 本地开发测试的时候，每次启动归零
	resetGroupUserList()
}

// 集群模式，重置用户今天次数
func resetGroupUserList() {
	log.Println("重置今日用户抽奖次数")
	cacheObj := datasource.InstanceCache()
	for i := 0; i < userFrameSize; i++ {
		key := fmt.Sprintf("day_users_%d", i)
		cacheObj.Do("DEL", key)
	}

	// IP当天的统计数，整点归零，设置定时器
	duration := comm.NextDayDuration()
	time.AfterFunc(duration, resetGroupUserList) //等待时间段d过去，然后调用func
}

// 今天的用户抽奖次数递增，返回递增后的数值
func IncrUserLuckyNum(uid int) int64 {
	i := uid % userFrameSize
	// 集群的redis统计数递增
	key := fmt.Sprintf("day_users_%d", i) //这里是两个哈希表
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("HINCRBY", key, uid, 1) //返回键为uid的值
	if err != nil {
		log.Println("user_day_lucky redis HINCRBY key=", key,
			", uid=", uid, ", err=", err)
		return math.MaxInt32
	} else {
		num := rs.(int64)
		return num
	}
}

// 从给定的数据直接初始化用户的参与次数
func InitUserLuckyNum(uid int, num int64) {
	if num <= 1 {
		return
	}
	i := uid % userFrameSize
	key := fmt.Sprintf("day_users_%d", i)
	cacheObj := datasource.InstanceCache()
	_, err := cacheObj.Do("HSET", key, uid, num)
	if err != nil {
		log.Println("user_day_lucky redis HSET key=", key,
			", uid=", uid, ", err=", err)
	}
}
