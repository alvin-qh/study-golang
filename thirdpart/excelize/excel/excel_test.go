package excel

import (
	"strconv"
	"study/thirdpart/excelize/data"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xuri/excelize/v2"
)

// 测试将 `DataFrame` 实例写入 Excel 文件
//
// 将数据写入 Excel 文件的单元格, 需要两个参数:
//   - Sheet 页名称: 例如 `"Sheet1"`
//   - 单元格名称: Excel 用 `A`, `B`, ..., `AA`, `AB` 的字母索引表示列名, 用从 `1` 开始的数字索引表示行名, 组合起来就是单元格名称, 例如: `A1`, `AB12`
func TestWriteDataFrameToExcel(t *testing.T) {
	// 测试写入 Excel
	t.Run("WriteExcel", func(t *testing.T) {
		// 创建 Excel 实例
		f := excelize.NewFile()
		defer f.Close()

		// 获取第一个 Sheet 的名称
		sheet := f.GetSheetName(0)

		// 读取 JSON 文件, 获取 `DataFrame` 实例
		frame := data.LoadJSON("../files/simple_data.json")

		// 将 `DataFrame` 列名写入 Excel 的第一行
		for i, n := range frame.Names() {
			// 将列索引, 行索引转化为 Excel 单元格名称
			// 第 `i+1` 列, 第 `1` 行
			col, err := excelize.CoordinatesToCellName(i+1, 1)
			assert.NoError(t, err)

			// 将数据写入指定单元格
			err = f.SetCellValue(sheet, col, n)
			assert.NoError(t, err)
		}

		// 遍历 `DataFrame` 的行列, 将数据写入 Excel 文件
		for i := 0; i < frame.Nrow(); i++ {
			for j := 0; j < frame.Ncol(); j++ {
				// 根据行列计算单元格名称
				// 第 `j+1` 列, 第 `i+2` 行 (即从第二行开始写入数据)
				col, err := excelize.CoordinatesToCellName(j+1, i+2)
				assert.NoError(t, err)

				// 获取 `DataFrame` 中指定行列的元素
				elem := frame.Elem(i, j)
				// 将元素写入 Excel 对应单元格
				err = f.SetCellValue(sheet, col, elem.Val())
				assert.NoError(t, err)
			}
		}

		// 将 Excel 写入指定的文件
		err := f.SaveAs("../files/01.xlsx")
		assert.NoError(t, err)
	})

	// 测试读取写入的 Excel
	t.Run("ReadExcel", func(t *testing.T) {
		// 打开 Excel 文件
		f, err := excelize.OpenFile("../files/01.xlsx")
		assert.NoError(t, err)

		defer f.Close()

		// 获取第一个 Sheet 的名称
		sheet := f.GetSheetName(0)

		// 获取行集对象
		rows, err := f.Rows(sheet)
		assert.NoError(t, err)

		// 读取第一行
		assert.True(t, rows.Next())

		// 获取第一行的所有列, 即标题列
		names, err := rows.Columns()
		assert.NoError(t, err)

		// 确认获取的标题列
		assert.Equal(t, []string{"A", "B", "C", "D", "E"}, names)

		// 读取 JSON 文件, 获取 `DataFrame` 实例
		frame := data.LoadJSON("../files/simple_data.json")

		// 遍历剩余的行
		for i := 0; rows.Next(); i++ {
			// 获取当前行的所有列
			cols, err := rows.Columns()
			assert.NoError(t, err)

			// 遍历列数据
			for j, col := range cols {
				// 从 `DataFrame` 获取对应数据
				elem := frame.Elem(i, j)

				// 确认 Excel 当前行的每列值
				val, err := strconv.ParseFloat(col, 64)
				assert.NoError(t, err)
				assert.Equal(t, elem.Val(), val)
			}
		}
	})
}
