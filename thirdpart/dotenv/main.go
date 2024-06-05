package main

// 通过导入 "github.com/joho/godotenv/autoload" 包可以启用 dotenv 的懒加载功能,
// 当使用 `os.Getenv` 获取环境变量时, 会自动加载当前工作路径下的 `.env` 文件
import (
	"fmt"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	keyA := os.Getenv("KEY_A")
	keyB := os.Getenv("KEY_B")

	fmt.Printf("Key A is: %s, and key B is: %s\n", keyA, keyB)
}
