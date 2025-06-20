// 测试创建序列实例
//
// 可以通过四种方式创建序列
//  1. 使用 `series.New` 函数创建一个空序列, 需要设置序列的类型和名称;
//  2. 使用 `series.Ints`, `series.Floats`, `series.Bools` 和 `series.Strings` 函数创建特定类型的序列;
//  3. 使用 `series.Empty` 函数创建一个类型和所给序列相同但长度为 `0` 的新序列;
//  4. 使用 `series.Copy` 函数创建一个类型和内容均和所给序列相同的新序列
package series

import (
	"math"
	"testing"

	"thirdpart/gota/utils"

	"github.com/go-gota/gota/series"
	"github.com/stretchr/testify/assert"
)

// 测试通过 `series.New` 函数创建数据序列
//
// 通过 `series.New` 方法可以创建 `series.Int`, `series.Float`, `series.Bool` 和 `series.String` 类型的序列
//
// 需要在创建时指定序列的类型, 另外需指定序列的名称
func TestCreateSeries_New(t *testing.T) {
	// 创建一个 Int 类型空序列
	s := series.New([]int{}, series.Int, "N1")
	assert.NoError(t, s.Error())

	assert.Equal(t, "N1", s.Name)
	assert.Equal(t, series.Int, s.Type())
	assert.Equal(t, 0, s.Len())
	assert.Equal(t, []int{}, utils.WithNoError(t, s.Int))

	// 创建一个只包含一个元素的序列
	s = series.New(1.0, series.Float, "N2")
	assert.NoError(t, s.Error())

	assert.Equal(t, "N2", s.Name)
	assert.Equal(t, series.Float, s.Type())
	assert.Equal(t, 1, s.Len())
	assert.Equal(t, []float64{1.0}, s.Float())

	// 创建一个包含所给切片元素的序列
	s = series.New([]bool{true, false, false, true}, series.Bool, "N3")
	assert.NoError(t, s.Error())

	assert.Equal(t, "N3", s.Name)
	assert.Equal(t, series.Bool, s.Type())
	assert.Equal(t, 4, s.Len())
	assert.Equal(t, []bool{true, false, false, true}, utils.WithNoError(t, s.Bool))
}

// 测试通过切片创建序列
//
// 可以通过 `series.Ints`, `series.Floats`, `series.Bools` 和 `series.Strings`
// 函数通过对应的切片创建特定类型的序列
//
// 通过 `series.Ints`, `series.Floats`, `series.Bools` 和 `series.Strings`
// 方法创建的序列没有名称, 需要通过 `Name` 属性设置
func TestCreateSeries_BySlice(t *testing.T) {
	// 创建一个整数类型序列
	s := series.Ints([]int{1, 2, 3, 4})
	assert.NoError(t, s.Error())

	assert.Equal(t, "", s.Name)
	assert.Equal(t, series.Int, s.Type())
	assert.Equal(t, 4, s.Len())
	assert.Equal(t, []int{1, 2, 3, 4}, utils.WithNoError(t, s.Int))

	// 创建一个浮点数类型序列
	s = series.Floats([]float64{1.1, 2.2, 3.3, 4.4})
	assert.NoError(t, s.Error())

	assert.Equal(t, "", s.Name)
	assert.Equal(t, series.Float, s.Type())
	assert.Equal(t, 4, s.Len())
	assert.Equal(t, []float64{1.1, 2.2, 3.3, 4.4}, s.Float())

	// 创建一个布尔类型序列
	s = series.Bools([]int{1, 0, 0, 1})
	assert.NoError(t, s.Error())

	assert.Equal(t, "", s.Name)
	assert.Equal(t, series.Bool, s.Type())
	assert.Equal(t, 4, s.Len())
	assert.Equal(t, []bool{true, false, false, true}, utils.WithNoError(t, s.Bool))

	// 创建一个字符串类型序列
	s = series.Strings([]string{"A", "B", "C", "D"})
	assert.NoError(t, s.Error())

	assert.Equal(t, "", s.Name)
	assert.Equal(t, series.String, s.Type())
	assert.Equal(t, 4, s.Len())
	assert.Equal(t, []string{"A", "B", "C", "D"}, s.Records())
}

// 测试通过 `nil` 值创建序列
//
// 如果通过一个 `nil` 值创建序列, 则会创建一个长度为 `1` 且内容为 `NaN` 的序列
func TestCreateSeries_ByNil(t *testing.T) {
	// 如果为数值类型序列传递 nil 值, 则相当于序列中包含一个 `math.NaN()` 值
	s := series.Ints(nil)
	assert.NoError(t, s.Error())

	assert.Equal(t, series.Int, s.Type())
	assert.Equal(t, 1, s.Len())
	assert.True(t, math.IsNaN(s.Float()[0]))

	// 如果为字符串类型序列传递 nil 值, 则相当于序列中包含一个 `"NaN"` 字符串
	s = series.Strings(nil)
	assert.NoError(t, s.Error())

	assert.Equal(t, series.String, s.Type())
	assert.Equal(t, 1, s.Len())
	assert.Equal(t, []string{"NaN"}, s.Records())

	// 如果为布尔类型序列传递 nil 值, 则相当于序列中包含一个 `"NaN"` 字符串, 且无法转为 Bool 类型
	s = series.Bools(nil)
	assert.NoError(t, s.Error())

	assert.Equal(t, series.Bool, s.Type())
	assert.Equal(t, 1, s.Len())
	assert.EqualError(t, utils.WithError(t, s.Bool), "can't convert NaN to bool")
}

// 创建一个空序列实例
//
// 通过现有序列的 `Empty` 方法可以创建一个和原序列类型, 名称一致, 但长度为零的空序列
func TestCreateSeries_Empty(t *testing.T) {
	s1 := series.Strings([]string{"A", "B", "C", "D"})
	s1.Name = "N1"

	assert.NoError(t, s1.Error())

	// 创建一个和 `s1` 序列名称, 类型相同但长度为 `0` 的序列
	s2 := s1.Empty()
	assert.NoError(t, s2.Error())

	assert.Equal(t, "N1", s2.Name)
	assert.Equal(t, series.String, s2.Type())
	assert.Equal(t, 0, s2.Len())
}

// 创建序列副本
//
// 可以通过现有序列的 `Copy` 方法创建其副本, 副本的名称, 类型及内容和原序列一致
func TestCreateSeries_Copy(t *testing.T) {
	s1 := series.Strings([]string{"A", "B", "C", "D"})
	s1.Name = "N1"

	assert.NoError(t, s1.Error())

	// 创建 `s1` 序列的副本, 副本的名称, 类型及内容和原序列一致
	s2 := s1.Copy()
	assert.NoError(t, s2.Error())

	assert.Equal(t, "N1", s2.Name)
	assert.Equal(t, series.String, s2.Type())
	assert.Equal(t, 4, s2.Len())
	assert.Equal(t, s1.Records(), s2.Records())
}
