package dotenv

import "github.com/joho/godotenv"

// 模块初始化函数
func init() {
	// 初始化 dotenv 环境变量
	// 该函数将环境变量文件中的内容进行加载, 之后就可以通过 `os.Getenv()` 函数读取
	//
	// 可以同时加载多个环境变量文件
	godotenv.Load("../.env")
}
