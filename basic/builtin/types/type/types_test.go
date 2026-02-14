package type_test

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 定义结构体, 所有成员字段为 public
type User struct {
	Id     int
	Name   string
	Gender rune
}

// 测试结构体类型
func TestType_StructType(t *testing.T) {
	// 测试初始化结构体实例
	t.Run("initialize struct instance", func(t *testing.T) {
		// 定义结构体变量并进行初始化
		var u User = User{
			Id:     1,
			Name:   "Alvin",
			Gender: 'M',
		}

		// 查看变量类型
		assert.Equal(t, reflect.Struct, reflect.TypeOf(u).Kind())
		assert.Equal(t, "User", reflect.TypeOf(u).Name())

		assert.Equal(t, 1, u.Id)
		assert.Equal(t, "Alvin", u.Name)
		assert.Equal(t, 'M', u.Gender)
	})

	// 测试获取或设置结构体实例的字段值
	t.Run("get and set properties of struct instance", func(t *testing.T) {
		// 定义结构体并按默认值初始化字段值
		u := User{}

		// 查看变量类型
		assert.Equal(t, reflect.Struct, reflect.TypeOf(u).Kind())
		assert.Equal(t, "User", reflect.TypeOf(u).Name())

		// 为结构体字段赋值
		u.Id = 2
		u.Name = "Emma"
		u.Gender = 'F'

		assert.Equal(t, 2, u.Id)
		assert.Equal(t, "Emma", u.Name)
		assert.Equal(t, 'F', u.Gender)
	})

	// 测试结构体指针类型
	t.Run("pointer of struct instance", func(t *testing.T) {
		// 通过 new 操作符产生 User 类型的指针变量
		pu := new(User)

		// 查看变量类型
		assert.Equal(t, reflect.Ptr, reflect.TypeOf(pu).Kind())
		assert.Equal(t, reflect.Ptr, reflect.TypeOf(pu).Kind())
		assert.Equal(t, "User", reflect.TypeOf(pu).Elem().Name())

		// 记录原始指针
		pSrc := pu

		*pu = User{
			Id:     1,
			Name:   "Alvin",
			Gender: 'M',
		}

		assert.Equal(t, pSrc, pu)
		assert.Equal(t, 1, pu.Id)
		assert.Equal(t, "Alvin", pu.Name)
		assert.Equal(t, 'M', pu.Gender)

		// 指针指向 User 变量结构体
		pu = &User{
			Id:     2,
			Name:   "Emma",
			Gender: 'F',
		}

		assert.NotEqual(t, pSrc, pu)
		assert.Equal(t, 2, pu.Id)
		assert.Equal(t, "Emma", pu.Name)
		assert.Equal(t, 'F', pu.Gender)
	})
}

// 测试类型转换
//
// Go 语言的类型转化基于非常简单的规则: 值类型转换
//
// Go 语言具备非常简单的类型系统: 基本值类型, 结构体类型, 指针类型和 `interface{}` 类型
//
// 基本值类型包括数字类型和布尔类型,
//   - 数字类型包括: int8~64, float32~64, complex64~128, rune, 数值类型之间可以直接转换;
//   - `bool` 类型只能是 `true`, `false`
//
// 结构体之间无法进行转换, 只能依赖接口对类型进行处理
func TestType_TypeConvert(t *testing.T) {
	// 测试值类型的强制类型转换
	//
	// 对于值类型, 可以通过类型运算符进行类型转换
	//
	// 类型转换是赋值的一种副作用, 即在内存间进行数值复制的时候, 对数值做了一次类型变更操作, 例如将 8byte 数值复制到 4byte 空间中,
	// 所以类型转换的过程中可能会丢失精度
	t.Run("forced type conversion", func(t *testing.T) {
		var v1 float64 = float64(123.456)
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
	})


	t.Run("type convert by switch statement", func(t *testing.T) {

	})

	// 测试指定类型切片和 `interface{}` 类型切片的转换
	//
	// 在 Go 语言中, 一般不推荐使用 `interface{}` 类型的切片, 即 `[]interface{}` 类型, 如果要用泛化类型表示数组,
	// 则使用 `interface{}` 直接表示即可
	//
	// 注意, 要在 `[]interface{}` 类型切片和其它类型切片间转换, 则需要通过一个 `O(n)` 复杂度的循环才能完成
	t.Run("slice type conversion", func(t *testing.T) {
		// 定义 interface{} 类型变量
		var v interface{} = []int{1, 2, 3, 4, 5}
		assert.Equal(t, reflect.Slice, reflect.TypeOf(v).Kind())

		// 将 interface{} 类型转为指定类型切片类型
		s, ok := AnyToSlice[int](v)
		assert.True(t, ok)
		assert.Equal(t, []int{1, 2, 3, 4, 5}, s)

		// 将指定类型的切片转为 `[]interface{}` 类型
		vs := TypedSliceToAnySlice([]int{1, 2, 3, 4, 5})
		assert.Len(t, vs, 5)

		// 将 `[]interface{}` 类型切片转为指定类型
		s, ok = AnySliceToTypedSlice[int](vs)
		assert.True(t, ok)
		assert.Equal(t, []int{1, 2, 3, 4, 5}, s)
	})
}
