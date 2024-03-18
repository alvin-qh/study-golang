package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试 `interface{}` 类型和其它类型的转换
//
// 在 Go 语言中, `interface{}` 类型可以表示为任意类型, 即任何类型实例都可以转为 `interface{}` 类型,
// 转换过程通过强制转换运算即可, 例如:
//
//	v := 100  // 定义一个 int32 类型变量
//	i := interface{}(v)  // 将 int32 类型转为 interface{} 类型
//
// 而从 `interface{}` 转为其原始类型则需要通过 Go 反射机制, Go 提供了如下 3 中方式:
// 方式 1: 通过 `.` 运算符进行转换, 转换失败会抛出 Panic 异常
//
//	i := interface{}(100)
//	n := i.(int)
//	// todo: 使用转换为 int 类型的 n 变量
//
// 方式 2: 通过 `.` 运算符进行转换, 并通过返回值判断是否转换成功
//
//	i := interface{}(100)
//	if n, ok := i.(int); ok {
//	   // todo: 使用转换为 int 类型的 n 变量
//	}
//
// 方式 3: 通过 `switch` 语法进行类型匹配
//
//	i := interface{}(100)
//	switch i.(type) {
//	case int32:
//	  // todo: 使用转换为 int32 类型的 n 变量
//	case float32:
//	  // todo: 使用转换为 float32 类型的 n 变量
//	case float64:
//	  // todo: 使用转换为 float64 类型的 n 变量
//	}
func TestInterfacesConvert(t *testing.T) {
	v := int(100)
	assert.IsType(t, int(0), v)

	// 将变量转为 interface{} 类型
	i := interface{}(v)
	assert.IsType(t, int(0), i)

	// 将 interface{} 类型转为 int 类型, 转换失败会 Panic
	v = i.(int)
	assert.Equal(t, v, 100)

	// 将 interface{} 类型转为 int 类型, 返回 bool 类型变量表示是否转换成功
	if v, ok := i.(int); ok {
		assert.Equal(t, v, 100)
	}

	// 通过 switch 判断变量的类型
	switch i.(type) {
	case int:
		assert.True(t, true)
	case float32:
	case float64:
		assert.Fail(t, "Invalid type touched")
	}
}

func TestInterfaceAsArgument(t *testing.T) {
	v1 := NewValue(100)

	tp, n := ExplainValue(v1)
	assert.Equal(t, tp, I32)
	assert.Equal(t, n, 100)

	v2 := NewValue(123.123)

	tp, n = ExplainValue(&v2)
	assert.Equal(t, tp, F64)
	assert.Equal(t, n, 123.123)
}
