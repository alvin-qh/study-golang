package callerstate

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试获取当前调用函数的调用栈情况
//
// 通过 `Where` 函数返回一个 `CallerState` 实例, 表示了当前调用函数的堆栈情况
func TestCallerState_Where(t *testing.T) {
	dir, err := os.Getwd()
	assert.Nil(t, err)

	dir = filepath.Join(dir, "callerstate_test.go")

	// 输出当前调用信息
	cs, err := Where()
	assert.Nil(t, err)
	assert.Equal(t, "study/basic/runtime/callerstate.TestCallerState_Where", cs.FuncName)
	assert.Equal(t, 21, cs.LineNo)
	assert.Equal(t, dir, cs.FileName)

	func() {
		// 输出当前调用信息
		cs, err = Where()
		assert.Nil(t, err)
		assert.Equal(t, "study/basic/runtime/callerstate.TestCallerState_Where.func1", cs.FuncName)
		assert.Equal(t, 29, cs.LineNo)
		assert.Equal(t, dir, cs.FileName)
	}()

	// 输出当前调用信息
	cs, err = Where()
	assert.Nil(t, err)
	assert.Equal(t, "study/basic/runtime/callerstate.TestCallerState_Where", cs.FuncName)
	assert.Equal(t, 37, cs.LineNo)
	assert.Equal(t, dir, cs.FileName)

	assert.Equal(t, "study/basic/runtime/callerstate.TestCallerState_Where:"+dir+"(37)", cs.String())
}

// 获取调用方堆栈信息
//
// 通过 `ListStackInfo` 函数可以获取到调用方的整体堆栈信息, 该函数返回 `[]*CallerState` 切片, 表示调用栈每帧的情况
//
// 结果中的第 `0` 项即为当前调用函数的栈信息
func TestCallerState_ListStackInfo(t *testing.T) {
	dir, err := os.Getwd()
	assert.Nil(t, err)

	dir = filepath.Join(dir, "callerstate_test.go")

	// 输出当前调用信息
	cs := ListStackInfo(10)
	assert.Equal(t, "study/basic/runtime/callerstate.TestCallerState_ListStackInfo", cs[0].FuncName)
	assert.True(t, cs[0].LineNo >= 58 && cs[0].LineNo <= 59)
	assert.Equal(t, dir, cs[0].FileName)

	assert.Equal(t, "testing.tRunner", cs[1].FuncName)
	assert.Regexp(t, `.+?[\\/]testing[\\/]testing.go`, cs[1].FileName)

	func() {
		// 输出当前调用信息
		cs := ListStackInfo(10)
		assert.Equal(t, "study/basic/runtime/callerstate.TestCallerState_ListStackInfo.func1", cs[0].FuncName)
		assert.Equal(t, 68, cs[0].LineNo)
		assert.Equal(t, dir, cs[0].FileName)

		assert.Equal(t, "study/basic/runtime/callerstate.TestCallerState_ListStackInfo.func1", cs[1].FuncName)
		assert.Regexp(t, `.+?[\\/]basic[\\/]runtime[\\/]callerstate[\\/]callerstate_test.go`, cs[1].FileName)
	}()

	// 输出当前调用信息
	cs = ListStackInfo(10)
	assert.Equal(t, "study/basic/runtime/callerstate.TestCallerState_ListStackInfo", cs[0].FuncName)
	assert.True(t, cs[0].LineNo >= 78 && cs[0].LineNo <= 79)
	assert.Equal(t, dir, cs[0].FileName)

	assert.Equal(t, "testing.tRunner", cs[1].FuncName)
	assert.Regexp(t, `.+?[\\/]src[\\/]testing[\\/]testing.go`, cs[1].FileName)
}

// 测试获取当前文件名称
func TestCallerState_GetCurrentGoFile(t *testing.T) {
	dir, err := os.Getwd()
	assert.Nil(t, err)

	dir = filepath.Join(dir, "callerstate_test.go")

	// 获取当前文件路径名称
	fn, err := GetCurrentGoFile()
	assert.Nil(t, err)
	assert.Equal(t, dir, fn)
}
