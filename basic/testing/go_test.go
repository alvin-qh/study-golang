package testing

import (
	"flag"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	_value = 0
)

// 测试主函数
//
// 测试主函数会在所有测试执行前执行, 并在该函数内部通过 `m.Run()` 启动实际的测试
func TestMain(m *testing.M) {
	// 在所有测试前执行
	setup()

	code := m.Run()

	// 在所有测试后执行
	tearDown()

	os.Exit(code)
}

// 在每个测试执行前执行
func setup() {
	_value = 100
}

// 在每个测试执行后执行
func tearDown() {
	_value = 0
}

// 测试 `TestMain` 函数
//
// 如果一个 Go 测试文件中包含 `TestMain` 函数, 则 Go 测试程序会测试执行前自动执行该函数
// (仅限当前文件中的测试)
//
// 在 `TestMain` 函数中, 必须通过参数 `m` (`testing.M` 类型) 来启动实际的测试, 所以
// 可以在实际测试执行前后进行一些处理
func TestTestMainFunc(t *testing.T) {
	assert.Equal(t, 100, _value)
}

// 测试执行时的全局表质量
//
// 通过 `flag` 包的 `Lookup` 方法, 可以查找名为 `"test.v"` 的标记, 如果该标记存在,
// 则表明当前运行环境为测试环境
func TestLookupFlagIfRunTesting(t *testing.T) {
	assert.NotNil(t, flag.Lookup("test.v"))
}
