package base

import (
	"basic-frame/modules/base/api"
	"github.com/gin-gonic/gin"
)

// RouterBase 模块路由实例,需要将该实例注册到全局路由实例中
var RouterBase = &Router{
	SystemConfigAPI: api.SystemConfigApi,
	ButtonAPI:       api.ButtonAPI,
	RestfulAPI:      api.ResApi,
}

// Router 模块路由管理结构体,需要将该结构体注册到全局路由结构体
type Router struct {
	SystemConfigAPI *api.SystemConfig
	ButtonAPI       *api.Button
	RestfulAPI      *api.RestfulApi
}

// Register 模块路由注册
func (a *Router) Register(api *gin.RouterGroup) {
	SystemConfigGroup := api.Group("system-config")
	{
		SystemConfigGroup.GET("", a.SystemConfigAPI.Query)
		SystemConfigGroup.POST("", a.SystemConfigAPI.Create)
	}

	buttonGroup := api.Group("button")
	{
		buttonGroup.POST("", a.ButtonAPI.Create)
	}

	ResApiGroup := api.Group("res-api")
	{
		ResApiGroup.POST("", a.RestfulAPI.Create)
	}
}
