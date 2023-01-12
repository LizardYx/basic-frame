package schema

import (
	"basic-frame/util/common"
	"sort"
	"time"
)

type SecurityLevel struct {
	ID        uint64    `json:"id"`                      // 唯一标识
	Name      string    `json:"name" binding:"required"` // 安全级别名称
	Sequence  int       `json:"sequence"`                // 排序值
	Status    int       `json:"status"`                  // 状态(1:启用 2:禁用)
	Memo      string    `json:"memo"`                    // 备注
	Creator   uint64    `json:"creator"`                 // 创建者
	CreatedAt time.Time `json:"created_at"`              // 创建时间
	UpdatedAt time.Time `json:"updated_at"`              // 更新时间
	Roles     Roles     `json:"roles"`                   // 安全级别绑定的角色
}

func (a SecurityLevel) Init() *SecurityLevel {
	if len(a.Roles) == 0 {
		a.Roles = make(Roles, 0)
	} else {
		for index, _ := range a.Roles {
			a.Roles[index].Init()
		}
	}
	return &a
}

type SecurityLevelQueryParam struct {
	common.PaginationParam
	ID          uint64 `form:"id"`           // 安全级别ID
	IDs         string `form:"ids"`          // 安全级别ID集合(逗号分隔)
	RoleIDs     string `form:"role_ids"`     // 角色ID集合(逗号分隔)
	Name        string `form:"name"`         // 安全级别名称(模糊查询)
	Status      int    `form:"status"`       // 状态(1:启用 2:禁用)
	ShowDetails bool   `form:"show_details"` // 是否显示角色信息
	QueryValue  string `form:"queryValue"`   // 模糊搜索(搜索 安全级别名称/备注)
	FindAll     bool   `form:"find_all"`     // 是否查找所有数据
}

type SecurityLevels []*SecurityLevel

type SecurityLevelQueryResult struct {
	Data       SecurityLevels
	PageResult *common.PaginationResult
}

func (a SecurityLevels) GetIDs() []uint64 {
	l := make([]uint64, len(a))
	for i, j := range a {
		l[i] = j.ID
	}
	return l
}

// Sort 通过Sequence对安全等级进行排序
func (a SecurityLevels) Sort() *SecurityLevels {
	sort.SliceStable(a, func(i, j int) bool {
		return a[i].Sequence > a[j].Sequence
	})
	return &a
}
