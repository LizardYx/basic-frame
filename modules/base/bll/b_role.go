package bll

import "C"
import (
	"basic-frame/modules/base/dao/model"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/consts"
	"basic-frame/util/ginx/errors"
	"basic-frame/util/mysql"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
)

var RoleBll = &Role{
	UserModel:          model.UserModel,
	UserGroupModel:     model.UserGroupModel,
	PositionModel:      model.PositionModel,
	OrganizationModel:  model.OrganizationModel,
	RoleModel:          model.RoleModel,
	MenuModel:          model.MenuModel,
	ButtonModel:        model.ButtonModel,
	SecurityLevelModel: model.SecurityLevelModel,
}

type Role struct {
	UserModel          *model.User
	UserGroupModel     *model.UserGroup
	PositionModel      *model.Position
	OrganizationModel  *model.Organization
	RoleModel          *model.Role
	MenuModel          *model.Menu
	ButtonModel        *model.Button
	SecurityLevelModel *model.SecurityLevel
}

func (a *Role) Query(c *gin.Context, params schema.RoleQueryParam) (*schema.RoleQueryResult, error) {
	return a.RoleModel.Query(params)
}

// Get 获取角色基本信息
func (a *Role) Get(c *gin.Context, id uint64) (*schema.Role, error) {
	// 获取角色详情
	item, err := a.RoleModel.Get(id)
	if err != nil {
		return nil, errors.WithMessage(err, "获取角色信息失败")
	} else if item == nil {
		return nil, errors.New("角色不存在")
	}

	return item, nil
}

// GetPre 获取角色所有信息(包含菜单、按钮、禁用字段)
func (a *Role) GetPre(c *gin.Context, id uint64) (*schema.RolePre, error) {
	// 获取角色详情
	item, err := a.RoleModel.GetPre(id)
	if err != nil {
		return nil, errors.WithMessage(err, "获取角色信息失败")
	} else if item == nil {
		return nil, errors.New("角色不存在")
	}

	return item, nil
}

// GetPreWithSelect 获取角色所有信息(包含菜单、按钮、禁用字段。以及菜单、按钮是否为全选)
func (a *Role) GetPreWithSelect(c *gin.Context, id uint64) (*schema.RolePre, error) {
	// 获取角色所有信息
	item, err := a.GetPre(c, id)
	if err != nil {
		return nil, err
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

func (a *Role) Create(c *gin.Context, item schema.RolePre) (*common.IDResult, error) {
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

// Update 更新角色基本信息
func (a *Role) Update(c *gin.Context, id uint64, item schema.Role) error {
	// 检查角色是否存在
	if _, err := a.Get(c, id); err != nil {
		return err
	}

	// 检查角色名称是否存在
	if RoleQueryResult, err := a.Query(c, schema.RoleQueryParam{
		Name: item.Name,
	}); err != nil {
		return errors.WithMessage(err, "检查角色名称是否存在失败")
	} else if len(RoleQueryResult.Data) != 0 && RoleQueryResult.Data[0].ID != id {
		return errors.New("角色名称已存在")
	}

	// 更新角色基本信息
	return a.RoleModel.UpdateByID(id, map[string]interface{}{
		"name":     item.Name,
		"sequence": item.Sequence,
		"status":   item.Status,
		"memo":     item.Memo,
	})
}

// UpdateRoleMenu 更新角色关联的菜单
func (a *Role) UpdateRoleMenu(c *gin.Context, id uint64, items schema.Menus) error {
	// 检查角色是否存在
	if _, err := a.Get(c, id); err != nil {
		return err
	}

	// 更新角色关联的菜单
	return a.RoleModel.UpdateRoleMenu(id, items)
}

// UpdateRoleButton 更新角色关联的按钮
func (a *Role) UpdateRoleButton(c *gin.Context, id uint64, items schema.Buttons) error {
	// 检查角色是否存在
	if _, err := a.Get(c, id); err != nil {
		return err
	}

	// 更新角色关联的按钮
	return a.RoleModel.UpdateRoleButton(id, items)
}

// UpdateRoleDisabledFields 更新角色关联的禁用字段
func (a *Role) UpdateRoleDisabledFields(c *gin.Context, id uint64, items schema.DisabledFields) error {
	// 检查角色是否存在
	if _, err := a.Get(c, id); err != nil {
		return err
	}

	// 更新角色关联的禁用字段
	return a.RoleModel.UpdateRoleDisabledFields(id, items)
}

// UpdateDetails 更新角色基本信息、角色和菜单的关联、角色和按钮的关联、角色和可禁用字段的关联
func (a *Role) UpdateDetails(c *gin.Context, item schema.RolePre) error {
	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 检查前端传入的参数是否允许
		if err := a.RoleDetailInfoCheck(c, &item); err != nil {
			return err
		}

		// 更新角色基本信息
		if err := a.Update(c, item.ID, *item.ToRole()); err != nil {
			return err
		}

		// 更新角色和菜单的关联
		if err := a.UpdateRoleMenu(c, item.ID, item.Menus); err != nil {
			return err
		}

		// 更新角色和按钮的关联
		if err := a.UpdateRoleButton(c, item.ID, item.Buttons); err != nil {
			return err
		}

		// 更新角色和可禁用字段的关联
		if err := a.UpdateRoleDisabledFields(c, item.ID, item.DisabledFields); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return nil
}

// UpdateAuditorType 给角色添加审核类型
func (a *Role) UpdateAuditorType(c *gin.Context, id uint64, item schema.UpdateAuditorTypeParam) error {
	return mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 检查该审核类型是否已被其它角色使用
		if roleQueryResult, err := a.RoleModel.Query(schema.RoleQueryParam{
			AuditorTypes: common.IntSliceToString([]int{item.AuditorType}, ","),
			FindAll:      true,
		}); err != nil {
			return errors.WithMessage(err, "检查审核类型失败")
		} else if len(roleQueryResult.Data) != 0 {
			// 审核类型已被其它角色使用，将该审核类型从其它角色中移除
			for _, roleInfo := range roleQueryResult.Data {
				var newAuditorTypes []uint64
				oldAuditorTypes := common.SplitStringToUint64(roleInfo.AuditorTypes, ",")
				for _, auditorType := range oldAuditorTypes {
					if auditorType != uint64(item.AuditorType) {
						newAuditorTypes = append(newAuditorTypes, auditorType)
					}
				}
				if err = a.RoleModel.UpdateByID(roleInfo.ID, map[string]interface{}{
					"auditor_types": common.UintSliceToString(newAuditorTypes, ","),
				}); err != nil {
					return errors.WithMessage(err, "更新角色审核类型失败")
				}
			}
		}

		// 获取角色基本信息
		roleInfo, err := a.Get(c, id)
		if err != nil {
			return err
		}

		// 设置审核类型给角色
		newAuditorTypes := common.SplitStringToUint64(roleInfo.AuditorTypes, ",")
		if !common.ContainsUint64(newAuditorTypes, uint64(item.AuditorType)) {
			newAuditorTypes = append(newAuditorTypes, uint64(item.AuditorType))
		}
		if err = a.RoleModel.UpdateByID(id, map[string]interface{}{
			"auditor_types": common.UintSliceToString(newAuditorTypes, ","),
		}); err != nil {
			return errors.WithMessage(err, "更新角色审核类型失败")
		}
		return nil
	})
}

// UserAddRole 给用户添加角色
func (a *Role) UserAddRole(c *gin.Context, roleID uint64, userIDs []uint64) error {
	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 获取角色基本信息
		item, err := a.Get(c, roleID)
		if err != nil {
			return err
		}

		// 给用户添加角色
		for _, userID := range userIDs {
			// 用户添加角色
			if err = a.UserModel.AppendUserRoles(userID, schema.Roles{item}); err != nil {
				return errors.WithMessage(err, "用户添加角色失败")
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return nil
}

// UserRemoveRole 移除用户的指定角色
func (a *Role) UserRemoveRole(c *gin.Context, roleID uint64, userIDs []uint64) error {
	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 获取角色基本信息
		roleInfo, err := a.Get(c, roleID)
		if err != nil {
			return err
		}

		// 移除用户的指定角色
		for _, userID := range userIDs {
			if err = a.UserModel.UserRemoveRole(userID, *roleInfo); err != nil {
				return errors.WithMessage(err, "用户角色移除失败")
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return nil
}

// Delete 删除角色
func (a *Role) Delete(c *gin.Context, id uint64) error {
	// 获取角色基本信息
	if oldItem, err := a.Get(c, id); err != nil {
		return err
	} else if oldItem.AuditorTypes != "" {
		// 检查是否是审核角色
		return errors.New("审核角色不能删除")
	}

	// 检查是否有组织绑定了该角色
	if OrgQueryResult, err := a.OrganizationModel.Query(schema.OrganizationQueryParam{
		RoleID: id,
	}); err != nil {
		return errors.WithMessage(err, "检查角色是否被组织使用失败")
	} else if len(OrgQueryResult.Data) != 0 {
		var orgNames []string
		for _, item := range OrgQueryResult.Data {
			orgNames = append(orgNames, item.Name)
		}
		return errors.New(fmt.Sprintf("该角色已被组织: %s 使用", common.StringSliceToString(orgNames, ",")))
	}

	// 检查是否有职位绑定了该角色
	if PositionQueryResult, err := a.PositionModel.Query(schema.PositionQueryParam{
		RoleID: id,
	}); err != nil {
		return errors.WithMessage(err, "检查角色是否被职位使用失败")
	} else if len(PositionQueryResult.Data) != 0 {
		var positionNames []string
		for _, position := range PositionQueryResult.Data {
			positionNames = append(positionNames, position.Name)
		}
		return errors.New(fmt.Sprintf("该角色已被职位: %s 使用", common.StringSliceToString(positionNames, ",")))
	}

	// 检查是否有用户组绑定了该角色
	if UserGroupQueryResult, err := a.UserGroupModel.Query(schema.UserGroupQueryParam{
		RoleID: id,
	}); err != nil {
		return errors.WithMessage(err, "检查角色是否被用户组使用失败")
	} else if len(UserGroupQueryResult.Data) != 0 {
		var userGroupNames []string
		for _, userGroup := range UserGroupQueryResult.Data {
			userGroupNames = append(userGroupNames, userGroup.Name)
		}
		return errors.New(fmt.Sprintf("该角色已被用户组: %s 使用", common.StringSliceToString(userGroupNames, ",")))
	}

	// 检查是否有用户绑定了该角色
	if UserQueryResult, err := a.UserModel.Query(schema.UserQueryParam{
		RoleID:       id,
		OmitPassword: true,
	}); err != nil {
		return errors.WithMessage(err, "检查角色是否被用户使用失败")
	} else if len(UserQueryResult.Data) != 0 {
		var userNames []string
		for _, user := range UserQueryResult.Data {
			userNames = append(userNames, user.UserName)
		}
		return errors.New(fmt.Sprintf("该角色已被用户: %s 使用", common.StringSliceToString(userNames, ",")))
	}

	// 检查是否有安全级别绑定了该角色
	if SecurityLevelQueryResult, err := a.SecurityLevelModel.Query(schema.SecurityLevelQueryParam{
		RoleIDs: common.UintSliceToString([]uint64{id}, ","),
		FindAll: true,
	}); err != nil {
		return errors.WithMessage(err, "检查角色是否被安全级别使用失败")
	} else if len(SecurityLevelQueryResult.Data) != 0 {
		var securityLevelNames []string
		for _, securityLevel := range SecurityLevelQueryResult.Data {
			securityLevelNames = append(securityLevelNames, securityLevel.Name)
		}
		return errors.New(fmt.Sprintf("该角色已被安全级别: %s 使用", common.StringSliceToString(securityLevelNames, ",")))
	}

	// 删除角色
	if err := a.RoleModel.Delete(id); err != nil {
		return errors.WithMessage(err, "删除角色失败")
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return nil
}

// ---------------------------------------- Params  Validate --------------------------------------

func (a *Role) RoleDetailInfoCheck(c *gin.Context, item *schema.RolePre) error {
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
			IDs:    common.UintSliceToString(item.Menus.GetIDs(), ","),
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
