package errors

import (
	"errors"
	"fmt"
)

// 定义具有 Wrapped(error) 函数的接口类型
// Wrapped 函数的作用是存入另一个 error 对象，表示当前的错误是由另一个错误引发
type Wrapped interface {
	Wrap(caused error)
}

// 自定义错误类型，并实现 Wrapable 接口 Unwrapped 接口

// 定义错误类型结构体
type LengthError struct {
	length   int
	expected int
	caused   error // 引发当前错误的另一个错误对象
}

// 为 LengthError 实现 error 接口，获取错误信息
func (e *LengthError) Error() string {
	return fmt.Sprintf("invalid length: %d, expected %d", e.length, e.expected)
}

// 为 LengthError 实现接口，使能够获取上一层错误信息（即 caused 字段值），以能够调用 errors.Unwrap(err) 函数
func (e *LengthError) Unwrap() error {
	return e.caused
}

// 为 LengthError 实现 Wrapped 接口，传入引发当前错误的另一个错误
func (e *LengthError) Wrap(caused error) {
	e.caused = caused
}

// 创建一组错误值，用于传递和比较
// 一般情况下，对于已知错误，都应该在集中位置，事先定义好错误值，这样的好处：1. 执行效率会高一些；2. 可以对得到的错误对象进行比较，以判断其错误原因
var (
	ErrorName = errors.New("invalid name")
)

// 自定义错误
type EmptyError struct {
	name string
}

// 为 EmptyError 实现 error 接口
func (e *EmptyError) Error() string {
	return fmt.Sprintf("%s is required", e.name)
}
