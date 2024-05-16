package utils

import (
	"cmp"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试二分查找
func TestSlices_BinarySearch(t *testing.T) {
	// 对递增有序序列执行二分查找
	t.Run("slices.BinarySearch", func(t *testing.T) {
		// 生成有序递增序列
		ns := Range(1, 1000, 1)

		// 确认切片有序递增
		assert.True(t, slices.IsSorted(ns))

		// 对递增有序序列执行二分查找
		v, ok := slices.BinarySearch(ns, 100)
		assert.True(t, ok)
		assert.Equal(t, 99, v)
	})

	// 测试通过回调进行二分查找
	// 如果序列不是递增有序, 或序列元素无法直接通过 `>` 和 `<` 比较,
	// 则可以使用二分查找的回调函数版本, 回调函数用于对序列中的两个元素进行比较
	t.Run("slices.BinarySearchFunc", func(t *testing.T) {
		// 生成有序递减序列
		ns := Range(1, 1000, 1)
		slices.Reverse(ns)

		// 确认序列有序递减
		assert.True(t, slices.IsSortedFunc(ns, func(e1, e2 int) int {
			return cmp.Compare(e2, e1)
		}))

		// 指定排序方法的二分查找法
		v, ok := slices.BinarySearchFunc(ns, 100, func(e1, e2 int) int {
			return cmp.Compare(e2, e1)
		})
		assert.True(t, ok)
		assert.Equal(t, 899, v)
	})
}

// 测试收缩切片中未使用的 Capacity 空间
//
// 如果一个切片通过 `append` 函数添加完元素后, 不会再进行修改, 则可以对其进行收缩来节省内存空间
func TestSlices_Clip(t *testing.T) {
	// 声明一个 Capacity 为 100 的切片
	s := make([]int, 0, 100)

	// 向切片中添加 20 个元素, 查看切片的长度和 Capacity 值
	s = append(s, Range(0, 20, 1)...)
	assert.Len(t, s, 20)
	assert.Equal(t, 100, cap(s))

	// 收缩切片, 去除未使用的 Capacity 空间
	s = slices.Clip(s)
	assert.Len(t, s, 20)
	assert.Equal(t, 20, cap(s))
	assert.Equal(t, Range(0, 20, 1), s)
}

// 测试对切片 Capacity 值进行扩容
//
// 可对切片的 Capacity 值扩容到 `n`, 扩容后的切片可以支持至少 `n` 次 `append` 操作而无需重新分配内存
func TestSlices_Grow(t *testing.T) {
	s := []int{1, 2, 3}
	assert.Len(t, s, 3)
	assert.Equal(t, 3, cap(s))

	s = slices.Grow(s, 100)
	assert.Len(t, s, 3)
	assert.Equal(t, 112, cap(s))

	assert.PanicsWithValue(t, "cannot be negative", func() {
		_ = slices.Grow(s, -1)
	})
}

type Value struct {
	S []int
}

// 测试复制切片
//
// `Clone` 函数对切片执行"浅拷贝"复制
func TestSlices_Clone(t *testing.T) {
	s := []Value{
		{[]int{1}},
		{[]int{2}},
		{[]int{3}},
		{[]int{4}},
		{[]int{5}},
	}

	// 对切片进行复制
	sc := slices.Clone(s)
	assert.Len(t, sc, 5)
	assert.Equal(t, s, sc)

	// 修改复制后的切片元素值, 确认原切片元素也同步修改, 表明为浅拷贝
	sc[1].S[0] *= 10
	assert.Equal(t, []int{20}, s[1].S)
}

// 测试切片去重
//
// 可以被去重的切片, 其重复元素必须为相邻元素, 即 `[1, 1, 3, 3, 2, 2, 6, 6, 4, 4, 4]`,
// 不相邻的元素即使相同也无法去重
//
// 使相同元素相邻的简便方法即对切片排序
func TestSlices_Compact(t *testing.T) {
	// 测试相同相邻元素去重
	t.Run("slices.Compact", func(t *testing.T) {
		s := []int{1, 1, 3, 3, 3, 2, 2, 5, 5, 5, 4, 4}

		s = slices.Compact(s)
		assert.Equal(t, []int{1, 3, 2, 5, 4}, s)

		s = []int{1, 1, 2, 3, 3, 2, 4, 5, 4, 4, 5, 6, 6}

		s = slices.Compact(s)
		assert.Equal(t, []int{1, 2, 3, 2, 4, 5, 4, 5, 6}, s)
	})

	// 测试通过回调进行去重
	// 如果序列元素无法直接通过 `>` 和 `<` 比较, 则可以使用去重的回调函数版本,
	// 回调函数用于对序列中的两个元素进行比较
	t.Run("slices.CompactFunc", func(t *testing.T) {
		s := []int{1, 1, 2, 3, 3, 2, 4, 5, 4, 4, 5, 6, 6}

		// 历史元素缓冲
		cache := make(map[int]struct{})

		// 通过回调函数对切片中的不相邻重复项进行去重
		s = slices.CompactFunc(s, func(i, j int) bool {
			// 判断相邻元素是否相等
			if i == j {
				return true
			}

			// 判断元素和之前的元素是否相等
			if _, ok := cache[i]; ok {
				return true
			}

			// 将未出现的元素加入缓冲
			cache[i] = struct{}{}
			return false
		})
		assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, s)
	})
}

// 测试两个切片比较
//
// 所谓切片比较, 即用第一个切片的每个元素和第二个切片同位置元素进行比较
//
// 如果可供比较的元素都相等, 则进一步比较切片的长度
func TestSlices_Compare(t *testing.T) {
	s1 := []int{1, 2, 3}

	// 测试两个切片比较
	t.Run("slices.Compare", func(t *testing.T) {
		// 元素完全相同的两个切片相等
		s2 := []int{1, 2, 3}
		assert.Equal(t, 0, slices.Compare(s1, s2))

		// 相同位置元素值不同, 则切片的比较结果为该位置元素值的比较结果
		s2 = []int{1, 3, 3}
		assert.Equal(t, -1, slices.Compare(s1, s2))

		s2 = []int{1, 1, 3}
		assert.Equal(t, 1, slices.Compare(s1, s2))

		// 相同位置元素不同后, 则无需比较其它元素
		s2 = []int{2}
		assert.Equal(t, -1, slices.Compare(s1, s2))

		s2 = []int{0}
		assert.Equal(t, 1, slices.Compare(s1, s2))

		// 两个切片长度不同, 且对应位置的元素值相同, 则切片比较结果即切片长度的比较结果
		s2 = []int{1, 2, 3, 4}
		assert.Equal(t, -1, slices.Compare(s1, s2))

		s2 = []int{1, 2}
		assert.Equal(t, 1, slices.Compare(s1, s2))
	})

	// 如果切片元素无法直接通过 `>` 和 `<` 比较,
	// 则可以使用切片比较的回调函数版本, 回调函数用于对两个切片中相同位置元素进行比较
	t.Run("slices.CompareFunc", func(t *testing.T) {
		s2 := []int{10, 20, 30}
		assert.Equal(t, 0, slices.CompareFunc(s1, s2, func(e1, e2 int) int {
			return cmp.Compare(e1*10, e2)
		}))
	})
}

// 测试按顺序连接多个切片
func TestSlices_Concat(t *testing.T) {
	s1 := []int{1, 2, 3}
	s2 := []int{4, 5, 6}
	s3 := []int{7, 8, 9}

	// 按顺序连接三个切片
	s := slices.Concat(s1, s2, s3)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9}, s)
}

// 测试切片是否包含指定元素
func TestSlices_Contains(t *testing.T) {
	s := []int{1, 2, 3, 4}

	// 测试切片是否包含指定元素
	t.Run("slices.Contains", func(t *testing.T) {
		assert.True(t, slices.Contains(s, 3))
		assert.False(t, slices.Contains(s, 5))
	})

	// 通过回调函数结果确认切片中是否包含指定元素
	t.Run("slices.Contains", func(t *testing.T) {
		assert.True(t, slices.ContainsFunc(s, func(e int) bool {
			return e == 3
		}))
	})
}

// 测试删除切片中的元素
func TestSlices_Delete(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// 通过下标删除切片中的元素
	t.Run("slices.Delete", func(t *testing.T) {
		sc := slices.Clone(s)

		// 删除下标为 2 的元素
		sc = slices.Delete(sc, 2, 3)
		assert.Equal(t, []int{1, 2, 4, 5, 6, 7, 8, 9, 10}, sc)

		sc = slices.Clone(s)

		// 删除下标为 2, 3 的元素
		sc = slices.Delete(sc, 2, 4)
		assert.Equal(t, []int{1, 2, 5, 6, 7, 8, 9, 10}, sc)

		sc = slices.Clone(s)

		// 删除下标 2 及之后的所有元素
		sc = slices.Delete(sc, 2, len(sc))
		assert.Equal(t, []int{1, 2}, sc)
	})

	// 通过回调函数删除切片中符合条件的元素
	t.Run("slices.DeleteFunc", func(t *testing.T) {
		sc := slices.Clone(s)

		// 删除值为偶数的元素
		sc = slices.DeleteFunc(sc, func(e int) bool {
			return e%2 == 0
		})
		assert.Equal(t, []int{1, 3, 5, 7, 9}, sc)
	})
}

// 测试判读两个切片是否相等
//
// 判读两个切片是否相等, 即对比两个切片对应位置元素是否相等
func TestSlices_Equal(t *testing.T) {
	s1 := []int{1, 2, 3}

	// 测试两个切片是否相等
	t.Run("slices.Equal", func(t *testing.T) {
		s2 := []int{1, 2, 3}
		assert.True(t, slices.Equal(s1, s2))

		s2 = []int{1, 2, 4}
		assert.False(t, slices.Equal(s1, s2))
	})

	// 通过回调函数判读两个切片是否相等
	t.Run("slices.EqualFunc", func(t *testing.T) {
		s2 := []int{10, 20, 30}

		// 通过回调函数比较两个切片的对应元素
		assert.True(t, slices.EqualFunc(s1, s2, func(e1, e2 int) bool {
			return cmp.Compare(e1*10, e2) == 0
		}))
	})
}

// 测试获取切片中指定元素的下标
//
// 对切片中元素进行查找, 如果查找成功, 则返回对应元素的下标, 否则返回 -1 表示查找失败
func TestSlices_Index(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// 通过元素值查找对应下标
	t.Run("slices.Index", func(t *testing.T) {
		// 查找元素 3 在切片中首次出现的下标
		idx := slices.Index(s, 3)
		assert.Equal(t, 2, idx)

		// 查找不存在的元素, 返回 -1
		idx = slices.Index(s, 0)
		assert.Equal(t, -1, idx)
	})

	// 通过回调函数查找对应下标
	t.Run("slices.IndexFunc", func(t *testing.T) {
		// 查找可被 3 整除的第一个元素下标
		idx := slices.IndexFunc(s, func(e int) bool {
			return e%3 == 0
		})
		assert.Equal(t, 2, idx)

		// 查找不存在元素, 返回 -1
		idx = slices.IndexFunc(s, func(e int) bool {
			return e == 0
		})
		assert.Equal(t, -1, idx)
	})
}

// 测试在切片中插入元素
//
// 可以在切片的指定索引前插入指定的一个或多个元素
//
// 可以在切片的末尾 (相当于 最大索引 + 1) 插入指定的一个或多个元素, 相当于 `append` 操作
//
// 如果插入位置 < 0 或者 > 切片长度, 则会引发 Panic
func TestSlices_Insert(t *testing.T) {
	s := []int{1, 2, 3}

	// 在下标 1 之前插入 2 个元素
	s = slices.Insert(s, 1, 10, 20)
	assert.Equal(t, []int{1, 10, 20, 2, 3}, s)

	// 在最后一个元素前插入 3 个元素
	s = slices.Insert(s, len(s)-1, 30, 40, 50)
	assert.Equal(t, []int{1, 10, 20, 2, 30, 40, 50, 3}, s)

	// 在末尾插入 1 个元素
	s = slices.Insert(s, len(s), 60)
	assert.Equal(t, []int{1, 10, 20, 2, 30, 40, 50, 3, 60}, s)

	// 测试插入下标不存在时的情况
	assert.PanicsWithError(t, "runtime error: slice bounds out of range [100:9]", func() {
		_ = slices.Insert(s, 100, 70)
	})
}

// 判读切片是否有序
func TestSlices_IsSorted(t *testing.T) {
	// 判读切片是否有序递增
	t.Run("slices.IsSorted", func(t *testing.T) {
		s := []int{1, 2, 3}
		assert.True(t, slices.IsSorted(s))

		s = []int{1, 3, 2}
		assert.False(t, slices.IsSorted(s))
	})

	// 通过回调函数结果判读切片是否有序
	t.Run("slices.IsSortedFunc", func(t *testing.T) {
		s := []int{1, 2, 3}
		assert.True(t, slices.IsSortedFunc(s, func(e1, e2 int) int {
			return cmp.Compare(e1, e2)
		}))

		s = []int{3, 2, 1}
		assert.True(t, slices.IsSortedFunc(s, func(e1, e2 int) int {
			return cmp.Compare(e2, e1)
		}))
	})
}

// 测试获取切片中的最大值或最小值
func TestSlices_Max_Min(t *testing.T) {
	// 获取切片中的最大值
	t.Run("slices.Max", func(t *testing.T) {
		s := []int{1, 2, 3}
		assert.Equal(t, 3, slices.Max(s))
	})

	// 通过回调函数进行比较, 获取切片中的最大值
	// 如果切片元素本身不支持通过 > 进行比较, 则可以使用该函数
	t.Run("slices.MaxFunc", func(t *testing.T) {
		s := []int{1, 2, 3}
		assert.Equal(t, 3, slices.MaxFunc(s, func(e1, e2 int) int {
			return cmp.Compare(e1, e2)
		}))
	})

	// 获取切片中的最小值
	t.Run("slices.Min", func(t *testing.T) {
		s := []int{1, 2, 3}
		assert.Equal(t, 1, slices.Min(s))
	})

	// 通过回调函数进行比较, 获取切片中的最小值
	// 如果切片元素本身不支持通过 < 进行比较, 则可以使用该函数
	t.Run("slices.MinFunc", func(t *testing.T) {
		s := []int{1, 2, 3}
		assert.Equal(t, 1, slices.MinFunc(s, func(e1, e2 int) int {
			return cmp.Compare(e1, e2)
		}))
	})
}

// 测试对切片进行排序
func TestSlices_Sort(t *testing.T) {
	// 按照递增顺序排序
	t.Run("slices.Sort", func(t *testing.T) {
		s := []int{3, 1, 2}

		// 按递增顺序排序
		slices.Sort(s)
		assert.Equal(t, []int{1, 2, 3}, s)
	})

	// 根据回调函数进行排序
	// 对于按递减顺序排序, 或者切片元素不支持 `>` 或 `<` 比较时, 可以使用该函数
	t.Run("slices.SortFunc", func(t *testing.T) {
		s := []int{3, 1, 2}

		// 通过回调函数, 按递减顺序排序
		slices.SortFunc(s, func(e1, e2 int) int {
			return cmp.Compare(e2, e1)
		})
		assert.Equal(t, []int{3, 2, 1}, s)

		s = []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

		// 通过回调函数, 按奇数在前, 偶数在后排序
		slices.SortFunc(s, func(e1, e2 int) int {
			if e1%2 == 0 && e2%2 != 0 {
				return 1
			}
			if e1%2 != 0 && e2%2 == 0 {
				return -1
			}
			return 0
		})
		assert.Equal(t, []int{1, 3, 5, 7, 9, 2, 4, 6, 8, 10}, s)
	})

	// 稳定排序
	// 稳定排序和 `slices.SortFunc` 函数类似, 但排序结果是稳定的
	// 所谓排序结果稳定, 即排序后可以保持原序列中相等元素的顺序不变
	t.Run("slices.SortStableFunc", func(t *testing.T) {
		s := []float64{1.1, 1.2, 1.3, 3.1, 3.2, 3.3, 2.1, 2.2, 2.3, 5.1, 5.2, 5.3, 4.1, 4.2, 4.3}

		// 忽略切片元素小数部分进行排序
		// 使用非稳定排序, 排序后原序列中相等元素的顺序发生了改变, 如 5.2 和 5.3 的顺序发生变化
		slices.SortFunc(s, func(e1, e2 float64) int {
			return cmp.Compare(int(e2), int(e1))
		})
		assert.Equal(t, []float64{5.1, 5.3, 5.2, 4.3, 4.2, 4.1, 3.1, 3.3, 3.2, 2.2, 2.3, 2.1, 1.1, 1.3, 1.2}, s)

		s = []float64{1.1, 1.2, 1.3, 3.1, 3.2, 3.3, 2.1, 2.2, 2.3, 5.1, 5.2, 5.3, 4.1, 4.2, 4.3}

		// 忽略切片元素小数部分进行排序
		// 使用稳定排序, 排序后原序列中相等元素的顺序保持不变
		slices.SortStableFunc(s, func(e1, e2 float64) int {
			return cmp.Compare(int(e2), int(e1))
		})
		assert.Equal(t, []float64{5.1, 5.2, 5.3, 4.1, 4.2, 4.3, 3.1, 3.2, 3.3, 2.1, 2.2, 2.3, 1.1, 1.2, 1.3}, s)
	})
}

// 测试替换切片中指定下标的元素值
func TestSlices_Replace(t *testing.T) {
	s := []int{1, 2, 3}

	// 替换下标为 1 的一个元素
	s = slices.Replace(s, 1, 2, 20)
	assert.Equal(t, []int{1, 20, 3}, s)

	// 替换下标 1, 2 的两个元素
	s = slices.Replace(s, 1, 2, 30, 40)
	assert.Equal(t, []int{1, 30, 40, 3}, s)
}

// 测试切片反转
func TestSlices_Reverse(t *testing.T) {
	s := []int{1, 2, 3}

	// 反转切片
	slices.Reverse(s)
	assert.Equal(t, []int{3, 2, 1}, s)
}

// 测试切片乱序
func TestSlices_Shuffle(t *testing.T) {
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// 乱序切片, 确认切片元素顺序被打乱
	Shuffle(s, 100)
	assert.NotEqual(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, s)

	// 重新排序乱序后的切片, 恢复原始顺序
	slices.Sort(s)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, s)
}
