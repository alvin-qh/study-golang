package callstack

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试获取调用栈信息
func TestGetCallStack(t *testing.T) {
	// 测试在调用函数中获取调用堆栈

	// 获取调用堆栈
	s := CallStack()

	// 按行拆分为调用栈信息
	lines := strings.Split(s, string(lineBreak))

	// 第 1 行为协程信息, 形如: goroutine 6 [running]:
	assert.Regexp(t, `goroutine \d+ \[running\]:`, lines[0])
	// 第 2 行为当前调用函数, 形如: basic/runtime/callstack.CallStack()
	assert.Equal(t, "study/basic/runtime/callstack.CallStack()", lines[1])
	// 第 3 行为第 2 行的详细说明, 形如: .../basic/runtime/callstack/call_stack.go:17 +0x45
	assert.Regexp(t, `.+?/basic/runtime/callstack/callstack.go:\d+ \+0x[a-f0-9]+`, lines[2])

	assert.Regexp(t, `study/basic/runtime/callstack.TestGetCallStack\(0x[a-f0-9]+\)`, lines[3])
	assert.Regexp(t, `.+?/basic/runtime/callstack/callstack_test.go:\d+ \+0x[a-f0-9]+`, lines[4])

	// 测试在 defer 中获取调用堆栈

	panicFunc := func() {
		panic(errors.New("test panic"))
	}

	// 处理 panic 的 defer 调用
	defer func() {
		assert.Error(t, recover().(error))
		s := CallStack()

		lines = strings.Split(s, string(lineBreak))
		assert.Regexp(t, `goroutine \d+ \[running\]:`, lines[0])

		assert.Equal(t, "study/basic/runtime/callstack.CallStack()", lines[1])
		assert.Regexp(t, `.+?/basic/runtime/callstack/callstack.go:\d+ \+0x[a-f0-9]+`, lines[2])

		assert.Regexp(t, `study/basic/runtime/callstack.TestGetCallStack.func2()`, lines[3])
		assert.Regexp(t, `.+?/basic/runtime/callstack/callstack_test.go:\d+ \+0x[a-f0-9]+`, lines[4])

		assert.Regexp(t, `panic\(\{0x[a-f0-9?]+, 0x[a-f0-9?]+\}\)`, lines[5]) // panic 发生的位置
	}()

	panicFunc() // 调用引发 panic 的函数
}
