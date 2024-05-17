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

// 在使用 Value Context 时, Key 的类型推荐使用自定义类型
type ContextKey string

// 测试通过 Context 实例携带键值对
//
// 可以向 Context 中存储一些键值对, 所有 goroutine 可以通过该 Context 对象访问这些键值对
func TestContext_Value(t *testing.T) {
	// 为 Context 实例设定键值对
	ctx := context.WithValue(context.Background(), ContextKey("num"), 100)
	ctx = context.WithValue(ctx, ContextKey("name"), "Alvin")

	var wg sync.WaitGroup
	wg.Add(1)

	// 启动 goroutine, 在其中通过 Context 实例获取键值对
	go func() {
		defer wg.Done()

		// 根据 key 获取 Value
		num := ctx.Value(ContextKey("num")).(int)
		name := ctx.Value(ContextKey("name")).(string)

		assert.Equal(t, 100, num)
		assert.Equal(t, "Alvin", name)
	}()

	// 等待 goroutine 结束
	wg.Wait()
}

// 测试通过 Context 实例发送取消信号
//
// 通过 context.WithCancel 函数创建一个 Context 实例并传递给 goroutine, 可以在 goroutine
// 中通过 Context.Done() 方法来检查是否被取消, 如果被取消, 则 goroutine 可以通过 context.Done()
// 返回的 chan 接收到信号, 从而令 goroutine 退出
func TestContext_Cancel(t *testing.T) {
	// 创建可取消 Context 实例
	// 返回一个上下文实例即一个 cancel 函数, 通过该函数可以发送取消指令
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)

	// 记录 goroutine 启动的时间
	now := time.Now()

	// 启动 goroutine1
	go func() {
		defer wg.Done()

	exit:
		for {
			select {
			case <-ctx.Done():
				// Context 接收到取消信号, 退出 goroutine
				break exit
			case <-time.After(10 * time.Millisecond):
				// 等待 10ms, 如果还未有取消信号, 则重新循环
			}
		}
	}()

	// 启动另一个 goroutine2, 在其中发送取消信号
	go func() {
		// 等待 100ms 后, 执行取消函数发送取消信号
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	// 等待 goroutine 结束
	wg.Wait()

	// 计算 goroutine1 的整体执行时间, 应该大于 100ms,
	// 这也是 goroutine1 从启动到收到取消信号的时间
	s := time.Since(now)
	assert.GreaterOrEqual(t, s.Milliseconds(), int64(100))
}

// 测试具备超时功能的 Context
//
// 通过 context.WithTimeout 函数创建一个 Context 实例并传递给 goroutine, 可以在 goroutine
// 中通过 Context.Done() 方法来检查是否被取消, 如果被取消, 则 goroutine 可以通过 context.Done()
// 返回的 chan 接收到信号, 从而令 goroutine 退出
func TestContext_Timeout(t *testing.T) {
	// 创建具备超时功能的 Context 实例, 返回 Context 实例及取消函数
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	// 记录 goroutine 启动的时间
	now := time.Now()

	// 启动 goroutine
	go func() {
		defer wg.Done()

	exit:
		for {
			select {
			case <-ctx.Done():
				// Context 超时或接收到取消信号时, 退出 goroutine
				break exit
			case <-time.After(10 * time.Millisecond):
				// 等待 10ms, 如果还未有取消信号, 则重新循环
			}
		}
	}()

	wg.Wait()

	// 计算 goroutine 的整体执行时间, 应该大于 100ms, 即 100ms 后 Context 超时, goroutine 结束
	s := time.Since(now)
	assert.GreaterOrEqual(t, s.Milliseconds(), int64(100))
}

// 在 Context 取消或超时时执行异步回调函数
//
// 如果一个 Context 可以被取消或允许超时, 则可以为其绑定一个 `AfterFunc` 回调函数, 当 Context 被取消或超时后,
// 该回调函数会在一个新的 goroutine 中调用
//
// `context.AfterFunc` 函数会返回一个 `stop` 函数, 用于停止回调函数执行
func TestContext_AfterFunc(t *testing.T) {
	// 创建一个可取消 Context 实例
	ctx, cancel := context.WithCancel(context.Background())

	// 定义等待组, 包含两个任务
	var wg sync.WaitGroup
	wg.Add(2)

	// 定义回调函数, 在 Context 被取消后执行
	// 返回的 stop 函数用于停止该回调函数执行
	stop := context.AfterFunc(ctx, func() {
		// 完成一个任务
		wg.Done()
	})

	// 应该通过 defer 保证 stop 函数的调用
	defer stop()

	// 启动 goroutine, 并在 Context 取消后结束
	go func() {
		<-ctx.Done()
		wg.Done()
	}()

	// 通过 Context 发送取消信号
	cancel()

	// 等待所有 goroutine 结束
	wg.Wait()
}
