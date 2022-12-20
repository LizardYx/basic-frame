package entity

import (
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"gorm.io/gorm"
	"time"
)

type SchemaTagManage schema.TagManage

func (a SchemaTagManage) ToTagManage() *TagManage {
	item := new(TagManage)
	_ = common.Copy(a, item)
	return item
}

type TagManage struct {
	ID        uint64         `gorm:"primaryKey,autoIncrement;"`
	Type      string         `gorm:"index;default:'';not null;"` // 标签类型
	KeyName   string         `gorm:"index"`                      // 标签名称
	Value     string         `gorm:"index;default:'';not null;"` // 标签值
	Creator   uint64         `gorm:"index"`
	CreatedAt time.Time      `gorm:"index"`
	UpdatedAt time.Time      `gorm:"index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (a TagManage) ToSchemaTagManage() *schema.TagManage {
	item := new(schema.TagManage)
	_ = common.Copy(a, item)
	return item
}

type TagManages []*TagManage

func (a TagManages) ToSchemaTagManages() schema.TagManages {
	list := make(schema.TagManages, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaTagManage()
	}
	return list
}
