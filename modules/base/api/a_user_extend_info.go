package api

import (
	"basic-frame/modules/base/bll"
	"basic-frame/modules/base/schema"
	"basic-frame/util/ginx"
	"github.com/gin-gonic/gin"
)

var UserExtendInfoApi = &UserExtendInfo{
	UserExtendInfoBll: bll.UserExtendInfoBll,
}

type UserExtendInfo struct {
	UserExtendInfoBll *bll.UserExtendInfo
}

// Update 更新用户扩展信息
func (a *UserExtendInfo) Update(c *gin.Context) {
	var item schema.UserExtendInfo
	if err := ginx.ParseJSON(c, &item); err != nil {
		ginx.ResError(c, item, err)
		return
	}

	if err := a.UserExtendInfoBll.Update(c, ginx.ParseParamID(c, "id"), item); err != nil {
		ginx.ResError(c, item, err)
		return
	}
	ginx.ResOperateSuccess(c, item)
}
