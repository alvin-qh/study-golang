package bufio

import (
	"encoding/binary"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试零拷贝将字符串转为字节切片
func TestZeroCopyStringToBytes(t *testing.T) {
	s := "hello world, 大家好"

	bs := zeroCopyStringToBytes(s)
	assert.Equal(t, len(s), len(bs))
	assert.Equal(t, s, string(bs))
}

// 测试零拷贝将字节切片转为字符串
func TestZeroCopyBytesToString(t *testing.T) {
	bs := []byte("hello world, 大家好")

	s := zeroCopyBytesToString(bs)
	assert.Equal(t, len(bs), len(s))
	assert.Equal(t, string(bs), s)
}

// 测试 BufferIO 写入和读取数据
func TestBufferIO(t *testing.T) {
	bio := New(65, binary.BigEndian)

	bio.WriteByte(0xFF)
	bio.WriteInt8(-0x7F)
	assert.Equal(t, 2, bio.Position())

	bio.WriteInt16(-0x7FFF)
	bio.WriteUint16(0xFFFF)
	assert.Equal(t, 6, bio.Position())

	bio.WriteInt(-0x7FFFFFFFFFFFFFFF)
	assert.Equal(t, 14, bio.Position())

	bio.WriteInt32(-0x7FFFFFFF)
	bio.WriteUInt32(0xFFFFFFFF)
	assert.Equal(t, 22, bio.Position())

	bio.WriteInt64(-0x7FFFFFFFFFFFFFFF)
	bio.WriteUInt64(0xFFFFFFFFFFFFFFFF)
	assert.Equal(t, 38, bio.Position())

	bio.WriteRune('好')
	assert.Equal(t, 41, bio.Position())

	bio.Write([]byte("hello"))
	assert.Equal(t, 46, bio.Position())

	bio.WriteString(" world!")
	assert.Equal(t, 53, bio.Position())

	bio.WriteFloat32(123.123)
	assert.Equal(t, 57, bio.Position())

	bio.WriteFloat64(123123.123111)
	assert.Equal(t, 65, bio.Position())

	bio.Seek(0, io.SeekStart)
	assert.Equal(t, 0, bio.Position())

	{
		n, err := bio.ReadByte()
		assert.NoError(t, err)
		assert.Equal(t, byte(0xFF), n)
		assert.Equal(t, 1, bio.Position())
	}
	{
		n, err := bio.ReadInt8()
		assert.NoError(t, err)
		assert.Equal(t, int8(-0x7F), n)
		assert.Equal(t, 2, bio.Position())
	}
	{
		n, err := bio.ReadInt16()
		assert.NoError(t, err)
		assert.Equal(t, int16(-0x7FFF), n)
		assert.Equal(t, 4, bio.Position())
	}
	{
		n, err := bio.ReadUInt16()
		assert.NoError(t, err)
		assert.Equal(t, uint16(0xFFFF), n)
		assert.Equal(t, 6, bio.Position())
	}
	{
		n, err := bio.ReadInt()
		assert.NoError(t, err)
		assert.Equal(t, -0x7FFFFFFFFFFFFFFF, n)
		assert.Equal(t, 14, bio.Position())
	}
	{
		n, err := bio.ReadInt32()
		assert.NoError(t, err)
		assert.Equal(t, int32(-0x7FFFFFFF), n)
		assert.Equal(t, 18, bio.Position())
	}
	{
		n, err := bio.ReadUInt32()
		assert.NoError(t, err)
		assert.Equal(t, uint32(0xFFFFFFFF), n)
		assert.Equal(t, 22, bio.Position())
	}
	{
		n, err := bio.ReadInt64()
		assert.NoError(t, err)
		assert.Equal(t, int64(-0x7FFFFFFFFFFFFFFF), n)
		assert.Equal(t, 30, bio.Position())
	}
	{
		n, err := bio.ReadUInt64()
		assert.NoError(t, err)
		assert.Equal(t, uint64(0xFFFFFFFFFFFFFFFF), n)
		assert.Equal(t, 38, bio.Position())
	}
	{
		r, n, err := bio.ReadRune()
		assert.NoError(t, err)
		assert.Equal(t, '好', r)
		assert.Equal(t, 3, n)
		assert.Equal(t, 41, bio.Position())
	}
	{
		bs := make([]byte, 5)
		n, err := bio.Read(bs)
		assert.NoError(t, err)
		assert.Equal(t, []byte("hello"), bs)
		assert.Equal(t, 5, n)
		assert.Equal(t, 46, bio.Position())
	}
	{
		s, err := bio.ReadString(7)
		assert.NoError(t, err)
		assert.Equal(t, " world!", s)
		assert.Equal(t, 53, bio.Position())
	}
	{
		f, err := bio.ReadFloat32()
		assert.NoError(t, err)
		assert.Equal(t, float32(123.123), f)
		assert.Equal(t, 57, bio.Position())
	}
	{
		f, err := bio.ReadFloat64()
		assert.NoError(t, err)
		assert.Equal(t, float64(123123.123111), f)
		assert.Equal(t, 65, bio.Position())
	}

	{
		bio.Seek(-20, io.SeekEnd)
		bio.WriteString("hello")

		dst := make([]byte, 5)
		bio.ReadAt(dst, 45)
		assert.Equal(t, []byte("hello"), dst)
	}

	{
		bio.WriteAt([]byte("world"), 15)

		bio.Seek(15, io.SeekStart)
		s, err := bio.ReadString(5)

		assert.NoError(t, err)
		assert.Equal(t, "world", s)
	}
}
