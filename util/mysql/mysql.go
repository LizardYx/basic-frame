package mysql

import (
	BaseEntity "basic-frame/modules/base/dao/entity"
	"basic-frame/util/common"
	"basic-frame/util/logger"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	GormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"log"
	"os"
	"strings"
	"time"
)

var DB *gorm.DB

// InitMysql 初始化Mysql数据库
func InitMysql() {
	var err error

	if DB, err = gorm.Open(mysql.Open(common.SysConfig.Mysql.GetDBConnString()), &gorm.Config{
		Logger: gormLogger(),
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   fmt.Sprintf("%s_", strings.ToLower(common.SysConfig.AppName)), // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: false,                                                         // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
	}); err != nil {
		logger.Log.Warn("mysql connect error: ", err)
		log.Fatal("mysql connect error: ", err)
		return
	}

	if err := autoMigrate(DB); err != nil {
		logger.Log.Warn("database autoMigrate error: ", err)
		log.Fatal("database autoMigrate error: ", err)
		return
	}
}

// autoMigrate 自动迁移 schema，保持 schema 是最新的。
func autoMigrate(db *gorm.DB) error {
	db = db.Set("gorm:table_options", "ENGINE=InnoDB")
	entities := make([]interface{}, 0)
	entities = append(entities, BaseEntity.GetTables()...)

	return db.AutoMigrate(entities...)
}

// gormLogger gorm默认logger实现，打印慢 SQL 和错误
func gormLogger() GormLogger.Interface {
	return GormLogger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		GormLogger.Config{
			SlowThreshold: time.Second,     // 慢 SQL 阈值
			LogLevel:      GormLogger.Info, // Log level
			Colorful:      true,            // 禁用彩色打印
		},
	)
}

// Paginate 分页查询
func Paginate(db *gorm.DB, params common.PaginationParam, out interface{}) (*common.PaginationResult, error) {
	var count int64
	err := db.Count(&count).Error
	if err != nil {
		return nil, err
	}

	if params.OnlyCount {
		// 仅查询count
		return &common.PaginationResult{Total: count}, nil
	} else if !params.Pagination {
		// 查询所有数据
		err = db.Find(out).Error
		return nil, err
	} else {
		// 分页查询
		if params.Current == 0 {
			params.Current = 1
		} else if params.PageSize <= 0 {
			params.PageSize = 10
		}
		db = db.Offset((params.Current - 1) * params.PageSize).Limit(params.PageSize)
		err = db.Find(out).Error
		return &common.PaginationResult{
			Total:    count,
			Current:  params.Current,
			PageSize: params.PageSize,
		}, nil
	}
}

// FindOne 查询单条数据
func FindOne(db *gorm.DB, out interface{}) (bool, error) {
	result := db.First(out)
	if err := result.Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// Check 检查数据是否存在
func Check(db *gorm.DB) (bool, error) {
	var count int64
	result := db.Count(&count)
	if err := result.Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
