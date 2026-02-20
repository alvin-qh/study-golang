package pointer_test

import (
	"study/basic/builtin/types/pointer"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// 测试指针的相等性
//
// 对于两个指针吧变量, 通过 `==` 运算符比较的时其存储的内存地址是否相同, 而非其指向的值是否相同
//
// 指针变量的比较逻辑和 C/C++ 语言中的指针比较逻辑相同
func TestPointer_Equality(t *testing.T) {
	// 定义一个 int 变量
	pn1 := new(100)
	pn2 := new(100)

	// 确认两个指针变量 pn1 和 pn2 存储的内存地址不同, 但它们指向的值相同, 都为 100
	assert.False(t, pn1 == pn2)
	assert.True(t, *pn1 == *pn2)

	// 注意, 通过 assert 断言比较两个指针变量时, 实际上比较的是其指向的变量, 故 pn1 和 pn2 被认为是相等的, 因为它们指向的变量值相同, 都为 100
	assert.Equal(t, pn1, pn1)

	// 如果要比较两个指针变量是否存储相同的内存地址, 则需要使用 `assert.Same` 方法, 该方法比较的是两个指针变量的内存地址是否相同, 而非其指向的值是否相同
	assert.NotSame(t, pn1, pn2)
}

// 测试将指针转为 `unsafe.Pointer` 类型
//
// `unsafe.Pointer` 类型表示一个纯粹的 "指针", 即一个无类型的内存地址, 可以转为其它类型指针
func TestUnsafe_Pointer(t *testing.T) {
	// 定义一个 int 变量
	n := 100

	// 定义 int* 类型指针, 指向 int 变量 n
	pn := &n

	// 将 int* 类型指针转为 unsafe.Pointer 类型
	p := unsafe.Pointer(pn)

	// 确认将 unsafe.Pointer 类型指针转回 int* 类型指针, 且指向的值为 100
	assert.Equal(t, (*int)(p), pn)
	assert.Equal(t, 100, *(*int)(p))
}

// 测试指针移动
//
// 当把指针转为 `unsafe.Pointer` 类型后, 即可将其进一步转为 `uintptr` 类型并进行加减操作
//
// `uintptr` 类型只能按照 1 字节的步长进行移动, 相当于 C 语言的 `void*` 类型或 `unsigned char*` 类型
//
// 注: 在 Intel 处理器上, 数值按大端方式存储, 即"BigEndian", 即 `0x12345678` 在内存中为 `0x78`, `0x56`, `0x34`, `0x12`
func TestUnsafe_PointerMovement(t *testing.T) {
	// 定义一个 int64 变量
	n := 0x1234567890ABCDEF

	// 通过 uint32 类型指针访问 int64 变量
	t.Run("access int64 variable by uint32 pointer", func(t *testing.T) {
		// 将 int* 转为 unsafe.Pointer 类型
		pn := unsafe.Pointer(&n)

		// 将 pn 转为 *uint32 类型, 相当于取 n 变量的前 4 字节
		i1 := (*uint32)(pn)

		// 相当于取 n 变量的后 4 字节
		i2 := pointer.PtrAdd((*uint32)(pn), 4)

		// 确认通过 uint32 类型指针访问 int64 变量的前 4 字节和后 4 字节, 分别为 0x90ABCDEF 和 0x12345678
		assert.Equal(t, uint32(0x90ABCDEF), *i1)
		assert.Equal(t, uint32(0x12345678), *i2)

		// 将 pn 转为 *uint32 类型, 相当于取 n 变量的前 2 字节
		s1 := (*uint16)(pn)

		// 相当于取 n 变量的第 3~4 字节
		s2 := pointer.PtrAdd((*uint16)(pn), 2)

		// 相当于取 n 变量的第 5~6 字节
		s3 := pointer.PtrAdd((*uint16)(pn), 4)

		// 相当于取 n 变量的第 7~8 字节
		s4 := pointer.PtrAdd((*uint16)(pn), 6)

		// 确认通过 uint16 类型指针访问 int64 变量的前 2 字节、第 3~4 字节、第 5~6 字节和第 7~8 字节, 分别为 0xCDEF、0x90AB、0x5678 和 0x1234
		assert.Equal(t, uint16(0xCDEF), *s1)
		assert.Equal(t, uint16(0x90AB), *s2)
		assert.Equal(t, uint16(0x5678), *s3)
		assert.Equal(t, uint16(0x1234), *s4)
	})

	// 通过 uint8 类型指针访问 int64 变量
	t.Run("access int64 variable by uint8 pointer", func(t *testing.T) {
		// 将 int* 转为 unsafe.Pointer 类型
		pn := unsafe.Pointer(&n)

		// 将 pn 转为 *uint8 类型, 相当于取 n 变量的前 1 字节
		b1 := (*uint8)(pn)

		// 相当于取 n 变量第 2 字节
		b2 := pointer.PtrAdd((*uint8)(pn), 1)

		// 相当于取 n 变量第 3 字节
		b3 := pointer.PtrAdd((*uint8)(pn), 2)

		// 相当于取 n 变量第 4 字节
		b4 := pointer.PtrAdd((*uint8)(pn), 3)

		// 相当于取 n 变量第 5 字节
		b5 := pointer.PtrAdd((*uint8)(pn), 4)

		// 相当于取 n 变量第 6 字节
		b6 := pointer.PtrAdd((*uint8)(pn), 5)

		// 相当于取 n 变量第 7 字节
		b7 := pointer.PtrAdd((*uint8)(pn), 6)

		// 相当于取 n 变量第 8 字节
		b8 := pointer.PtrAdd((*uint8)(pn), 7)

		// 确认通过 uint8 类型指针访问 int64 变量的前 1~8 字节, 值为 0xEFCDAB9078563412
		assert.Equal(t, uint8(0xEF), *b1)
		assert.Equal(t, uint8(0xCD), *b2)
		assert.Equal(t, uint8(0xAB), *b3)
		assert.Equal(t, uint8(0x90), *b4)
		assert.Equal(t, uint8(0x78), *b5)
		assert.Equal(t, uint8(0x56), *b6)
		assert.Equal(t, uint8(0x34), *b7)
		assert.Equal(t, uint8(0x12), *b8)
	})
}

type Value struct {
	I1 int8   // 实际占用 2 字节
	I2 int16  // 实际占用 2 字节
	I3 int32  // 实际占用 4 字节, I1, I2, I3 对齐到 8 字节
	I4 int64  // 实际占用 8 字节
	A  []rune // 实际占用 24 字节, 为切片结构体占用, 参考 `reflect.SliceHeader`
}

// 测试获取结构体类型的内存布局
func TestUnsafe_StructLayout(t *testing.T) {
	v := Value{}

	// 获取结构体内存大小, 以字节为单位
	size := unsafe.Sizeof(v)
	assert.Equal(t, 40, int(size))

	// 计算结构体各字段内存大小, 以字节为单位
	// 计算结果比结构体整体大小少 1 字节, 这是由于内存对齐导致的
	// 为了内存对齐 (8 字节), `I1` 字段实际占用内存 2 字节
	assert.Equal(t, 39, int(
		unsafe.Sizeof(v.I1)+
			unsafe.Sizeof(v.I2)+
			unsafe.Sizeof(v.I3)+
			unsafe.Sizeof(v.I4)+
			unsafe.Sizeof(v.A),
	))

	// 获取结构体内存对齐方式, 以 8 字节对齐
	align := unsafe.Alignof(v)
	assert.Equal(t, 8, int(align))

	// 获取结构体各字段相对于结构体地址的偏移量
	offI1 := unsafe.Offsetof(v.I1)
	assert.Equal(t, 0, int(offI1))

	offI2 := unsafe.Offsetof(v.I2)
	assert.Equal(t, 2, int(offI2))

	offI3 := unsafe.Offsetof(v.I3)
	assert.Equal(t, 4, int(offI3))

	offI4 := unsafe.Offsetof(v.I4)
	assert.Equal(t, 8, int(offI4))

	offA := unsafe.Offsetof(v.A)
	assert.Equal(t, 16, int(offA))

	_ = v
}

// 测试结构体指针
func TestUnsafe_PointerOfStruct(t *testing.T) {
	v := Value{
		I1: 1,
		I2: 2,
		I3: 3,
		I4: 4,
		A:  []rune("abc"),
	}

	// 获取结构体各字段相对结构体变量地址的偏移量
	offs, err := pointer.FieldOffsets(v)
	assert.Nil(t, err)

	// 获取结构体的指针
	pv := unsafe.Pointer(&v)

	// 已结构体地址为基础, 通过指针移动访问结构体各字段
	assert.Equal(t, int8(1), *pointer.PtrAdd((*int8)(pv), offs["I1"]))
	assert.Equal(t, int16(2), *pointer.PtrAdd((*int16)(pv), offs["I2"]))
	assert.Equal(t, int32(3), *pointer.PtrAdd((*int32)(pv), offs["I3"]))
	assert.Equal(t, int64(4), *pointer.PtrAdd((*int64)(pv), offs["I4"]))
	assert.Equal(t, []rune{'a', 'b', 'c'}, *pointer.PtrAdd((*[]rune)(pv), offs["A"]))
}
