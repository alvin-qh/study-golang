package conv

import "unsafe"

// 对指针值加上指定的偏移量
func PtrAdd[P ~*E, E any](ptr P, offset uintptr) P {
	p := unsafe.Pointer(ptr)
	return P(unsafe.Pointer(uintptr(p) + offset))
}
