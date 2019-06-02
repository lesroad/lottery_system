/**
 * 同一个IP抽奖，每天的操作限制，本地或者redis缓存
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

const ipFrameSize = 2 //将大的散列表分成两段小的

func init() {
	//// IP当天的统计数，整点归零，设置定时器
	//duration := comm.NextDayDuration()
	//time.AfterFunc(duration, resetGroupIpList)

	// 本地开发测试的时候，每次启动归零
	resetGroupIpList()
}

// 重置单机IP今天次数
func resetGroupIpList() {
	log.Println("重置IP抽奖次数")
	cacheObj := datasource.InstanceCache()
	for i := 0; i < ipFrameSize; i++ {
		key := fmt.Sprintf("day_ips_%d", i)
		cacheObj.Do("DEL", key)
	}
	
	// IP当天的统计数，整点归零，设置定时器
	duration := comm.NextDayDuration()
	time.AfterFunc(duration, resetGroupIpList)
}

// 今天的IP抽奖次数递增，返回递增后的数值
func IncrIpLuckyNum(strIp string) int64 {
	ip := comm.Ip4toInt(strIp)
	i := ip % ipFrameSize
	// 集群的redis统计数递增
	return incrServIpLucyNum(i, ip)
}

func incrServIpLucyNum(i, ip int64) int64 {
	key := fmt.Sprintf("day_ips_%d", i)
	cacheObj := datasource.InstanceCache()
	rs, err := cacheObj.Do("HINCRBY", key, ip, 1) //key(ip)+1
	if err != nil {
		log.Println("ip_day_lucky redis HINCRBY err=", err)
		return math.MaxInt32
	} else {
		return rs.(int64)
	}
}
