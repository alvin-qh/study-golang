package arrays_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试使用数组指针
//
// 通过 `&<数组变量>` 可以获取数组的指针, 通过 `*<数组指针>` 可以得到数组本身
//
// 另外, 和 C 语言不同, Go 语言不支持通过 `+/-` 运算操作数组指针, 仍是通过下标来访问指定数组元素
func TestPointer_PointerOfArray(t *testing.T) {
	// 创建一个长度为 3 的整型数组
	a := [...]int{1, 2, 3}

	// 获取数组的指针
	pa := &a

	// 确认通过数组指针可访问到数组本身
	assert.Equal(t, [...]int{1, 2, 3}, *pa)

	// 确认通过数组指针和数组下标, 也可以访问数组元素, 像操作数组一样
	assert.Equal(t, 1, pa[0])

	// 确认通过数组指针和下标, 可以修改数组元素
	pa[0] = 10
	assert.Equal(t, [...]int{10, 2, 3}, a)
}
