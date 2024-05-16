package channellock

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试 Lock
func TestChanLock(t *testing.T) {
	// 创建等待组, 可等待 3 个任务
	wg := sync.WaitGroup{}
	wg.Add(3)

	// 初始化空 context 对象
	ctx := context.TODO()

	// 产生锁对象
	lock := New()

	start := time.Now()

	for i := 0; i < 3; i++ {
		// 启动协程, 由于所得关系, 所以三个协程按顺序执行
		// 总执行时间是三个协程执行时间之和
		go func() {
			defer wg.Done()

			// 锁定
			locked := lock.Lock(ctx)
			// 判断是否锁定成功
			assert.True(t, locked)

			if locked {
				// 结束后解锁
				defer lock.Unlock()

				// 等待一段时间
				time.Sleep(time.Millisecond * 100)
			}
		}()
	}

	wg.Wait()

	d := time.Since(start)
	// 总体执行时间超过 300ms
	assert.GreaterOrEqual(t, d, time.Millisecond*300)
}

// 测试锁超时
func TestChanLockTimeout(t *testing.T) {
	// 产生锁对象
	lock := New()

	// 锁定
	locked := lock.Lock(context.TODO())
	// 锁定成功
	assert.True(t, locked)

	// 函数退出时解锁
	defer lock.Unlock()

	// 创建 2 秒超时的 context 对象
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*200)
	defer cancel()

	start := time.Now()

	// 再次锁定, 因为没有解锁, 所以无法进入锁
	locked = lock.Lock(ctx)
	assert.False(t, locked) // 锁定失败

	d := time.Since(start)
	// 执行时间超过 2 秒
	assert.GreaterOrEqual(t, d, time.Millisecond*200)
}

// 测试加锁解锁
func TestChanLockAndUnlock(t *testing.T) {
	// 产生锁对象
	lock := New()
	defer lock.Unlock()

	ctx := context.TODO()

	n := 0

	start := time.Now()

	wg := sync.WaitGroup{}
	wg.Add(10)

	for i := 0; i < 5; i++ {
		// 读协程
		go func() {
			defer wg.Done()

			lock.Lock(ctx)
			defer lock.Unlock()

			n--
			time.Sleep(time.Millisecond * 10)
		}()

		// 写协程
		go func() {
			defer wg.Done()

			lock.Lock(ctx)
			defer lock.Unlock()

			n++
			time.Sleep(time.Millisecond * 10)
		}()
	}

	wg.Wait()

	d := time.Since(start)
	assert.GreaterOrEqual(t, d, time.Millisecond*100)

	assert.Equal(t, 0, n)
}
