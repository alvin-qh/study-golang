package conv_test

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 定义结构体
type User struct {
	Id     int
	Name   string
	Gender rune
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

// 测试 `any` 类型的转换
//
// `any` 类型是 `interface{}` 类型的别名
//
// `any` 类型相当于 "任意类型", 及任意类型都可以转为 `any` 类型, 且 `any` 类型可以转回为其原始类型
func TestConv_ConvertAnyType(t *testing.T) {
	var v any

	// 将值类型转为 any 类型, 并确认变量 v 的实际类型为 `reflect.Int` 类型
	v = int(10)
	assert.Equal(t, reflect.Int, reflect.TypeOf(v).Kind())
	assert.Equal(t, 10, v)

	// 将结构体转为 any 类型, 并确认变量 v 的实际类型为 `reflect.Struct` 类型
	v = User{
		Id:     1,
		Name:   "Alvin",
		Gender: 'F',
	}
	assert.Equal(t, reflect.Struct, reflect.TypeOf(v).Kind())
	assert.Equal(t, "User", reflect.TypeOf(v).Name())

	// 将 any 类型变量转为其原始类型变量
	// 转换返回两个值, 第一个为是转换后的值, 第二个表示是否转换成功 (第二个返回值可以忽略)
	u, ok := v.(User)
	assert.True(t, ok)
	assert.Equal(t, User{
		Id:     1,
		Name:   "Alvin",
		Gender: 'F',
	}, u)

	// 指针类型转换为 any 类型, 并确认变量 v 的实际类型为 `reflect.Pointer` 类型
	v = &User{
		Id:     1,
		Name:   "Alvin",
		Gender: 'F',
	}
	assert.Equal(t, reflect.Pointer, reflect.TypeOf(v).Kind())
	assert.Equal(t, reflect.Struct, reflect.TypeOf(v).Elem().Kind())
	assert.Equal(t, "User", reflect.TypeOf(v).Elem().Name())

	// 将 any 类型转为指针类型
	// 转换返回两个值, 第一个为是转换后的值, 第二个表示是否转换成功 (第二个返回值可以忽略)
	pu, ok := v.(*User)
	assert.True(t, ok)
	assert.Equal(t, User{
		Id:     1,
		Name:   "Alvin",
		Gender: 'F',
	}, *pu)

	// 如果类型转换时只返回一个值, 则转换失败会抛出 Panic
	assert.Panics(t, func() {
		// v 是 User* 类型而非 User 类型, 故转换会失败
		u := v.(User)
		assert.Equal(t, User{}, u)
	})
}

// 测试值类型的强制类型转换
//
// 对于值类型, 可以通过类型运算符进行类型转换
//
// 类型转换是赋值的一种副作用, 即在内存间进行数值复制的时候, 对数值做了一次类型变更操作, 例如将 8byte 数值复制到 4byte 空间中,
// 所以类型转换的过程中可能会丢失精度
func TestConv_ForceConvert(t *testing.T) {
	var v1 float64 = 123.456
	var v2 int64 = int64(v1)

	// 转换前后的两个变量不相同
	assert.NotEqual(t, v2, v1)
	assert.NotSame(t, &v2, &v1)

	var v3 int32 = int32(v2)

	// 转换前后的两个变量不相同
	assert.NotEqual(t, v3, v2)
	assert.NotSame(t, &v3, &v2)

	// 转换前后的两个变量值相同但类型不同
	assert.EqualValues(t, v2, v3)
}

// 测试利用 `switch` 语句进行类型转换
func TestConv_ConvertWithSwitch(t *testing.T) {
	// 创建不同类型的对象
	obj := makeDifferentTypeObject()

	// 通过 switch 语句进行类型转换
	// 每个分支用于判断 `v` 变量的一种类型, 如果类型匹配到具体分支, 则 `vv` 变量是该类型的值
	switch val := obj.(type) {
	case int:
		assert.Equal(t, 100, val)
	case string:
		assert.Equal(t, "Hello", val)
	case User:
		assert.Equal(t, User{
			Id:     1,
			Name:   "Alvin",
			Gender: 'M',
		}, val)
	default:
		assert.Fail(t, "unknown type")
	}
}
