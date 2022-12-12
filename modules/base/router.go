package base

import (
	"basic-frame/modules/base/api"
	"github.com/gin-gonic/gin"
)

// RouterBase 模块路由实例,需要将该实例注册到全局路由实例中
var RouterBase = &Router{
	SystemConfigAPI:  api.SystemConfigApi,
	DisabledFieldApi: api.DisabledFieldApi,
	RestfulApi:       api.RestfulApiApi,
	ButtonApi:        api.ButtonApi,
}

// Router 模块路由管理结构体,需要将该结构体注册到全局路由结构体
type Router struct {
	SystemConfigAPI  *api.SystemConfig
	DisabledFieldApi *api.DisabledField
	RestfulApi       *api.RestfulApi
	ButtonApi        *api.Button
}

// Register 模块路由注册
func (a *Router) Register(api *gin.RouterGroup) {

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

	SystemConfigGroup := api.Group("base/system-config")
	{
		SystemConfigGroup.GET("", a.SystemConfigAPI.Query)
		SystemConfigGroup.PUT(":id", a.SystemConfigAPI.Update)
	}
}
