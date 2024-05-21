package semaphore

import (
	"context"
	"runtime"
	"study/basic/builtin/slice/utils"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/semaphore"
)

func init() {
	runtime.GOMAXPROCS(0)
}

// 测试信号量
//
// 本例创建具有 10 个通过逐个释放信号量
func TestSemaphore_Weighted(t *testing.T) {
	ctx := context.Background()

	// 声明一个具备 10 个值的信号量实例
	sem := semaphore.NewWeighted(10)
	// 将 10 个信号量值全部占用掉
	sem.Acquire(ctx, 10)

	var wg sync.WaitGroup
	wg.Add(2)

	// 启动 goroutine, 对 10 个信号量值逐一释放
	go func() {
		n := 10

		// 循环 10 次, 每次释放一个信号量
		for n > 0 {
			time.Sleep(10 * time.Millisecond)
			sem.Release(1)

			n--
		}

		wg.Done()
	}()

	// 上一次获取信号量的时间
	last := time.Now()
	intervals := make([]int64, 0, 10)

	// 启动 goroutine, 对 10 个信号量逐一占用
	go func() {
		n := 10

		// 循环 10 次, 每次占用一个信号量
		for n > 0 {
			sem.Acquire(ctx, 1)

			intervals = append(intervals, time.Since(last).Milliseconds())
			last = time.Now()

			n--
		}

		wg.Done()
	}()

	wg.Wait()

	// 确认占用了 10 个信号量
	assert.Len(t, intervals, 10)

	// 确认当前所有的信号量都已经被占用
	assert.False(t, sem.TryAcquire(1))

	// 确认每次释放和占用信号量的时间间隔都为 10 毫秒
	assert.Equal(t, utils.Repeat(10, int64(10)), intervals)
}
