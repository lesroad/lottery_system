package controllers

import (
	"iris项目/my_lottery/conf"
	"iris项目/my_lottery/models"
	"iris项目/my_lottery/services"
)

func (api *IndexController) prize(prizeCode int, limitBlack bool) *models.ObjGiftPrize {
	var prizeGift *models.ObjGiftPrize
	giftList := services.NewGiftService().GetAllUse(true)
	for _, gift := range giftList {
		if gift.PrizeCodeA <= prizeCode &&
			gift.PrizeCodeB >= prizeCode {
			// 中奖编码区间满足条件，说明可以中奖
			if !limitBlack || gift.Gtype < conf.GtypeGiftSmall { //如果非实物奖直接发，实物奖需要看是不是在黑名单外
				prizeGift = &gift
				break
			}
		}
	}
	return prizeGift
}
