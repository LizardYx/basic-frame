package casbin_adapter

import (
	"basic-frame/modules/base/dao/model"
	"github.com/casbin/casbin/v2"
	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"time"
)

type CasbinAdapter struct {
	ButtonModel *model.Button
}

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

func (a *CasbinAdapter) LoadPolicy(model casbinModel.Model) error {
	// load casbin rule
	// TODO:

	return nil
}

// InitCasbin 初始化Casbin
func InitCasbin(adapter persist.Adapter) (*casbin.SyncedEnforcer, error) {
	// 初始化Casbin模型
	m, err := casbinModel.NewModelFromString(modelString)
	if err != nil {
		return nil, err
	}

	// 创建Enforcer对象
	e, err := casbin.NewSyncedEnforcer(m)
	if err != nil {
		return nil, err
	}

	err = e.InitWithModelAndAdapter(e.GetModel(), adapter)
	if err != nil {
		return nil, err
	}
	// 自动更新所有权限
	e.StartAutoLoadPolicy(time.Duration(60) * time.Second)
	return e, nil
}
