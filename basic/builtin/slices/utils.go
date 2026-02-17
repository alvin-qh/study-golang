package slices

import (
	"math/rand"
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

// 返回从 `start` 到 `stop` 且步长为 `step` 的整数切片
func Range(start, stop, step int) []int {
	a := make([]int, 0, (stop-start)/step+1)

	for i := start; i < stop; i += step {
		a = append(a, i)
	}

	return a
}

// 交换切片中指定两个位置的元素
func SwapElement[T ~[]E, E any](s T, i, j int) T {
	tmp := s[i]
	s[i] = s[j]
	s[j] = tmp

	return s
}

// 打乱切片的顺序
func Shuffle[T ~[]E, E any](s T, times int) T {
	l := len(s)

	for times > 0 {
		i := rand.Intn(l)
		j := rand.Intn(l)
		SwapElement(s, i, j)

		times--
	}
	return s
}

// 获取一个将 `val` 元素重复 `n` 次的切片
func Repeat[T any](n int, val T) []T {
	s := make([]T, n)
	for i := range n {
		s[i] = val
	}
	return s
}

// 将所给的 A 类型切片转为 B 类型切片
func Map[T any, R any](src []T, mapfn func(T) R) []R {
	rs := make([]R, len(src))
	for i, v := range src {
		rs[i] = mapfn(v)
	}
	return rs
}
