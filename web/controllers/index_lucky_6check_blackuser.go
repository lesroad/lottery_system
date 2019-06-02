package controllers

import (
	"iris项目/my_lottery/models"
	"iris项目/my_lottery/services"
	"time"
)

func (api *IndexController) checkBlackUser(uid int) (bool, *models.LtUser) {
	info := services.NewUserService().Get(uid)
	if info != nil && info.Blacktime > int(time.Now().Unix()) {
		// 黑名单存在并且有效
		return false, info
	}
	return true, info
}
