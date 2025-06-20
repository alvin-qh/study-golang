package waitgroup

import (
	"basic/testing/assertion"
	"sync"
	"testing"
	"time"
)

// 测试等待组
//
// 等待组用于等待多个任务结束, 即同步多个任务的最终状态
func TestWaitGroup_Wait(t *testing.T) {
	// 定义等待组
	var wg sync.WaitGroup

	// 循环 10 次, 每次循环启动一个 goroutine
	for i := 0; i < 10; i++ {
		// 为等待组增加一个等待任务
		wg.Add(1)

		// 启动协程, 等待 10ms 后完成等待任务
		go func() {
			// 协程结束后完成等待任务
			defer wg.Done()

			time.Sleep(10 * time.Millisecond)
		}()
	}

	start := time.Now()

	// 等待所有任务完成
	wg.Wait()

	// 确定任务完成时长 100ms
	assertion.DurationMatch(t, 10*time.Millisecond, time.Since(start))
}
