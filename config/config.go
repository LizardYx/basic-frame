package config

import (
	"basic-frame/util/common"
	"basic-frame/util/consts"
	"fmt"
	"gopkg.in/ini.v1"
)

// LoadSystemConfigFile 加载系统配置文件
func LoadSystemConfigFile() (string, error) {
	// 读取系统配置文件
	cfg, err := ini.Load(consts.SystemConfigFileName)
	if err != nil {
		// 未找到系统配置文件，创建默认系统配置文件
		var errMessage string
		if errMessage, err = createSystemConfigFile(); err != nil {
			return errMessage, err
		} else {
			cfg, _ = ini.Load(consts.SystemConfigFileName)
		}
	}
	common.SysConfig.RunMode = cfg.Section("").Key("RunMode").String()
	common.SysConfig.AppName = cfg.Section("").Key("AppName").String()
	common.SysConfig.AppVersion = cfg.Section("").Key("AppVersion").MustFloat64()
	common.SysConfig.DefaultLang = cfg.Section("").Key("DefaultLang").String()
	// 读取mysql配置部分
	if err = cfg.Section("Mysql").MapTo(&common.SysConfig.Mysql); err != nil {
		return "failed to load Mysql config :", err
	}
	// 读取webServer配置部分
	if err = cfg.Section("WebServer").MapTo(&common.SysConfig.WebServer); err != nil {
		return "failed to load WebServer config :", err
	}
	// 读取JWTAuth配置部分
	if err = cfg.Section("JWTAuth").MapTo(&common.SysConfig.JWTAuth); err != nil {
		return "failed to load JWTAuth config :", err
	}
	// 读取Logger配置部分
	if err = cfg.Section("Logger").MapTo(&common.SysConfig.Logger); err != nil {
		return "failed to load Logger config :", err
	}
	// 读取SmtpServer配置部分
	if err = cfg.Section("SmtpServer").MapTo(&common.SysConfig.SmtpServer); err != nil {
		return "failed to load SmtpServer config :", err
	}
	return "", nil
}

// CreateSystemConfigFile 创建默认的系统配置文件
func createSystemConfigFile() (string, error) {
	var err error

	cfg := ini.Empty()
	defaultSection := cfg.Section("")
	defaultSection.Comment = getRunModelDesc()
	defaultSection.NewKey("RunMode", consts.RunModeDev)
	defaultSection.NewKey("AppName", consts.AppName)
	defaultSection.NewKey("AppVersion", consts.AppVersion)
	defaultSection.NewKey("DefaultLang", consts.DefaultLang)

	// 创建默认mysql配置
	var mysqlSection *ini.Section
	if mysqlSection, err = cfg.NewSection("Mysql"); err != nil {
		return fmt.Sprintf("init Mysql config to %s failed: ", consts.SystemConfigFileName), err
	} else {
		mysqlSection.NewKey("Host", consts.MysqlHost)
		mysqlSection.NewKey("Port", consts.MysqlPort)
		mysqlSection.NewKey("User", consts.MysqlUser)
		mysqlSection.NewKey("Password", consts.MysqlPassword)
		mysqlSection.NewKey("DBName", consts.MysqlDBName)
	}

	// 创建默认webServer配置
	var webServerSection *ini.Section
	if webServerSection, err = cfg.NewSection("WebServer"); err != nil {
		return fmt.Sprintf("init WebServer config to %s failed: ", consts.SystemConfigFileName), err
	} else {
		webServerSection.NewKey("Host", consts.WebServerHost)
		webServerSection.NewKey("Port", consts.WebServerPort)
		webServerSection.NewKey("HttpsMode", consts.UseHttpsMode)
		webServerSection.NewKey("HttpsCrtFile", consts.HttpsCrtFile)
		webServerSection.NewKey("HttpsKeyFile", consts.HttpsKeyFile)
	}

	// 创建默认JWTAuth配置
	var JWTAuthSection *ini.Section
	if JWTAuthSection, err = cfg.NewSection("JWTAuth"); err != nil {
		return fmt.Sprintf("init JWTAuth config to %s failed: ", consts.SystemConfigFileName), err
	} else {
		JWTAuthSection.NewKey("SigningMethod", consts.JWTAuthSigningMethod)
		JWTAuthSection.NewKey("SigningKey", consts.JWTAuthSigningKey)
		JWTAuthSection.NewKey("Expired", consts.JWTAuthExpired)
	}

	// 创建默认Logger配置
	var LoggerSection *ini.Section
	if LoggerSection, err = cfg.NewSection("Logger"); err != nil {
		return fmt.Sprintf("init Logger config to %s failed: ", consts.SystemConfigFileName), err
	} else {
		LoggerSection.NewKey("RotationTime", consts.LoggerRotationTime)
		LoggerSection.NewKey("MaxAge", consts.LoggerMaxAge)
	}

	// 创建默认SmtpServer配置
	var SmtpServerSection *ini.Section
	if SmtpServerSection, err = cfg.NewSection("SmtpServer"); err != nil {
		return fmt.Sprintf("init SmtpServer config to %s failed: ", consts.SystemConfigFileName), err
	} else {
		SmtpServerSection.NewKey("Enable", consts.SmtpEnable)
		SmtpServerSection.NewKey("EmailServerName", consts.SmtpServerName)
		SmtpServerSection.NewKey("SmtpHost", consts.SmtpHost)
		SmtpServerSection.NewKey("SmtpPort", consts.SmtpPort)
		SmtpServerSection.NewKey("SmtpPassword", consts.SmtpPassword)
	}

	if err = cfg.SaveTo(consts.SystemConfigFileName); err != nil {
		return fmt.Sprintf("Create %s file failed: ", consts.SystemConfigFileName), err
	}
	return "", nil
}

func getRunModelDesc() string {
	return fmt.Sprintf("运行模式(%s:开发模式,%s:测试模式,%s:线上模式)", consts.RunModeDev, consts.RunModeTest, consts.RunModePro)
}
