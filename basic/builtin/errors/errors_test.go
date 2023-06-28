package errors

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// # 创建错误对象
//
// Go 语言遵循 "Error is a value" 的理念, 通过是否返回 `error` 对象来确定一个函数调用是否成功
//
// 所有错误都是 error 接口的对象, 保存着错误原因
//
// 错误可以传递, 即通过 `fmt.Errorf(...)` 函数, 通过 `%w` 占位符并传递 `error` 对象, 可以得到一个新 `error` 对象, 是前一个 `error` 对象的"包装"
//
// 包装对象可以通过 `errors.Unwrap` 函数可以获取"被包装"的 `error` 对象, 相当于获取错误链条的上一个环节
//
// 通过 `errors.Is` 函数可以判断一个错误链上是否有指定的错误对象, 即沿着错误链条逐级向上查找, 直到找到目标错误对象
func TestError(t *testing.T) {
	err1 := errors.New("invalid name")            // 创建一个 error 类型的错误对象
	assert.Error(t, err1)                         // 断言是否是 error 类型变量
	assert.Equal(t, "invalid name", err1.Error()) // 获取错误信息

	err2 := fmt.Errorf("error: %w", err1) // 通过 fmt.Errorf 创建一个带有错误信息的错误对象, %w 表示包装 (Wrap) 另一个错误对象
	assert.Error(t, err2)
	assert.Equal(t, "error: invalid name", err2.Error())

	// 通过 errors.Unwrap 函数可以获取到一个错误对象中包装的错误对象
	assert.Same(t, err1, errors.Unwrap(err2)) // err2 中包装了 err1
	// 通过 errors.Is 函数可以判断错误链条中是否最终包装了指定的错误
	assert.True(t, errors.Is(err2, err1)) // 判断 err2 中是否包装了 err1
}

// # 测试使用预定义错误
//
// Go 语言错误对象的一个最佳实践即, 将所有可以预期的错误对象进行预定义, 并在代码中返回这些错误, 而非返回临时构建的错误对象
func TestPredefinedError(t *testing.T) {
	// 调用函数, 返回错误对象
	err := causeLengthError()

	// 判断错误是否为预期类型
	assert.Same(t, ErrLength, err)    // 通过引用比较确认错误是否为预期类型
	assert.ErrorIs(t, err, ErrLength) // 通过 errors.Is 判断错误是否为预期类型
}

// # 自定义错误类型
//
// 通过自定义的 `CustomError` 类型对象表示的错误, 具备更多的错误信息, 且可以作为错误传递链条的一环, 包装另一个错误对象
//
// 对于一个 `error` 接口类型对象, 其原始类型是否为某个自定义错误类型, Go 语言提供了两种方法来进行判断:
//   - 通过 `err.(*CustomError)` 语法进行类型转换 (或 `switch err.(type)` 语法进行类型选择), 查看 `err` 变量类型是否为所期待的自定义错误类型;
//   - 通过 `errors.As` 函数, 第 1 个参数为 `error` 类型变量, 第 2 个参数为自定义错误类型的指针, 此时对于参数 1 链条上任意错误类型和参数 2 一致, 则返回 `true`;
func TestCustomErrorType(t *testing.T) {
	// 使用自定义错误类型
	err := NewCustomError("custom error", 100, ErrLength) // 定义 error 接口引用, 引用到 LengthError 类型对象上
	assert.Error(t, err)                                  // 判断自定义类型是否为 error 类型

	// 可以通过类型转换将错误对象转为其原始类型
	cErr, ok := err.(*CustomError)
	assert.True(t, ok)
	assert.Equal(t, 100, cErr.val) // 转为原始类型后, 即可获取类型中定义的其它属性
	assert.Equal(t, "custom error, error value=100", cErr.Error(), cErr.msg)

	assert.ErrorIs(t, err, ErrLength)             // 具备 Unwrap 方法的 error 对象可以使用 errors.Is 方法进行判断
	assert.Same(t, ErrLength, errors.Unwrap(err)) // 同时可以通过 errors.Unwrap 方法获取被包装的错误对象

	// 也可以可以通过 errors.As 函数来判断 err 错误的类型是否和 cErr 相同
	cErr = &CustomError{}                 // 这段代码为 cErr 变量换个类型就容易理解: var err error = &CustomError{}, 将 CustomError 对象地址转为 error 类型
	assert.True(t, errors.As(err, &cErr)) // errors.As(err, any(&err)), 第 2 个参数要求为 error 对象的地址, 所以 err 变量相当于是 CustomError 对象指针的指针

	// errors.As 的一个增强功能即可以沿着错误链条逐个查找和指定类型相同的错误
	wrapErr := fmt.Errorf("wrap error %w", err)
	assert.True(t, errors.As(wrapErr, &cErr))
}

// # 测试 Go 异常处理
//
// 除了返回 `error` 对象外, Go 还提供了运行时异常的抛出和处理机制
//
// 被调用的函数内部通过调用 `panic` 函数抛出一个变量 (可以是 `error` 变量)
//
// 调用方函数内部可以通过 `defer` 关键字修饰一个函数调用, 当 `panic` 调用发生后, 该 `defer` 调用随即被调用
//
// `defer` 调用内部可以通过 `recover` 函数获取通过 `panic` 抛出的信息
func TestDeferAndPanicError(t *testing.T) {
	defer func() { // 当异常发生后, 该 defer 调用会被执行
		r := recover()                    // 捕获异常
		if err := r.(error); err != nil { // 判断异常类型
			// 处理异常
			assert.ErrorIs(t, err, ErrLength)
		}
		// 异常处理完毕, TestDeferAndPanicError 函数随之结束
	}()

	panicError()                         // 调用函数并抛出异常
	assert.Fail(t, "Cannot be run here") // 这里不会被调用
}
