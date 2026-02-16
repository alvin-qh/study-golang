package mutex_test

import (
	"runtime"
	"study/basic/testing/assertion"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	runtime.GOMAXPROCS(0)
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

	var n int32 = 0

	// 进行一次读锁定
	mut.RLock()

	// 在进行一次读锁定, 可以看到未发生阻塞或等待
	// 即多次读锁定会共享, 除第一次真实加锁外, 其余读锁定只会在同一把锁上加标记
	mut.RLock()

	// 操作一次共享资源
	n++

	// 解除一次读锁定
	mut.RUnlock()

	// 启动 goroutine, 加读锁后 100ms 后解除锁
	go func() {
		defer mut.RUnlock()

		// 休眠 100ms 后, 在函数结束后解除一次读锁定
		// 至此, 两次读锁都被解除
		time.Sleep(100 * time.Millisecond)

		// 操作一次共享资源
		n++
	}()

	start := time.Now()

	// 进行写锁定
	// 如果互斥锁已经加了读锁, 则必须等所有的读锁都解锁后才能加写锁
	mut.Lock()
	defer mut.Unlock()

	// 确认写锁必须在所有读锁解除后才能成功, 即 100ms 后成功
	assertion.DurationMatch(t, 100*time.Millisecond, time.Since(start))

	// 确认共享资源操作次数
	assert.Equal(t, int32(2), n)
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

		time.Sleep(10 * time.Millisecond)
	}()

	start := time.Now()

	// 通过 sync.RWMutex 实例加写锁
	mut.Lock()
	defer mut.Unlock()

	// 确认写锁必须在读锁解除后才能成功
	assertion.DurationMatch(t, 10*time.Millisecond, time.Since(start))
}
