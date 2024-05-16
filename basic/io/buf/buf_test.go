package buf

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试创建缓冲区实例
func TestBytes_NewBuffer(t *testing.T) {
	// 产生一个初始长度为 0 的 Buffer 实例
	buf := bytes.NewBuffer([]byte{})
	assert.Equal(t, 0, buf.Len()) // 长度为 0
	assert.Equal(t, 0, buf.Cap())

	// 产生一个初始长度为 100 的 Buffer 实例
	buf = bytes.NewBuffer(make([]byte, 100))
	assert.Equal(t, 100, buf.Len()) // 长度为 100
	assert.Equal(t, 100, buf.Cap())

	// 产生一个初始长度为 0, 初始容量为 100 的 Buffer 实例
	buf = bytes.NewBuffer(make([]byte, 0, 100))
	assert.Equal(t, 0, buf.Len()) // 长度为 100
	assert.Equal(t, 100, buf.Cap())

	// 以已有内容的 bytes 集合初始化 Buffer 实例
	buf = bytes.NewBuffer([]byte(`Hello World!`))
	assert.Equal(t, 12, buf.Len()) // 长度和初始内容一致
	assert.Equal(t, 12, buf.Cap())

	// 读取内容到指定的 byte 值 (包含指定 byte 值), 返回字符串
	s, err := buf.ReadString(byte('!'))
	assert.Nil(t, err)
	assert.Equal(t, `Hello World!`, s)
}

// 测试扩大缓冲区容积
//
// 可以将缓冲区容积 (Capacity) 扩大到至少为要求的长度
func TestBytes_BufferGrow(t *testing.T) {
	// 产生一个初始长度为 0, 初始容量为 100 的 Buffer 实例
	buf := bytes.NewBuffer(make([]byte, 0, 100))
	assert.Equal(t, 0, buf.Len())   // 缓冲区长度为 0
	assert.Equal(t, 100, buf.Cap()) // 确认缓冲区容积

	// 扩大容积
	buf.Grow(200)
	assert.Equal(t, 0, buf.Len())            // 缓冲区长度仍为 0
	assert.GreaterOrEqual(t, buf.Cap(), 200) // 缓冲区容积增大到至少 200
}

// 测试截断缓冲区
//
// 截断缓冲区, 即将缓冲区存储的内容截断为指定长度, 缓冲区总体容积不变
func TestBytes_BufferTruncate(t *testing.T) {
	// 产生一个初始长度为 100 的 Buffer 实例
	buf := bytes.NewBuffer(make([]byte, 100))
	assert.Equal(t, 100, buf.Len()) // 缓冲区长度为 100
	assert.Equal(t, 100, buf.Cap()) // 确认缓冲区容积

	// 截断操作
	buf.Truncate(20)
	assert.Equal(t, 20, buf.Len())  // 缓冲区长度截断为 20
	assert.Equal(t, 100, buf.Cap()) // 缓冲区容积不变
}

// 测试重置缓冲区
//
// 重置缓冲区即将缓冲区内容截断为 `0`
func TestBytes_BufferReset(t *testing.T) {
	// 产生一个初始长度为 100 的 Buffer 实例
	buf := bytes.NewBuffer(make([]byte, 100))
	assert.Equal(t, 100, buf.Len()) // 缓冲区长度为 100
	assert.Equal(t, 100, buf.Cap()) // 确认缓冲区容积

	// 重置操作
	buf.Reset()
	assert.Equal(t, 0, buf.Len())   // 缓冲区内容被清空, 长度为 0
	assert.Equal(t, 100, buf.Cap()) // 缓冲区容积不变
}

// 测试缓冲区的写入和读取
//
// 数据读写依赖 `bytes.Buffer` 类型, 其实现了 `io.Reader` 和 `io.Writer` 接口, 可以同时进行读写操作
//
// 字符串类型通过编码为 utf8 编码, 进行读取和写入
//
// 对于基本类型数据, 需要借助 `binary` 包, 转换为字节串后进行读写
func TestBytes_BufferWriteRead(t *testing.T) {
	// 创建空 Buffer 以供写入
	buf := bytes.NewBuffer(make([]byte, 0, 20))

	// 写入操作
	{
		// 写入 byte 集合
		n, err := buf.Write([]byte{1, 2, 3, 4, 5})
		assert.Nil(t, err)
		assert.Equal(t, 5, n)         // 写入 5 字节
		assert.Equal(t, 5, buf.Len()) // 缓存内容 5 字节

		// 写入单个 byte
		err = buf.WriteByte(6)
		assert.Nil(t, err)
		assert.Equal(t, 6, buf.Len()) // 共写入 5 + 1 = 6 字节

		// 写入字符串
		n, err = buf.WriteString("Hello World")
		assert.Nil(t, err)
		assert.Equal(t, 11, n)         // 写入 11 字节
		assert.Equal(t, 17, buf.Len()) // 共写入 11 + 6 = 17 字节

		// 写入字符
		n, err = buf.WriteRune('好')
		assert.Nil(t, err)
		assert.Equal(t, 3, n)          // 写入 11 字节
		assert.Equal(t, 20, buf.Len()) // 共写入 17 + 3 = 20 字节
	}

	// 读取操作
	{
		// 接收读取结果的 bytes
		data := make([]byte, 5)
		n, err := buf.Read(data)
		assert.Nil(t, err)
		assert.Equal(t, 5, n)
		assert.Equal(t, []byte{1, 2, 3, 4, 5}, data) // 读取 5 字节内容

		// 读取单个 byte
		b, err := buf.ReadByte()
		assert.Nil(t, err)
		assert.Equal(t, byte(6), b)

		// 读取字符串到指定 byte
		s, err := buf.ReadString('d')
		assert.Nil(t, err)
		assert.Equal(t, "Hello World", s)

		// 读取单个字符
		c, n, err := buf.ReadRune()
		assert.Nil(t, err)
		assert.Equal(t, 3, n)
		assert.Equal(t, '好', c)
	}
}

// 测试缓冲区内容写入另一个缓冲区
func TestBytes_BufferWriteTo(t *testing.T) {
	// 创建两个缓冲实例
	buf1 := bytes.NewBufferString("Hello World")
	buf2 := bytes.NewBuffer([]byte{})

	// 将一个缓冲区的内容写入另一个缓冲区
	n, err := buf1.WriteTo(buf2)
	assert.Nil(t, err)
	assert.Equal(t, int64(11), n)
	assert.Equal(t, "Hello World", buf2.String())
}

// 测试缓冲区内容读取另一个缓冲区
func TestBytes_BufferReadFrom(t *testing.T) {
	// 创建两个缓冲实例
	buf1 := bytes.NewBufferString("Hello World")
	buf2 := bytes.NewBuffer([]byte{})

	// 从一个缓冲区读取到另一个缓冲区
	n, err := buf2.ReadFrom(buf1)
	assert.Nil(t, err)
	assert.Equal(t, int64(11), n)
	assert.Equal(t, "Hello World", buf2.String())
}
