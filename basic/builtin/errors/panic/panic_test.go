package panic

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 抛出 Panic 的函数
func raisePanic(err error) {
	if err != nil {
		panic(err)
	}
}

// 测试从 Panic 中恢复的错误代码执行
//
// 在一个函数中, 如果预计会引发 Panic, 且不希望中断代码执行, 则可以使用 defer 函数恢复 Panic
//
// 因为 defer 会在函数结束前执行, 所以可以在 defer 函数中通过 `recover` 函数从 Panic 中恢复代码执行,
// `recover` 函数的返回值可以转换为 Panic 抛出的类型, 从而捕获到 Panic 引发的的错误
func TestError_RecoverPanic(t *testing.T) {
	var panicErr error = nil

	// 调用函数, 在该函数内引发了 Panic, 且通过 defer 函数中调用 `recover` 函数恢复函数调用,
	// 所以该函数可以正常返回
	r := func() string {
		defer func() {
			r := recover()
			if e, ok := r.(error); ok {
				panicErr = e
			}
		}()

		raisePanic(fmt.Errorf("test error"))
		return "OK"
	}()

	// 代码执行到这里, 表示 Panic 并未中断代码, 代码已经恢复

	// 未执行到函数返回, 所以未接受到函数返回值
	assert.Equal(t, "", r)

	// 查看 Panic 恢复后捕获的异常信息
	assert.NotNil(t, panicErr)
	assert.Equal(t, "test error", panicErr.Error())
}
