package long

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试 long.Long 类型
// long.Long 类型相当于为 int64 类型起了个别名, 根据 Go 语言规范, 重新定义的类型即可为其定义相关的类型函数
// Go 语言不支持运算符重载, 所以定义了 long.Long 类型严格上和 int64 类型并不是同一类型 (虽然值是完全一致的), 赋值时需要类型转化

// 测试 Compare 函数
func TestCompare(t *testing.T) {
	l1 := Long(100)
	assert.True(t, l1.Compare(l1) == 0)

	l2 := Long(200)
	assert.True(t, l1.Compare(l2) < 0)

	l3 := Long(50)
	assert.True(t, l1.Compare(l3) > 0)

	l4 := l1
	assert.True(t, l1.Compare(l4) == 0)
}

// 测试 String 函数
func TestToString(t *testing.T) {
	var l Long
	assert.Equal(t, "0", l.String())

	l = Long(100)
	assert.Equal(t, "100", l.String())
}
