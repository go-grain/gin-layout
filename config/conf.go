package config

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
	"os"
	"time"
)

//go:embed config.yaml
var fileYaml string

type Gin struct {
	Host  string `mapstructure:"host" json:"host" yaml:"host"`
	Model string `mapstructure:"model" json:"model" yaml:"model"`
}

type JWT struct {
	SecretKey         string `mapstructure:"secret_key" json:"secret_key" yaml:"secret_key"`
	ExpirationSeconds int64  `mapstructure:"expiration_seconds" json:"expiration_seconds" yaml:"expiration_seconds"`
	Issuer            string `mapstructure:"issuer" json:"issuer" yaml:"issuer"`
}

type DataBase struct {
	Driver   string          `mapstructure:"driver" json:"driver" yaml:"driver"`
	LogLevel logger.LogLevel `mapstructure:"log_level" json:"log_level" yaml:"log_level"`
	MySql    struct {
		Source string `mapstructure:"source" json:"source" yaml:"source"`
	} `mapstructure:"mysql" json:"mysql" yaml:"mysql"`
	Redis struct {
		UserName     string        `mapstructure:"user_name" json:"user_name" yaml:"user_name"`
		Password     string        `mapstructure:"password" json:"password" yaml:"password"`
		Addr         string        `mapstructure:"addr" json:"addr" yaml:"addr"`
		DB           int           `mapstructure:"db" json:"db" yaml:"db"`
		ReadTimeout  time.Duration `mapstructure:"read_timeout" json:"read_timeout" yaml:"read_timeout"`
		WriteTimeout time.Duration `mapstructure:"write_timeout" json:"write_timeout" yaml:"write_timeout"`
	} `yaml:"redis"`
}

type Log struct {
	Level     zapcore.Level `mapstructure:"level" json:"level" yaml:"level"`
	LogPath   string        `mapstructure:"log_path" json:"log_path" yaml:"log_path"`
	SplitSize int           `mapstructure:"split_size" json:"split_size" yaml:"split_size"`
}

type Config struct {
	Gin      Gin      `mapstructure:"gin" json:"gin" yaml:"gin"`
	Log      Log      `mapstructure:"log" json:"log" yaml:"log"`
	DataBase DataBase `mapstructure:"database" json:"database" yaml:"database"`
	JWT      JWT      `mapstructure:"jwt" json:"jwt" yaml:"jwt"`
	Role     Role     `mapstructure:"role" json:"role" yaml:"role"`
}

type Role struct {
	DefaultRole      string `mapstructure:"default_role" json:"default_role" yaml:"default_role"`
	DefaultAdminRole string `mapstructure:"default_admin_role" json:"default_admin_role" yaml:"default_admin_role"`
}

func InitConfig() (*Config, error) {
	// 如果没有配置文件 则导出一份默认配置文件到本地
	_, err := os.Stat("config/config.yaml")
	if os.IsNotExist(err) {
		err = os.MkdirAll("config", os.ModePerm)
		if err != nil {
			return nil, errors.New("配置文件初始化失败")
		}
		err = os.WriteFile("config/config.yaml", []byte(fileYaml), os.ModePerm)
		if err != nil {
			return nil, errors.New("配置文件初始化失败")
		}
	}

	var conf *viper.Viper
	conf = viper.New()
	conf.SetConfigFile("config/config.yaml")
	err = conf.ReadInConfig()
	if err != nil {
		panic(any(err.Error()))
		return nil, err
	}

	c := Config{}
	err = conf.Unmarshal(&c)
	if err != nil {
		panic(any(err.Error()))
		return nil, err
	}

	conf.WatchConfig()
	conf.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("配置文件已更改: ", e.Name)
		if err = conf.Unmarshal(&c); err != nil {
			fmt.Println(err)
		}
	})
	return &c, nil
}
