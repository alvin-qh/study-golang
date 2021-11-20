package buf

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试写入和读取字符串
func TestCreateBuffer(t *testing.T) {
	// 产生一个初始长度为 0 的 Buffer 对象
	buf := bytes.NewBuffer([]byte{})
	assert.Equal(t, 0, buf.Len()) // 长度为 0
	assert.Equal(t, 0, buf.Cap())

	// 产生一个初始长度为 100 的 Buffer 对象
	buf = bytes.NewBuffer(make([]byte, 100))
	assert.Equal(t, 100, buf.Len()) // 长度为 100
	assert.Equal(t, 100, buf.Cap())

	// 产生一个初始长度为 0，初始容量为 100 的 Buffer 对象
	buf = bytes.NewBuffer(make([]byte, 0, 100))
	assert.Equal(t, 0, buf.Len()) // 长度为 100
	assert.Equal(t, 100, buf.Cap())

	// 以已有内容的 bytes 集合初始化 Buffer 对象
	buf = bytes.NewBuffer([]byte(`Hello World!`))
	assert.Equal(t, 12, buf.Len()) // 长度和初始内容一致
	assert.Equal(t, 12, buf.Cap())

	s, err := buf.ReadString(byte('!')) // 读取内容到指定的 byte 值（包含指定 byte 值），返回字符串
	assert.NoError(t, err)
	assert.Equal(t, `Hello World!`, s)
}

// 测试 Buffer 的读写
// 数据读写依赖 bytes.Buffer 类型，其实现了 io.Reader 和 io.Writer 接口，可以同时进行读写操作
// 字符串类型通过编码为 utf8 编码，进行读取和写入
// 对于 int, float, bool, rune, slice 等类型，需要借助 binary 包，转换为 byte 类型后进行读写
func TestBufferRW(t *testing.T) {
	// 创建空 Buffer 以供写入
	buf := bytes.NewBuffer(make([]byte, 0, 20))

	// 写入操作
	n, err := buf.Write([]byte{1, 2, 3, 4, 5}) // 写入 byte 集合
	assert.NoError(t, err)
	assert.Equal(t, 5, n)         // 写入 5 字节
	assert.Equal(t, 5, buf.Len()) // 缓存内容 5 字节

	err = buf.WriteByte(6) // 写入单个 byte
	assert.NoError(t, err)
	assert.Equal(t, 6, buf.Len()) // 共写入 5 + 1 = 6 字节

	n, err = buf.WriteString("Hello World") // 写入字符串
	assert.NoError(t, err)
	assert.Equal(t, 11, n)         // 写入 11 字节
	assert.Equal(t, 17, buf.Len()) // 共写入 11 + 6 = 17 字节

	n, err = buf.WriteRune('好') // 写入字符
	assert.NoError(t, err)
	assert.Equal(t, 3, n)          // 写入 11 字节
	assert.Equal(t, 20, buf.Len()) // 共写入 17 + 3 = 20 字节

	// 读取操作
	data := make([]byte, 5) // 接收读取结果的 bytes
	n, err = buf.Read(data)
	assert.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, []byte{1, 2, 3, 4, 5}, data) // 读取 5 字节内容

	b, err := buf.ReadByte() // 读取单个 byte
	assert.NoError(t, err)
	assert.Equal(t, byte(6), b)

	s, err := buf.ReadString('d') // 读取字符串到指定 byte
	assert.NoError(t, err)
	assert.Equal(t, "Hello World", s)

	c, n, err := buf.ReadRune() // 读取单个字符
	assert.NoError(t, err)
	assert.Equal(t, 3, n)
	assert.Equal(t, '好', c)
}

// 测试 Buffer 对象的相互复制
func TestBufferCopy(t *testing.T) {
	buf := bytes.NewBufferString("Hello World")
	assert.Equal(t, "Hello World", buf.String())

	// 将 buf 拷贝到 newBuf
	newBuf := bytes.NewBuffer([]byte{}) // 创建目标对象
	n, err := newBuf.ReadFrom(buf)      // 将源对象拷贝到目标对象
	assert.NoError(t, err)
	assert.Equal(t, int64(11), n)
	assert.Equal(t, "Hello World", newBuf.String())

	// 将 newBuf 拷贝到 buf
	newBuf.WriteString(", Welcome")

	n, err = newBuf.WriteTo(buf)
	assert.NoError(t, err)
	assert.Equal(t, int64(20), n)
	assert.Equal(t, "Hello World, Welcome", buf.String())
}
