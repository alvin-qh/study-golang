package errors

import (
	"errors"
	"fmt"
)

// # 定义异常包装接口
type Wrapped interface {
	// # 包装上一层异常对象
	//
	// 参数:
	//   - `caused`: 引发当前异常的前一个异常对象
	Wrap(caused error)
}

// # 定义自定义错误类型
//
// 表示长度不正确的错误信息
type LengthError struct {
	length   int
	expected int
	caused   error // 引发当前错误的另一个错误对象
}

// # 获取错误信息
//
// 为 `LengthError` 类型实现 `error` 接口
//
// 返回错误信息字符串
func (e *LengthError) Error() string {
	return fmt.Sprintf("invalid length: %d, expected %d", e.length, e.expected)
}

// 为 `LengthError` 实现接口, 使能够获取上一层错误信息 (即 `caused` 字段值), 以能够调用 `errors.Unwrap` 函数
func (e *LengthError) Unwrap() error {
	return e.caused
}

// # 传入引发当前错误的另一个错误
//
// 为 `LengthError` 类型实现 `Wrapped` 接口
func (e *LengthError) Wrap(caused error) {
	e.caused = caused
}

// 创建一组错误值, 用于传递和比较
// 一般情况下, 对于已知错误, 都应该在集中位置, 事先定义好错误值, 这样的好处: 1. 执行效率会高一些; 2. 可以对得到的错误对象进行比较, 以判断其错误原因
var (
	ErrName = errors.New("invalid name")
)

// # 自定义错误
type EmptyError struct {
	name string
}

// 为 `EmptyError` 实现 `error` 接口
func (e *EmptyError) Error() string {
	return fmt.Sprintf("%s is required", e.name)
}
