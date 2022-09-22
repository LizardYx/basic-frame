package common

// SystemConfig 系统基础配置信息
type SystemConfig struct {
	AppName    string  // App名称
	AppVersion float64 // App版本号
	RunMode    string  // App运行模式
	Cors       bool    // 是否允许跨域
	Mysql      Mysql
	WebServer  WebServer
	JWTAuth    JWTAuth
	SmtpServer SmtpServer
}

// Mysql mysql数据库配置信息
type Mysql struct {
	Host     string // 数据库主机地址
	Port     int    // 数据库连接端口
	User     string // 数据库用户名
	Password string // 数据库密码
	DBName   string // 数据库名称
}

// WebServer Web服务配置信息
type WebServer struct {
	Host      string // Web服务IP地址
	Port      int    // Web服务端口
	HttpsMode bool   // 是否启用Https
	CrtFile   string // Https Crt证书
	KeyFile   string // Https私钥
}

// JWTAuth JWTAuth认证配置信息
type JWTAuth struct {
	SigningMethod string // 加密方式
	SigningKey    string // 密钥
	Expired       int    // 过期时间
}

// SmtpServer 邮件服务配置信息
type SmtpServer struct {
	EmailServerName string // 邮件服务名称
	SmtpHost        string // 邮件服务IP地址
	SmtpPort        int    // 邮件服务端口
	SmtpPassword    string // 邮件服务密码
}
