package schema

import (
	"basic-frame/util/common"
	"sort"
	"time"
)

// TODO: 等待用户组和用户扩展信息完成
type User struct {
	ID        uint64    `json:"id"`                           // 唯一标识
	UserName  string    `json:"user_name" binding:"required"` // 用户名称
	Password  string    `json:"password"`                     // 用户密码
	Status    int       `json:"status"`                       // 状态(1:启用 2:禁用)
	Sequence  int       `json:"sequence"`                     // 排序值
	Creator   uint64    `json:"creator"`                      // 创建者
	CreatedAt time.Time `json:"created_at"`                   // 创建时间
	UpdatedAt time.Time `json:"updated_at"`                   // 更新时间
	//ExtendInfo    UserExtendInfo `json:"extend_info"`                  // 用户扩展信息
	Organizations Organizations `json:"organizations"` // 用户所属的组织结构
	Positions     Positions     `json:"positions"`     // 用户所属的职位
	Roles         Roles         `json:"roles"`         // 用户所属的角色(直接绑定给用户的角色)
	//UserGroups    UserGroups     `json:"user_groups"`                  // 用户所属的的用户组
}

type UserQueryParam struct {
	common.PaginationParam
	ID             uint64 `form:"id"`               // 用户ID
	IDs            string `form:"ids"`              // 用户ID集合(逗号分隔)
	UserName       string `form:"user_name"`        // 用户名称
	Status         int    `form:"status"`           // 状态(1:启用 2:禁用)
	OrgID          uint64 `form:"org_id"`           // 组织ID
	OrgIDs         string `form:"org_ids"`          // 组织ID集合(逗号分隔)
	PositionID     uint64 `form:"position_id"`      // 职位ID
	PositionIDs    string `form:"position_ids"`     // 职位ID集合(逗号分隔)
	AuditorTypes   string `form:"auditor_types"`    // 审核类型(逗号分隔)
	RoleID         uint64 `form:"role_id"`          // 角色ID
	RoleIDs        string `form:"role_ids"`         // 角色ID集合(逗号分隔)
	UserGroupID    uint64 `form:"user_group_id"`    // 用户组ID
	UserGroupIDs   string `form:"user_group_ids"`   // 用户组ID集合(逗号分隔)
	ShowExtendInfo bool   `form:"show_extend_info"` // 是否显示扩赞信息
	ShowDetails    bool   `form:"show_details"`     // 是否显示详情
	OmitPassword   bool   `form:"omit_password"`    // 是否隐藏密码
	SequenceSort   int    `form:"sequence_sort"`    // 按权重排序(1:升序排序 2:降序排序)
	QueryValue     string `form:"query_value"`      // 模糊搜索(搜索 用户名称/用户昵称/移动手机/QQ账号/邮箱账号)
	FindAll        bool   `form:"find_all"`         // 是否查找所有
	FindDeleted    bool   `form:"find_deleted"`     // 是否查找已经删除的用户
}

type SubUserQueryParam struct {
	CurrentOrg bool `form:"current_org"` // 是否包含当前组织的用户
	SonOrg     bool `form:"son_org"`     // 是否包含子组织的用户
}

type LoginParam struct {
	UserName string `json:"user_name" binding:"required"` // 用户名
	Password string `json:"password"  binding:"required"` // 密码(md5加密)
}

type UserQueryResult struct {
	Data       Users
	PageResult *common.PaginationResult
}

// UpdatePasswordParam 更新密码请求参数
type UpdatePasswordParam struct {
	OldPassword string `json:"old_password" binding:"required"` // 旧密码(md5加密)
	NewPassword string `json:"new_password" binding:"required"` // 新密码(md5加密)
}

// LoginTokenInfo 登录令牌信息
type LoginTokenInfo struct {
	AccessToken string `json:"access_token"` // 访问令牌
	TokenType   string `json:"token_type"`   // 令牌类型
	ExpiresAt   int64  `json:"expires_at"`   // 令牌到期时间戳
}

type LoginRes struct {
	UserInfo  User           `json:"user_info"`
	TokenInfo LoginTokenInfo `json:"token_info"`
}

func (a User) Init() *User {
	if len(a.Organizations) == 0 {
		a.Organizations = make(Organizations, 0)
	} else {
		a.Organizations = *a.Organizations.Init()
	}

	if len(a.Positions) == 0 {
		a.Positions = make(Positions, 0)
	}

	if len(a.Roles) == 0 {
		a.Roles = make(Roles, 0)
	} else {
		a.Roles = *a.Roles.Init()
	}

	// TODO: 等待用户组完成
	//if len(a.UserGroups) == 0 {
	//	a.UserGroups = make(UserGroups, 0)
	//}
	return &a
}

type Users []*User

func (a Users) ToMap() map[uint64]*User {
	m := make(map[uint64]*User)
	for _, i := range a {
		m[i.ID] = i
	}
	return m
}

func (a Users) GetIDs() []uint64 {
	l := make([]uint64, len(a))
	for i, j := range a {
		l[i] = j.ID
	}
	return l
}

func (a Users) Init() *Users {
	items := make(Users, 0)
	for _, user := range a {
		items = append(items, user.Init())
	}
	return &items
}

func (a Users) GetUserNames() []string {
	var userNames []string

	for _, user := range a {
		userNames = append(userNames, user.UserName)
	}
	return userNames
}

func (a Users) GetUserRealNames() []string {
	var realNames []string

	// TODO: 等待扩展信息完成
	//for _, user := range a {
	//	realNames = append(realNames, user.ExtendInfo.RealName)
	//}
	return realNames
}

// SortUser 按权重排序用户
func (a Users) SortUser() *Users {
	sort.Slice(a, func(i, j int) bool {
		return a[i].Sequence > a[j].Sequence
	})
	return &a
}
