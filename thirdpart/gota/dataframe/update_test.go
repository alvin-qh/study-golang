package dataframe

import (
	"fmt"
	"testing"

	df "github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/stretchr/testify/assert"
)

// 通过 Map 创建 `DataFrame` 实例
func makeDFForUpdate(t *testing.T) df.DataFrame {
	frame := LoadSliceMap(map[string][]any{
		"Id":       {1, 2, 3, 4, 5, 6},
		"Name":     {"Alvin", "Emma", "Lucy", "Tom", "Arthur", "Jenny"},
		"Gender":   {"M", "F", "F", "M", "M", "F"},
		"Age":      {42, 38, 51, 34, 21, 33},
		"Salaries": {32000, 18000, 9500, 10000, 7850, 9200},
	})

	assert.NoError(t, frame.Error())

	fmt.Println(frame)
	return frame
}

// 测试修改 `DataFrame` 的列名
func TestUpdateDataFrame_Names(t *testing.T) {
	frame := makeDFForUpdate(t)

	// 重新设定每一列的列名
	err := frame.SetNames("Id_", "Name_", "Gender_", "Age_", "Salaries_")
	assert.NoError(t, err)

	assert.Equal(t, []string{"Id_", "Name_", "Gender_", "Age_", "Salaries_"}, frame.Names())

	fmt.Println(frame)
}

// 测试修改 `DataFrame` 数据
//
// 每次修改 `DataFrame` 一行数据, 返回和被修改 `DataFrame` 本身相同的实例
func TestUpdateDataFrame_Set(t *testing.T) {
	frame := makeDFForUpdate(t)

	// 修改第一行数据
	rs := frame.Set(0, df.New(
		series.Ints(43),
		series.Strings("M"),
		series.Ints(1),
		series.Strings("Alvin_"),
		series.Ints(28000),
	))
	assert.NoError(t, rs.Error())
	assert.Equal(t, frame, rs)

	fmt.Println(rs)

	// 确认修改后的数据
	assert.Equal(t, [][]string{
		{"Age", "Gender", "Id", "Name", "Salaries"},
		{"43", "M", "1", "Alvin_", "28000"},
		{"38", "F", "2", "Emma", "18000"},
		{"51", "F", "3", "Lucy", "9500"},
		{"34", "M", "4", "Tom", "10000"},
		{"21", "M", "5", "Arthur", "7850"},
		{"33", "F", "6", "Jenny", "9200"},
	}, rs.Records())

	// 修改第二行和第四行数据
	rs = frame.Set(
		[]int{1, 3},
		df.New(
			series.Ints([]int{39, 35}),
			series.Strings([]string{"F", "M"}),
			series.Ints([]int{2, 4}),
			series.Strings([]string{"Emma_", "Tom_"}),
			series.Ints([]int{16000, 8000}),
		),
	)
	assert.NoError(t, rs.Error())
	assert.Equal(t, frame, rs)

	fmt.Println(rs)

	// 确认修改后的数据
	assert.Equal(t, [][]string{
		{"Age", "Gender", "Id", "Name", "Salaries"},
		{"43", "M", "1", "Alvin_", "28000"},
		{"39", "F", "2", "Emma_", "16000"},
		{"51", "F", "3", "Lucy", "9500"},
		{"35", "M", "4", "Tom_", "8000"},
		{"21", "M", "5", "Arthur", "7850"},
		{"33", "F", "6", "Jenny", "9200"},
	}, rs.Records())
}

// 测试依据指定列对 `DataFrame` 数据进行排序
//
// 排序需要指定依据哪一列的数据进行排序 (正序, 倒序), 对应 `dataframe.Sort` 和 `dataframe.RevSort` 函数
//
// 排序结果为包含新顺序数据的 `DataFrame` 实例
func TestUpdateDataFrame_Arrange(t *testing.T) {
	frame := makeDFForUpdate(t)

	// 先按 `Gender` 字段排序, 在 `Gender` 相同时按 `Age` 排序
	rs := frame.Arrange(df.Sort("Gender"), df.RevSort("Age"))
	assert.NoError(t, rs.Error())

	fmt.Println(rs)

	// 确认排序结果
	assert.Equal(t, [][]string{
		{"Age", "Gender", "Id", "Name", "Salaries"},
		{"51", "F", "3", "Lucy", "9500"},
		{"38", "F", "2", "Emma", "18000"},
		{"33", "F", "6", "Jenny", "9200"},
		{"42", "M", "1", "Alvin", "32000"},
		{"34", "M", "4", "Tom", "10000"},
		{"21", "M", "5", "Arthur", "7850"},
	}, rs.Records())
}

// 测试合并两个 `DataFrame` 实例
//
// 合并结果为同时包含被合并的两个 `DataFrame` 数据的新 `DataFrame` 实例
func TestUpdateDataFrame_Cbind(t *testing.T) {
	frame := makeDFForUpdate(t)

	// 本例将当前 `DataFrame` 和其副本进行合并
	rs := frame.CBind(frame.Copy())

	assert.NoError(t, rs.Error())

	fmt.Println(rs)

	// 确认合并结果
	assert.Equal(t, [][]string{
		{
			"Age_0", "Gender_0", "Id_0", "Name_0", "Salaries_0",
			"Age_1", "Gender_1", "Id_1", "Name_1", "Salaries_1",
		},
		{"42", "M", "1", "Alvin", "32000", "42", "M", "1", "Alvin", "32000"},
		{"38", "F", "2", "Emma", "18000", "38", "F", "2", "Emma", "18000"},
		{"51", "F", "3", "Lucy", "9500", "51", "F", "3", "Lucy", "9500"},
		{"34", "M", "4", "Tom", "10000", "34", "M", "4", "Tom", "10000"},
		{"21", "M", "5", "Arthur", "7850", "21", "M", "5", "Arthur", "7850"},
		{"33", "F", "6", "Jenny", "9200", "33", "F", "6", "Jenny", "9200"},
	}, rs.Records())
}

// 对 `DataFrame` 数据列进行映射
//
// 通过一个回调函数处理 `DataFrame` 中的每个 `Series`, 并返回新的 `Series`, 最终由这些新 `Series`
// 组成的新 `DataFrame` 实例
func TestUpdateDataFrame_Capply(t *testing.T) {
	frame := makeDFForUpdate(t)

	rs := frame.Capply(func(s series.Series) series.Series {
		s = s.Copy()

		if s.Name == "Salaries" {
			for i := 0; i < s.Len(); i++ {
				elem := s.Elem(i)
				s.Set(i, series.Floats(elem.Float()*0.8))
			}
		}
		return s
	})
	assert.NoError(t, rs.Error())

	fmt.Println(rs)

	// 确认 `DataFrame` 中各个 `Series` 的处理结果
	assert.Equal(t, [][]string{
		{"Age", "Gender", "Id", "Name", "Salaries"},
		{"42", "M", "1", "Alvin", "25600"},
		{"38", "F", "2", "Emma", "14400"},
		{"51", "F", "3", "Lucy", "7600"},
		{"34", "M", "4", "Tom", "8000"},
		{"21", "M", "5", "Arthur", "6280"},
		{"33", "F", "6", "Jenny", "7360"},
	}, rs.Records())
}

// 测试修改 DataFrame 的一列数据
//
// 可通过 `Mutate` 方法修改 `DataFrame` 中的一列数据, 并返回修改后的 `DataFrame`
//
// 如果修改的序列名称在 `DataFrame` 中存在, 则取代 `DataFrame` 中的同名列, 否则新增一列
func TestUpdateDataFrame_Mutate(t *testing.T) {
	frame := makeDFForUpdate(t)

	// 修改 DataFrame 的一列数据
	rs := frame.Mutate(series.New([]int{10, 20, 30, 40, 50, 60}, series.Int, "Id"))
	assert.NoError(t, rs.Error())

	fmt.Println(rs)

	// 确认修改结果
	assert.Equal(t, [][]string{
		{"Age", "Gender", "Id", "Name", "Salaries"},
		{"42", "M", "10", "Alvin", "32000"},
		{"38", "F", "20", "Emma", "18000"},
		{"51", "F", "30", "Lucy", "9500"},
		{"34", "M", "40", "Tom", "10000"},
		{"21", "M", "50", "Arthur", "7850"},
		{"33", "F", "60", "Jenny", "9200"},
	}, rs.Records())

	// 为 DataFrame 添加一列数据
	rs = frame.Mutate(series.New([]bool{true, false, false, true, true, true}, series.Bool, "Regular"))
	assert.NoError(t, rs.Error())

	fmt.Println(rs)

	// 确认修改结果
	assert.Equal(t, [][]string{
		{"Age", "Gender", "Id", "Name", "Salaries", "Regular"},
		{"42", "M", "1", "Alvin", "32000", "true"},
		{"38", "F", "2", "Emma", "18000", "false"},
		{"51", "F", "3", "Lucy", "9500", "false"},
		{"34", "M", "4", "Tom", "10000", "true"},
		{"21", "M", "5", "Arthur", "7850", "true"},
		{"33", "F", "6", "Jenny", "9200", "true"},
	}, rs.Records())
}
