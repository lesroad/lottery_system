package comm

import (
	"crypto/md5"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"iris项目/my_lottery/conf"
	"iris项目/my_lottery/models"
	"net"
)

// 得到客户端IP地址
func ClientIP(request *http.Request) string {
	host, _, _ := net.SplitHostPort(request.RemoteAddr)
	return host
}

// 跳转URL
func Redirect(writer http.ResponseWriter, url string) {
	writer.Header().Add("Location", url)
	writer.WriteHeader(http.StatusFound)
}

// 从cookie中得到当前登录的用户
func GetLoginUser(request *http.Request) *models.ObjLoginuser {
	c, err := request.Cookie("lottery_loginuser")
	//没有cookie
	if err != nil {
		return nil
	}
	params, err := url.ParseQuery(c.Value) //返回Values类型的字典
	if err != nil {
		return nil
	}
	uid, err := strconv.Atoi(params.Get("uid"))
	if err != nil || uid < 1 {
		return nil
	}
	// Cookie最长使用时长
	now, err := strconv.Atoi(params.Get("now"))
	if err != nil || NowUnix()-now > 86400*30 { //超过30天cookie失效
		return nil
	}
	//// IP修改了是不是要重新登录
	//ip := params.Get("ip")
	//if ip != ClientIP(request) {
	//	return nil
	//}
	// 构建login对象
	loginuser := &models.ObjLoginuser{}
	loginuser.Uid = uid
	loginuser.Username = params.Get("username")
	loginuser.Now = now
	loginuser.Ip = ClientIP(request)
	loginuser.Sign = params.Get("sign")
	if err != nil {
		log.Println("fuc_web GetLoginUser Unmarshal ", err)
		return nil
	}
	sign := createLoginuserSign(loginuser)
	if sign != loginuser.Sign { //签名不一致
		log.Println("fuc_web GetLoginUser createLoginuserSign not sign", sign, loginuser.Sign)
		return nil
	}

	return loginuser
}

// 根据登录用户信息生成签名字符串
func createLoginuserSign(loginuser *models.ObjLoginuser) string {
	str := fmt.Sprintf("uid=%d&username=%s&secret=%s", loginuser.Uid, loginuser.Username, conf.CookieSecret)
	sign := fmt.Sprintf("%x", md5.Sum([]byte(str))) //转化16进制
	return sign
}

// 将登录的用户信息设置到cookie中
func SetLoginuser(writer http.ResponseWriter, loginuser *models.ObjLoginuser) {
	if loginuser == nil || loginuser.Uid < 1 {
		c := &http.Cookie{
			Name:   "lottery_loginuser",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		}
		http.SetCookie(writer, c)
		return
	}
	if loginuser.Sign == "" {
		loginuser.Sign = createLoginuserSign(loginuser) //生成签名
	}
	params := url.Values{}
	params.Add("uid", strconv.Itoa(loginuser.Uid))
	params.Add("username", loginuser.Username)
	params.Add("now", strconv.Itoa(loginuser.Now))
	params.Add("ip", loginuser.Ip)
	params.Add("sign", loginuser.Sign)
	c := &http.Cookie{
		Name:  "lottery_loginuser",
		Value: params.Encode(),
		Path:  "/",
	}
	http.SetCookie(writer, c)
}
