package conf

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试 `Default` 函数
//
// 确认获取配置项时, 如果配置项不存在, 则返回默认值
func TestDefault(t *testing.T) {
	v, err := Default("not.exist.key", 100)
	assert.NoError(t, err)
	assert.Equal(t, 100, v)

	v, err = Default("server.cors.max-age", 0)
	assert.NoError(t, err)
	assert.Equal(t, 86400, v)

	_, err = Default("server.cors.max-age", "Unknown")
	assert.EqualError(t, err, "value by key \"server.cors.max-age\" not match type \"string\"")
}
