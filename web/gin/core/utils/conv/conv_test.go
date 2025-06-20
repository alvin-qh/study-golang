package conv

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 定义一个实现 `fmt.Stringer` 接口的类型
type point struct {
	x, y int
}

// 实现 `fmt.Stringer` 接口的 `String` 方法
//
// 返回:
//   - 对象的字符串形式
func (p *point) String() string {
	return fmt.Sprintf("[%v,%v]", p.x, p.y)
}

// 测试 `AnyToString` 函数
//
// 验证将各种类型变量转为字符串后的结果
func TestAnyToString(t *testing.T) {
	assert.Equal(t, "100", AnyToString(100))
	assert.Equal(t, "1.2345", AnyToString(1.2345))
	assert.Equal(t, "true", AnyToString(true))
	assert.Equal(t, "100", AnyToString(uint(100)))
	assert.Equal(t, "Hello", AnyToString("Hello"))
	assert.Equal(t, "(10+20i)", AnyToString(complex(10, 20)))
	assert.Equal(t, "[10,20]", AnyToString(&point{10, 20}))
}
