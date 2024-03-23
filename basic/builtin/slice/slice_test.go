package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 创建切片变量
func TestCreateSlice(t *testing.T) {
	// 创建切片类型变量, 此时切片为空
	var s []int
	assert.Nil(t, s)
	assert.Equal(t, 0, len(s)) // nil 的长度为 0

	// 创建一个空的切片
	s = []int{}
	assert.Len(t, s, 0) // 为 nil 的切片长度为 0

	// 向切片中添加元素
	s = append(s, 0)
	s = append(s, 1, 2, 3)
	assert.EqualValues(t, []int{0, 1, 2, 3}, s)

	// 创建长度为 5 的切片并初始化元素
	s = []int{1, 2, 3, 4, 5}
	assert.EqualValues(t, []int{1, 2, 3, 4, 5}, s)
}

// 使用 `make` 函数创建切片
func TestMakeSlice(t *testing.T) {
	// 通过 make 函数初始化切片, 初始长度 3
	s := make([]int, 3)
	assert.EqualValues(t, []int{0, 0, 0}, s) // 切片初始长度为 3

	// 通过下标给切片赋值
	s[0], s[1], s[2] = 100, 200, 300
	assert.EqualValues(t, []int{100, 200, 300}, s)
}

// 向切片中添加数据
func TestAppendToSlice(t *testing.T) {
	// 创建一个 len=0, cap=10 的切片
	s := make([]int, 0, 10)

	// 向切片中添加元素
	s = append(s, 100)
	assert.EqualValues(t, []int{100}, s)

	// 向切片中添加多个元素
	s = append(s, 200, 300, 400)
	assert.EqualValues(t, []int{100, 200, 300, 400}, s)

	assert.Equal(t, 4, len(s))
	assert.Equal(t, 10, cap(s))
}

// 切片部分截取
func TestCutSubPiecesFromSlice(t *testing.T) {
	arr := []int{1, 2, 3, 4, 5}

	// 切片为数组的前 2 个元素
	s := arr[:2]
	assert.EqualValues(t, []int{1, 2}, s)

	// 切片为数组的后 3 个元素
	s = arr[2:]
	assert.EqualValues(t, []int{3, 4, 5}, s)

	// 切片为数组的第 3 个元素
	s = arr[2:3]
	assert.EqualValues(t, []int{3}, s)
}

// 从切片中"删除"元素
//
// 切片本身不具备删除元素的操作, 可以通过新建切片并忽略要删除元素的方式进行
func TestRemoveItemFromSlice(t *testing.T) {
	s := []int{1, 2, 3, 4, 5}

	// 先取前 2 个元素的切片, 在其之上添加后 2 个元素, 相当于删除第 3 个元素
	// s[3:]... 相当于将数组展开成若干参数项
	s = append(s[:2], s[3:]...)

	assert.EqualValues(t, []int{1, 2, 4, 5}, s)
}

// 切片的引用特性
//
// 和数组不同, 切片变量的特性是引用, 所以赋值操作只能赋值切片的引用, 而不会产生新的切片
func TestReferenceVariableOfSlice(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2 := s1

	assert.Equal(t, s1, s2) // 赋值运算符会传递切片的引用

	// 两个引用指向了同一个切片
	s2[1] = 20
	assert.Equal(t, s1, s2)
	assert.EqualValues(t, []int{1, 20, 3}, s1)
	assert.EqualValues(t, []int{1, 20, 3}, s2)
}

// 复制切片
// 通过 copy 函数, 可以将一个切片的元素复制到另一个切片中
func TestCopySlice(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2 := make([]int, len(s1))

	// 将 s1 的元素复制到 s2 中
	copy(s2, s1)
	assert.EqualValues(t, []int{1, 2, 3}, s2) // copy 会复制切片的内容

	// 因为 s1 长度较小, 所以会复制 s1 的全部元素, 并保留 s2 的多余元素
	s2 = make([]int, 4)
	copy(s2, s1)
	assert.EqualValues(t, []int{1, 2, 3, 0}, s2)

	// 因为 s2 长度较小, 所以会复制 s1 中和 s2 长度匹配的那部分, 其余的不复制
	s2 = make([]int, 2)
	copy(s2, s1)
	assert.EqualValues(t, []int{1, 2}, s2)
}

// 测试多维切片
func TestSlice(t *testing.T) {
	// 创建 2 维切片
	s := make([][]int, 0)

	s = append(s, []int{1, 2, 3})
	assert.EqualValues(t, [][]int{{1, 2, 3}}, s)

	s = append(s, []int{4, 5, 6}, []int{7, 8, 9})
	assert.EqualValues(t, [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}, s)
}

// 测试切片 cap 的增长率
func TestSliceCapGrowUp(t *testing.T) {
	var s []int
	assert.Equal(t, len(s), cap(s), 0)

	c := 1
	for i := 0; i < 20; i++ {
		// 当切片不为空时, 每当 cap 和 len 相等, 在添加元素时, 会分配原有长度 2 倍的空间作为新的 cap
		if len(s) > 0 && len(s) == cap(s) {
			c = len(s) * 2
		}

		s = append(s, i)
		assert.Equal(t, len(s), i+1)
		assert.Equal(t, cap(s), c)
	}
}
