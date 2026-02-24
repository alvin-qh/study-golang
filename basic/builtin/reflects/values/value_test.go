package values_test

import (
	"reflect"
	"study/basic/builtin/reflects"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试用结构体
type User struct {
	Id     int
	Name   string
	Gender rune
}

func TestReflect_ValueOf(t *testing.T) {
	var obj any

	// 通过反射, 尝试从 any 变量获取其存储的整型值
	t.Run("For Integer Value", func(t *testing.T) {
		// 为 any 类型变量赋予整型值 100
		obj = 100

		// 从 any 变量通过反射获取反射值对象
		tv := reflect.ValueOf(obj)

		// 确认反射值对象的值类型为 int
		assert.Equal(t, reflect.TypeFor[int](), tv.Type())

		// 确认反射值对象有效
		assert.True(t, tv.IsValid())

		// 确认反射值对象的值不为零值
		assert.False(t, tv.IsZero())

		// 确认反射值对象的值为整型
		assert.True(t, tv.CanInt())

		// 确认通过反射值获取变量本身的整数值为 100, 因为变量本身值为整数
		assert.Equal(t, int64(100), tv.Int())

		// 确认反射值对象的值不是浮点型
		assert.False(t, tv.CanFloat())

		// 确认无法通过反射值获取变量本身的浮点型值, 因为变量并不是浮点数
		assert.PanicsWithError(t, "reflect: call of reflect.Value.Float on int Value", func() { tv.Float() })

		// 确认反射值对象的值不是指针类型
		assert.False(t, tv.CanAddr())

		// 确认无法通过反射值获取变量本身的浮点型值, 因为变量并不是浮点数
		assert.PanicsWithValue(t, "reflect.Value.Addr of unaddressable value", func() { tv.Addr() })

		// 确认反射值对象的值不可修改
		assert.False(t, tv.CanSet())

		// 确认反射值对象的值可转为 any 类型
		assert.True(t, tv.CanInterface())
	})
}

// 通过反射读取对应变量的值
//
// 通过 `reflect.ValueOf` 用于获取一个变量 (`any` 类型) 的值反射
func TestReflect_GetValue(t *testing.T) {
	// 定义 any 类型变量
	var obj any

	// 令 any 类型变量赋值为 100
	obj = 100

	// 获取变量的值反射实例
	tv := reflect.ValueOf(obj)

	// 确认变量的实际类型为整型
	assert.Equal(t, reflect.Int, tv.Type().Kind())

	// 确认变量的值是否为 100
	assert.True(t, tv.CanInt())
	assert.False(t, tv.CanFloat())
	assert.False(t, tv.CanAddr())
	assert.True(t, tv.CanInterface())
	assert.True(t, tv.CanSet())

	tv.Set(reflect.Value{})

	assert.Equal(t, 100, int(tv.Int()))

	// 定义 `interface{}` 类型变量, 值为 `user` 类型结构体
	obj = User{Id: 1, Name: "Alvin", Gender: 'M'}

	// 获取变量的值反射实例
	tv = reflect.ValueOf(obj)
	assert.Equal(t, "study/basic/builtin/reflects/reflect_test.User[struct]", reflects.GetFullTypeName(tv.Type()))

	// 根据名称获取 `Id` 字段的值, 并转为 `int` 类型
	assert.Equal(t, 1, int(tv.FieldByName("Id").Int()))

	// 根据名称获取 `Name` 字段的值, 并转为 `string` 类型
	assert.Equal(t, "Alvin", tv.FieldByName("Name").String())

	// 根据名称获取 `Gender` 字段的值, 并转为 `rune` 类型
	assert.Equal(t, 'M', rune(tv.FieldByName("Gender").Int()))

	// 配合类型反射实例, 对结构体变量进行反射遍历
	names := []string{"Id", "Name", "Gender"}
	values := []any{1, "Alvin", 'M'}

	tp := reflect.TypeOf(obj)

	// 获取实例字段总数
	for i := 0; i < tp.NumField(); i++ {
		// 通过 类型反射 实例, 获取第 `i` 个字段的 类型
		field := tp.Field(i)
		assert.Equal(t, names[i], field.Name)

		// 通过 值反射 实例, 获取第 `i` 个字段的 值
		value := tv.Field(i)

		// 将所有字段值都获取为 `interface{}` 类型
		assert.EqualValues(t, values[i], value.Interface())
	}
}
