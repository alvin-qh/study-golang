package builtin

import (
	"basic/builtin/types"
	"reflect"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

func TestStructType(t *testing.T) {
	// 定义结构体变量并进行初始化
	var u1 types.User = types.User{Id: 1, Name: "Alvin", Gender: 'M'}
	assert.Equal(t, "User", reflect.TypeOf(u1).Name()) // 查看变量类型
	assert.Equal(t, 1, u1.Id)
	assert.Equal(t, "Alvin", u1.Name)
	assert.Equal(t, 'M', u1.Gender)

	// 定义结构体并跳过初始化
	u1 = types.User{}
	u1.Id = 2 // 为结构体字段赋值
	u1.Name = "Emma"
	u1.Gender = 'F'
	assert.Equal(t, 2, u1.Id)
	assert.Equal(t, "Emma", u1.Name)
	assert.Equal(t, 'F', u1.Gender)

	// 通过 new 操作符产生 User 类型的指针变量
	pu1 := new(types.User)
	(*pu1) = u1 // 为指针指向的结构体赋值
	assert.Equal(t, 2, pu1.Id)
	assert.Equal(t, "Emma", pu1.Name)
	assert.Equal(t, 'F', pu1.Gender)

	// 赋值语句可以 copy 结构体
	u2 := u1
	u2.Id = 3
	assert.Equal(t, 2, u1.Id) // userCopy 和 user 两个变量表示两个不同的 User 对象
	assert.Equal(t, 3, u2.Id)

	pu1 = &u1 // 指针指向 user 变量结构体
	pu1.Id = 3
	assert.Equal(t, 3, u1.Id) // pUser 指向 user 结构体
	assert.Equal(t, 3, pu1.Id)
}

// go 语言的类型转化基于非常简单个规则：值类型转换
// go 语言具备非常简单的类型系统：基本值类型，结构体类型，指针类型和 interface{} 类型
// 基本值类型包括：数字类型和布尔类型，数字类型包括：int8~64, float32~64, complex64~128, rune，数值类型之间可以直接转换, bool 类型只能是 true, false
// 结构体之间无法进行转换，只能依赖接口对类型进行处理
func TestTypeConvert(t *testing.T) {
	// 对于值类型，可以通过类型运算符进行类型转换
	// 类型转换是赋值的一种副作用，即在内存间进行数值复制的时候，对数值做了一次类型变更操作，例如将 8byte 数值复制到 4byte 空间中
	// 所以类型转换的过程中可能会 丢失精度
	var v1 float64 = float64(123.456)
	var v2 int64 = int64(v1)
	assert.NotEqual(t, v2, v1) // 转换前后的两个变量不相同
	assert.NotSame(t, v2, v1)

	var v3 int32 = int32(v2)
	assert.NotEqual(t, v3, v2) // 转换前后的两个变量不相同
	assert.NotSame(t, v3, v2)
	assert.EqualValues(t, v2, v3) // 转换前后的两个变量值相同但类型不同

	// 对于指针类型，不同类型的指针类型不能直接互相转化，需要借助 unsafe 包转换为 uintptr 类型后间接转换
	// go 语言指针实际上都是 64 位整型地址，指针类型实际上表达的是指针指向的内存空间存储数据的类型
	// 如果对指针类型进行任意转换，但并未改变指针所指向地址的数据类型，可能会导致指针解引的时候出现错误
	var pv1 *float64 = &v1 // 获取变量地址，赋予指针类型变量
	assert.Equal(t, 123.456, *pv1)

	var pv2 *int64 = (*int64)(unsafe.Pointer(&v1))    // 要将 *float64 类型变量强制转换为 *int64，需要通过 unsafe.Pointer 进行
	assert.Equal(t, int64(4638387860618067575), *pv2) // 转换后解引指针，得到错误结果，说明指针指向的变量类型错误

	// 任何类型都可以转为 interface{} 类型，即 go 语言特有的 “空接口” 类型
	// interface{} 类型可以看作是“未知”类型，可以通过“类型断言”转化为其原始类型
	var obj interface{} = v1     // 令 obj 为 float64 类型
	if v1, ok := obj.(int); ok { // 类型断言，ok 为 true 时，v1 为转换后的结果，否则 v1 不可用
		assert.Equal(t, 123.456, v1)
	}

	// 也可以通过类型 switch 对类型进行判断
	switch val := obj.(type) { // 对 obj 进行 type 判断，并赋值给 val 变量
	case float64:
		assert.Equal(t, 123.456, val) // 当 type 为 float64 类型时
	case int64, int32:
		assert.Fail(t, "Cannot run here")
	default:
		assert.Fail(t, "Cannot run here")
	}
}

// 指针运算
// 默认情况下，go 语言不允许指针执行运算操作
// 要对指针进行运算，需要使用 unsafe.Pointer 类型，并将得到的结果转为 uintptr 类型
// 要将指针运算结果转回原始指针类型，需要对 uintptr 类型再次作转换为 unsafe.Pointer 类型，在转换回所需指针类型
// 和 C 语言不同的是，uintptr 类型的运算不能按照类型本身的大小进行计算，只能按 byte 大小进行，所以需要测量类型的大小
func TestPointerCalculation(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5}

	// 获取 切片 第 1 个元素的地址，即整个切片的地址
	pn := &slice[0]
	assert.Equal(t, 1, *pn)

	// 将指针移动一个 int 大小，指针指向切片第 2 个元素的地址
	pn = (*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&slice[0])) + unsafe.Sizeof(slice[0])))
	assert.Equal(t, 2, *pn)
}
