package casbin_adapter

import (
	"basic-frame/modules/base/dao/model"
	baseSchema "basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/consts"
	"basic-frame/util/ginx/errors"
	"fmt"
	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"time"
)

var modelString = `
				[request_definition]
				r = sub, obj, act
				
				[policy_definition]
				p = sub, obj, act
				
				[role_definition]
				g = _, _
				
				[policy_effect]
				e = some(where (p.eft == allow))
				
				[matchers]
				m = g(r.sub, p.sub) == true \
					&& keyMatch2(r.obj, p.obj) == true \
					&& regexMatch(r.act, p.act) == true \
					|| r.sub == "root"
					`

// InitCasbin 初始化Casbin
func InitCasbin() error {
	// 初始化Casbin模型
	m, err := casbinModel.NewModelFromString(modelString)
	if err != nil {
		return err
	}

	// 创建Enforcer对象
	e, err := casbin.NewSyncedEnforcer(m)
	if err != nil {
		return err
	}

	err = e.InitWithModelAndAdapter(e.GetModel(), &CasbinAdapter{})
	if err != nil {
		return err
	}
	// 自动更新所有权限
	e.StartAutoLoadPolicy(time.Duration(consts.CasbinPolicyReloadTime) * time.Second)
	common.SysConfig.CasbinSyncEnforcer = e
	return nil
}

// ---------------------------------------- Casbin Adapter --------------------------------------

var _ persist.Adapter = (*CasbinAdapter)(nil)

type CasbinAdapter struct {
	UserModel         *model.User
	OrganizationModel *model.Organization
	RoleModel         *model.Role
	MenuModel         *model.Menu
	ButtonModel       *model.Button
}

// LoadPolicy 加载casbin权限策略
func (a *CasbinAdapter) LoadPolicy(model casbinModel.Model) error {
	// 获取角色的权限策略
	if err := a.LoadRolePolicy(model); err != nil {
		return errors.WithMessage(err, "加载角色的权限策略失败")
	}

	// 获取用户的权限策略
	if err := a.loadUserPolicy(model); err != nil {
		return errors.WithMessage(err, "加载用户的权限策略失败")
	}
	return nil
}

// LoadRolePolicy 加载角色权限策略
func (a *CasbinAdapter) LoadRolePolicy(model casbinModel.Model) error {
	// 获取所有的角色，以及角色关联的页面、按钮信息
	roleQueryResult, err := a.RoleModel.Query(baseSchema.RoleQueryParam{
		Status:      consts.BaseStatusEnable,
		ShowDetails: true,
		FindAll:     true,
	})
	if err != nil {
		return errors.WithStack(err)
	} else if len(roleQueryResult.Data) != 0 {
		// 遍历所有的角色，获取页面、按钮关联的Api信息
		for _, roleInfo := range roleQueryResult.Data {
			var restfulApis baseSchema.RestfulApis

			// 获取角色关联页面的Api集合
			if len(roleInfo.Menus) != 0 {
				if menuTrees, err := a.MenuModel.GetRoleRestfulApis(roleInfo.Menus.GetIDs()); err != nil {
					return errors.WithStack(err)
				} else if len(*menuTrees) != 0 {
					// 获取菜单页面关联的所有Api(去重)
					restfulApis = menuTrees.GetRestfulApis()
				}
			}

			// 获取角色关联按钮的Api集合
			if len(roleInfo.Buttons) != 0 {
				if buttonPres, err := a.ButtonModel.GetButtonRestfulApis(roleInfo.Buttons.GetIDs()); err != nil {
					return errors.WithStack(err)
				} else if len(*buttonPres) != 0 {
					// 获取按钮关联的所有Api(去重)
					for _, buttonPre := range *buttonPres {
						buttonPre.GetRestfulApis(&restfulApis)
					}
				}
			}

			// 遍历所有的restfulApi，生成角色与Api接口的权限映射
			if len(restfulApis) != 0 {
				for _, restfulApi := range restfulApis {
					line := fmt.Sprintf("p,%d,%s,%s", roleInfo.ID, restfulApi.Path, restfulApi.Method)
					if err = persist.LoadPolicyLine(line, model); err != nil {
						return errors.WithStack(err)
					}
				}
			}
		}
	}
	return nil
}

// loadUserPolicy 加载用户权限策略
func (a *CasbinAdapter) loadUserPolicy(model casbinModel.Model) error {
	// 获取所有启用的用户
	userQueryResult, err := a.UserModel.Query(baseSchema.UserQueryParam{
		Status:       consts.BaseStatusEnable,
		ShowDetails:  true,
		OmitPassword: true,
		FindAll:      true,
	})
	if err != nil {
		return errors.WithStack(err)
	} else if len(userQueryResult.Data) != 0 {
		// 遍历所有的用户
		for _, userInfo := range userQueryResult.Data {
			// 获取用户的所有角色ID
			var roleIds []uint64

			// 获取组织结构的角色ID，如果用户有该组织结构的权限，那么就自动会拥有子组织的权限
			if len(userInfo.Organizations) != 0 {
				organizationQueryResult, err := a.OrganizationModel.Query(baseSchema.OrganizationQueryParam{
					IDs:         common.UintSliceToString(userInfo.Organizations.GetIDs(), ","),
					Status:      consts.BaseStatusEnable,
					ShowDetails: true,
					FindAll:     true,
				})
				if err != nil {
					return errors.WithStack(err)
				} else {
					organizationQueryResult.Data.GetRoleIds(&roleIds)
				}
			}
			userInfo.Positions.GetRoleIds(&roleIds)
			userInfo.UserGroups.GetRoleIds(&roleIds)
			roleIds = append(roleIds, userInfo.Roles.GetIDs()...)

			if len(roleIds) != 0 {
				// 角色去重
				newRoleIds := common.RemoveRepeatedElement(roleIds)
				for _, roleID := range newRoleIds {
					line := fmt.Sprintf("g,%d,%d", userInfo.ID, roleID)
					if err = persist.LoadPolicyLine(line, model); err != nil {
						return errors.WithStack(err)
					}
				}
			}
		}
	}
	return nil
}

// SavePolicy saves all policy rules to the storage.
func (a *CasbinAdapter) SavePolicy(model casbinModel.Model) error {
	return nil
}

// AddPolicy adds a policy rule to the storage.
// This is part of the Auto-Save feature.
func (a *CasbinAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemovePolicy removes a policy rule from the storage.
// This is part of the Auto-Save feature.
func (a *CasbinAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return nil
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
// This is part of the Auto-Save feature.
func (a *CasbinAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return nil
}
