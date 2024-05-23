//go:build !windows

package filelock

import (
	"study/basic/testing/assertion"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	LOCK_FILE_NAME = ".lock"
)

// 测试互斥文件锁
//
// 本例中, 启动了三个异步任务, 包括:
//   - 其中两个异步任务用于等待文件锁, 并在等待成功后进行加锁, 执行完后退出锁
//   - 第三个异步任务用于等待前两个异步任务结束
func TestFileLock_LockUnlock(t *testing.T) {
	// 定义接收任务结果的通道
	ch := make(chan time.Duration, 100)
	defer close(ch)

	// 定义任务函数, 用于后续异步执行
	task := func(start time.Time) {
		// 由于任务中会对文件锁进行加锁操作, 所以所有任务均顺序执行, 无法并发

		// 获取文件锁
		fl := New(LOCK_FILE_NAME, false)

		// 进行互斥锁
		err := fl.XLock()
		assert.Nil(t, err)

		// 任务结束后自动解锁, 解除锁临界区
		defer fl.Unlock()

		// 进入锁临界区, 模拟计算时间 100ms
		time.Sleep(100 * time.Millisecond)

		// 发送结果, 结果为
		ch <- time.Since(start)
	}

	// 创建文件锁并进行锁定, 之后启动异步任务

	// 由于互斥文件锁的存在, 异步任务无法进入锁临界区, 在加锁位置等待
	// 此时解锁, 任务方可继续执行
	fl := New(LOCK_FILE_NAME, false)
	defer fl.Close()

	// 加锁
	err := fl.XLock()
	assert.Nil(t, err)

	// 启动两个异步任务模拟锁的使用
	go task(time.Now())
	go task(time.Now())

	// 解锁操作, 此时异步任务可以继续执行
	err = fl.Unlock()
	assert.Nil(t, err)

	// 记录异步任务结果的集合, 长度满足 2 时表示全部任务执行完毕
	r := make([]time.Duration, 0, 2)

	// 执行
exit:
	for len(r) < 2 {
		select {
		case duration := <-ch:
			r = append(r, duration)
		case <-time.After(time.Second):
			assert.Fail(t, "timeout")
			break exit
		}
	}

	assertion.Between(t, r[0].Milliseconds(), int64(100), int64(120))
	assertion.Between(t, r[1].Milliseconds(), int64(100+100), int64(100+100+20))
}
