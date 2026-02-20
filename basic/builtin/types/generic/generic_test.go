package generic_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 定义泛型方法
//
// 该方法接受 `int` 和 `float64` 类型参数, 并返回相同类型结果
func Add[T int | float64](a, b T) T {
	// 返回两个泛型参数的和
	return a + b
}

// 测试 `Add` 泛型函数
func TestGeneric_Add(t *testing.T) {
	var r any

	// 测试 int 作为参数和返回值类型
	// [int] 在参数类型明确时可以省略, 编译器可以自行推断类型
	r = Add /*[int]*/ (1, 2)
	assert.IsType(t, int(0), r)
	assert.Equal(t, 3, r)

	// 测试 float64 作为参数和返回值类型
	// 可以增加 [float64] 来明确泛型类型
	r = Add(1.1, 1.2)
	assert.IsType(t, float64(0), r)
	assert.Equal(t, float64(2.3), r)
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
	// 返回两个泛型参数的差
	return a - b
}

// 测试 `Subtract` 泛型函数
func TestGeneric_Subtract(t *testing.T) {
	var r any

	// 测试 int 作为参数和返回值类型
	r = Subtract(1, 2)
	assert.IsType(t, int(0), r)
	assert.Equal(t, -1, r)

	// 测试 complex64 作为参数和返回值类型
	// 默认情况下 complex 函数返回 complex128 类型, 所以这里需要明确泛型类型
	r = Subtract[complex64](complex(1, 5), complex(2, 6))
	assert.IsType(t, complex64(0), r)
	assert.Equal(t, complex64((-1 + -1i)), r)
}

// 将表示切片的 `any` 类型转换为切片类型
func ToSlice[T any](obj any) ([]T, error) {
	// 将 any 类型参数转换为 T 类型切片
	s, ok := obj.([]T)
	if !ok {
		return nil, fmt.Errorf("invalid type")
	}

	// 返回转换后的切片
	return s, nil
}

// 测试 `any` 类型转切片
//
// 这里的 `any` 类型值本身类型应为切片类型, 才能转换回原本的切片类型对象
func TestGeneric_ToSlice(t *testing.T) {
	// 定义 any 类型变量, 值为一个整型切片
	var s any = []int{1, 2, 3, 4}

	// 调用泛型函数, 将 any 类型变量转换为 int 类型切片
	r, err := ToSlice[int](s)

	// 确认转换成功, 并且结果正确
	assert.Nil(t, err)
	assert.Equal(t, []int{1, 2, 3, 4}, r)
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
	// 将泛型参数转换为 int64 类型, 并返回其字符串表示
	return strconv.FormatInt(int64(n), 10)
}

// 测试泛型约束及其衍生类型
func TestGeneric_Itoa(t *testing.T) {
	// 测试 int 类型参数转换为字符串
	t.Run("int", func(t *testing.T) {
		// 定义 int 类型变量
		n := int(100)

		// 调用泛型函数, 将 int 类型变量转换为字符串
		s := Itoa(n)

		// 确认转换结果正确
		assert.Equal(t, "100", s)
	})

	// 测试 int 类型的衍生类型参数转换为字符串
	t.Run("derived type", func(t *testing.T) {
		// 定义 int 类型的衍生类型 NN, 该类型也满足泛型约束, 因此也可以作为参数调用泛型函数
		type NN int

		// 为 NN 类型变量赋值
		var nn NN = 200

		// 调用泛型函数, 将 NN 类型变量转换为字符串
		s := Itoa(nn)

		// 确认转换结果正确
		assert.Equal(t, "200", s)
	})
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
	// 遍历切片 `s` 的每个元素, 将其值设置为 `v`
	for i := range s {
		s[i] = v
	}
}

// 测试用一个值填充切片的所有元素
func TestGeneric_Fill(t *testing.T) {
	// 定义一个切片类型的衍生类型 `Ints`, 该类型满足泛型约束, 因此也可以作为参数调用泛型函数
	type Ints []int

	// 创建一个 `Ints` 类型的切片, 长度为 3
	ns := make(Ints, 3)

	// 调用泛型函数, 用 100 填充切片的所有元素
	Fill(ns, 100)

	// 确认切片的所有元素都被正确填充为 100
	assert.Equal(t, Ints{100, 100, 100}, ns)
}
