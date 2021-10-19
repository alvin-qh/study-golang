// GO 集合操作
// GO 集合包括：数组、切片和列表

package builtin

import (
	"container/list"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试数组
func TestArray(t *testing.T) {
	// 创建长度为 10 的数组
	var array1 [5]int
	length := len(array1)                                    // 获取数组长度
	assert.Equal(t, 5, length)                               // len 函数用于测量数组的长度
	assert.ElementsMatch(t, [...]int{0, 0, 0, 0, 0}, array1) // 数组元素的初始值均为 0

	// 遍历数组
	for i := 0; i < len(array1); i++ {
		array1[i] = i + 1 // 使用下标访问数组元素
	}
	assert.ElementsMatch(t, [...]int{1, 2, 3, 4, 5}, array1) // 数组实际的结果值

	// 通过 range 方法遍历数组，其中 index 为数组下标，value 为对应的元素值
	for index, value := range array1 {
		assert.Equal(t, index+1, value) // 查看遍历的结果值
	}

	// 可以通过赋值的方式创建数组
	array2 := [3]int{1} // 初始化数组的前 1 个元素
	assert.Len(t, array2, 3)
	assert.ElementsMatch(t, [...]int{1, 0, 0}, array2) // 除了显式初始化的元素外，其余元素值为 0

	// 隐式声明数组长度，初始化全部数组元素
	// '...' 表示数组长度由初始化元素个数决定
	// 注意：如果省略了省略号，则表示一个 切片 而非 数组
	array3 := [...]int{1, 2, 3}
	assert.Len(t, array3, 3)
	assert.ElementsMatch(t, [...]int{1, 2, 3}, array3)

	// 初始化多维数组
	array4 := [9][9]int{}
	assert.Len(t, array4, 9)    // 数组的第 1 维长度
	assert.Len(t, array4[0], 9) // 数组的第 2 维长度

	// 通过循环给数组赋值
	for i := 0; i < len(array4); i++ {
		for j := 0; j < len(array4[i]); j++ {
			array4[i][j] = (i + 1) * (j + 1)
		}

		assert.ElementsMatch(t, [...]int{
			(i + 1) * 1,
			(i + 1) * 2,
			(i + 1) * 3,
			(i + 1) * 4,
			(i + 1) * 5,
			(i + 1) * 6,
			(i + 1) * 7,
			(i + 1) * 8,
			(i + 1) * 9, // 注意，这里要多一个 逗号，表示参数换行
		}, array4[i])
	}

	// 定义一个元素为任意类型的数组
	array5 := [...]Any{"Hello", 1, false}
	assert.Len(t, array5, 3)
	assert.ElementsMatch(t, [...]interface{}{"Hello", 1, false}, array5)

	// 定义一个元素为 interface{} 类型的数组，相当于 []Any{...}
	array6 := [...]interface{}{"Hello", 1, false}
	assert.Len(t, array6, 3)
	assert.ElementsMatch(t, [...]Any{"Hello", 1, false}, array6)

	// 测试数组的指针
	array7 := [...]int{1, 2, 3}
	pArray7 := &array7                                   // var pArray7 *[3]Any = &array7，获取数组的指针
	assert.ElementsMatch(t, [...]int{1, 2, 3}, *pArray7) // 解引指针，获取数组

	pArray7[0] = 10 // 通过指针改变数组的元素值
	assert.ElementsMatch(t, [...]int{10, 2, 3}, array7)

	// 表示数组的拷贝
	array8 := array7
	array8[0] = 100
	assert.False(t, array7[0] == array8[0]) // array7 和 array8 是两个不同的数组
}

// 测试切片
// 切片表示一个 可变长度 的数组的 引用类型
// 切片有两个长度 len, cap，前者表示切片的元素个数，后者表示切片的实际长度，超过 cap 后，切片需要重新分配
func TestSlice(t *testing.T) {
	// 创建切片类型 空 变量
	var slice []int // slice := []int(nil)，创建一个切片类型变量
	assert.Nil(t, slice)
	assert.Equal(t, 0, len(slice)) // nil 的长度为 0

	// 创建一个空的切片
	slice = []int{}
	assert.Len(t, slice, 0)

	slice = append(slice, 0)
	assert.Len(t, slice, 1)        // 向切片中添加元素
	slice = append(slice, 1, 2, 3) //向切片中添加多个元素
	assert.Len(t, slice, 4)
	assert.ElementsMatch(t, slice, [...]int{0, 1, 2, 3})

	// 通过 make 函数初始化切片，初始长度 3，最大长度 5
	slice = make([]int, 3)
	assert.Equal(t, 3, len(slice)) // 数组初始长度为 3
	assert.ElementsMatch(t, []int{0, 0, 0}, slice)
	assert.Equal(t, 3, cap(slice)) // 数组初始容积为 5，即

	slice[0], slice[1], slice[2] = 100, 200, 300 // 通过下标给切片赋值
	assert.ElementsMatch(t, []int{100, 200, 300}, slice)

	slice = append(slice, 400)     // 向切片中添加元素
	assert.Equal(t, len(slice), 4) // 切片的长度变为 4
	assert.Equal(t, cap(slice), 6) // 数组 cap 增长为 6
	assert.ElementsMatch(t, slice, []int{100, 200, 300, 400})

	slice = append(slice, 500, 600, 700) // 向切片中添加多个元素
	assert.Equal(t, len(slice), 7)       // 切片的长度变为 7
	assert.Equal(t, cap(slice), 12)      // 数组 cap 增长为 12

	array := [...]int{1, 2, 3, 4, 5}

	// 从数组（或另一个切片）中获取切片
	slice = array[:2]
	assert.Len(t, slice, 2)

	// ints := IntArray(array[:])
	// assert.Equal(t, len(ints), len(array))

	// ints = IntArray(array[0:])
	// assert.Equal(t, len(ints), len(array))

	// ints = IntArray(array[1:2])
	// assert.Equal(t, len(ints), 1)
	// assert.Equal(t, ints[0], 2)

	// ints = IntArray(array[1:3])
	// assert.Equal(t, len(ints), 2)
	// assert.Equal(t, ints[0], 2)
	// assert.Equal(t, ints[1], 3)

	// ints = IntArray(array[:3])
	// assert.Equal(t, len(ints), 3)
	// assert.Equal(t, ints[0], 1)
	// assert.Equal(t, ints[1], 2)
	// assert.Equal(t, ints[2], 3)
}

func TestSliceGrowUp(t *testing.T) {
	var ints []int
	assert.Equal(t, len(ints), 0)
	assert.Equal(t, cap(ints), 0)

	c := 1
	for i := 0; i < 20; i++ {
		if len(ints) > 0 && len(ints) == cap(ints) {
			c = len(ints) * 2
		}

		ints = append(ints, i)
		assert.Equal(t, len(ints), i+1)
		assert.Equal(t, cap(ints), c)
	}
}

func TestSliceShare(t *testing.T) {
	a := []int{1, 2, 3}
	b := a
	assert.Equal(t, &a, &b)

	a[1] = 100
	assert.Equal(t, &a, &b)
	assert.Equal(t, b[1], 100)

	a = append(a, 200)
	assert.NotEqual(t, &a, &b)
}

func TestNewInts(t *testing.T) {
	ints := NewInts(10, 100)
	assert.Equal(t, len(ints), 10)
	assert.Equal(t, cap(ints), 100)
}

func TestInts_Append(t *testing.T) {
	ints := IntArray{1, 2, 3, 4}
	ints.Append(5)
	assert.Equal(t, ints[4], 5)
	assert.Equal(t, len(ints), 5)
}

func TestInts_Remove(t *testing.T) {
	ints := IntArray{1, 2, 3, 4}
	ints.Remove(2)
	assert.Equal(t, ints[2], 4)
	assert.Equal(t, len(ints), 3)

	ints.Remove(0)
	assert.Equal(t, ints[1], 4)
	assert.Equal(t, len(ints), 2)

	ints.Remove(len(ints) - 1)
	assert.Equal(t, ints[0], 2)
	assert.Equal(t, len(ints), 1)
}

func TestInts_Clear(t *testing.T) {
	ints := IntArray{1, 2, 3, 4}
	ints.Clear()

	assert.Equal(t, len(ints), 0)
}

func TestInts_Size(t *testing.T) {
	ints := IntArray{1, 2, 3, 4}
	assert.Equal(t, ints.Size(), 4)

	ints.Clear()
	assert.Equal(t, ints.Size(), 0)
}

func TestList(t *testing.T) {
	lst := list.New()
	lst.PushBack(1)
	assert.Equal(t, lst.Len(), 1)

	lst.PushBack("Hello")
	assert.Equal(t, lst.Len(), 2)

	type Any interface{}

	array := make([]Any, lst.Len())
	for iter, i := lst.Front(), 0; iter != nil; iter, i = iter.Next(), i+1 {
		array[i] = iter.Value
	}
	assert.Equal(t, len(array), 2)
}

func TestListAt(t *testing.T) {
	lst := ListAssign(1, 2, 3, "Hello")
	assert.Equal(t, ListAt(lst, 0), 1)
	assert.Equal(t, ListAt(lst, 1), 2)
	assert.Equal(t, ListAt(lst, 3), "Hello")

	fn := func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}

	defer fn()
	ListAt(lst, 4)
}

func TestListToSlice(t *testing.T) {
	lst := ListAssign(1, 2, 3, "Hello")
	array := ListToSlice(lst)
	assert.Equal(t, len(array), 4)
	assert.ElementsMatch(t, array, []Any{1, 2, 3, "Hello"})
}
