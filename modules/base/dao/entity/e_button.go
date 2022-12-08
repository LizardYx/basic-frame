package entity

import (
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"gorm.io/gorm"
	"time"
)

type SchemaButton schema.Button

func (a SchemaButton) ToButton() *Button {
	item := new(Button)
	_ = common.Copy(a, item)
	return item
}

type SchemaButtons schema.Buttons

func (a SchemaButtons) ToButton() *Buttons {
	item := new(Buttons)
	_ = common.Copy(a, item)
	return item
}

type SchemaButtonPre schema.ButtonPre

func (a SchemaButtonPre) ToButton() *Button {
	item := new(Button)
	_ = common.Copy(a, item)
	return item
}

type Button struct {
	ID          uint64         `gorm:"primaryKey,autoIncrement;"`
	UUID        string         `gorm:"index;unique;"`              // 前端组装菜单需要的
	BtnID       int            `gorm:"index;"`                     // 前端识别按钮用
	Name        string         `gorm:"index;default:'';not null;"` // 按钮名称
	Icon        string         `gorm:""`                           // 按钮图标
	Class       string         `gorm:""`                           // 按钮样式
	MenuID      uint64         `gorm:"index;"`                     // 菜单ID
	Sequence    int            `gorm:"index;"`                     // 排序值
	ParentID    *uint64        `gorm:"index;"`                     // 父级按钮ID
	ShowStatus  int            `gorm:"index;default:2;not null;"`  // 显示状态(1:显示 2:隐藏)
	Status      int            `gorm:"index;default:2;not null;"`  // 状态(1:启用 2:禁用)
	Memo        string         `gorm:"size:1024;"`                 // 备注
	Creator     uint64         `gorm:"index"`
	CreatedAt   time.Time      `gorm:"index"`
	UpdatedAt   time.Time      `gorm:"index"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	RestfulApis []RestfulApi   `gorm:"many2many:button_restful_api;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // 按钮调用的Api
	SonButtons  []Button       `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`         // 子按钮
}

func (a Button) ToSchemaButton() *schema.Button {
	item := new(schema.Button)
	_ = common.Copy(a, item)
	return item.Init()
}

func (a Button) ToSchemaButtonPre() *schema.ButtonPre {
	item := new(schema.ButtonPre)
	_ = common.Copy(a, item)
	return item
}

type Buttons []*Button

func (a Buttons) ToSchemaButtons() schema.Buttons {
	list := make(schema.Buttons, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaButton()
	}
	return list
}

func (a Buttons) ToSchemaButtonPres() schema.ButtonPres {
	list := make(schema.ButtonPres, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaButtonPre().Init()
	}
	return list
}
