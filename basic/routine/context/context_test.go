package context

import (
	"context"
	"errors"
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
// 通过 `context.WithCancel` 函数创建一个 Context 实例和一个 `CancelFunc` 类型的取消函数, 调用取消
// 函数可以通过 `Context` 对象发送一个取消信号, 该取消信号可以通过 `Context.Done()` 函数返回的 `chan`
// 来接受 (接收到的内容为 `struct{}`, 没有实际意义)
//
// Context 实例作为上下文需要传递到对应的所有 goroutine 中, 并在 goroutine 中接收 `Context.Done()`
// 返回的 `chan` 传递的信息
//
// 一旦在 goroutine 中接收到 Context 的取消信号, 就意味着 goroutine 不应该继续执行下去了, 此时应该
// 结束 goroutine 函数
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

// 设置和获取 Context 发送取消信号的原因
//
// 通过 `context.WithCancelCause` 函数可以返回一个 Context 实例和一个 `CancelCauseFunc` 类型
// 的取消函数
//
// 和 `context.WithCancel` 函数的使用方式基本类似, 区别仅在于 `CancelCauseFunc` 类型取消函数具备
// 一个 `error` 类型的参数, 表示取消原因
//
// 当取消函数调用后, 可以通过 `context.Cause` 函数从 Context 实例中获取表示取消原因的 `error` 实例
func TestContext_CancelReason(t *testing.T) {
	// 设置一个可取消并可设置取消原因的 Context 实例
	// 返回的取消函数可以设置一个 error 类型的实例表示取消原因
	ctx, cancel := context.WithCancelCause(context.Background())

	var wg sync.WaitGroup
	wg.Add(2)

	// 启动 goroutine, 并等待 Context 实例的取消信号
	go func() {
		// 等待取消信号
		<-ctx.Done()
		wg.Done()
	}()

	// 启动 goroutine, 在等待后通过 Context 实例发送取消信号并设置取消原因
	go func() {
		time.Sleep(100 * time.Millisecond)

		// 通过 Context 实例发送取消信号并设置取消原因
		cancel(errors.New("wait timeout"))
		wg.Done()
	}()

	wg.Wait()

	// 获取 Context 实例中的取消原因
	assert.EqualError(t, context.Cause(ctx), "wait timeout")
}

// 测试通过 Context 实例设定超时时间
//
// 通过 `context.WithTimeout` 函数可以设置一个超时时间, 并返回一个 Context 实例和一个 `CancelFunc` 类型
// 的取消函数
//
// 取消函数的作用和 `context.WithCancel` 函数返回的取消函数作用一致, 一般需要用 `defer` 关键字进行调用, 以
// 保证这个取消函数最终会被调用
//
// 当超时时间到达后, Context 实例会通过 `Context.Done()` 函数返回的 `chan` 发送一个信号 (和调用取消函数发
// 送的信号一致), 并通过该信号令 goroutine 退出
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

// 为超时设置一个原因
//
// 可以通过 `context.WithTimeoutCause` 函数设置超时时间并指定一个 `error` 类型实例表示超时原因,
// 当超时取消的信号发送后, 可以通过 `context.Cause` 函数从 Context 实例中获取到表示超时原因的 `error` 实例
func TestContext_TimeoutReason(t *testing.T) {
	// 创建具备超时功能的 Context 实例, 设定超时原因, 返回 Context 实例及取消函数
	ctx, cancel := context.WithTimeoutCause(context.Background(), 100*time.Millisecond, errors.New("wait to long"))
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

	// 确定 Context 实例中的超时原因
	assert.EqualError(t, context.Cause(ctx), "wait to long")
}

// 测试通过 Context 实例设定截至时间
//
// 通过 `context.WithDeadline` 函数可以设置一个截至时间, 并返回一个 Context 实例和一个 `CancelFunc` 类型
// 的取消函数
//
// 整个流程和使用 `context.WithTimeout` 函数基本类似, 只是 `context.WithDeadline` 函数设置的是截至时间,
// 是一个明确的日期时间; 而 `context.WithTimeout` 函数则设置的是一个从当前时间开始计算的时长
func TestContext_Deadline(t *testing.T) {
	// 定义一个未来的时间表示截至时间
	future := time.Now().Add(100 * time.Millisecond)

	// 创建一个具备截止时间的 Context 实例, 并返回一个取消函数
	ctx, cancel := context.WithDeadline(context.Background(), future)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	// 记录 goroutine 启动的时间
	now := time.Now()

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

	// 获取上下文中设置的截至时间
	dl, ok := ctx.Deadline()
	assert.True(t, ok)
	assert.Equal(t, future, dl)

	wg.Wait()

	// 计算 goroutine 的整体执行时间, 应该大于 100ms, 即 100ms 后 Context 截至, goroutine 结束
	s := time.Since(now)
	assert.GreaterOrEqual(t, s.Milliseconds(), int64(100))
}
func TestContext_DeadlineReason(t *testing.T) {
	// 定义一个未来的时间表示截至时间
	future := time.Now().Add(100 * time.Millisecond)

	// 创建一个具备截止时间的 Context 实例, 设定截至原因, 并返回一个取消函数
	ctx, cancel := context.WithDeadlineCause(context.Background(), future, errors.New("time is up"))
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	// 记录 goroutine 启动的时间
	now := time.Now()

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

	// 获取上下文中设置的截至时间
	dl, ok := ctx.Deadline()
	assert.True(t, ok)
	assert.Equal(t, future, dl)

	wg.Wait()

	// 计算 goroutine 的整体执行时间, 应该大于 100ms, 即 100ms 后 Context 截至, goroutine 结束
	s := time.Since(now)
	assert.GreaterOrEqual(t, s.Milliseconds(), int64(100))

	// 确定 Context 实例中的截至原因
	assert.EqualError(t, context.Cause(ctx), "time is up")
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
