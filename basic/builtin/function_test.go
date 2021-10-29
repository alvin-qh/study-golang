package builtin

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 定义具备单一返回值的函数
func pow(n int) float64 { return math.Pow(float64(n), 2) }

// 定义同时返回两个值的函数
func sqrt(n int) (float64, error) {
	if n <= 0 {
		return 0, fmt.Errorf("invalid number") // 返回 0 和错误值
	}
	return math.Sqrt(float64(n)), nil // 返回结果和无错误
}

// 单一命名返回值
// 可以给函数的返回值命名，在函数体内部，命名返回值相当于一个变量，给其赋的值将会作为函数的返回值
func add1(a int, b int) (sum int) {
	sum = a + b // 给明名返回值赋值
	if sum == 0 {
		return sum // 可以在 return 语句中显示的书写返回的值
	}
	return // 也可以省略书写返回值，但 return 语句不能省略
}

// 对于同时返回多个值的情况，可以为每个返回值命名
func add2(a int, b int) (sum int, sub int) {
	sum = a + b
	sub = a - b
	if sum != sub {
		return sum, sub // 显式返回所有值
	}
	return // 省略显式返回值，但命名返回值已经赋值作为返回值
}

// 如果某个参数之前的参数类型一致，则可以将这部分参数用 ',' 分隔，并仅在最后一个参数上声明类型
//
// 可以混合命名返回值和匿名返回值，但匿名返回值必须在命名返回值之前
// 此时 return 语句必须显式返回所需的值
func add2_1(a, b int) (int, sub int) {
	sub = a - b
	return a + b, sub
}

// 不定参数
// 在参数名后使用 ... 表示该参数为不定参数，不定参数可以支持 0 或 多个 指定类型的参数
// 不定参数在函数体内部表现为一个切片
func add3(a ...int) (sum int) {
	sum = 0
	for _, n := range a {
		sum += n
	}
	return
}

// 测试函数的参数和返回值
func TestFunctionArgsAndReturns(t *testing.T) {
	// 单一返回值
	r := pow(4)
	assert.Equal(t, 16.0, r)

	// 多返回值
	r, err := sqrt(16)
	assert.Nil(t, err)
	assert.Equal(t, 4.0, r)

	r, err = sqrt(-16)
	assert.Error(t, err)
	assert.Equal(t, 0.0, r)

	r1 := add1(10, 20)
	assert.Equal(t, 30, r1)

	r1, r2 := add2(10, 20)
	assert.Equal(t, 30, r1)
	assert.Equal(t, -10, r2)

	r1, r2 = add2_1(10, 20)
	assert.Equal(t, 30, r1)
	assert.Equal(t, -10, r2)

	// 不定参数
	r1 = add3(1, 2, 3, 4, 5)
	assert.Equal(t, 15, r1)
}

// 定义具有函数类型参数的函数，即将参数类型作为参数传递
func callback(fn func(int, int) int, a, b int) int {
	r := fn(a, b)
	return r
}

// 将 函数 作为 函数类型值 返回
func getExecutor() (executor func(a, b int) int) {
	executor = add1
	return
}

// 可以定义 函数 类型，表示一类特定 参数 和 返回值 的函数
// 通过函数类型定义的变量，可以引用到类型符合的函数上（即存储该函数地址）
// 通过函数类型变量，可以调用函数，就和使用函数名调用函数类似
// 函数类型可以定义为指针类型，但这不表示“指向函数的指针”，实际上，函数变量本身就是“指向函数的指针”；指针类型的函数变量表示“指向另一个函数变量的指针”
// 函数类型
func TestTypeOfFunction(t *testing.T) {
	// 定义函数类型
	type FuncType = func(a, b int) (r int) // 定义函数类型，包括其参数和返回值
	var func1 FuncType = add1              // 定义函数类型变量并复制

	r := func1(12, 13) // 通过函数变量执行函数
	assert.Equal(t, 25, r)

	// 定义函数指针类型
	// 函数指针实际上是 指向函数变量 的指针，而不是指向函数的指针（函数变量原本就表示指向函数的指针）
	type FuncTypePtr = *func(a, b int) (r int) // 定义函数指针类型
	var func2 FuncTypePtr = &func1             // 指针指向一个函数变量的地址, func2 = &add2 这种写法是不被允许的，这一点和 C 语言不同

	r = (*func2)(12, 13) // 函数指针变量解引后即和普通函数变量使用相同
	assert.Equal(t, 25, r)

	// 函数变量的另一个作用是：结合闭包和匿名函数的使用
	var func3 FuncType = func(x, y int) int { return x + y }
	r = func3(1, 2)
	assert.Equal(t, 3, r)

	// 当然，对于无需传递的匿名函数，则可以自调用
	r = func(x, y int) int { return int(math.Pow(float64(x), float64(y))) }(4, 2)
	assert.Equal(t, 16, r)

	// 将函数作为参数传递给具有函数类型参数的函数
	r = callback(add1, 10, 20)
	assert.Equal(t, 30, r)

	// 将函数类型值作为返回值返回
	fn := getExecutor()
	r = fn(10, 20)
	assert.Equal(t, 30, r)
}
