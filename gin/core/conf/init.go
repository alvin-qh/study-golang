package conf

import (
	"encoding/json"
	"fmt"
	"strings"
	"study-gin/core/utils"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// 配置结构体
type _Config struct {
	// 服务器配置
	Server struct {
		// 服务地址
		Address string `mapstructure:"address"`

		// 跨域配置
		Cors struct {
			Enable           bool `mapstructure:"enable"`
			AllowOrigin      any  `mapstructure:"allow-origin"`
			AllowMethods     any  `mapstructure:"allow-methods"`
			AllowHeaders     any  `mapstructure:"allow-headers"`
			ExposeHeaders    any  `mapstructure:"expose-headers"`
			AllowCredentials bool `mapstructure:"allow-credentials"`
			MaxAge           bool `mapstructure:"max-age"`
		} `mapstructure:"cors"`
	} `mapstructure:"server"`

	// 日志配置
	Logger struct {
		File       string `mapstructure:"file"`
		Level      string `mapstructure:"level"`
		ShowCaller bool   `mapstructure:"show-caller"`
		MaxSize    int    `mapstructure:"max-size"`
		Compress   bool   `mapstructure:"compress"`
		MaxAge     int    `mapstructure:"max-age"`
		MaxBackups int    `mapstructure:"max-backup"`
	} `mapstructure:"logger"`
}

var (
	Config = &_Config{}
)

func configToJson() string {
	data, err := json.MarshalIndent(viper.AllSettings(), "", "  ")
	if err != nil {
		log.Fatalf("invalid config caused %v", err)
	}
	return fmt.Sprintf("\n%v", string(data))
}

func setDefaultConfig() {
	viper.SetDefault("server.address", "0.0.0.0:8080")

	viper.SetDefault("server.cors.allow-origin", "*")
	viper.SetDefault("server.cors.allow-methods", []string{
		"GET",
		"PUT",
		"POST",
		"DELETE",
		"OPTIONS",
		"HEAD",
		"PATCH",
	})
	viper.SetDefault("server.cors.allow-headers", "*")
	viper.SetDefault("server.cors.expose-headers", "*")
	viper.SetDefault("server.cors.allow-credentials", true)
	viper.SetDefault("server.cors.max-age", 86400)

	viper.SetDefault("logger.file", "logs/server.log")
	viper.SetDefault("logger.level", "INFO")
	viper.SetDefault("logger.show-caller", true)
	viper.SetDefault("logger.max-size", 100)
	viper.SetDefault("logger.compress", true)
	viper.SetDefault("logger.max-age", 30)
	viper.SetDefault("logger.max-backup", 100)
}

func Init(filename string) {
	setDefaultConfig()

	viper.SetConfigType("yaml")
	viper.SetConfigFile(filename)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("cannot load config file \"%v\", caused %v", filename, err)
	}
	if err := viper.Unmarshal(Config); err != nil {
		log.Fatalf("cannot load config file \"%v\", caused %v", filename, err)
	}
	log.Infof("config file \"%v\" read, content is %v", filename, configToJson())
}

type _Value interface {
	~string |
		~int |
		~int64 |
		~int32 |
		~float32 |
		~float64
}

func Default[T _Value](val T, defVal T) T {
	switch t := (any(val)).(type) {
	case string:
		if len(t) == 0 {
			return defVal
		}
		return val
	default:
		if t == 0 {
			return defVal
		}
		return val
	}
}

func ToString(val any, sep string) (result string) {
	switch sval := val.(type) {
	case string:
		result = sval
	case []string:
		result = strings.Join(sval, sep)
	case []any:
		result = strings.Join((func() []string {
			lst := make([]string, len(sval))
			for i, v := range sval {
				lst[i] = utils.AnyToString(v)
			}
			return lst
		})(), sep)
	default:
		log.Fatalf("invalid config value \"%v\"", val)
	}
	return
}
