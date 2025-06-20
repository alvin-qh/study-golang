package pool

import (
	"basic/builtin/slice/utils"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试通过任务池执行大量并发任务
func TestTaskPool_Execute(t *testing.T) {
	// 创建任务池, 共 10 个并发任务
	pool := NewTaskPool[string, int](10)

	// 执行完毕后, 关闭任务池
	defer pool.Close()

	// 存储执行结果的切片
	rs := make([]int, 0)
	var mux sync.Mutex

	// 为 100 个任务创建等待组
	var wg sync.WaitGroup
	wg.Add(100)

	// 启动 goroutine, 分配任务
	go func() {
		// 创建可执行任务
		exec := pool.Worker(func(arg string) (int, error) {
			return strconv.Atoi(arg)
		})

		// 共执行 100 个任务
		for i := 0; i < 100; i++ {
			// 执行任务
			exec(
				strconv.Itoa(i+1),
				func(result int) {
					// 任务执行成功, 记录结果
					mux.Lock()
					defer mux.Unlock()

					rs = append(rs, result)

					// 任务完成
					wg.Done()
				},
				func(err error) {
					// 任务执行失败, 记录错误信息
					assert.NoError(t, err)

					// 任务完成
					wg.Done()
				},
			)
		}
	}()

	// 等待所有任务结束
	wg.Wait()

	// 查看任务结果
	assert.Len(t, rs, 100)
	assert.ElementsMatch(t, utils.Range(1, 101, 1), rs)
}

// 测试任务池的优雅关闭
//
// 任务池关闭后, 可等待所有 goroutine 结束
func TestTaskPool_CloseAndWait(t *testing.T) {
	// 创建任务池, 共 10 个并发任务
	pool := NewTaskPool[string, int](10)

	// 关闭任务池并等待当前任务结束
	defer pool.CloseAndWait()

	// 启动 goroutine, 分配任务
	go func() {
		// 创建可执行任务
		exec := pool.Worker(
			func(arg string) (int, error) {
				// 阻塞任务
				time.Sleep(1 * time.Second)
				return strconv.Atoi(arg)
			},
		)

		// 共执行 100 个任务
		for i := 0; i < 100; i++ {
			// 执行任务
			exec(
				strconv.Itoa(i+1),
				func(result int) {
					assert.Fail(t, "cannot run here")
				},
				func(err error) {
					assert.Fail(t, "cannot run here")
				},
			)
		}
	}()
}
