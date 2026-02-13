package function_test

import (
	"study/basic/builtin/functions/function"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试单返回值
//
// 其中 `Pow` 函数 1 个返回值
func TestFunction_SingleReturn(t *testing.T) {
	r := function.Pow(4)
	assert.Equal(t, 16.0, r)
}

// 测试多返回值
//
// 其中 `Sqrt` 函数具备 2 个返回值, 第 2 个返回值表示是否有错误
func TestFunction_MultiReturn(t *testing.T) {
	r, ok := function.Sqrt(16)
	assert.True(t, ok)
	assert.Equal(t, 4.0, r)

	r, ok = function.Sqrt(-16)
	assert.False(t, ok)
	assert.Equal(t, 0.0, r)
}

// 测试命名返回值
//
// 其中 `NumAdd` 函数具备命名返回值
func TestFunction_SingleNamedReturn(t *testing.T) {
	r := function.NumAdd(10, 20)

	assert.Equal(t, 30, r)
}

// 测试多个命名返回值
//
// 其中 `NumAddAndSub` 函数具备 2 个命名返回值
func TestFunction_MultiNamedReturn(t *testing.T) {
	r1, r2 := function.NumAddAndSub(20, 10)

	assert.Equal(t, 30, r1)
	assert.Equal(t, 10, r2)
}

// 测试函数参数定义的简要形式
//
// 其中 `NumAddAndSubForm2` 函数具备 2 个同类型参数, 以简单形式书写
func TestFunction_SimpleArgsForm(t *testing.T) {
	r1, r2 := function.NumAddAndSubForm2(20, 10)

	assert.Equal(t, 30, r1)
	assert.Equal(t, 10, r2)
}

// 测试不定参数
//
// 其中 `AddForVarargs` 函数具备不定数量参数
func TestFunction_IndefiniteArgs(t *testing.T) {
	r := function.AddForVarargs(1, 2, 3, 4)
	assert.Equal(t, 10, r)

	r = function.AddForVarargs(1, 2, 3, 4, 5)
	assert.Equal(t, 15, r)
}

// 测试函数作为变量
func TestFunction_AsVariable(t *testing.T) {
	type FuncType = func(a, b int) (r int) // 定义函数类型, 包括其参数和返回值
	var f1 FuncType = function.NumAdd     // 定义函数类型变量并复制

	r := f1(10, 20)
	assert.Equal(t, 30, r)
}

// 测试函数作为参数传递
//
// 其中 `Callback` 函数的第一个参数类型是 `func(a, b int) int` 函数类型, 可以接受一个函数作为参数
func TestFunction_AsArgument(t *testing.T) {
	r := function.Callback(function.NumAdd, 10, 20)
	assert.Equal(t, 30, r)

	r = function.Callback(func(a, b int) int { return a * b }, 10, 20)
	assert.Equal(t, 200, r)
}

// 测试从函数返回一个函数类型返回值
//
// 其中 `GetExecutor` 函数返回一个 `func(a, b int) int` 函数类型的返回值
func TestFunction_AsReturnValue(t *testing.T) {
	f := function.GetExecutor()

	r := f(10, 20)
	assert.Equal(t, 30, r)
}
