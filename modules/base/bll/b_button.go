package bll

import (
	"basic-frame/modules/base/dao/model"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/ginx/errors"
	"basic-frame/util/mysql"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

var ButtonBll = &Button{
	ButtonModel: model.ButtonModel,
}

type Button struct {
	Enforcer        *casbin.SyncedEnforcer
	ButtonModel     *model.Button
	RestfulApiModel *model.RestfulApi
}

func (a *Button) Query(c *gin.Context, params schema.ButtonQueryParam) (*schema.ButtonQueryResult, error) {
	return a.ButtonModel.Query(params)
}

func (a *Button) Create(c *gin.Context, item schema.ButtonPre) (*common.IDResult, error) {
	// 初始化按钮UUID
	item.UUID = common.GetUUID()

	// 按钮参数验证
	if err := a.BtnParamsCheck(&item); err != nil {
		return nil, err
	}

	// 创建按钮
	return a.ButtonModel.CreateButtonPre(item)
}

func (a *Button) Update(c *gin.Context, id uint64, item schema.Button) error {
	// 检查按钮是否存在
	oldItem, err := a.ButtonModel.Get(id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.New("未找到该按钮")
	}

	// 参数检查
	item.Name = strings.TrimSpace(item.Name)
	if item.Name == "" {
		return errors.New("按钮名称不能为空")
	}
	if item.ParentID != nil {
		parentID := uint64(0)
		if *item.ParentID == parentID {
			item.ParentID = nil
		}
	}

	// 更新按钮信息
	return a.ButtonModel.UpdateByID(id, map[string]interface{}{
		"btn_id":      item.BtnID,
		"select":      item.Select,
		"name":        item.Name,
		"icon":        item.Icon,
		"class":       item.Class,
		"menu_id":     item.MenuID,
		"sequence":    item.Sequence,
		"parent_id":   item.ParentID,
		"show_status": item.ShowStatus,
		"status":      item.Status,
		"memo":        item.Memo,
	})
}

func (a *Button) BatchUpdate(c *gin.Context, items schema.Buttons) error {
	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := a.Update(c, item.ID, *item); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	LoadCasbinPolicy(c, a.Enforcer)
	return err
}

func (a *Button) UpdateButtonRestfulApis(c *gin.Context, id uint64, item schema.RestfulApis) error {
	oldItem, err := a.ButtonModel.Get(id)
	if err != nil {
		return err
	} else if oldItem == nil {
		return errors.New("未找到该按钮")
	}
	return a.ButtonModel.UpdateButtonRestfulApis(id, item)
}

func (a *Button) UpdateButtonPre(c *gin.Context, item schema.ButtonPre) error {
	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		if err := a.Update(c, item.ID, *item.ToSchemaButton()); err != nil {
			return err
		}
		if err := a.UpdateButtonRestfulApis(c, item.ID, item.RestfulApis); err != nil {
			return err
		}
		if len(item.RestfulApis) != 0 {
			for _, restfulApiInfo := range item.RestfulApis {
				if restfulApiInfo.ID != 0 {
					if err := a.RestfulApiModel.UpdateByID(restfulApiInfo.ID, map[string]interface{}{
						"method": restfulApiInfo.Method,
						"path":   restfulApiInfo.Path,
						"memo":   restfulApiInfo.Memo,
					}); err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	LoadCasbinPolicy(c, a.Enforcer)
	return err
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
