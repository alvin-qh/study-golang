package conf

import (
	"os"
	"path"
	"testing"
	"viper/logging"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

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

// 在所有测试前执行
func TestMain(m *testing.M) {
	// 初始化日志
	logging.Setup()

	// 初始化配置文件路径

	// 获取当前工作路径
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	// 获取演示配置文件路径
	confPath := path.Join(cwd, "../demo")
	// 设置配置文件路径
	SetConfigPaths([]string{confPath})

	os.Exit(m.Run())
}

// 测试读取 json 配置文件
func TestReadJsonConfig(t *testing.T) {
	// 读取 conf.json 配置文件
	err := ReadJsonConfig("conf.json")
	assert.NoError(t, err)

	// 确认配置文件读取正确
	assert.Equal(t, "localhost", viper.GetString("host.address"))
	assert.Equal(t, []int{5799, 6029}, viper.GetIntSlice("host.ports"))
	assert.Equal(t, "127.0.0.1", viper.GetString("database.metric.host"))
	assert.Equal(t, 3099, viper.GetInt("database.metric.port"))
}

func TestUnmarshalConfig(t *testing.T) {
	// 读取 conf.yml 配置文件
	err := ReadYamlConfig("conf.yml")
	assert.NoError(t, err)

	conf := CacheConfig{}
	err = viper.Unmarshal(&conf)
	assert.NoError(t, err)

	assert.Equal(t, 100, conf.Cache.Redis.MaxItems)
	assert.Equal(t, 64, conf.Cache.Redis.ItemSize)
	assert.Equal(t, 200, conf.Cache.MemCache.MaxItems)
	assert.Equal(t, 80, conf.Cache.MemCache.ItemSize)
}
