package reflect_test

import (
	"reflect"
	"study/basic/builtin/reflects"
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
	assert.Equal(t, ".int[int]", reflects.GetFullTypeName(tv.Type()))
	assert.Equal(t, 100, int(tv.Int()))

	obj = &User{Id: 1, Name: "Alvin", Gender: 'M'}

	// obj 为结构体指针, 获取指针类型变量的值实例
	tv = reflect.ValueOf(obj)
	assert.Equal(t, ".[ptr]", reflects.GetFullTypeName(tv.Type()))
	assert.Equal(t, uintptr(unsafe.Pointer(obj.(*User))), tv.Pointer())

	// 获取指针指向的实例值实例
	tv = tv.Elem()
	assert.Equal(t, "study/basic/builtin/reflects/reflect_test.User[struct]", reflects.GetFullTypeName(tv.Type()))
	assert.Equal(t, 1, int(tv.FieldByName("Id").Int()))
	assert.Equal(t, "Alvin", tv.FieldByName("Name").String())
	assert.Equal(t, 'M', rune(tv.FieldByName("Gender").Int()))

	obj = User{Id: 1, Name: "Alvin", Gender: 'M'}

	// obj 为 interface{} 指针, 获取类型为指针类型
	tv = reflect.ValueOf(&obj)
	assert.Equal(t, ".[ptr]", reflects.GetFullTypeName(tv.Type()))

	// 指针的值为实例地址
	assert.Equal(t, uintptr(unsafe.Pointer(&obj)), tv.Pointer())

	// 获取指针指向的实例值
	tv = tv.Elem()
	assert.Equal(t, ".[interface]", reflects.GetFullTypeName(tv.Type()))

	// 其值为 User 实例
	assert.Equal(t, obj.(User), tv.Interface().(User))

	// 再次从 interface{} 类型解除引用, 获取其原始值
	tv = tv.Elem()
	assert.Equal(t, "study/basic/builtin/reflects/reflect_test.User[struct]", reflects.GetFullTypeName(tv.Type()))
	assert.Equal(t, 1, int(tv.FieldByName("Id").Int()))
	assert.Equal(t, "Alvin", tv.FieldByName("Name").String())
	assert.Equal(t, 'M', rune(tv.FieldByName("Gender").Int()))
}

// 通过反射读取 "切片" 实例值
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
	assert.Equal(t, ".[slice]", reflects.GetFullTypeName(tv.Type()))

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

// 通过反射读取 "字典" 实例值
func TestReflect_ValueOfMap(t *testing.T) {
	// 定义字典类型变量, 类型转为 any 类型
	var obj any = map[string]any{"a": 1, "b": "Hello", "c": false}

	// 获取实例的 值反射 结果
	tv := reflect.ValueOf(obj)
	assert.Equal(t, ".[map]", reflects.GetFullTypeName(tv.Type()))

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
	err := reflects.SetValueByReflect(&n, 200)
	assert.Nil(t, err)
	assert.Equal(t, 200, n)

	s := "Hello"

	// 将变量 s 的值设置为 "OK"
	err = reflects.SetValueByReflect(&s, "OK")
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
	var fn any = Add

	// 获取函数变量类型
	tp := reflect.TypeOf(fn)
	assert.Equal(t, ".[func]", reflects.GetFullTypeName(tp))

	// 获取函数变量的反射值
	tv := reflect.ValueOf(fn)

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
