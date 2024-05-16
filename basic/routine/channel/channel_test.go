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

// 测试基本的 Channel 对象
//
// routine A 可以通过 channel 可以按顺序发送数据, routine B 可以通过 channel 接收该数据
func TestSimpleChannel(t *testing.T) {
	// 创建一个 channel 对象, 字符串类型, 缓冲区 100 个元素
	ch := make(chan string, 100)
	defer close(ch)

	// 异步执行函数
	go func() {
		// 休眠一段时间
		time.Sleep(time.Millisecond * 100)

		// 向 channel 中写入字符串
		ch <- "Hello"
	}()

	now := time.Now()

	// 等待从 channel 中读取数据
	s, ok := <-ch
	d := time.Since(now)

	assert.True(t, ok)
	assert.Equal(t, "Hello", s)

	// 100ms 后接收到数据, 之前处于阻塞状态
	assert.GreaterOrEqual(t, d, time.Millisecond*100)
}

// 测试无缓冲的 channel 对象
//
// 如果 channel 对象不具备缓冲, 则会阻塞发送方, 直到接收方读取了发送的数据
func TestNoCacheChannelBlocked(t *testing.T) {
	d := time.Duration(0)

	wg := sync.WaitGroup{}
	wg.Add(1)

	// 第二个参数为 0 或缺省, 表示 channel 无缓冲
	ch := make(chan string)

	// 异步函数, 向 channel 中发送数据
	go func() {
		start := time.Now()
		defer func() {
			// 获取 channel 被关闭的异常
			err, ok := recover().(error)
			assert.True(t, ok)
			assert.Equal(t, "send on closed channel", err.Error())

			d = time.Since(start)
			wg.Done()
		}()

		// 发送数据, 此时由于没有任何接收方, 所以发送会被阻塞
		ch <- "Hello"
	}()

	<-time.After(time.Millisecond * 100)
	// 等待一段时间后关闭 channel, 并不接收 channel 中的数据
	close(ch)

	wg.Wait()

	// 之前代码执行时间应该和等待时间相同
	assert.GreaterOrEqual(t, d, time.Millisecond*100)
}

// 测试无缓冲 channel 发送不阻塞情况
func TestNoCacheChannelNoBlocking(t *testing.T) {
	d := time.Duration(0)

	wg := sync.WaitGroup{}
	wg.Add(1)

	ch := make(chan string)

	// 异步函数, 向 chan 中发送数据
	go func() {
		defer wg.Done()

		start := time.Now()

		// 发送数据, 此时由于立即有接收方接收数据, 所以发送不会阻塞
		ch <- "Hello"

		d = time.Since(start)
	}()

	// 接收 channel 中的数据, 此时发送方立即完成
	s := <-ch
	close(ch)
	assert.Equal(t, s, "Hello")

	wg.Wait()
	// 之前的代码执行时间很短暂
	assert.Less(t, d, time.Millisecond)
}

// 测试具备缓冲的 channel 对象
func TestCachedChannelNoBlocking(t *testing.T) {
	d := time.Duration(0)

	// 具备缓冲的 channel 对象
	// 在缓冲区没写满之前, 数据发送不被阻塞

	wg := sync.WaitGroup{}
	wg.Add(1)

	// 第二个参数指定了缓冲区大小, 即缓存多少个发送对象
	ch := make(chan string, 1)

	// 异步函数, 向 channel 中发送数据
	go func() {
		start := time.Now()
		defer func() {
			d = time.Since(start)
			wg.Done()
		}()

		// 发送数据, 此时由于没有任何接收方, 所以发送会被阻塞
		ch <- "Hello"
	}()

	// 等待一段时间后关闭 channel, 并不接收 channel 中的数据
	<-time.After(time.Millisecond * 100)
	close(ch)

	wg.Wait()

	// 队列不阻塞, 所以执行时间很短暂
	assert.LessOrEqual(t, d, time.Millisecond)
}

func TestCachedChannelBlocked(t *testing.T) {
	d := time.Duration(0)
	// 缓冲区写满后, 数据发送被阻塞, 此时 channel 的数据必须被消费掉, 否则无法写入新数据

	wg := sync.WaitGroup{}
	wg.Add(1)

	// 第二个参数指定了缓冲区大小, 即缓存多少个发送对象
	ch := make(chan string, 1)

	// 异步函数, 向 channel 中发送数据
	go func() {
		start := time.Now()
		defer func() {
			// 获取 channel 被关闭的异常
			err, ok := recover().(error)
			assert.True(t, ok)
			assert.Equal(t, "send on closed channel", err.Error())

			d = time.Since(start)
			wg.Done()
		}()

		// 发送数据, 此时 channel 缓冲已满, 发送被阻塞
		ch <- "Hello"
		ch <- "World"
	}()

	// 等待一段时间后关闭 channel, 并不接收 channel 中的数据
	<-time.After(time.Millisecond * 100)
	close(ch)

	wg.Wait()

	// 之前代码执行时间应该和等待时间相同
	assert.GreaterOrEqual(t, d, time.Millisecond*100)

	// 缓冲区写满后被消费, 发送不阻塞

	wg = sync.WaitGroup{}
	wg.Add(1)

	// 第二个参数指定了缓冲区大小, 即缓存多少个发送对象
	ch = make(chan string, 1)

	// 异步函数, 向 channel 中发送数据
	go func() {
		start := time.Now()
		defer func() {
			recover()

			d = time.Since(start)
			wg.Done()
		}()

		// 发送数据, 此时由于没有任何接收方, 所以发送会被阻塞
		ch <- "Hello"
		ch <- "World"
	}()

	// 接收 channel 中的数据, 此时发送方立即完成
	s := <-ch
	assert.Equal(t, "Hello", s)
	close(ch)

	wg.Wait()
	// 代码执行时间短暂
	assert.LessOrEqual(t, d, time.Millisecond)
}
