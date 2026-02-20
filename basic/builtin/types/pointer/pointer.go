package pointer

import (
	"errors"
	"reflect"
	"unsafe"
)

var (
	ErrNotStruct = errors.New("not struct")
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

// 获取结构体各字段的 Offset 值
//
// 参数:
// - `v`: 结构体变量或结构体指针变量
// 返回值:
// - `map[string]uintptr`: 字段名称与 Offset 值的映射
// - `error`: 如果参数不是结构体类型，则返回错误
func FieldOffsets(v any) (map[string]uintptr, error) {
	// 获取变量的类型信息
	t := reflect.TypeOf(v)

	// 如果变量是指针类型，则获取指针指向的元素类型
	if t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	// 如果变量不是结构体类型，则返回错误
	if t.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}

	// 遍历结构体的字段，获取字段的 Offset 值
	r := make(map[string]uintptr)
	for f := range t.Fields() {
		r[f.Name] = f.Offset
	}

	// 返回字段 Offset 值的映射
	return r, nil
}
