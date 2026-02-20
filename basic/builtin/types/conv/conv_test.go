package conv_test

import (
	"math/rand"
	"study/basic/builtin/types/conv"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 定义结构体
type User struct {
	Id     int
	Name   string
	Gender rune
}

// 测试 `any` 类型的转换
//
// `any` 类型是 `interface{}` 类型的别名
//
// `any` 类型相当于 "任意类型", 及任意类型都可以转为 `any` 类型, 且 `any` 类型可以转回为其原始类型
func TestConv_ConvertAnyType(t *testing.T) {
	// 定义 any 类型变量
	var v any

	// 将值类型转为 any 类型
	v = int(10)

	// 确认变量 v 的实际类型为 `int` 类型, 且 v 的值为 10
	assert.IsType(t, int(0), v)
	assert.Equal(t, 10, v)

	// 将结构体转为 any 类型
	v = User{
		Id:     1,
		Name:   "Alvin",
		Gender: 'F',
	}

	// 确认变量 v 的实际类型为 `User` 类型, 且 v 的值为 User{1, "Alvin", 'F'}
	assert.IsType(t, User{}, v)
	assert.Equal(t, User{1, "Alvin", 'F'}, v)

	// 将 any 类型变量转为其原始类型变量, 并确认转换成功
	// 转换返回两个值, 第一个为是转换后的值, 第二个表示是否转换成功 (第二个返回值可以忽略)
	_, ok := v.(User)
	assert.True(t, ok)

	// 将结构体指针转为 any 类型
	v = &User{
		Id:     1,
		Name:   "Alvin",
		Gender: 'F',
	}

	// 确认变量 v 的实际类型为 `*User` 类型, 且 v 的值为 &User{1, "Alvin", 'F'}
	assert.IsType(t, &User{}, v)
	assert.Equal(t, &User{1, "Alvin", 'F'}, v)

	// 将 any 类型变量转为其原始类型变量, 并确认转换成功
	pu, ok := v.(*User)
	assert.True(t, ok)

	// 确认变量 pu 的实际类型为 `*User` 类型, 且 pu 的值为 &User{1, "Alvin", 'F'}
	assert.Equal(t, User{
		Id:     1,
		Name:   "Alvin",
		Gender: 'F',
	}, *pu)

	// 测试转换失败的情况, 即将 any 类型变量转为一个与其实际类型不匹配的类型
	assert.Panics(t, func() {
		// v 是 User* 类型而非 User 类型, 故转换会失败
		_ = v.(User)
	})
}

// 测试值类型的强制类型转换
//
// 对于值类型, 可以通过类型运算符进行类型转换
//
// 类型转换是赋值的一种副作用, 即在内存间进行数值复制的时候, 对数值做了一次类型变更操作, 例如将 8byte 数值复制到 4byte 空间中,
// 所以类型转换的过程中可能会丢失精度
func TestConv_ForceConvert(t *testing.T) {
	// 定义一个 float64 类型的变量, 并将其转换为 int64 类型
	var v1 float64 = 123.456
	var v2 int64 = int64(v1)

	// 确认转换前后的两个变量值不同, 因为在类型转换过程中丢失了小数部分, 同时两个变量所在的内存地址不同
	assert.NotEqualValues(t, v1, v2)
	assert.NotSame(t, &v2, &v1)

	// 创建一个 int32 类型的变量, 并将其转换为 int64 类型
	var v3 int32 = int32(v2)

	// 确认转换前后的两个变量值相同, 因为 int32 类型的数值可以完全表示 int64 类型的数值, 同时两个变量所在的内存地址不同
	assert.EqualValues(t, v2, v3)
	assert.NotSame(t, &v3, &v2)
}

// 创建类型不同的对象
func makeDifferentTypeObject() any {
	var obj any

	// 根据 0~2 的随机数结果, 创建不同类型的对象
	switch rand.Intn(3) {
	case 0:
		obj = 100 // 创建整数类型对象
	case 1:
		obj = "Hello" // 创建字符串类型对象
	case 2:
		obj = User{ // 创建结构体类型对象
			Id:     1,
			Name:   "Alvin",
			Gender: 'M',
		}
	}
	return obj
}

// 测试利用 `switch` 语句进行类型转换
func TestConv_ConvertWithSwitch(t *testing.T) {
	// 创建不同类型的对象
	obj := makeDifferentTypeObject()

	// 通过 switch 语句进行类型转换
	// 每个分支用于判断 `v` 变量的一种类型, 如果类型匹配到具体分支, 则 `vv` 变量是该类型的值
	switch val := obj.(type) {
	case int:
		// 确认变量 v 的实际类型为 `int` 类型, 且 v 的值为 100
		assert.Equal(t, 100, val)
	case string:
		// 确认变量 v 的实际类型为 `string` 类型, 且 v 的值为 "Hello"
		assert.Equal(t, "Hello", val)
	case User:
		// 确认变量 v 的实际类型为 `User` 类型, 且 v 的值为 User{1, "Alvin", 'M'}
		assert.Equal(t, User{
			Id:     1,
			Name:   "Alvin",
			Gender: 'M',
		}, val)
	default:
		// 如果没有匹配到任何分支, 则说明 obj 的类型未知, 这时应该让测试失败
		assert.Fail(t, "unknown type")
	}
}

// 测试指定类型切片和 `any` 类型切片的转换
//
// 在 Go 语言中, 一般不推荐使用 `any` 类型的切片, 即 `[]any` 类型, 如果要用泛化类型表示数组,
// 则使用 `any` 直接表示即可
//
// 注意, 要在 `[]any` 类型切片和其它类型切片间转换, 则需要通过一个 `O(n)` 复杂度的循环才能完成
func TestConv_SliceTypeConversion(t *testing.T) {
	// 定义 any 类型变量, 并确认其实际类型为 `[]int` 类型
	var v any = []int{1, 2, 3, 4, 5}
	assert.IsType(t, []int{}, v)

	// 将 any 类型转为指定类型切片类型, 并确认转换成功, 且转换后的切片值为 []int{1, 2, 3, 4, 5}
	s, ok := conv.AnyToSlice[int](v)
	assert.True(t, ok)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, s)

	// 将指定类型的切片转为 `[]any` 类型, 并确认转换后的切片长度为 5
	vs := conv.TypedSliceToAnySlice([]int{1, 2, 3, 4, 5})
	assert.Equal(t, []any{1, 2, 3, 4, 5}, vs)

	// 将 `[]any` 类型切片转为指定类型, 并确认转换成功, 且转换后的切片值为 []int{1, 2, 3, 4, 5}
	s, ok = conv.AnySliceToTypedSlice[int](vs)
	assert.True(t, ok)
	assert.Equal(t, []int{1, 2, 3, 4, 5}, s)
}
