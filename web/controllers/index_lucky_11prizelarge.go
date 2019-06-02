package controllers

import (
	"iris项目/my_lottery/comm"
	"iris项目/my_lottery/models"
)

func (c *IndexController) prizeLarge(ip string,
	loginuser *models.ObjLoginuser,
	userinfo *models.LtUser,
	blackipInfo *models.LtBlackip) {
	// TO DO
	// userService := services.NewUserService()
	// blackipService := services.NewBlackipService()
	nowTime := comm.NowUnix()
	blackTime := 30 * 86400 //30天
	// 更新用户的黑名单信息
	if userinfo == nil || userinfo.Id <= 0 {
		userinfo = &models.LtUser{
			Id:         loginuser.Uid,
			Username:   loginuser.Username,
			Blacktime:  nowTime + blackTime,
			SysCreated: nowTime,
			SysIp:      ip,
		}
		c.ServiceUser.Create(userinfo)
	} else {
		userinfo = &models.LtUser{
			Id:         loginuser.Uid,
			Blacktime:  nowTime + blackTime,
			SysUpdated: nowTime,
		}
		c.ServiceUser.Update(userinfo, nil)
	}
	// 更新要IP的黑名单信息
	if blackipInfo == nil || blackipInfo.Id <= 0 {
		blackipInfo = &models.LtBlackip{
			Ip:         ip,
			Blacktime:  nowTime + blackTime,
			SysCreated: nowTime,
		}
		c.ServiceBlackip.Create(blackipInfo)
	} else {
		blackipInfo.Blacktime = nowTime + blackTime
		blackipInfo.SysUpdated = nowTime
		c.ServiceBlackip.Update(blackipInfo, nil)
	}
}
