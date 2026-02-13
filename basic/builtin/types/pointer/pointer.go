package pointer

import (
	"errors"
	"reflect"
	"unsafe"
)

var (
	ErrNotStruct = errors.New("not struct")
)

// 对指针值加上指定的偏移量
func PtrAdd[P ~*E, E any](ptr P, offset uintptr) P {
	p := unsafe.Pointer(ptr)
	return P(unsafe.Pointer(uintptr(p) + offset))
}

// 获取结构体各字段的 Offset 值
func FieldOffsets(v any) (map[string]uintptr, error) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil, ErrNotStruct
	}

	r := make(map[string]uintptr)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		r[f.Name] = f.Offset
	}
	return r, nil
}
