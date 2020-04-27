package wxccserver

// Account 表示一个公众号账号
type Account struct {
	AppID        string
	AppSecret    string
	SignatureKey string `default:""` // 为安全性，每个account可以独立一个签名key
}
