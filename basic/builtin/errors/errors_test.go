package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 创建错误变量
// go 语言遵循 “error is a value” 的理念，通过是否返回 error 来确定一个函数调用是否成功
// 所有错误都是 error 接口的对象，保存着错误信息
func TestMakeError(t *testing.T) {
	err := errors.New("invalid name")            // 创建一个 error 类型的错误对象
	assert.Error(t, err)                         // 断言是否是 error 类型变量
	assert.Equal(t, "invalid name", err.Error()) // 获取错误信息

	err = fmt.Errorf("error: %v", err) // 通过 fmt.Errorf 创建一个带有错误信息的错误对象，%v 表示该错误信息中引用另一个错误的信息
	assert.Error(t, err)
	assert.Equal(t, "error: invalid name", err.Error())
}

// 自定义错误类型
// error 接口只定义了 Error() string 方法，所以所有实现了该方法的结构体，均可视为 error 类型，作为错误值使用
func TestCustomErrorType(t *testing.T) {
	// 使用自定义错误类型
	var err error = &LengthError{length: 12, expected: 20}          // 定义 error 接口引用，引用到 LengthError 类型对象上
	assert.Error(t, err)                                            // 判断自定义类型是否为 error 类型
	assert.Equal(t, "invalid length: 12, expected 20", err.Error()) // 获取错误信息

	// 所有 error 类型的变量都可以转化为其原始错误类型，转换方式同 go 语言接口引用类型转换语法
	if e := err.(*LengthError); e != nil { // 将 error 变量的引用转为 *LengthError 类型，e 即为 LengthError 引用
		assert.Error(t, e)
		assert.Equal(t, "invalid length: 12, expected 20", e.Error())
	} else {
		assert.Fail(t, "Cannot run here")
	}

	if _, ok := err.(*LengthError); ok { // 类型转换也可以通过 ok 这个哨兵返回值来判断
		assert.Error(t, err)
		assert.Equal(t, "invalid length: 12, expected 20", err.Error())
	} else {
		assert.Fail(t, "Cannot run here")
	}
}

// 测试错误值的比较和类型判断
// 如果为错误类型增加一个指向另一个错误的引用（参见 LengthError::caused 字段），则组成了错误链表（error chains），可以以此回溯到最初引发错误的那个错误
// 通过 errors.Is(err, target, error) 函数，可以判断一个错误是否为指定错误，该函数会不断调用错误对象的 Unwrap 函数，相当于沿着错误链表回溯，直到错误匹配，否则返回 false
// 通过 errors.As(err error, target *error) 函数，可以判断两个参数是否是相同类型的错误值，并把前一个参数传递给后一个参数。注意：target 相当于是 error 引用的指针
// error.As 方法也会通过调用错误的 Unwrap 函数，不断回溯错误链，直到找到类型匹配的项，否则返回 false
func TestErrorIsOrAs(t *testing.T) {
	// 测试错误值比较
	var err error = ErrName                 // 产生一个错误对象
	assert.True(t, errors.Is(err, ErrName)) // 判断错误值是否等于 ErrName
	assert.ErrorIs(t, err, ErrName)         // errors.Is 的断言

	// 测试沿错误链表回溯并比较
	// 这是产生错误链的简单方式，即无需额外定义错误类型，直接进行 Wrap 操作即可
	err = fmt.Errorf("%w, name is: %s", ErrName, "alvin") // 产生一个错误，并包装 ErrName 错误，相当于 err -> ErrName 的链
	assert.True(t, errors.Is(err, ErrName))               // 判断错误值是否等于或包含 ErrName 错误，即通过 Unwrap 不断回溯比较
	assert.ErrorIs(t, err, ErrName)                       // errors.Is 的断言
	assert.Equal(t, ErrName, errors.Unwrap(err))          // 显示调用 errors.Unwrap 函数，获取错误

	// 判断错误类型
	var targetErr1 *LengthError

	err = fmt.Errorf("invalid length")           // 定义错误
	assert.False(t, errors.As(err, &targetErr1)) // 确认 err 变量不是 LengthError 类型（或错误链上也没有 LengthError 类型）

	err = &LengthError{length: 10, expected: 20} // 定义 LengthError 类型的错误
	assert.True(t, errors.As(err, &targetErr1))  // 确认 err 变量是 LengthError 类型，并将 err 的引用传递给 targetErr1 变量
	assert.Same(t, err, targetErr1)              // 现在 err 和 targetErr1 具备相同的引用
	assert.Equal(t, "invalid length: 10, expected 20", targetErr1.Error())

	err = fmt.Errorf("%w caused", &LengthError{length: 10, expected: 20})  // 包装 LengthError，产生 err -> LengthError 的错误链
	assert.True(t, errors.As(err, &targetErr1))                            // 确认 err 的链上有 LengthError 类型的错误值（通过不断调用 Unwrap 函数），并将符合类型的错误值传递给 targetErr1 变量
	assert.Equal(t, "invalid length: 10, expected 20", targetErr1.Error()) // targetErr1 引用到了被包装的错误值上

	var targetErr2 *EmptyError
	err = &LengthError{length: 10, expected: 20}            // 定义 err 变量引用到 LengthError 类型
	err.(Wrapped).Wrap(&EmptyError{name: "name"})           // 手动调用 Wrap 函数，为 LengthError 类型的错误传入另一个错误，相当于 LengthError -> EmptyError 的链
	assert.True(t, errors.As(err, &targetErr2))             // 确认 err 的链上具有 EmptyError 类型的错误值，并传递该错误值引用到 targetErr2 变量
	assert.Equal(t, "name is required", targetErr2.Error()) // 确认 targetErr2 是引用到 EmptyError 类型的变量

	err = errors.Unwrap(err) // 手动调用错误的 Unwrap 函数，获取前一个错误
	assert.IsType(t, &EmptyError{}, err)
}

// 抛出异常的函数
func PanicError() {
	panic(ErrName) // 抛出异常，panic 调用会终止当前函数调用，即 panic 之后的代码不会被调用
	// 调用终止后，代码执行会跳转到调用方的 defer 调用上，如果该 defer 调用具备 recover 捕获
	// 则会捕获该异常，并结束调用方函数，否则异常会继续上更上一级的调用者抛出
}

func TestDeferAndPanicError(t *testing.T) {
	defer func() { // 当异常发生后，该 defer 调用会被执行
		r := recover()                    // 捕获异常
		if err := r.(error); err != nil { // 判断异常类型
			// 处理异常
			assert.ErrorIs(t, err, ErrName)
		}

		// 异常处理完毕，TestDeferAndPanicError 函数随之结束
	}()

	PanicError() // 调用函数并抛出异常
	assert.Fail(t, "Cannot be run here")
}
