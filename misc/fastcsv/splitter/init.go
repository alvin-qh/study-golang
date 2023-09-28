package splitter

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path"
	"syscall"

	"golang.org/x/exp/slices"
	"golang.org/x/sync/errgroup"
)

var (
	// 读缓存
	READ_BUF_SIZE = syscall.Getpagesize() * 128

	// 写缓存
	WRITE_BUF_SIZE = syscall.Getpagesize() * 128
)

var (
	// 分隔符
	SEP = []byte(",")

	// 换行符
	RET = []byte("\n")

	// BOM 头
	BOM = "\uFEFF"
)

// 定义一个分割器结构体
type Splitter struct {
	filename string       // 文件名
	outDir   string       // 输出文件路径
	columns  [][]byte     // 列集合
	records  [][][]byte   // 记录集合
	ci       *columnIndex // 必要列索引
}

// 打开输入 csv 文件
//
// 参数:
//   - `filename` (`string`): 要打开的输入 csv 文件名
//   - `distFolder` (`string`): 保存输出 csv 文件的目录
//   - `commonColumns` (`...string`): 输出 csv 文件中要包含的公共列
//
// 返回:
//   - `*Splitter`: 分隔器对象
//   - `error`: 错误对象
func Open(filename string, distFolder string, commonColumns ...string) (*Splitter, error) {
	// 如果目标目录不存在, 则创建该目录
	if _, err := os.Stat(distFolder); errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(distFolder, os.ModePerm); err != nil {
			return nil, err
		}
	}
	// 返回结构体对象
	return &Splitter{
		filename: filename,
		outDir:   distFolder,
		ci:       newColumnIndex(commonColumns),
	}, nil
}

// 合并数据
//
// 参数:
//   - `records` (`[][]byte`): 要写入的数据集合
//   - `record` (`[]byte`): 要写入的单个数据
//
// 返回:
//   - `[][]byte`: 参数合并后的结果
func mergeData(records [][]byte, record []byte) [][]byte {
	r := make([][]byte, 0, len(records)+1)
	return append(append(r, records...), record)
}

// 将读取的内容写入各自的输出文件
//
// 返回:
//   - `error`: 错误对象
func (c *Splitter) splitToFiles() error {
	// 错误组等待对象
	// 该对象用于等待所有的协程执行结束, 并在任意协程返回错误后, 通知其它协程中断
	// 如果任意协程返回错误, 则该对象等待结果为该错误
	g, ctx := errgroup.WithContext(context.Background())

	// 遍历所有列
	for i, col := range c.columns {
		// 如果列包含在公共列中, 则不输出文件
		if slices.Contains(c.ci.colIndex, i) {
			continue
		}

		// 记录本地变量用于协程闭包函数内部使用
		name := col
		index := i

		// 启动协程并加入等待组
		g.Go(func() error {
			// 为当前列创建输出文件
			w, err := newWriter(path.Join(c.outDir, fmt.Sprintf("%v.csv", string(name))))
			if err != nil {
				return err
			}
			defer w.Close()

			// 写入 csv 表头信息, 包括 公共列 + 当前列
			w.Write(mergeData(c.ci.colList, name)...)

			// 逐条遍历记录, 写入指定列的记录
			for _, row := range c.records {
				select {
				case <-ctx.Done(): // 查看协程是否中断
					return nil
				default:
				}

				// 写入 csv 内容数据, 包括 公共列 + 当前列
				w.Write(mergeData(c.ci.Records(row), row[index])...)
			}
			return nil
		})
	}
	// 等待协程
	return g.Wait()
}

// 执行 csv 文件分割
//
// 将所给的 csv 文件按列名分割为多个文件, 每个文件中包括 "公共列" + "任一列"
func (c *Splitter) Split() error {
	r, err := newReader(c.filename)
	if err != nil {
		return err
	}

	c.columns, c.records, err = r.ReadCSV()
	if err != nil {
		return err
	}

	if err := c.ci.Map(c.columns); err != nil {
		return err
	}

	if err := c.splitToFiles(); err != nil {
		return err
	}

	return nil
}
