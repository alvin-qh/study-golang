package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 处理带有错误返回值的函数调用
//
// Args:
//   - `t` 测试实例指针
//   - `f` 带有错误返回值的函数指针
//
// Returns:
//   - `T` 函数返回值
func WithNoError[T any](t *testing.T, f func() (T, error)) (v T) {
	v, err := f()
	assert.NoError(t, err)
	return
}

// 处理带有错误返回值的函数调用
//
// Args:
//   - `t` 测试实例指针
//   - `f` 带有错误返回值的函数指针
//
// Returns:
//   - `err` 错误实例
func WithError[T any](t *testing.T, f func() (T, error)) (err error) {
	_, err = f()
	return
}
