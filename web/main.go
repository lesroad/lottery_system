package main

import (
	"fmt"
	"iris项目/my_lottery/bootstrap"
	"iris项目/my_lottery/web/routes"
)

var port = 8080

func newApp() *bootstrap.Bootstrapper {
	app := bootstrap.New("抽奖系统", "lesroad")
	app.Bootstrap()
	app.Configure(routes.Configure)
	return app
}

func main() {
	app := newApp()
	app.Listen(fmt.Sprintf(":%d", port))
}
