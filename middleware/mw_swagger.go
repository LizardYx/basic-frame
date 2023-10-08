package middleware

import (
	"basic-frame/docs"
	"basic-frame/util/common"
	"fmt"
)

func SwaggerMiddleware() {
	docs.SwaggerInfo.Title = common.SysConfig.AppName
	docs.SwaggerInfo.Version = fmt.Sprintf("%.5f", common.SysConfig.AppVersion)
	docs.SwaggerInfo.Host = common.SysConfig.WebServer.Host
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}
