package conf

import (
	"fmt"
	"os"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	DEF_CONF_FILE = "./application.yaml" // 默认配置文件名称
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
			MaxAge           int  `mapstructure:"max-age"`
		} `mapstructure:"cors"`

		Template struct {
			Enable        bool   `mapstructure:"enable"`
			TemplatesPath string `mapstructure:"templates-path"`
			StaticPath    string `mapstructure:"static-path"`
			StaticBaseURI string `mapstructure:"static-base-uri"`
		} `mapstructure:"template"`
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
	// 实例化配置结构体实例
	Config = &_Config{}
)

// 将配置内容转为字符串
func configToString() string {
	// data, err := json.MarshalIndent(viper.AllSettings(), "", "  ")
	// if err != nil {
	//    log.Fatalf("invalid config caused %v", err)
	// }
	// return fmt.Sprintf("\n%v", string(data))

	var fn func(string, map[string]any) string
	fn = func(prefix string, m map[string]any) string {
		sb := strings.Builder{}

		for k, v := range m {
			if sv, ok := v.(map[string]any); ok {
				sb.WriteString(fn(fmt.Sprintf("%v%v.", prefix, k), sv))
			} else {
				sb.WriteString(fmt.Sprintf("%v%v", prefix, k))
				sb.WriteString("=")
				sb.WriteString(fmt.Sprintf("%v\n", v))
			}
		}
		return sb.String()
	}

	return fn("\t", viper.AllSettings())
}

// 设置默认配置项
//
// 默认配置项是指: 当配置文件不包含指定项时, 该项默认的值
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

	viper.SetDefault("server.template.enable", false)
	viper.SetDefault("server.template.templates-path", "templates")
	viper.SetDefault("server.template.static-path", "static")
	viper.SetDefault("server.template.static-base-uri", "/static")

	viper.SetDefault("logger.file", "")
	viper.SetDefault("logger.level", "INFO")
	viper.SetDefault("logger.show-caller", true)
	viper.SetDefault("logger.max-size", 100)
	viper.SetDefault("logger.compress", true)
	viper.SetDefault("logger.max-age", 30)
	viper.SetDefault("logger.max-backup", 100)
}

// 初始化配置项
func init() {
	// 设置默认配置项
	setDefaultConfig()

	// 读取配置文件名称
	// 1. 先从环境变量中获取 CONF 变量作为配置文件名
	// 2. 如果环境变量不存在, 则使用默认的配置文件名
	filename, ok := os.LookupEnv("CONF")
	if !ok {
		filename = DEF_CONF_FILE
	}

	// 设置配置文件格式为 yaml 格式
	viper.SetConfigType("yaml")
	// 设置配置文件名称
	viper.SetConfigFile(filename)
	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Errorf("cannot load config file \"%v\", caused %v", filename, err)
	}
	// 将配置文件内容反序列化到 _Config 类型结构体变量
	if err := viper.Unmarshal(Config); err != nil {
		log.Fatalf("cannot load config file \"%v\", caused %v", filename, err)
	}
	log.Infof("config file \"%v\" read, content is\n%v", filename, configToString())
}

// 获取配置项, 如果配置项不存在, 则返回默认值
//
// 参数:
//   - `key` (`string`): 配置项键名称
//   - `def` (`T`): 配置项默认值
//
// 返回:
//   - `val` (`T`): 配置项的值, 如果配置项不存在, 则为 `def` 参数的值
//   - `err` (`error`): 错误对象
func Default[T any](key string, def T) (val T, err error) {
	// 从配置中获取指定 key 的值
	v := viper.Get(key)
	if v == nil {
		// 如果 key 不存在, 则返回默认值
		val = def
	} else {
		// 如果 key 存在, 则将配置项转为指定类型
		var ok bool
		val, ok = v.(T)
		if !ok {
			err = fmt.Errorf("value by key \"%v\" not match type \"%v\"", key, reflect.ValueOf(def).Type().Name())
		}
	}
	return
}
