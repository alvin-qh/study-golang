package error_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 自定义错误结构体
//
// 实现了 `error` 接口的结构体即可看作错误类型
type CustomError struct {
	err error
	tm  time.Time
}

// 增加 `Error` 方法, 实现 `error` 接口
//
// 注意, 安装 Go 语言规范, `Error` 方法需要依照错误实例 (而非错误实例指针) 进行定义, 即:
//
// 错误定义:
//
//	func (e *CustomError) Error() string {
//	    return fmt.Sprintf("custom error: %v, raised at %v", e.err, e.tm.Format(time.RFC3339))
//	}
//
// 正确定义:
//
//	func (e CustomError) Error() string {
//	    return fmt.Sprintf("custom error: %v, raised at %v", e.err, e.tm.Format(time.RFC3339))
//	}
func (e CustomError) Error() string {
	return fmt.Sprintf("custom error: %v, raised at %v", e.err, e.tm.Format(time.RFC3339))
}

// 返回自定义错误
func causeCustomError(tm *time.Time) error { return CustomError{ErrType, *tm} }

// 测试自定义错误类型
//
// 可通过 `errors.As(err, target)` 函数判断一个 `err` 参数是否为 `target` 指针指定的错误类型, 并且在类型匹配时,
// 将错误实例赋值到 `target` 指针指向的错误类型实例
func TestError_CustomError(t *testing.T) {
	now := time.Now()

	// 调用函数, 返回自定义错误实例
	err := causeCustomError(&now)

	var target CustomError
	// 判断自定义错误类型是否为 target 参数类型, 并且在类型匹配时, 将错误实例赋值给 target 实例
	assert.True(t, errors.As(err, &target))
	assert.EqualError(t, target, "custom error: invalid type, raised at "+now.Format(time.RFC3339))
}
