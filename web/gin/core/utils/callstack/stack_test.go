package callstack

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试 `source` 函数
//
// 确认当 `lines` 参数为 `nil` 时, 返回预设值
func TestSourceIfLinesIsNil(t *testing.T) {
	line := source(nil, 2)
	assert.Equal(t, "???", string(line))
}

// 测试 `source` 函数
//
// 确认当 `n` 参数为 `0` 时, 返回预设值
func TestSourceIfNIsZero(t *testing.T) {
	line := source([][]byte{
		[]byte("abcd"),
	}, 0)

	assert.Equal(t, "???", string(line))
}

// 测试 `source` 函数
//
// 确认参数 `n` 的数值超出参数 `lines` 的范围, 返回预设值
func TestSourceIfNOutofRange(t *testing.T) {
	line := source([][]byte{
		[]byte("abcd"),
	}, 100)

	assert.Equal(t, "???", string(line))
}

// 测试 `source` 函数
//
// 确认返回正确索引的行
func TestSource(t *testing.T) {
	line := source([][]byte{
		[]byte("aaaa"),
		[]byte("bbbb"),
		[]byte("cccc"),
		[]byte("dddd"),
	}, 2)

	assert.Equal(t, "bbbb", string(line))
}

// 测试 `function` 函数
//
// 确认传入 `pc` 参数无效时, 返回预设结果
func TestFunctionWithInvalidPC(t *testing.T) {
	f := function(0x1234)
	assert.Equal(t, "???", string(f))
}

// 测试 `function` 函数
//
// 确认传入调用地址, 返回调用函数名称
func TestFunction(t *testing.T) {
	pc, _, _, ok := runtime.Caller(0)
	assert.True(t, ok)

	f := function(pc)
	assert.Equal(t, "TestFunction", string(f))
}

// 测试 `CallStack` 函数
//
// 确认 `CallStack` 函数返回当前的调用栈内容
func TestCallStack(t *testing.T) {
	stack := CallStack(0)

	assert.Contains(t, string(stack), "/core/utils/callstack/stack.go")
	assert.Contains(t, string(stack), "gin/core/utils/callstack/stack_test.go")
	assert.Contains(t, string(stack), "TestCallStack: stack := CallStack(0)")
}
