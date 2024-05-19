package utils

import (
	"math/rand"
)

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
	for i := 0; i < n; i++ {
		s[i] = val
	}
	return s
}
