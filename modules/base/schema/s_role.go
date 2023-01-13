package schema

import (
	"basic-frame/util/common"
	"time"
)

type Role struct {
	ID             uint64         `json:"id"`                      // 唯一标识
	Name           string         `json:"name" binding:"required"` // 角色名称
	Sequence       int            `json:"sequence"`                // 排序值
	Type           int            `json:"type" binding:"required"` // 角色类型(1.用户角色 2.组织角色 3.职位角色 4.用户组角色)
	AuditorTypes   string         `json:"auditor_types"`           // 审核类型(逗号分隔)
	Status         int            `json:"status"`                  // 状态(1:启用 2:禁用)
	Memo           string         `json:"memo"`                    // 备注
	Creator        uint64         `json:"creator"`                 // 创建者
	CreatedAt      time.Time      `json:"created_at"`              // 创建时间
	UpdatedAt      time.Time      `json:"updated_at"`              // 更新时间
	Menus          Menus          `json:"menus"`                   // 角色能使用的菜单
	Buttons        Buttons        `json:"buttons"`                 // 角色能使用的按钮
	DisabledFields DisabledFields `json:"disabled_fields"`         // 角色禁用的字段
}

func (a Role) Init() *Role {
	a.Menus = *a.Menus.Init()
	a.Buttons = *a.Buttons.Init()
	a.DisabledFields = *a.DisabledFields.Init()
	return &a
}

type RoleQueryParam struct {
	common.PaginationParam
	ID           uint64 `form:"id"`            // 角色ID
	IDs          string `form:"ids"`           // 角色ID集合(逗号分隔)
	Name         string `form:"name"`          // 角色名称
	Type         int    `form:"type"`          // 角色类型
	Types        string `form:"types"`         // 角色类型集合(逗号分隔)
	AuditorTypes string `form:"auditor_types"` // 审核类型集合(逗号分隔)
	Status       int    `form:"status"`        // 状态(1:启用 2:禁用)
	Memo         string `form:"memo"`          // 备注(模糊查询)
	ShowDetails  bool   `form:"show_details"`  // 是否显示菜单、按钮、可禁用字段信息
	QueryValue   string `form:"query_value"`   // 模糊搜索(搜索 角色名称/备注)
	FindAll      bool   `form:"find_all"`      // 是否查找所有数据
}

type UpdateAuditorTypeParam struct {
	AuditorType int `json:"auditor_type"` // 审核类型
}

type Roles []*Role

type RoleQueryResult struct {
	Data       Roles
	PageResult *common.PaginationResult
}

func (a Roles) GetIDs() []uint64 {
	l := make([]uint64, len(a))
	for i, j := range a {
		l[i] = j.ID
	}
	return l
}

func (a Roles) Init() *Roles {
	items := make(Roles, 0)

	for _, role := range a {
		items = append(items, role.Init())
	}
	return &items
}
