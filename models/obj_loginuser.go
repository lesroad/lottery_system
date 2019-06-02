package models

// 站点中与浏览器交互的用户模型
type ObjLoginuser struct {
	Uid      int
	Username string
	Now      int
	Ip       string
	Sign     string //签名数据
}
