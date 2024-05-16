package pointer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 函数变量指针
//
// Go 语言的函数本身是引用类型, 所以不存在函数指针这一概念, 但如果一个函数复制给变量后, 则可以获取该变量的地址
func TestFunction_Pointer(t *testing.T) {
	// 定义函数变量
	f := Add

	// 定义函数变量指针类型
	type FuncPtr = *func(a, b int) int

	// 函数变量地址赋值给函数变量
	var pf FuncPtr = &f

	r := (*pf)(10, 20)
	assert.Equal(t, 30, r)
}
