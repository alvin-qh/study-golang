package io

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试 byte 类型数据的读写操作
// byte 类型读写依赖 bytes.Buffer 类型，其实现了 io.Reader 和 io.Writer 接口，可以同时进行读写操作
// 字符串类型通过编码为 utf8 编码，进行读取和写入
// int, float, bool, rune, slice 等类型，需要借助 binary 包，转换为 byte 类型后进行读写
func TestBytesWriterAndReader(t *testing.T) {
	// 写入 bytes.Buffer 对象
	buf := bytes.NewBuffer([]byte{}) // 产生一个初始长度为 0 的 Buffer 对象进行写入

	count, err := buf.Write([]byte(`Hello World`)) // 写入编码过的字符串
	assert.NoError(t, err)                         // 确认写入成功
	assert.Equal(t, 11, count)                     // 写入 11 字节

	binary.Write(buf, binary.BigEndian, int64(100))                  // 写入 int64 类型，大端模式
	binary.Write(buf, binary.BigEndian, true)                        // 写入 bool 类型，大端模式
	binary.Write(buf, binary.LittleEndian, []float64{1.1, 2.2, 3.3}) // 写入 float64 切片，小端模式

	data := make([]byte, 4)
	count = binary.PutVarint(data, int64(100))        // 将 int64 以变体形式存入 byte 数组。变体 （varint）可以根据数值的大小变化编码长度，可以节省存储空间
	assert.Equal(t, 2, count)                         // 变体长度为 2，较 int64 原本长度（长度8）减少 6 个字节
	binary.Write(buf, binary.BigEndian, data[:count]) // 将变体写入 Buffer 对象

	count = binary.PutUvarint(data, uint64(123456))   // 将 uint64 以变体形式存入 byte 数组
	assert.Equal(t, 3, count)                         // 变体长度为 2，较 uint64 原本长度（长度8）减少 6 个字节
	binary.Write(buf, binary.BigEndian, data[:count]) // 将变体写入 Buffer 对象

	assert.Equal(t, 49, buf.Len()) // 共写入49字节
	assert.Equal(t, 84, buf.Cap())

	buf.Grow(buf.Cap() - buf.Len() + 1) //
	assert.Equal(t, 49, buf.Len())      // 共写入49字节
	assert.Equal(t, 204, buf.Cap())

	// 从 bytes.Buffer 进行读取
	buf = bytes.NewBuffer(buf.Bytes()) // 从写入结果产生新的 Buffer 对象。buf.Bytes() 返回一个 []byte 切片，包括 Buffer 对象中所有内容

	data = make([]byte, 11)
	count, err = buf.Read(data)
	assert.NoError(t, err)
	assert.Equal(t, 11, count)
	assert.Equal(t, "Hello World", string(data))

	var num int64
	err = binary.Read(buf, binary.BigEndian, &num)
	assert.NoError(t, err)
	assert.Equal(t, int64(100), num)

	var b bool
	err = binary.Read(buf, binary.BigEndian, &b)
	assert.NoError(t, err)
	assert.Equal(t, true, b)

	arr := make([]float64, 3)
	err = binary.Read(buf, binary.LittleEndian, &arr)
	assert.NoError(t, err)
	assert.Equal(t, []float64{1.1, 2.2, 3.3}, arr)

	n, err := binary.ReadVarint(buf)
	assert.NoError(t, err)
	assert.Equal(t, int64(100), n)

	un, err := binary.ReadUvarint(buf)
	assert.NoError(t, err)
	assert.Equal(t, uint64(123456), un)
}
