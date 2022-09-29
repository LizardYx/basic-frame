package consts

// 系统基础配置信息
const (
	AppName              = "BasicFrame"          // 系统名称
	AppVersion           = "0.1"                 // 系统版本号
	RunMode              = "DEBUG"               // 运行模式
	DefaultLang          = LanguageEn            // 系统默认语言
	SystemConfigFileName = "./config/config.ini" // 系统基础配置文件
)

// Mysql配置信息
const (
	MysqlHost     = "0.0.0.0"     // mysql主机地址
	MysqlPort     = "3306"        // mysql连接端口
	MysqlUser     = "root"        // mysql连接用户名
	MysqlPassword = ""            // mysql连接密码
	MysqlDBName   = "basic_frame" // mysql数据库名称
)

// WebServer配置信息
const (
	WebServerHost = "0.0.0.0" // Web服务IP地址
	WebServerPort = "8080"    // Web服务端口
	UseHttpsMode  = "false"   // Web服务是否启用HTTPS
	HttpsCrtFile  = ""        // HTTPS Crt证书
	HttpsKeyFile  = ""        // HTTPS私钥
)

// JWTAuth配置信息
const (
	JWTAuthSigningMethod = "HS256"       // JWTAuth加密方式
	JWTAuthSigningKey    = "basic-frame" // JWTAuth签名
	JWTAuthExpired       = "604800"      // JWTAuth认证过期时间(秒)
)

// SmtpServer邮件服务配置信息
const (
	SmtpEnable     = "false" // 是否启用Smtp邮件服务
	SmtpServerName = ""      // 邮件服务名称
	SmtpHost       = ""      // 邮件服务IP地址
	SmtpPort       = ""      // 邮件服务端口
	SmtpPassword   = ""      // 邮件服务密码
)

// 超管用户信息
const (
	RootName     = "root"    // 超管账号
	RootPassword = "1234asd" // 超管密码
)

// i18n配置
const (
	LanguageZh = "zh_CN" // 国际化(中文)
	LanguageEn = "en_US" // 国际化(英文)
)
