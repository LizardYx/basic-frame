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
	params := ginx.ParamsID{
		ID: ginx.ParseParamID(c, "id"),
	}
	err := a.DisabledFieldBll.Delete(c, params.ID)
	if err != nil {
		ginx.ResError(c, params, err)
		return
	}
	ginx.ResOperateSuccess(c, params)
}
