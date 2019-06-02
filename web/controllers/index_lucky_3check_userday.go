package controllers

import (
	"fmt"
	"iris项目/my_lottery/conf"
	"iris项目/my_lottery/models"
	"iris项目/my_lottery/services"
	"iris项目/my_lottery/web/utils"
	"log"
	"strconv"
	"time"
)

func (c *IndexController) checkUserday(uid int, num int64) bool {
	userdayService := services.NewUserdayService()
	userdayInfo := userdayService.GetUserToday(uid) //返回一个用户今天抽奖的对象
	if userdayInfo != nil && userdayInfo.Uid == uid {
		// 今天存在抽奖记录
		if userdayInfo.Num >= conf.UserPrizeMax {
			if int(num) < userdayInfo.Num {
				//以数据库中的为准
				utils.InitUserLuckyNum(uid, int64(userdayInfo.Num))
			}
			return false
		} else {
			userdayInfo.Num++
			if int(num) < userdayInfo.Num {
				utils.InitUserLuckyNum(uid, int64(userdayInfo.Num))
			}
			err103 := userdayService.Update(userdayInfo, nil)
			if err103 != nil {
				log.Println("index_lucky_check_userday ServiceUserDay.Update "+
					"err103=", err103)
			}
		}
	} else {
		// 创建今天的用户参与记录
		y, m, d := time.Now().Date()
		strDay := fmt.Sprintf("%d%02d%02d", y, m, d)
		day, _ := strconv.Atoi(strDay)
		userdayInfo = &models.LtUserday{
			Uid:        uid,
			Day:        day,
			Num:        1,
			SysCreated: int(time.Now().Unix()),
		}
		err103 := userdayService.Create(userdayInfo)
		if err103 != nil {
			log.Println("index_lucky_check_userday ServiceUserDay.Create "+
				"err103=", err103)
		}
		utils.InitUserLuckyNum(uid, 1)
	}
	return true
}
