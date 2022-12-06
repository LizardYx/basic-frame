package main

import (
	"basic-frame/middleware"
	"basic-frame/util/common"
	"basic-frame/util/consts"
	"basic-frame/util/i18n/Localizer"
	"basic-frame/util/logger"
	"context"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// InitGinEngine 初始化gin引擎
func InitGinEngine() *gin.Engine {
	// 设置gin运行模式
	gin.SetMode(common.SysConfig.GetGinRunModel())

	// 初始化gin基础配置
	app := gin.Default()

	// 要在路由组之前全局使用「跨域中间件」, 否则OPTIONS会返回404
	app.Use(middleware.Cors())

	// 初始化Webserver日志记录
	app.Use(middleware.GinLoggerMiddleware())

	// 初始化i18n国际化
	Localizer.I18n = Localizer.GetLocalizer(consts.DefaultLang)

	// 注册模块路由
	if err := globalRouter.Register(app); err != nil {
		logger.Log.Fatalf("router register error: %v", err)
	}
	return app
}

// InitWebServer WebServer初始化
func InitWebServer() {
	webServer := &http.Server{
		Addr:         common.SysConfig.WebServer.GetWebServerAddr(),
		Handler:      InitGinEngine(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	// 初始化Websocket中间件
	go middleware.Manager.Start()

	// 根据配置文件启动Webserver
	go func() {
		if common.SysConfig.WebServer.HttpsMode {
			webServer.ListenAndServeTLS(common.SysConfig.WebServer.HttpsCrtFile, common.SysConfig.WebServer.HttpsKeyFile)
		} else {
			webServer.ListenAndServe()
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	logger.Log.Info("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := webServer.Shutdown(ctx); err != nil {
		logger.Log.Fatal("Server Shutdown Failed: ", err)
	}
	logger.Log.Info("Server exiting")
}
