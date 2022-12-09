package middleware

import (
	"basic-frame/util/consts"
	"basic-frame/util/ginx"
	"basic-frame/util/ginx/errors"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// CasbinMiddleware 接口权限检查
func CasbinMiddleware(enforcer *casbin.SyncedEnforcer, skippers ...SkipperFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 是否为不做权限检查的接口
		if SkipHandler(c, skippers...) {
			c.Next()
			return
		}

		// 如果是超管账号
		if ginx.GetUserName(c) == consts.AdminName {
			// 超管账号,不验证权限
			c.Next()
			return
		}

		// 是否有匹配项
		p := c.Request.URL.Path
		m := c.Request.Method
		userID := ginx.GetUserID(c)
		if ok, err := enforcer.Enforce(strconv.FormatUint(userID, 10), p, m); err != nil {
			ginx.ResError(c, "", errors.WithStack(err))
			c.Abort()
			return
		} else if !ok {
			ginx.ResError(c, "", errors.NewIANAResponse(http.StatusUnauthorized, ""))
			c.Abort()
			return
		}
		c.Next()
	}
}
