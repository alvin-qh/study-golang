package binary

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试数值的二进制字节序
func TestBinary_ByteOrder(t *testing.T) {
	// 小端字节序
	t.Run("LittleEndian", func(t *testing.T) {
		// 用于存放 byte 的集合
		data := make([]byte, 14)

		// 获取小端字节序实例, 序列化 3 个整数
		le := binary.LittleEndian
		le.PutUint16(data[0:], 100) // 16 位整数, 2 字节
		le.PutUint32(data[2:], 200) // 32 位整数, 4 字节
		le.PutUint64(data[6:], 400) // 64 位整数, 8 字节, 共 14 字节

		// 读取小端字节序数据
		assert.Equal(t, uint16(100), le.Uint16(data[0:])) // 16 位整数, 2 字节
		assert.Equal(t, uint32(200), le.Uint32(data[2:])) // 32 位整数, 4 字节
		assert.Equal(t, uint64(400), le.Uint64(data[6:])) // 64 位整数, 8 字节
	})

	// 大端字节序
	t.Run("BigEndian", func(t *testing.T) {
		// 用于存放 byte 的集合
		data := make([]byte, 14)

		// 获取大端字节序实例, 序列化 3 个整数
		be := binary.BigEndian
		be.PutUint16(data[0:], 1000) // 16 位整数, 2 字节
		be.PutUint32(data[2:], 2000) // 32 位整数, 4 字节
		be.PutUint64(data[6:], 4000) // 64 位整数, 8 字节, 共 14 字节

		// 读取大端字节序数据
		assert.Equal(t, uint16(1000), be.Uint16(data[0:])) // 16 位整数, 2 字节
		assert.Equal(t, uint32(2000), be.Uint32(data[2:])) // 32 位整数, 4 字节
		assert.Equal(t, uint64(4000), be.Uint64(data[6:])) // 64 位整数, 8 字节
	})
}

// 测试字节流的读取和写入
//
// 利用 binary 可以向 `io.Writer` 接口类型写入任意类型数据, 或从 `io.Reader`
// 接口类型读取任意类型数据
//
// binary 内部通过 BigEndian 和 LittleEndian 对象在数据和 bytes 之间进行转换
func TestBinary_WriteRead(t *testing.T) {
	// 创建缓冲区
	buf := bytes.NewBuffer([]byte{})

	// 写入操作
	{
		// 写入字符串
		err := binary.Write(buf, binary.BigEndian, []byte(`Hello World`))
		assert.Nil(t, err)
		assert.Equal(t, 11, buf.Len()) // 缓存内容 11 字节

		// 写入整数值
		err = binary.Write(buf, binary.BigEndian, int64(100))
		assert.Nil(t, err)
		assert.Equal(t, 19, buf.Len()) // 共写入 11 + 8 = 19 字节, 增加 int64 = 8 字节

		// 写入 bool 值
		err = binary.Write(buf, binary.BigEndian, true) // 写入 bool 类型, 大端模式
		assert.Nil(t, err)
		assert.Equal(t, 20, buf.Len()) // 共写入 19 + 1 = 20 字节, 增加 bool = 1 字节

		// 写入切片值
		err = binary.Write(buf, binary.LittleEndian, []float64{1.1, 2.2, 3.3}) // 写入 float64 切片, 小端模式
		assert.Nil(t, err)
		assert.Equal(t, 44, buf.Len()) // 共写入 20 + 24 = 44 字节, 增加 float64 * 3 = 24 字节
	}

	// 读取操作
	{
		// 读取 byte 集合
		data := make([]byte, 11) // 接收读取结果的 bytes
		err := binary.Read(buf, binary.BigEndian, data)
		assert.Nil(t, err)
		assert.Equal(t, `Hello World`, string(data)) // 读取 11 字节内容

		// 读取整数
		var num int64
		err = binary.Read(buf, binary.BigEndian, &num)
		assert.Nil(t, err)
		assert.Equal(t, int64(100), num)

		// 读取 bool 值
		var b bool
		err = binary.Read(buf, binary.BigEndian, &b)
		assert.Nil(t, err)
		assert.True(t, b)

		// 读取切片
		s := make([]float64, 3)
		err = binary.Read(buf, binary.LittleEndian, s)
		assert.Nil(t, err)
		assert.Equal(t, []float64{1.1, 2.2, 3.3}, s)
	}
}

// 测试可变长度整数
//
// `varint` 和 `varuint` 表示可变长度整数
//
// 使用可变长度整数可以减少对存储空间的消耗, 存储将根据数值的实际大小变化存储长度
func TestBinary_Variant(t *testing.T) {
	// 存放二进制的 bytes
	data := make([]byte, 5)

	// 写入可变数据
	{
		// 将 int64 以 "可变长度" (变体) 形式存入 byte 数组
		n := binary.PutVarint(data, int64(100))
		// 变体长度为 2, 较 int64 原本长度 (长度 8) 减少 6 个字节
		assert.Equal(t, 2, n)

		// 将 uint64 以变体形式存入 byte 数组
		n = binary.PutUvarint(data[n:], uint64(123456))
		// 变体长度为 3, 较 uint64 原本长度 (长度 8) 减少 6 个字节
		assert.Equal(t, 3, n)
	}

	// 读取可变数据
	{
		reader := bytes.NewReader(data)

		// 读取一个 int64 变体
		v, err := binary.ReadVarint(reader)
		assert.Nil(t, err)
		assert.Equal(t, int64(100), v)

		// 读取 uint64 变体
		uv, err := binary.ReadUvarint(reader)
		assert.Nil(t, err)
		assert.Equal(t, uint64(123456), uv)
	}
}
