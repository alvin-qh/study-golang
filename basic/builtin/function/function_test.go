package function

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试单返回值
// Pow 函数 1 个返回值
func TestPow(t *testing.T) {
	r := Pow(4)
	assert.Equal(t, 16.0, r)
}

// 测试多返回值
// Sqrt 函数具备 2 个返回值，第 2 个返回值表示是否有错误
func TestSqrt(t *testing.T) {
	r, err := Sqrt(16)
	assert.NoError(t, err)
	assert.Equal(t, 4.0, r)

	r, err = Sqrt(-16)
	assert.ErrorIs(t, err, ErrInvalidNumber)
	assert.Equal(t, 0.0, r)
}

// 测试命名返回值
// NumAdd 函数具备命名返回值
func TestNumAdd(t *testing.T) {
	r := NumAdd(10, 20)

	assert.Equal(t, 30, r)
}

// 测试多个命名返回值
// NumAddAndSub 函数具备 2 个命名返回值
func TestNumAddAndSub(t *testing.T) {
	r1, r2 := NumAddAndSub(20, 10)

	assert.Equal(t, 30, r1)
	assert.Equal(t, 10, r2)
}

// 测试函数参数定义的简要形式
func TestNumAddAndSubForm2(t *testing.T) {
	r1, r2 := NumAddAndSubForm2(20, 10)

	assert.Equal(t, 30, r1)
	assert.Equal(t, 10, r2)
}

// 测试不定参数
func TestAddForVarargs(t *testing.T) {
	r := AddForVarargs(1, 2, 3, 4)
	assert.Equal(t, 10, r)

	r = AddForVarargs(1, 2, 3, 4, 5)
	assert.Equal(t, 15, r)
}

// 测试函数作为变量
func TestFunctionAsVariable(t *testing.T) {
	type FuncType = func(a, b int) (r int) // 定义函数类型，包括其参数和返回值
	var f1 FuncType = NumAdd               // 定义函数类型变量并复制

	r := f1(10, 20)
	assert.Equal(t, 30, r)
}

// 测试函数作为参数传递
// Callback 函数的第一个参数类型是 func(a, b int) int，可以接受一个函数作为参数
func TestFunctionAsArgument(t *testing.T) {
	r := Callback(NumAdd, 10, 20)
	assert.Equal(t, 30, r)

	r = Callback(func(a, b int) int { return a * b }, 10, 20)
	assert.Equal(t, 200, r)
}

// 测试从函数返回一个函数类型返回值
// GetExecutor 函数返回一个 func(a, b int) int 类型的返回值
func TestFunctionAsReturnValue(t *testing.T) {
	f := GetExecutor()

	r := f(10, 20)
	assert.Equal(t, 30, r)
}

// 函数变量指针
// Go 语言的函数本身是引用类型，所以不存在函数指针这一概念
// 但如果一个函数复制给变量后，则可以获取该变量的地址
func TestFunctionPointer(t *testing.T) {
	// 定义函数变量
	f := NumAdd

	// 定义函数变量指针类型
	type FuncPtr = *func(a, b int) int
	var pf FuncPtr = &f // 函数变量地址赋值给函数变量

	r := (*pf)(10, 20)
	assert.Equal(t, 30, r)
}

// 函数闭包
// 闭包函数定义了一个作用域，但也可以使用外部作用域定义的变量
func TestClosureFunction(t *testing.T) {
	x := 10
	y := 20

	// 在闭包内使用外部变量值
	f1 := func(z int) int { return z + x + y }
	r := f1(30)
	assert.Equal(t, 60, r)

	// 在闭包内修改外部变量值
	f2 := func(a, b int) {
		x = a
		y = b
	}
	f2(100, 200)
	assert.Equal(t, 100, x)
	assert.Equal(t, 200, y)

	// 定义闭包并直接执行
	r = func() int { return x * y }()
	assert.Equal(t, 20000, r)
}
