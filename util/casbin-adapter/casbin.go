package casbin_adapter

import (
	"basic-frame/modules/base/dao/model"
	"basic-frame/util/common"
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
	e.StartAutoLoadPolicy(time.Duration(60) * time.Second)
	common.SysConfig.CasbinSyncEnforcer = e
	return nil
}

// ---------------------------------------- Casbin Adapter --------------------------------------

var _ persist.Adapter = (*CasbinAdapter)(nil)

type CasbinAdapter struct {
	ButtonModel *model.Button
}

func (a *CasbinAdapter) LoadPolicy(model casbinModel.Model) error {
	// load casbin rule
	// TODO:

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
