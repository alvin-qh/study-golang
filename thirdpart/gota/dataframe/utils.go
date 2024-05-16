package dataframe

import (
	df "github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
)

// 将一个 DataFrame 的所有元素读取为列的切片
//
// 列切片的内容为
//
//	[
//	  [第 1 列...],
//	  [第 2 列...],
//	  [第 3 列...],
//	  ...
//	]
func AllColumnsVals(df *df.DataFrame) [][]any {
	vals := make([][]any, 0, df.Ncol())

	for i := 0; i < df.Ncol(); i++ {
		col := make([]any, 0, df.Nrow())
		for j := 0; j < df.Nrow(); j++ {
			elem := df.Elem(j, i)
			col = append(col, elem.Val())
		}
		vals = append(vals, col)
	}

	return vals
}

// 将一个 DataFrame 的所有元素读取为行的切片
//
// 行切片的内容为
//
//	[
//	  [第 1 行...],
//	  [第 2 行...],
//	  [第 3 行...],
//	  ...
//	]
func AllRowsVals(df *df.DataFrame) [][]any {
	vals := make([][]any, 0, df.Ncol())

	for i := 0; i < df.Nrow(); i++ {
		col := make([]any, 0, df.Nrow())
		for j := 0; j < df.Ncol(); j++ {
			elem := df.Elem(i, j)
			col = append(col, elem.Val())
		}
		vals = append(vals, col)
	}

	return vals
}

// 获取指定名称或索引的整列数据
//
// `columns` 参数为一个整数或字符串组成的不定参数, 可以表示列名称或列序号, 例如
//
//	GetColumnsVals(df, "A", 1, "C")
func GetColumnsVals(frame *df.DataFrame, columns ...any) []series.Series {
	// 存储结果的切片
	cols := make([]series.Series, 0, len(columns))

	// 获取所有列名称
	colNames := frame.Names()

	// 遍历要获取的所有列
	for _, column := range columns {
		// 根据参数类型, 获取列数据
		switch c := column.(type) {
		case string:
			// 对于参数类型为字符串, 则直接使用列名称获取列数据
			cols = append(cols, frame.Col(c))
		case int:
			// 对于参数类型为整数, 则使用列序号获取列名称, 再获取列数据
			cols = append(cols, frame.Col(colNames[c]))
		}
	}
	return cols
}

// 获取指定索引的整行数据
func GetRowsVals(frame *df.DataFrame, rows ...int) [][]any {
	// 存储结果的切片
	rs := make([][]any, 0, len(rows))

	// 遍历所有指定行
	for _, row := range rows {
		// 生成保存行数据的切片
		r := make([]any, frame.Ncol())
		// 将当前行数据存入切片
		for col := 0; col < frame.Ncol(); col++ {
			r[col] = frame.Elem(row, col).Val()
		}
		rs = append(rs, r)
	}
	return rs
}

// 选择指定列中指定值所在的整行数据
//
// 通过 `column` 参数表示的列名称或列序号找到指定列, 通过 `value` 参数在该列中找到指定行, 返回该行数据
func SelectRows(frame *df.DataFrame, column any, value any) *df.DataFrame {
	// 根据列名或列序号找到整列数据
	col := GetColumnsVals(frame, column)[0]

	// 获取所有列名称
	colNames := frame.Names()

	// 保存结果的切片
	ses := make([]series.Series, 0)

	// 遍历所找到列的每一行
	for row := 0; row < col.Len(); row++ {
		// 如果某行匹配到指定的值, 表示找到指定行
		if col.Elem(row).Val() == value {
			// 遍历该行的所有列
			for col := 0; col < frame.Ncol(); col++ {
				// 保存指定列的指定行数据
				elem := frame.Elem(row, col)
				ses = append(ses, series.New(elem.Val(), elem.Type(), colNames[col]))
			}
			break
		}
	}

	rs := df.New(ses...)
	return &rs
}

// 求 `map[string][]any` 类型 Map 中所有值的最大长度
func maxLenOfMapValues(m map[string][]any) int {
	max := 0
	for _, vs := range m {
		if len(vs) > max {
			max = len(vs)
		}
	}
	return max
}

// 将 `map[string][]any` 类型 Map 转换为 `[]map[string]any` 类型 Map
//
// 例如:
//
//	map[string][]any{"A": {1, 2, 3}, "B": {4, 5, 6}} => []map[string]any{{A: 1, B: 4}, {A: 2, B: 5}, {A: 3, B: 6}}
//
// 后者可作为产生 `DataFrame` 的参数使用
func MapConvert(m map[string][]any) []map[string]any {
	max := maxLenOfMapValues(m)

	r := make([]map[string]any, max)
	for k, vs := range m {
		for i, v := range vs {
			if r[i] == nil {
				r[i] = make(map[string]any)
			}
			r[i][k] = v
		}
	}
	return r
}

// 将 `map[string][]any` 类型 Map 转换为 `DataFrame` 实例
func LoadSliceMap(m map[string][]any) df.DataFrame {
	return df.LoadMaps(MapConvert(m))
}
