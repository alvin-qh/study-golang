package value

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试 `Default` 函数
//
// 确认如果 `val` 参数为 `nil` 时, 将返回 `def` 参数的值
func TestDefault(t *testing.T) {
	assert.Equal(t, 100, Default(nil, 100))
	assert.Equal(t, 0, Default(0, 100))
	assert.Equal(t, "OK", Default(nil, "OK"))
	assert.Equal(t, "", Default("", "Hello"))
}

// 测试用结构体类型
type Point struct {
	x int
	y int
}

// 测试 `Join` 函数
//
// 确认可以将任意类型元素的数组连接为字符串
func TestArrayToString(t *testing.T) {
	assert.Equal(t, "1,2,3", Join([]int{1, 2, 3}, ","))
	assert.Equal(t, "A,B,C", Join([]string{"A", "B", "C"}, ","))
	assert.Equal(t, "{1 2},{3 4},{5 6}", Join([]Point{{1, 2}, {3, 4}, {5, 6}}, ","))
}

// 测试 `JoinAny` 函数
//
// 确认可以将任意类型数组连接为字符串
func TestJoinAnyToString(t *testing.T) {
	assert.Equal(t, "Hello", JoinAny("Hello", ","))
	assert.Equal(t, "A,B,C", JoinAny([]string{"A", "B", "C"}, ","))
	assert.Equal(t, "1,2,3", JoinAny([]int{1, 2, 3}, ","))
	assert.Equal(t, "{1 2},{3 4},{5 6}", Join([]Point{{1, 2}, {3, 4}, {5, 6}}, ","))
}
