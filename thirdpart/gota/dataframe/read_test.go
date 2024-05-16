package dataframe

import (
	"fmt"
	"testing"

	df "github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"github.com/stretchr/testify/assert"
)

// 通过 Map 创建 `DataFrame` 实例
func makeDFForRead(t *testing.T) (map[string][]any, df.DataFrame) {
	dm := map[string][]any{
		"A": {0.1, 1.1, 2.1, 3.1, 4.1},
		"B": {0.2, 1.2, 2.2, 3.2, 4.2},
		"C": {0.3, 1.3, 2.3, 3.3, 4.3},
		"D": {0.4, 1.4, 2.4, 3.4, 4.4},
		"E": {0.5, 1.5, 2.5, 3.5, 4.5},
	}

	frame := LoadSliceMap(dm)

	assert.NoError(t, frame.Error())

	fmt.Println(frame)
	return dm, frame
}

func TestReadDataFrame_Basic(t *testing.T) {
	_, frame := makeDFForRead(t)

	// 获取列名
	ns := frame.Names()
	assert.Equal(t, []string{"A", "B", "C", "D", "E"}, ns)

	// 获取行数和列数
	assert.Equal(t, 5, frame.Ncol())
	assert.Equal(t, 5, frame.Nrow())

	// 获取 DataFrame 的秩, 即行数和列数组成的元组
	r, c := frame.Dims()
	assert.Equal(t, 5, r)
	assert.Equal(t, 5, c)
}

// 测试以字符串二维切片形式获取 DataFrame 中的所有数据
//
// 通过 `DataFrame` 实例的 `Records` 可以获取一个字符串类型二维切片, 以行为单位获取 `DataFrame` 元素值
func TestReadDataFrame_Record(t *testing.T) {
	_, frame := makeDFForRead(t)

	// 以字符串二维切片形式获取整个 DataFrame,
	// 二维切片的每一项为 DataFrame 中的一行数据
	records := frame.Records()
	assert.Equal(t, [][]string{
		{"A", "B", "C", "D", "E"},                                    // Names 行
		{"0.100000", "0.200000", "0.300000", "0.400000", "0.500000"}, // 第一行
		{"1.100000", "1.200000", "1.300000", "1.400000", "1.500000"}, // 第二行
		{"2.100000", "2.200000", "2.300000", "2.400000", "2.500000"}, // 第三行
		{"3.100000", "3.200000", "3.300000", "3.400000", "3.500000"}, // 第四行
		{"4.100000", "4.200000", "4.300000", "4.400000", "4.500000"}, // 第五行
	}, records)
}

// 测试获取 DataFrame 中指定列数据
//
// 可通过 `Col` 方法根据列名称返回一个 `Series` 实例, 其中包含了指定列的数据
func TestReadDataFrame_Col(t *testing.T) {
	_, frame := makeDFForRead(t)

	// 获取指定名称的列数据
	s := frame.Col("A")
	assert.NoError(t, s.Error())
	assert.Equal(t, []float64{0.1, 1.1, 2.1, 3.1, 4.1}, s.Float())

	s = frame.Col("E")
	assert.NoError(t, s.Error())
	assert.Equal(t, []float64{0.5, 1.5, 2.5, 3.5, 4.5}, s.Float())
}

// 测试通过行列元素获取 `DataFrame` 中指定位置的元素
func TestReadDataFrame_Elem(t *testing.T) {
	dm, frame := makeDFForRead(t)

	// 获取 DataFrame 每个元素, 即通过行列索引获取指定位置元素
	for col := 0; col < frame.Ncol(); col++ {
		colName := frame.Names()[col]

		for row := 0; row < frame.Nrow(); row++ {
			// 获取指定位置元素
			elem := frame.Elem(row, col)
			// 以浮点类型获取元素值
			assert.Equal(t, dm[colName][row], elem.Float())
		}
	}
}

// 计算 `DataFrame` 中每一列的统计值
//
// 可以为 `DataFrame` 的每一列计算平均值 (`mean`), 中值 (`median`), 标准差 (`stddev`), 最小值 (`min`),
// 25% 分位数 (`25%`), 50% 分位数 (`50%`), 75% 分位数 (`75%`), 最大值  (`max`)
func TestReadDataFrame_Describe(t *testing.T) {
	_, frame := makeDFForRead(t)

	// 获取 DataFrame 每一列的的统计信息
	describ := frame.Describe()
	assert.NoError(t, describ.Error())

	fmt.Println(describ)

	assert.Equal(t, [][]any{
		{"mean", 2.1, 2.2, 2.3, 2.4, 2.5},   // 每列平均值
		{"median", 2.1, 2.2, 2.3, 2.4, 2.5}, // 每列中值
		{"stddev", 1.5811388300841895, 1.5811388300841898, 1.5811388300841895, 1.5811388300841898, 1.5811388300841898}, // 每列标准差
		{"min", 0.1, 0.2, 0.3, 0.4, 0.5}, // 每列最小值
		{"25%", 1.1, 1.2, 1.3, 1.4, 1.5}, // 每列 25% 分位数
		{"50%", 2.1, 2.2, 2.3, 2.4, 2.5}, // 每列 50% 分位数
		{"75%", 3.1, 3.2, 3.3, 3.4, 3.5}, // 每列 75% 分位数
		{"max", 4.1, 4.2, 4.3, 4.4, 4.5}, // 每列最大值
	}, AllRowsVals(&describ))
}

// 测试根据列索引获取 `DataFrame` 中指定列的数据
//
// 可通过 `Select` 方法根据列索引返回一个 `DataFrame` 实例, 其中包含了指定列的数据
//
// 可以通过整数, 整数切片, 字符串, 字符串切片来表示列索引, 获取对应列的数据
func TestReadDataFrame_Select(t *testing.T) {
	_, frame := makeDFForRead(t)

	// 通过列索引获取指定列的内容, 返回包含所选列的 DataFrame 实例
	rs := frame.Select(0)
	assert.NoError(t, rs.Error())

	fmt.Println(rs)

	assert.Equal(t, [][]any{
		{0.1, 1.1, 2.1, 3.1, 4.1}, // A
	}, AllColumnsVals(&rs))

	// 通过列名称获取指定列的内容, 返回包含所选列的 DataFrame 实例
	rs = frame.Select("B")
	assert.NoError(t, rs.Error())

	fmt.Println(rs)

	assert.Equal(t, [][]any{
		{0.2, 1.2, 2.2, 3.2, 4.2}, // B
	}, AllColumnsVals(&rs))

	// 通过列索引获取若干列内容, 返回包含所选列的 DataFrame 实例
	rs = frame.Select([]int{0, 2, 4})
	assert.NoError(t, rs.Error())

	fmt.Println(rs)

	assert.Equal(t, [][]any{
		{0.1, 1.1, 2.1, 3.1, 4.1}, // A
		{0.3, 1.3, 2.3, 3.3, 4.3}, // C
		{0.5, 1.5, 2.5, 3.5, 4.5}, // E
	}, AllColumnsVals(&rs))

	// 通过列名称获取若干列内容, 返回包含所选列的 DataFrame 实例
	rs = frame.Select([]string{"A", "C", "E"})
	assert.NoError(t, rs.Error())

	fmt.Println(rs)

	assert.Equal(t, [][]any{
		{0.1, 1.1, 2.1, 3.1, 4.1}, // A
		{0.3, 1.3, 2.3, 3.3, 4.3}, // C
		{0.5, 1.5, 2.5, 3.5, 4.5}, // E
	}, AllColumnsVals(&rs))
}

// 通过条件对数据进行过滤, 获取过滤后的结果组成的 DataFrame
//
// 过滤需要指定条件, 即 `dataframe.F` 结构体实例, 其中:
//
// `Colname`, `Colidx` 指定过滤列名称或索引
//
// `Comparator` 指定过滤条件, 过滤条件包括:
//   - `series.Eq` 指定等于条件
//   - `series.Neq` 指定不等于条件
//   - `series.Greater` 指定大于条件
//   - `series.GreaterEq` 指定大于等于条件
//   - `series.Less` 指定小于条件
//   - `series.LessEq` 指定小于等于条件
//   - `series.In` 指定包含在列表中的条件
//
// `Comparando` 指定过滤条件的值
//
// 多次调用 `Filter` 方法相当于多个条件进行 `AND` 组合
//
// 过滤结果为包含过滤后数据的 `DataFrame` 实例
func TestReadDataFrame_Filter(t *testing.T) {
	_, frame := makeDFForRead(t)

	// 通过多次调用 `Filter` 方法相当于多个条件进行 `AND` 组合
	rs := frame.Filter(
		df.F{Colname: "A", Comparator: series.GreaterEq, Comparando: 2.0},
	).Filter(
		df.F{Colname: "E", Comparator: series.Less, Comparando: 4},
	)
	assert.NoError(t, rs.Error())

	fmt.Println(rs)

	// 确认过滤结果
	assert.Equal(t, [][]any{
		{2.1, 3.1}, // A
		{2.2, 3.2}, // B
		{2.3, 3.3}, // C
		{2.4, 3.4}, // D
		{2.5, 3.5}, // E
	}, AllColumnsVals(&rs))

	// 在一个 `Filter` 方法中传入多个条件, 相当于多个条件的 `OR` 组合
	rs = frame.Filter(
		df.F{Colname: "A", Comparator: series.Eq, Comparando: 2.1},
		df.F{Colname: "B", Comparator: series.Greater, Comparando: 4},
	)
	assert.NoError(t, rs.Error())

	fmt.Println(rs)

	assert.Equal(t, [][]any{
		{2.1, 4.1}, // A
		{2.2, 4.2}, // B
		{2.3, 4.3}, // C
		{2.4, 4.4}, // D
		{2.5, 4.5}, // E
	}, AllColumnsVals(&rs))
}

// 测试对 `DataFrame` 中指定列进行聚合排序
//
// 所谓聚合过滤, 即可以通过指定的逻辑运算符 (`And`, `Or`) 组合多个过滤条件
//
// 过滤结果为包含过滤后数据的 `DataFrame` 实例
func TestReadDataFrame_FilterAggregation(t *testing.T) {
	_, frame := makeDFForRead(t)

	// 也可以通过 `FilterAggregation` 方法和 `AND` 操作符明确进行 `AND` 操作
	rs := frame.FilterAggregation(
		df.And,
		df.F{Colname: "A", Comparator: series.GreaterEq, Comparando: 2.0},
		df.F{Colname: "E", Comparator: series.Less, Comparando: 4},
	)
	assert.NoError(t, rs.Error())

	fmt.Println(rs)

	assert.Equal(t, [][]any{
		{2.1, 3.1}, // A
		{2.2, 3.2}, // B
		{2.3, 3.3}, // C
		{2.4, 3.4}, // D
		{2.5, 3.5}, // E
	}, AllColumnsVals(&rs))

	// 也可以通过 `FilterAggregation` 方法和 `OR` 操作符明确进行 `OR` 操作
	rs = frame.FilterAggregation(
		df.Or,
		df.F{Colname: "A", Comparator: series.Eq, Comparando: 2.1},
		df.F{Colname: "B", Comparator: series.Greater, Comparando: 4},
	)
	assert.NoError(t, rs.Error())

	fmt.Println(rs)

	assert.Equal(t, [][]any{
		{2.1, 4.1}, // A
		{2.2, 4.2}, // B
		{2.3, 4.3}, // C
		{2.4, 4.4}, // D
		{2.5, 4.5}, // E
	}, AllColumnsVals(&rs))
}
