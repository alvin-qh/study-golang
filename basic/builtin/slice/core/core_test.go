package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试定义切片类型变量
func TestSlice_Define(t *testing.T) {
	// 创建切片类型变量, 此时切片为 nil, 长度为 0
	var s []int
	assert.Nil(t, s)
	assert.Equal(t, 0, len(s))

	// 创建一个空的切片, 此时切片非 nil, 长度为 0
	s = []int{}
	assert.NotNil(t, s)
	assert.Equal(t, 0, len(s))

	// 创建切片并初始化元素
	s = []int{1, 2, 3, 4, 5}
	assert.NotNil(t, s)
	assert.Equal(t, 5, len(s))
	assert.Equal(t, []int{1, 2, 3, 4, 5}, s)
}

// 测试通过 `make` 函数创建切片
func TestSlice_Make(t *testing.T) {
	// 创建切片, 初始长度 3
	s := make([]int, 3)
	assert.Equal(t, 3, len(s))
	assert.Equal(t, 3, cap(s))
	assert.Equal(t, []int{0, 0, 0}, s)

	// 通过下标索引设置切片元素
	s[0] = 1
	s[1] = 2
	s[2] = 3
	assert.Equal(t, []int{1, 2, 3}, s)

	// 创建切片, 长度为 0, Capacity 为 10
	s = make([]int, 0, 10)
	assert.Equal(t, 0, len(s))
	assert.Equal(t, 10, cap(s))

	// 访问切片时, 下标越界会导致 Panic
	assert.PanicsWithError(t, "runtime error: index out of range [0] with length 0", func() {
		s[0] = 1
	})
}

// 切片部分截取, 形成新切片
//
// 新切片的起始下标为 `0`, 长度为截取的元素个数, Capacity 为原切片的 Capacity - 截取的元素起始索引, 例如:
//
//	s := []int{1, 2, 3, 4, 5}
//	sc := s[1:3]
//
// 假设 `s` 的 Capacity 为 `8`, 则 `sc` 为 [2, 3], `len(sc)` 为 `2`, Capacity 为 `8 - 1` 为 `7`
func TestSlice_Cut(t *testing.T) {
	// 定义长度为 0, Capacity 为 10 的切片
	s := make([]int, 0, 10)
	s = append(s, 1, 2, 3, 4, 5)

	// 截取前 2 个元素
	// 结果切片的长度为 `2`, Capacity 为 `10 - 0` 为 `0`
	sc := s[:2]
	assert.Equal(t, 2, len(sc))
	assert.Equal(t, 10, cap(sc))
	assert.Equal(t, []int{1, 2}, sc)

	// 截取后 3 个元素
	// 结果切片的长度为 `3`, Capacity 为 `10 - 2` 为 `7`
	sc = s[2:]
	assert.Equal(t, 3, len(sc))
	assert.Equal(t, 8, cap(sc))
	assert.Equal(t, []int{3, 4, 5}, sc)

	// 截取第 3 个元素
	// 结果切片的长度为 `1`, Capacity 为 `10 - 2` 为 `8`
	sc = s[2:3]
	assert.Equal(t, 1, len(sc))
	assert.Equal(t, 8, cap(sc))
	assert.Equal(t, []int{3}, sc)

	// 截取下标 2~4 的元素, 并通过第三个值指定 Capacity 的最大值为 `6`
	// 结果切片的长度为 `3`, Capacity 为 `6 - 2` 为 `4`
	sc = s[2:len(s):6]
	assert.Equal(t, 3, len(sc))
	assert.Equal(t, 4, cap(sc))
	assert.Equal(t, []int{3, 4, 5}, sc)

	// 获取一个同类型空数组, 即长度为 `0`, Capacity 为 `0` 的切片
	sc = s[:0:0]
	assert.Equal(t, []int{}, sc)
}

// 测试通过 `append` 函数追加元素
func TestSlice_Append(t *testing.T) {
	// 测试为 Nil 切片追加元素
	t.Run("append to nil slice", func(t *testing.T) {
		// 创建 Nil 切片
		var s []int
		assert.Nil(t, s)

		// Nil 切片的长度为 0
		assert.Equal(t, 0, len(s))
		// Nil 切片的 Capacity 为 0
		assert.Equal(t, 0, cap(s))

		// 向空切片追加元素
		s = append(s, 1)
		assert.Equal(t, 1, len(s))
		assert.Equal(t, 1, cap(s))
		assert.Equal(t, []int{1}, s)

		// 向空切片中继续添加元素
		// 当元素长度超过 Capability 时, 会自动扩大 Capability 长度
		s = append(s, 2, 3, 4, 5, 6, 7)
		assert.Equal(t, 7, len(s))
		assert.Equal(t, 8, cap(s))
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7}, s)
	})

	// 测试向空切片中追加元素
	t.Run("append to empty slice", func(t *testing.T) {
		s := make([]int, 0, 5)

		// 向空切片追加多个元素
		// 在元素数量小于 Capacity 时, 不会增加 Capacity 长度
		s = append(s, 1)
		assert.Len(t, s, 1)
		assert.Equal(t, 5, cap(s))
		assert.Equal(t, []int{1}, s)

		// 继续追加多个元素
		// 当元素长度超过 Capability 时, 会自动扩大 Capability 长度
		s = append(s, 2, 3, 4, 5, 6)
		assert.Len(t, s, 6)
		assert.Equal(t, 10, cap(s))
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, s)
	})

	// 测试通过 `append` 函数删除切片元素
	t.Run("delete slice element by append", func(t *testing.T) {
		s := []int{1, 2, 3, 4, 5}

		// 删除下标为 2 的元素, 相当于切片中将下标小于 2 的元素和下标大于 2 的元素合并在一起
		s = append(s[:2], s[3:]...)
		assert.Equal(t, []int{1, 2, 4, 5}, s)
	})
}

// 切片的引用特性
//
// 和数组不同, 切片变量的特性是引用, 所以赋值操作只能赋值切片的引用, 而不会产生新的切片
func TestSlice_Reference(t *testing.T) {
	s1 := []int{1, 2, 3}

	// 赋值运算符会传递切片的引用
	s2 := s1

	// 两个变量是不同变量
	assert.NotSame(t, &s1, &s2)

	// 两个切片的元素地址相同, 即元素位于相同内存空间
	assert.Same(t, &s1[0], &s2[0])

	// 改变其中一个切片变量的元素, 另一个切片变量也同步修改
	s2[1] = 20
	assert.Equal(t, []int{1, 20, 3}, s1)

	// 通过 append 函数追加元素, 此时会产生新切片
	// 返回新切片, 之后 s1 和 s2 不再引用同一个切片
	s2 = append(s2, 4)
	assert.NotSame(t, &s1[0], &s2[0])
	assert.NotEqual(t, s2, s1)
}

// 测试复制切片
//
// 通过 Go 语言提供的 `copy` 函数可以将一个切片的元素复制到另一个切片中
//
// 复制的元素数量为两个切片中的最小长度, 即如果将 `s1` 复制到 `s2`, 则:
//   - 如果 `len(s1)` >= `len(s2)` 则从 `s1` 向 `s2` 复制 `len(s2)` 个元素
//   - 如果 `len(s1)` <= `len(s2)` 则从 `s1` 向 `s2` 复制 `len(s1)` 个元素
func TestSlice_Copy(t *testing.T) {
	s1 := []int{1, 2, 3}

	// 创建一个和 s1 长度相同的数组
	s2 := make([]int, len(s1))

	// 将 s1 切片复制到 s2 切片中
	copy(s2, s1)
	assert.Equal(t, s2, s1)

	// 将 s1 切片下标范围 1~2 的元素 (即 [2, 3]) 复制到 s2 切片中
	// s2 从 [<1>, <2>, 3] 变为 [<2>, <3>, 3]
	copy(s2, s1[1:])
	assert.Equal(t, []int{2, 3, 3}, s2)

	// 将 s1 切片下标范围 1~2 的元素 (即 [2, 3]) 复制到 s2 下标 1 开始的切片中
	// s2 从 [2, <3>, <3>] 变为 [2, <2>, <3>]
	copy(s2[1:], s1[1:])
	assert.Equal(t, []int{2, 2, 3}, s2)

	// 将 s1 切片的元素 (即 [1, 2, 3]) 复制到 s2 下标范围 1~1 切片中 (即 [2])
	// s2 从 [2, <2>, 3] 变为 [2, <1>, 3]
	copy(s2[1:2], s1)
	assert.Equal(t, []int{2, 1, 3}, s2)
}
