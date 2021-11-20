package binary

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试整数和 bytes 间转换
func TestValueToBytes(t *testing.T) {
	// 用于存放 byte 的集合
	data := make([]byte, 28)

	be := binary.BigEndian
	assert.Equal(t, "binary.BigEndian", be.GoString())
	assert.Equal(t, "BigEndian", be.String())

	le := binary.LittleEndian
	assert.Equal(t, "binary.LittleEndian", le.GoString())
	assert.Equal(t, "LittleEndian", le.String())

	// 存入整数
	le.PutUint16(data[0:], 100)
	le.PutUint32(data[2:], 200)
	le.PutUint64(data[6:], 400)

	be.PutUint16(data[14:], 1000)
	be.PutUint32(data[16:], 2000)
	be.PutUint64(data[20:], 4000)

	// 获取数据
	u16 := le.Uint16(data[0:])
	assert.Equal(t, int16(100), int16(u16))

	u32 := le.Uint32(data[2:])
	assert.Equal(t, int32(200), int32(u32))

	u64 := le.Uint64(data[6:])
	assert.Equal(t, int64(400), int64(u64))

	u16 = be.Uint16(data[14:])
	assert.Equal(t, int16(1000), int16(u16))

	u32 = be.Uint32(data[16:])
	assert.Equal(t, int32(2000), int32(u32))

	u64 = be.Uint64(data[20:])
	assert.Equal(t, int64(4000), int64(u64))
}

// 利用 binary 可以向 io.Writer 写入任意类型数据，或从 io.Reader 读取任意类型数据
// binary 内部通过 BigEndian 和 LittleEndian 对象在数据和 bytes 之间进行转换
func TestBufferRW(t *testing.T) {
	// 写入操作

	// 创建空 Buffer 以供写入
	buf := bytes.NewBuffer([]byte{})

	// 写入字符串
	err := binary.Write(buf, binary.BigEndian, []byte(`Hello World`))
	assert.NoError(t, err)
	assert.Equal(t, 11, buf.Len()) // 缓存内容 11 字节

	// 写入整数值
	err = binary.Write(buf, binary.BigEndian, int64(100))
	assert.NoError(t, err)
	assert.Equal(t, 19, buf.Len()) // 共写入 11 + 8 = 19 字节，增加 int64 = 8 字节

	// 写入 bool 值
	err = binary.Write(buf, binary.BigEndian, true) // 写入 bool 类型，大端模式
	assert.NoError(t, err)
	assert.Equal(t, 20, buf.Len()) // 共写入 19 + 1 = 20 字节，增加 bool = 1 字节

	// 写入切片值
	err = binary.Write(buf, binary.LittleEndian, []float64{1.1, 2.2, 3.3}) // 写入 float64 切片，小端模式
	assert.NoError(t, err)
	assert.Equal(t, 44, buf.Len()) // 共写入 20 + 24 = 44 字节，增加 float64 * 3 = 24 字节

	// 读取操作

	// 读取 byte 集合
	data := make([]byte, 11) // 接收读取结果的 bytes
	err = binary.Read(buf, binary.BigEndian, data)
	assert.NoError(t, err)
	assert.Equal(t, `Hello World`, string(data)) // 读取 11 字节内容

	// 读取整数
	var num int64
	err = binary.Read(buf, binary.BigEndian, &num)
	assert.NoError(t, err)
	assert.Equal(t, int64(100), num)

	// 读取 bool 值
	var b bool
	err = binary.Read(buf, binary.BigEndian, &b)
	assert.NoError(t, err)
	assert.True(t, b)

	// 读取切片
	s := make([]float64, 3)
	err = binary.Read(buf, binary.LittleEndian, s)
	assert.NoError(t, err)
	assert.Equal(t, []float64{1.1, 2.2, 3.3}, s)
}

// varint 和 varuint 表示可变长度整数
// 使用可变长度整数可以减少对存储空间的消耗，存储将根据数值的实际大小变化存储长度
func TestVarint(t *testing.T) {
	// 存放二进制的 bytes
	data := make([]byte, 5)

	// 写入 varuint
	n := binary.PutVarint(data, int64(100)) // 将 int64 以“可变长度”（变体）形式存入 byte 数组。变体 （varint）可以根据数值的大小变化编码长度，可以节省存储空间
	assert.Equal(t, 2, n)                   // 变体长度为 2，较 int64 原本长度（长度8）减少 6 个字节

	n = binary.PutUvarint(data[n:], uint64(123456)) // 将 uint64 以变体形式存入 byte 数组
	assert.Equal(t, 3, n)                           // 变体长度为 2，较 uint64 原本长度（长度8）减少 6 个字节

	reader := bytes.NewReader(data)

	v, err := binary.ReadVarint(reader) // 读取一个 int64 变体
	assert.NoError(t, err)
	assert.Equal(t, int64(100), v) // 比较写入的变体和读取的变体

	uv, err := binary.ReadUvarint(reader) // 读取 uint64 变体
	assert.NoError(t, err)
	assert.Equal(t, uint64(123456), uv) // 比较写入的变体和读取的变体
}
