package bll

import (
	"basic-frame/modules/base/dao/model"
	"basic-frame/util/ginx/errors"
	"github.com/gin-gonic/gin"
)

var DisabledFieldBll = &DisabledField{
	DisabledFieldModel: model.DisabledFieldModel,
}

type DisabledField struct {
	DisabledFieldModel *model.DisabledField
}

func (a *DisabledField) Delete(c *gin.Context, id uint64) error {
	// 检查可禁用字段是否存在
	if item, err := a.DisabledFieldModel.Get(id); err != nil {
		return errors.WithMessage(err, "检查可禁用字段是否存在失败")
	} else if item == nil {
		return errors.New("未找到该可禁用字段")
	}

	// 删除可禁用字段
	if err := a.DisabledFieldModel.Delete(id); err != nil {
		return errors.WithMessage(err, "删除可禁用字段失败")
	}
	return nil
}
