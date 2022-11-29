package entity

import (
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"gorm.io/gorm"
	"time"
)

type SchemaSystemConfig schema.SystemConfig

func (a SchemaSystemConfig) ToSystemConfig() *SystemConfig {
	item := new(SystemConfig)
	_ = common.Copy(a, item)
	return item
}

type SystemConfig struct {
	ID            uint64         `gorm:"primaryKey,autoIncrement;"`  // 唯一标识
	WebServerHost string         `gorm:"index;default:'';not null;"` // Web服务IP地址
	WebServerPort uint64         `gorm:"index;default:0;not null;"`  // Web服务端口
	HttpsMode     bool           `gorm:"index"`                      // 是否启用Https
	HttpsCrtFile  string         `gorm:"index"`                      // Https Crt证书
	HttpsKeyFile  string         `gorm:"index"`                      // Https私钥
	MenuVersion   float64        `gorm:"index"`                      // 前端菜单版本号(只能通过定时任务，由菜单json中的版本号字段修改)
	CreatedAt     time.Time      `gorm:"index"`                      // 创建时间
	UpdatedAt     time.Time      `gorm:"index"`                      // 更新时间
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

func (a SystemConfig) ToSchemaSystemConfig() schema.SystemConfig {
	item := new(schema.SystemConfig)
	_ = common.Copy(a, item)
	return *item
}

func (a SystemConfig) ToSchemaSystemConfigPre() schema.SystemConfigPre {
	item := new(schema.SystemConfigPre)
	_ = common.Copy(a, item)
	return *item
}

type SystemConfigs []*SystemConfig

func (a SystemConfigs) ToSchemaSystemConfigs() schema.SystemConfigs {
	item := new(schema.SystemConfigs)
	_ = common.Copy(a, item)
	return *item
}

func (a SystemConfigs) ToSchemaSystemConfigPres() schema.SystemConfigPres {
	item := new(schema.SystemConfigPres)
	_ = common.Copy(a, item)
	return *item.Init()
}
