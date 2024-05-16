package routine

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
// 互斥锁会产生一个临界区, 只有一个线程可以进入临界区, 其它的线程进入等待, 直到进入临界区的线程退出临界区
func TestMutex(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(200)

	var mut sync.Mutex // 定义互斥锁对象

	num := 0

	// 协程函数, 用于增加公共变量 (表示写)
	increment := func() {
		defer wg.Done()

		mut.Lock()         // 加锁, 进入临界区
		defer mut.Unlock() // 解锁, 退出临界区

		num += 1 // 操作公共变量
	}

	// 协程函数, 用于减少公共变量 (表示读)
	read := func() {
		defer wg.Done()

		mut.Lock()
		defer mut.Unlock() // 解锁, 退出临界区

		num -= 1 // 操作公共变量
	}

	// 执行 100 次写操作和读操作
	// 同步执行可以保证每个写操作都有对应的读操作
	for i := 0; i < 100; i++ {
		go increment()
		go read()
	}

	wg.Wait()

	// 读写次数平衡, 所以结果为 0
	assert.Equal(t, 0, num)
}

// 测试读写互斥锁
// RWMutex 具备 RLock/RUnlock 和 Lock/Unlock 函数, 额外的 RLock/RUnlock 用于读锁, 在读多于写的操作中, 可以提高执行效率
func TestRWMutex(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(200)

	var mut sync.RWMutex // 定义互斥锁对象

	num := 0

	// 协程函数, 用于增加公共变量的值 (表示写)
	increment := func() {
		defer wg.Done()

		mut.Lock() // 加 X 锁, 进入临界区
		defer mut.Unlock()

		num += 1
	}

	// 协程函数, 用于减少公共变量的值 (表示读)
	read := func() {
		defer wg.Done()

		mut.RLock() // 加 S 锁, 读锁只对写锁做阻塞处理, 对同为读锁不做处理, 所以在读多写少的环境下, 读锁的并发性更好
		defer mut.RUnlock()

		num -= 1
	}

	// 执行 100 次写操作和读操作
	// 同步执行可以保证每个写操作都有对应的读操作
	for i := 0; i < 100; i++ {
		go increment()
		go read()
	}

	wg.Wait()

	// 读写次数平衡, 所以结果为 0
	assert.Equal(t, 0, num)
}
