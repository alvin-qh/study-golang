package routine

import (
	"context"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	runtime.GOMAXPROCS(0)
}

// 测试基本的 Channel 对象
// 线程 A 可以通过 channel 可以按顺序发送数据，线程 B 可以通过 channel 接收该数据
func TestSimpleChannel(t *testing.T) {
	// 创建一个 channel 对象，字符串类型，缓冲区 100 个元素
	ch := make(chan string, 100)
	defer close(ch)

	// 异步执行函数
	go func() {
		time.Sleep(time.Second) // 休眠一段时间
		ch <- "Hello"           // 向 channel 中写入字符串
	}()

	now := time.Now()

	// 等待从 channel 中读取数据
	s, ok := <-ch
	assert.True(t, ok)
	assert.Equal(t, "Hello", s)

	d := time.Since(now)
	assert.GreaterOrEqual(t, d, time.Second) // 1 秒后接收到数据，之前处于阻塞状态
}

// 测试无缓冲的 channel 对象
// 如果 channel 对象不具备缓冲，则会阻塞发送方，直到接收方读取了发送的数据
func TestNoCacheChannel(t *testing.T) {
	d := time.Duration(0)

	// 测试无缓冲 channel 发送阻塞情况

	wg := sync.WaitGroup{}
	wg.Add(1)

	ch := make(chan string) // 第二个参数为 0 或缺省，表示 channel 无缓冲

	go func() { // 异步函数，向 channel 中发送数据
		start := time.Now()
		defer func() {
			err, ok := recover().(error) // 获取 channel 被关闭的异常
			assert.True(t, ok)
			assert.Equal(t, "send on closed channel", err.Error())

			d = time.Since(start)
			wg.Done()
		}()
		ch <- "Hello" // 发送数据，此时由于没有任何接收方，所以发送会被阻塞
	}()

	<-time.After(time.Second) // 等待一段时间后关闭 channel，并不接收 channel 中的数据
	close(ch)

	wg.Wait()
	assert.GreaterOrEqual(t, d, time.Second) // 之前代码执行时间应该和等待时间相同

	// 测试无缓冲 channel 发送不阻塞情况

	wg = sync.WaitGroup{}
	wg.Add(1)

	ch = make(chan string)

	go func() { // 异步函数，向 chan 中发送数据
		defer wg.Done()

		start := time.Now()
		ch <- "Hello" // 发送数据，此时由于立即有接收方接收数据，所以发送不会阻塞

		d = time.Since(start)
	}()

	s := <-ch // 接收 channel 中的数据，此时发送方立即完成
	close(ch)
	assert.Equal(t, s, "Hello")

	wg.Wait()
	assert.Less(t, d, time.Millisecond) // 之前的代码执行时间很短暂
}

// 测试具备缓冲的 channel 对象
func TestCachedChannel(t *testing.T) {
	d := time.Duration(0)

	// 具备缓冲的 channel 对象
	// 在缓冲区没写满之前，数据发送不被阻塞

	wg := sync.WaitGroup{}
	wg.Add(1)

	ch := make(chan string, 1) // 第二个参数指定了缓冲区大小，即缓存多少个发送对象

	go func() { // 异步函数，向 channel 中发送数据
		start := time.Now()
		defer func() {
			d = time.Since(start)
			wg.Done()
		}()
		ch <- "Hello" // 发送数据，此时由于没有任何接收方，所以发送会被阻塞
	}()

	<-time.After(time.Second) // 等待一段时间后关闭 channel，并不接收 channel 中的数据
	close(ch)

	wg.Wait()
	assert.LessOrEqual(t, d, time.Millisecond) // 队列不阻塞，所以执行时间很短暂

	// 缓冲区写满后，数据发送被阻塞，此时 channel 的数据必须被消费掉，否则无法写入新数据

	wg = sync.WaitGroup{}
	wg.Add(1)

	ch = make(chan string, 1) // 第二个参数指定了缓冲区大小，即缓存多少个发送对象

	go func() { // 异步函数，向 channel 中发送数据
		start := time.Now()
		defer func() {
			err, ok := recover().(error) // 获取 channel 被关闭的异常
			assert.True(t, ok)
			assert.Equal(t, "send on closed channel", err.Error())

			d = time.Since(start)
			wg.Done()
		}()
		ch <- "Hello"
		ch <- "World" // 发送数据，此时 channel 缓冲已满，发送被阻塞
	}()

	<-time.After(time.Second) // 等待一段时间后关闭 channel，并不接收 channel 中的数据
	close(ch)

	wg.Wait()
	assert.GreaterOrEqual(t, d, time.Second) // 之前代码执行时间应该和等待时间相同

	// 缓冲区写满后被消费，发送不阻塞

	wg = sync.WaitGroup{}
	wg.Add(1)

	ch = make(chan string, 1) // 第二个参数指定了缓冲区大小，即缓存多少个发送对象

	go func() { // 异步函数，向 channel 中发送数据
		start := time.Now()
		defer func() {
			recover()

			d = time.Since(start)
			wg.Done()
		}()
		ch <- "Hello"
		ch <- "World" // 发送数据，此时由于没有任何接收方，所以发送会被阻塞
	}()

	s := <-ch // 接收 channel 中的数据，此时发送方立即完成
	assert.Equal(t, "Hello", s)
	close(ch)

	wg.Wait()
	assert.LessOrEqual(t, d, time.Millisecond) // 代码执行时间短暂
}

// 测试 Lock
func TestChanLock(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(3)

	ctx := context.TODO() // 初始化空 context 对象

	l := NewLock() // 产生锁对象

	start := time.Now()

	for i := 0; i < 3; i++ {
		// 启动协程，由于所得关系，所以三个协程按顺序执行
		// 总执行时间是三个协程执行时间之和
		go func() {
			defer wg.Done()

			locked := l.Lock(ctx)  // 锁定
			assert.True(t, locked) // 判断是否锁定成功
			if locked {
				defer l.Unlock()        // 结束后解锁
				time.Sleep(time.Second) // 等待一段时间
			}
		}()
	}

	wg.Wait()

	d := time.Since(start)
	assert.GreaterOrEqual(t, d, time.Second*3) // 总体执行时间超过 3 秒
}

// 测试锁超时
func TestChanLockTimeout(t *testing.T) {
	l := NewLock() // 产生锁对象

	locked := l.Lock(context.TODO()) // 锁定
	defer l.Unlock()                 // 函数退出时解锁
	assert.True(t, locked)           // 锁定成功

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2) // 创建 2 秒超时的 context 对象
	defer cancel()

	start := time.Now()

	locked = l.Lock(ctx)    // 再次锁定，因为没有解锁，所以无法进入锁
	assert.False(t, locked) // 锁定失败

	d := time.Since(start)
	assert.GreaterOrEqual(t, d, time.Second*2) // 执行时间超过 2 秒
}

// 测试加锁解锁
func TestChanLockAndUnlock(t *testing.T) {
	l := NewLock() // 产生锁对象
	defer l.Unlock()

	ctx := context.TODO()

	n := 0

	start := time.Now()

	wg := sync.WaitGroup{}
	wg.Add(200)

	for i := 0; i < 100; i++ {
		// 读协程
		go func() {
			defer wg.Done()

			l.Lock(ctx)
			defer l.Unlock()

			n--
			time.Sleep(time.Millisecond * 10)
		}()

		// 写协程
		go func() {
			defer wg.Done()

			l.Lock(ctx)
			defer l.Unlock()

			n++
			time.Sleep(time.Millisecond * 10)
		}()
	}

	wg.Wait()

	d := time.Since(start)
	assert.GreaterOrEqual(t, d, time.Second*2)

	assert.Equal(t, 0, n)
}

// 测试通过 channel 创建的生成器
func TestChanGenerator(t *testing.T) {
	// 创建生成器对象
	g := NewGenerator(func(ch chan interface{}) interface{} { // 传入生成函数作为参数
		n := 0
		for {
			ch <- n
			n++
		}
	})
	defer g.Close()

	x := 0

	// 通过生成器生成 100 个数据
	// 通过 Next 函数，返回生成器生成出的数据
	for n, err := g.Next(); err == nil && n.(int) < 100; n, err = g.Next() { // 获取生成器下一个数据
		assert.NoError(t, err)
		assert.Equal(t, x, n.(int)) // 检查生成数据
		x++
	}
	assert.Equal(t, 100, x)
	x++

	// 继续生成 100 个数据
	// 通过 Range 函数，返回一个 channel，通过对其进行 range 操作进行生成
	for n := range g.Range() {
		assert.Equal(t, x, n.(int)) // 检查生成数据
		if x > 100 {
			break
		}
		x++
	}
}
