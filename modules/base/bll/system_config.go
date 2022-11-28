package bll

import (
	"basic-frame/modules/base/dao/model"
	"basic-frame/modules/base/schema"
	"context"
)

var SystemConfigBll = &SystemConfig{
	SystemConfigModel: model.SystemConfigModel,
}

type SystemConfig struct {
	SystemConfigModel *model.SystemConfig
}

func (a *SystemConfig) Query(params schema.SystemConfigQueryParam) (*schema.SystemConfigQueryResult, error) {
	// 创建系统基础配置项
	return a.SystemConfigModel.Query(params)
}

func (a *SystemConfig) Create(ctx context.Context, item schema.SystemConfig) (uint64, error) {
	// 参数检查
	if err := a.ParamsValidate(ctx, item); err != nil {
		return 0, err
	}

	// 创建系统基础配置项
	return a.SystemConfigModel.Create(item)
}

// ---------------------------------------- Params  Validate --------------------------------------

func (a *SystemConfig) ParamsValidate(ctx context.Context, item schema.SystemConfig) error {
	// 检查系统基础配置项是否存在
	if systemConfigQueryResult, err := a.Query(schema.SystemConfigQueryParam{}); err != nil {
		return err
	} else if len(systemConfigQueryResult.Data) != 0 {
		return err
	}

	return nil
}
