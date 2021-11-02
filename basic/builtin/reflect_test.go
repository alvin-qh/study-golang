package builtin

import (
	"basic/builtin/types"
	"fmt"
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// 通过反射获取对象类型
// reflect.TypeOf 用于获取一个变量 (interface{} 类型) 的 类型反射
func TestReflectGetType(t *testing.T) {
	// 定义 interface{} 类型变量，实际类型为 int64 类型
	var obj interface{} = int64(100)

	tp := reflect.TypeOf(obj)                                   // 获取变量的 类型反射 对象
	assert.Equal(t, ".int64[int64]", types.GetFullTypeName(tp)) // 变量的类型是 "int64"

	// 定义 interface{} 类型变量，实际类型为 types.User 类型
	obj = types.User{Id: 1, Name: "Alvin", Gender: 'M'}

	tp = reflect.TypeOf(obj)                                                       // 获取变量的 类型反射 对象
	assert.Equal(t, "basic/builtin/types.User[struct]", types.GetFullTypeName(tp)) // 变量的类型是 "User"

	// 对于 指针 类型，类型的名字为 "", 种类为 reflect.Ptr
	tp = reflect.TypeOf(&obj)
	assert.Equal(t, ".[ptr]", types.GetFullTypeName(tp))

	tp = tp.Elem() // 获取指针指向对象的类型, 为 interface{} 类型，此处无法获取该 interface{} 对象的原始类型
	assert.Equal(t, ".[interface]", types.GetFullTypeName(tp))

	// 对于 切片 类型，类型的名字为 "", 种类为 reflect.Slice
	tp = reflect.TypeOf([]int{1, 2, 3, 4})
	assert.Equal(t, ".[slice]", types.GetFullTypeName(tp))

	// 对于 数组 类型，类型的名字为 "", 种类为 reflect.Array
	tp = reflect.TypeOf([...]int{1, 2, 3, 4})
	assert.Equal(t, ".[array]", types.GetFullTypeName(tp))

	// 对于 字典 类型，类型的名字为 "", 种类为 reflect.Map
	tp = reflect.TypeOf(map[string]interface{}{"a": 1, "b": "Hello"})
	assert.Equal(t, ".[map]", types.GetFullTypeName(tp))
}

// 通过反射获取指针类型及其指向对象的类型
func TestReflectGetTypeFromPtr(t *testing.T) {
	n := 100
	var obj interface{} = &n

	tp := reflect.TypeOf(obj)
	assert.Equal(t, ".[ptr]", types.GetFullTypeName(tp))

	tp = tp.Elem()
	assert.Equal(t, ".int[int]", types.GetFullTypeName(tp))

	u := types.User{Id: 1, Name: "Alvin", Gender: 'M'}
	obj = &u

	tp = reflect.TypeOf(obj)
	assert.Equal(t, ".[ptr]", types.GetFullTypeName(tp))

	tp = tp.Elem()
	assert.Equal(t, "basic/builtin/types.User[struct]", types.GetFullTypeName(tp))
}

// 通过反射读取对象值
// reflect.ValueOf 用于获取一个变量 (interface{} 类型) 的 值反射
func TestReflectGetValue(t *testing.T) {
	// 定义 interface{} 类型变量，值为整型
	var obj interface{} = 100
	tv := reflect.ValueOf(obj) // 获取变量的 值反射 对象
	assert.Equal(t, ".int[int]", types.GetFullTypeName(tv.Type()))
	assert.Equal(t, 100, int(tv.Int())) // 通过反射获取 值

	// 定义 interface{} 类型变量，值为 User 类型结构体
	obj = types.User{Id: 1, Name: "Alvin", Gender: 'M'}
	tv = reflect.ValueOf(obj) // 获取变量的 值反射 对象
	assert.Equal(t, "basic/builtin/types.User[struct]", types.GetFullTypeName(tv.Type()))
	assert.Equal(t, 1, int(tv.FieldByName("Id").Int()))        // 根据 名称 获取 Id 字段的值，并转为 int 类型
	assert.Equal(t, "Alvin", tv.FieldByName("Name").String())  // 根据 名称 获取 Name 字段的值，并转为 string 类型
	assert.Equal(t, 'M', rune(tv.FieldByName("Gender").Int())) // 根据 名称 获取 Gender 字段的值，并转为 rune 类型

	// 配合 类型反射 对象，对 结构体 变量进行反射遍历
	names := []string{"Id", "Name", "Gender"}
	values := []interface{}{1, "Alvin", 'M'}

	tp := reflect.TypeOf(obj)
	for i := 0; i < tp.NumField(); i++ { // 获取对象字段总数
		field := tp.Field(i) // 通过 类型反射 对象，获取第 i 个字段的 类型
		assert.Equal(t, names[i], field.Name)

		value := tv.Field(i)                                // 通过 值反射 对象，获取第 i 个字段的 值
		assert.EqualValues(t, values[i], value.Interface()) // 将所有字段值都获取为 interface{} 类型
	}
}

// 通过反射读取指针及其指向的对象值
func TestReflectGetValueFromPtr(t *testing.T) {
	n := 100
	var obj interface{} = &n

	tv := reflect.ValueOf(obj)
	assert.Equal(t, uintptr(unsafe.Pointer(&n)), tv.Pointer())

	tv = tv.Elem()
	assert.Equal(t, 100, int(tv.Int()))

	u := types.User{Id: 1, Name: "Alvin", Gender: 'M'}
	obj = &u

	tv = reflect.ValueOf(obj)
	assert.Equal(t, ".[ptr]", types.GetFullTypeName(tv.Type()))
	assert.Equal(t, uintptr(unsafe.Pointer(&u)), tv.Pointer())

	tv = tv.Elem()
	assert.Equal(t, "basic/builtin/types.User[struct]", types.GetFullTypeName(tv.Type()))
	assert.Equal(t, 1, int(tv.FieldByName("Id").Int()))
	assert.Equal(t, "Alvin", tv.FieldByName("Name").String())
	assert.Equal(t, 'M', rune(tv.FieldByName("Gender").Int()))

	obj = u
	tv = reflect.ValueOf(&obj)
	assert.Equal(t, ".[ptr]", types.GetFullTypeName(tv.Type()))
	assert.Equal(t, uintptr(unsafe.Pointer(&obj)), tv.Pointer())

	tv = tv.Elem()
	assert.Equal(t, ".[interface]", types.GetFullTypeName(tv.Type()))
	assert.Equal(t, u, tv.Interface().(types.User))

	tv = tv.Elem()
	assert.Equal(t, "basic/builtin/types.User[struct]", types.GetFullTypeName(tv.Type()))
	assert.Equal(t, 1, int(tv.FieldByName("Id").Int()))
	assert.Equal(t, "Alvin", tv.FieldByName("Name").String())
	assert.Equal(t, 'M', rune(tv.FieldByName("Gender").Int()))
}

// 通过反射读取 切片 对象值
func TestReflectGetValueFromMap(t *testing.T) {
	// 通过反射操作 字典
	// 定义字典类型变量，类型转为 interface{} 类型
	var obj interface{} = map[string]interface{}{"a": 1, "b": "Hello", "c": false}

	tv := reflect.ValueOf(obj) // 获取对象的 值反射 结果
	assert.Equal(t, ".[map]", types.GetFullTypeName(tv.Type()))

	keys := tv.MapKeys()                   // 通过反射获取字典所有的 key 的集合
	assert.Len(t, keys, 3)                 // 获取 key 集合的长度
	assert.Equal(t, "a", keys[0].String()) // 获取每个 key 的值
	assert.Equal(t, "b", keys[1].String())
	assert.Equal(t, "c", keys[2].String())

	val := tv.MapIndex(reflect.ValueOf("a")) // 通过反射，根据 key 的值获取 value，注意，这里的 key 必须是 值反射 对象
	assert.Equal(t, 1, val.Interface().(int))

	val = tv.MapIndex(reflect.ValueOf("b"))
	assert.Equal(t, "Hello", val.Interface().(string))

	val = tv.MapIndex(reflect.ValueOf("c"))
	assert.Equal(t, false, val.Interface().(bool))

	iter := tv.MapRange() // 获取 key/value 对的迭代器
	for iter.Next() {
		k := iter.Key()                                                          // 通过迭代器获取 key
		v := iter.Value()                                                        // 通过迭代器获取 value
		assert.Equal(t, v.Interface(), obj.(map[string]interface{})[k.String()]) // 获取 key 的实际值 和 value 的实际值
	}
}

// 通过反射读取 字典 对象值
func TestReflectGetValueFromSlice(t *testing.T) {
	// 通过反射操作切片
	// 定义一个切片对象，类型转为 interface{} 类型
	var obj interface{} = []interface{}{1, "Hello", false}

	tv := reflect.ValueOf(obj) // 获取对象的 值反射 结果
	assert.Equal(t, ".[slice]", types.GetFullTypeName(tv.Type()))

	len := tv.Len() // 通过 值反射 对象获取切片长度
	assert.Equal(t, 3, len)

	val := tv.Index(0) // 获取切片指定下标的值
	assert.Equal(t, 1, val.Interface().(int))

	val = tv.Index(1)
	assert.Equal(t, "Hello", val.Interface().(string))

	val = tv.Index(2)
	assert.Equal(t, false, val.Interface().(bool))

	tv = tv.Slice(0, 2) // 通过 值反射 对象进行切片操作

	len = tv.Len() // 获取切片长度
	assert.Equal(t, 2, len)

	val = tv.Index(0) // 获取切片指定下标的值
	assert.Equal(t, 1, val.Interface().(int))

	val = tv.Index(1)
	assert.Equal(t, "Hello", val.Interface().(string))
}

func TestReflectSetValue(t *testing.T) {
	var obj interface{} = 100

	tv := reflect.ValueOf(&obj)
	assert.Equal(t, ".[ptr]", types.GetFullTypeName(tv.Type()))

	tv = tv.Elem()
	assert.Equal(t, ".[interface]", types.GetFullTypeName(tv.Type()))

	ptr := tv.Addr().Interface().(*interface{})
	*ptr = 200
	assert.Equal(t, obj, 200)

	obj = 100
	tv = reflect.ValueOf(&obj)
	assert.Equal(t, ".[ptr]", types.GetFullTypeName(tv.Type()))
	assert.Equal(t, 100, int(tv.Elem().Interface().(int)))

	tv = tv.Elem()
	assert.Equal(t, ".[interface]", types.GetFullTypeName(tv.Type()))

	tv.Set(reflect.ValueOf(200))
	assert.Equal(t, 200, obj)

	n := 100
	obj = &n
	tv = reflect.ValueOf(&obj)
	assert.Equal(t, ".[ptr]", types.GetFullTypeName(tv.Type()))

	tv = tv.Elem()
	assert.Equal(t, ".[interface]", types.GetFullTypeName(tv.Type()))

	obj = types.User{Id: 1, Name: "Alvin", Gender: 'M'}

	tv = reflect.ValueOf(&obj)
	assert.Equal(t, reflect.Ptr, tv.Type().Kind()) // 确认是指针类型

	tv = tv.Elem()
	assert.Equal(t, reflect.Interface, tv.Type().Kind()) // 确认是 interface{} 类型

	tv = tv.Elem()
	assert.Equal(t, reflect.Struct, tv.Type().Kind())                                                        // 确认是 interface{} 类型
	assert.Equal(t, "basic/builtin/types.User", fmt.Sprintf("%v.%v", tv.Type().PkgPath(), tv.Type().Name())) // 确认是 interface{} 类型

	// tmp := tv.Elem()
	// tmp = tv.Elem()
	// assert.Equal(t, "Alvin", tmp.FieldByName("Name").String())

	// tmp = reflect.New(tv.Elem().Elem().Type())
	// tmp.Elem().Set(tv.Elem().Elem())
	// tmp.Elem().FieldByName("Name").SetString("Emma")
	// tv.Elem().Set(tmp.Elem())
}
