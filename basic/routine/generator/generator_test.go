package generator

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func init() {
	runtime.GOMAXPROCS(0)
}

// 测试通过 channel 创建的生成器
func TestChanGenerator(t *testing.T) {
	// 创建生成器对象

	// 传入生成函数作为参数
	g := New(func(ch chan int) int {
		n := 0
		for {
			ch <- n
			n++
		}
	})
	defer g.Close()

	x := 0

	// 通过生成器生成 100 个数据
	// 通过 Next 函数, 返回生成器生成出的数据, 并获取生成器下一个数据
	for n, err := g.Next(0); err == nil && n < 100; n, err = g.Next(0) {
		assert.Nil(t, err)

		// 检查生成数据
		assert.Equal(t, x, n)
		x++
	}

	assert.Equal(t, 100, x)
	x++

	// 继续生成 100 个数据
	// 通过 Range 函数, 返回一个 channel, 通过对其进行 range 操作进行生成
	for n := range g.Range() {
		// 检查生成数据
		assert.Equal(t, x, n)

		if x > 100 {
			break
		}
		x++
	}
}
