package errors

import (
	"errors"
	"fmt"
)

// 在 Go 语言中, 所有的异常都应该实现 `error` 接口, 即具备 `Error()` 方法
// `Error()` 方法的返回值是 `string` 类型, 表示错误信息

// 创建一组错误值, 用于传递和比较
// 一般情况下, 对于已知错误, 都应该在集中位置, 事先定义好错误值, 这样的好处: 1. 执行效率会高一些; 2. 可以对得到的错误对象进行比较, 以判断其错误原因
var (
	ErrLength = errors.New("invalid length")
)

// # 返回预定义错误对象
//
// 该方法返回预定义的错误对象, 以便于错误接收方对错误类型进行判断
//
// 返回 `ErrLength` 全局错误变量
func causeLengthError() error {
	return ErrLength
}

// # 自定义错误类型
//
// 有时候仅仅通过字符串无法清晰的表达错误信息, 此时可以自定义错误类型
//
// 自定义错误类型只需提供 `Error()` 方法, 即实现了 `error` 接口, 就可以被认为是错误类型
//
// 自定义错误类型如果提供了 `Unwrap()` 方法, 即可作为错误包装对象, 包装其它错误对象, 从而完成错误的传递, 可以使用 `errors.Unwrap` 函数和 `errors.Is` 函数进行操作
type CustomError struct {
	msg string
	val any
	err error
}

// # 新建自定义错误对象
//
// 参数:
//   - `err`: 错误信息字符串
//   - `value`: 错误对象中包含的其它错误信息
//   - `caused`: 引发当前错误的上一级错误对象, 对其进行包装
//
// 返回: `CustomError` 类型错误对象
func NewCustomError(err string, value any, caused error) error {
	return &CustomError{msg: err, val: value, err: caused}
}

// # 获取错误信息
//
// 返回描述错误信息的字符串
//
// 返回: 错误信息字符串
func (c *CustomError) Error() string {
	return fmt.Sprintf("%v, error value=%v", c.msg, c.val)
}

// # 获取被包装的错误对象
//
// 该自定义错误类型可以传递错误对象, 如果在创建时指定了被包装的错误对象, 则这里返回该被包装对象
//
// 返回: 被包装的错误对象
func (c *CustomError) Unwrap() error {
	return c.err
}

// # 抛出异常的函数
//
// 调用 `panic` 函数会终止当前函数调用, 代码会跳转到该函数调用方的 `defer` 调用上, 该 `defer` 调用内部可通过 `recover` 函数获取该异常,
// 并结束调用方函数
//
// 如果调用方函数未处理这个异常, 则异常会继续上更上一级的调用者抛出
func panicError() {
	panic(ErrLength)
}
