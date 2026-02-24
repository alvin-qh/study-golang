package types_test

import (
	"reflect"
	"study/basic/builtin/reflects"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试反射用结构体
type User struct {
	Id     int
	Name   string
	Gender rune
}

// 通过反射获取变量的类型信息
//
// 通过 `reflect.TypeOf` 函数可以获取变量的类型信息, 返回 `reflect.Type` 类型实例, 基本信息包括:
// - 名称: `reflect.Type.Name()`: 获取变量所属类型的名称, 对于指针 (Pointer), 切片 (Slice), 字典 (Map) 等类型则返回空字符串
// - 类型: `reflect.Type.Kind()`: 获取变量所属类型的类型, 返回 `reflect.Kind` 类型枚举值, 参见 `reflect.Kind` 类型
// - 包路径: `reflect.Type.PkgPath()`: 获取变量所属类型所在的包路径, 对于内置类型则返回空字符串
func TestType_TypeOf(t *testing.T) {
	// 定义一个 any 类型变量
	var obj any

	// 令 any 类型变量存储整型值
	obj = int64(100)

	// 获取 any 变量的类型反射实例
	tp := reflect.TypeOf(obj)

	// 确认 any 变量的实际类型名称为 `int64`, 类型为 `reflect.Int64`, 包路径为空字符串 (内置类型)
	assert.Equal(t, "int64", tp.Name())
	assert.Equal(t, reflect.Int64, tp.Kind())
	assert.Equal(t, "", tp.PkgPath())

	// 令 any 变量存储指针类型值
	obj = &obj

	// 获取 any 变量的类型反射实例
	tp = reflect.TypeOf(obj)

	// 确认 any 变量的实际类型名称为空字符串 (指针类型), 类型为 `reflect.Pointer`, 包路径为空字符串 (内置类型)
	assert.Equal(t, "", tp.Name())
	assert.Equal(t, reflect.Pointer, tp.Kind())
	assert.Equal(t, "", tp.PkgPath())

	// 令 any 变量存储切片类型值
	obj = make([]string, 0)

	// 获取 any 变量的类型反射实例
	tp = reflect.TypeOf(obj)

	// 确认 any 变量的实际类型名称为空字符串 (切片类型), 类型为 `reflect.Slice`, 包路径为空字符串 (内置类型)
	assert.Equal(t, "", tp.Name())
	assert.Equal(t, reflect.Slice, tp.Kind())
	assert.Equal(t, "", tp.PkgPath())

	// 令 any 变量存储字典类型值
	obj = make(map[string]string, 0)

	// 获取 any 变量的类型反射实例
	tp = reflect.TypeOf(obj)

	// 确认 any 变量的实际类型名称为空字符串 (字典类型), 类型为 `reflect.Map`, 包路径为空字符串 (内置类型)
	assert.Equal(t, "", tp.Name())
	assert.Equal(t, reflect.Map, tp.Kind())
	assert.Equal(t, "", tp.PkgPath())

	// 令 any 存储自定义结构体类型值
	obj = User{}

	// 获取 any 变量的类型反射实例
	tp = reflect.TypeOf(obj)

	// 确认 any 变量的实际类型名称为结构体类型名称, 类型为 `reflect.Struct`, 包路径为结构体所在包路径
	assert.Equal(t, "User", tp.Name())
	assert.Equal(t, reflect.Struct, tp.Kind())
	assert.Equal(t, "study/basic/builtin/reflects/types_test", tp.PkgPath())
}

// 测试 `reflects.GetFullTypeName` 函数, 获取变量所属类型的限定名
func TestType_GetFullTypeName(t *testing.T) {
	// 定义 any 类型变量
	var obj any

	// 令 any 变量存储整型值, 确认其类型限定名称
	obj = int64(100)
	assert.Equal(t, ".int64[int64]", reflects.GetValueFullTypeName(obj))

	// 令 any 存储指针类型值, 确认其类型限定名称
	obj = &obj
	assert.Equal(t, ".[ptr]", reflects.GetValueFullTypeName(obj))

	// 令 any 变量存储切片类型值, 确认其类型限定名称
	obj = make([]string, 0)
	assert.Equal(t, ".[slice]", reflects.GetValueFullTypeName(obj))

	// 令 any 变量存储字典类型值
	obj = make(map[string]string, 0)
	assert.Equal(t, ".[map]", reflects.GetValueFullTypeName(obj))

	// 令 any 存储结构体类型值, 确认其限定名称
	obj = User{}
	assert.Equal(t, "study/basic/builtin/reflects/types_test.User[struct]", reflects.GetValueFullTypeName(obj))
}

// 通过泛型获取类型名称
//
// 通过 `reflect.TypeFor[T]()` 函数可以通过泛型方式获取类型 `T` 的类型对象
//
// 通过泛型的方法, 可有效的减少反射带来的性能损失
func TestType_TypeFor(t *testing.T) {
	tp := reflect.TypeFor[int]()
	assert.Equal(t, ".int[int]", reflects.GetFullTypeName(tp))

	// 获取指针变量的反射实例
	tp = reflect.TypeFor[*any]()
	assert.Equal(t, ".[ptr]", reflects.GetFullTypeName(tp))

	// 获取指针变量的反射实例
	tp = tp.Elem()
	assert.Equal(t, ".[interface]", reflects.GetFullTypeName(tp))

	// 获取切片变量的反射实例
	tp = reflect.TypeFor[[]int]()
	assert.Equal(t, ".[slice]", reflects.GetFullTypeName(tp))

	// 获取数组变量的反射实例
	tp = reflect.TypeFor[[4]int]()
	assert.Equal(t, ".[array]", reflects.GetFullTypeName(tp))

	// 获取 Map 变量的反射实例
	tp = reflect.TypeFor[map[string]any]()
	assert.Equal(t, ".[map]", reflects.GetFullTypeName(tp))

	// 获取结构体变量的反射实例
	tp = reflect.TypeFor[User]()
	assert.Equal(t, "study/basic/builtin/reflects/types_test.User[struct]", reflects.GetFullTypeName(tp))
}

// 通过反射获取指针类型及其指向实例的类型
//
// 若 `reflect.Type` 类型是一个指针类型, 则可以通过 `reflect.Type.Elem()` 函数获取该指针指向的实例类型
func TestType_ElemOfPointerType(t *testing.T) {
	// 定义一个整型变量
	n := 100

	// obj 保存指向变量的指针
	var obj any = &n

	// 获取 any 变量的实际类型, 确认变量的类型为指针类型
	tp := reflect.TypeOf(obj)
	assert.Equal(t, reflect.Pointer, tp.Kind())

	// 获取指针所指向的值的实际类型, 确认其类型即类型名称
	tp = tp.Elem()
	assert.Equal(t, reflect.Int, tp.Kind())
	assert.Equal(t, "int", tp.Name())

	// obj 保存指向结构体的指针
	obj = &User{}

	// 获取指针变量的类型, 确认其类型为指针类型
	tp = reflect.TypeOf(obj)
	assert.Equal(t, reflect.Pointer, tp.Kind())

	// 获取指针所指向的实例类型, 确认其类型即类型名称
	tp = tp.Elem()
	assert.Equal(t, reflect.Struct, tp.Kind())
	assert.Equal(t, "User", tp.Name())
}
