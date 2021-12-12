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
	assert.NoError(t, err)

	dir = filepath.Join(dir, "caller_test.go")

	// 输出当前调用信息
	cs, err := Where()
	assert.NoError(t, err)
	assert.Equal(t, "basic/runtime/caller.TestGetCallerInfo", cs.FuncName)
	assert.Equal(t, 19, cs.LineNo)
	assert.Equal(t, dir, cs.FileName)

	func() {
		// 输出当前调用信息
		cs, err = Where()
		assert.NoError(t, err)
		assert.Equal(t, "basic/runtime/caller.TestGetCallerInfo.func1", cs.FuncName)
		assert.Equal(t, 27, cs.LineNo)
		assert.Equal(t, dir, cs.FileName)
	}()

	// 输出当前调用信息
	cs, err = Where()
	assert.NoError(t, err)
	assert.Equal(t, "basic/runtime/caller.TestGetCallerInfo", cs.FuncName)
	assert.Equal(t, 35, cs.LineNo)
	assert.Equal(t, dir, cs.FileName)

	assert.Equal(t, "basic/runtime/caller.TestGetCallerInfo:"+dir+"(35)", cs.String())
}

// 测试通过闭包获取
func TestGetCallerStackInfo(t *testing.T) {
	dir, err := os.Getwd()
	assert.NoError(t, err)

	dir = filepath.Join(dir, "caller_test.go")

	// 输出当前调用信息
	cs := ListStackInfo()
	assert.Equal(t, "basic/runtime/caller.TestGetCallerStackInfo", cs[0].FuncName)
	assert.Equal(t, 52, cs[0].LineNo)
	assert.Equal(t, dir, cs[0].FileName)

	assert.Equal(t, "testing.tRunner", cs[1].FuncName)
	assert.Regexp(t, ".+?/src/testing/testing.go", cs[1].FileName)

	func() {
		// 输出当前调用信息
		cs := ListStackInfo()
		assert.Equal(t, "basic/runtime/caller.TestGetCallerStackInfo.func1", cs[0].FuncName)
		assert.Equal(t, 62, cs[0].LineNo)
		assert.Equal(t, dir, cs[0].FileName)

		assert.Equal(t, "basic/runtime/caller.TestGetCallerStackInfo", cs[1].FuncName)
		assert.Regexp(t, ".+?/basic/runtime/caller/caller_test.go", cs[1].FileName)
	}()

	// 输出当前调用信息
	cs = ListStackInfo()
	assert.Equal(t, "basic/runtime/caller.TestGetCallerStackInfo", cs[0].FuncName)
	assert.Equal(t, 72, cs[0].LineNo)
	assert.Equal(t, dir, cs[0].FileName)

	assert.Equal(t, "testing.tRunner", cs[1].FuncName)
	assert.Regexp(t, ".+?/src/testing/testing.go", cs[1].FileName)
}

// 测试获取当前文件名称
func TestGetCurrentGoFile(t *testing.T) {
	dir, err := os.Getwd()
	assert.NoError(t, err)

	dir = filepath.Join(dir, "caller_test.go")

	fn, err := GetCurrentGoFile()
	assert.NoError(t, err)
	assert.Equal(t, dir, fn)
}
