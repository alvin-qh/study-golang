package arrays_test

import (
	"study/basic/builtin/arrays"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试创建数组
//
// 通过 `[<n>]<type>` 或 `[...]<type>{<v1>, <v2>, ...}` 声明一个数组, 对于 Go 语言来说, 数组必须具备确定的长度,
// 所以要么给数组一个长度值, 要么初始化时指定数组的元素
//
// 对于数组, 可以通过 `len` 函数来获取其长度
//
// Go 语言的数组元素均为值类型, 故数组元素会自动初始化为其类型的默认值
func TestArray_Create(t *testing.T) {
	// 创建长度为 5 的整型类型数组
	var arr [5]int

	// 获取数组长度
	cnt := len(arr)

	// 确认数组长度为 5
	assert.Equal(t, 5, cnt)

	// 确认数组元素的初始值均为 0
	assert.Equal(t, [...]int{0, 0, 0, 0, 0}, arr)
}

// 测试通过循环遍历数组
//
// 可以通过 `for` 循环产生数组下标, 并通过下标访问数组元素
func TestArray_Index(t *testing.T) {
	// 声明长度为 5 的数组
	var arr [5]int

	// 通过下标遍历数组
	for i := 0; i < len(arr); i++ { // lint:ignore
		// 使用下标访问数组元素
		arr[i] = i + 1
	}

	// 数组实际的结果值
	assert.Equal(t, [...]int{1, 2, 3, 4, 5}, arr)
}

// 测试数组的 `range` 操作
//
// Go 语言支持通过 `for n, v := range <array>` 的语法对数组进行迭代, 其中:
// - `n` 为迭代过程中每个元素的下标,
// - `v` 表示迭代过程中每个元素的值
func TestArray_ForInRange(t *testing.T) {
	// [...] 表示数组实际长度由后续初始化列表中的元素数量决定
	arr := [...]int{1, 2, 3, 4, 5}

	// 遍历数组, i 为当前遍历数组项的下标, v 为值
	for i, v := range arr {
		// 确认数组每一项的值正确
		assert.Equal(t, i+1, v)
	}
}

// 测试数组赋值
//
// 如果声明数组变量时, 同时指定了数组长度和初始化元素列表, 则初始化列表中的元素必须小于等于指定的数组长度
//
// 如果元素初始化列表中元素的数量小于指定的数组长度, 则其余未指定初始化值的元素会自动取默认值
func TestArray_Assign(t *testing.T) {
	// 初始化数组的前 1 个元素, 后续元素会自动初始化为其元素类型默认值 (本例中为 0)
	arr := [3]int{1}

	// 确认数组长度为 3
	assert.Len(t, arr, 3)

	// 确认数组元素值, 除了显式初始化的元素外, 其余元素值为 0
	assert.Equal(t, [...]int{1, 0, 0}, arr)
}

// 测试多维数组
//
// Go 语言的多维数组和大多数语言类似, 需要同时指定各个维度的长度 (或初始化元素列表)
func TestArray_MultiDim(t *testing.T) {
	// 声明一个 9 x 9 的整型数组
	arr := [9][9]int{}

	// 确认数组的第 1 维长度, 长度为 9
	assert.Len(t, arr, 9)

	// 确认数组的第 2 维长度, 长度为 9
	assert.Len(t, arr[0], 9)

	// 通过循环给数组赋值
	for i := range len(arr) {
		for j := range len(arr[i]) {
			// 为数组的第 i 行第 j 列赋值
			arr[i][j] = (i + 1) * (j + 1)
		}
	}

	// 确认数组元素每一项的值正确
	assert.Equal(t, [...][9]int{
		{1, 2, 3, 4, 5, 6, 7, 8, 9},
		{2, 4, 6, 8, 10, 12, 14, 16, 18},
		{3, 6, 9, 12, 15, 18, 21, 24, 27},
		{4, 8, 12, 16, 20, 24, 28, 32, 36},
		{5, 10, 15, 20, 25, 30, 35, 40, 45},
		{6, 12, 18, 24, 30, 36, 42, 48, 54},
		{7, 14, 21, 28, 35, 42, 49, 56, 63},
		{8, 16, 24, 32, 40, 48, 56, 64, 72},
		{9, 18, 27, 36, 45, 54, 63, 72, 81},
	}, arr)
}

// 测试任意类型数组项
//
// 如果数组元素类型为 `any`, 则意味着这个数组可以存储各类元素值
func TestArray_AnyType(t *testing.T) {
	// 声明一个项类型为 any 类型的数组, 即数组项可以为任意类型
	arr := [...]any{"Hello", 1, false}

	// 确认数组长度为 3
	assert.Len(t, arr, 3)

	// 获取数组每一项类型名称
	// 注意: 因为 arrays.GetTypeNameOfArrayElement() 函数的参数是一个切片, 所以这里需要使用 [:] 来获取数组的切片
	types := arrays.GetTypeNameOfArrayElement(arr[:])

	// 确认数组每一项的类型
	assert.Equal(t, []string{"string", "int", "bool"}, types)
}

// 测试数组复制
//
// Go 语言的赋值运算符 `=` 对于数组变量, 相当于"复制", 但是 Go 语言的数组复制是 "Copy on Write" 模式的,
// 即赋值后, 新数组变量和原数组变量指向同一个数组, 但修改其中任意一个后, 就会产生新数组, 以保证一个数组变量的修改不会影响到另一个
func TestArray_Clone(t *testing.T) {
	// 创建一个长度为 3 的整型数组
	arr := [...]int{1, 2, 3}

	// 复制数组
	arrDup := arr

	// 确认数组复制前后的值相同, 且数组的指针也相同
	// 因为此时尚未对复制后的数组进行写操作, 故 arr 和 arrDup 实际的内存地址相同
	assert.Equal(t, arr, arrDup)
	assert.Equal(t, &arr, &arrDup)

	// 对复制后的数组进行写操作
	arrDup[0] = 100

	// 确认数组复制前后的值不同, 且数组的指针也不同
	// 即在发生写操作时, 会真正触发数组在内存中的复制行为
	assert.NotEqual(t, arr, arrDup)
	assert.NotEqual(t, &arr, &arrDup)
}
