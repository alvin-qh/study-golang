package generic

import (
	"fmt"
	"strconv"
)

// 定义泛型方法
//
// 该方法接受 `int` 和 `float64` 类型参数, 并返回相同类型结果
func Add[T int | float64](a, b T) T {
	return a + b
}

// 定义泛型接口类
//
// 该接口类型可以表示一个数值, 即 `int`, `int8`, `int32`, `int64`, `float32`, `float64`, `complex64` 类型中的任意一个
//
// 注意, 这样定义的接口类型只能用于泛型, 而不能用于变量或函数参数定义
type Number interface {
	int | int8 | int32 | int64 | float32 | float64 | complex64
}

// 通过 `Number` 接口定义泛型类型
//
// 由此, 所有能被 `Number` 接口表示的类型都可以作为该方法的参数和返回值
func Subtract[T Number](a, b T) T {
	return a + b
}

// 将表示切片的 `interface{}` 类型转换为切片类型
func ToSlice[T any](obj any) ([]T, error) {
	s, ok := obj.([]T)
	if !ok {
		return nil, fmt.Errorf("invalid type")
	}
	return s, nil
}

// 泛型类型自动推断
//
// 在约束类型前增加 `~` 表示可自动匹配该类型的所有衍生类型, 例如:
//
//	var n int = 1
//	s := Itoa(n)
//
//	type N int
//	var n2 N = 1
//	s := Itoa(n2)
//
// 所以, 虽然将泛型约束定义为 `int`, 但也可以接收 `N` 这样从 `int` 类型衍生而来的类型
//
// 注意, `~` 只能用于基本类型, 例如 `int32`, `float64`, `string`, `[]int` 等
func Itoa[T ~int](n T) string {
	return strconv.FormatInt(int64(n), 10)
}

// 泛型类型自动推断
//
// 如果泛型用于函数参数, 则可以通过优化泛型定义令 Go 编译器更好的推断泛型类型, 例如:
//
//	func Fill[S ~[]T, T any](s S, v T)
//
// 其泛型定义表示 `S` 可以为 `T` 类型切片或切片类型别名, 而 `T` 可以为任意类型,
// 这样就可以通过切片类型别名和元素类型分别推断泛型类型, 例如:
//
//	type Ints []int
//	ns := make(Ints, 3)
//	Fill(ns, 100)
//
// Go 可以自行推断出 `S` 类型为 `int` 类型
func Fill[S ~[]T, T any](s S, v T) {
	for i := range s {
		s[i] = v
	}
}
