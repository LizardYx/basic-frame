package schema

import (
	"basic-frame/util/common"
	"time"
)

type UserGroup struct {
	ID        uint64    `json:"id"`                      // 唯一标识
	Name      string    `json:"name" binding:"required"` // 用户组名称
	RoleID    uint64    `json:"role_id"`                 // 用户组的基础角色
	Sequence  int       `json:"sequence"`                // 排序值
	Status    int       `json:"status"`                  // 状态(1:启用 2:禁用)
	Memo      string    `json:"memo"`                    // 备注
	Creator   uint64    `json:"creator"`                 // 创建者
	CreatedAt time.Time `json:"created_at"`              // 创建时间
	UpdatedAt time.Time `json:"updated_at"`              // 更新时间
}

type UserGroupQueryParam struct {
	common.PaginationParam
	IDs        string `form:"ids"`         // 用户组ID集合(逗号分隔)
	Name       string `form:"name"`        // 用户组名称
	RoleID     uint64 `form:"role_id"`     // 用户组的角色
	Status     int    `form:"status"`      // 状态(1:启用 2:禁用)
	Memo       string `form:"memo"`        // 备注(模糊查询)
	QueryValue string `form:"query_value"` // 模糊搜索(模糊查询:用户组名称、备注)
	FindAll    bool   `form:"find_all"`    // 是否查找所有
}

type UserGroupQueryResult struct {
	Data       UserGroups
	PageResult *common.PaginationResult
}

type UserGroups []*UserGroup

func (a UserGroups) GetIDs() []uint64 {
	l := make([]uint64, len(a))
	for i, j := range a {
		l[i] = j.ID
	}
	return l
}

func (a UserGroups) GetRoleIds(roleIDs *[]uint64) {
	for _, userGroup := range a {
		if !common.ContainsUint64(*roleIDs, userGroup.RoleID) {
			*roleIDs = append(*roleIDs, userGroup.RoleID)
		}
	}
}
