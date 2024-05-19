package generic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试 `GenericIntFloatAdd` 泛型函数
func TestGeneric_Add(t *testing.T) {
	var r interface{}

	// 测试 int 作为参数和返回值类型
	// [int] 在参数类型明确时可以省略, 编译器可以自行推断类型
	r = Add(1, 2)
	assert.Equal(t, 3, r)

	// 测试 float64 作为参数和返回值类型
	// 可以增加 [float64] 来明确泛型类型
	r = Add(1.1, 1.2)
	assert.Equal(t, float64(2.3), r)
}

// 测试 `GenericAdd` 泛型函数
func TestGeneric_Subtract(t *testing.T) {
	var r interface{}

	// 测试 int 作为参数和返回值类型
	r = Subtract(1, 2)
	assert.Equal(t, 3, r)

	// 测试 complex64 作为参数和返回值类型
	// 默认情况下 complex 函数返回 complex128 类型, 所以这里需要明确泛型类型
	r = Subtract[complex64](complex(1, 5), complex(2, 6))
	assert.Equal(t, complex64((3 + 11i)), r)
}

// 测试 `interface{}` 类型转切片
//
// 这里的 `interface{}` 类型值本身应为切片类型, 才能转换回原本的切片类型
func TestGeneric_ToSlice(t *testing.T) {
	var s any = []int{1, 2, 3, 4}

	r, err := ToSlice[int](s)
	assert.Nil(t, err)
	assert.Equal(t, []int{1, 2, 3, 4}, r)
}

// 测试泛型约束及其衍生类型
func TestGeneric_Itoa(t *testing.T) {
	// 定义 int 类型变量, 调用泛型函数
	n := int(100)
	s := Itoa(n)
	assert.Equal(t, "100", s)

	type NN int

	// 定义衍生类型 N 类型变量, 调用泛型函数
	var nn NN = 200
	s = Itoa(nn)
	assert.Equal(t, "200", s)
}

// 测试用一个值填充切片的所有元素
func TestGeneric_Fill(t *testing.T) {
	type Ints []int

	ns := make(Ints, 3)

	Fill(ns, 100)
	assert.Equal(t, Ints{100, 100, 100}, ns)
}
