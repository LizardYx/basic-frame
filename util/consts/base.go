package consts

// Base模块状态控制
const (
	BaseStatusEnable   = 1 // 启动状态
	BaseStatusDisabled = 2 // 禁用状态
)

// 角色类型
const (
	RoleTypeForUser      = 1 // 用户使用的角色
	RoleTypeForOrg       = 2 // 组织使用的角色
	RoleTypeForPosition  = 3 // 职位使用的角色
	RoleTypeForUserGroup = 4 // 用户组使用的角色
)

var RoleTypes = []uint64{RoleTypeForUser, RoleTypeForOrg, RoleTypeForPosition, RoleTypeForUserGroup}
