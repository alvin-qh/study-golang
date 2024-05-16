package caller

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试通过闭包获取
func TestGetCallerInfo(t *testing.T) {
	dir, err := os.Getwd()
	assert.Nil(t, err)

	dir = filepath.Join(dir, "callerstate_test.go")

	// 输出当前调用信息
	cs, err := Where()
	assert.Nil(t, err)
	assert.Equal(t, "study/basic/runtime/callerstate.TestGetCallerInfo", cs.FuncName)
	assert.Equal(t, 19, cs.LineNo)
	assert.Equal(t, dir, cs.FileName)

	func() {
		// 输出当前调用信息
		cs, err = Where()
		assert.Nil(t, err)
		assert.Equal(t, "study/basic/runtime/callerstate.TestGetCallerInfo.func1", cs.FuncName)
		assert.Equal(t, 27, cs.LineNo)
		assert.Equal(t, dir, cs.FileName)
	}()

	// 输出当前调用信息
	cs, err = Where()
	assert.Nil(t, err)
	assert.Equal(t, "study/basic/runtime/callerstate.TestGetCallerInfo", cs.FuncName)
	assert.Equal(t, 35, cs.LineNo)
	assert.Equal(t, dir, cs.FileName)

	assert.Equal(t, "study/basic/runtime/callerstate.TestGetCallerInfo:"+dir+"(35)", cs.String())
}

// 测试通过闭包获取
func TestGetCallerStackInfo(t *testing.T) {
	dir, err := os.Getwd()
	assert.Nil(t, err)

	dir = filepath.Join(dir, "callerstate_test.go")

	// 输出当前调用信息
	cs := ListStackInfo(10)
	assert.Equal(t, "study/basic/runtime/callerstate.TestGetCallerStackInfo", cs[0].FuncName)
	assert.True(t, cs[0].LineNo >= 52 && cs[0].LineNo <= 53)
	assert.Equal(t, dir, cs[0].FileName)

	assert.Equal(t, "testing.tRunner", cs[1].FuncName)
	assert.Regexp(t, `.+?[\\/]testing[\\/]testing.go`, cs[1].FileName)

	func() {
		// 输出当前调用信息
		cs := ListStackInfo(10)
		assert.Equal(t, "study/basic/runtime/callerstate.TestGetCallerStackInfo.func1", cs[0].FuncName)
		assert.Equal(t, 63, cs[0].LineNo)
		assert.Equal(t, dir, cs[0].FileName)

		assert.Equal(t, "study/basic/runtime/callerstate.TestGetCallerStackInfo", cs[1].FuncName)
		assert.Regexp(t, `.+?[\\/]basic[\\/]runtime[\\/]callerstate[\\/]callerstate_test.go`, cs[1].FileName)
	}()

	// 输出当前调用信息
	cs = ListStackInfo(10)
	assert.Equal(t, "study/basic/runtime/callerstate.TestGetCallerStackInfo", cs[0].FuncName)
	assert.True(t, cs[0].LineNo >= 72 && cs[0].LineNo <= 73)
	assert.Equal(t, dir, cs[0].FileName)

	assert.Equal(t, "testing.tRunner", cs[1].FuncName)
	assert.Regexp(t, `.+?[\\/]src[\\/]testing[\\/]testing.go`, cs[1].FileName)
}

// 测试获取当前文件名称
func TestGetCurrentGoFile(t *testing.T) {
	dir, err := os.Getwd()
	assert.Nil(t, err)

	dir = filepath.Join(dir, "callerstate_test.go")

	// 获取当前文件路径名称
	fn, err := GetCurrentGoFile()
	assert.Nil(t, err)
	assert.Equal(t, dir, fn)
}
