package schema

import (
	"basic-frame/util/common"
	"sort"
	"time"
)

type Organization struct {
	ID               uint64        `json:"id"`                      // 唯一标识
	Name             string        `json:"name" binding:"required"` // 组织名称
	RoleID           uint64        `json:"role_id"`                 // 组织的基础角色
	Sequence         int           `json:"sequence"`                // 排序值
	ParentID         *uint64       `json:"parent_id"`               // 父级组织ID
	Status           int           `json:"status"`                  // 状态(1:启用 2:禁用)
	Memo             string        `json:"memo"`                    // 备注
	Creator          uint64        `json:"creator"`                 // 创建者
	CreatedAt        time.Time     `json:"created_at"`              // 创建时间
	UpdatedAt        time.Time     `json:"updated_at"`              // 更新时间
	Positions        Positions     `json:"positions"`               // 组织的职位列表
	SonOrganizations Organizations `json:"son_organizations"`       // 下属组织列表
}

type OrganizationQueryParam struct {
	common.PaginationParam
	ID            uint64 `form:"id"`
	IDs           string `form:"ids"`     // 组织ID集合(逗号分隔)
	RoleID        uint64 `form:"role_id"` // 角色ID
	ParentID      uint64 `form:"parent_id"`
	Status        int    `form:"status"`         // 状态(1:启用 2:禁用)
	ShowPositions bool   `form:"show_positions"` // 是否显示组织的职位列表
	ShowDetails   bool   `form:"show_details"`   // 是否显示组织树
	QueryValue    string `form:"queryValue"`     // 模糊搜索(搜索 组织名称/备注)
	FindAll       bool   `form:"find_all"`       // 是否查询所有数据
}

// GetOrgIDs 获取组织的ID集合,包含子组织(已去重)
func (a Organization) GetOrgIDs(items *[]uint64) {
	if !common.ContainsUint64(*items, a.ID) {
		*items = append(*items, a.ID)
		if len(a.SonOrganizations) != 0 {
			for _, sonOrgInfo := range a.SonOrganizations {
				sonOrgInfo.GetOrgIDs(items)
			}
		}
	}
}

// SetCreator 设置组织、子组织、职位的创建人
func (a Organization) SetCreator(creator uint64) *Organization {
	if creator != 0 {
		a.Creator = creator
		if len(a.SonOrganizations) != 0 {
			for index, sonOrganization := range a.SonOrganizations {
				a.SonOrganizations[index] = sonOrganization.SetCreator(creator)
			}
		}
		if len(a.Positions) != 0 {
			for index, _ := range a.Positions {
				a.Positions[index].Creator = creator
			}
		}
	}
	return &a
}

func (a Organization) Init() *Organization {
	if a.ParentID == nil {
		parentID := uint64(0)
		a.ParentID = &parentID
	}
	if len(a.Positions) == 0 {
		a.Positions = make(Positions, 0)
	}
	if len(a.SonOrganizations) == 0 {
		a.SonOrganizations = make(Organizations, 0)
	} else {
		for index, sonOrganization := range a.SonOrganizations {
			a.SonOrganizations[index] = sonOrganization.Init()
		}
	}
	return &a
}

// GetRoleIds 获取组织、子组织的角色ID。isTop为True时，还包含当前组织的职位角色
func (a Organization) GetRoleIds(roleIDs *[]uint64, isTop bool) {
	if a.RoleID != 0 {
		*roleIDs = append(*roleIDs, a.RoleID)
	}
	if !isTop && len(a.Positions) != 0 {
		for _, position := range a.Positions {
			*roleIDs = append(*roleIDs, position.RoleID)
		}
	}
	if len(a.SonOrganizations) != 0 {
		for _, sonOrganization := range a.SonOrganizations {
			sonOrganization.GetRoleIds(roleIDs, false)
		}
	}
	return
}

// GetPositionIds 获取组织、子组织的职位ID集合
func (a Organization) GetPositionIds(positionIDs *[]uint64) {
	*positionIDs = append(*positionIDs, a.Positions.GetIDs()...)
	if len(a.SonOrganizations) != 0 {
		for _, sonOrganization := range a.SonOrganizations {
			sonOrganization.GetPositionIds(positionIDs)
		}
	}
	return
}

// SortOrganization 对组织、子组织的职位进行排序
func (a Organization) SortOrganization() *Organization {
	if a.ID != 0 {
		if len(a.Positions) != 0 {
			sort.SliceStable(a.Positions, func(i, j int) bool {
				return a.Positions[i].Sequence > a.Positions[j].Sequence
			})
		}
		if len(a.SonOrganizations) != 0 {
			for _, sonOrgInfo := range a.SonOrganizations {
				sonOrgInfo.SortOrganization()
			}
		}
	}
	return &a
}

type Organizations []*Organization

type OrganizationQueryResult struct {
	Data       Organizations
	PageResult *common.PaginationResult
}

func (a Organizations) GetIDs() []uint64 {
	l := make([]uint64, len(a))
	for i, j := range a {
		l[i] = j.ID
	}
	return l
}

// GetSonOrgIDs 获取子组织的组织ID集合(已去重)
func (a Organizations) GetSonOrgIDs() []uint64 {
	var sonOrgIDs []uint64
	for _, orgInfo := range a {
		for _, sonOrgInfo := range orgInfo.SonOrganizations {
			sonOrgInfo.GetOrgIDs(&sonOrgIDs)
		}
	}
	return sonOrgIDs
}

// SetCreator 设置组织、子组织、职位的创建人
func (a Organizations) SetCreator(creator uint64) *Organizations {
	if creator != 0 {
		for _, organization := range a {
			organization.SetCreator(creator)
		}
	}
	return &a
}

func (a Organizations) Init() *Organizations {
	items := make(Organizations, 0)

	for _, organization := range a {
		items = append(items, organization.Init())
	}
	return &items
}

// GetRoleIds 获取组织、子组织、当前组织职位的角色ID
func (a Organizations) GetRoleIds(roleIDs *[]uint64) {
	for _, organization := range a {
		organization.GetRoleIds(roleIDs, true)
	}
}

// SortOrganizations 对所有数据进行排序
func (a Organizations) SortOrganizations() *Organizations {
	if len(a) != 0 {
		for index, OrgInfo := range a {
			if len(OrgInfo.SonOrganizations) != 0 {
				OrgInfo.SonOrganizations.SortOrganizations()
			}
			if len(a[index].Positions) != 0 {
				sort.SliceStable(a[index].Positions, func(i, j int) bool {
					return a[index].Positions[i].Sequence > a[index].Positions[j].Sequence
				})
			}
		}
		sort.SliceStable(a, func(i, j int) bool {
			return a[i].Sequence > a[j].Sequence
		})
	}
	return &a
}

// OrganizationTree TODO: 等待用户表完成
type OrganizationTree struct {
	ID        uint64    `json:"id"`                      // 唯一标识
	Name      string    `json:"name" binding:"required"` // 组织名称
	RoleID    uint64    `json:"role_id"`                 // 组织的基础角色
	Sequence  int       `json:"sequence"`                // 排序值
	ParentID  *uint64   `json:"parent_id"`               // 父级组织ID
	Status    int       `json:"status"`                  // 状态(1:启用 2:禁用)
	Memo      string    `json:"memo"`                    // 备注
	Creator   uint64    `json:"creator"`                 // 创建者
	CreatedAt time.Time `json:"created_at"`              // 创建时间
	UpdatedAt time.Time `json:"updated_at"`              // 更新时间
	Positions Positions `json:"positions"`               // 组织的职位列表
	//Users                Users             `json:"users"`                   // 组织用户列表
	SonOrganizationTrees OrganizationTrees `json:"son_organization_trees"` // 下属组织列表
}

type OrganizationTrees []*OrganizationTree

func (a Organization) ToSchemaOrgTree() *OrganizationTree {
	item := new(OrganizationTree)
	_ = common.Copy(a, item)
	if len(a.SonOrganizations) != 0 {
		for _, sonOrgInfo := range a.SonOrganizations {
			item.SonOrganizationTrees = append(item.SonOrganizationTrees, sonOrgInfo.ToSchemaOrgTree())
		}
	}
	return item.Init()
}

func (a Organizations) ToSchemaOrgTrees() *OrganizationTrees {
	OrgTrees := new(OrganizationTrees)
	for _, organization := range a {
		*OrgTrees = append(*OrgTrees, organization.ToSchemaOrgTree())
	}
	return OrgTrees
}

func (a OrganizationTree) Init() *OrganizationTree {
	if a.ParentID == nil {
		parentID := uint64(0)
		a.ParentID = &parentID
	}
	if len(a.Positions) == 0 {
		a.Positions = make(Positions, 0)
	}
	// TODO: 等待用户表完成
	//if len(a.Users) == 0 {
	//	a.Users = make(Users, 0)
	//}
	if len(a.SonOrganizationTrees) == 0 {
		a.SonOrganizationTrees = make(OrganizationTrees, 0)
	} else {
		for index, sonOrganization := range a.SonOrganizationTrees {
			a.SonOrganizationTrees[index] = sonOrganization.Init()
		}
	}
	return &a
}

// TODO: 等待用户表完成
//func (a OrganizationTrees) AddUserToTree(users Users) {
//	for _, OrgTreeInfo := range a {
//		for _, userInfo := range users {
//			newUser := OrgTreeInfo.AddUserToTree(*userInfo)
//			if newUser.ID != 0 {
//				OrgTreeInfo.Users = append(OrgTreeInfo.Users, &newUser)
//			}
//		}
//		if len(OrgTreeInfo.SonOrganizationTrees) != 0 {
//			OrgTreeInfo.SonOrganizationTrees.AddUserToTree(users)
//		}
//	}
//}

// TODO: 等待用户表完成
//func (a OrganizationTree) AddUserToTree(user User) User {
//	var newUser User
//	for _, organizationInfo := range user.Organizations {
//		if organizationInfo.ID == a.ID {
//			newUser = user
//			newUser.Positions = make(Positions, 0)
//		}
//	}
//	if newUser.ID != 0 {
//		for _, positionInfo := range user.Positions {
//			if positionInfo.OrganizationID == a.ID {
//				newUser.Positions = append(newUser.Positions, positionInfo)
//			}
//		}
//	}
//	return newUser
//}
