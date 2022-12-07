package base

import (
	"basic-frame/modules/base/api"
	"github.com/gin-gonic/gin"
)

// RouterBase 模块路由实例,需要将该实例注册到全局路由实例中
var RouterBase = &Router{
	SystemConfigAPI:  api.SystemConfigApi,
	RestfulApi:       api.RestfulApiApi,
	DisabledFieldApi: api.DisabledFieldApi,
}

// Router 模块路由管理结构体,需要将该结构体注册到全局路由结构体
type Router struct {
	SystemConfigAPI  *api.SystemConfig
	RestfulApi       *api.RestfulApi
	DisabledFieldApi *api.DisabledField
}

// Register 模块路由注册
func (a *Router) Register(api *gin.RouterGroup) {
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
