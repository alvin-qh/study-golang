// context.Context 对象作为参数, 传入协程函数, 用于对协程进行控制, 包括: 设置超时时间, 取消协程等
//
// Deadline: 返回 context.Context 被取消的时间, 也就是完成工作的截止日期;
// Done: 返回一个 Channel, 这个 Channel 会在当前工作完成或者上下文被取消后关闭, 多次调用 Done 方法会返回同一个 Channel;
// Err: 返回 context.Context 结束的原因, 它只会在 Done 方法对应的 Channel 关闭时返回非空的值:
//
//	如果 context.Context 被取消, 会返回 Canceled 错误;
//	如果 context.Context 超时, 会返回 DeadlineExceeded 错误;
//
// Value: 从 context.Context 中获取键对应的值, 对于同一个上下文来说, 多次调用 Value 并传入相同的 Key 会返回相同的结果,
// 该方法可以用来传递请求特定的数据
package context

import (
	"context"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	runtime.GOMAXPROCS(0)
}

// 测试通过 Context 实例发送取消信号
//
// 可以通过向 goroutine 传递一个 Context 实例, 可以在 goroutine 中通过 Context.Done() 方法来检查是否被取消,
// 如果被取消, 则 goroutine 可以通过 context.Done() 返回的 chan 接收到信号, 从而令 goroutine 退出
func TestContext_Cancel(t *testing.T) {
	// 创建可取消 Context 实例
	ctx, cancel := context.WithCancel(context.Background())

	now := time.Now()

	var wg sync.WaitGroup
	wg.Add(1)

	// 启动 goroutine
	go func() {
	exit:
		for {
			select {
			// 等待 context 发送 cancel 信号
			case <-ctx.Done():
				wg.Done()
				// 退出循环
				break exit
			default:
				<-time.After(10 * time.Millisecond)
			}
		}
	}()

	time.Sleep(100 * time.Millisecond)
	cancel()

	wg.Wait()

	s := time.Since(now)
	assert.GreaterOrEqual(t, s.Milliseconds(), int64(100))
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
