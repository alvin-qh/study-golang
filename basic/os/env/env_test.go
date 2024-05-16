package env

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 根据环境变量的名称获取对应的环境变量值
func TestOS_SetAndGetEnv(t *testing.T) {
	// 设置环境变量
	err := os.Setenv("TEST_ENV", "Golang-Demo")
	assert.Nil(t, err)

	// 获取环境变量
	env := os.Getenv("TEST_ENV")
	assert.Equal(t, "Golang-Demo", env)
}

// 测试查找环境变量
//
// `os.LookupEnv` 函数用于查找指定名称的环境变量, 和 `os.Getenv` 函数的区别为,
// `os.LookupEnv` 会额外返回一个 bool 值来表示指定名称的环境变量是否被设置
//
// 这对于查找一些可能会设置为空字符串的环境变量有用
func TestOS_LookupEnv(t *testing.T) {
	// 设置环境变量
	err := os.Setenv("TEST_ENV", "Golang-Demo")
	assert.Nil(t, err)

	env, ok := os.LookupEnv("TEST_ENV")
	assert.True(t, ok)
	assert.Equal(t, "Golang-Demo", env)

	os.Unsetenv("TEST_ENV")

	_, ok = os.LookupEnv("TEST_ENV")
	assert.False(t, ok)
}

// 测试获取所有环境变量
//
// 参考 `sys.Environ` 函数
func TestOS_Environ(t *testing.T) {
	// 设置环境变量
	err := os.Setenv("TEST_ENV", "Golang-Demo")
	assert.Nil(t, err)

	// 获取所有环境变量
	envs := Environ()
	assert.NotEmpty(t, envs)

	// 查看是否包含了指定的环境变量
	env, ok := envs["TEST_ENV"]
	assert.True(t, ok)
	assert.Equal(t, "Golang-Demo", env)

	// 遍历结果, 确认每个环境变量正确
	for n, v := range envs {
		assert.Equal(t, os.Getenv(n), v)
	}
}

// 测试删除指定名称的环境变量
func TestOS_Unsetenv(t *testing.T) {
	// 设置环境变量
	err := os.Setenv("TEST_ENV", "Golang-Demo")
	assert.Nil(t, err)

	// 查看设置的环境变量
	assert.Equal(t, "Golang-Demo", os.Getenv("TEST_ENV"))

	// 删除指定名称的环境变量
	err = os.Unsetenv("TEST_ENV")
	assert.Nil(t, err)

	// 查看是否删除成功
	assert.Equal(t, "", os.Getenv("TEST_ENV"))
}

// 测试清空当前进程的所有环境变量
//
// 通过 `os.Clearenv` 函数可以清空当前的所有环境变量, 注意: 所谓清空只针对于当前进程, 其它进程不受影响
func TestOS_Clearenv(t *testing.T) {
	// 清空所有环境变量
	os.Clearenv()

	envs := os.Environ()
	assert.Empty(t, envs)
}

// 测试利用环境变量格式化字符串
//
// 可以在字符串中通过 `$NAME` 或者 `${NAME}` 占位符表示名称为 `NAME` 的环境变量,
// 格式化后, 占位符会被替换为实际的环境变量值
func TestOS_ExpandEnv(t *testing.T) {
	// 设置环境变量
	err := os.Setenv("TEST_ENV", "Golang-Demo")
	assert.Nil(t, err)

	// 利用环境变量格式化字符串
	s := os.ExpandEnv("TEST_ENV=$TEST_ENV")
	assert.Equal(t, "TEST_ENV=Golang-Demo", s)
}

// 测试通过一个回调函数替换字符串中的占位符
//
// 可以在字符串中通过 `$NAME` 或者 `${NAME}` 占位符表示名称为 `NAME` 的环境变量,
// 该占位符会传递给回调函数, 并替换为回调函数返回的字符串
func TestOS_Expand(t *testing.T) {
	// 设置环境变量
	err := os.Setenv("TEST_ENV", "Golang-Demo")
	assert.Nil(t, err)

	// 利用环境变量格式化字符串
	s := os.Expand("TEST_ENV=$TEST_ENV", func(name string) string {
		assert.Equal(t, "TEST_ENV", name)
		return os.Getenv(name)
	})
	assert.Equal(t, "TEST_ENV=Golang-Demo", s)
}
