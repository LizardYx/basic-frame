package bll

import (
	"basic-frame/middleware"
	"basic-frame/modules/base/dao/model"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/consts"
	"basic-frame/util/ginx"
	"basic-frame/util/ginx/errors"
	"basic-frame/util/mysql"
	"basic-frame/util/secret"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

var UserBll = &User{
	ButtonModel:         model.ButtonModel,
	MenuModel:           model.MenuModel,
	RoleModel:           model.RoleModel,
	OrganizationModel:   model.OrganizationModel,
	PositionModel:       model.PositionModel,
	UserModel:           model.UserModel,
	UserGroupModel:      model.UserGroupModel,
	UserExtendInfoModel: model.UserExtendInfoModel,
}

type User struct {
	ButtonModel         *model.Button
	MenuModel           *model.Menu
	RoleModel           *model.Role
	OrganizationModel   *model.Organization
	PositionModel       *model.Position
	UserModel           *model.User
	UserGroupModel      *model.UserGroup
	UserExtendInfoModel *model.UserExtendInfo
}

func (a *User) Login(c *gin.Context, item schema.LoginParam) (*schema.LoginTokenInfo, *schema.User, error) {
	// 验证用户名密码
	user, err := a.Verify(c, item.UserName, item.Password)
	if err != nil {
		return nil, nil, err
	}
	user.Password = ""

	// 生成令牌
	newTokenInfo := schema.LoginTokenInfo{}
	if token, err := middleware.GenerateToken(user.ID, user.UserName); err != nil {
		return nil, nil, errors.WithMessage(err, "生成令牌失败")
	} else {
		expiredTime, _ := strconv.Atoi(consts.JWTAuthExpired)
		newTokenInfo = schema.LoginTokenInfo{
			AccessToken: token,
			TokenType:   consts.JWTAuthType,
			ExpiresAt:   time.Now().Add(time.Duration(expiredTime) * time.Second).Unix(),
		}
	}
	return &newTokenInfo, user, nil
}

func (a *User) Query(c *gin.Context, params schema.UserQueryParam) (*schema.UserQueryResult, error) {
	// 如果要查询某种审核类型的用户，先获取审核角色ID
	if params.AuditorTypes != "" {
		if roleQueryResult, err := a.RoleModel.Query(schema.RoleQueryParam{
			AuditorTypes: params.AuditorTypes,
			Status:       consts.BaseStatusEnable,
			FindAll:      true,
		}); err != nil {
			return nil, errors.WithMessage(err, "获取用户列表失败")
		} else if len(roleQueryResult.Data) != 0 {
			var roleIDs []uint64
			if params.RoleIDs != "" {
				roleIDs = common.SplitStringToUint64(params.RoleIDs, ",")
			}
			roleIDs = append(roleIDs, roleQueryResult.Data.GetIDs()...)
			params.RoleIDs = common.UintSliceToString(roleIDs, ",")
		}
	}

	return a.UserModel.Query(params)
}

func (a *User) GetSubUsers(c *gin.Context, userID uint64, params schema.SubUserQueryParam) (*schema.UserQueryResult, error) {
	users := make(schema.Users, 0)

	// 获取当前用户的信息
	userInfo, err := a.Get(c, userID)
	if err != nil {
		return nil, errors.WithMessage(err, "获取当前用户信息失败")
	}

	// 获取 拥有用户同组织子职位的用户
	if params.CurrentOrg && len(userInfo.Positions) != 0 {
		var sonPositionIDs []uint64
		// 获取所有的职位
		PositionQueryResult, err := a.PositionModel.Query(schema.PositionQueryParam{
			Status:  consts.BaseStatusEnable,
			FindAll: true,
		})
		if err != nil {
			return nil, errors.WithMessage(err, "获取职位列表失败")
		}

		// 获取用户同组织的子职位ID
		for _, item := range userInfo.Positions {
			for _, position := range PositionQueryResult.Data {
				// 如果职位属于同一组织且权重小于用户的职位
				if item.OrganizationID == position.OrganizationID && position.Sequence < item.Sequence &&
					!common.ContainsUint64(sonPositionIDs, position.ID) {
					sonPositionIDs = append(sonPositionIDs, position.ID)
				}
			}
		}

		// 获取拥有该子职位的用户
		if len(sonPositionIDs) != 0 {
			if UserQueryResult, err := a.UserModel.Query(schema.UserQueryParam{
				Status:       consts.BaseStatusEnable,
				PositionIDs:  common.UintSliceToString(sonPositionIDs, ","),
				FindAll:      true,
				OmitPassword: true,
			}); err != nil {
				return nil, errors.WithMessage(err, "获取用户信息失败")
			} else {
				for _, user := range UserQueryResult.Data {
					if len(users) == 0 {
						users = append(users, user)
						continue
					}
					if !common.ContainsUint64(users.GetIDs(), user.ID) {
						users = append(users, user)
					}
				}
			}
		}
	}

	return &schema.UserQueryResult{
		PageResult: nil,
		Data:       *users.SortUser(),
	}, nil
}

// Init 创建超管用户
func (a *User) Init(c *gin.Context) error {
	// 检查超管用户是否创建
	if UserQueryResult, err := a.Query(c, schema.UserQueryParam{
		UserName:     consts.AdminName,
		OmitPassword: true,
	}); err != nil {
		return errors.WithMessage(err, "检查超管用户是否创建失败")
	} else if len(UserQueryResult.Data) != 0 {
		// 超管已创建
		return nil
	} else {
		// 超管未创建
		if _, err = a.Create(c, schema.User{
			UserName: consts.AdminName,
			Password: consts.AdminPassword,
			Status:   consts.BaseStatusEnable,
		}); err != nil {
			return errors.WithMessage(err, "创建超管用户失败")
		}
	}
	return nil
}

func (a *User) Get(c *gin.Context, id uint64) (*schema.User, error) {
	// 获取用户信息
	item, err := a.UserModel.Get(id)
	if err != nil {
		return nil, errors.WithMessage(err, "获取用户信息失败")
	} else if item == nil {
		return nil, errors.New("用户不存在")
	}

	item.Password = ""
	return item, nil
}

func (a *User) Create(c *gin.Context, item schema.User) (*common.IDResult, error) {
	// 检查组织、职位、角色、用户组是否被禁用
	if err := a.UserParamsCheck(c, &item); err != nil {
		return nil, err
	}

	item.Password = secret.SHA1String(item.Password)
	var IDResult *common.IDResult
	var err error
	_ = mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 创建用户
		if IDResult, err = a.UserModel.Create(item); err != nil {
			return err
		}

		// 创建用户扩展信息
		if _, err = a.UserExtendInfoModel.Create(schema.UserExtendInfo{UserID: IDResult.ID, RealName: item.UserName}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return IDResult, nil
}

// Update 更新用户基础信息
func (a *User) Update(c *gin.Context, id uint64, item schema.User) error {
	// 检查用户是否存在
	if _, err := a.Get(c, id); err != nil {
		return err
	}

	// 更新用户信息
	if err := a.UserModel.UpdateByID(id, map[string]interface{}{
		"status":   item.Status,
		"sequence": item.Sequence,
	}); err != nil {
		return err
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return nil
}

func (a *User) UpdateStatus(c *gin.Context, id uint64, newStatus int) error {
	// 检查用户是否存在
	if oldItem, err := a.Get(c, id); err != nil {
		return err
	} else if oldItem.Status == newStatus {
		if newStatus == consts.BaseStatusEnable {
			return errors.New("用户已启用，请勿重复操作")
		} else {
			return errors.New("用户已禁用，请勿重复操作")
		}
	}

	// 更新用户状态
	return a.UserModel.UpdateByID(id, map[string]interface{}{
		"status": newStatus,
	})
}

func (a *User) UpdatePassword(c *gin.Context, id uint64, item schema.UpdatePasswordParam) error {
	item.NewPassword = strings.TrimSpace(item.NewPassword)

	// 检查用户是否存在
	if oldItem, err := a.Get(c, id); err != nil {
		return err
	} else if secret.SHA1String(item.OldPassword) != oldItem.Password {
		return errors.New("旧密码不正确")
	} else if item.NewPassword == "" {
		return errors.New("新密码不能为空")
	}

	// 更新用户信息
	item.NewPassword = secret.SHA1String(item.NewPassword)
	return a.UserModel.UpdateByID(id, map[string]interface{}{
		"password": item.NewPassword,
	})
}

func (a *User) Delete(c *gin.Context, id uint64) error {
	// 检查用户是否存在
	if oldItem, err := a.Get(c, id); err != nil {
		return err
	} else if oldItem.UserName == consts.AdminName {
		// 检查是否是超管
		return errors.New("超管用户不允许删除")
	}

	// 删除用户、用户扩展信息
	if err := a.UserModel.Delete(id); err != nil {
		return err
	}
	LoadCasbinPolicy(c, common.SysConfig.CasbinSyncEnforcer)
	return nil
}

// ---------------------------------------- User Organization --------------------------------------

// GetMenuTreesInfo 获取菜单树顶层结构
func (a *User) GetMenuTreesInfo(MenuTrees *schema.MenuTrees, parentId *uint64, menus schema.Menus, buttons schema.Buttons) error {
	for _, menu := range menus {
		if menu.ShowStatus == consts.BaseShowStatusDisabled {
			continue
		}
		if *menu.ParentID == *parentId {
			newMenuTree := schema.MenuTree{
				Menu:        *menu,
				RestfulApis: make(schema.RestfulApis, 0),
				Buttons:     make(schema.ButtonPres, 0),
				SonMenus:    make(schema.MenuTrees, 0),
			}
			// 获取该菜单的按钮树
			parentId := uint64(0)
			if err := a.GetButtonTreesInfo(&newMenuTree.Buttons, newMenuTree.ID, &parentId, buttons); err != nil {
				return err
			}
			// 获取每个顶层菜单的子菜单树以及菜单包含的按钮树
			if err := a.GetMenuTreesInfo(&newMenuTree.SonMenus, &newMenuTree.ID, menus, buttons); err != nil {
				return err
			}
			*MenuTrees = append(*MenuTrees, &newMenuTree)
		}
	}
	return nil
}

func (a *User) GetButtonTreesInfo(buttonTrees *schema.ButtonPres, menuId uint64, parentId *uint64, buttons schema.Buttons) error {
	for _, button := range buttons {
		if button.ShowStatus == consts.BaseShowStatusDisabled {
			continue
		}
		if *button.ParentID == *parentId && button.MenuID == menuId {
			newButtonTree := schema.ButtonPre{
				Button:      *button,
				RestfulApis: make(schema.RestfulApis, 0),
				SonButtons:  make(schema.ButtonPres, 0),
			}
			if err := a.GetButtonTreesInfo(&newButtonTree.SonButtons, menuId, &newButtonTree.ID, buttons); err != nil {
				return err
			}
			*buttonTrees = append(*buttonTrees, &newButtonTree)
		}
	}
	return nil
}

// GetMenuTree 获取用户菜单树
func (a *User) GetMenuTree(c *gin.Context) (*schema.MenuTrees, error) {
	var menus schema.Menus
	var buttons schema.Buttons
	menuTrees := make(schema.MenuTrees, 0)
	if ginx.GetUserName(c) == consts.AdminName {
		// 如果是超管用户,获取用户所有的菜单和按钮
		if menuQueryResult, err := a.MenuModel.Query(schema.MenuQueryParam{
			Status:     consts.BaseStatusEnable,
			ShowStatus: consts.BaseShowStatusEnable,
		}); err != nil {
			return &menuTrees, err
		} else {
			menus = menuQueryResult.Data
		}
		if buttonQueryResult, err := a.ButtonModel.Query(schema.ButtonQueryParam{
			ShowStatus: consts.BaseShowStatusEnable,
			Status:     consts.BaseStatusEnable,
		}); err != nil {
			return &menuTrees, err
		} else {
			buttons = buttonQueryResult.Data
		}
	} else {
		// 如果是普通用户，获取用户所有的角色ID
		roleIds, err := a.GetUserAllRoleIds(c, ginx.GetUserID(c))
		if err != nil {
			return &menuTrees, err
		} else if len(*roleIds) == 0 {
			return &menuTrees, err
		}
		// 根据角色获取菜单列表、按钮列表
		if RoleQueryResult, err := a.RoleModel.Query(schema.RoleQueryParam{
			IDs:         common.UintSliceToString(*roleIds, ","),
			ShowDetails: true,
		}); err != nil {
			return &menuTrees, err
		} else {
			for _, roleInfo := range RoleQueryResult.Data {
				// 获取用户菜单列表
				for _, newMenuInfo := range roleInfo.Menus {
					if len(menus) == 0 {
						menus = append(menus, newMenuInfo)
					} else {
						for index, menuInfo := range menus {
							if menuInfo.ID == newMenuInfo.ID {
								break
							}
							if index == (len(menus) - 1) {
								menus = append(menus, newMenuInfo)
							}
						}
					}
				}
				// 获取用户按钮列表
				for _, newButtonInfo := range roleInfo.Buttons {
					if len(buttons) == 0 {
						buttons = append(buttons, newButtonInfo)
					} else {
						for index, buttonInfo := range buttons {
							if buttonInfo.ID == newButtonInfo.ID {
								break
							}
							if index == (len(buttons) - 1) {
								buttons = append(buttons, newButtonInfo)
							}
						}
					}
				}
			}
		}
	}
	// 解析菜单列表和按钮列表，组成菜单树
	parentId := uint64(0)
	if err := a.GetMenuTreesInfo(&menuTrees, &parentId, menus, buttons); err != nil {
		return &menuTrees, err
	}
	return menuTrees.Init().SortMenuTrees(), nil
}

func (a *User) GetUserAllRoleIds(c *gin.Context, userID uint64) (*[]uint64, error) {
	// 获取用户信息，包含组织、职位、角色、用户组
	item, err := a.Get(c, userID)
	if err != nil {
		return nil, err
	}

	// 获取用户的角色集合
	var roleIds []uint64
	if len(item.Organizations) != 0 {
		// 获取组织结构的角色ID，如果用户有该组织结构的权限，那么就自动会拥有子组织的权限
		if OrgQueryResult, err := a.OrganizationModel.Query(schema.OrganizationQueryParam{
			IDs:         common.UintSliceToString(item.Organizations.GetIDs(), ","),
			Status:      consts.BaseStatusEnable,
			ShowDetails: true,
		}); err != nil {
			return nil, err
		} else {
			OrgQueryResult.Data.GetRoleIds(&roleIds)
		}
	}
	item.Positions.GetRoleIds(&roleIds)
	item.UserGroups.GetRoleIds(&roleIds)
	roleIds = append(roleIds, item.Roles.GetIDs()...)

	// 角色ID去重
	newRoleIds := common.RemoveRepeatedElement(roleIds)
	return &newRoleIds, nil
}

// UpdateUserOrganizations 更新用户关联的组织信息
func (a *User) UpdateUserOrganizations(c *gin.Context, id uint64, items schema.Organizations) error {
	return a.UserModel.ReplaceUserOrganizations(id, items)
}

// UpdateUserPositions 更新用户关联的职位信息
func (a *User) UpdateUserPositions(c *gin.Context, id uint64, items schema.Positions) error {
	return a.UserModel.ReplaceUserPositions(id, items)
}

// UpdateUserRoles 更新用户关联的角色信息
func (a *User) UpdateUserRoles(c *gin.Context, id uint64, items schema.Roles) error {
	return a.UserModel.ReplaceUserRoles(id, items)
}

// UpdateUserUserGroup 更新用户关联的用户组信息
func (a *User) UpdateUserUserGroup(c *gin.Context, id uint64, items schema.UserGroups) error {
	return a.UserModel.ReplaceUserUserGroup(id, items)
}

// ---------------------------------------- User Permission --------------------------------------

// UpdateUserPermission 更新用户权限
func (a *User) UpdateUserPermission(c *gin.Context, id uint64, item schema.User) error {
	err := mysql.DB.Transaction(func(tx *gorm.DB) error {
		// 检查用户是否存在
		oldItem, err := a.Get(c, id)
		if err != nil {
			return err
		} else {
			item.ID = oldItem.ID
			item.Password = oldItem.Password
		}

		// 检查组织、职位、角色、用户组是否被禁用
		if err = a.UserParamsCheck(c, &item); err != nil {
			return err
		}

		// 更新用户权重
		if err = a.Update(c, id, item); err != nil {
			return err
		}

		// 更新用户和组织的关联关系
		if err = a.UpdateUserOrganizations(c, id, item.Organizations); err != nil {
			return err
		}

		// 更新用户和职位的关联关系
		if err = a.UpdateUserPositions(c, id, item.Positions); err != nil {
			return err
		}

		// 更新用户和角色的关联关系
		if err = a.UpdateUserRoles(c, id, item.Roles); err != nil {
			return err
		}

		// 更新用户和用户组的关联关系
		if err = a.UpdateUserUserGroup(c, id, item.UserGroups); err != nil {
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

// BatchUpdateUserPermission 批量更新用户权限
func (a *User) BatchUpdateUserPermission(c *gin.Context, items schema.Users) error {
	return mysql.DB.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := a.UpdateUserPermission(c, item.ID, *item); err != nil {
				return err
			}
		}
		return nil
	})
}

// ---------------------------------------- Params Validate --------------------------------------

func (a *User) UserParamsCheck(c *gin.Context, item *schema.User) error {
	// 检查用户名或密码是否为空
	item.UserName = strings.TrimSpace(item.UserName)
	item.Password = strings.TrimSpace(item.Password)
	if item.UserName == "" || item.Password == "" {
		return errors.New("用户名或密码不能为空")
	}

	// 检查用户名是否存在
	UserQueryResult, err := a.Query(c, schema.UserQueryParam{
		UserName:     item.UserName,
		OmitPassword: true,
	})
	if err != nil {
		return errors.WithMessage(err, "检查用户名是否存在失败")
	} else if len(UserQueryResult.Data) != 0 {
		if UserQueryResult.Data[0].ID != item.ID {
			return errors.New("该用户已存在")
		}
	}

	// 检查组织是否被禁用
	if len(item.Organizations) != 0 {
		if OrgQueryResult, err := a.OrganizationModel.Query(schema.OrganizationQueryParam{
			IDs:    common.UintSliceToString(item.Organizations.GetIDs(), ","),
			Status: consts.BaseStatusDisabled,
		}); err != nil {
			return err
		} else if len(OrgQueryResult.Data) != 0 {
			var orgName []string
			for _, orgInfo := range OrgQueryResult.Data {
				orgName = append(orgName, orgInfo.Name)
			}
			return errors.New(fmt.Sprintf("组织: %s 已被禁用", common.StringSliceToString(orgName, ",")))
		}
	}

	// 检查职位是否被禁用
	if len(item.Positions) != 0 {
		if PositionQueryResult, err := a.PositionModel.Query(schema.PositionQueryParam{
			IDs:    common.UintSliceToString(item.Positions.GetIDs(), ","),
			Status: consts.BaseStatusDisabled,
		}); err != nil {
			return err
		} else if len(PositionQueryResult.Data) != 0 {
			var positionNames []string
			for _, position := range PositionQueryResult.Data {
				positionNames = append(positionNames, position.Name)
			}
			return errors.New(fmt.Sprintf("职位: %s 已被禁用", common.StringSliceToString(positionNames, ",")))
		}
	}

	// 检查角色是否被禁用、类型是否正确
	if len(item.Roles) != 0 {
		if RoleQueryResult, err := a.RoleModel.Query(schema.RoleQueryParam{
			IDs:    common.UintSliceToString(item.Roles.GetIDs(), ","),
			Status: consts.BaseStatusDisabled,
		}); err != nil {
			return err
		} else if len(RoleQueryResult.Data) != 0 {
			var disabledRoleNames []string
			var ErrTypeRoleNames []string
			for _, roleInfo := range RoleQueryResult.Data {
				if roleInfo.Status == consts.BaseStatusDisabled {
					disabledRoleNames = append(disabledRoleNames, roleInfo.Name)
				}
				if roleInfo.Type != consts.RoleTypeForUser {
					ErrTypeRoleNames = append(ErrTypeRoleNames, roleInfo.Name)
				}
			}
			if len(disabledRoleNames) != 0 {
				return errors.New(fmt.Sprintf("角色: %s 已被禁用", common.StringSliceToString(disabledRoleNames, ",")))
			}
			if len(ErrTypeRoleNames) != 0 {
				return errors.New(fmt.Sprintf("角色: %s 的角色类型错误", common.StringSliceToString(ErrTypeRoleNames, ",")))
			}
		}
	}

	// 检查用户组是否被禁用
	if len(item.UserGroups) != 0 {
		if UserGroupQueryResult, err := a.UserGroupModel.Query(schema.UserGroupQueryParam{
			IDs:    common.UintSliceToString(item.UserGroups.GetIDs(), ","),
			Status: consts.BaseStatusDisabled,
		}); err != nil {
			return errors.WithMessage(err, "检查用户组是否被禁用失败")
		} else if len(UserGroupQueryResult.Data) != 0 {
			var userGroupNames []string
			for _, userGroup := range UserGroupQueryResult.Data {
				userGroupNames = append(userGroupNames, userGroup.Name)
			}
			return errors.New(fmt.Sprintf("角色组: %s 已被禁用", common.StringSliceToString(userGroupNames, ",")))
		}
	}
	return nil
}

// Verify 检查用户是否存在
func (a *User) Verify(c *gin.Context, userName, password string) (*schema.User, error) {
	UserQueryResult, err := a.Query(c, schema.UserQueryParam{
		UserName: userName,
	})
	if err != nil {
		return nil, errors.WithMessage(err, "检查用户是否存在失败")
	} else if len(UserQueryResult.Data) == 0 {
		return nil, errors.New("未找到该用户")
	}

	item := UserQueryResult.Data[0]
	if item.Password != secret.SHA1String(password) {
		return nil, errors.New("用户密码不正确，请使用正确的密码")
	} else if item.Status != consts.BaseStatusDisabled {
		return nil, errors.New("用户已被禁用")
	}
	return item, nil
}
