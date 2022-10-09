package main

import (
	"basic-frame/config"
	"basic-frame/middleware"
	"basic-frame/util/consts"
	"basic-frame/util/i18n/Localizer"
	"basic-frame/util/logger"
	"basic-frame/util/mysql"
	"log"
)

func main() {
	// 初始化i18n
	Localizer.I18n = Localizer.GetLocalizer(consts.DefaultLang)
	// 加载系统配置文件
	if errMessage, err := config.LoadSystemConfigFile(); err != nil {
		log.Fatal(errMessage, err)
	}

	// 初始化日志中间件
	logger.InitLogger()

	// DB连接
	mysql.InitMysql()

	go middleware.Manager.Start()

	select {}
}
