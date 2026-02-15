package channel_test

import (
	"runtime"
	"strconv"
	"study/basic/testing/assertion"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func init() {
	runtime.GOMAXPROCS(0)
}

// 测试基本的 `chan` 实例
//
// routine A 可以通过 `chan` 可以按顺序发送数据, routine B 可以通过 `chan` 接收该数据
func TestChan_Simple(t *testing.T) {
	// 创建一个 chan 实例, 字符串类型, 缓冲区 100 个元素
	ch := make(chan string, 100)
	defer close(ch)

	// 异步执行函数
	go func() {
		// 休眠一段时间
		time.Sleep(10 * time.Millisecond)

		// 向 chan 中写入字符串
		ch <- "Hello"
	}()

	now := time.Now()

	// 等待从 chan 中读取数据
	s, ok := <-ch
	d := time.Since(now)

	assert.True(t, ok)
	assert.Equal(t, "Hello", s)

	// 100ms 后接收到数据, 之前处于阻塞状态
	assertion.DurationMatch(t, 10*time.Millisecond, d)
}

// 测试无缓冲的 `chan` 实例
//
// 如果 `chan` 实例不具备缓冲, 则会阻塞发送方, 直到接收方读取了发送的数据
func TestChan_Blocked(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	// 第二个参数为 0 或缺省, 表示 chan 无缓冲
	ch := make(chan string)

	var d time.Duration

	// 异步函数, 向 chan 中发送数据
	go func() {
		start := time.Now()

		defer func() {
			// 获取 chan 被关闭的异常
			err, ok := recover().(error)
			assert.True(t, ok)
			assert.Equal(t, "send on closed channel", err.Error())

			d = time.Since(start)
			wg.Done()
		}()

		// 发送数据, 此时由于没有任何接收方, 所以发送会被阻塞
		ch <- "Hello"
	}()

	<-time.After(time.Millisecond * 10)
	// 等待一段时间后关闭 chan, 并不接收 chan 中的数据
	close(ch)

	wg.Wait()

	// 之前代码执行时间应该和等待时间相同
	assertion.DurationMatch(t, 10*time.Millisecond, d)
}

// 测试无缓冲 chan 发送不阻塞情况
//
// 如果 chan 本身不具备缓冲, 则需要保证接收方可以及时从 chan 中读取
func TestChan_NonBlocked(t *testing.T) {
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

	// 接收 chan 中的数据, 此时发送方立即完成
	s := <-ch
	assert.Equal(t, s, "Hello")

	close(ch)

	wg.Wait()
	// 之前的代码执行时间很短暂
	assert.Less(t, d, time.Millisecond)
}

// 测试具备缓冲的 chan 实例
//
// 具备缓冲区的 chan 可以在读取不及时的情况下, 仍有一定的缓冲保证发送不阻塞
func TestChan_CheckedNonBlocked(t *testing.T) {
	d := time.Duration(0)

	// 具备缓冲的 chan 实例
	// 在缓冲区没写满之前, 数据发送不被阻塞

	wg := sync.WaitGroup{}
	wg.Add(1)

	// 第二个参数指定了缓冲区大小, 即缓存多少个发送实例
	ch := make(chan string, 1)

	// 异步函数, 向 chan 中发送数据
	go func() {
		start := time.Now()
		defer func() {
			d = time.Since(start)
			wg.Done()
		}()

		// 发送数据, 此时由于没有任何接收方, 所以发送会被阻塞
		ch <- "Hello"
	}()

	// 等待一段时间后关闭 chan, 并不接收 chan 中的数据
	<-time.After(time.Millisecond * 100)
	close(ch)

	wg.Wait()

	// 队列不阻塞, 所以执行时间很短暂
	assert.LessOrEqual(t, d, time.Millisecond)
}

// 测试具备缓冲的 chan 在缓冲写满后仍会发生阻塞
func TestChan_CheckedBlocked(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(1)

	// 第二个参数指定了缓冲区大小, 即缓存多少个发送实例
	// 缓冲区写满后, 数据发送被阻塞, 此时 chan 的数据必须被消费掉, 否则无法写入新数据
	ch := make(chan string, 1)

	var d time.Duration

	// 异步函数, 向 chan 中发送数据
	go func() {
		start := time.Now()
		defer func() {
			// 获取 chan 被关闭的异常
			err, ok := recover().(error)
			assert.True(t, ok)
			assert.Equal(t, "send on closed channel", err.Error())

			d = time.Since(start)
			wg.Done()
		}()

		// 发送数据, 此时 chan 缓冲已满, 发送被阻塞
		ch <- "Hello"
		ch <- "World"
	}()

	// 等待一段时间后关闭 chan, 并不接收 chan 中的数据
	<-time.After(time.Millisecond * 10)
	close(ch)

	wg.Wait()

	// 之前代码执行时间应该和等待时间相同
	assertion.DurationMatch(t, 10*time.Millisecond, d)

	// 缓冲区写满后被消费, 发送不阻塞

	wg = sync.WaitGroup{}
	wg.Add(1)

	// 第二个参数指定了缓冲区大小, 即缓存多少个发送实例
	ch = make(chan string, 1)

	// 异步函数, 向 chan 中发送数据
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

	// 接收 chan 中的数据, 此时发送方立即完成
	s := <-ch
	assert.Equal(t, "Hello", s)
	close(ch)

	wg.Wait()
	// 代码执行时间短暂
	assert.LessOrEqual(t, d, time.Millisecond)
}

// 返回一个 chan, 可以读取通过 fn 函数不断的产生数据
func Generator[T any](fn func(c chan T)) chan T {
	c := make(chan T)
	go func() {
		fn(c)
		close(c)
	}()
	return c
}

// 测试通过 range 关键字读取 chan 中的数据
//
// 可以通过 range 关键字读取一个 chan 中的数据, 知道 chan 被关闭
func TestChan_Range(t *testing.T) {
	// 产生一个 chan 实例, 并定义向 chan 中写数据的函数
	gen := Generator(func(ch chan string) {
		for i := range 10 {
			ch <- strconv.Itoa(i)
		}
	})

	r := make([]string, 0, 100)

	// 通过 range 关键字逐个从 chan 中读取数据, 直到 chan 被关闭
	for s := range gen {
		r = append(r, s)
	}

	// 确认读取的数据正确
	assert.Len(t, r, 10)
	assert.Equal(t, []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}, r)
}
