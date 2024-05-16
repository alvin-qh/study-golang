package series

import (
	"math"
	"study/thirdpart/gota/utils"
	"testing"

	"github.com/go-gota/gota/series"
	"github.com/stretchr/testify/assert"
)

// 测试获取序列的基本属性
//
// 序列的基本属性包括序列的名称, 类型和长度
func TestReadSeries_Basic(t *testing.T) {
	// 创建一个具备名称, Int 类型的序列
	s := series.New([]int{1, 2, 3, 4, 5}, series.Int, "N1")
	assert.NoError(t, s.Error())

	// 获取序列的基本属性, 包括名称, 类型及长度
	assert.Equal(t, "N1", s.Name)
	assert.Equal(t, series.Int, s.Type())
	assert.Equal(t, 5, s.Len())
}

// 测试获取序列指定索引的元素值
func TestReadSeries_Val(t *testing.T) {
	s := series.Ints([]int{1, 2, 3, 4, 5})
	assert.NoError(t, s.Error())

	for i := 0; i < s.Len(); i++ {
		assert.Equal(t, i+1, s.Val(i))
	}
}

// 测试获取序列指定索引位置的元素实例
func TestReadSeries_Elem(t *testing.T) {
	s := series.Ints([]int{1, 2, 3, 4, 5})
	assert.NoError(t, s.Error())

	// 获取序列长度
	assert.Equal(t, 5, s.Len())

	// 获取序列的前两个元素
	e1 := s.Elem(0)
	e2 := s.Elem(1)

	// 获取元素的类型
	assert.Equal(t, series.Int, e1.Type())

	// 判断元素值是否为 `NaN`
	assert.False(t, e1.IsNA())

	// 获取元素不同类型的值
	assert.Equal(t, 1, e1.Val())                         // interface{}
	assert.Equal(t, 1, utils.WithNoError(t, e1.Int))     // int
	assert.Equal(t, 1.0, e1.Float())                     // float64
	assert.Equal(t, "1", e1.String())                    // string
	assert.Equal(t, true, utils.WithNoError(t, e1.Bool)) // bool

	// 进行两个元素比较
	assert.False(t, e2.Eq(e1))       // e1 == e2
	assert.True(t, e2.Neq(e1))       // e1 != e2
	assert.True(t, e2.Greater(e1))   // e2 > e1
	assert.True(t, e2.GreaterEq(e1)) // e2 >= e1
	assert.True(t, e1.Less(e2))      // e1 < e2
	assert.True(t, e1.LessEq(e2))    // e1 <= e2

	// 复制元素
	ec := e1.Copy()
	assert.Equal(t, ec.Type(), e1.Type())
	assert.Equal(t, ec.Val(), e1.Val())

	// 为元素设置值, 设置后原序列也会同步改变
	e1.Set(100)
	e2.Set(200)
	assert.Equal(t, []int{100, 200, 3, 4, 5}, utils.WithNoError(t, s.Int))
}

// 测试将序列转为切片
//
// 可以通过 `Int`, `Float`, `Records`, `Bool` 方法将序列转为对应类型的切片
func TestReadSeries_ToSlice(t *testing.T) {
	s := series.Ints([]int{1, 0, 1, 1, 0})
	assert.NoError(t, s.Error())

	var r any

	// 序列转为整数切片
	r, err := s.Int()
	assert.NoError(t, err)
	assert.Equal(t, []int{1, 0, 1, 1, 0}, r)

	// 序列转为浮点数切片
	r = s.Float()
	assert.Equal(t, []float64{1.0, 0.0, 1.0, 1.0, 0.0}, r)

	// 序列转为字符串切片
	r = s.Records()
	assert.Equal(t, []string{"1", "0", "1", "1", "0"}, r)

	// 序列转为布尔切片
	r, err = s.Bool()
	assert.NoError(t, err)
	assert.Equal(t, []bool{true, false, true, true, false}, r)
}

// 测试序列的计算结果
//
// 序列支持的计算包括: 累加 (`Sum`), 最大值 (`Max`), 最小值 (`Min`), 平均值 (`Mean`),
// 标准差 (`StdDev`), 中位数 (`Median`)
func TestReadSeries_Calculate(t *testing.T) {
	s := series.Ints([]int{1, 2, 3, 4, 5, 6})
	assert.NoError(t, s.Error())

	// 获取整列数据的累加和, 结果为 float64 类型
	r := s.Sum()
	assert.Equal(t, 21.0, r)

	// 获取整列数据的最大值, 结果为 float64 类型
	r = s.Max()
	assert.Equal(t, 6.0, r)

	// 获取整列数据的最小值, 结果为 float64 类型
	r = s.Min()
	assert.Equal(t, 1.0, r)

	// 获取整列数据的平均值, 结果为 float64 类型
	r = s.Mean()
	assert.Equal(t, 3.5, r)

	// 计算序列的标准差结果
	r = s.StdDev()
	assert.Equal(t, 1.8708286933869707, r)

	// 求序列的中值
	r = s.Median()
	assert.Equal(t, 3.5, r)

	// 计算序列的 P 分数位 (用于计算正态分布)
	// 本例中找到一个数值 P (即 3.0), 令小于 P 的数占序列总个数的 40% (约 2 个元素), 大于等于 P 的数占序列总个数 60% (约 4 个元素)
	r = s.Quantile(0.4)
	assert.Equal(t, 3.0, r)

	s = series.Strings([]string{"A", "B", "C", "D", "E"})

	// 获取字符串类型序列中最小值
	sr := s.MinStr()
	assert.Equal(t, "A", sr)

	// 获取字符串类型序列中最大值
	sr = s.MaxStr()
	assert.Equal(t, "E", sr)
}

// 测试对序列元素排序
//
// 通过序列的 `Order` 方法可以对其进行排序, 返回排序后序列的索引, 例如:
//
//	s := series.Ints([]any{1, math.NaN(), 2, 3, 4, 5})
//	ids := s.Order(true)
//
// 表示按逆序获取序列排序后的索引, 结果为: [5, 4, 3, 2, 0, 1], 表示按该索引顺序获取序列元素, 结果是逆序的
//
// 如果元素中包括 `NaN`, 则无论正序逆序, 该元素都会在末尾
func TestReadSeries_Order(t *testing.T) {
	s := series.Ints([]any{1, math.NaN(), 2, 3, 4, 5})
	assert.NoError(t, s.Error())

	// 按元素的排序结果获取序列元素的索引值, 即值小的元素索引在前, 值大的元素索引在后, 如果元素值为 NaN, 则排在末尾, 参数用来控制排序的顺序
	ids := s.Order(true)
	assert.Equal(t, []int{5, 4, 3, 2, 0, 1}, ids)

	// 根据排序得到的索引获取序列元素值
	rs := s.Subset(ids)
	assert.NoError(t, rs.Error())
	assert.Equal(t, []string{"5", "4", "3", "2", "1", "NaN"}, rs.Records())
}

// 测试获取序列的子集
func TestReadSeries_Subset(t *testing.T) {
	s := series.Floats([]float64{1, 2, 3, 4, 5})
	assert.NoError(t, s.Error())

	// 获取索引为 1 的子集
	rs := s.Subset(1)
	assert.NoError(t, rs.Error())
	assert.Equal(t, []float64{2}, rs.Float())

	// 获取索引为 2,3,4 的子集
	rs = s.Subset([]int{2, 3, 4})
	assert.NoError(t, rs.Error())
	assert.Equal(t, []float64{3, 4, 5}, rs.Float())
}

// 测试获取序列的片段
//
// 可以获取序列的两个索引位置之间的部分, 即序列的片段
func TestReadSeries_Slice(t *testing.T) {
	s := series.New([]int{1, 2, 3, 4, 5, 6}, series.Int, "N1")
	assert.NoError(t, s.Error())

	// 获取索引 `[3, 5)` 之间的元素, 即序列 `[4, 5]`
	ss := s.Slice(3, 5)
	assert.NoError(t, ss.Error())

	// 序列片段的类型和原序列一致
	assert.Equal(t, series.Int, ss.Type())

	// 确认序列片段的元素值
	assert.Equal(t, []int{4, 5}, utils.WithNoError(t, ss.Int))
}

// 测试在序列上进行滑动窗口操作
//
// 可以通过序列设置指定元素数量的滑动窗口, 例如:
//
//	s := series.Floats([]float64{1, 2, 3, 4, 5, 6})
//	win := s.Rolling(3)
//
// 上述操作会得到一个滑动窗口, 该滑动窗口将序列分为了 `[1, 2, 3]`, `[2, 3, 4]`, `[3, 4, 5]`, `[4, 5, 6]` 几部分
//
// 通过滑动窗口实例, 可以计算各个窗口的平均值 (`Mean`) 和标准差 (`StdDev`)
func TestReadSeries_SlideWindow(t *testing.T) {
	s := series.Floats([]float64{1, 2, 3, 4, 5, 6})
	assert.NoError(t, s.Error())

	// 获取序列数据的一个窗口, 本例中没 3 个值为一个窗口, 则窗口结果为 [1, 2, 3], [2, 3, 4], [3, 4, 5], [4, 5, 6]
	win := s.Rolling(3)

	// 求各个窗口的平均值, 因为从序列的第 3 个元素开始才能形成窗口, 所以结果的前两个值为 NaN
	r := win.Mean()
	assert.NoError(t, r.Error())
	assert.Equal(t, []float64{2, 3, 4, 5}, r.Float()[2:])

	// 求各个窗口的标准差, 因为从序列的第 3 个元素开始才能形成窗口, 所以结果的前两个值为 NaN
	r = win.StdDev()
	assert.NoError(t, r.Error())
	assert.Equal(t, []float64{1, 1, 1, 1}, r.Float()[2:])
}
