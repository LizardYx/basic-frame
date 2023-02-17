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
	"strings"
)

var ButtonBll = &Button{
	ButtonModel:     model.ButtonModel,
	RestfulApiModel: model.RestfulApiModel,
}

type Button struct {
	Enforcer        *casbin.SyncedEnforcer
	ButtonModel     *model.Button
	RestfulApiModel *model.RestfulApi
}

func (a *Button) Query(c *gin.Context, params schema.ButtonQueryParam) (*schema.ButtonQueryResult, error) {
	return a.ButtonModel.Query(params)
}

// Create 创建按钮和按钮关联的RestfulApi(不会创建子按钮)
func (a *Button) Create(c *gin.Context, item schema.ButtonPre) (*common.IDResult, error) {
	// 初始化按钮UUID
	a.InitUUID(&item, true)

	// 按钮参数验证
	if err := a.BtnPreParamsCheck(&item); err != nil {
		return nil, err
	}

	// 创建按钮
	return a.ButtonModel.CreateButtonPre(item)
}

func (a *Button) Get(c *gin.Context, id uint64) (*schema.Button, error) {
	// 获取指定按钮信息
	item, err := a.ButtonModel.Get(id)
	if err != nil {
		return nil, errors.WithMessage(err, "获取按钮信息失败")
	} else if item == nil {
		return nil, errors.New("未找到该按钮")
	}

	return item, nil
}

// GetButtonAndRestfulApis 获取按钮和按钮的RestfulApi信息
func (a *Button) GetButtonAndRestfulApis(c *gin.Context, id uint64) (*schema.ButtonPre, error) {
	// 获取指定按钮信息
	item, err := a.ButtonModel.GetButtonRestfulApis(id)
	if err != nil {
		return nil, errors.WithMessage(err, "获取按钮信息失败")
	} else if item == nil {
		return nil, errors.New("未找到该按钮")
	}

	return item, nil
}

// Update 更新按钮基础信息
func (a *Button) Update(c *gin.Context, id uint64, item schema.Button) error {
	// 参数检查
	if err := a.BtnParamsCheck(&item); err != nil {
		return err
	}

	// 检查按钮是否存在
	if _, err := a.Get(c, id); err != nil {
		return err
	}

	// 更新按钮信息
	return a.ButtonModel.UpdateByID(id, map[string]interface{}{
		"btn_id":      item.BtnID,
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

// BatchUpdate 批量更新按钮基础信息
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
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return err
}

// UpdateButtonRestfulApis 更新按钮关联的Restful接口信息
func (a *Button) UpdateButtonRestfulApis(c *gin.Context, id uint64, items schema.RestfulApis) error {
	// 检查按钮是否存在
	buttonPre, err := a.GetButtonAndRestfulApis(c, id)
	if err != nil {
		return err
	}

	return mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 新增Restful接口信息,并更新按钮和Restful接口的关联
		if err = a.ButtonModel.UpdateButtonRestfulApis(id, items); err != nil {
			return errors.WithMessage(err, "新增Restful接口信息失败")
		}

		// 更新Restful接口信息
		for _, item := range items {
			if item.ID != 0 {
				if err = a.RestfulApiModel.UpdateByID(item.ID, map[string]interface{}{
					"method": item.Method,
					"path":   item.Path,
					"memo":   item.Memo,
				}); err != nil {
					return errors.WithMessage(err, "更新Restful接口信息失败")
				}
			}
		}

		// 删除Restful接口信息
		var currentApiUUIDs = items.GetUUID()
		for _, restfulApi := range buttonPre.RestfulApis {
			if !common.ContainsString(currentApiUUIDs, restfulApi.UUID) {
				// 如果该Restful接口UUID不存在了，则删除
				if err = a.RestfulApiModel.Delete(restfulApi.ID); err != nil {
					return errors.WithMessage(err, "删除Restful接口信息失败")
				}
			}
		}
		return nil
	})
}

// UpdateButtonPre 更新按钮基础信息、按钮关联的Restful接口信息、指定Restful接口基础信息
func (a *Button) UpdateButtonPre(c *gin.Context, item schema.ButtonPre) error {
	// 初始化按钮信息
	a.InitUUID(&item, false)

	// 按钮参数验证
	if err := a.BtnPreParamsCheck(&item); err != nil {
		return err
	}

	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 更新按钮基础信息
		if err := a.Update(c, item.ID, *item.ToSchemaButton()); err != nil {
			return err
		}

		// 更新按钮关联的Restful接口信息
		if err := a.UpdateButtonRestfulApis(c, item.ID, item.RestfulApis); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return err
}

// Delete 删除按钮和按钮调用的Api
func (a *Button) Delete(c *gin.Context, id uint64) error {
	// 检查按钮是否有子项
	if ButtonQueryResult, err := a.ButtonModel.Query(schema.ButtonQueryParam{
		PaginationParam: common.PaginationParam{
			OnlyCount: true,
		},
		ParentID: id,
	}); err != nil {
		return errors.WithMessage(err, "检查按钮是否有子项失败")
	} else if ButtonQueryResult.PageResult.Total != 0 {
		return errors.New("有子按钮，请勿删除")
	}

	// 获取按钮信息和按钮的RestfulApi信息
	buttonPre, err := a.GetButtonAndRestfulApis(c, id)
	if err != nil {
		return err
	}

	// 删除按钮和按钮关联的Api信息
	err = mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 删除按钮的RestfulApi信息
		for _, restfulApi := range buttonPre.RestfulApis {
			if err = a.RestfulApiModel.Delete(restfulApi.ID); err != nil {
				return err
			}
		}

		// 删除按钮信息
		if err = a.ButtonModel.Delete(id); err != nil {
			return errors.WithMessage(err, "删除按钮失败")
		}
		return nil
	})
	if err != nil {
		return err
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return nil
}

// ---------------------------------------- Params  Validate --------------------------------------

// BtnParamsCheck 按钮参数校验
func (a *Button) BtnParamsCheck(item *schema.Button) error {
	// 检查按钮名称
	item.Name = strings.TrimSpace(item.Name)
	if item.Name == "" {
		return errors.New("按钮名称不能为空")
	}

	// 检查按钮的父按钮ID
	if item.ParentID != nil {
		parentID := uint64(0)
		if *item.ParentID == parentID {
			item.ParentID = nil
		}
	}
	return nil
}

// BtnPreParamsCheck 按钮参数校验
func (a *Button) BtnPreParamsCheck(item *schema.ButtonPre) error {
	item.Name = strings.TrimSpace(item.Name)
	if item.Name == "" {
		return errors.New("按钮名称不能为空")
	}
	if len(item.SonButtons) != 0 {
		for index, _ := range item.SonButtons {
			if err := a.BtnPreParamsCheck(item.SonButtons[index]); err != nil {
				return err
			}
		}
	}
	return nil
}

// InitUUID 初始化按钮、restfulApi 的UUID数据
func (a *Button) InitUUID(item *schema.ButtonPre, isCreate bool) {
	if item.UUID == "" {
		item.ID = 0
		item.UUID = common.GetUUID()
	} else if isCreate {
		item.ID = 0
	}
	for _, restfulApi := range item.RestfulApis {
		if restfulApi.UUID == "" {
			restfulApi.ID = 0
			restfulApi.UUID = common.GetUUID()
		} else if isCreate {
			restfulApi.ID = 0
		}
	}
	for _, sonButton := range item.SonButtons {
		a.InitUUID(sonButton, isCreate)
	}
}
