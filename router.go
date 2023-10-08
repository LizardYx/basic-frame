package main

import (
	"basic-frame/middleware"
	"basic-frame/modules/base"
	"basic-frame/util/common"
	"basic-frame/util/consts"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// 全局路由实例
var globalRouter = Router{
	BaseRouter: base.RouterBase,
}

// Router 全局路由管理
type Router struct {
	BaseRouter *base.Router
}

// Register 注册路由
func (a *Router) Register(app *gin.Engine) error {
	a.RegisterAPI(app)
	return nil
}

// RegisterAPI 注册各个模块的API
func (a *Router) RegisterAPI(app *gin.Engine) {
	g := app.Group("/api")
	// jwtauth 登陆认证
	g.Use(middleware.UserJWTAuth(
		middleware.AllowPathPrefixSkipper("/api/v1/login"),
		middleware.AllowPathPrefixSkipper("/api/v1/base/permission-tree/download"),
	))

	// Casbin 权限认证
	g.Use(middleware.CasbinMiddleware(
		common.SysConfig.CasbinSyncEnforcer,
		middleware.AllowPathPrefixSkipper("/api/v1/login"),
		middleware.AllowPathPrefixSkipper("/api/v1/base/permission-tree/download"),
	))

	// 增加请求频率限制中间件
	// TODO:

	v1 := g.Group("/v1")
	{
		a.BaseRouter.Register(v1)
	}
	if common.SysConfig.RunMode != consts.RunModeRelease {
		// swagger/index.html
		app.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}
}
