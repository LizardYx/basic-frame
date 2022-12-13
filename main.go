package main

import (
	"basic-frame/config"
	baseBll "basic-frame/modules/base/bll"
	casbin_adapter "basic-frame/util/casbin-adapter"
	"basic-frame/util/logger"
	"basic-frame/util/mysql"
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
)

func main() {
	// 加载系统配置文件
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

	// 初始化Casbin
	if err := casbin_adapter.InitCasbin(); err != nil {
		logger.Log.Warningf("%v", err)
		fmt.Printf("%v\n", err)
	}

	// 初始化菜单数据
	go baseBll.MenuBll.InitData()

	// 注册定时任务
	RegisterCronTask(context.Background())

	// 初始化Webserver
	InitWebServer()
}

// RegisterCronTask 注册定时任务
func RegisterCronTask(ctx context.Context) {
	c := cron.New()

	// 任务调度
	//if _, err := c.AddFunc("@every 3s", Task); err != nil {
	//	panic(err)
	//}

	c.Start()
}
