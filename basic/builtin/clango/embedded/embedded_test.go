package embedded_test

import (
	"study/basic/builtin/clango/embedded"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试调用内置 C 代码
//
// 测试 `clango.go` 文件中通过注释内嵌的 C 代码
func TestCLang_Embedded(t *testing.T) {
	ptr := embedded.CreateCString("Hello World!")
	defer embedded.FreeCString(ptr)

	s := embedded.ConvertCString(ptr)
	assert.Equal(t, s, "Hello World!")

	embedded.ShowCString(ptr)
}
