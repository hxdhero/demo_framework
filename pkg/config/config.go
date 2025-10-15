package config

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"lls_api/pkg/util/validate"
	"log"
	"strings"
)

const (
	EnvProd = "prod"
	EnvUat  = "uat"
	EnvDev  = "dev"
)

type Config struct {
	Env  string
	Http struct {
		Host string
		Port int `validate:"required"`
	}
	DB struct {
		Dsn         string `mapstructure:"dsn"`
		MaxIdleConn int    `mapstructure:"max_idle_conn"`
		MaxOpenConn int    `mapstructure:"max_open_conn"`
		Debug       bool   `mapstructure:"debug"`
	}
	Redis struct {
		Prefix string `mapstructure:"prefix" validate:"required"`
		Host   string `mapstructure:"host" validate:"required"`
		Pwd    string `mapstructure:"pwd" validate:"required"`
		DB     int    `mapstructure:"db" `
	}
	OSS struct {
		AccessKeyID     string `mapstructure:"access_key_id" validate:"required"`
		AccessKeySecret string `mapstructure:"access_key_secret" validate:"required"`
		Endpoint        string `mapstructure:"endpoint" validate:"required"`
		Bucket          string `mapstructure:"bucket" validate:"required"`
		Prefix          string
	}
	BaiDuOCR struct {
		ClientId string `mapstructure:"client_id"`
		Secret   string `mapstructure:"secret"`
	} `mapstructure:"baidu_ocr"`
	ServerlessForFilesKey string `mapstructure:"serverless_for_files_key" validate:"required"`
	JWT                   struct {
		Secret string `mapstructure:"secret"`
		Exp    int    `mapstructure:"exp"`
	} `mapstructure:"jwt"`
	Notification struct {
		Feishu struct {
			AppID      string `mapstructure:"app_id"`
			Secret     string `mapstructure:"secret"`
			ErrChatID  string `mapstructure:"err_chat_id"`
			TemplateID string `mapstructure:"template_id"`
		} `mapstructure:"feishu"`
	} `mapstructure:"notification"`
}

var (
	C *Config
)

func InitConfig() {
	if C != nil {
		return
	}
	v := viper.New()
	// 设置config文件路径
	appendConfigPath(v)

	// 加载基础配置
	v.SetConfigName("config")
	v.SetConfigType("yml")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	log.Default().Println("加载了配置文件: " + v.ConfigFileUsed())

	// 加载本地配置
	v.SetConfigName("config_local")
	v.SetConfigType("yml")
	err = v.MergeInConfig()
	if err != nil {
		var configNotFoundErr viper.ConfigFileNotFoundError
		if errors.As(err, &configNotFoundErr) {
			fmt.Println("WARN!!! 没有找到config_local文件: " + configNotFoundErr.Error())
		} else {
			panic(err)
		}
	} else {
		log.Default().Println("加载了配置文件: " + v.ConfigFileUsed())
	}
	// 优先级最高,环境变量
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	c := &Config{}
	if err := v.Unmarshal(c); err != nil {
		panic(err)
	}
	if err := validate.V.Struct(c); err != nil {
		panic(err)
	}

	// 全局config 赋值
	C = c
	return
}

func appendConfigPath(v *viper.Viper) {
	paths := []string{
		"./config",
		"../config",
		"../../config",
		"../../../config",
	}
	for _, e := range paths {
		v.AddConfigPath(e)
	}
}
