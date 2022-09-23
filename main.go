package main

import (
	"basic-frame/config"
	"log"
)

func main() {
	// 加载系统配置文件
	if errMessage, err := config.LoadSystemConfigFile(); err != nil {
		log.Fatal(errMessage, err)
	}
	// 初始化日志中间件

	// 初始化i18n

	// DB连接

	select {}
}
