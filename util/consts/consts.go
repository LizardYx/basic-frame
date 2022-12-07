package consts

// 系统默认基础配置信息
const (
	AppName              = "BasicFrame"          // 系统名称
	AppVersion           = "0.1"                 // 系统版本号
	RunMode              = RunModeDev            // 系统默认运行模式
	DefaultLang          = LanguageEn            // 系统默认语言
	SystemConfigFileName = "./config/config.ini" // 系统基础配置文件
)

// 运行模式
const (
	RunModeDebug   = "Debug"   // 系统运行模式(debug)
	RunModeDev     = "Develop" // 系统运行模式(开发)
	RunModeTest    = "Test"    // 系统运行模式(测试)
	RunModeRelease = "Release" // 系统运行模式(线上)
)

// Mysql默认配置信息
const (
	MysqlHost     = "0.0.0.0"     // mysql主机地址
	MysqlPort     = "3306"        // mysql连接端口
	MysqlUser     = "root"        // mysql连接用户名
	MysqlPassword = ""            // mysql连接密码
	MysqlDBName   = "basic_frame" // mysql数据库名称
)

// WebServer默认配置信息
const (
	WebServerHost = "0.0.0.0" // Web服务IP地址
	WebServerPort = "8080"    // Web服务端口
	UseHttpsMode  = "false"   // Web服务是否启用HTTPS
	HttpsCrtFile  = ""        // HTTPS Crt证书
	HttpsKeyFile  = ""        // HTTPS私钥
)

// JWTAuth默认配置信息
const (
	JWTAuthSecretKeyLen = 16       // JWTAuth认证密钥长度
	JWTAuthExpired      = "604800" // JWTAuth认证过期时间(秒)
)

// Logger默认配置信息
const (
	LoggerRotationTime = "24"  // 设置日志分割的时间(单位:小时)
	LoggerMaxAge       = "168" // 设置文件清理前的最长保存时间(单位:小时)
)

// SmtpServer邮件服务默认配置信息
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

// Websocket
const (
	WSEventIDConnectSuccess = 1  // Websocket连接成功
	WSEventIDConnectClosed  = -1 // Websocket连接关闭
)

// Restful api
const (
	ApiOperateSuccess = "operate success"
)
