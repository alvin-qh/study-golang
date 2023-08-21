package conf

import "github.com/spf13/viper"

func SetConfigPaths(configPaths []string) {
	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}
}

func ReadJsonConfig(filename string) error {
	viper.SetConfigName(filename)
	viper.SetConfigType("json")
	return viper.ReadInConfig()
}
