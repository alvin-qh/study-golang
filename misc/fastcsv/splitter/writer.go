package splitter

import (
	"bufio"
	"bytes"
	"os"
)

// 文件写入器结构体
type writer struct {
	w *bufio.Writer
	f *os.File
}

// 创建写入器结构体对象
//
// 参数:
//   - `filename` (`string`): 要写入的文件名称
//
// 返回:
//   - `*writer`: 写入器结构体对象
//   - `error`: 错误对象
func newWriter(filename string) (*writer, error) {
	// 创建文件
	file, err := os.Create(filename)
	if err != nil {
		return nil, err
	}

	// 创建写入文件的缓冲写入器对象
	w := bufio.NewWriterSize(file, WRITE_BUF_SIZE)
	return &writer{
		w: w,
		f: file,
	}, nil
}

// 写入记录
//
// 参数:
//   - `records` (`...[]byte`): 要写入的记录集合
//
// 返回:
//   - `error`: 错误对象
func (w *writer) Write(records ...[]byte) error {
	// 写入一行内容, 内容通过指定的分隔符分隔
	if _, err := w.w.Write(bytes.Join(records, SEP)); err != nil {
		return err
	}
	// 写入换行符合
	_, err := w.w.Write(RET)
	return err
}

// 关闭写入器对象
func (w *writer) Close() {
	w.w.Flush()
	w.f.Close()
}
