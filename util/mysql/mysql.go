package mysql

import (
	"basic-frame/util/common"
	"basic-frame/util/logger"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"strings"
)

var DB *gorm.DB

// InitMysql 初始化Mysql数据库
func InitMysql() {
	var err error

	if DB, err = gorm.Open("mysql", common.SysConfig.Mysql.GetDBConnString()); err != nil {
		logger.Log.Warn("mysql connect error: ", err)
		log.Fatal("mysql connect error: ", err)
		return
	} else if DB.Error != nil {
		logger.Log.Warn("database error: ", err)
		log.Fatal("database error: ", err)
		return
	}
	defer DB.Close()
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return strings.ToLower(common.SysConfig.AppName) + "_" + defaultTableName
	}
	DB.AutoMigrate()
}
