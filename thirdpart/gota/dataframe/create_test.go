package dataframe

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	df "github.com/go-gota/gota/dataframe"
	"github.com/stretchr/testify/assert"
	"gonum.org/v1/gonum/mat"
)

// 测试从 CSV 文件中读取内容并生成 DataFrame 实例
func TestCreateDataFrame_ReadCSV(t *testing.T) {
	// 读取文件, 获取 io 实例
	f, err := os.Open("../data/data.csv")
	assert.NoError(t, err)

	// 通过文件 io 实例读取 CSV 内容, 生成 DataFrame 实例
	frame := df.ReadCSV(f, df.DetectTypes(true), df.HasHeader(true))
	assert.NoError(t, frame.Error())

	fmt.Println(frame)

	// CSV 中包含 5 行 5 列数据
	assert.Equal(t, 5, frame.Ncol())
	assert.Equal(t, 5, frame.Nrow())

	// 读取 DataFrame 中的每个元素, 返回一个 float64 类型二维切片
	vals := AllColumnsVals(&frame)
	assert.Equal(t, [][]any{
		{0.1, 1.1, 2.1, 3.1, 4.1}, // A
		{0.2, 1.2, 2.2, 3.2, 4.2}, // B
		{0.3, 1.3, 2.3, 3.3, 4.3}, // C
		{0.4, 1.4, 2.4, 3.4, 4.4}, // D
		{0.5, 1.5, 2.5, 3.5, 4.5}, // E
	}, vals)
}

// 测试从 JSON 字符串中读取内容并生成 DataFrame 实例
//
// 要求: JSON 内容为一个数组, 数组的每项为一个字典, Key 为列名, Value 为列值, 每个字典表示一行数据
func TestCreateDataFrame_ReadJSON(t *testing.T) {
	jsonStr := `[
        {
            "A": 0.1,
            "B": 0.2,
            "C": 0.3,
            "D": 0.4,
            "E": 0.5
        },
        {
            "A": 1.1,
            "B": 1.2,
            "C": 1.3,
            "D": 1.4,
            "E": 1.5
        },
        {
            "A": 2.1,
            "B": 2.2,
            "C": 2.3,
            "D": 2.4,
            "E": 2.5
        },
        {
            "A": 3.1,
            "B": 3.2,
            "C": 3.3,
            "D": 3.4,
            "E": 3.5
        },
        {
            "A": 4.1,
            "B": 4.2,
            "C": 4.3,
            "D": 4.4,
            "E": 4.5
        }
    ]`

	// 将字符串转为字节串并产生一个内存缓冲
	buf := bytes.NewBuffer([]byte(jsonStr))

	// 读取缓冲中的 JSON 内容, 生成 DataFrame 实例
	frame := df.ReadJSON(buf, df.DetectTypes(true))
	assert.NoError(t, frame.Error())

	fmt.Println(frame)

	// JSON 中包含 5 行 5 列数据
	assert.Equal(t, 5, frame.Ncol())
	assert.Equal(t, 5, frame.Nrow())

	// 读取 DataFrame 中的每个元素, 返回一个 float64 类型二维切片
	vals := AllColumnsVals(&frame)
	assert.Equal(t, [][]any{
		{0.1, 1.1, 2.1, 3.1, 4.1}, // A
		{0.2, 1.2, 2.2, 3.2, 4.2}, // B
		{0.3, 1.3, 2.3, 3.3, 4.3}, // C
		{0.4, 1.4, 2.4, 3.4, 4.4}, // D
		{0.5, 1.5, 2.5, 3.5, 4.5}, // E
	}, vals)
}

// 测试从 HTML 中读取 DataFrame 实例
//
// 如果 HTML 中包含 `<table>` 元素, 则获取其中 `<tbody>` 元素内的数据, 生成 `DataFrame` 实例
//
// 一个 HTML 中可以包含多个 `<table>` 元素, 故 `dataframe.ReadHTML` 函数返回的是一个 `DataFrame` 实例的集合
//
// 注意: 要获取的数据必须包含在 `<table>` 或 `<tbody>` 元素中, 包括表头内容, 所以如果表头在 `<thead>` 元素中, 则
// 读取 HTML 时需加上 `HasHeader(false)` 选型, 即读取表格时不包含表头, 否则表格的第一行数据会被忽略
func TestCreateDataFrame_ReadHTML(t *testing.T) {
	// 读取文件, 获取 io 实例
	f, err := os.Open("../data/data.html")
	assert.NoError(t, err)

	// 通过文件 io 实例读取 CSV 内容, 生成 DataFrame 实例
	frames := df.ReadHTML(f, df.DetectTypes(true), df.HasHeader(false))
	assert.Len(t, frames, 1)

	frame := frames[0]
	assert.NoError(t, frame.Error())

	fmt.Println(frame)

	// HTML 表格中包含 5 行 5 列数据
	assert.Equal(t, 5, frame.Ncol())
	assert.Equal(t, 5, frame.Nrow())

	// 读取 DataFrame 中的每个元素, 返回一个 float64 类型二维切片
	vals := AllColumnsVals(&frame)
	assert.Equal(t, [][]any{
		{0.1, 1.1, 2.1, 3.1, 4.1}, // A
		{0.2, 1.2, 2.2, 3.2, 4.2}, // B
		{0.3, 1.3, 2.3, 3.3, 4.3}, // C
		{0.4, 1.4, 2.4, 3.4, 4.4}, // D
		{0.5, 1.5, 2.5, 3.5, 4.5}, // E
	}, vals)
}

// 测试从 Map 集合结构创建 DataFrame 实例
//
// 要求 Map 结构中的 Key 为列名, Value 为列值, 每个 Map 表示一行数据
func TestCreateDataFrame_LoadMaps(t *testing.T) {
	ms := []map[string]any{
		{
			"A": 0.1,
			"B": 0.2,
			"C": 0.3,
			"D": 0.4,
			"E": 0.5,
		},
		{
			"A": 1.1,
			"B": 1.2,
			"C": 1.3,
			"D": 1.4,
			"E": 1.5,
		},
		{
			"A": 2.1,
			"B": 2.2,
			"C": 2.3,
			"D": 2.4,
			"E": 2.5,
		},
		{
			"A": 3.1,
			"B": 3.2,
			"C": 3.3,
			"D": 3.4,
			"E": 3.5,
		},
		{
			"A": 4.1,
			"B": 4.2,
			"C": 4.3,
			"D": 4.4,
			"E": 4.5,
		},
	}

	// 从 Map 结构中读取 `DataFrame` 实例
	frame := df.LoadMaps(ms, df.DetectTypes(true))
	assert.NoError(t, frame.Error())

	fmt.Println(frame)

	// Map 中包含 5 行 5 列数据
	assert.Equal(t, 5, frame.Ncol())
	assert.Equal(t, 5, frame.Nrow())

	// 读取 DataFrame 中的每个元素, 返回一个 float64 类型二维切片
	vals := AllColumnsVals(&frame)
	assert.Equal(t, [][]any{
		{0.1, 1.1, 2.1, 3.1, 4.1}, // A
		{0.2, 1.2, 2.2, 3.2, 4.2}, // B
		{0.3, 1.3, 2.3, 3.3, 4.3}, // C
		{0.4, 1.4, 2.4, 3.4, 4.4}, // D
		{0.5, 1.5, 2.5, 3.5, 4.5}, // E
	}, vals)
}

// 结构体中的字段可作为 DataFrame 的一个元素, 其中:
//
//  1. 字段名为 DataFrame 的列名
//  2. 字段类型为作为 DataFrame 的元素的类型
//  3. 如果要忽略一个字段, 则:
//     - 为该字段添加 `dataframe:"-"` 标签;
//     - 将字段名改为小写, 即私有字段;
//  4. 可以通过 `dataframe:"名称,类型"` 来显式指定字段名和字段类型
type column struct {
	A float64 `dataframe:"A,float64"`
	B float64 `dataframe:"A,float64"`
	C float64 `dataframe:"A,float64"`
	D float64 `dataframe:"A,float64"`
	E float64 `dataframe:"A,float64"`
}

// 测试从结构体集合中创建 DataFrame 实例
//
// 每个结构体表示一行数据
func TestCreateDataFrame_LoadStructs(t *testing.T) {
	s := []column{
		{
			A: 0.1,
			B: 0.2,
			C: 0.3,
			D: 0.4,
			E: 0.5,
		},
		{
			A: 1.1,
			B: 1.2,
			C: 1.3,
			D: 1.4,
			E: 1.5,
		},
		{
			A: 2.1,
			B: 2.2,
			C: 2.3,
			D: 2.4,
			E: 2.5,
		},
		{
			A: 3.1,
			B: 3.2,
			C: 3.3,
			D: 3.4,
			E: 3.5,
		},
		{
			A: 4.1,
			B: 4.2,
			C: 4.3,
			D: 4.4,
			E: 4.5,
		},
	}

	// 从结构体集合中创建 DataFrame 实例
	frame := df.LoadStructs(s)
	assert.NoError(t, frame.Error())

	fmt.Println(frame)

	// 结构体中包含 5 行 5 列数据
	assert.Equal(t, 5, frame.Ncol())
	assert.Equal(t, 5, frame.Nrow())

	// 读取 DataFrame 中的每个元素, 返回一个 float64 类型二维切片
	vals := AllColumnsVals(&frame)
	assert.Equal(t, [][]any{
		{0.1, 1.1, 2.1, 3.1, 4.1}, // A
		{0.2, 1.2, 2.2, 3.2, 4.2}, // B
		{0.3, 1.3, 2.3, 3.3, 4.3}, // C
		{0.4, 1.4, 2.4, 3.4, 4.4}, // D
		{0.5, 1.5, 2.5, 3.5, 4.5}, // E
	}, vals)
}

// 测试通过矩阵实例创建 `DataFrame` 实例
func TestCreateDataFrame_LoadMatrix(t *testing.T) {
	// 创建一个 5x5 的矩阵 (5r, 5c)
	m := mat.NewDense(5, 5, []float64{
		0.1, 0.2, 0.3, 0.4, 0.5,
		1.1, 1.2, 1.3, 1.4, 1.5,
		2.1, 2.2, 2.3, 2.4, 2.5,
		3.1, 3.2, 3.3, 3.4, 3.5,
		4.1, 4.2, 4.3, 4.4, 4.5,
	})

	// 从矩阵中创建 DataFrame 实例
	frame := df.LoadMatrix(m)
	assert.NoError(t, frame.Error())

	frame.SetNames("A", "B", "C", "D", "E")
	assert.Equal(t, []string{"A", "B", "C", "D", "E"}, frame.Names())

	fmt.Println(frame)

	// 矩阵中包含 5 行 5 列数据
	assert.Equal(t, 5, frame.Ncol())
	assert.Equal(t, 5, frame.Nrow())

	// 读取 DataFrame 中的每个元素, 返回一个 float64 类型二维切片
	vals := AllColumnsVals(&frame)
	assert.Equal(t, [][]any{
		{0.1, 1.1, 2.1, 3.1, 4.1}, // A
		{0.2, 1.2, 2.2, 3.2, 4.2}, // B
		{0.3, 1.3, 2.3, 3.3, 4.3}, // C
		{0.4, 1.4, 2.4, 3.4, 4.4}, // D
		{0.5, 1.5, 2.5, 3.5, 4.5}, // E
	}, vals)
}

// 通过 Map 创建 `DataFrame` 实例
func makeDFForCreate(t *testing.T) df.DataFrame {
	frame := LoadSliceMap(map[string][]any{
		"A": {0.1, 1.1, 2.1, 3.1, 4.1},
		"B": {0.2, 1.2, 2.2, 3.2, 4.2},
		"C": {0.3, 1.3, 2.3, 3.3, 4.3},
		"D": {0.4, 1.4, 2.4, 3.4, 4.4},
		"E": {0.5, 1.5, 2.5, 3.5, 4.5},
	})

	assert.NoError(t, frame.Error())

	fmt.Println(frame)
	return frame
}

// 测试复制 `DataFrame` 实例
//
// 通过 `Copy` 方法复制一个 `DataFrame` 实例, 结果和原 `DataFrame` 完全一致
func TestCreateDataFrame_Copy(t *testing.T) {
	frame := makeDFForCreate(t)

	copiedFrame := frame.Copy()
	assert.NoError(t, copiedFrame.Error())

	fmt.Println(copiedFrame)

	assert.Equal(t, frame.Types(), copiedFrame.Types())
	assert.Equal(t, frame.Ncol(), copiedFrame.Ncol())
	assert.Equal(t, frame.Nrow(), copiedFrame.Nrow())
	assert.Equal(t, frame.Names(), copiedFrame.Names())
	assert.Equal(t, frame.Records(), copiedFrame.Records())
}
