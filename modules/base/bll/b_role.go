package bll

import "C"
import (
	"basic-frame/modules/base/dao/model"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/consts"
	"basic-frame/util/ginx/errors"
	"github.com/gin-gonic/gin"
	"strings"
)

var RoleBll = &Role{
	RoleModel:   model.RoleModel,
	MenuModel:   model.MenuModel,
	ButtonModel: model.ButtonModel,
}

type Role struct {
	RoleModel   *model.Role
	MenuModel   *model.Menu
	ButtonModel *model.Button
}

func (a *Role) Query(c *gin.Context, params schema.RoleQueryParam) (*schema.RoleQueryResult, error) {
	return a.RoleModel.Query(params)
}

func (a *Role) Get(c *gin.Context, id uint64) (*schema.Role, error) {
	// 获取角色详情
	item, err := a.RoleModel.Get(id)
	if err != nil {
		return nil, errors.WithMessage(err, "获取角色信息失败")
	} else if item == nil {
		return nil, errors.New("角色不存在")
	}

	// 获取菜单树信息
	menuTrees, err := a.MenuModel.QueryMenuTreeForCreateRole()
	if err != nil {
		return nil, errors.WithMessage(err, "获取获取菜单树信息失败")
	}

	// 检查角色的菜单是否被选中
	for _, menuInfo := range item.Menus {
		// 获取菜单的所有子菜单ID集合
		var menuIDs []uint64
		menuTrees.GetMenuIDsByMenuID(menuInfo.ID, &menuIDs)
		// 获取菜单的所有按钮ID集合
		var buttonIDs []uint64
		menuTrees.GetButtonIDsByMenuID(menuInfo.ID, &buttonIDs)
		// 如果菜单的所有子菜单、所有按钮以及子按钮都被选中
		if a.IsSelectAll(menuIDs, item.Menus.GetIDs()) && a.IsSelectAll(buttonIDs, item.Buttons.GetIDs()) {
			menuInfo.Select = true
		}
	}

	// 检查角色的按钮是否被选中
	for _, buttonInfo := range item.Buttons {
		// 获取按钮的子按钮ID集合
		buttonIDs := menuTrees.GetSonButtonIDs(buttonInfo.MenuID, buttonInfo.ID)
		// 如果按钮的所有子按钮都被选中
		if a.IsSelectAll(buttonIDs, item.Buttons.GetIDs()) {
			buttonInfo.Select = true
		}
	}
	return item, nil
}

func (a *Role) Create(c *gin.Context, item schema.Role) (*common.IDResult, error) {
	// 检查前端传入的参数是否允许
	if err := a.RoleDetailInfoCheck(c, &item); err != nil {
		return &common.IDResult{}, err
	}

	// 创建角色
	IDResult, err := a.RoleModel.Create(item)
	if err != nil {
		return &common.IDResult{}, errors.WithMessage(err, "创建角色失败")
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return IDResult, nil
}

func (a *Role) Update(c *gin.Context, id uint64, item schema.Role) error {
	// 检查角色是否存在
	oldItem, err := a.RoleModel.Get(id)
	if err != nil {
		return errors.WithMessage(err, "检查角色是否存在失败")
	} else if oldItem == nil {
		return errors.New("角色不存在")
	}

	// 检查角色名称是否存在
	if RoleQueryResult, err := a.Query(c, schema.RoleQueryParam{
		Name: item.Name,
	}); err != nil {
		return errors.WithMessage(err, "检查角色名称是否存在失败")
	} else if len(RoleQueryResult.Data) != 0 {
		if RoleQueryResult.Data[0].ID != id {
			return errors.New("角色名称已存在")
		}
	}

	// 更新角色基本信息
	return a.RoleModel.UpdateByID(id, map[string]interface{}{
		"name":     item.Name,
		"sequence": item.Sequence,
		"status":   item.Status,
		"memo":     item.Memo,
	})
}

func (a *Role) UpdateRoleMenu(c *gin.Context, id uint64, items schema.Menus) error {
	// 检查角色是否存在
	oldItem, err := a.RoleModel.Get(id)
	if err != nil {
		return errors.WithMessage(err, "检查角色是否存在失败")
	} else if oldItem == nil {
		return errors.New("角色不存在")
	}

	return a.RoleModel.UpdateRoleMenu(id, items)
}

// ---------------------------------------- Params  Validate --------------------------------------

func (a *Role) RoleDetailInfoCheck(c *gin.Context, item *schema.Role) error {
	// 检查角色名称是否为空
	item.Name = strings.TrimSpace(item.Name)
	if item.Name == "" {
		return errors.New("角色名称不能空")
	}

	// 检查角色名称是否存在
	RoleQueryResult, err := a.Query(c, schema.RoleQueryParam{
		Name: item.Name,
	})
	if err != nil {
		return errors.WithMessage(err, "检查角色名称是否存在失败")
	} else if len(RoleQueryResult.Data) != 0 {
		if RoleQueryResult.Data[0].ID != item.ID {
			return errors.New("角色名称已存在")
		}
	}

	// 检查前端传的菜单是否有被禁用的
	if len(item.Menus) != 0 {
		if MenuQueryResult, err := a.MenuModel.Query(schema.MenuQueryParam{
			PaginationParam: common.PaginationParam{
				OnlyCount: true,
			},
			IDs:    item.Menus.GetIDs(),
			Status: consts.BaseStatusDisabled,
		}); err != nil {
			return errors.WithMessage(err, "检查菜单是否被禁用失败")
		} else if MenuQueryResult.PageResult.Total != 0 {
			return errors.New("被禁用的菜单不能创建角色")
		}
	}

	// 检查前端传的按钮是否有被禁用的
	if len(item.Buttons) != 0 {
		if ButtonQueryResult, err := a.ButtonModel.Query(schema.ButtonQueryParam{
			PaginationParam: common.PaginationParam{
				OnlyCount: true,
			},
			Status: consts.BaseStatusDisabled,
			IDs:    common.UintSliceToString(item.Buttons.GetIDs(), ","),
		}); err != nil {
			return errors.WithMessage(err, "检查按钮是否被禁用失败")
		} else if ButtonQueryResult.PageResult.Total != 0 {
			return errors.New("被禁用的按钮不能创建角色")
		}
	}

	// 检查角色类型是否可用
	if !common.ContainsUint64(consts.RoleTypes, uint64(item.Type)) {
		return errors.New("请检查角色类型是否正确")
	}

	return nil
}

func (a *Role) IsSelectAll(allIDs, currentIDs []uint64) bool {
	var isSelect bool
	if len(allIDs) == 0 {
		isSelect = true
	} else {
		for index, id := range allIDs {
			for i, currentID := range currentIDs {
				if currentID == id {
					break
				}
				if i == (len(currentIDs) - 1) {
					return isSelect
				}
			}
			if index == (len(allIDs) - 1) {
				isSelect = true
			}
		}
	}
	return isSelect
}
