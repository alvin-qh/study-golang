package main

import (
	"path"
	"study/thirdpart/viper/logging"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"os"
)

const (
	ConfPath = "demo"
	ConfFile = "conf.json"
)

func main() {
	// 初始化日志
	logging.Setup()

	// 获取当前工作目录
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Cannot get current path", cwd)
	}

	// 生成配置文件路径
	confPath := path.Join(cwd, ConfPath)
	log.Infof("Get path of config files \"%v\"", confPath)

	// 设置配置文件路径
	viper.AddConfigPath("conf")
	// 设置配置文件类型
	viper.SetConfigType("json")
	// 设置配置文件名称
	viper.SetConfigName("conf.json")

	// 读取配置文件
	if err = viper.ReadInConfig(); err != nil {
		log.Fatal("Cannot read config file", err)
	}

	// 输出配置文件内容
	log.Infof("host.address=%v", viper.GetString("host.address"))
	log.Infof("host.ports=%v", viper.GetIntSlice("host.ports"))
	log.Infof("database.metric.host=%v", viper.GetString("database.metric.host"))
	log.Infof("database.metric.port=%v", viper.GetInt("database.metric.port"))
}
