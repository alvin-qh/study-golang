package strings

import (
	"reflect"
	"unsafe"
)

// 对指针值加上指定的偏移量，返回新的指针地址
//
// 泛型参数:
// - `P“: 指针类型，必须是指向 `E` 类型的指针
// - `E“: 任意类型
// 函数参数:
// - `ptr“: 原始指针值
// - `offset`: 要增加的字节偏移量
// 返回值: 加上偏移量后的新指针
func PtrAdd[P ~*E, E any](ptr P, offset uintptr) P {
	p := unsafe.Pointer(ptr)
	return P(unsafe.Add(p, offset))
}

// 将切片转换为字符串
//
// 该转换使用零拷贝方式, 将切片数据就地转为字符串
//
// 该方法使用的方式已经过时, 新版本 Go 语言应使用 `unsafe.SliceData` 配合 `unsafe.String` 函数完成
func BytesToStringLegacy(b []byte) string {
	// 定义字符串变量并获取字符串头结构体地址
	s := ""
	strH := (*reflect.StringHeader)(unsafe.Pointer(&s))

	// 获取切片头结构体地址
	sliceH := (*reflect.SliceHeader)(unsafe.Pointer(&b))

	// 将切片数据设置到字符串中
	strH.Data = sliceH.Data
	strH.Len = sliceH.Len

	return s
}

// 将字符串转换为切片
//
// 该转换使用零拷贝方式, 将字符串数据就地转为切片实例
//
// 该方法使用的方式已经过时, 新版本 Go 语言应使用 `unsafe.StringData` 配合 `unsafe.Slice` 函数完成
func StringToBytesLegacy(s string) []byte {
	// 定义字节切片变量
	bs := []byte{}
	sliceH := (*reflect.SliceHeader)(unsafe.Pointer(&bs))

	// 获取字符串头结构体地址
	strH := (*reflect.StringHeader)(unsafe.Pointer(&s))

	// 将字符串数据设置到字节切片中
	sliceH.Data = strH.Data
	sliceH.Len = strH.Len

	return bs
}
