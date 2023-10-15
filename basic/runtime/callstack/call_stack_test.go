package callstack

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCallStack(t *testing.T) {
	// 测试在调用函数中获取调用堆栈

	s := CallStack() // 获取调用堆栈

	lines := strings.Split(s, string(lineBreak))             // 分解为行集合
	assert.Regexp(t, `goroutine \d+ \[running\]:`, lines[0]) // 第 1 行为协程信息, 形如: goroutine 6 [running]:

	assert.Equal(t, `study-golang/basic/runtime/callstack.CallStack()`, lines[1])                             // 第 2 行为当前调用函数, 形如: basic/runtime/callstack.CallStack()
	assert.Regexp(t, `\s+/[a-zA-Z\-/]+/basic/runtime/callstack/call_stack.go:\d+ \+0x[a-zA-Z0-9]+`, lines[2]) // 第 3 行为第 2 行的详细说明, 形如: .../basic/runtime/callstack/call_stack.go:17 +0x45

	assert.Regexp(t, `study-golang/basic/runtime/callstack.TestGetCallStack\(0x[a-zA-Z0-9]+\)`, lines[3])
	assert.Regexp(t, `\s+/[a-zA-Z\-/]+/basic/runtime/callstack/call_stack_test.go:\d+ \+0x[a-zA-Z0-9]+`, lines[4])

	// 测试在 defer 中获取调用堆栈

	panicFunc := func() {
		panic(errors.New("test panic"))
	}

	defer func() { // 处理 panic 的 defer 调用
		assert.Error(t, recover().(error))
		s := CallStack()

		lines = strings.Split(s, string(lineBreak))
		assert.Regexp(t, `goroutine \d+ \[running\]:`, lines[0])

		assert.Equal(t, "study-golang/basic/runtime/callstack.CallStack()", lines[1])
		assert.Regexp(t, `\s+/[a-zA-Z\-/]+/basic/runtime/callstack/call_stack.go:\d+ \+0x\d+`, lines[2])

		assert.Regexp(t, `study-golang/basic/runtime/callstack.TestGetCallStack.func2()`, lines[3])
		assert.Regexp(t, `\s+/[a-zA-Z\-/]+/basic/runtime/callstack/call_stack_test.go:\d+ \+0x\d+`, lines[4])

		assert.Regexp(t, `panic\(\{0x[a-zA-Z0-9?]+, 0x[a-zA-Z0-9?]+\}\)`, lines[5]) // panic 发生的位置
	}()

	panicFunc() // 调用引发 panic 的函数
}
