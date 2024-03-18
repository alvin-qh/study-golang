package clango

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试调用内置 C 代码
//
// 测试 `clango.go` 文件中通过注释内嵌的 C 代码
func TestCallEmbeddedC(t *testing.T) {
	ptr := CreateCString("Hello World!")
	defer FreeCString(ptr)

	s := ConvertCString(ptr)
	assert.Equal(t, s, "Hello World!")

	ShowCString(ptr)
}

// 测试通过 `.h` 文件引入的外部 C 代码
//
// 测试 `clango.go` 文件中通过 C 头文件引入的外部 C 代码
func TestCallExternalC(t *testing.T) {
	pt1 := CreatePoint(10.0, 20.0)
	assert.Equal(t, 10.0, pt1.GetX())
	assert.Equal(t, 20.0, pt1.GetY())

	pt2 := CreatePoint(30.0, 50.0)

	dis := pt1.Distance(pt2)
	assert.Equal(t, 36.06, dis)
}
