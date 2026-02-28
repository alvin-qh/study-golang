package values_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Date struct {
	Year  int
	Month int
	Day   int
}

// 测试用结构体
type User struct {
	Id       int
	Name     string
	Gender   rune
	Birthday Date
}

func TestReflect_ValueOf(t *testing.T) {
	var obj any

	// 测试通过反射, 尝试从 any 变量获取其存储的整型值
	t.Run("For Integer Value", func(t *testing.T) {
		// 为 any 类型变量赋予整型值 100
		obj = 100

		// 从 any 变量通过反射获取反射值对象
		tv := reflect.ValueOf(obj)

		// 确认反射值对象的值类型为 int
		assert.Equal(t, reflect.TypeFor[int](), tv.Type())

		// 确认反射值对象有效
		assert.True(t, tv.IsValid())

		// 确认反射值对象的值为整型
		assert.True(t, tv.CanInt())

		// 确认通过反射值获取变量本身的整数值为 100, 因为变量本身值为整数
		assert.Equal(t, int64(100), tv.Int())
	})

	// 测试通过反射, 尝试从 any 变量获取其存储的整型值
	t.Run("For Float Value", func(t *testing.T) {
		// 为 any 类型变量赋予整型值 100
		obj = 100.001

		// 从 any 变量通过反射获取反射值对象
		tv := reflect.ValueOf(obj)

		// 确认反射值对象的值类型为 float64
		assert.Equal(t, reflect.TypeFor[float64](), tv.Type())

		// 确认反射值对象有效
		assert.True(t, tv.IsValid())

		// 确认反射值对象的值不为零值
		assert.False(t, tv.IsZero())

		// 确认反射值对象的值为浮点型
		assert.True(t, tv.CanFloat())

		// 确认通过反射值获取变量本身的值为 100.001, 因为变量本身值为浮点数
		assert.Equal(t, 100.001, tv.Float())
	})

	// 测试通过反射, 判断一个变量是否为 0 值
	//
	// 无论是整型还是浮点型, 都可以通过反射值对象的 `.IsZero()` 方法判断一个该反射值对应的变量是否为 0 值
	t.Run("Check If Zero", func(t *testing.T) {
		// 将 any 类型变量设置为整数 0 值
		obj = 0

		// 从 any 变量通过反射获取反射值对象, 确认反射值对象表示整型, 且其值为零值
		tv := reflect.ValueOf(obj)
		assert.True(t, tv.CanInt())
		assert.True(t, tv.IsZero())

		// 将 any 类型变量设置为浮点型 0 值
		obj = 0.0

		// 从 any 变量通过反射获取反射值对象, 确认反射值对象表示浮点型, 且其值为零值
		tv = reflect.ValueOf(obj)
		assert.True(t, tv.CanFloat())
		assert.True(t, tv.IsZero())
	})

	// 测试通过反射, 尝试从 any 变量获取其存储的字符串值
	//
	// 注意: 任何类型变量的反射值对象, 都可以通过 `.String()` 方法获取该反射值对应的变量的字符串表示
	t.Run("For String Type", func(t *testing.T) {
		// 为 any 类型变量设置字符串值
		obj = "hello"

		// 从 any 变量通过反射获取反射值对象, 确认该反射值对象表示字符串, 且其值为 "hello"
		tv := reflect.ValueOf(obj)
		assert.Equal(t, reflect.TypeFor[string](), tv.Type())
		assert.Equal(t, "hello", tv.String())

		// 为 any 类型变量设置整型值
		obj = 100

		// 从 any 变量通过反射获取反射值对象, 并获取其字符串值, 对应整型变量的反射值对象, 确认其字符串值是 "<int Value>"
		tv = reflect.ValueOf(obj)
		assert.Equal(t, "<int Value>", tv.String())

		// 为 any 类型变量设置切片对象
		obj = []string{"One", "Two", "Three"}

		// 从 any 变量通过反射获取反射值对象, 并获取其字符串值, 对应切片变量的反射值对象, 确认字符串值为 "<[]string Value>"
		tv = reflect.ValueOf(obj)
		assert.Equal(t, "<[]string Value>", tv.String())
	})

	// 测试通过反射, 尝试从 any 变量获取其存储的切片对象
	t.Run("For Slice Type", func(t *testing.T) {
		// 为 any 类型变量设置切片对象
		obj = []string{"One", "Two", "Three"}

		// 从 any 变量通过反射获取反射值对象, 确认该反射值对象表示切片
		tv := reflect.ValueOf(obj)
		assert.Equal(t, reflect.TypeFor[[]string](), tv.Type())

		// 确认通过反射值获取切片的长度值为 3
		assert.Equal(t, 3, tv.Len())

		// 获取切片的第 1 个元素, 确认该元素为字符串 "One"
		tiv := tv.Index(0)
		assert.Equal(t, "One", tiv.String())

		// 获取切片的第 2 个元素, 确认该元素为字符串 "Two"
		tiv = tv.Index(1)
		assert.Equal(t, "Two", tiv.String())

		// 获取切片的第 3 个元素, 确认该元素为字符串 "Three"
		tiv = tv.Index(2)
		assert.Equal(t, "Three", tiv.String())

		// 获取切片的第 4 个元素, 确认会 Panic 错误, 表示数组下标越界
		assert.PanicsWithValue(t, "reflect: slice index out of range", func() { tv.Index(3) })
	})

	// 测试通过反射, 尝试从 any 变量获取其存储的结构体对象
	t.Run("For Struct Type", func(t *testing.T) {
		// 为 any 类型变量设置结构体对象
		obj = User{Id: 1, Name: "Alvin", Gender: 'M', Birthday: Date{1981, 3, 17}}

		// 获取结构体的类型信息
		tt := reflect.TypeFor[User]()

		// 从 any 变量通过反射获取反射值对象, 确认该反射值对象表示切片
		tv := reflect.ValueOf(obj)
		assert.Equal(t, reflect.TypeFor[User](), tv.Type())

		// 通过 .Field(n) 方法获取结构体的字段值, 确认字段值正确
		assert.Equal(t, "Id", tt.Field(0).Name)
		assert.Equal(t, int64(1), tv.Field(0).Int())

        // 通过 .Field(n) 方法获取结构体的字段值, 确认字段值正确
		assert.Equal(t, "Alvin", tv.Field(1).String())
		assert.Equal(t, 'M', tv.Field(2).Interface().(rune))

		// 通过 .Field(m).Field(n) 方法获取嵌套结构体的字段值, 确认字段值正确
		assert.Equal(t, int64(1981), tv.Field(3).Field(0).Int())
		assert.Equal(t, int64(3), tv.Field(3).Field(1).Int())
		assert.Equal(t, int64(17), tv.Field(3).Field(2).Int())

		// 通过 .FieldByName(name) 方法获取结构体的字段值, 确认字段值正确
		assert.Equal(t, int64(1), tv.FieldByName("Id").Int())
		assert.Equal(t, "Alvin", tv.FieldByName("Name").String())
		assert.Equal(t, 'M', tv.FieldByName("Gender").Interface().(rune))

		// 通过 .FieldByName(name1).FieldByName(name2) 方法获取嵌套结构体的字段值, 确认字段值正确
		assert.Equal(t, int64(1981), tv.FieldByName("Birthday").FieldByName("Year").Int())
		assert.Equal(t, int64(3), tv.FieldByName("Birthday").FieldByName("Month").Int())
		assert.Equal(t, int64(17), tv.FieldByName("Birthday").FieldByName("Day").Int())

		// 通过 .FieldByIndex({n}) 方法获取结构体的字段值, 确认字段值正确
		assert.Equal(t, int64(1), tv.FieldByIndex([]int{0}).Int())
		assert.Equal(t, "Alvin", tv.FieldByIndex([]int{1}).String())
		assert.Equal(t, 'M', tv.FieldByIndex([]int{2}).Interface().(rune))

		// 通过 .FieldByIndex({n1, n2}) 方法获取结构体的字段值, 确认字段值正确
		assert.Equal(t, int64(1981), tv.FieldByIndex([]int{3, 0}).Int())
		assert.Equal(t, int64(3), tv.FieldByIndex([]int{3, 1}).Int())
		assert.Equal(t, int64(17), tv.FieldByIndex([]int{3, 2}).Int())

		// 通过 .FieldByNameFunc(func) 方法获取属性名为 Id 的字段, 并确认字段值
		fd := tv.FieldByNameFunc(func(name string) bool { return name == "Id" })
		assert.Equal(t, int64(1), fd.Int())

		// 通过 .FieldByNameFunc(func) 方法获取属性名为 Name 的字段, 并确认字段值
		fd = tv.FieldByNameFunc(func(name string) bool { return name == "Name" })
		assert.Equal(t, "Alvin", fd.String())

		// 通过 .FieldByNameFunc(func) 方法获取属性名为 Gender 的字段, 并确认字段值
		fd = tv.FieldByNameFunc(func(name string) bool { return name == "Gender" })
		assert.Equal(t, 'M', fd.Interface().(rune))

		// 通过 .FieldByNameFunc(func) 方法获取属性名为 Birthday 的字段, 并确认字段值
		fd = tv.FieldByNameFunc(func(name string) bool { return name == "Birthday" })

		tv.Fields()
	})

	// 测试通过反射获变量值时使用了错误类型, 会导致 Panic 错误
	t.Run("For Wrong Type", func(t *testing.T) {
		// 为 any 类型变量设置整数值
		obj = 100

		// 获取 any 类型变量的值反射实例, 确认通过反射值获取浮点数时会导致 Panic 错误 (变量原本类型为 int)
		tv := reflect.ValueOf(obj)
		assert.PanicsWithError(t, "reflect: call of reflect.Value.Float on int Value", func() { tv.Float() })

		// 为 any 类型变量设置浮点数值
		obj = 100.001

		// 获取 any 类型变量的值反射实例, 确认通过反射值获取整型时会导致 Panic 错误 (变量本类型为 float64)
		tv = reflect.ValueOf(obj)
		assert.PanicsWithError(t, "reflect: call of reflect.Value.Int on float64 Value", func() { tv.Int() })

		// 为 any 类型变量设置字符串值
		obj = "hello"

		// 获取 any 类型变量的值反射实例, 确认通过反射值获取字符串时会导致 Panic 错误 (变量本类型不是指针)
		tv = reflect.ValueOf(obj)
		assert.PanicsWithValue(t, "reflect.Value.Addr of unaddressable value", func() { tv.Addr() })
	})
}
