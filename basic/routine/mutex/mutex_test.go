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
// 互斥锁 `sync.Mutex` 实例的 `Lock` 方法表示加锁, 同一时间一个 `sync.Mutex` 实例只能被加锁一次
//
// 加锁成功后的代码范围称为"临界区", 也就意味着只有一个 goroutine 可以进入临界区
//
// 进入临界区的代码必须执行 `sync.Mutex` 实例的 `Unlock` 方法进行解锁, 解锁也称为退出临界区, 直到解锁后吗, 另一个
// `Lock` 方法才能执行成功, 在此之前该 `Lock` 方法会被一直阻塞
//
// `sync.Mutex` 实例实现了 `sync.Locker` 接口, 即具备 `Lock` 和 `Unlock` 方法
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
// 通过 `sync.Mutex` 实例的 `TryLock` 方法加锁时, 如果互斥锁已经占用, 则不会阻塞等待, 直接返回 `false` 表示加锁失败
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
// 和 `sync.Mutex` 类型相比, `sync.RWMutex` 增加了 `RLock` 方法, 该方法表示加读锁
//
// 所谓的"读锁", 表示用于锁定资源进行读操作的锁, 在读锁解锁前, 资源应该是只读的
//
// 对同一个互斥锁多次加读锁不会造成阻塞, 也不会实质性的多次加锁, 只是在互斥锁上标记了多次读锁定, 当所有读锁都解锁后, 读锁解除
//
// 读锁不会和读锁互斥, 只会和写锁 (即通过 `sync.RWMutex` 实例的 `Lock` 方法加的锁) 互斥, 即:
//   - 如果加读锁前已经具备写锁, 则加锁会被阻塞, 直到写锁被解除;
//   - 如果加写锁前已经具备读锁或写锁, 则加锁会被阻塞, 直到所有的读锁或写锁都解除;
//
// `sync.RWMutex` 类型的读锁性能较高, 但写锁的性能不如 `sync.Mutex` 类型, 所以 `sync.RWMutex` 类型的应用场景为读多写少的情况,
// 如果写较多或读写差不多, 则应选择 `sync.Mutex` 锁
//
// `sync.RWMutex` 类型也实现了 `sync.Locker` 接口, 但对应的是 `Lock` 和 `Unlock` 方法, 接口不涉及 `RLock` 和 `RUnlock` 方法
func TestRWMutex_LockAndUnlock(t *testing.T) {
	// 声明读写互斥锁
	var mut sync.RWMutex

	// 进行一次读锁定
	mut.RLock()

	// 在进行一次读锁定, 可以看到未发生阻塞或等待
	// 即多次读锁定会共享, 除第一次真实加锁外, 其余读锁定只会在同一把锁上加标记
	mut.RLock()

	// 解除一次读锁定
	mut.RUnlock()

	// 启动 goroutine, 加读锁后 100ms 后解除锁
	go func() {
		defer mut.RUnlock()

		// 休眠 100ms 后, 在函数结束后解除一次读锁定
		// 至此, 两次读锁都被解除
		time.Sleep(100 * time.Millisecond)
	}()

	start := time.Now()

	// 进行写锁定
	// 如果互斥锁已经加了读锁, 则必须等所有的读锁都解锁后才能加写锁
	mut.Lock()
	defer mut.Unlock()

	assert.GreaterOrEqual(t, time.Since(start).Milliseconds(), int64(100))
}

// 测试读锁的非阻塞加锁方式
//
// 和 `sync.Mutex` 类型的 `TryLock` 方法类似, `sync.RWMutex` 类型也具备一个 `TryRLock` 方法, 可以以非阻塞的方式加读锁,
// 如果加锁失败 (即锁被占用), 则返回 `false`
func TestRWMutex_TryLockAndUnlock(t *testing.T) {
	// 声明读写互斥锁
	var mut sync.RWMutex

	// 尝试加读锁, 此次加锁成功
	r := mut.TryRLock()
	assert.True(t, r)

	// 尝试再次加读锁, 此次加锁也成功
	// 即可以多次加读锁
	r = mut.TryRLock()
	assert.True(t, r)

	// 尝试加写锁, 此次加锁失败, 已被读锁占用
	r = mut.TryLock()
	assert.False(t, r)

	// 进行两次读解锁, 则所有读锁都被解除
	mut.RUnlock()
	mut.RUnlock()

	// 再次尝试加写锁, 加锁成功
	r = mut.TryLock()
	assert.True(t, r)

	defer mut.Unlock()

	// 再次尝试加读锁, 加锁失败, 已被前一个写锁占用
	r = mut.TryRLock()
	assert.False(t, r)

	// 再次尝试加写锁, 加锁失败, 已被前一个写锁占用
	r = mut.TryLock()
	assert.False(t, r)
}

// 测试纯读锁
//
// 通过 `sync.RWMutex` 类型的 `Locker` 方法可以返回一个 `sync.Locker` 接口的实例, 接口的 `Lock` 和 `Unlock` 方法,
// 分别对应 `sync.RWMutex` 类型的 `RLock` 和 `RUnlock` 方法
func TestRWMutex_RLocker(t *testing.T) {
	// 声明读写互斥锁
	var mut sync.RWMutex

	// 获取读锁 sync.Locker 接口
	lk := mut.RLocker()

	// 通过 sync.Locker 接口实例加读锁
	lk.Lock()

	// 启动 goroutine, 100ms 后通过 sync.Locker 接口解除读锁
	go func() {
		defer lk.Unlock()

		time.Sleep(100 * time.Millisecond)
	}()

	start := time.Now()

	// 通过 sync.RWMutex 实例加写锁
	mut.Lock()
	defer mut.Unlock()

	// 确认写锁必须在读锁解除后才能成功
	assert.GreaterOrEqual(t, time.Since(start).Milliseconds(), int64(100))
}
