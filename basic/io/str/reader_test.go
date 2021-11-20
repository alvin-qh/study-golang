package str

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	STRING_READER_DATA = `<html>
    <body>
        <div>Hello World</div>
    </body>
</html>`
)

// 测试 strings.Reader
// strings.Reader 接收一个字符串参数，返回一个 Reader 对象，该对象实现了如下接口：
//  io.Reader, io.ReaderAt, io.ByteReader, io.ByteScanner, io.RuneReader, io.RuneScanner, io.Seeker, io.WriterTo
func TestStringReader(t *testing.T) {
	// 实例化一个新的 Reader 对象，用于对字符串内容进行读取操作
	reader := strings.NewReader(STRING_READER_DATA)
	assert.Equal(t, 68, reader.Len()) // io.Reader 接口函数，共 68 字节可以读取

	// 实例化一个 68 字节切片，用于接收读取内容
	data := make([]byte, reader.Len())

	// 读取前 7 个字节（0~6）
	n, err := reader.Read(data[:7]) // io.Reader 接口函数，读取指定长度的 byte 集合
	assert.NoError(t, err)
	assert.Equal(t, 7, n)
	assert.Equal(t, "<html>\n", string(data[:7]))

	// 跳过 7 个字节，读取接下来的 11 个字节（7~17）
	n, err = reader.ReadAt(data[7:18], 7) // io.ReaderAt 接口函数，从流的任意位置开始，读取指定长度的 byte 集合，ReadAt 函数不会移动读取指针
	assert.NoError(t, err)
	assert.Equal(t, 11, n)
	assert.Equal(t, "    <body>\n", string(data[7:18]))

	// 移动读取指针
	// 因为 io.ReaderAt 不会移动读取指针，指针停留在 6 的位置
	// 从当前位置移动 11 个字节，指针到达 17 的位置
	reader.Seek(11, io.SeekCurrent) // io.Seeker 接口函数，移动读取指针

	// 读取接下来的 31 个 byte
	n, err = ReadBytes(reader, data[18:], 31)
	assert.NoError(t, err)
	assert.Equal(t, 31, n)
	assert.Equal(t, "        <div>Hello World</div>\n", string(data[18:49]))

	// 读取接下来的 19 个字符
	n, err = ReadRune(reader, data[49:], 19)
	assert.NoError(t, err)
	assert.Equal(t, 19, n)
	assert.Equal(t, "    </body>\n</html>", string(data[49:]))

	// 查看 Reader 中剩余未读取的数据
	assert.Nil(t, reader.UnreadRune()) // io.RuneScanner 接口函数，获取尚未读取的 rune，该函数必须在 ReadRune() 之后执行，之间不能插入其它 Reader 操作
	assert.Nil(t, reader.UnreadByte()) // io.ByteScanner 接口函数，获取尚未读取的 bytes

	// 对比整个读取结果和源数据
	assert.Equal(t, STRING_READER_DATA, string(data))
}
