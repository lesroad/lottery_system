package controllers

import (
	"iris项目/my_lottery/models"
	"time"
)

// 验证当前用户的IP是否存在黑名单限制
func (c *IndexController) checkBlackip(ip string) (bool, *models.LtBlackip) {
	info := c.ServiceBlackip.GetByIp(ip)
	if info == nil || info.Ip == "" {
		return true, nil
	}
	if info.Blacktime > int(time.Now().Unix()) {
		// IP黑名单存在，而且还在黑名单有效期内
		return false, info
	}
	return true, info
}
