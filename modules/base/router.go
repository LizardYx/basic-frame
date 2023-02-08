package base

import (
	"basic-frame/modules/base/api"
	"github.com/gin-gonic/gin"
)

// RouterBase 模块路由实例,需要将该实例注册到全局路由实例中
var RouterBase = &Router{
	SystemConfigApi:   api.SystemConfigApi,
	TagManageApi:      api.TagManageApi,
	SecurityLevelApi:  api.SecurityLevelApi,
	DisabledFieldApi:  api.DisabledFieldApi,
	RestfulApi:        api.RestfulApiApi,
	ButtonApi:         api.ButtonApi,
	MenuApi:           api.MenuApi,
	RoleApi:           api.RoleApi,
	PositionApi:       api.PositionApi,
	OrganizationApi:   api.OrganizationApi,
	UserExtendInfoApi: api.UserExtendInfoApi,
	UserGroupApi:      api.UserGroupApi,
	UserApi:           api.UserApi,
}

// Router 模块路由管理结构体,需要将该结构体注册到全局路由结构体
type Router struct {
	SystemConfigApi   *api.SystemConfig
	TagManageApi      *api.TagManage
	SecurityLevelApi  *api.SecurityLevel
	DisabledFieldApi  *api.DisabledField
	RestfulApi        *api.RestfulApi
	ButtonApi         *api.Button
	MenuApi           *api.Menu
	RoleApi           *api.Role
	PositionApi       *api.Position
	OrganizationApi   *api.Organization
	UserExtendInfoApi *api.UserExtendInfo
	UserGroupApi      *api.UserGroup
	UserApi           *api.User
}

// Register 模块路由注册
func (a *Router) Register(api *gin.RouterGroup) {

	api.POST("login", a.UserApi.Login)
	api.POST("logout", a.UserApi.Logout)
	api.GET("refresh-token", a.UserApi.RefreshToken)

	UserGroup := api.Group("base/users")
	{
		UserGroup.GET("", a.UserApi.Query)
		UserGroup.GET(":id", a.UserApi.Get)
		UserGroup.POST("", a.UserApi.Create)
		UserGroup.PUT(":id", a.UserApi.Update)
		UserGroup.PUT(":id/enable", a.UserApi.EnableUser)
		UserGroup.PUT(":id/disabled", a.UserApi.DisabledUser)
		UserGroup.PUT("", a.UserApi.BatchUpdateUserPermission)
		UserGroup.PUT(":id/permission", a.UserApi.UpdateUserPermission)
		UserGroup.DELETE(":id", a.UserApi.Delete)
	}

	UserGroupsGroup := api.Group("base/user-groups")
	{
		UserGroupsGroup.GET("", a.UserGroupApi.Query)
		UserGroupsGroup.GET(":id", a.UserGroupApi.Get)
		UserGroupsGroup.GET(":id/users", a.UserGroupApi.GetUserGroupUsers)
		UserGroupsGroup.POST("", a.UserGroupApi.Create)
		UserGroupsGroup.PUT(":id", a.UserGroupApi.Update)
		UserGroupsGroup.PUT(":id/user-join", a.UserGroupApi.UserJoinUserGroup)
		UserGroupsGroup.PUT(":id/user-remove", a.UserGroupApi.UserGroupRemoveUser)
		UserGroupsGroup.DELETE(":id", a.UserGroupApi.Delete)
	}

	UserExtendInfoGroup := api.Group("base/user-extend-infos")
	{
		UserExtendInfoGroup.PUT(":id", a.UserExtendInfoApi.Update)
	}

	OrganizationGroup := api.Group("base/organizations")
	{
		OrganizationGroup.GET("", a.OrganizationApi.Query)
		OrganizationGroup.GET("/organization-tree/basic-info/:id", a.OrganizationApi.Get)
		OrganizationGroup.GET("/organization-tree/create-notifications", a.OrganizationApi.GetOrgTreeForCreateNotifications)
		OrganizationGroup.GET("/organization-tree/create-user", a.OrganizationApi.GetOrganizationTreeForCreateUser)
		OrganizationGroup.GET("/organization-tree/all", a.OrganizationApi.GetOrganizationTreeWithUser)
		OrganizationGroup.GET("/organization-tree/edit", a.OrganizationApi.GetOrganizationTree)
		OrganizationGroup.POST("", a.OrganizationApi.Create)
		OrganizationGroup.PUT(":id/user-join", a.OrganizationApi.UserJoinOrganization)
		OrganizationGroup.PUT(":id/user-remove", a.OrganizationApi.UserRemoveOrganization)
		OrganizationGroup.PUT(":id", a.OrganizationApi.Update)
		api.PUT("base/organizations", a.OrganizationApi.UpdateOrganizations)
		OrganizationGroup.DELETE(":id", a.OrganizationApi.Delete)
	}

	PositionGroup := api.Group("base/positions")
	{
		PositionGroup.GET("", a.PositionApi.Query)
		PositionGroup.GET(":id", a.PositionApi.Get)
		PositionGroup.POST("", a.PositionApi.Create)
		PositionGroup.PUT(":id", a.PositionApi.Update)
		PositionGroup.PUT(":id/user-join", a.PositionApi.PositionAddUser)
		PositionGroup.PUT(":id/user-remove", a.PositionApi.PositionRemoveUser)
		PositionGroup.DELETE(":id", a.PositionApi.Delete)
	}

	SecurityLevelGroup := api.Group("base/security-levels")
	{
		SecurityLevelGroup.GET("", a.SecurityLevelApi.Query)
		SecurityLevelGroup.GET(":id", a.SecurityLevelApi.Get)
		SecurityLevelGroup.POST("", a.SecurityLevelApi.Create)
		SecurityLevelGroup.PUT(":id", a.SecurityLevelApi.Update)
		SecurityLevelGroup.DELETE(":id", a.SecurityLevelApi.Delete)
	}

	RoleGroup := api.Group("base/roles")
	{
		RoleGroup.GET("", a.RoleApi.Query)
		RoleGroup.GET(":id", a.RoleApi.Get)
		RoleGroup.GET(":id/permission-tree", a.RoleApi.GetRolePermissionTree)
		RoleGroup.POST("", a.RoleApi.Create)
		RoleGroup.PUT(":id", a.RoleApi.Update)
		RoleGroup.PUT(":id/auditor_type", a.RoleApi.UpdateAuditorType)
		RoleGroup.PUT(":id/user-join", a.RoleApi.UserAddRole)
		RoleGroup.PUT(":id/user-remove", a.RoleApi.UserRemoveRole)
		RoleGroup.DELETE(":id", a.RoleApi.Delete)
	}

	MenuGroup := api.Group("base/permission-tree")
	{
		MenuGroup.GET("/create-role", a.MenuApi.GetPermissionTreeForCreateRole)
		MenuGroup.GET("/edit", a.MenuApi.GetPermissionTree)
		MenuGroup.GET("/download", a.MenuApi.DownloadPermissionTree)
		MenuGroup.PUT("", a.MenuApi.UpdatePermissionTree)
		MenuGroup.PUT("/basic-info-update", a.MenuApi.BatchUpdateMenus)
		MenuGroup.DELETE("/:id", a.MenuApi.Delete)
	}

	ButtonGroup := api.Group("base/button")
	{
		ButtonGroup.POST("", a.ButtonApi.Create)
		ButtonGroup.PUT(":id", a.ButtonApi.Update)
		ButtonGroup.PUT("", a.ButtonApi.BatchUpdate)
		ButtonGroup.DELETE(":id", a.ButtonApi.Delete)
	}

	DisabledFieldGroup := api.Group("base/disabled_field")
	{
		DisabledFieldGroup.DELETE(":id", a.DisabledFieldApi.Delete)
	}

	TagManageGroup := api.Group("base/tag-manages")
	{
		TagManageGroup.GET("", a.TagManageApi.Query)
		TagManageGroup.GET(":id", a.TagManageApi.Get)
		TagManageGroup.POST("", a.TagManageApi.Create)
		TagManageGroup.PUT(":id", a.TagManageApi.Update)
		TagManageGroup.DELETE(":id", a.TagManageApi.Delete)
	}

	SystemConfigGroup := api.Group("base/system-config")
	{
		SystemConfigGroup.GET("", a.SystemConfigApi.Query)
		SystemConfigGroup.PUT(":id", a.SystemConfigApi.Update)
	}
}
