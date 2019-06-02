package controllers

import (
	"fmt"
	"iris项目/my_lottery/comm"
	"iris项目/my_lottery/conf"
	"iris项目/my_lottery/models"
	"iris项目/my_lottery/web/utils"
	"log"
)

// localhost:8080/lucky
func (c *IndexController) GetLucky() map[string]interface{} {
	rs := make(map[string]interface{})
	rs["code"] = 0
	rs["msg"] = ""
	// 1 验证登录用户
	loginuser := comm.GetLoginUser(c.Ctx.Request())
	if loginuser == nil || loginuser.Uid < 1 {
		rs["code"] = 101
		rs["msg"] = "请先登录，再来抽奖"
		return rs
	}

	// 2 用户抽奖分布式锁定
	ok := utils.LockLucky(loginuser.Uid)
	if ok {
		defer utils.UnlockLucky(loginuser.Uid)
	} else {
		rs["code"] = 102
		rs["msg"] = "正在抽奖，请稍后重试"
		return rs
	}

	// 3 验证用户今日参与次数（需要有数据表，不能是大概值）
	//缓存
	userDayNum := utils.IncrUserLuckyNum(loginuser.Uid)
	if userDayNum > conf.UserPrizeMax {
		rs["code"] = 103
		rs["msg"] = "今日抽奖次数已用完，请明天再来"
		return rs
	} else {
		ok = c.checkUserday(loginuser.Uid, userDayNum) //mysql中验证
		if !ok {
			rs["code"] = 103
			rs["msg"] = "今日抽奖次数已用完，请明天再来"
			return rs
		}
	}

	// 4 验证IP今日的参与次数（只在缓存中更新，抽奖次数可以是大概值）
	ip := comm.ClientIP(c.Ctx.Request())
	ipDayNum := utils.IncrIpLuckyNum(ip)
	if ipDayNum > conf.IpLimitMax {
		rs["code"] = 104
		rs["msg"] = "相同IP参与次数太多，明天再来参与吧"
		return rs
	}

	limitBlack := false
	if ipDayNum > conf.IpPrizeMax {
		limitBlack = true
	}
	// 5 验证IP黑名单(缓存+数据库)
	var blackipInfo *models.LtBlackip
	if !limitBlack {
		ok, blackipInfo = c.checkBlackip(ip)
		// IP黑名单存在，而且还在黑名单有效期内
		if !ok {
			fmt.Println("黑名单中的IP", ip, blackipInfo)
			limitBlack = true
		}
	}

	// 6 验证用户黑名单(缓存+数据库)
	var userInfo *models.LtUser
	if !limitBlack {
		ok, userInfo = c.checkBlackUser(loginuser.Uid)
		if !ok {
			fmt.Println("黑名单中的用户", loginuser.Uid, userInfo)
			limitBlack = true
		}
	}

	// 7 获得抽奖编码
	prizeCode := comm.Random(10000)
	// 8 匹配奖品是否中奖并发放不限量奖品（虚拟币）
	prizeGift := c.prize(prizeCode, limitBlack) //得到奖品对象
	if prizeGift == nil || (prizeGift.PrizeNum > 0 && prizeGift.LeftNum <= 0) {
		rs["code"] = 205
		rs["msg"] = "很遗憾，没有中奖，请下次再试"
		return rs
	}

	// 9 限量奖品发放
	if prizeGift.PrizeNum > 0 {
		if utils.GetGiftPoolNum(prizeGift.Id) <= 0 { //先从奖品池得到奖品数量判断是否可以发奖
			rs["code"] = 206
			rs["msg"] = "很遗憾，没有中奖，请下次再试"
			return rs
		}
		ok = utils.PrizeGift(prizeGift.Id, prizeGift.LeftNum) //再更新库存并返回结果
		if !ok {
			rs["code"] = 207
			rs["msg"] = "很遗憾，没有中奖，请下次再试"
			return rs
		}
	}

	// 10 不同编码的优惠券发放
	if prizeGift.Gtype == conf.GtypeCodeDiff {
		code := utils.PrizeCodeDiff(prizeGift.Id, c.ServiceCode) //更新库存并返回抽到的券编码
		if code == "" {
			rs["code"] = 208
			rs["msg"] = "很遗憾，没有中奖，请下次再试"
			return rs
		}
		prizeGift.Gdata = code
	}

	// 11 记录中奖纪录
	result := models.LtResult{
		GiftId:     prizeGift.Id,
		GiftName:   prizeGift.Title,
		GiftType:   prizeGift.Gtype,
		Uid:        loginuser.Uid,
		Username:   loginuser.Username,
		PrizeCode:  prizeCode,
		GiftData:   prizeGift.Gdata,
		SysCreated: comm.NowUnix(),
		SysIp:      ip,
		SysStatus:  0,
	}
	err := c.ServiceResult.Create(&result)
	if err != nil {
		log.Println("index_lucky.GetLucky ServiceResult.Create ", result,
			", error=", err)
		rs["code"] = 209
		rs["msg"] = "很遗憾，没有中奖，请下次再试"
		return rs
	}
	if prizeGift.Gtype == conf.GtypeGiftLarge {
		// 如果获得了实物大奖，需要将用户、IP设置成黑名单一段时间
		c.prizeLarge(ip, loginuser, userInfo, blackipInfo)
	}
	// 12 返回抽奖结果
	rs["gift"] = prizeGift
	return rs
}
