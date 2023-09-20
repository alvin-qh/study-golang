package conf

import (
	"encoding/json"
	"fmt"
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
			AllowedOrigins any `mapstructure:"allowed-origins"`
			AllowedMethods any `mapstructure:"allowed-methods"`
		} `mapstructure:"cors"`
	} `mapstructure:"server"`

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

func Load(filename string) (*Config, error) {
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

func ToString(val any) string {
	switch sval := val.(type) {
	case string:
		return sval
	case []any:
		return strings.Join((func() []string {
			result := make([]string, len(sval))
			for i, v := range sval {
				switch sv := v.(type) {
				case string:
					result[i] = sv
				default:
					result[i] = fmt.Sprintf("%v", sv)
				}
			}
			return result
		})(), ",")
	default:
		log.Panicf("invalid config value \"%v\"", val)
		return ""
	}
}
