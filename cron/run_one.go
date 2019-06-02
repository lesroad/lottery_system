package cron

import (
	"iris项目/my_lottery/comm"
	"iris项目/my_lottery/services"
	"iris项目/my_lottery/web/utils"
	"log"
	"time"
)

/**
 * 只需要一个应用运行的服务
 * 全局的服务
 */
func ConfigureAppOneCron() {
	// 每5分钟执行一次，奖品的发奖计划到期的时候，需要重新生成发奖计划
	go resetAllGiftPrizeData()
	// 每分钟执行一次，根据发奖计划，把奖品数量放入奖品池
	go distributionAllGiftPool()
}

// 重置所有奖品的发奖计划
// 每5分钟执行一次
func resetAllGiftPrizeData() {
	giftService := services.NewGiftService()
	list := giftService.GetAll(false)
	nowTime := comm.NowUnix()
	for _, giftInfo := range list {
		if giftInfo.PrizeTime != 0 &&
			(giftInfo.PrizeData == "" || giftInfo.PrizeEnd <= nowTime) {
			// 立即执行
			log.Println("更新发奖计划 giftInfo=", giftInfo)
			utils.ResetGiftPrizeData(&giftInfo, giftService)
			// 更新缓存数据
			giftService.GetAll(true)
		}
	}
	// 每5分钟执行一次
	time.AfterFunc(5*time.Minute, resetAllGiftPrizeData)
}

// 根据发奖计划，把奖品数量放入奖品池
// 每分钟执行一次
func distributionAllGiftPool() {
	num := utils.DistributionGiftPool()
	log.Println("奖品池的数量, num=", num)

	// 每分钟执行一次
	time.AfterFunc(time.Minute, distributionAllGiftPool)
}
