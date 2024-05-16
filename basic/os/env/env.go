package env

import (
	"os"
	"strings"
)

// 获取当前系统环境变量
//
// `os.Environ` 函数返回 `[]string` 结果, 其中包含所有的系统环境变量,
// 每一条表示一个环境变量定义, 变量名称和值通过 `=` 连接, 例如: `PATH=/usr/bin:/bin`
//
// 返回包含环境变量名称和值的 Map
//
// 注意, 所有环境变量名称都为大写字母
func Environ() map[string]string {
	r := make(map[string]string)

	for _, e := range os.Environ() {
		kv := strings.SplitN(e, "=", 2)
		r[kv[0]] = kv[1]
	}
	return r
}
