package controllers

import (
	"fmt"
	"iris项目/my_lottery/comm"
	"iris项目/my_lottery/conf"
	"iris项目/my_lottery/models"
	"iris项目/my_lottery/services"
	"iris项目/my_lottery/web/utils"

	"strings"

	"github.com/kataras/iris"
	"github.com/kataras/iris/mvc"
)

type AdminCodeController struct {
	Ctx            iris.Context
	ServiceUser    services.UserService
	ServiceGift    services.GiftService
	ServiceCode    services.CodeService
	ServiceResult  services.ResultService
	ServiceUserday services.UserdayService
	ServiceBlackip services.BlackipService
}

func (c *AdminCodeController) Get() mvc.Result {
	giftId := c.Ctx.URLParamIntDefault("gift_id", 0)
	page := c.Ctx.URLParamIntDefault("page", 1)
	size := 100
	pagePrev := ""
	pageNext := ""
	// 数据列表
	var datalist []models.LtCode
	var num int
	var cacheNum int
	if giftId > 0 { //指定某个券
		datalist = c.ServiceCode.Search(giftId)
		num, cacheNum = utils.GetCacheCodeNum(giftId, c.ServiceCode)
	} else {
		datalist = c.ServiceCode.GetAll(page, size)
	}
	total := (page - 1) + len(datalist)
	// 数据总数
	if len(datalist) >= size {
		if giftId > 0 {
			total = int(c.ServiceCode.CountByGift(giftId))
		} else {
			total = int(c.ServiceCode.CountAll())
		}
		pageNext = fmt.Sprintf("%d", page+1)
	}
	if page > 1 {
		pagePrev = fmt.Sprintf("%d", page-1)
	}
	return mvc.View{
		Name: "admin/code.html",
		Data: iris.Map{
			"Title":    "管理后台",
			"Channel":  "code",
			"GiftId":   giftId,
			"Datalist": datalist,
			"Total":    total,
			"PagePrev": pagePrev,
			"PageNext": pageNext,
			"CodeNum":  num,
			"CacheNum": cacheNum,
			//如果不一样需要手动点击页面重新整理
		},
		Layout: "admin/layout.html",
	}
}

// 导入不同编码的券
func (c *AdminCodeController) PostImport() {
	giftId := c.Ctx.URLParamIntDefault("gift_id", 0)
	fmt.Println("PostImport giftId=", giftId)
	if giftId < 1 {
		c.Ctx.Text("没有指定奖品ID，无法进行导入，<a href='' onclick='history.go(-1);return false;'>返回</a>")
		return
	}
	gift := c.ServiceGift.Get(giftId, true)
	if gift == nil || gift.Gtype != conf.GtypeCodeDiff {
		c.Ctx.HTML("奖品信息不存在或者奖品类型不是差异化优惠券，无法进行导入，<a href='' onclick='history.go(-1);return false;'>返回</a>")
		return
	}
	codes := c.Ctx.PostValue("codes")
	now := comm.NowUnix()
	list := strings.Split(codes, "\n")
	sucNum := 0
	errNum := 0
	for _, code := range list {
		code := strings.TrimSpace(code)
		if code != "" {
			data := &models.LtCode{
				GiftId:     giftId,
				Code:       code,
				SysCreated: now,
			}
			//导入数据库
			err := c.ServiceCode.Create(data)
			if err != nil {
				errNum++
			} else {
				// sucNum++
				//导入缓存
				ok := utils.ImportCacheCodes(giftId, code)
				if ok {
					sucNum++
				} else {
					errNum++
				}
			}
		}
	}
	c.Ctx.HTML(fmt.Sprintf("成功导入 %d 条，导入失败 %d 条，<a href='/admin/code?gift_id=%d'>返回</a>", sucNum, errNum, giftId))
}

func (c *AdminCodeController) GetDelete() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		c.ServiceCode.Delete(id)
	}
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/admin/code"
	}
	return mvc.Response{
		Path: refer,
	}
}

//恢复券
func (c *AdminCodeController) GetReset() mvc.Result {
	id, err := c.Ctx.URLParamInt("id")
	if err == nil {
		c.ServiceCode.Update(&models.LtCode{Id: id, SysStatus: 0}, []string{"sys_status"})
	}
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/admin/code"
	}
	return mvc.Response{
		Path: refer,
	}
}

// 重新整理优惠券的数据，如果是本地服务，也需要启动时加载
// 重新整理缓存数据
func (c *AdminCodeController) GetRecache() {
	refer := c.Ctx.GetHeader("Referer")
	if refer == "" {
		refer = "/admin/code"
	}
	id, err := c.Ctx.URLParamInt("id")
	if id < 1 || err != nil {
		rs := fmt.Sprintf("没有指定优惠券所属的奖品id, <a href='%s'>返回</a>", refer)
		c.Ctx.HTML(rs)
		return
	}
	sucNum, errNum := utils.RecacheCodes(id, c.ServiceCode)

	rs := fmt.Sprintf("sucNum=%d, errNum=%d, <a href='%s'>返回</a>", sucNum, errNum, refer)
	c.Ctx.HTML(rs)
}
