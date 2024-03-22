//go:build !windows

package filelock

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	LOCK_FILE = "./.lock"
)

// 测试互斥文件锁
func TestFileXLockAndUnlock(t *testing.T) {
	// 记录 task 是否完成的 map
	tasks := map[string]bool{
		"task1": false,
		"task2": false,
	}

	// 等待 task 全部完成的等待组
	wg := &sync.WaitGroup{}

	// 添加两个等待任务
	wg.Add(2)

	// 定义任务函数, 用于后续异步执行
	task := func(name string) {
		// 由于任务中会对文件锁进行加锁操作, 所以所有任务均顺序执行, 无法并发

		// 获取文件锁
		fl := New(LOCK_FILE, false)

		// 进行互斥锁
		err := fl.XLock()
		assert.NoError(t, err)

		// 任务结束后自动解锁, 解除锁临界区
		defer fl.Unlock()

		// 进入锁临界区

		// 标记当前任务完成
		tasks[name] = true

		// 通知等待组
		wg.Done()
	}

	// 创建文件锁并进行锁定, 之后启动异步任务
	// 由于互斥文件锁的存在, 异步任务无法进入锁临界区, 在加锁位置等待
	// 此时解锁, 任务方可继续执行
	fl := New(LOCK_FILE, false)
	defer os.Remove(LOCK_FILE)

	err := fl.XLock() // 加锁
	assert.NoError(t, err)

	go task("task1") // 启动异步任务
	go task("task2")

	// 先尝试对等待组等待 100ms, 无法成功, 表示任务都无法结束
	ok := WaitTimeout(wg, time.Millisecond*100)
	assert.False(t, ok)

    // 解锁操作, 此时异步任务可以继续执行
	err = fl.Unlock()
	assert.NoError(t, err)

	wg.Wait() // 继续对等待组进行等待

    // 等待成功, 判断任务是否正确执行
	assert.True(t, tasks["task1"])
	assert.True(t, tasks["task2"])
}
