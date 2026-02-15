package arrays_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试创建数组
//
// 通过 `[<n>]<type>` 或 `[...]<type>{<v1>, <v2>, ...}` 声明一个数组, 对于 Go 语言来说, 数组必须具备确定的长度,
// 所以要么给数组一个长度值, 要么初始化时指定数组的元素
//
// 对于数组, 可以通过 `len` 函数来获取其长度
func TestArray_Create(t *testing.T) {
	var arr [5]int  // 创建长度为 5 的数组
	cnt := len(arr) // 获取数组长度

	assert.Equal(t, 5, cnt)                       // len 函数用于测量数组的长度
	assert.Equal(t, [...]int{0, 0, 0, 0, 0}, arr) // 数组元素的初始值均为 0
}

// 测试通过循环遍历数组
//
// 和大多数语言类似, 可以通过 `for` 循环, 并通过下标来遍历数组的每个元素
func TestArray_Index(t *testing.T) {
	var arr [5]int // 声明长度为 5 的数组

	// 遍历数组
	for i := 0; i < len(arr); i++ {
		arr[i] = i + 1 // 使用下标访问数组元素
	}
	assert.Equal(t, [...]int{1, 2, 3, 4, 5}, arr) // 数组实际的结果值
}

// 测试数组的 `range` 操作
//
// Go 语言支持通过 `for n, v := range <array>` 的语法对数组进行迭代, 其中 `n` 为迭代过程中每个元素的下标,
// `v` 表示迭代过程中每个元素的值
func TestArray_ForInRange(t *testing.T) {
	arr := [...]int{1, 2, 3, 4, 5} // [...] 表示数组实际长度

	// 遍历数组
	for i, v := range arr { // i 为当前遍历数组项的下标, v 为值
		assert.Equal(t, i+1, v) // 查看遍历的结果值
	}
}

// 测试数组赋值
//
// 如果声明数组变量时, 同时指定了数组长度和初始化元素列表, 则初始化列表中的元素必须小于等于指定的数组长度
//
// 如果元素初始化列表中元素的数量小于指定的数组长度, 则其余未指定初始化值的元素会自动取默认值
func TestArray_Assign(t *testing.T) {
	arr := [3]int{1} // 初始化数组的前 1 个元素
	assert.Len(t, arr, 3)
	assert.Equal(t, [...]int{1, 0, 0}, arr) // 除了显式初始化的元素外, 其余元素值为 0
}

// 测试多维数组
//
// Go 语言的多维数组和大多数语言类似, 需要同时指定各个维度的长度 (或初始化元素列表)
func TestArray_MultiDim(t *testing.T) {
	arr := [9][9]int{}
	assert.Len(t, arr, 9)    // 数组的第 1 维长度
	assert.Len(t, arr[0], 9) // 数组的第 2 维长度

	// 通过循环给数组赋值
	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr[i]); j++ {
			arr[i][j] = (i + 1) * (j + 1)
		}

		// 比较每一维数组值
		assert.Equal(t, [...]int{
			(i + 1) * 1,
			(i + 1) * 2,
			(i + 1) * 3,
			(i + 1) * 4,
			(i + 1) * 5,
			(i + 1) * 6,
			(i + 1) * 7,
			(i + 1) * 8,
			(i + 1) * 9,
		}, arr[i])
	}
}

// 测试任意类型数组项
//
// 如果数组元素类型为 `interface{}`, 则意味着这个数组可以存储各类元素值
func TestArray_AnyType(t *testing.T) {
	arr := [...]any{"Hello", 1, false} // 声明一个项类型为 interface{} 类型的数组, 即数组项可以为任意类型

	assert.Len(t, arr, 3)
	assert.Equal(t, "string", reflect.TypeOf(arr[0]).Name())
	assert.Equal(t, "int", reflect.TypeOf(arr[1]).Name())
	assert.Equal(t, "bool", reflect.TypeOf(arr[2]).Name())
}

// 测试数组复制
//
// Go 语言的赋值运算符 `=` 对于数组变量, 相当于"复制", 但是 Go 语言的数组复制是 "Copy on Write" 模式的,
// 即赋值后, 新数组变量和原数组变量指向同一个数组, 但修改其中任意一个后, 就会产生新数组, 以保证一个数组变量的修改不会影响到另一个
func TestArray_Clone(t *testing.T) {
	arr := [...]int{1, 2, 3}

	arrDup := arr // 复制数组
	assert.Equal(t, arr, arrDup)
	assert.Equal(t, &arr, &arrDup)

	arrDup[0] = 100
	assert.NotEqual(t, arr, arrDup)   // 复制前后数组不同
	assert.NotEqual(t, &arr, &arrDup) // 复制前后数组不同
}
