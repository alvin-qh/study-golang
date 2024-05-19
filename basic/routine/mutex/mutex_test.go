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
// 互斥锁会产生一个临界区, 只有一个线程可以进入临界区, 其它的线程进入等待,
// 直到进入临界区的线程退出临界区
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

// 测试互斥锁
//
// 互斥锁会产生一个临界区, 只有一个线程可以进入临界区, 其它的线程进入等待,
// 直到进入临界区的线程退出临界区
func TestMutex_TryLockAndUnlock(t *testing.T) {
	// 定义互斥锁对象
	var mut sync.Mutex

	// 加锁, 进入临界区
	mut.Lock()

	// 启动一个 goroutine, 在 100ms 以后进行解锁
	go func() {
		time.Sleep(100 * time.Millisecond)
		mut.Unlock()
	}()

	start := time.Now()

	// 再次锁定, 此时会发生阻塞, 100ms 后解锁
	mut.Lock()
	defer mut.Unlock()

	// 确定再次锁定会超过 100ms, 因为 100ms 后才进行解锁
	assert.GreaterOrEqual(t, time.Since(start).Milliseconds(), int64(100))
}

// 测试读写互斥锁
//
// `RWMutex` 具备 `RLock`/`RUnlock` 和 `Lock`/`Unlock` 函数,
// 额外的 `RLock`/`RUnlock` 用于读锁, 在读多于写的操作中, 可以提高执行效率
func TestMutex_RWMutex(t *testing.T) {
	var wg sync.WaitGroup

	num := 0

	// 定义互斥锁对象
	var mut sync.RWMutex

	// 协程函数, 用于增加公共变量的值 (表示写)
	increment := func() {
		defer wg.Done()

		// 加 X 锁, 进入临界区
		mut.Lock()
		defer mut.Unlock()

		num += 1
	}

	// 协程函数, 用于减少公共变量的值 (表示读)
	read := func() {
		defer wg.Done()

		// 加 S 锁, 读锁只对写锁做阻塞处理, 对同为读锁不做处理, 所以在读多写少的环境下, 读锁的并发性更好
		mut.RLock()
		defer mut.RUnlock()

		num -= 1
	}

	// 执行 100 次写操作和读操作
	// 同步执行可以保证每个写操作都有对应的读操作
	for i := 0; i < 100; i++ {
		wg.Add(2)

		go increment()
		go read()
	}

	wg.Wait()

	// 读写次数平衡, 所以结果为 0
	assert.Equal(t, 0, num)
}
