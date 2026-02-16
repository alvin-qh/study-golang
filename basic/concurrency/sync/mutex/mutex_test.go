package mutex_test

import (
	"runtime"
	"sync"
	"testing"

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

	// 创建等待组对象, 用于等待 2 个任务完成
	var wg sync.WaitGroup
	wg.Add(2)

	// 创建一个 goroutine, 用于对 n 加 1
	// 通过互斥锁, 可以保证 n 的值永远是从 0 加 1
	go func() {
		defer wg.Done()

		for range 100 {
			// 对加法控制锁加锁, 进入加法临界区
			mutAdd.Lock()

			// 确认程序进入加法临界区时, n 的值一定为 0
			assert.Equal(t, 0, n)

			n++

			// 解锁减法控制锁, 此时减法 goroutine 可以执行一次操作
			mutSub.Unlock()
		}
	}()

	// 再创建一个 goroutine, 用于对 n 减 1
	// 通过互斥锁, 可以保证 n 的值永远是从 1 减 1
	go func() {
		defer wg.Done()

		for range 100 {
			// 对减法控制锁加锁, 进入减法临界区
			mutSub.Lock()

			// 确保 n 此时的值为 1
			assert.Equal(t, 1, n)

			n--

			// 解锁加法控制锁, 此时加法 goroutine 可以执行一次操作
			mutAdd.Unlock()
		}
	}()

	// 等待所有任务完成
	wg.Wait()

	// 确认程序退出时, n 的值为 0, 说明互斥锁确实在控制任务的执行顺序
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
