package conv

import (
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// 测试获取字符串数据指针
//
// `unsafe.StringData` 函数获取字符串数据指针, 即字符串结构中存储的字节数据地址
func TestUnsafe_StringData(t *testing.T) {
	// 获取字符串数据地址
	bs := unsafe.StringData("hello world")
	assert.Equal(t, byte('h'), *bs)

	// 通过移动指针访问字符串中其它的字节数据
	bs = PtrAdd(bs, 6)
	assert.Equal(t, byte('w'), *bs)
}

// 测试基于字节指针产生字符串变量
//
// `unsafe.String` 函数可以零拷贝方式, 基于一个指向连续字节数据的指针产生一个字符串变量
func TestUnsafe_String(t *testing.T) {
	bs := []byte("hello world")

	// 获取字节切片的数据指针, 转为字符串
	s := unsafe.String(unsafe.SliceData(bs), len(bs))
	assert.Equal(t, "hello world", s)

	// 获取字节切片的数据指针, 并将指定长度的数据转为字符串
	s = unsafe.String(unsafe.SliceData(bs), len(bs)-6)
	assert.Equal(t, "hello", s)

	// 获取字节切片的数据指针, 将指针移动 6 字节后, 将剩余部分转为字符串
	s = unsafe.String(PtrAdd(unsafe.SliceData(bs), 6), len(bs)-6)
	assert.Equal(t, "world", s)

	n := int64(7163384699739271026)

	// 将指向整数的指针转为字节指针, 并基于该指针将数据转为字符串
	s = unsafe.String((*byte)(unsafe.Pointer(&n)), unsafe.Sizeof(n))
	assert.Equal(t, "romantic", s)
}

// 测试基于字节切片产生字符串变量 (过时方式)
func TestUnsafe_BytesToStringLegacy(t *testing.T) {
	bs := []byte("hello world")

	s := BytesToStringLegacy(bs)
	assert.Equal(t, "hello world", s)
}

// 测试基于字符串产生字节切片变量 (过时方式)
func TestUnsafe_StringToBytesLegacy(t *testing.T) {
	s := "hello world"

	bs := StringToBytesLegacy(s)
	assert.Equal(t, []byte("hello world"), bs)
}
