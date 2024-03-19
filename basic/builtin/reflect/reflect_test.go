package reflect

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// 表示用户的结构体
type user struct {
	Id     int
	Name   string
	Gender rune
}

// 测试函数反射的简单加法函数
func Add(a, b int) (r int) {
	r = a + b
	return
}

// 通过反射获取对象类型
//
// `reflect.TypeOf` 用于获取一个变量 (`interface{}` 类型) 的 类型反射
func TestReflectGetType(t *testing.T) {
	// 定义 interface{} 类型变量, 实际类型为 int64 类型
	var obj interface{} = int64(100)

	tp := reflect.TypeOf(obj) // 获取变量的 类型反射 对象
	assert.Equal(             // 变量的类型是 "int64"
		t,
		".int64[int64]",
		GetFullTypeName(tp),
	)

	// 定义 interface{} 类型变量, 实际类型为 User 类型
	obj = user{}
	tp = reflect.TypeOf(obj) // 获取变量的 类型反射 对象
	assert.Equal(            // 变量的类型是 `user`
		t,
		"study-golang/basic/builtin/reflect.user[struct]",
		GetFullTypeName(tp),
	)

	// 对于 指针 类型, 类型的名字为 "", 种类为 `reflect.Ptr`
	tp = reflect.TypeOf(&obj)
	assert.Equal(t, ".[ptr]", GetFullTypeName(tp))

	tp = tp.Elem() // 获取指针指向对象的类型, 为 `interface{}` 类型, 此处无法获取该 `interface{}` 对象的原始类型
	assert.Equal(t, ".[interface]", GetFullTypeName(tp))

	// 对于 切片 类型, 类型的名字为 "", 种类为 `reflect.Slice`
	tp = reflect.TypeOf([]int{1, 2, 3, 4})
	assert.Equal(t, ".[slice]", GetFullTypeName(tp))

	// 对于 数组 类型, 类型的名字为 "", 种类为 `reflect.Array`
	tp = reflect.TypeOf([...]int{1, 2, 3, 4})
	assert.Equal(t, ".[array]", GetFullTypeName(tp))

	// 对于 字典 类型, 类型的名字为 "", 种类为 `reflect.Map`
	tp = reflect.TypeOf(map[string]interface{}{"a": 1, "b": "Hello"})
	assert.Equal(t, ".[map]", GetFullTypeName(tp))
}

// 通过反射获取指针类型及其指向对象的类型
//
// 若 `reflect.Type` 类型是一个指针类型, 则可以通过 `reflect.Type.Elem()` 函数获取该指针指向的对象类型
func TestReflectGetTypeFromPtr(t *testing.T) {
	n := 100

	// obj 保存指向变量的指针
	var obj interface{} = &n

	tp := reflect.TypeOf(obj) // 获取指针类型变量的类型
	assert.Equal(t, ".[ptr]", GetFullTypeName(tp))

	tp = tp.Elem() // 获取指针所指向的对象类型
	assert.Equal(t, ".int[int]", GetFullTypeName(tp))

	// obj 保存指向结构体的指针
	obj = &user{}

	tp = reflect.TypeOf(obj) // 获取指针变量的类型
	assert.Equal(t, ".[ptr]", GetFullTypeName(tp))

	tp = tp.Elem() // 获取指针所指向的对象类型
	assert.Equal(
		t,
		"study-golang/basic/builtin/reflect.user[struct]",
		GetFullTypeName(tp),
	)
}

// 通过反射读取对象值
//
// `reflect.ValueOf` 用于获取一个变量 (`interface{}` 类型) 的值反射
func TestReflectGetValue(t *testing.T) {
	// 定义 interface{} 类型变量, 值为整型
	var obj interface{} = 100

	tv := reflect.ValueOf(obj) // 获取变量的 值反射 对象
	assert.Equal(t, ".int[int]", GetFullTypeName(tv.Type()))
	assert.Equal(t, 100, int(tv.Int())) // 通过反射获取 值

	// 定义 `interface{}` 类型变量, 值为 `user` 类型结构体
	obj = user{Id: 1, Name: "Alvin", Gender: 'M'}

	tv = reflect.ValueOf(obj) // 获取变量的 值反射 对象
	assert.Equal(
		t,
		"study-golang/basic/builtin/reflect.user[struct]",
		GetFullTypeName(tv.Type()),
	)
	assert.Equal(t, 1, int(tv.FieldByName("Id").Int()))        // 根据 名称 获取 `Id` 字段的值, 并转为 `int` 类型
	assert.Equal(t, "Alvin", tv.FieldByName("Name").String())  // 根据 名称 获取 `Name` 字段的值, 并转为 `string` 类型
	assert.Equal(t, 'M', rune(tv.FieldByName("Gender").Int())) // 根据 名称 获取 `Gender` 字段的值, 并转为 `rune` 类型

	// 配合 类型反射 对象, 对 结构体 变量进行反射遍历
	names := []string{"Id", "Name", "Gender"}
	values := []interface{}{1, "Alvin", 'M'}

	tp := reflect.TypeOf(obj)
	for i := 0; i < tp.NumField(); i++ { // 获取对象字段总数
		field := tp.Field(i) // 通过 类型反射 对象, 获取第 `i` 个字段的 类型
		assert.Equal(t, names[i], field.Name)

		value := tv.Field(i) // 通过 值反射 对象, 获取第 `i` 个字段的 值
		assert.EqualValues(  // 将所有字段值都获取为 `interface{}` 类型
			t,
			values[i],
			value.Interface(),
		)
	}
}

// 通过反射读取指针及其指向的对象值
//
// 若 `reflect.Value` 引用了一个指针类型值, 则可通过 `reflect.Value.Elem()` 函数获取其指向对象的值
func TestReflectGetValueFromPtr(t *testing.T) {
	n := 100
	// obj 为变量指针
	var obj interface{} = &n

	tv := reflect.ValueOf(obj) // 获取指针类型变量的值对象
	assert.Equal(              // 其值为一个地址
		t,
		uintptr(unsafe.Pointer(&n)),
		tv.Pointer(),
	)
	tv = tv.Elem() // 获取指针指向对象的值对象
	assert.Equal(  // 其类型为 int 类型
		t,
		".int[int]",
		GetFullTypeName(tv.Type()),
	)
	assert.Equal(t, 100, int(tv.Int())) // 获取值

	obj = &user{Id: 1, Name: "Alvin", Gender: 'M'}
	// obj 为结构体指针
	tv = reflect.ValueOf(obj) // 获取指针类型变量的值对象
	assert.Equal(             // 其类型为一个指针类型
		t,
		".[ptr]",
		GetFullTypeName(tv.Type()),
	)
	assert.Equal( // 其值为一个地址, 指向 User 对象
		t,
		uintptr(unsafe.Pointer(obj.(*user))),
		tv.Pointer(),
	)

	tv = tv.Elem() // 获取指针指向的对象值对象
	assert.Equal(  // 其类型为结构体类型
		t,
		"study-golang/basic/builtin/reflect.user[struct]",
		GetFullTypeName(tv.Type()),
	)
	assert.Equal(t, 1, int(tv.FieldByName("Id").Int())) // 获取结构体各字段值
	assert.Equal(t, "Alvin", tv.FieldByName("Name").String())
	assert.Equal(t, 'M', rune(tv.FieldByName("Gender").Int()))

	obj = user{Id: 1, Name: "Alvin", Gender: 'M'}
	// obj 为 interface{} 指针
	tv = reflect.ValueOf(&obj)
	assert.Equal(t, ".[ptr]", GetFullTypeName(tv.Type()))        // 获取类型为指针类型
	assert.Equal(t, uintptr(unsafe.Pointer(&obj)), tv.Pointer()) // 指针的值为对象地址

	tv = tv.Elem()                                              // 获取指针指向的对象值
	assert.Equal(t, ".[interface]", GetFullTypeName(tv.Type())) // 其类型为 interface{} 类型
	assert.Equal(t, obj.(user), tv.Interface().(user))          // 其值为 User 对象

	tv = tv.Elem() // 再次从 interface{} 类型解除引用, 获取其原始值
	assert.Equal(  // 获取类型为 User 类型
		t,
		"study-golang/basic/builtin/reflect.user[struct]",
		GetFullTypeName(tv.Type()),
	)
	assert.Equal(t, 1, int(tv.FieldByName("Id").Int())) // 获取对象各字段值
	assert.Equal(t, "Alvin", tv.FieldByName("Name").String())
	assert.Equal(t, 'M', rune(tv.FieldByName("Gender").Int()))
}

// 通过反射读取"切片"对象值
//
// 若 `reflect.Value` 对象引用了一个切片类型值, 则可以通过 `reflect.Value` 对象提供的一组切片反射方法对其进行操作, 包括:
//   - `reflect.Value.Len()`
//   - `reflect.Value.Index(n)`
//   - `reflect.Value.Slice(m, n)`
func TestReflectGetValueFromSlice(t *testing.T) {
	// 通过反射操作切片
	// 定义一个切片对象, 类型转为 interface{} 类型
	var obj interface{} = []interface{}{1, "Hello", false}

	tv := reflect.ValueOf(obj) // 获取对象的 值反射 结果
	assert.Equal(t, ".[slice]", GetFullTypeName(tv.Type()))

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

// 通过反射读取"字典"对象值
func TestReflectGetValueFromMap(t *testing.T) {
	// 通过反射操作 字典
	// 定义字典类型变量, 类型转为 interface{} 类型
	var obj interface{} = map[string]interface{}{"a": 1, "b": "Hello", "c": false}

	tv := reflect.ValueOf(obj) // 获取对象的 值反射 结果
	assert.Equal(t, ".[map]", GetFullTypeName(tv.Type()))

	keys := tv.MapKeys()   // 通过反射获取字典所有的 key 的集合
	assert.Len(t, keys, 3) // 获取 key 集合的长度
	assert.ElementsMatch(  // 获取每个 key 的值
		t,
		[]string{"a", "b", "c"},
		[]string{keys[0].String(),
			keys[1].String(),
			keys[2].String()},
	)

	val := tv.MapIndex(reflect.ValueOf("a")) // 通过反射, 根据 key 的值获取 value, 注意, 这里的 key 必须是"值反射"对象
	assert.Equal(t, 1, val.Interface().(int))

	val = tv.MapIndex(reflect.ValueOf("b"))
	assert.Equal(t, "Hello", val.Interface().(string))

	val = tv.MapIndex(reflect.ValueOf("c"))
	assert.Equal(t, false, val.Interface().(bool))

	iter := tv.MapRange() // 获取 key/value 对的迭代器
	for iter.Next() {
		k := iter.Key()   // 通过迭代器获取 key
		v := iter.Value() // 通过迭代器获取 value
		assert.Equal(     // 获取 key 的实际值 和 value 的实际值
			t,
			v.Interface(),
			obj.(map[string]interface{})[k.String()],
		)
	}
}

// 测试通过反射设置变量值
func TestSetValueByReflectFunc(t *testing.T) {
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

// 测试函数调用的反射
func TestReflectFunction(t *testing.T) {
	// 将函数作为变量赋值给 interface{} 类型变量
	var f interface{} = Add

	// 获取函数变量类型
	tp := reflect.TypeOf(f)
	assert.Equal(t, ".[func]", GetFullTypeName(tp)) // PkgPath, Name 都为空, Kind 为 func

	// 获取函数变量的反射值
	tv := reflect.ValueOf(f)
	args := []reflect.Value{ // 构建调用函数的参数列表
		reflect.ValueOf(10),
		reflect.ValueOf(20),
	}
	// 通过反射调用函数, 获取返回值结果, 是一个 reflect.Value 类型的 slice
	r := tv.Call(args)

	// 校验返回值结果
	assert.Len(t, r, 1)
	assert.Equal(t, 30, r[0].Interface().(int))
}
