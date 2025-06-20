package dataframe

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"thirdpart/gota/utils"

	df "github.com/go-gota/gota/dataframe"
	"github.com/stretchr/testify/assert"
)

// 通过 Map 创建 `DataFrame` 实例
func makeDFForWrite(t *testing.T) df.DataFrame {
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

// 测试将 `DataFrame` 实例内容写入 CSV 文件
func TestWriteDataFrame_WriteCSV(t *testing.T) {
	w := utils.WithNoError(t, func() (*os.File, error) { return os.Create("data.csv") })
	defer os.Remove(w.Name())

	fw := makeDFForWrite(t)

	// 将 `DataFrame` 写入 CSV 文件
	err := fw.WriteCSV(bufio.NewWriter(w))
	assert.NoError(t, err)

	w.Close()

	r := utils.WithNoError(t, func() (*os.File, error) { return os.Open("data.csv") })

	fr := df.ReadCSV(bufio.NewReader(r))
	assert.NoError(t, fr.Error())
	fmt.Println(fr)

	assert.Equal(t, fw, fr)

	r.Close()
}

// 测试将 `DataFrame` 实例内容写入 JSON 文件
//
// JSON 文件内容为一个 Object 数组, 数组每一项为 `DataFrame` 一行数据
func TestWriteDataFrame_WriteJSON(t *testing.T) {
	w := utils.WithNoError(t, func() (*os.File, error) { return os.Create("data.json") })
	defer os.Remove(w.Name())

	fw := makeDFForWrite(t)

	// 将 `DataFrame` 写入 JSON 文件
	err := fw.WriteJSON(w)
	assert.NoError(t, err)

	w.Close()

	r := utils.WithNoError(t, func() (*os.File, error) { return os.Open("data.json") })

	fr := df.ReadJSON(r)
	assert.NoError(t, fr.Error())
	fmt.Println(fr)

	assert.Equal(t, fw, fr)

	r.Close()
}
