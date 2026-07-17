package common

var (
	PageSize          uint    = 10
	VisitorPageSize   uint    = 8
	Version           string  = "0.3.9"
	VisitorExpire     float64 = 600
	Upload            string  = "static/upload/"
	Dir               string  = "config/"
	MysqlConf         string  = Dir + "mysql.json"
	IsCompireTemplate bool    = false //是否编译静态模板到二进制
	PublicKefuName    string  = "線上客服"
	DefaultWelcome    string  = "您好，這裡是線上客服，請問有什麼可以幫您？"
	DefaultOffline    string  = "目前客服暫時離線，請留下訊息，我們會盡快回覆您。"
	DefaultAllNotice  string  = "歡迎使用線上客服。"
	KefuAvatar        string  = "/static/images/kf.png"
)
