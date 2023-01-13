package api

import (
	"basic-frame/modules/base/bll"
	"basic-frame/util/ginx"
	"github.com/gin-gonic/gin"
)

var DisabledFieldApi = &DisabledField{
	DisabledFieldBll: bll.DisabledFieldBll,
}

type DisabledField struct {
	DisabledFieldBll *bll.DisabledField
}

func (a *DisabledField) Delete(c *gin.Context) {
	if err := a.DisabledFieldBll.Delete(c, ginx.ParseParamID(c, "id")); err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}
