package wxccserver

// Config 配置
type Config struct {
	Port         int      `default:"9527"`  // 服务监听端口
	EnableTLS    bool     `default:"false"` // 启用TLS
	TLSCertFile  string   `default:""`
	TLSKeyFile   string   `default:""`
	AccountsFile string   `default:"./config/accounts.toml"` // 公众号账号配置文件
	AheadTimeout int      `default:"60"`                     // access token 和 jsapi ticket提前超时时间(单位: 秒);
	AllowIP      []string // 允许访问的ip; 为空表示不限制
	Signature    struct {
		Key     string // 默认签名密钥
		Timeout int    `default:"10"` // 签名超时时间(单位:秒)
	}
	Log struct {
		Level        string `default:"debug"`  // 日志记录级别
		Path         string `default:"./logs"` // 日志保存路径
		RotationTime int    `default:"24"`     // 日志分割时间(单位:时)
		MaxAge       int    `default:"168"`    // 日志最大保存时间(单位:时)
	}
}
