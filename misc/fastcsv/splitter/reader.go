package splitter

import (
	"bufio"
	"bytes"
	"io"
	"os"
)

// 文件读取器结构体
type reader struct {
	f *os.File
	r *bufio.Reader
}

// 创建文件读取结构体
//
// 参数:
//   - `filename` (`string`): 要读取的文件名
//
// 返回:
//   - `*reader`: `reader` 结构体指针
//   - `error`: 错误对象
func newReader(filename string) (*reader, error) {
	// 打开文件
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	// 通过预设的缓存大小创建文件读取器对象
	r := bufio.NewReaderSize(f, READ_BUF_SIZE)
	return &reader{
		f: f,
		r: r,
	}, nil
}

// 关闭 `reader` 读取器
//
// 实现 `io.Closer` 接口
func (r *reader) Close() {
	r.f.Close()
}

// 读取一行内容
//
// 所谓一行内容, 即两个 `\n` 字符之间的内容
//
// 返回:
//   - `[]byte`: 表示文件一行内容的
//   - `error`: 错误对象
func (r *reader) readLine() (line []byte, err error) {
	// 通过 `bufio.Reader` 读取一行数据
	tmpLine, isprefix, err := r.r.ReadLine()
	for isprefix && err == nil {
		// 如果读缓冲区满但尚未读完一行, 则继续读后续内容
		var bs []byte
		bs, isprefix, err = r.r.ReadLine()

		// 将读取内容追加到当前行内容上
		tmpLine = append(line, bs...)
	}
	if err != nil && err != io.EOF {
		return
	}

	// 将读取的内容进行复制
	line = make([]byte, len(tmpLine))
	copy(line, tmpLine)
	return
}

// 读取 csv 文件
//
// 返回:
//   - `columns` (`[][]byte`): csv 的表头 (第一行)
//   - `records` (`[][][]byte`): csv 的表头 (第一行)
//   - `err` ([][]byte): csv 的表头 (第一行)
func (r *reader) ReadCSV() (columns [][]byte, records [][][]byte, err error) {
	// 读取第一行内容 (表头)
	column, err := r.readLine()
	if err != nil {
		return
	}
	// 将 BOM 标识取消掉
	column = bytes.TrimLeft(column, BOM)
	// 通过分隔符将内容进行分割
	columns = bytes.Split(column, SEP)

	// 创建存储内容的集合
	records = make([][][]byte, 0, 5000)
	for {
		// 读取一行
		line, err := r.readLine()
		if err != nil {
			break
		}
		// 将内容分割后加入到集合
		records = append(records, bytes.Split(line, SEP))
	}
	return
}
