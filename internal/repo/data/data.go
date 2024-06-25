package data

import (
	"errors"
	"github.com/go-grain/gin-layout/config"
	sysModel "github.com/go-grain/gin-layout/model/system"
	"gorm.io/gorm"
)

const (
	// dbMySQL Gorm Drivers mysql || postgres || sqlite || sqlserver
	dbMySQL    string = "mysql"
	dbPostgres string = "postgres"
	dbTidb     string = "tidb"
)

var db *DB

type DB struct {
	DB *gorm.DB
}

func InitDB(conf config.Config) (*gorm.DB, error) {
	switch conf.DataBase.Driver {
	case dbMySQL:
		mysql, err := InitMysql(conf)
		if err != nil {
			return nil, err
		}
		return mysql, err
	default:
		return nil, errors.New("数据库配置有问题")
	}
}

func NewDB() *DB {
	return db
}

func (db *DB) autoMigrate() error {
	err := db.DB.AutoMigrate(
		sysModel.SysRole{},
		sysModel.SysUser{},
	)
	if err != nil {
		return err
	}

	return nil
}
