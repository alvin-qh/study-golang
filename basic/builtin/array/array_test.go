package array

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// # 测试创建数组
func TestCreateArray(t *testing.T) {
	var arr [5]int  // 创建长度为 5 的数组
	cnt := len(arr) // 获取数组长度

	assert.Equal(t, 5, cnt)                               // len 函数用于测量数组的长度
	assert.ElementsMatch(t, [...]int{0, 0, 0, 0, 0}, arr) // 数组元素的初始值均为 0
}

// # 测试通过循环遍历数组
func TestLoopForArray(t *testing.T) {
	var arr [5]int // 声明长度为 5 的数组

	// 遍历数组
	for i := 0; i < len(arr); i++ {
		arr[i] = i + 1 // 使用下标访问数组元素
	}
	assert.ElementsMatch(t, [...]int{1, 2, 3, 4, 5}, arr) // 数组实际的结果值
}

// # 测试数组的 `range` 操作
func TestRangeForArray(t *testing.T) {
	arr := [...]int{1, 2, 3, 4, 5} // [...] 表示数组实际长度

	// 遍历数组
	for i, v := range arr { // i 为当前遍历数组项的下标, v 为值
		assert.Equal(t, i+1, v) // 查看遍历的结果值
	}
}

// # 测试数组赋值
func TestAssignArray(t *testing.T) {
	arr := [3]int{1} // 初始化数组的前 1 个元素
	assert.Len(t, arr, 3)
	assert.ElementsMatch(t, [...]int{1, 0, 0}, arr) // 除了显式初始化的元素外, 其余元素值为 0
}

// # 测试多维数组
func TestMultiDimensionalArray(t *testing.T) {
	arr := [9][9]int{}
	assert.Len(t, arr, 9)    // 数组的第 1 维长度
	assert.Len(t, arr[0], 9) // 数组的第 2 维长度

	// 通过循环给数组赋值
	for i := 0; i < len(arr); i++ {
		for j := 0; j < len(arr[i]); j++ {
			arr[i][j] = (i + 1) * (j + 1)
		}

		// 比较每一维数组值
		assert.ElementsMatch(t, [...]int{
			(i + 1) * 1,
			(i + 1) * 2,
			(i + 1) * 3,
			(i + 1) * 4,
			(i + 1) * 5,
			(i + 1) * 6,
			(i + 1) * 7,
			(i + 1) * 8,
			(i + 1) * 9, // 注意, 这里要多一个 逗号, 表示参数换行
		}, arr[i])
	}
}

// # 测试任意类型数组项
func TestGenericArrayItem(t *testing.T) {
	arr := [...]interface{}{"Hello", 1, false} // 声明一个项类型为 interface{} 类型的数组, 即数组项可以为任意类型

	assert.Len(t, arr, 3)
	assert.Equal(t, "string", reflect.TypeOf(arr[0]).Name())
	assert.Equal(t, "int", reflect.TypeOf(arr[1]).Name())
	assert.Equal(t, "bool", reflect.TypeOf(arr[2]).Name())
}

// # 测试数组指针
func TestPointerOfArray(t *testing.T) {
	arr := [...]int{1, 2, 3}
	parr := &arr                                      // 获取数组的指针
	assert.ElementsMatch(t, [...]int{1, 2, 3}, *parr) // 解引指针, 获取数组

	parr[0] = 10 // 通过指针改变数组的元素值, 数组指针的使用方式和数组本身基本一致
	assert.ElementsMatch(t, [...]int{10, 2, 3}, arr)
}

// 测试数组复制
func TestCloneArray(t *testing.T) {
	arr := [...]int{1, 2, 3}

	arrCpy := arr // 复制数组
	assert.Equal(t, arr, arrCpy)

	arrCpy[0] = 100
	assert.NotEqual(t, arr, arrCpy) // 复制前后数组不同
}
