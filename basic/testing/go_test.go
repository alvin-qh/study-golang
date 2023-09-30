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
// 测试主函数会在所有测试执行前执行, 并通过其中的 `m.Run()` 执行实际的测试函数
func TestMain(m *testing.M) {
	// 在所有测试前执行
	setup()

	code := m.Run()

	// 在所有测试后执行
	tearDown()

	os.Exit(code)
}

func setup() {
	_value = 100
}

func tearDown() {
	_value = 0
}

// 验证 `TestMain` 函数执行
func TestTestMainFunc(t *testing.T) {
	assert.Equal(t, 100, _value)
}

// 验证当执行测试时, 会具有一个名为 `"test.v"` 的标记表示正在执行测试
func TestLookupFlagIfRunTesting(t *testing.T) {
	assert.NotNil(t, flag.Lookup("test.v"))
}
