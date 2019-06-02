package utils

import (
	"encoding/json"
	"log"
	"time"

	"iris项目/my_lottery/comm"
	"iris项目/my_lottery/conf"
	"iris项目/my_lottery/datasource"
	"iris项目/my_lottery/models"
	"iris项目/my_lottery/services"

	"github.com/garyburd/redigo/redis"
)

// func init() {
// 	// 本地开发测试的时候，每次重新启动，奖品池自动归零
// 	resetServGiftPool()
// }

// 重置一个奖品的发奖周期信息
// 奖品剩余数量也会重新设置为当前奖品数量
// 奖品的奖品池有效数量则会设置为空
// 奖品数量、发放周期等设置有修改的时候，也需要重置
// 【难点】根据发奖周期，重新更新发奖计划
func ResetGiftPrizeData(giftInfo *models.LtGift, giftService services.GiftService) {
	if giftInfo == nil || giftInfo.Id < 1 {
		return
	}
	id := giftInfo.Id
	nowTime := comm.NowUnix()
	// 不能发奖，不需要设置发奖周期
	if giftInfo.SysStatus == 1 || // 状态不对
		giftInfo.TimeBegin >= nowTime || // 开始时间不对
		giftInfo.TimeEnd <= nowTime || // 结束时间不对
		giftInfo.LeftNum <= 0 || // 剩余数不足
		giftInfo.PrizeNum <= 0 { // 总数不限制
		if giftInfo.PrizeData != "" { //发奖计划不为空则清空
			clearGiftPrizeData(giftInfo, giftService)
		}
		return
	}
	// 不限制发奖周期，直接把奖品数量全部设置到奖品池
	dayNum := giftInfo.PrizeTime
	if dayNum <= 0 {
		setGiftPool(id, giftInfo.LeftNum)
		return
	}

	// 对于设置发奖周期的奖品重新计算出来合适的奖品发放节奏
	// 奖品池的剩余数先设置为空
	setGiftPool(id, 0)

	// 每天的概率一样
	// 一天内24小时，每个小时的概率是不一样的
	// 一小时内60分钟的概率一样
	prizeNum := giftInfo.PrizeNum
	avgNum := prizeNum / dayNum

	// 每天可以分配到的奖品数量
	dayPrizeNum := make(map[int]int)
	// 平均分配，每天分到的奖品数量做分布
	if avgNum >= 1 && dayNum > 0 {
		for day := 0; day < dayNum; day++ {
			dayPrizeNum[day] = avgNum
		}
	}
	// 剩下的随机分配到任意哪天
	prizeNum -= dayNum * avgNum
	for prizeNum > 0 {
		prizeNum--
		day := comm.Random(dayNum)
		dayPrizeNum[day] += 1
	}
	// 每天的map，每小时的map，60分钟的数组，奖品数量
	prizeData := make(map[int]map[int][60]int)
	for day, num := range dayPrizeNum {
		//计算一天的发奖计划
		dayPrizeData := getGiftPrizeDataOneDay(num)
		prizeData[day] = dayPrizeData
	}
	// 将周期内每天、每小时、每分钟的数据 prizeData 格式化，再序列化保存到数据表,格式：([时间:数量])
	datalist := formatGiftPrizeData(nowTime, dayNum, prizeData)
	str, err := json.Marshal(datalist)
	if err != nil {
		log.Println("prizedata.ResetGiftPrizeData json error=", err)
	} else {
		// 保存奖品的分布计划数据
		info := &models.LtGift{
			Id:         giftInfo.Id,
			LeftNum:    giftInfo.PrizeNum,
			PrizeData:  string(str),
			PrizeBegin: nowTime,
			PrizeEnd:   nowTime + dayNum*86400,
			SysUpdated: nowTime,
		}
		err := giftService.Update(info, nil)
		if err != nil {
			log.Println("prizedata.ResetGiftPrizeData giftService.Update",
				info, ", error=", err)
		}
	}
}

// 清空奖品的发放计划
func clearGiftPrizeData(giftInfo *models.LtGift, giftService services.GiftService) {
	info := &models.LtGift{
		Id:        giftInfo.Id,
		PrizeData: "",
	}
	err := giftService.Update(info, []string{"prize_data"})
	if err != nil {
		log.Println("prizedata.clearGiftPrizeData giftService.Update",
			info, ", error=", err)
	}
	//奖品池也设为0
	setGiftPool(giftInfo.Id, 0)
}

// 设置奖品池的数量
func setGiftPool(id, num int) {
	key := "gift_pool"
	cacheObj := datasource.InstanceCache()
	_, err := cacheObj.Do("HSET", key, id, num)
	if err != nil {
		log.Println("prizedata.setServGiftPool error=", err)
	}
}

// 将给定的奖品数量分布到这一天的时间内
// 结构为： map[hour][minute]num
func getGiftPrizeDataOneDay(num int) map[int][60]int {
	rs := make(map[int][60]int)
	hourData := [24]int{}
	// 分别将奖品分布到24个小时内
	if num > 100 {
		// 奖品数量多的时候，直接按照百分比计算出来
		for _, h := range conf.PrizeDataRandomDayTime {
			hourData[h]++
		}
		for h := 0; h < 24; h++ {
			d := hourData[h]
			n := num * d / 100
			hourData[h] = n
			num -= n
		}
	}
	// 奖品数量少的时候，或者剩下了一些没有分配，需要用到随即概率来计算
	for num > 0 {
		num--
		// 通过随机数确定奖品落在哪个小时
		hourIndex := comm.Random(100)
		h := conf.PrizeDataRandomDayTime[hourIndex]
		hourData[h]++
	}
	// 将每个小时内的奖品数量分配到60分钟
	for h, hnum := range hourData {
		if hnum <= 0 {
			continue
		}
		minuteData := [60]int{}
		if hnum >= 60 {
			avgMinute := hnum / 60
			for i := 0; i < 60; i++ {
				minuteData[i] = avgMinute
			}
			hnum -= avgMinute * 60
		}
		// 剩下的数量不多的时候，随机到各分钟内
		for hnum > 0 {
			hnum--
			m := comm.Random(60)
			minuteData[m]++
		}
		rs[h] = minuteData
	}
	return rs
}

// 将prizeData格式化成具体到一个时间（分钟）的奖品数量
// 结构为： [day][hour][minute]num
// result: [][时间, 数量]
func formatGiftPrizeData(nowTime, dayNum int, prizeData map[int]map[int][60]int) [][2]int {
	rs := make([][2]int, 0)
	nowHour := time.Now().Hour()
	// 处理周期内每一天的计划
	for dn := 0; dn < dayNum; dn++ {
		dayData, ok := prizeData[dn]
		if !ok {
			continue
		}
		dayTime := nowTime + dn*86400
		// 处理周期内，每小时的计划
		for hn := 0; hn < 24; hn++ {
			hourData, ok := dayData[(hn+nowHour)%24]
			if !ok {
				continue
			}
			hourTime := dayTime + hn*3600
			// 处理周期内，每分钟的计划
			for mn := 0; mn < 60; mn++ {
				num := hourData[mn]
				if num <= 0 {
					continue
				}
				// 找到特定一个时间的计划数据
				minuteTime := hourTime + mn*60
				rs = append(rs, [2]int{minuteTime, num})
			}
		}
	}
	return rs
}

/**
 * 根据奖品的发奖计划，把设定的奖品数量放入奖品池
 * 需要每分钟执行一次
 * 【难点】定时程序，根据奖品设置的数据，更新奖品池的数据
 */
func DistributionGiftPool() int {
	totalNum := 0
	now := comm.NowUnix()
	giftService := services.NewGiftService()
	list := giftService.GetAll(false) //后台程序读的频率很低
	if list != nil && len(list) > 0 {
		for _, gift := range list {
			// 是否正常状态
			if gift.SysStatus != 0 {
				continue
			}
			// 是否限量产品
			if gift.PrizeNum < 1 {
				continue
			}
			// 时间段是否正常
			if gift.TimeBegin > now || gift.TimeEnd < now {
				continue
			}
			// 计划数据的长度太短，不需要解析和执行
			// 发奖计划，[[时间1,数量1],[时间2,数量2]]
			if len(gift.PrizeData) <= 7 { //光时间就超过7位
				continue
			}
			var cronData [][2]int
			err := json.Unmarshal([]byte(gift.PrizeData), &cronData)
			if err != nil {
				log.Println("prizedata.DistributionGiftPool Unmarshal error=", err)
			} else {
				index := 0
				giftNum := 0
				for i, data := range cronData {
					ct := data[0]
					num := data[1]
					if ct <= now {
						// 之前没有执行的数量，都要放进奖品池
						giftNum += num
						index = i + 1 //偏移量
					} else {
						break
					}
				}
				// 有奖品需要放入到奖品池
				if giftNum > 0 {
					incrGiftPool(gift.Id, giftNum)
					totalNum += giftNum
				}
				// 有计划数据被执行过，需要更新发奖计划
				if index > 0 {
					if index >= len(cronData) { //偏移到末尾
						cronData = make([][2]int, 0)
					} else {
						cronData = cronData[index:]
					}
					// 更新到数据库
					str, err := json.Marshal(cronData)
					if err != nil {
						log.Println("prizedata.DistributionGiftPool Marshal(cronData)", cronData, "error=", err)
					}
					columns := []string{"prize_data"}
					err = giftService.Update(&models.LtGift{
						Id:        gift.Id,
						PrizeData: string(str),
					}, columns)
					if err != nil {
						log.Println("prizedata.DistributionGiftPool giftService.Update error=", err)
					}
				}
			}
		}
		if totalNum > 0 {
			// 预加载缓存数据
			giftService.GetAll(true)
		}
	}
	return totalNum
}

// 根据计划数据，往奖品池增加奖品数量
func incrGiftPool(id, num int) int {
	key := "gift_pool"
	cacheObj := datasource.InstanceCache()
	rtNum, err := redis.Int64(cacheObj.Do("HINCRBY", key, id, num))
	if err != nil {
		log.Println("prizedata.incrServGiftPool error=", err)
		return 0
	}
	// 保证加入的库存数量正确的被加入到池中
	if int(rtNum) < num {
		// 加少了，补偿一次
		num2 := num - int(rtNum)
		rtNum, err = redis.Int64(cacheObj.Do("HINCRBY", key, id, num2))
		if err != nil {
			log.Println("prizedata.incrServGiftPool2 error=", err)
			return 0
		}
	}
	return int(rtNum)
}

// // 重置集群的奖品池
// func resetServGiftPool() {
// 	key := "gift_pool"
// 	cacheObj := datasource.InstanceCache()
// 	_, err := cacheObj.Do("DEL", key)
// 	if err != nil {
// 		log.Println("prizedata.resetServGiftPool DEL error=", err)
// 	}
// }
