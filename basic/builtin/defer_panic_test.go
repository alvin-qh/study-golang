package builtin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type List struct {
	slice []string
}

func (l *List) append(s string) {
	l.slice = append(l.slice, s)
}

// defer 调用会在函数执行完毕后依次调用，即函数中的其它代码都执行完毕并返回后，才会依次执行 defer 调用
// 调用顺序为 后进先出 LIFO，即越靠后的 defer 调用越先执行
func deferOrder(t *testing.T, lst *List) {
	lst.append("1") // 第 1 步

	defer lst.append("2") // 第 4 步
	defer lst.append("3") // 第 3 步

	lst.append("4") // 第 2 步

	assert.Equal(t, []string{"1", "4"}, lst.slice) // 此时由于函数尚未结束，所以 defer 调用尚未执行
}

// defer，即延迟执行语句，在一般情况下，用于资源的释放，所以 defer 调用会在函数的其它代码都执行后被调用
// 如果函数中包含多个 defer 调用，则会按照 LIFO 的顺序依次调用
func TestDeferCall(t *testing.T) {
	var lst = List{slice: make([]string, 0, 4)}

	deferOrder(t, &lst)
	assert.Equal(t, []string{"1", "4", "3", "2"}, lst.slice)
}
