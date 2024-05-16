package bufio

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"unicode/utf8"
	"unsafe"
)

const (
	sizeInt8  = 1
	sizeInt16 = 2
	sizeInt32 = 4
	sizeInt64 = 8
	sizeRune  = 4
)

const (
	INT_SIZE = unsafe.Sizeof(int(0))
)

// 可用于同时读写数据的缓存类型
//
// 该类型实现了如下接口:
//   - `io.Reader`
//   - `io.Writer`
//   - `io.ReadWriter`
//   - `io.Closer`
//   - `io.Seeker`
//   - `io.ByteReader`
//   - `io.ByteWriter`
//   - `io.StringWriter`
//   - `io.ReadFrom`
//   - `io.WriteTo`
//   - `fmt.Stringer`
type BufferIO struct {
	data  []byte           // 存储字节数据的切片
	order binary.ByteOrder // 整数存储字节顺序
	pos   int              // 当前读写位置
}

// 创建一个新实例
//
// 根据所给的 `size` 参数创建缓存, 并指定整数存储的字节序实例
func New(size int, order binary.ByteOrder) *BufferIO {
	return &BufferIO{
		data:  make([]byte, size),
		order: order,
		pos:   0,
	}
}

// 关闭当前缓存
func (b *BufferIO) Close() error {
	b.data = nil
	b.pos = 0
	return nil
}

// 获取当前缓存的字节切片
func (b *BufferIO) Bytes() []byte {
	return b.data
}

// 获取当前缓存的字符串内容
func (b *BufferIO) String() string {
	return string(b.data)
}

// 获取当前缓存的字节序
func (b *BufferIO) Order() binary.ByteOrder {
	return b.order
}

// 获取当前缓存的总大小
func (b *BufferIO) Size() int {
	return len(b.data)
}

// 获取当前缓存的读写位置
func (b *BufferIO) Position() int {
	return b.pos
}

// 获取当前读写位置开始, 剩余的可用字节数
func (b *BufferIO) remaining() int {
	return len(b.data) - b.pos
}

// 向当前缓存写入字节串
//
// 如果缓存没有剩余空间, 则写入失败, 返回 `io.EOF` 错误
//
// 如果写入成功, 则返回实际写入的字节数
func (b *BufferIO) Write(data []byte) (int, error) {
	if b.remaining() == 0 {
		return 0, io.EOF
	}

	l := copy(b.data[b.pos:], data)
	b.pos += l

	var err error = nil
	if l != len(data) {
		err = io.ErrShortWrite
	}
	return l, err
}

// 以 UTF-8 编码向当前缓存写入字节串
//
// 如果写入的字符串内容超过剩余可用字节数, 则写入失败, 返回 `io.EOF` 错误
//
// 如果写入成功, 则返回实际写入的字节数
func (b *BufferIO) WriteString(s string) (int, error) {
	bs := unsafe.Slice(unsafe.StringData(s), len(s))
	return b.Write(bs)
}

// 向当前缓存写入一个字节
//
// 如果缓存剩余空间不足一个字节, 则写入失败, 返回 `io.EOF` 错误
func (b *BufferIO) WriteByte(n byte) error {
	return b.WriteUint8(n)
}

// 向当前缓存写入一个 8 位整数
//
// 如果缓存剩余空间不足一个字节, 则写入失败, 返回 `io.EOF` 错误
func (b *BufferIO) WriteInt8(n int8) error {
	return b.WriteUint8(uint8(n))
}

// 向当前缓存写入一个 8 位无符号整数
//
// 如果缓存剩余空间不足一个字节, 则写入失败, 返回 `io.EOF` 错误
func (b *BufferIO) WriteUint8(n uint8) error {
	if b.remaining() < sizeInt8 {
		return io.EOF
	}

	b.data[b.pos] = n
	b.pos += sizeInt8
	return nil
}

// 向当前缓存写入 16 位整数
//
// 如果缓存剩余空间不足两个字节, 则写入失败, 返回 `io.EOF` 错误
func (b *BufferIO) WriteInt16(n int16) error {
	return b.WriteUint16(uint16(n))
}

// 向当前缓存写入 16 位无符号整数
//
// 如果缓存剩余空间不足两个字节, 则写入失败, 返回 `io.EOF` 错误
func (b *BufferIO) WriteUint16(n uint16) error {
	if b.remaining() < sizeInt16 {
		return io.EOF
	}

	b.order.PutUint16(b.data[b.pos:], n)
	b.pos += sizeInt16
	return nil
}

// 向当前缓存写入一个 UTF-8 编码字符
//
// 如果缓存剩余空间不足以容纳该字符, 则写入失败, 返回 `io.EOF` 错误
//
// 如果写入成功, 则返回该字符实际的字节长度 (1~4 字节)
func (b *BufferIO) WriteRune(r rune) (int, error) {
	size := utf8.RuneLen(r)
	if b.remaining() < size {
		return 0, io.EOF
	}

	n := utf8.EncodeRune(b.data[b.pos:], r)
	b.pos += n
	return n, nil
}

// 向当前缓存写入一个整数
//
// 对于 32 位系统, 表示写入一个 32 位整数, 而对于 64 位系统, 则表示写入一个 64 位整数
//
// 如果缓存剩余空间不足以容纳一个整数, 则写入失败, 返回 `io.EOF` 错误
func (b *BufferIO) WriteInt(n int) error {
	return b.WriteUInt(uint(n))
}

// 向当前缓存写入一个无符号整数
//
// 对于 32 位系统, 表示写入一个 32 位无符号整数, 而对于 64 位系统, 则表示写入一个 64 位无符合整数
//
// 如果缓存剩余空间不足以容纳一个整数, 则写入失败, 返回 `io.EOF` 错误
func (b *BufferIO) WriteUInt(n uint) error {
	if INT_SIZE == 4 {
		return b.WriteUInt32(uint32(n))
	}

	if INT_SIZE == 8 {
		return b.WriteUInt64(uint64(n))
	}
	panic(fmt.Errorf("unknown size of int type"))
}

// 向当前缓存写入一个 32 位整数
//
// 如果缓存剩余空间不足 4 个字节, 则写入失败, 返回 `io.EOF` 错误
func (b *BufferIO) WriteInt32(n int32) error {
	return b.WriteUInt32(uint32(n))
}

// 向当前缓存写入一个 32 位无符号整数
//
// 如果缓存剩余空间不足 4 个字节, 则写入失败, 返回 `io.EOF` 错误
func (b *BufferIO) WriteUInt32(n uint32) error {
	if b.remaining() < sizeInt32 {
		return io.EOF
	}

	b.order.PutUint32(b.data[b.pos:], n)
	b.pos += sizeInt32
	return nil
}

// 向当前缓存写入一个 64 位整数
//
// 如果缓存剩余空间不足 8 个字节, 则写入失败, 返回 `io.EOF` 错误
func (b *BufferIO) WriteInt64(n int64) error {
	return b.WriteUInt64(uint64(n))
}

// 向当前缓存写入一个 64 位无符号整数
//
// 如果缓存剩余空间不足 8 个字节, 则写入失败, 返回 `io.EOF` 错误
func (b *BufferIO) WriteUInt64(n uint64) error {
	if b.remaining() < sizeInt64 {
		return io.EOF
	}

	b.order.PutUint64(b.data[b.pos:], n)
	b.pos += sizeInt64
	return nil
}

// 从当前缓存读取指定长度的字节, 写入 `data` 参数表示的字节切片中
//
// 如果已经可读内容, 则返回 io.EOF 错误, 如果 `data` 长度大于剩余数据,
// 则返回 `io.ErrShortBuffer` 错误和读取到的字节数
func (b *BufferIO) Read(data []byte) (int, error) {
	if b.remaining() == 0 {
		return 0, io.EOF
	}

	size := copy(data, b.data[b.pos:])
	b.pos += size

	var err error = nil
	if size < len(data) {
		err = io.ErrShortBuffer
	}

	data = nil
	return size, err
}

// 从当前缓存读取指定长度的字符串
//
// 从缓存读取 `size` 长度的字节切片, 将其转换为字符串后返回
func (b *BufferIO) ReadString(size int) (string, error) {
	if b.remaining() < size {
		return "", io.EOF
	}

	endPos := b.pos + size

	s := unsafe.String(unsafe.SliceData(b.data[b.pos:endPos]), size)
	b.pos = endPos

	return s, nil
}

// 从当前缓存读取一个字节
//
// 如果当前缓存剩余不足一个字节, 则返回错误
func (b *BufferIO) ReadByte() (byte, error) {
	return b.ReadUint8()
}

// 从当前缓存读取一个 8 位整数
//
// 如果当前缓存剩余不足一个 8 位整数, 则返回错误
func (b *BufferIO) ReadInt8() (int8, error) {
	n, err := b.ReadUint8()
	return int8(n), err
}

// 从当前缓存读取一个 8 位无符号整数
//
// 如果当前缓存剩余不足一个 8 位无符号整数, 则返回错误
func (b *BufferIO) ReadUint8() (uint8, error) {
	if b.remaining() < sizeInt8 {
		return 0, io.EOF
	}

	b.pos += sizeInt8
	return b.data[b.pos-sizeInt8], nil
}

// 从当前缓存读取一个 16 位整数
//
// 如果当前缓存剩余不足一个 16 位整数, 则返回错误
func (b *BufferIO) ReadInt16() (int16, error) {
	n, err := b.ReadUInt16()
	return int16(n), err
}

// 从当前缓存读取一个 16 位无符号整数
//
// 如果当前缓存剩余不足一个 16 位无符号整数, 则返回错误
func (b *BufferIO) ReadUInt16() (uint16, error) {
	if b.remaining() < sizeInt16 {
		return 0, io.EOF
	}

	b.pos += sizeInt16
	return b.order.Uint16(b.data[b.pos-sizeInt16:]), nil
}

// 从当前缓存读取一个 UTF-8 字符
//
// 如果当前缓存剩余不足一个字符, 则返回错误
func (b *BufferIO) ReadRune() (rune, int, error) {
	r, size := utf8.DecodeRune(b.data[b.pos:])
	if size == 0 {
		return 0, 0, io.EOF
	}

	b.pos += size
	return r, size, nil
}

// 从当前缓存读取一个整数
//
// 对于 32 位系统, 表示读取一个 32 位整数, 而对于 64 位系统, 则表示读取一个 64 位整数
//
// 如果当前缓存剩余不足一个整数, 则返回错误
func (b *BufferIO) ReadInt() (int, error) {
	n, err := b.ReadUInt()
	return int(n), err
}

// 从当前缓存读取一个无符号整数
//
// 对于 32 位系统, 表示读取一个 32 位无符号整数, 而对于 64 位系统, 则表示读取一个 64 位无符号整数
//
// 如果当前缓存剩余不足一个整数, 则返回错误
func (b *BufferIO) ReadUInt() (uint, error) {
	if INT_SIZE == 4 {
		n, err := b.ReadUInt32()
		return uint(n), err
	}
	if INT_SIZE == 8 {
		n, err := b.ReadUInt64()
		return uint(n), err
	}
	panic(fmt.Errorf("unknown size of int type"))

}

// 从当前缓存读取一个 32 位整数
//
// 如果当前缓存剩余不足一个 32 位整数, 则返回错误
func (b *BufferIO) ReadInt32() (int32, error) {
	n, err := b.ReadUInt32()
	return int32(n), err
}

// 从当前缓存读取一个 32 位无符号整数
//
// 如果当前缓存剩余不足一个 32 位无符号整数, 则返回错误
func (b *BufferIO) ReadUInt32() (uint32, error) {
	if b.remaining() < sizeInt32 {
		return 0, io.EOF
	}

	b.pos += sizeInt32
	return b.order.Uint32(b.data[b.pos-sizeInt32:]), nil
}

// 从当前缓存读取一个 64 位整数
//
// 如果当前缓存剩余不足一个 64 位整数, 则返回错误
func (b *BufferIO) ReadInt64() (int64, error) {
	n, err := b.ReadUInt64()
	return int64(n), err
}

// 从当前缓存读取一个 64 位无符号整数
//
// 如果当前缓存剩余不足一个 64 位无符号整数, 则返回错误
func (b *BufferIO) ReadUInt64() (uint64, error) {
	if b.remaining() < sizeInt64 {
		return 0, io.EOF
	}

	b.pos += sizeInt64
	return b.order.Uint64(b.data[b.pos-sizeInt64:]), nil
}

// 移动读写位置
//
// `whence` 参数为移动的起始点, 包括:
//   - `io.SeekStart` 表示从当前缓存起始位置开始移动
//   - `io.SeekCurrent` 表示从当前读写位置开始移动
//   - `io.SeekEnd` 表示从缓存末尾位置开始移动
//
// `offset` 表示移动的距离 (单位为字节), 正数表示向缓存末尾方向移动, 负数表示向缓存起始方向移动
func (b *BufferIO) Seek(offset int64, whence int) (int64, error) {
	var pos int
	switch whence {
	case io.SeekStart:
		pos = 0
	case io.SeekEnd:
		pos = len(b.data)
	case io.SeekCurrent:
		pos = b.pos
	}

	pos += int(offset)
	if pos > len(b.data) || pos < 0 {
		return 0, io.EOF
	}

	b.pos = pos
	return int64(b.pos), nil
}

// 将数据写入当前缓存的指定位置
//
// 如果偏移量参数 (`off`) 超过当前缓存长度, 则返回 `io.EOF` 错误
//
// 如果缓存剩余部分不足以完全写入 `src` 全部内容, 则返回已经写入的长度及 `io.ErrShortWrite` 错误
func (b *BufferIO) WriteAt(src []byte, off int64) (n int, err error) {
	if off >= int64(b.Size()) {
		return 0, io.EOF
	}

	size := copy(b.data[off:], src)
	if size < len(src) {
		err = io.ErrShortWrite
	}
	return size, err
}

// 从当前缓存的指定位置读取数据
//
// 如果偏移量参数 (`off`) 超过当前缓存长度, 则返回 `io.EOF` 错误
//
// 如果缓存剩余部分不足以完全读取 `dst` 所需长度, 则返回已经读取的长度及 `io.ErrShortBuffer` 错误
func (b *BufferIO) ReadAt(dst []byte, off int64) (n int, err error) {
	if off >= int64(b.Size()) {
		return 0, io.EOF
	}

	size := copy(dst, b.data[off:])
	if size < len(dst) {
		err = io.ErrShortBuffer
	}
	return size, err
}

// 将当前缓存的全部内容写入到一个 `io.Writer` 实例中
func (b *BufferIO) WriteTo(w io.Writer) (int64, error) {
	n, err := w.Write(b.data)
	return int64(n), err
}

// 从一个 `io.Reader` 实例中读取数据并写入当前缓存
func (b *BufferIO) ReadFrom(r io.Reader) (int64, error) {
	n, err := r.Read(b.data)
	return int64(n), err
}

// 向当前缓存写入 32 位浮点数
//
// 如果当前缓存剩余不足一个 32 位浮点数, 则返回错误
func (b *BufferIO) WriteFloat32(f float32) error {
	bits := math.Float32bits(f)
	return b.WriteUInt32(bits)
}

// 向当前缓存写入 64 位浮点数
//
// 如果当前缓存剩余不足一个 64 位浮点数, 则返回错误
func (b *BufferIO) WriteFloat64(f float64) error {
	bits := math.Float64bits(f)
	return b.WriteUInt64(bits)
}

// 从当前缓存读取 32 位浮点数
//
// 如果当前缓存剩余不足一个 32 位浮点数, 则返回错误
func (b *BufferIO) ReadFloat32() (float32, error) {
	bits, err := b.ReadUInt32()
	return math.Float32frombits(bits), err
}

// 从当前缓存读取 64 位浮点数
//
// 如果当前缓存剩余不足一个 64 位浮点数, 则返回错误
func (b *BufferIO) ReadFloat64() (float64, error) {
	bits, err := b.ReadUInt64()
	return math.Float64frombits(bits), err
}

// 从缓存当前位置读取之后所有的字符串行
func (b *BufferIO) ReadLines(size int) ([]string, error) {
	if b.remaining() < size {
		return nil, io.EOF
	}

	endPos := b.pos + size

	buf := bytes.NewBuffer(b.data[b.pos:endPos])
	br := bufio.NewReader(buf)

	lines := make([]string, 0, 100)
	for {
		line, err := br.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				lines = append(lines, line)
				break
			}
			return nil, err
		}
		lines = append(lines, line[:len(line)-1])
	}

	b.pos = endPos
	return lines, nil
}
