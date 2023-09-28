package splitter

import (
	"bytes"
	"fmt"
)

// 表示列索引的结构体
// 该结构体用于存储必要列的列名和对应的索引
type columnIndex struct {
	colList  [][]byte
	colIndex []int
}

// 创建列索引结构体对象
//
// 参数:
//   - `necessaryCols` (`[]string`): 必要的列列名集合
//
// 返回:
//   - `*columnIndex`: 列索引结构体对象指针
func newColumnIndex(necessaryCols []string) *columnIndex {
	// 创建保存列名字节数组的集合
	cols := make([][]byte, len(necessaryCols))

	// 将传入的列名每一项转为字节数组后存储
	for n, c := range necessaryCols {
		cols[n] = []byte(c)
	}

	return &columnIndex{
		colList: cols,
	}
}

// 映射列索引
//
// 参数:
//   - `allColumns` (`[][]byte`): 所有列集合
//
// 返回:
//   - `error`: 错误信息
func (ci *columnIndex) Map(allColumns [][]byte) (err error) {
	// 创建列索引集合
	index := make([]int, 0, len(ci.colList))

	// 从所有列中查询所需列的索引并存储
	for _, nc := range ci.colList {
		found := false
		for n, ac := range allColumns {
			if bytes.Equal(nc, ac) {
				index = append(index, n)
				found = true
				break
			}
		}
		// 如果所需列不存在于所有列中, 则返回错误
		if !found {
			err = fmt.Errorf("column \"%v\" not exist", string(nc))
			return
		}
	}

	ci.colIndex = index
	return
}

// 从一行记录中筛选所需列的数据
//
// 参数:
//   - `row` (`[][]byte`): 表示一行的数据集合
//
// 返回:
//   - `[][]byte`: 指定索引的数据集合
func (ci *columnIndex) Records(row [][]byte) [][]byte {
	r := make([][]byte, len(ci.colIndex))
	for n, i := range ci.colIndex {
		r[n] = row[i]
	}
	return r
}
