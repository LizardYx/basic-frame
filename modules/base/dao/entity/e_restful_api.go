package entity

import (
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"gorm.io/gorm"
	"time"
)

type SchemaRestfulApi schema.RestfulApi

func (a SchemaRestfulApi) ToRestfulApi() *RestfulApi {
	item := new(RestfulApi)
	_ = common.Copy(a, item)
	return item
}

type SchemaRestfulApis schema.RestfulApis

func (a SchemaRestfulApis) ToRestfulApis() *RestfulApis {
	item := new(RestfulApis)
	_ = common.Copy(a, item)
	return item
}

type RestfulApi struct {
	ID        uint64         `gorm:"primaryKey,autoIncrement;"`
	UUID      string         `gorm:"index;"`                     // 接口UUID
	Method    string         `gorm:"index;default:'';not null;"` // 资源请求方式(支持正则)
	Path      string         `gorm:"index;default:'';not null;"` // 资源请求路径（支持/:id匹配）
	Memo      string         `gorm:"size:1024;"`                 // 备注
	Creator   uint64         `gorm:"index"`
	CreatedAt time.Time      `gorm:"index"`
	UpdatedAt time.Time      `gorm:"index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (a RestfulApi) ToSchemaRestfulApi() *schema.RestfulApi {
	item := new(schema.RestfulApi)
	_ = common.Copy(a, item)
	return item
}

type RestfulApis []*RestfulApi

func (a RestfulApis) ToSchemaRestfulApis() schema.RestfulApis {
	item := new(schema.RestfulApis)
	_ = common.Copy(a, item)
	return *item.Init()
}
