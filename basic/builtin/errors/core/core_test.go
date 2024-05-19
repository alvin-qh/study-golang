package core

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 定义所需的错误值
var (
	ErrType  = errors.New("invalid type")
	ErrValue = errors.New("invalid value")
)

// 测试不同 `error` 实例是否相等
//
// 错误实例中如果错误信息相同, 则两个错误实例相等
func TestError_Equal(t *testing.T) {
	assert.Equal(t, ErrType, errors.New("invalid type"))
	assert.Equal(t, ErrType, fmt.Errorf("invalid type"))
}

// 包装一个现有的错误实例
func wrapError(otherErr error) error {
	return fmt.Errorf("wrapper error: %w", otherErr)
}

// 测试错误实例的包装
//
// 可以依据一个现有错误实例生成一个新错误实例, 后者包装了前者
//
//	fmt.Errorf("wrapper error: %w", err)
//
// 对于一个已经被包装的错误实例, 可以通过 `errors.Is` 断言判断是否包含指定的错误实例
//
//	errors.Is(err, ErrType)
//
// 可以通过被包装的错误实例获取到其原始的错误实例
//
//	srcErr := errors.Unwrap(err)
func TestError_Wrap(t *testing.T) {
	err := wrapError(ErrType)
	assert.EqualError(t, err, "wrapper error: invalid type")

	// 判断指定错误是否已经包装了已知错误实例
	assert.True(t, errors.Is(err, ErrType))

	// ErrorIs 断言用于确认一个错误实例是否包装了指定的另一个错误实例
	assert.ErrorIs(t, err, ErrType)

	// 获取被包装
	uwErr := errors.Unwrap(err)
	assert.EqualError(t, uwErr, "invalid type")
	assert.ErrorIs(t, err, ErrType)
}

// 测试组合多个错误实例
//
// 通过 `errors.Join` 函数可以组合多个错误实例, 返回组合后的错误实例,
// 可以通过 `errors.Is` 函数判断组合错误实例中是否包含指定的错误
func TestError_Join(t *testing.T) {
	// 组合两个错误实例
	joinedErr := errors.Join(ErrType, ErrValue)
	// 组合后的错误信息同时包含被组合错误实例的全部信息, 通过 `\n` 分隔
	assert.EqualError(t, joinedErr, "invalid type\ninvalid value")

	// 可以通过 `errors.Is` 函数判断组合错误实例中是否包含指定的错误
	assert.True(t, errors.Is(joinedErr, ErrType))
	assert.True(t, errors.Is(joinedErr, ErrValue))

	// 可以通过包含 `Unwrap() []error` 接口的错误实例获取到被组合的错误实例
	errs := joinedErr.(interface{ Unwrap() []error }).Unwrap()
	assert.Equal(t, []error{ErrType, ErrValue}, errs)
}

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
func causeCustomError(tm *time.Time) error {
	return CustomError{ErrType, *tm}
}

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
