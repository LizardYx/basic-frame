package entity

import (
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"gorm.io/gorm"
	"time"
)

type SchemaMenu schema.Menu

func (a SchemaMenu) ToMenu() *Menu {
	item := new(Menu)
	_ = common.Copy(a, item)
	return item
}

type SchemaMenus schema.Menus

func (a SchemaMenus) ToMenu() *Menus {
	item := new(Menus)
	_ = common.Copy(a, item)
	return item
}

type SchemaMenuTree schema.MenuTree

func (a SchemaMenuTree) ToMenu() *Menu {
	item := new(Menu)
	_ = common.Copy(a, item)
	return item
}

type Menu struct {
	ID          uint64         `gorm:"primaryKey,autoIncrement;"`
	UUID        string         `gorm:"index;unique;"`              // 前端组装菜单需要的
	Name        string         `gorm:"index;default:'';not null;"` // 菜单名称
	Icon        string         `gorm:""`                           // 菜单图标
	Class       string         `gorm:""`                           // 菜单样式
	Router      string         `gorm:"index;default:'';not null;"` // 访问路由
	Sequence    int            `gorm:"index;"`                     // 排序值
	ParentID    *uint64        `gorm:"index;"`                     // 父级菜单ID
	ShowStatus  int            `gorm:"index;default:2;not null;"`  // 显示状态(1:显示 2:隐藏)
	Status      int            `gorm:"index;default:2;not null;"`  // 状态(1:启用 2:禁用)
	Memo        string         `gorm:"size:1024;"`                 // 备注
	Creator     uint64         `gorm:"index"`
	CreatedAt   time.Time      `gorm:"index"`
	UpdatedAt   time.Time      `gorm:"index"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	RestfulApis []RestfulApi   `gorm:"many2many:menu_restful_api;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"` // 页面调用的Api
	Buttons     []Button       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`                            // 页面的按钮
	SonMenus    []Menu         `gorm:"foreignKey:ParentID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`       // 子菜单
}

func (a Menu) ToSchemaMenu() *schema.Menu {
	item := new(schema.Menu)
	_ = common.Copy(a, item)
	return item.Init()
}

type Menus []*Menu

func (a Menus) ToSchemaMenus() schema.Menus {
	list := make(schema.Menus, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaMenu()
	}
	return list
}

func (a Menu) ToSchemaMenuTree() *schema.MenuTree {
	item := new(schema.MenuTree)
	_ = common.Copy(a, item)
	return item
}

func (a Menus) ToSchemaMenuTrees() schema.MenuTrees {
	list := make(schema.MenuTrees, len(a))
	for i, item := range a {
		list[i] = item.ToSchemaMenuTree().Init()
	}
	return *list.SortMenuTrees()
}
