package static

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试将文件内容嵌入全局变量中
func TestEmbed_GlobalStatic(t *testing.T) {
	assert.Equal(t, "Hello World, This content was embed as static resource\n", STATIC_DATA)
}

// 测试将文件内容嵌入全局变量, 读取该变量并进行 json 反序列化
func TestEmbed_LocalStatic(t *testing.T) {
	// 读取反序列化结果
	user, err := GetEmbedStaticString()

	assert.Nil(t, err)
	assert.Equal(t, "Alvin", user.Name)
	assert.Equal(t, "M", user.Gender)
	assert.Equal(t, 42, user.Age)
}
