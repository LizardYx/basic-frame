package bll

import (
	"basic-frame/modules/base/dao/model"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/ginx/errors"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

var ButtonBll = &Button{
	ButtonModel: model.ButtonModel,
}

type Button struct {
	Enforcer    *casbin.SyncedEnforcer
	ButtonModel *model.Button
}

func (a *Button) Query(c *gin.Context, params schema.ButtonQueryParam) (*schema.ButtonQueryResult, error) {
	return a.ButtonModel.Query(params)
}

// Delete 删除按钮和按钮调用的Api
func (a *Button) Delete(c *gin.Context, id uint64) error {
	// 检查按钮是否存在
	oldItem, err := a.ButtonModel.Get(id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.NewIANAResponse(http.StatusNotFound, "未找到该按钮")
	}

	// 检查按钮是否有子项
	if ButtonQueryResult, err := a.ButtonModel.Query(schema.ButtonQueryParam{
		PaginationParam: common.PaginationParam{
			OnlyCount: true,
		},
		ParentID: id,
	}); err != nil {
		return errors.NewIANAResponse(http.StatusInternalServerError, "删除菜单失败")
	} else if ButtonQueryResult.PageResult.Total != 0 {
		return errors.NewIANAResponse(http.StatusInternalServerError, "有子按钮，请勿删除")
	}
	LoadCasbinPolicy(c, a.Enforcer)
	return a.ButtonModel.Delete(id)
}

// ---------------------------------------- Params  Validate --------------------------------------

// BtnParamsCheck 按钮参数校验
func (a *Button) BtnParamsCheck(item *schema.ButtonPre) error {
	item.Name = strings.TrimSpace(item.Name)
	if item.Name == "" {
		return errors.New("按钮名称不能为空")
	}
	if len(item.SonButtons) != 0 {
		for index, _ := range item.SonButtons {
			if err := a.BtnParamsCheck(item.SonButtons[index]); err != nil {
				return err
			}
		}
	}
	return nil
}
