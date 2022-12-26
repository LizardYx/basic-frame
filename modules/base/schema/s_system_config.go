package schema

import (
	"basic-frame/util/common"
	"time"
)

// SystemConfig 前端请求使用的结构体
type SystemConfig struct {
	ID            uint64    `json:"id"`                                 // 唯一标识
	WebServerHost string    `json:"web_server_host" binding:"required"` // Web服务IP地址
	WebServerPort uint64    `json:"web_server_port" binding:"required"` // Web服务端口
	HttpsMode     bool      `json:"https_mode"`                         // 是否启用Https
	HttpsCrtFile  string    `json:"https_crt_file"`                     // Https Crt证书
	HttpsKeyFile  string    `json:"https_key_file"`                     // Https私钥
	MenuVersion   float64   `json:"menu_version"`                       // 前端菜单版本号(只能通过定时任务，由菜单json中的版本号字段修改)
	CreatedAt     time.Time `json:"created_at"`                         // 创建时间
	UpdatedAt     time.Time `json:"updated_at"`                         // 更新时间
}

type SystemConfigQueryParam struct {
	common.PaginationParam
}

type SystemConfigs []*SystemConfig

type SystemConfigQueryResult struct {
	Data       SystemConfigPres
	PageResult *common.PaginationResult
}

// ---------------------------------------- Response Struct --------------------------------------

// SystemConfigPre 接口返回的结构体
type SystemConfigPre struct {
	ID            uint64    `json:"id"`                                 // 唯一标识
	WebServerHost string    `json:"web_server_host" binding:"required"` // Web服务IP地址
	WebServerPort uint64    `json:"web_server_port" binding:"required"` // Web服务端口
	HttpsMode     bool      `json:"https_mode"`                         // 是否启用Https
	HttpsCrtFile  string    `json:"https_crt_file"`                     // Https Crt证书
	HttpsKeyFile  string    `json:"https_key_file"`                     // Https私钥
	MenuVersion   float64   `json:"menu_version"`                       // 前端菜单版本号(只能通过定时任务，由菜单json中的版本号字段修改)
	CreatedAt     time.Time `json:"created_at"`                         // 创建时间
	UpdatedAt     time.Time `json:"updated_at"`                         // 更新时间
}

type SystemConfigPres []*SystemConfigPre

func (a SystemConfigPres) Init() *SystemConfigPres {
	items := make(SystemConfigPres, 0)

	items = append(items, a...)
	return &items
}
