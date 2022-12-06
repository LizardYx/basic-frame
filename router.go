package main

import (
	"basic-frame/middleware"
	"basic-frame/modules/base"
	"github.com/gin-gonic/gin"
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
	g.Use(middleware.UserJWTAuth(middleware.AllowMethodAndPathPrefixSkipper("PUT", "/api/v1/system-config")))
	// jwtauth
	// TODO:

	// Casbin
	// TODO:

	// 增加请求频率限制中间件
	// TODO:

	v1 := g.Group("/v1")
	{
		a.BaseRouter.Register(v1)
	}
}
