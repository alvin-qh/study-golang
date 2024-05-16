package series

import (
	"study/thirdpart/gota/utils"
	"testing"

	"github.com/go-gota/gota/series"
	"github.com/stretchr/testify/assert"
)

// 测试在序列中追加元素
func TestUpdateSeries_Append(t *testing.T) {
	s := series.New([]int{1, 2, 3}, series.Int, "N1")
	assert.NoError(t, s.Error())

	// 在序列末尾追加新的元素值
	s.Append([]int{4, 5, 6})
	assert.NoError(t, s.Error())

	// 确认序列中已包括了新元素
	assert.Equal(t, 6, s.Len())
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, utils.WithNoError(t, s.Int))
}

// 测试合并两个序列
func TestUpdateSeries_Concat(t *testing.T) {
	s1 := series.New([]int{1, 2, 3}, series.Int, "N1")
	s2 := series.New([]int{4, 5, 6}, series.Int, "N2")

	// 合并两个序列
	s := s1.Concat(s2)
	assert.NoError(t, s.Error())

	// 合并后的序列名和第一个序列相同
	assert.Equal(t, "N1", s.Name)

	// 确认返回值包括两个序列的值
	assert.Equal(t, 6, s.Len())
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, utils.WithNoError(t, s.Int))
}

// 测试设置序列中的元素值
//
// 可以设置指定索引位置的元素值, 也可以设置一系列指定位置的元素值
func TestUpdateSeries_Set(t *testing.T) {
	s := series.New([]int{1, 2, 3}, series.Int, "N1")
	assert.NoError(t, s.Error())

	// 设置指定位置的值
	s.Set(0, series.Ints(100))
	assert.Equal(t, []int{100, 2, 3}, utils.WithNoError(t, s.Int))

	// 设置一系列指定位置的对应值
	s.Set([]int{1, 2}, series.Ints([]int{200, 400}))
	assert.Equal(t, []int{100, 200, 400}, utils.WithNoError(t, s.Int))
}

// 测试序列的映射
//
// 通过序列的 `Map` 方法可以将序列中的元素值映射为另一个值, 并将所有映射后的结果组成新序列
//
// 映射后的序列和原序列具有相同的类型和长度, 只是元素值被映射为新值
func TestReadSeries_Map(t *testing.T) {
	s := series.New([]int{1, 2, 3, 4, 5, 6}, series.Int, "N1")
	assert.NoError(t, s.Error())

	// 将序列中的元素值进行计算后组成新的序列
	r := s.Map(func(e series.Element) series.Element {
		// 复制当前元素
		r := e.Copy()

		// 获取当前处理的元素值
		n, err := e.Int()

		// 如果元素值为奇数或 NaN, 则将元素值设置为 0
		if err != nil || n%2 != 0 {
			r.Set(0)
		}
		return r
	})
	assert.NoError(t, r.Error())

	// 映射后的序列和原序列长度一致
	assert.Equal(t, s.Len(), r.Len())

	// 映射后的序列和原序列类型一致
	assert.Equal(t, series.Int, r.Type())

	// 确认映射结果
	assert.Equal(t, []int{0, 2, 0, 4, 0, 6}, utils.WithNoError(t, r.Int))
}
