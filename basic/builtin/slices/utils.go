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
	// 将指针转换为 `unsafe.Pointer` 类型，以便进行地址运算
	p := unsafe.Pointer(ptr)

	// 使用 `unsafe.Add` 函数将指针加上指定的偏移量，并返回新的指针
	return P(unsafe.Add(p, offset))
}

// 返回从 `start` 到 `stop` 且步长为 `step` 的整数切片
func Range(start, stop, step int) []int {
	// 创建一个切片对象, 切片的 Capacity 为 (stop-start)/step+1
	a := make([]int, 0, (stop-start)/step+1)

	// 使用 for 循环从 `start` 开始，依次增加 `step`，直到达到 `stop`，将每个值添加到切片中
	for i := start; i < stop; i += step {
		a = append(a, i)
	}

	// 返回生成的整数切片
	return a
}

// 交换切片中指定两个位置的元素
func SwapElement[T ~[]E, E any](s T, i, j int) T {
	// 交换切片中下标为 `i` 和 `j` 的元素
	tmp := s[i]
	s[i] = s[j]
	s[j] = tmp

	return s
}

// 打乱切片的顺序
func Shuffle[T ~[]E, E any](s T, times int) T {
	// 获取切片的长度
	l := len(s)

	// 循环 times 次，每次随机交换两个元素
	for range times {
		// 生成两个随机下标 `i` 和 `j`，范围在 0 到切片长度之间
		i := rand.Intn(l)
		j := rand.Intn(l)
		if i != j {
			// 交换切片中下标为 `i` 和 `j` 的元素
			SwapElement(s, i, j)
		}
	}

	// 返回打乱顺序后的切片
	return s
}

// 获取一个将 `val` 元素重复 `n` 次的切片
func Repeat[T any](n int, val T) []T {
	// 创建一个切片对象, 切片长度为 n
	s := make([]T, n)

	// 使用 for 循环将 `val` 元素重复 `n` 次，并赋给切片的每个元素
	for i := range n {
		s[i] = val
	}

	// 返回生成的切片
	return s
}

// 将所给的 A 类型切片转为 B 类型切片
func Map[T any, R any](src []T, mapfn func(T) R) []R {
	// 创建一个切片对象, 切片长度与 src 相同
	rs := make([]R, len(src))

	// 使用 for 循环将 src 中的每个元素映射到 rs 中
	for i, v := range src {
		rs[i] = mapfn(v)
	}

	// 返回映射后的切片
	return rs
}
