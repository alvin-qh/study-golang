package closure

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试函数闭包
//
// 闭包函数定义了一个作用域, 但也可以使用外部作用域定义的变量
func TestFunction_Closure(t *testing.T) {
	x := 10
	y := 20

	// 在闭包内使用外部变量值
	f1 := func(z int) int { return z + x + y }
	r := f1(30)
	assert.Equal(t, 60, r)

	// 在闭包内修改外部变量值
	f2 := func(a, b int) {
		x = a
		y = b
	}

	f2(100, 200)
	assert.Equal(t, 100, x)
	assert.Equal(t, 200, y)

	// 定义闭包并直接执行
	r = func() int { return x * y }()
	assert.Equal(t, 20000, r)
}
