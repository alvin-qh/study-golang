package filelock

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const LOCK_FILE = "./.lock"

// 测试互斥文件锁
func TestFileXLockAndUnlock(t *testing.T) {
	defer os.Remove(LOCK_FILE)

	// 记录 task 是否完成的 map
	tasks := map[string]bool{
		"task1": false,
		"task2": false,
	}

	// 等待 task 全部完成的等待组
	wg := &sync.WaitGroup{}
	wg.Add(2) // 添加两个任务

	// 异步 task 函数
	// 由于任务中会对文件锁进行加锁操作，所以所有任务均顺序执行，无法并发
	task := func(name string) {
		fl := New(LOCK_FILE, false) // 创建文件锁
		defer fl.Unlock()           // 任务结束后自动解锁，解除锁临界区

		err := fl.XLock() // 进行互斥锁
		assert.NoError(t, err)

		// 进入锁临界区

		tasks[name] = true // 标记当前任务完成
		wg.Done()          // 通知等待组
	}

	// 创建文件锁并进行锁定，之后启动异步任务
	// 由于互斥文件锁的存在，异步任务无法进入锁临界区，在加锁位置等待
	// 此时解锁，任务方可继续执行
	fl := New(LOCK_FILE, false) // 创建文件锁

	err := fl.XLock() // 加锁
	assert.NoError(t, err)

	go task("task1") // 启动异步任务
	go task("task2")

	ok := WaitTimeout(wg, time.Second*2) // 先尝试对等待组等待 2 秒，无法成功，表示任务都无法结束
	assert.False(t, ok)

	err = fl.Unlock() // 解锁操作，此时异步任务可以继续执行
	assert.NoError(t, err)

	wg.Wait() // 继续对等待组进行等待

	assert.True(t, tasks["task1"]) // 等待成功，判断任务是否正确执行
	assert.True(t, tasks["task2"])
}
