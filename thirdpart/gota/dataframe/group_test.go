package dataframe

import (
	"fmt"
	"testing"

	df "github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/stretchr/testify/assert"
)

// 通过 Map 创建 `DataFrame` 实例
func makeDFForGrouping(t *testing.T) df.DataFrame {
	frame := LoadSliceMap(map[string][]any{
		"Name":     {"Alvin", "Emma", "Lucy", "Tom", "Arthur"},
		"Gender":   {"M", "F", "F", "M", "M"},
		"Age":      {42, 38, 51, 34, 21},
		"Salaries": {32000, 18000, 9500, 10000, 7850},
	})
	assert.NoError(t, frame.Error())

	fmt.Println(frame)
	return frame
}

// 测试 `DataFrame` 分组
//
// 分组通过 `dataframe.GroupBy` 函数进行, 返回 `dataframe.Groups` 实例, 根据所给字段 (一到多个) 将整个 DataFrame 分为指定的组,
// 聚合通过 `dataframe.Groups` 实例的 `Aggregation` 方法进行, 用于对分组的指定列数据进行统计
func TestDataFrameGroup_GroupBy(t *testing.T) {
	frame := makeDFForGrouping(t)

	// 通过 `Gender` 字段进行分组
	gs := frame.GroupBy("Gender")
	assert.NoError(t, gs.Err)

	// 分组结果中包括 2 个分组
	gf := gs.GetGroups()
	assert.Len(t, gf, 2)

	// 获取名称为 `"F"` 的分组
	g, ok := gf["F"]
	assert.True(t, ok)
	assert.NoError(t, g.Error())

	// 确认 `"F"` 分组包含内容
	fmt.Println(g)
	assert.Equal(t, [][]any{
		{38, 51},         // Age
		{"F", "F"},       // Gender
		{"Emma", "Lucy"}, // Name
		{18000, 9500},    // Salaries
	}, AllColumnsVals(&g))

	// 获取名称为 `"M"` 的分组
	g, ok = gf["M"]
	assert.True(t, ok)
	assert.NoError(t, g.Error())

	// 确认 `"M"` 分组包含内容
	fmt.Println(g)
	assert.Equal(t, [][]any{
		{42, 34, 21},               // Age
		{"M", "M", "M"},            // Gender
		{"Alvin", "Tom", "Arthur"}, // Name
		{32000, 10000, 7850},       // Salaries
	}, AllColumnsVals(&g))
}

// 测试 `DataFrame` 的分组聚合
//
// 聚合方式包括:
//   - `dataframe.Aggregation_COUNT` 计数
//   - `dataframe.Aggregation_SUM` 求和
//   - `dataframe.Aggregation_MAX` 最大值
//   - `dataframe.Aggregation_MIN` 最小值
//   - `dataframe.Aggregation_MEDIAN` 中位数
//   - `dataframe.Aggregation_STD` 标准差
//   - `dataframe.Aggregation_MEAN` 平均值
func TestDataFrameGroup_Aggregation(t *testing.T) {
	frame := makeDFForGrouping(t)

	// 通过 `Gender` 字段进行分组
	gs := frame.GroupBy("Gender")
	assert.NoError(t, gs.Err)

	// 对根据 `Gender` 字段的分组结果进行聚合, 计算各分组 `Age` 和 `Salaries` 字段的各种统计结果
	// 需要为分组结果指定列名称
	agg := gs.Aggregation([]df.AggregationType{
		df.Aggregation_COUNT,  // Age
		df.Aggregation_SUM,    // Salaries
		df.Aggregation_MAX,    // Age
		df.Aggregation_MIN,    // Age
		df.Aggregation_MEDIAN, // Salaries
		df.Aggregation_STD,    // Salaries
		df.Aggregation_MEAN,   // Age
	}, []string{"Age", "Salaries", "Age", "Age", "Salaries", "Salaries", "Age"})

	assert.NoError(t, agg.Error())

	fmt.Println(agg)

	// 对分组结果进行过滤
	// 过滤 `Gender` 字段为 `"F"` 的聚合结果
	r := agg.Filter(df.F{Colname: "Gender", Comparator: series.Eq, Comparando: "F"})
	assert.NoError(t, r.Error())

	fmt.Println(r)

	// 确认 `Gender` 为 `"F"` 的统计结果
	m := r.Maps()
	assert.Len(t, m, 1)
	assert.Equal(t, 2.0, m[0]["Age_COUNT"])
	assert.Equal(t, 51.0, m[0]["Age_MAX"])
	assert.Equal(t, 38.0, m[0]["Age_MIN"])
	assert.Equal(t, 44.5, m[0]["Age_MEAN"])
	assert.Equal(t, 27500.0, m[0]["Salaries_SUM"])
	assert.Equal(t, 13750.0, m[0]["Salaries_MEDIAN"])
	assert.Equal(t, 6010.407640085654, m[0]["Salaries_STD"])

	// 过滤 `Gender` 字段为 `"M"` 的聚合结果
	mAgg := agg.Filter(df.F{Colname: "Gender", Comparator: series.Eq, Comparando: "M"})
	assert.NoError(t, agg.Error())

	fmt.Println(mAgg)

	// 确认 `Gender` 为 `"M"` 的统计结果
	m = r.Maps()
	assert.Len(t, m, 1)
	assert.Equal(t, 2.0, m[0]["Age_COUNT"])
	assert.Equal(t, 51.0, m[0]["Age_MAX"])
	assert.Equal(t, 38.0, m[0]["Age_MIN"])
	assert.Equal(t, 44.5, m[0]["Age_MEAN"])
	assert.Equal(t, 27500.0, m[0]["Salaries_SUM"])
	assert.Equal(t, 13750.0, m[0]["Salaries_MEDIAN"])
	assert.Equal(t, 6010.407640085654, m[0]["Salaries_STD"])
}
