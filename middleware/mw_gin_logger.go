package middleware

import (
	"basic-frame/util/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// GinLoggerMiddleware gin的日志中间件，记录restful请求信息
func GinLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		fields := make(map[string]interface{})
		start := time.Now()
		method := c.Request.Method
		if method == http.MethodPost || method == http.MethodPut {
			// TODO: 还需要记录这2种请求的参数
		}
		c.Next()
		timeConsuming := time.Since(start).Nanoseconds() / 1e6
		fields["user_agent"] = c.GetHeader("User-Agent")
		logger.Log.WithFields(fields).Infof("%s(%s) | %d(%dms)", c.Request.RequestURI, method, c.Writer.Status(), timeConsuming)
	}
}
