package conf

import "github.com/spf13/viper"

func SetConfigPaths(configPaths []string) {
	for _, path := range configPaths {
		viper.AddConfigPath(path)
	}
}

func readConfig(name string, ctype string) error {
	viper.SetConfigName(name)
	viper.SetConfigType(ctype)
	return viper.ReadInConfig()
}

func ReadJsonConfig(name string) error {
	return readConfig(name, "json")
}

func ReadYamlConfig(name string) error {
	return readConfig(name, "yaml")
}
