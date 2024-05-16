package bufio

import (
	"encoding/binary"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试 BufferIO 写入和读取数据
func TestBufferIO_WriteRead(t *testing.T) {
	// 定义用于读写的 IO 对象
	bio := New(65, binary.BigEndian)

	// 测试数据写入
	{
		bio.WriteByte(0xFF)
		bio.WriteInt8(-0x7F)
		assert.Equal(t, 2, bio.Position()) // 2

		bio.WriteInt16(-0x7FFF)
		bio.WriteUint16(0xFFFF)
		assert.Equal(t, 6, bio.Position()) // 4 + 2

		bio.WriteInt(-0x7FFFFFFFFFFFFFFF)
		assert.Equal(t, 14, bio.Position()) // 8 + 6

		bio.WriteInt32(-0x7FFFFFFF)
		bio.WriteUInt32(0xFFFFFFFF)
		assert.Equal(t, 22, bio.Position()) // 8 + 14

		bio.WriteInt64(-0x7FFFFFFFFFFFFFFF)
		bio.WriteUInt64(0xFFFFFFFFFFFFFFFF)
		assert.Equal(t, 38, bio.Position()) // 16 + 22

		n, err := bio.WriteRune('好')
		assert.Nil(t, err)
		assert.Equal(t, 3, n)
		assert.Equal(t, 41, bio.Position()) // 38 + 3

		bio.Write([]byte("hello"))
		assert.Equal(t, 46, bio.Position()) // 41 + 5

		bio.WriteString(" world")
		assert.Equal(t, 52, bio.Position()) // 46+6

		bio.WriteFloat32(123.123)
		assert.Equal(t, 56, bio.Position()) // 52 + 4

		bio.WriteFloat64(123123.123111)
		assert.Equal(t, 64, bio.Position()) // 56 + 8
	}

	// 移动读写位置
	bio.Seek(0, io.SeekStart)
	assert.Equal(t, 0, bio.Position())

	// 测试数据读取
	{
		n1, err := bio.ReadByte()
		assert.Nil(t, err)
		assert.Equal(t, byte(0xFF), n1)
		assert.Equal(t, 1, bio.Position()) // 1

		n2, err := bio.ReadInt8()
		assert.Nil(t, err)
		assert.Equal(t, int8(-0x7F), n2)
		assert.Equal(t, 2, bio.Position()) // 1 + 1

		n3, err := bio.ReadInt16()
		assert.Nil(t, err)
		assert.Equal(t, int16(-0x7FFF), n3)
		assert.Equal(t, 4, bio.Position()) // 2 + 2

		n4, err := bio.ReadUInt16()
		assert.Nil(t, err)
		assert.Equal(t, uint16(0xFFFF), n4)
		assert.Equal(t, 6, bio.Position()) // 2 + 4

		n5, err := bio.ReadInt()
		assert.Nil(t, err)
		assert.Equal(t, int(-0x7FFFFFFFFFFFFFFF), n5)
		assert.Equal(t, 14, bio.Position()) // 8 + 6

		n6, err := bio.ReadInt32()
		assert.Nil(t, err)
		assert.Equal(t, int32(-0x7FFFFFFF), n6)
		assert.Equal(t, 18, bio.Position()) // 4 + 14

		n7, err := bio.ReadUInt32()
		assert.Nil(t, err)
		assert.Equal(t, uint32(0xFFFFFFFF), n7)
		assert.Equal(t, 22, bio.Position()) // 4 + 18

		n8, err := bio.ReadInt64()
		assert.Nil(t, err)
		assert.Equal(t, int64(-0x7FFFFFFFFFFFFFFF), n8)
		assert.Equal(t, 30, bio.Position()) // 8 + 22

		n9, err := bio.ReadUInt64()
		assert.Nil(t, err)
		assert.Equal(t, uint64(0xFFFFFFFFFFFFFFFF), n9)
		assert.Equal(t, 38, bio.Position()) // 8 + 30

		// 读取字符, 返回字符, 字符字节数和错误信息
		r, n, err := bio.ReadRune()
		assert.Nil(t, err)
		assert.Equal(t, 3, n)
		assert.Equal(t, '好', r)
		assert.Equal(t, 41, bio.Position()) // 38 + 3

		bs := make([]byte, 5)

		n, err = bio.Read(bs)
		assert.Nil(t, err)
		assert.Equal(t, 5, n)
		assert.Equal(t, []byte("hello"), bs)
		assert.Equal(t, 46, bio.Position()) // 5 + 41

		bio.Seek(-5, io.SeekCurrent) // 46 - 5

		s, err := bio.ReadString(11)
		assert.Nil(t, err)
		assert.Equal(t, "hello world", s)
		assert.Equal(t, 52, bio.Position()) // 11 + 41

		f1, err := bio.ReadFloat32()
		assert.Nil(t, err)
		assert.Equal(t, float32(123.123), f1)
		assert.Equal(t, 56, bio.Position()) // 4 + 52

		f2, err := bio.ReadFloat64()
		assert.Nil(t, err)
		assert.Equal(t, float64(123123.123111), f2)
		assert.Equal(t, 64, bio.Position()) // 8 + 56
	}
}

// 测试在指定字节位置写入和读取数据
func TestBufferIO_WriteReadAt(t *testing.T) {
	// 定义用于读写的 IO 对象
	bio := New(65, binary.BigEndian)

	// 在第 41 字节位置写入数据
	n, err := bio.WriteAt([]byte("Hello World"), 41)
	assert.Nil(t, err)
	assert.Equal(t, 11, n)

	// 读写位置不发生变化
	assert.Equal(t, 0, bio.Position())

	bs := make([]byte, 11)

	// 从第 41 字节位置读取数据
	n, err = bio.ReadAt(bs, 41)
	assert.Nil(t, err)
	assert.Equal(t, []byte("Hello World"), bs)
	assert.Equal(t, 11, n)

	// 读写位置不发生变化
	assert.Equal(t, 0, bio.Position())
}

// 测试按行读取指定长度数据
func TestBufferIO_ReadLines(t *testing.T) {
	// 定义用于读写的 IO 对象
	bio := New(65, binary.BigEndian)

	// 在第 41 字节位置写入数据
	_, err := bio.WriteAt([]byte("Hello\nWorld"), 41)
	assert.Nil(t, err)

	// 移动读写位置到 41 字节
	bio.Seek(41, io.SeekStart)

	// 按行读取指定长度数据
	lines, err := bio.ReadLines(11)
	assert.Nil(t, err)
	assert.Equal(t, []string{"Hello", "World"}, lines)
}
