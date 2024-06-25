package data

import (
	"errors"
	"fmt"
	"github.com/go-grain/gin-layout/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func InitMysql(conf config.Config) (*gorm.DB, error) {
	out := &os.File{}
	if conf.Gin.Model == "debug" {
		out = os.Stdout
	} else {
		var err error
		out, err = os.OpenFile("log/mysql.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			out = os.Stdout
		}
	}
	newLogger := logger.New(
		log.New(out, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,    // 慢 SQL 阈值
			LogLevel:                  conf.DataBase.LogLevel,    // 日志级别
			IgnoreRecordNotFoundError: true,                      // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  conf.Gin.Model == "debug", // 禁用彩色打印
		},
	)

	mysqlConfig := mysql.Config{
		DSN:                       conf.DataBase.MySql.Source, // DSN data source name
		DefaultStringSize:         191,                        // string 类型字段的默认长度
		DisableDatetimePrecision:  false,                      // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                       // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                       // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                      // 根据版本自动配置
	}
	gormDB, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{Logger: newLogger})
	if err != nil {
		return nil, errors.New(err.Error())
	}

	sqlDB, _ := gormDB.DB()
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(500)
	sqlDB.SetConnMaxIdleTime(time.Second * 5)
	sqlDB.SetConnMaxLifetime(time.Hour)

	db = &DB{DB: gormDB}
	err = db.autoMigrate()
	if err != nil {
		fmt.Println("MySQL AutoMigrate error", err.Error())
		return nil, err
	}

	fmt.Println("初始化MySql成功")

	return gormDB, nil
}
