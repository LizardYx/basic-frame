package common

var (
	SysConfig = new(SystemConfig)
)

// SystemConfig 系统基础配置信息
type SystemConfig struct {
	RunMode     string  // App运行模式
	AppName     string  // App名称
	AppVersion  float64 // App版本号
	DefaultLang string  // 系统默认语言
	MenuVersion float64 // 系统菜单版本号
	Mysql       Mysql
	WebServer   WebServer
	JWTAuth     JWTAuth
	Logger      LoggerConf
	SmtpServer  SmtpServer
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
	Host         string // Web服务IP地址
	Port         int    // Web服务端口
	HttpsMode    bool   // 是否启用Https
	HttpsCrtFile string // Https Crt证书
	HttpsKeyFile string // Https私钥
}

// JWTAuth JWTAuth认证配置信息
type JWTAuth struct {
	SigningMethod string // 加密方式
	SigningKey    string // 密钥
	Expired       int    // 过期时间
}

type LoggerConf struct {
	RotationTime int // 设置日志切割时间间隔(单位:小时)
	MaxAge       int // 设置文件清理前的最长保存时间(单位:小时)
}

// SmtpServer 邮件服务配置信息
type SmtpServer struct {
	Enable          bool   // 是否启用邮件服务
	EmailServerName string // 邮件服务名称
	SmtpHost        string // 邮件服务IP地址
	SmtpPort        int    // 邮件服务端口
	SmtpPassword    string // 邮件服务密码
}

// PaginationParam 分页查询条件
type PaginationParam struct {
	Pagination bool `form:"-"`                    // 是否使用分页查询
	OnlyCount  bool `form:"-"`                    // 是否仅查询count
	Current    int  `form:"current,default=1"`    // 当前页
	PageSize   int  `form:"pageSize,default=100"` // 页大小
}

// PaginationResult 分页查询结果
type PaginationResult struct {
	Total    int64 `json:"total"`
	Current  int   `json:"current"`
	PageSize int   `json:"pageSize"`
}

// ResponseError 定义响应错误
type ResponseError struct {
	Code       int    // 错误码
	Message    string // 错误消息
	StatusCode int    // 响应状态码
	ERR        error  // 响应错误
}

// IDResult 响应唯一标识
type IDResult struct {
	ID uint64 `json:"id,string"`
}
