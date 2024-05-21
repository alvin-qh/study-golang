package routine

import (
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	runtime.GOMAXPROCS(0)
}

// 测试互斥锁
//
// 互斥锁会产生一个临界区, 只有一个线程可以进入临界区, 其它的线程进入等待, 直到进入临界区的线程退出临界区
func TestMutex_LockAndUnlock(t *testing.T) {
	// 定义两个锁, 一个用于控制加法, 一个用于控制减法
	var mutAdd, mutSub sync.Mutex

	// 将控制减法的互斥锁提前锁定, 否则设定的时序就无法达成
	mutSub.Lock()

	n := 0

	var wg sync.WaitGroup
	wg.Add(2)

	// goroutine1, 用于对 n 加 1
	// 通过互斥锁, 可以保证 n 的值永远是从 0 加 1
	go func() {
		defer wg.Done()

		for i := 0; i < 100; i++ {
			// 锁定加法控制锁, 进入加法临界区
			mutAdd.Lock()
			// 确认进入加法临界区是, n 的值一定为 0
			assert.Equal(t, 0, n)

			n++
			// 解锁减法控制锁, 此时减法 goroutine 可以执行一次操作
			mutSub.Unlock()
		}
	}()

	// goroutine2, 用于对 n 减 1
	// 通过互斥锁, 可以保证 n 的值永远是从 1 减 1
	go func() {
		defer wg.Done()

		for i := 0; i < 100; i++ {
			// 对控制减法的互斥锁加锁
			mutSub.Lock()
			// 确保 n 此时的值为 1
			assert.Equal(t, 1, n)

			n--
			// 对控制
			mutAdd.Unlock()
		}
	}()

	wg.Wait()

	assert.Equal(t, 0, n)
}

// 测试互斥锁的非阻塞锁
//
// 通过 `Mutex.TryLock` 方法加锁时, 如果互斥锁已经被锁, 则不会阻塞等待, 直接返回 `false` 表示加锁失败
func TestMutex_TryLockAndUnlock(t *testing.T) {
	// 定义互斥锁对象
	var mut sync.Mutex

	// 尝试第一次加锁, 此时加锁成功
	r := mut.TryLock()
	assert.True(t, r)

	// 在已加锁的基础上, 尝试第二次加锁, 此时加锁失败
	r = mut.TryLock()
	assert.False(t, r)

	// 进行解锁
	mut.Unlock()

	// 再次尝试加锁, 此时加锁成功
	r = mut.TryLock()
	assert.True(t, r)
}

// 测试读写互斥锁
//
// 和 `Mutex` 类型相比, `RWMutex` 增加了 `RLock` 方法, 该方法表示读锁, 读锁在同一个 goroutine 内不会重复加锁
//
// 另外, 读锁
func TestMutex_RWMutex(t *testing.T) {
	var mut sync.RWMutex

	mut.RLock()

	mut.RLock()

	// 启动 goroutine, 加读锁后 100ms 后解除锁
	go func() {
		defer mut.RUnlock()

		mut.RUnlock()
		time.Sleep(100 * time.Millisecond)
	}()

	start := time.Now()

	mut.Lock()
	defer mut.Unlock()

	assert.GreaterOrEqual(t, time.Since(start).Milliseconds(), int64(100))
}
