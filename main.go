package main

import (
	"basic-frame/config"
	baseBll "basic-frame/modules/base/bll"
	"basic-frame/util/logger"
	"basic-frame/util/mysql"
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
)

func main() {
	// 加载系统配置文件
	// 初始化日志中间件
	// 连接DB
	// 读取并更新DB中的系统基础配置信息
	// 读取DB中的菜单版本、系统默认语言、SMTP服务信息、到缓存中
	// 初始化超管用户
	if errMessage, err := config.LoadSystemConfigFile(); err != nil {
		log.Fatal(errMessage, err)
	}

	// 初始化日志中间件
	logger.InitLogger()

	// DB连接
	mysql.InitMysql()

	// 初始化系统基础配置
	if err := baseBll.SystemConfigBll.Init(); err != nil {
		logger.Log.Warningf("%v", err)
		fmt.Printf("%v\n", err)
	}

	// 注册定时任务
	RegisterCronTask(context.Background())

	// 初始化Webserver
	InitWebServer()
}

// RegisterCronTask 注册定时任务
func RegisterCronTask(ctx context.Context) {
	c := cron.New()

	// 读取菜单文件
	// TODO:

	// 任务调度
	//if _, err := c.AddFunc("@every 3s", Task); err != nil {
	//	panic(err)
	//}

	c.Start()
}
