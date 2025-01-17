package static

import (
	_ "embed"
	"encoding/json"
)

// 通过 `go:embed` 指令将当前目录下的文件在编译时嵌入全局变量
// `go:embed` 使用的约束包括: 1. 使用 `var` 修饰的变量; 2. 必须为全局变量
// 一旦文件嵌入到全局变量中, 则其和普通全局变量无区别, 直接使用即可
//
// 文件可以嵌入为 `string` 类型或 `[]byte` 类型
var (
	// 嵌入为字符串类型
	//go:embed static1.txt
	STATIC_DATA string

	// 嵌入为字节数组类型
	//go:embed static2.json
	STATIC_JSON []byte
)

// 用于解析 JSON 字符串的结构体
type User struct {
	Name   string `json:"name"`
	Gender string `json:"gender"`
	Age    int    `json:"age"`
}

// 获取嵌入的字符串 (JSON 格式) 转化为对象并返回
func GetEmbedStaticString() (*User, error) {
	var u User

	// 将嵌入字符串作为 JSON 格式进行反序列化
	if err := json.Unmarshal(STATIC_JSON, &u); err != nil {
		return nil, err
	}
	return &u, nil
}
