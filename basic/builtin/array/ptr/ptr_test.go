package ptr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试数组指针
//
// 通过 `&<数组变量>` 可以获取数组的指针, 通过 `*<数组指针>` 可以得到数组本身
//
// 另外, 和 C 语言不同, Go 语言不支持通过 `+/-` 运算操作数组指针, 仍是通过下标来访问指定数组元素
func TestArray_Pointer(t *testing.T) {
	a := [...]int{1, 2, 3}

	// 获取数组的指针
	pa := &a
	assert.Equal(t, [...]int{1, 2, 3}, *pa)

	// 要通过数组指针访问数组元素, 仍是使用数组下标
	assert.Equal(t, 1, pa[0])

	// 通过指针改变数组的元素值, 数组指针的使用方式和数组本身基本一致
	pa[0] = 10
	assert.Equal(t, [...]int{10, 2, 3}, a)
}
