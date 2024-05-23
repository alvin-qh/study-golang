package cond

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

// 测试条件量
//
// 条件量 `sync.Cond` 表示同步原语中的 `condition`
//
// 和 `sync.Mutex` 类型不同, `sync.Cond` 类型允许通过一个实例, 在一个 goroutine 中控制另一个 goroutine,
// 而前者只能做到所有相关的 goroutine 轮流进入临界区
//
// 通过 `sync.Cond` 实例的 `Signal` 方法发送一个信号, 此时所有在该实例等待的 (执行 `Wait` 方法并阻塞) 的 goroutine
// 中的一个, 会立即结束等待
//
// 注意: `sync.Cond` 实例的 `Wait` 方法必须在一个相关临界区内执行, 所以在创建 `sync.Cond` 实例时, 需要为其关联一个 `sync.Mutex`
// 实例, 并在 `sync.Cond` 实例执行 `Wait` 方法前, 通过该实例的 `L` 字段进入临界区
func TestCond_NewCond(t *testing.T) {
	// 创建互斥锁, 创建 `syncCond` 实例需要借助互斥锁
	var mut sync.Mutex

	// 创建 `sync.Cond` 实例, 关联之前创建的 `sync.Mutex` 实例
	// 这里关联的实际上是一个 `sync.Locker` 接口实例, 所以 `sync.RWMutex` 等类型实例也可以使用
	cond := sync.NewCond(&mut)

	// 声明一个接收结果的通道
	ch := make(chan int64)

	// 启动 goroutine, 等待 `sync.Cond` 实例的信号, 并计算等待时长
	go func() {
		// 进入临界区
		cond.L.Lock()
		defer cond.L.Unlock()

		start := time.Now()

		// 等待信号发送
		cond.Wait()

		// 计算信号等待时长
		ch <- time.Since(start).Milliseconds()
	}()

	// 休眠 100ms 后, 向 `sync.Cond` 实例发送信号
	// 此时 goroutine 中的等待会结束
	time.Sleep(100 * time.Millisecond)
	// 发送信号
	cond.Signal()

	// 接收结果, 即 goroutine 等待时长
	since := <-ch
	assertion.Between(t, since, int64(100), int64(120))
}

// 测试条件量信号
//
// 通过 `sync.Cond` 实例, 每个 goroutine 可以调用 `Wait` 方法等待一个信号, 而调用 `Signal` 方法可以发送一次信号,
// 即结束一个 goroutine 的等待, 其余 goroutine 仍处于继续等待的状态
func TestCond_SignalOneByOne(t *testing.T) {
	var mut sync.Mutex

	// 创建条件量
	cond := sync.NewCond(&mut)

	// 用于输出结果的信道实例
	ch := make(chan []int64)
	defer close(ch)

	start := time.Now()

	// 启动 10 个 goroutine
	for i := 0; i < 10; i++ {
		// 每个 goroutine 中, 等待信号, 并记录信号等待时间
		go func(id int64) {
			// 进入临界区
			cond.L.Lock()
			defer cond.L.Unlock()

			// 等待信号发送
			cond.Wait()

			// 计算信号等待时长
			ch <- []int64{id, time.Since(start).Milliseconds()}
		}(int64(i + 1))
	}

	// 启动 goroutine, 分别通过条件量发送 10 次信号
	// 每次信号唤醒一个等待, 则 10 次信号所有的 goroutine 都将被唤醒
	go func() {
		for i := 0; i < 10; i++ {
			// 休眠 10 毫秒后发送信号
			time.Sleep(10 * time.Millisecond)
			cond.Signal()
		}
	}()

	var totalId, totalTime = int64(0), int64(0)

	count := 0
	// 通过信道获取所有 goroutine 的执行结果
	for r := range ch {
		// 计算每个 goroutine 函数 id 参数的总和
		totalId += r[0]

		// 记录每个 goroutine 等待时长的总和
		totalTime += r[1]

		count++
		// 接收 10 次结果后退出循环
		if count == 10 {
			break
		}
	}

	// 确认每个 goroutine 都依次完成等待， 并发送结果
	// 各协程 ID 和 1+2+3+...+10
	assert.Equal(t, int64(55), totalId)
	// 各协程等待时长和 10+20+30+...+100
	assertion.Between(t, totalTime, int64(550), int64(580))
}

// 测试条件量信号广播
//
// 通过 `sync.Cond` 类型的 `Broadcast` 方法, 可以将信号一次性发送给所有在该条件量等待的 goroutine,
// 而不是类似 `Signal` 方法那样一次只发送一个信号
func TestCond_Broadcast(t *testing.T) {
	var mut sync.Mutex

	// 创建条件量
	cond := sync.NewCond(&mut)

	// 用于输出结果的信道实例
	ch := make(chan []int64)
	defer close(ch)

	start := time.Now()

	// 启动 10 个 goroutine
	for i := 0; i < 10; i++ {
		// 每个 goroutine 中, 等待信号, 并记录信号等待时间
		go func(id int64) {
			// 进入临界区
			cond.L.Lock()
			defer cond.L.Unlock()

			// 等待信号发送
			cond.Wait()

			// 计算信号等待时长
			ch <- []int64{id, time.Since(start).Milliseconds()}
		}(int64(i))
	}

	// 启动 goroutine, 等待 10ms 后广播信号
	go func() {
		// 休眠 10 毫秒后广播信号
		time.Sleep(10 * time.Millisecond)
		cond.Broadcast()
	}()

	var totalId, totalTime = int64(0), int64(0)

	count := 0
	// 通过信道获取所有 goroutine 的执行结果
	for r := range ch {
		// 计算每个 goroutine 函数 id 参数的总和
		totalId += r[0]

		// 记录每个 goroutine 等待时长的总和
		totalTime += r[1]

		count++
		// 接收 10 次结果后退出循环
		if count == 10 {
			break
		}
	}

	// 确认每个 goroutine 都依次完成等待， 并发送结果
	// 各协程 ID 和 1+2+3+...+10
	assert.Equal(t, int64(45), totalId)
	// 各协程等待时长和 10+10+10+...+10
	assertion.Between(t, totalTime, int64(100), int64(120))
}
