package schema

import (
	"basic-frame/util/common"
	"time"
)

type Position struct {
	ID             uint64    `json:"id"`                                 // 唯一标识
	Name           string    `json:"name" binding:"required"`            // 职位名称
	RoleID         uint64    `json:"role_id"`                            // 职位的角色
	OrganizationID uint64    `json:"organization_id" binding:"required"` // 职位的组织ID
	Sequence       int       `json:"sequence"`                           // 排序值
	Status         int       `json:"status"`                             // 状态(1:启用 2:禁用)
	Memo           string    `json:"memo"`                               // 备注
	Creator        uint64    `json:"creator"`                            // 创建者
	CreatedAt      time.Time `json:"created_at"`                         // 创建时间
	UpdatedAt      time.Time `json:"updated_at"`                         // 更新时间

}

type PositionQueryParam struct {
	common.PaginationParam
	IDs            string `form:"ids"`             // 职位ID集合(逗号分隔)
	Name           string `form:"name"`            // 职位名称
	RoleID         uint64 `form:"role_id"`         // 职位的角色
	OrganizationID uint64 `form:"organization_id"` // 职位的组织ID
	Status         int    `form:"status"`          // 状态(1:启用 2:禁用)
	Memo           string `form:"memo"`            // 备注(模糊查询)
	SequenceSort   int    `form:"sequence_sort"`   // 按权重值排序(1:升序排序  2:降序排序)
	QueryValue     string `form:"queryValue"`      // 模糊搜索(模糊查询:职位名称、备注)
	FindAll        bool   `form:"find_all"`        // 是否查找所有数据
}

type Positions []*Position

type PositionQueryResult struct {
	Data       Positions
	PageResult *common.PaginationResult
}

func (a Positions) GetIDs() []uint64 {
	l := make([]uint64, len(a))
	for i, j := range a {
		l[i] = j.ID
	}
	return l
}

func (a Positions) GetRoleIds(roleIDs *[]uint64) {
	for _, position := range a {
		if !common.ContainsUint64(*roleIDs, position.RoleID) {
			*roleIDs = append(*roleIDs, position.RoleID)
		}
	}
}

func (a Positions) GetNames() []string {
	var names []string
	for _, position := range a {
		names = append(names, position.Name)
	}
	return names
}

func (a Positions) GetOrgIDs() []uint64 {
	var orgIDs []uint64

	for _, item := range a {
		if !common.ContainsUint64(orgIDs, item.OrganizationID) {
			orgIDs = append(orgIDs, item.OrganizationID)
		}
	}
	return orgIDs
}
