package dotenv

import "github.com/joho/godotenv"

// 模块初始化函数
func init() {
	// 初始化 dotenv 环境变量
	// 该函数将环境变量文件中的内容进行加载, 之后就可以通过 `os.Getenv()` 函数读取
	//
	// 如果环境变量文件位于工作目录, 且文件名为 `.env`, 则可以省略参数;
	// 如果环境变量文件不在工作目录, 但文件名为 `.env`, 则参数中可省略文件名
	// 如果环境变量文件不在工作目录, 且文件名不为 `.env`, 则参数需要表明环境变量文件的完全路径名
	//
	// 可以同时加载多个环境变量文件
	//
	// godotenv.Load("../.env")
	godotenv.Load("../")
}
