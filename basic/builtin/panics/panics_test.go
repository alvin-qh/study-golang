package panics

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// defer，即延迟执行语句，在一般情况下，用于资源的释放，所以 defer 调用会在函数的其它代码都执行后被调用
// 如果函数中包含多个 defer 调用，则会按照 LIFO 的顺序依次调用
func TestDeferCall(t *testing.T) {
	// defer 调用会在函数执行完毕后依次调用，即函数中的其它代码都执行完毕并返回后，才会依次执行 defer 调用
	// 调用顺序为 后进先出 LIFO，即越靠后的 defer 调用越先执行
	deferOrder := func(l *List) {
		l.Append("1") // 第 1 步

		defer l.Append("2") // 第 4 步
		defer l.Append("3") // 第 3 步

		l.Append("4")                                // 第 2 步
		assert.Equal(t, []string{"1", "4"}, l.slice) // 此时由于函数尚未结束，所以 defer 调用尚未执行
	}

	l := List{slice: make([]string, 0, 4)}

	deferOrder(&l)
	assert.Equal(t, []string{"1", "4", "3", "2"}, l.slice)
}
