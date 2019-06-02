package bootstrap

import (
	"time"

	"iris项目/my_lottery/conf"
	"iris项目/my_lottery/cron"

	"github.com/gorilla/securecookie"
	"github.com/kataras/iris"
	"github.com/kataras/iris/sessions"
)

const (
	//站点对外目录
	StaticAssets = "./public/"
	Favicon      = "favicon.ico"
)

//定义配置器
type Configurator func(*Bootstrapper)

// 使用Go内建的嵌入机制(匿名嵌入)，允许类型之前共享代码和数据
// （Bootstrapper继承和共享 iris.Application ）
// 参考文章： https://hackthology.com/golangzhong-de-mian-xiang-dui-xiang-ji-cheng.html

type Bootstrapper struct {
	*iris.Application
	AppName      string
	AppOwner     string
	AppSpawnDate time.Time

	Sessions *sessions.Sessions
}

//实例化Bootstrapper
func New(appName, appOwner string, cfgs ...Configurator) *Bootstrapper {
	b := &Bootstrapper{
		AppName:      appName,
		AppOwner:     appOwner,
		AppSpawnDate: time.Now(),
		Application:  iris.New(),
	}

	for _, cfg := range cfgs {
		cfg(b) //配置Bootstrapper
	}
	return b
}

//初始化模板
func (b *Bootstrapper) SetupViews(viewsDir string) {
	htmlEngine := iris.HTML(viewsDir, ".html").Layout("shared/layout.html") //从 "./views" 文件夹加载所以的模板,其中扩展名为“.html”并解析它们
	// 每次重新加载模版（线上关闭它）
	htmlEngine.Reload(true)
	// 给模版内置各种定制的方法
	htmlEngine.AddFunc("FromUnixtimeShort", func(t int) string {
		dt := time.Unix(int64(t), int64(0))     //创建本地时间
		return dt.Format(conf.SysTimeformShort) //根据指定的格式返回dt代表的时间点的格式化文本表示
	})
	htmlEngine.AddFunc("FromUnixtime", func(t int) string {
		dt := time.Unix(int64(t), int64(0))
		return dt.Format(conf.SysTimeform)
	})
	b.RegisterView(htmlEngine)
}

// SetupSessions initializes the sessions, optionally.
func (b *Bootstrapper) SetupSessions(expires time.Duration, cookieHashKey, cookieBlockKey []byte) {
	b.Sessions = sessions.New(sessions.Config{
		Cookie:   "SECRET_SESS_COOKIE_" + b.AppName, //session id
		Expires:  expires,                           //过期时间
		Encoding: securecookie.New(cookieHashKey, cookieBlockKey),
	})
}

//// SetupWebsockets prepares the websocket server.
//func (b *Bootstrapper) SetupWebsockets(endpoint string, onConnection websocket.ConnectionFunc) {
//	ws := websocket.New(websocket.Config{})
//	ws.OnConnection(onConnection)
//
//	b.Get(endpoint, ws.Handler())
//	b.Any("/iris-ws.js", func(ctx iris.Context) {
//		ctx.Write(websocket.ClientSource)
//	})
//}

//异常处理
func (b *Bootstrapper) SetupErrorHandlers() {
	b.OnAnyErrorCode(func(ctx iris.Context) {
		err := iris.Map{
			"app":     b.AppName,
			"status":  ctx.GetStatusCode(),
			"message": ctx.Values().GetString("message"),
		}

		if jsonOutput := ctx.URLParamExists("json"); jsonOutput {
			ctx.JSON(err)
			return
		}

		ctx.ViewData("Err", err)
		ctx.ViewData("Title", "Error")
		ctx.View("shared/error.html")
	})
}

//配置
func (b *Bootstrapper) Configure(cs ...Configurator) {
	for _, c := range cs {
		c(b)
	}
}

//切换任务
func (b *Bootstrapper) setupCron() {
	// 每过多久执行一次奖品发放计划更新
	cron.ConfigureAppOneCron()
}

//初始化Bootstrapper
func (b *Bootstrapper) Bootstrap() *Bootstrapper {
	b.SetupViews("./views") //设置模板
	b.SetupSessions(24*time.Hour,
		[]byte("the-big-and-secret-fash-key-here"),
		[]byte("lot-secret-of-characters-big-too"),
	)
	b.SetupErrorHandlers()                                         //异常信息
	b.Favicon(StaticAssets + Favicon)                              //默认图标
	b.StaticWeb(StaticAssets[1:len(StaticAssets)-1], StaticAssets) //静态站点
	b.setupCron()                                                  //启动其他任务
	// b.Use(recover.New())                                           // context.Handler 类型 每一个请求都会先执行此方法 app.Use(context.Handler)
	// b.Use(logger.New())
	return b
}

//监听
func (b *Bootstrapper) Listen(addr string, cfgs ...iris.Configurator) {
	b.Run(iris.Addr(addr), cfgs...)
}
