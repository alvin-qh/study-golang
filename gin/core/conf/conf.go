package conf

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// 配置结构体
type Config struct {
	// 服务器配置
	Server struct {
		// 服务地址
		Address []string `mapstructure:"address"`

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

const (
	CONF_FILE = "./application.yaml"
)

func configToJson() string {
	data, err := json.MarshalIndent(viper.AllSettings(), "", "  ")
	if err != nil {
		panic(err)
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

func Load(filename string) (*Config, error) {
	setDefaultConfig()

	viper.SetConfigType("yaml")
	viper.SetConfigFile(filename)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	log.Infof("config file \"%v\" read, content is %v", filename, configToJson())

	conf := &Config{}
	err := viper.Unmarshal(conf)

	return conf, err
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

func ToString(val any, sep string) string {
	switch sval := val.(type) {
	case string:
		return sval
	case []string:
		return strings.Join(sval, sep)
	case []any:
		return strings.Join((func() []string {
			result := make([]string, len(sval))
			for i, v := range sval {
				switch sv := v.(type) {
				case string:
					result[i] = sv
				case int:
					result[i] = strconv.FormatInt(int64(sv), 10)
				case int64:
					result[i] = strconv.FormatInt(sv, 10)
				case byte:
					result[i] = strconv.FormatUint(uint64(sv), 10)
				case uint:
					result[i] = strconv.FormatUint(uint64(sv), 10)
				case uint64:
					result[i] = strconv.FormatUint(uint64(sv), 10)
				default:
					result[i] = fmt.Sprintf("%v", sv)
				}
			}
			return result
		})(), sep)
	default:
		log.Panicf("invalid config value \"%v\"", val)
		return ""
	}
}
