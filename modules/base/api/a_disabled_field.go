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

// Delete 删除禁用字段
//
//	@Summary		删除禁用字段
//	@Description	删除禁用字段
//	@Tags			DisabledField
//	@Param			id	path		int	true	"字段ID"
//	@Success		200	{object}	ginx.OperateResult
//	@Failure		500	{object}	ginx.OperateResult
//	@Router			/api/v1/base/disabled_field/:id [delete]
func (a *DisabledField) Delete(c *gin.Context) {
	if err := a.DisabledFieldBll.Delete(c, ginx.ParseParamID(c, "id")); err != nil {
		ginx.ResError(c, "", err)
		return
	}
	ginx.ResOperateSuccess(c, "")
}
