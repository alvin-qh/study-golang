package time

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试暂停当前线程 (协程)
func TestSleep(t *testing.T) {
	tm := time.Now()

	// 休眠 50ms
	time.Sleep(50 * time.Millisecond)

	// 计算当前时间和指定时间的差值, 相当于 time.Now().Sub(tm)
	d := time.Since(tm)
	assert.Equal(t, 50, int(d.Milliseconds()))
}

// 测试定时信号
// 定时信号可以在指定时间后, 通过 channel 发送一次性的定时信号
func TestTimeAfter(t *testing.T) {
	tm1 := time.Now()

	// 50ms 后发送信号
	c := time.After(50 * time.Millisecond)

	// 等待信号到达
	tm2 := <-c

	// 计算从发送信号到接收信号的时间差
	d := tm2.Sub(tm1)
	assert.Equal(t, 50, int(d.Milliseconds()))
}

// 测试定时回调
// 定时回调可以在指定时间后回调指定函数
func TestTimeAfterFunc(t *testing.T) {
	// 定义一个 channel
	ch := make(chan struct{})
	defer close(ch)

	tm := time.Now()

	// 50ms 后回调函数
	time.AfterFunc(50*time.Millisecond, func() {
		// 发送一个空的信号
		ch <- struct{}{}
	})

	// 等待信号到达
	<-ch

	// 计算函数多久后进行回调
	d := time.Since(tm)
	assert.Equal(t, 50, int(d.Milliseconds()))
}

// 测试周期性定时消息
func TestTicker(t *testing.T) {
	tm1 := time.Now()

	// 每隔 50ms 发送一次信号
	tk := time.NewTicker(50 * time.Millisecond)
	defer tk.Stop()

	n := 0

	// 获取两次信号
	for n < 2 {
		// 等待信号到达
		tm2 := <-tk.C

		// 计算每次信号到达的时间
		d := tm2.Sub(tm1)
		assert.Equal(t, (n+1)*50, int(d.Milliseconds()))

		n++
	}
}

// 测试一次性定时器
func TestTimer(t *testing.T) {
	tm1 := time.Now()

	// 50ms 后发送定时器信号
	ti := time.NewTimer(50 * time.Millisecond)
	defer ti.Stop()

	// 等待定时器信号
	tm2 := <-ti.C

	// 计算信号到达时间
	d := tm2.Sub(tm1)
	assert.Equal(t, 50, int(d.Milliseconds()))
}

// Timer 和 After 的异同
//    Timer 和 After 均可以在指定时间后发送一个信号, 以达到定时任务的效果
//    Timer (或者 Ticker) 具备更丰富的操作手段, 包括 Stop (打断计时) 和 Reset (重新计时)
//    After 的功能比较简单
