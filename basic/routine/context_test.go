package routine

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	runtime.GOMAXPROCS(0)
}

// 测试 Cancel Context
func TestCancelContext(t *testing.T) {
	// 正常结束协程
	ctx, _, ch := CreateCancelContext() // 创建 context 对象

	go CancelContextHandler(ctx) // 调用协程函数
	ch <- "OK"                   // 向 ch 发送内容, 令协程正常结束

	r := <-ch // 从 ch 获取结果
	close(ch)
	assert.Equal(t, "Completed", r)

	// 通过 Context cancel 函数结束协程
	ctx, cancel, ch := CreateCancelContext() // 创建 context 对象, 返回 cancel 函数

	go CancelContextHandler(ctx)
	cancel() // 调用 cancel 函数结束协程

	r = <-ch // 从 ch 获取结果
	assert.Equal(t, "Canceled", r)
}

// 测试 Timeout Context
func TestTimeoutContext(t *testing.T) {
	// 正常结束协程
	ctx, _, ch := CreateTimeoutContext(time.Second) // 创建 context 对象

	go TimeoutHandler(ctx, time.Second) // 调用协程函数
	ch <- "OK"                          // 向 chan 发送数据, 令协程正常结束

	r := <-ch // 从 ch 获取结果
	close(ch)
	assert.Equal(t, "Completed by chan", r) // 表示正常结束

	// 在 Context 超时前结束协程
	ctx, _, ch = CreateTimeoutContext(time.Second) // 创建 context 对象

	go TimeoutHandler(ctx, 100*time.Millisecond) // 调用协程函数
	r = <-ch                                     // 从 ch 获取结果
	close(ch)
	assert.Equal(t, "Completed by timer", r)

	// 在 Context 超时之后结束协程
	ctx, _, ch = CreateTimeoutContext(100 * time.Millisecond)

	go TimeoutHandler(ctx, 200*time.Millisecond) // 调用协程函数

	r = <-ch // 从 ch 获取结果
	close(ch)
	assert.Equal(t, "Canceled by timeout", r)

	// 通过 Context cancel 函数结束协程
	ctx, cancel, ch := CreateTimeoutContext(time.Second)

	go TimeoutHandler(ctx, time.Second) // 调用协程函数
	cancel()

	r = <-ch // 从 ch 获取结果
	close(ch)
	assert.Equal(t, "Canceled by cancel called", r)
}
