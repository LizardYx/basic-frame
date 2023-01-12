package bll

import (
	"basic-frame/modules/base/dao/model"
	"basic-frame/modules/base/schema"
	"basic-frame/util/common"
	"basic-frame/util/ginx/errors"
	"github.com/gin-gonic/gin"
)

var SystemConfigBll = &SystemConfig{
	SystemConfigModel: model.SystemConfigModel,
}

type SystemConfig struct {
	SystemConfigModel *model.SystemConfig
}

// Init 初始化系统基础配置项
func (a *SystemConfig) Init() error {
	// 检查系统基础配置项是否存在
	if systemConfigQueryResult, err := a.Query(&gin.Context{}, schema.SystemConfigQueryParam{}); err != nil {
		return errors.WithMessage(err, "检查系统基础配置项是否存在失败")
	} else if len(systemConfigQueryResult.Data) == 0 {
		// 如果数据库中没有系统基础配置项,参数检查
		if common.SysConfig.WebServer.Host == "" || common.SysConfig.WebServer.Port == 0 {
			return errors.New("Web服务的主机地址、端口不能为空")
		}

		// 创建系统基础配置项
		if _, err = a.SystemConfigModel.Create(schema.SystemConfig{
			WebServerHost: common.SysConfig.WebServer.Host,
			WebServerPort: common.SysConfig.WebServer.Port,
			HttpsMode:     common.SysConfig.WebServer.HttpsMode,
			HttpsCrtFile:  common.SysConfig.WebServer.HttpsCrtFile,
			HttpsKeyFile:  common.SysConfig.WebServer.HttpsKeyFile,
			MenuVersion:   common.SysConfig.MenuVersion,
		}); err != nil {
			return errors.WithMessage(err, "检查系统基础配置项是否存在失败")
		}
	} else {
		// 如果数据库中有系统基础配置项,将数据写入全局缓存中
		basicConfig := systemConfigQueryResult.Data[0]
		common.SysConfig.WebServer = common.WebServer{
			Host:         basicConfig.WebServerHost,
			Port:         basicConfig.WebServerPort,
			HttpsMode:    basicConfig.HttpsMode,
			HttpsCrtFile: basicConfig.HttpsCrtFile,
			HttpsKeyFile: basicConfig.HttpsKeyFile,
		}
		common.SysConfig.MenuVersion = basicConfig.MenuVersion
	}
	return nil
}

func (a *SystemConfig) Query(c *gin.Context, params schema.SystemConfigQueryParam) (*schema.SystemConfigQueryResult, error) {
	// 创建系统基础配置项
	return a.SystemConfigModel.Query(params)
}

func (a *SystemConfig) Update(c *gin.Context, id uint64, item schema.SystemConfig) error {
	// 参数检查
	if err := a.ParamsValidate(item); err != nil {
		return err
	}

	// 检查系统基础配置项是否存在
	if oldItem, err := a.SystemConfigModel.Get(id); err != nil {
		return errors.WithMessage(err, "获取系统基础配置项信息失败")
	} else if oldItem == nil {
		return errors.New("系统基础配置项不存在")
	}

	// 更新系统基础配置项
	if err := a.SystemConfigModel.UpdateByID(id, map[string]interface{}{
		"web_server_host": item.WebServerHost,
		"web_server_port": item.WebServerPort,
		"https_mode":      item.HttpsMode,
		"https_crt_file":  item.HttpsCrtFile,
		"https_key_file":  item.HttpsKeyFile,
	}); err != nil {
		return errors.WithMessage(err, "更新系统基础配置项失败")
	}

	return nil
}

// ---------------------------------------- Params  Validate --------------------------------------

// ParamsValidate 参数检查
func (a *SystemConfig) ParamsValidate(item schema.SystemConfig) error {
	// 如果启用Https模式
	if item.HttpsMode {
		if item.WebServerHost == "" || item.WebServerPort == 0 {
			return errors.New("启用HTTPS模式时,Web服务IP地址、端口 不能为空")
		}
	}
	return nil
}
