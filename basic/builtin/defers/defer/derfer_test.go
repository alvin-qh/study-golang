package defers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试 defer 关键字
//
// Go 语言的 defer 关键字用于函数内部, 目的是将一个函数调用加入到一个栈结构中,
// 当当前函数结束后, 按出栈顺序依次调用通过 defer 关键字加入的函数调用
//
// 如此一来, 就可以保证 defer 关键字指定的函数调用会在当前函数结束后一定被自动
// 调用, 而无需在所有函数出口位置 (return, panic) 手动调用函数
//
// defer 关键字一般用于打开资源的释放
func TestDefer_InFunction(t *testing.T) {
	count := 0

	for range 100 {
		// 在循环中执行函数, 每次函数执行结束后, 函数内部的 defer 都会被执行
		func() {
			// 函数结束后进行调用, 将 count 变量加 1
			defer func() { count++ }()
		}()
	}
	assert.Equal(t, 100, count)

	count = 0

	func() {
		// 尽管在每个循环的范围内都使用了 defer
		// defer 后的函数也不会在每次循环结束后执行
		for range 100 {
			defer func() { count++ }()
		}

		// 循环结束后, 没有任何一个 defer 执行过
		assert.Equal(t, 0, count)
	}()

	// 只有函数结束后, 所有的 defer 才会被执行
	assert.Equal(t, 100, count)
}
