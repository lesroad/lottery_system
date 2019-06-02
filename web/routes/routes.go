package routes

import (
	"iris项目/my_lottery/bootstrap"
	"iris项目/my_lottery/services"
	"iris项目/my_lottery/web/controllers"
	"iris项目/my_lottery/web/middleware"

	"github.com/kataras/iris/mvc"
)

func Configure(b *bootstrap.Bootstrapper) {
	giftService := services.NewGiftService()
	userdayService := services.NewUserdayService()
	userService := services.NewUserService()
	blackipService := services.NewBlackipService()
	codeService := services.NewCodeService()
	resultService := services.NewResultService()

	index := mvc.New(b.Party("/"))
	index.Register(giftService, userdayService, userService, blackipService, codeService, resultService)
	index.Handle(new(controllers.IndexController))

	admin := mvc.New(b.Party("/admin"))
	admin.Router.Use(middleware.BasicAuth)
	admin.Register(giftService, userdayService, userService, blackipService, codeService, resultService)
	admin.Handle(new(controllers.AdminController))

	adminResult := admin.Party("/result")
	adminResult.Register(resultService)
	adminResult.Handle(new(controllers.AdminResultController))

	adminGift := admin.Party("/gift")
	adminGift.Register(giftService)
	adminGift.Handle(new(controllers.AdminGiftController))

	adminCode := admin.Party("/code")
	adminCode.Register(codeService)
	adminCode.Handle(new(controllers.AdminCodeController))

	adminUser := admin.Party("/user")
	adminUser.Register(userService)
	adminUser.Handle(new(controllers.AdminUserController))

	adminBlackip := admin.Party("/blackip")
	adminBlackip.Register(blackipService)
	adminBlackip.Handle(new(controllers.AdminBlackipController))
}
