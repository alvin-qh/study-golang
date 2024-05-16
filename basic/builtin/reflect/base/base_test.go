package base

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// 测试用结构体
type User struct {
	Id     int
	Name   string
	Gender rune
}

// 通过反射获取变量类型
//
// 通过 `reflect.TypeOf` 函数可以获取变量的类型
func TestReflect_TypeOf(t *testing.T) {
	// 定义实际类型为 `int64` 类型
	var obj any = int64(100)

	// 获取变量的类型反射实例
	tp := reflect.TypeOf(obj)
	assert.Equal(t, ".int64[int64]", GetFullTypeName(tp))

	// 定义实际类型为 `User` 类型
	obj = User{}

	// 获取变量的反射实例
	tp = reflect.TypeOf(obj)
	assert.Equal(t, "study/basic/builtin/reflect/base.User[struct]", GetFullTypeName(tp))

	// 获取指针变量的反射实例
	tp = reflect.TypeOf(&obj)
	assert.Equal(t, ".[ptr]", GetFullTypeName(tp))

	// 获取指针变量的反射实例
	tp = tp.Elem()
	assert.Equal(t, ".[interface]", GetFullTypeName(tp))

	// 获取切片变量的反射实例
	tp = reflect.TypeOf([]int{1, 2, 3, 4})
	assert.Equal(t, ".[slice]", GetFullTypeName(tp))

	// 获取数组变量的反射实例
	tp = reflect.TypeOf([...]int{1, 2, 3, 4})
	assert.Equal(t, ".[array]", GetFullTypeName(tp))

	// 获取 Map 变量的反射实例
	tp = reflect.TypeOf(map[string]any{"a": 1, "b": "Hello"})
	assert.Equal(t, ".[map]", GetFullTypeName(tp))
}

// 通过反射获取指针类型及其指向实例的类型
//
// 若 `reflect.Type` 类型是一个指针类型, 则可以通过 `reflect.Type.Elem()` 函数获取该指针指向的实例类型
func TestReflect_Elem(t *testing.T) {
	n := 100

	// obj 保存指向变量的指针
	var obj any = &n

	// 获取指针类型变量的类型
	tp := reflect.TypeOf(obj)
	assert.Equal(t, ".[ptr]", GetFullTypeName(tp))

	// 获取指针所指向的实例类型
	tp = tp.Elem()
	assert.Equal(t, ".int[int]", GetFullTypeName(tp))

	// obj 保存指向结构体的指针
	obj = &User{}

	// 获取指针变量的类型
	tp = reflect.TypeOf(obj)
	assert.Equal(t, ".[ptr]", GetFullTypeName(tp))

	// 获取指针所指向的实例类型
	tp = tp.Elem()
	assert.Equal(t, "study/basic/builtin/reflect/base.User[struct]", GetFullTypeName(tp))
}

// 通过反射读取实例值
//
// 通过 `reflect.ValueOf` 用于获取一个变量 (`interface{}` 类型) 的值反射
func TestReflect_ValueOf(t *testing.T) {
	// 定义 interface{} 类型变量, 值为整型
	var obj any = 100

	// 获取变量的 值反射 实例
	tv := reflect.ValueOf(obj)
	assert.Equal(t, ".int[int]", GetFullTypeName(tv.Type()))

	// 通过反射获取值
	assert.Equal(t, 100, int(tv.Int()))

	// 定义 `interface{}` 类型变量, 值为 `user` 类型结构体
	obj = User{Id: 1, Name: "Alvin", Gender: 'M'}

	// 获取变量的值反射实例
	tv = reflect.ValueOf(obj)
	assert.Equal(t, "study/basic/builtin/reflect/base.User[struct]", GetFullTypeName(tv.Type()))

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

// 通过反射读取指针及其指向的实例值
//
// 若 `reflect.Value` 引用了一个指针类型值, 则可通过 `reflect.Value.Elem` 方法获取其指向实例的值
func TestReflect_ValueOfPtr(t *testing.T) {
	n := 100

	// obj 为变量指针
	var obj any = &n

	// 获取指针类型变量的值实例
	tv := reflect.ValueOf(obj)
	assert.Equal(t, uintptr(unsafe.Pointer(&n)), tv.Pointer())

	// 获取指针指向实例的值实例
	tv = tv.Elem()
	assert.Equal(t, ".int[int]", GetFullTypeName(tv.Type()))
	assert.Equal(t, 100, int(tv.Int()))

	obj = &User{Id: 1, Name: "Alvin", Gender: 'M'}

	// obj 为结构体指针, 获取指针类型变量的值实例
	tv = reflect.ValueOf(obj)
	assert.Equal(t, ".[ptr]", GetFullTypeName(tv.Type()))
	assert.Equal(t, uintptr(unsafe.Pointer(obj.(*User))), tv.Pointer())

	// 获取指针指向的实例值实例
	tv = tv.Elem()
	assert.Equal(t, "study/basic/builtin/reflect/base.User[struct]", GetFullTypeName(tv.Type()))
	assert.Equal(t, 1, int(tv.FieldByName("Id").Int()))
	assert.Equal(t, "Alvin", tv.FieldByName("Name").String())
	assert.Equal(t, 'M', rune(tv.FieldByName("Gender").Int()))

	obj = User{Id: 1, Name: "Alvin", Gender: 'M'}

	// obj 为 interface{} 指针, 获取类型为指针类型
	tv = reflect.ValueOf(&obj)
	assert.Equal(t, ".[ptr]", GetFullTypeName(tv.Type()))

	// 指针的值为实例地址
	assert.Equal(t, uintptr(unsafe.Pointer(&obj)), tv.Pointer())

	// 获取指针指向的实例值
	tv = tv.Elem()
	assert.Equal(t, ".[interface]", GetFullTypeName(tv.Type()))

	// 其值为 User 实例
	assert.Equal(t, obj.(User), tv.Interface().(User))

	// 再次从 interface{} 类型解除引用, 获取其原始值
	tv = tv.Elem()
	assert.Equal(t, "study/basic/builtin/reflect/base.User[struct]", GetFullTypeName(tv.Type()))
	assert.Equal(t, 1, int(tv.FieldByName("Id").Int()))
	assert.Equal(t, "Alvin", tv.FieldByName("Name").String())
	assert.Equal(t, 'M', rune(tv.FieldByName("Gender").Int()))
}

// 通过反射读取"切片"实例值
//
// 若 `reflect.Value` 实例引用了一个切片类型值, 则可以通过 `reflect.Value` 实例提供的一组切片反射方法对其进行操作, 包括:
//   - `reflect.Value.Len()`
//   - `reflect.Value.Index(n)`
//   - `reflect.Value.Slice(m, n)`
func TestReflect_ValueOfSlice(t *testing.T) {
	// 定义一个切片实例, 类型转为 interface{} 类型
	var obj any = []any{1, "Hello", false}

	// 获取实例的 值反射 结果
	tv := reflect.ValueOf(obj)
	assert.Equal(t, ".[slice]", GetFullTypeName(tv.Type()))

	// 通过值反射实例获取切片长度
	len := tv.Len()
	assert.Equal(t, 3, len)

	// 获取切片指定下标的值
	val := tv.Index(0)
	assert.Equal(t, 1, val.Interface().(int))

	val = tv.Index(1)
	assert.Equal(t, "Hello", val.Interface().(string))

	val = tv.Index(2)
	assert.Equal(t, false, val.Interface().(bool))

	// 通过值反射实例进行切片操作
	tv = tv.Slice(0, 2)

	// 获取切片长度
	len = tv.Len()
	assert.Equal(t, 2, len)

	// 获取切片指定下标的值
	val = tv.Index(0)
	assert.Equal(t, 1, val.Interface().(int))

	val = tv.Index(1)
	assert.Equal(t, "Hello", val.Interface().(string))
}

// 通过反射读取"字典"实例值
func TestReflect_ValueOfMap(t *testing.T) {
	// 定义字典类型变量, 类型转为 any 类型
	var obj any = map[string]any{"a": 1, "b": "Hello", "c": false}

	// 获取实例的 值反射 结果
	tv := reflect.ValueOf(obj)
	assert.Equal(t, ".[map]", GetFullTypeName(tv.Type()))

	// 通过反射获取字典所有的 key 的集合
	keys := tv.MapKeys()
	assert.Len(t, keys, 3)
	assert.ElementsMatch(t, []string{"a", "b", "c"}, []string{keys[0].String(), keys[1].String(), keys[2].String()})

	// 通过反射, 根据 key 的值获取 Value, 注意, 这里的 Key 必须是"值反射"实例
	val := tv.MapIndex(reflect.ValueOf("a"))
	assert.Equal(t, 1, val.Interface().(int))

	val = tv.MapIndex(reflect.ValueOf("b"))
	assert.Equal(t, "Hello", val.Interface().(string))

	val = tv.MapIndex(reflect.ValueOf("c"))
	assert.Equal(t, false, val.Interface().(bool))

	// 获取 Key/Value 对的迭代器
	iter := tv.MapRange()
	for iter.Next() {
		// 通过迭代器获取 Key
		k := iter.Key()
		// 通过迭代器获取 Value
		v := iter.Value()

		// 获取 Key 的实际值和 Value 的实际值
		assert.Equal(t, v.Interface(), obj.(map[string]any)[k.String()])
	}
}

// 测试通过反射设置变量值
func TestReflect_SetValueByReflect(t *testing.T) {
	n := 100

	// 将变量 n 的值设置为 200
	err := SetValueByReflect(&n, 200)
	assert.Nil(t, err)
	assert.Equal(t, 200, n)

	s := "Hello"

	// 将变量 s 的值设置为 "OK"
	err = SetValueByReflect(&s, "OK")
	assert.Nil(t, err)
	assert.Equal(t, "OK", s)
}

// 测试函数反射的简单加法函数
func Add(a, b int) (r int) {
	r = a + b
	return
}

// 测试通过反射调用函数
func TestReflect_Call(t *testing.T) {
	// 将函数作为变量赋值给 interface{} 类型变量
	var f any = Add

	// 获取函数变量类型
	tp := reflect.TypeOf(f)
	assert.Equal(t, ".[func]", GetFullTypeName(tp))

	// 获取函数变量的反射值
	tv := reflect.ValueOf(f)

	// 构建调用函数的参数列表
	args := []reflect.Value{
		reflect.ValueOf(10),
		reflect.ValueOf(20),
	}

	// 通过反射调用函数, 获取返回值结果, 是一个 `reflect.Value` 类型的切片, 表示一到多个返回值
	r := tv.Call(args)

	// 校验返回值结果
	assert.Len(t, r, 1)
	assert.Equal(t, 30, r[0].Interface().(int))
}
