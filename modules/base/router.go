package base

import (
	"basic-frame/modules/base/api"
	"github.com/gin-gonic/gin"
)

// RouterBase 模块路由实例,需要将该实例注册到全局路由实例中
var RouterBase = &Router{
	SystemConfigApi:  api.SystemConfigApi,
	TagManageApi:     api.TagManageApi,
	SecurityLevelApi: api.SecurityLevelApi,
	DisabledFieldApi: api.DisabledFieldApi,
	RestfulApi:       api.RestfulApiApi,
	ButtonApi:        api.ButtonApi,
	MenuApi:          api.MenuApi,
	RoleApi:          api.RoleApi,
}

// Router 模块路由管理结构体,需要将该结构体注册到全局路由结构体
type Router struct {
	SystemConfigApi  *api.SystemConfig
	TagManageApi     *api.TagManage
	SecurityLevelApi *api.SecurityLevel
	DisabledFieldApi *api.DisabledField
	RestfulApi       *api.RestfulApi
	ButtonApi        *api.Button
	MenuApi          *api.Menu
	RoleApi          *api.Role
}

// Register 模块路由注册
func (a *Router) Register(api *gin.RouterGroup) {

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
		ButtonGroup.PUT(":id/batch-update", a.ButtonApi.BatchUpdate)
		ButtonGroup.DELETE(":id", a.ButtonApi.Delete)
	}

	DisabledFieldGroup := api.Group("base/disabled_field")
	{
		DisabledFieldGroup.DELETE(":id", a.DisabledFieldApi.Delete)
	}

	SecurityLevelGroup := api.Group("base/security-levels")
	{
		SecurityLevelGroup.GET("", a.SecurityLevelApi.Query)
		SecurityLevelGroup.GET(":id", a.SecurityLevelApi.Get)
		SecurityLevelGroup.POST("", a.SecurityLevelApi.Create)
		SecurityLevelGroup.PUT(":id", a.SecurityLevelApi.Update)
		SecurityLevelGroup.DELETE(":id", a.SecurityLevelApi.Delete)
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
