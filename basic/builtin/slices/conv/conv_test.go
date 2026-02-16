package conv_test

import (
	"study/basic/builtin/slices"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// 获取切片的数据项指针
//
// `unsafe.SliceData` 函数提供了获取切片类型数据指针 (即指向切片中第一个元素地址) 的方法
func TestUnsafe_SliceData(t *testing.T) {
	// 转换字符串, 得到字节串
	bs := []byte("hello world")

	// 获取字节切片的数据指针
	ptr := unsafe.SliceData(bs)
	assert.Equal(t, &bs[0], ptr)
	assert.Equal(t, byte('h'), *ptr)

	// 移动指针, 指向切片不同元素
	ptr = slices.PtrAdd(ptr, 2)
	assert.Equal(t, byte('l'), *ptr)

	ptr = slices.PtrAdd(ptr, uintptr(len(bs)-1-2))
	assert.Equal(t, byte('d'), *ptr)
}

// 从数据指针还原切片
//
// `unsafe.Slice` 函数通过零拷贝方式, 在现有连续数据指针的基础上形成切片类型变量
func TestUnsafe_Slice(t *testing.T) {
	bs := []byte("hello world")

	// 将字节切片转为指针
	ptr := unsafe.SliceData(bs)

	// 将字节指针向后移动 6 字节
	ptr = slices.PtrAdd(ptr, 6)

	// 将移动过的指针还原为切片
	s := unsafe.Slice(ptr, len(bs)-6)
	assert.Equal(t, []byte("world"), s)

	n := int64(0x12345678ABCDEF90)

	// 将 64 位整数地址转为字节指针
	ptr = (*byte)(unsafe.Pointer(&n))

	// 将指针转为切片, 确认切片中包含 8 个字节, 为 64 位整数的 8 个字节
	s = unsafe.Slice(ptr, 8)
	assert.Equal(t, []byte{0x90, 0xEF, 0xCD, 0xAB, 0x78, 0x56, 0x34, 0x12}, s)
}
