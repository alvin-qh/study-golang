package once_test

import (
	"study/basic/testing/assertion"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试只能执行一次的函数
//
// 如果有一个函数需要在 goroutine 中执行, 但该函数只能执行一次 (例如初始化函数), 则此时有三种方案:
//   - 在所有 goroutine 执行前找一个机会执行该函数, 这样会导致该函数不是异步执行;
//   - 通过原子标志位标记该函数已经执行过, 这样会导致代码比较复杂;
//   - 通过 `sync.OnceFunc` 函数保障目标函数, 这样得到的函数永远只会执行一次, 但被包装的原始函数不受影响
//
// `sync.OnceFunc` 函数返回的结果是同步的, 即多个 goroutine 同时执行该函数时, 只会有一个 goroutine
// 执行成功
func TestOnce_OnceFunc(t *testing.T) {
	// 定义一个计数器, 用于统计被包装函数的执行次数
	count := 0

	// 定义 Once 方法
	fn := sync.OnceFunc(func() {
		count++
	})

	// 定义等待组对象, 用于等待全部任务完成
	var wg sync.WaitGroup

	// 启动 10 个 goroutine
	for range 10 {
		// 启动 goroutine, 在其中执行一次 OnceFunc
		wg.Go(func() { fn() })
	}

	// 等待任务结束
	wg.Wait()

	// 确认最终 OnceFunc 只被执行了一次
	assert.Equal(t, 1, count)
}

// 测试永远返回相同值的函数
//
// `sync.OnceValue` 函数用于包装一个具备返回值的函数, 返回一个函数结果, 该返回结果函数无论调用多少次,
// 其返回值永远和该函数第一次调用的返回值一致
//
// `sync.OnceValue` 返回的结果也是同步的, 即无论多少 goroutine 调用该函数, 永远只有第一个调用成功的 goroutine 返回值生效
func TestOnce_OnceValue(t *testing.T) {
	// 包装函数, 被包装的函数返回当前时间
	fn := sync.OnceValue(func() time.Time {
		return time.Now()
	})

	// 定义等待组对象, 用于等待全部任务完成
	var wg sync.WaitGroup

	// 保持每次调用函数结果
	rs := make([]time.Time, 10)

	// 启动 10 个 goroutine
	for i := range rs {
		// 启动 goroutine, 在其中重复调用前面产生的函数
		wg.Go(func() {
			time.Sleep(10 * time.Millisecond)
			rs[i] = fn()
		})
	}

	// 等待所有任务结束
	wg.Wait()

	// 确认所有函数调用结果都是一样的
	assertion.All(t, rs, fn())
}

// 测试永远返回相同值的函数
//
// `sync.OnceValues` 的作用和 `sync.OnceValue` 类似, 只是包装的函数需返回两个返回值
func TestOnce_OnceValues(t *testing.T) {
	start := time.Now()

	// 包装函数, 被包装的函数返回当前时间
	fn := sync.OnceValues(func() (time.Time, int) {
		return time.Now(), int(time.Since(start).Milliseconds())
	})

	// 定义等待组对象, 用于等待全部任务完成
	var wg sync.WaitGroup

	// 保持每次调用函数结果
	rs := make([][2]any, 10)

	// 启动 10 个 goroutine
	for i := 0; i < len(rs); i++ {
		wg.Add(1)

		// 启动 goroutine, 在其中重复调用前面产生的函数
		go func() {
			defer wg.Done()

			time.Sleep(10 * time.Millisecond)
			d, n := fn()
			rs[i] = [2]any{d, n}
		}()
	}

	// 等待所有任务结束
	wg.Wait()

	// 确认所有函数调用结果都是一样的
	d, n := fn()
	assertion.All(t, rs, [2]any{d, n})
}
