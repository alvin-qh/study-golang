package generic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 定义 int 类型的衍生类型, 此时泛型参数必须通过 `~int` 才能同时覆盖 `int` 及其衍生类型
type Integer int

// 测试 `GenericIntFloatAdd` 泛型函数
func TestGenericIntFloatAdd(t *testing.T) {
	var r interface{}

	// 测试 int 作为参数和返回值类型
	// [int] 在参数类型明确时可以省略, 编译器可以自行推断类型
	r = GenericIntFloatAdd[int](1, 2)
	assert.Equal(t, 3, r)

	// 测试 Integer 作为参数和返回值类型
	// 这里的 [Integer] 不能省略, 因为仅通过参数无法推断类型为 Integer
	r = GenericIntFloatAdd[Integer](1, 2)
	assert.Equal(t, Integer(3), r)

	// 测试 float64 作为参数和返回值类型
	// 可以增加 [float64] 来明确泛型类型
	r = GenericIntFloatAdd(1.1, 1.2)
	assert.Equal(t, float64(2.3), r)
}

// 测试 `GenericAdd` 泛型函数
func TestGenericAdd(t *testing.T) {
	var r interface{}

	// 测试 int 作为参数和返回值类型
	r = GenericAdd(1, 2)
	assert.Equal(t, 3, r)

	// 测试 complex64 作为参数和返回值类型
	// 默认情况下 complex 函数返回 complex128 类型, 所以这里需要明确泛型类型
	r = GenericAdd[complex64](complex(1, 5), complex(2, 6))
	assert.Equal(t, complex64((3 + 11i)), r)
}

func TestGenericSlice(t *testing.T) {
	var s GenericSlice[int] = []int{1, 2, 3, 4}
	assert.EqualValues(t, s, []int{1, 2, 3, 4})
}
