package conf

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"

	"study/thirdpart/viper/logging"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

// 在所有测试前执行
func TestMain(m *testing.M) {
	// 初始化日志
	logging.Setup()

	os.Exit(m.Run())
}

// 测试读取 json 配置文件
func TestReadJsonConfig(t *testing.T) {
	// 设置读取当前路径下的 json 配置文件
	viper.SetConfigType("json")
	viper.SetConfigName("conf.json")
	viper.AddConfigPath(".")

	// 也可以直接设置配置文件的路径和名称, 取代上面三行代码
	// viper.SetConfigFile("./conf.json")

	// 读取 conf.json 配置文件
	err := viper.ReadInConfig()
	assert.NoError(t, err)

	// 确认配置文件读取正确
	assert.Equal(t, "localhost", viper.GetString("host.address"))
	assert.Equal(t, []int{5799, 6029}, viper.GetIntSlice("host.ports"))
	assert.Equal(t, "127.0.0.1", viper.GetString("database.metric.host"))
	assert.Equal(t, 3099, viper.GetInt("database.metric.port"))
}

// 用于反序列化 `conf.yml` 配置文件的结构体
type CacheConfig struct {
	Cache struct {
		Redis struct {
			MaxItems int `mapstructure:"max-items"`
			ItemSize int `mapstructure:"item-size"`
		} `mapstructure:"redis"`

		MemCache struct {
			MaxItems int `mapstructure:"max-items"`
			ItemSize int `mapstructure:"item-size"`
		} `mapstructure:"memcache"`
	} `mapstructure:"cache"`
}

// 测试对配置文件进行反序列化
func TestUnmarshalConfig(t *testing.T) {
	// 设置读取当前路径下的 yaml 配置文件
	viper.SetConfigType("yaml")
	viper.SetConfigName("conf.yml")
	viper.AddConfigPath(".")

	// 也可以直接设置配置文件的路径和名称, 取代上面三行代码
	// viper.SetConfigFile("./conf.yml")

	// 读取 conf.yml 配置文件
	err := viper.ReadInConfig()
	assert.NoError(t, err)

	conf := CacheConfig{}
	err = viper.Unmarshal(&conf)
	assert.NoError(t, err)

	assert.Equal(t, 100, conf.Cache.Redis.MaxItems)
	assert.Equal(t, 64, conf.Cache.Redis.ItemSize)
	assert.Equal(t, 200, conf.Cache.MemCache.MaxItems)
	assert.Equal(t, 80, conf.Cache.MemCache.ItemSize)
}

// 测试使用独立配置对象
func TestConfObject(t *testing.T) {
	// 实例化 viper 对象
	conf := viper.NewWithOptions(
		viper.KeyDelimiter("."), // 设置配置属性名层级的分隔符, 默认为 `.`
	)

	// 设置读取当前路径下的 ini 配置文件
	conf.SetConfigName("conf.toml")
	conf.SetConfigType("toml")
	conf.AddConfigPath(".")

	// 也可以直接设置配置文件的路径和名称, 取代上面三行代码
	// conf.SetConfigFile("./conf.ini")

	err := conf.ReadInConfig()
	assert.NoError(t, err)

	assert.Equal(t, "localhost", conf.GetString("host.address"))
	assert.Equal(t, []int{5799, 6029}, conf.GetIntSlice("host.ports"))

	assert.Equal(t, "198.0.0.1", conf.GetString("database.warehouse_host"))
}

// 测试保存配置文件
func TestSaveConfigFile(t *testing.T) {
	// 实例化 viper 对象
	conf := viper.New()

	// 设置配置信息
	conf.Set("logging.level", "DEBUG")
	conf.Set("logging.layout", "[%t{%Y-%M-%dT%h%m%s}] [%level] {%msg}%n")
	conf.Set("logging.appender", map[string]any{
		"file": map[string]any{
			"path":     "./file.log",
			"rolling":  true,
			"max-size": "100MB",
		},
		"console": map[string]any{
			"output": "system.out",
		},
	})

	// 设置配置对应的文件名
	conf.SetConfigFile("./new_conf.yml")
	// 存储配置信息到文件中
	conf.WriteConfig()
	defer os.Remove("./new_conf.yml")

	// 读取文件函数
	readAll := func(f *os.File) string {
		buf := make([]byte, 1024)
		res := bytes.Buffer{}
		for {
			n, err := f.Read(buf)
			if err != nil && !errors.Is(err, io.EOF) {
				panic(err)
			}
			if n == 0 {
				break
			}
			res.Write(buf[:n])
		}
		return res.String()
	}

	// 打开生成的配置文件
	f, err := os.Open("./new_conf.yml")
	// 确认文件存在
	assert.NoError(t, err)

	defer f.Close()

	// 读取生成的配置文件
	s := readAll(f)
	// 确认文件内容
	assert.Contains(t, s, "output: system.out")

	// 将配置信息存储到另一个文件中
	conf.WriteConfigAs("./new_conf.json")
	defer os.Remove("./new_conf.json")

	// 缺省生成的配置文件存在
	_, err = os.Stat("./new_conf.json")
	assert.NoError(t, err)
}
