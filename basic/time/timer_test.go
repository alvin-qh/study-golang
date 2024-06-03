package time

import (
	"study/basic/testing/assertion"
	"testing"
	"time"
)

// 测试暂停当前 goroutine
func TestTime_Sleep(t *testing.T) {
	tm := time.Now()

	// 休眠 50ms
	time.Sleep(120 * time.Millisecond)

	// 计算当前时间和指定时间的差值, 相当于 time.Now().Sub(tm)
	assertion.Between(t, time.Since(tm).Milliseconds(), int64(120), int64(140))
}

// 测试定时信号
//
// 定时信号可以在指定时间后, 通过 channel 发送一次性的定时信号
func TestTime_After(t *testing.T) {
	tm1 := time.Now()

	// 50ms 后发送信号
	c := time.After(120 * time.Millisecond)

	// 等待信号到达
	tm2 := <-c

	// 计算从发送信号到接收信号的时间差
	assertion.Between(t, tm2.Sub(tm1).Milliseconds(), int64(120), int64(140))
}

// 测试定时回调
//
// 定时回调可以在指定时间后回调指定函数
func TestTime_AfterFunc(t *testing.T) {
	// 定义一个 channel
	ch := make(chan struct{})
	defer close(ch)

	// 50ms 后回调函数
	time.AfterFunc(120*time.Millisecond, func() {
		// 发送一个空的信号
		ch <- struct{}{}
	})

	tm := time.Now()

	// 等待信号到达
	<-ch

	// 计算函数多久后进行回调
	assertion.Between(t, time.Since(tm).Milliseconds(), int64(120), int64(140))
}

// 测试周期性定时消息
func TestTime_NewTicker(t *testing.T) {
	tm1 := time.Now()

	// 每隔 50ms 发送一次信号
	tk := time.NewTicker(120 * time.Millisecond)
	defer tk.Stop()

	n := 0

	// 获取两次信号
	for n < 2 {
		// 等待信号到达
		tm2 := <-tk.C

		// 计算每次信号到达的时间
		assertion.Between(t, tm2.Sub(tm1).Milliseconds(), int64((n+1)*120), int64((n+1)*140))
		n++
	}
}

// 测试一次性定时器
func TestTime_NewTimer(t *testing.T) {
	tm1 := time.Now()

	// 50ms 后发送定时器信号
	ti := time.NewTimer(120 * time.Millisecond)
	defer ti.Stop()

	// 等待定时器信号
	tm2 := <-ti.C

	// 计算信号到达时间
	assertion.Between(t, tm2.Sub(tm1).Milliseconds(), int64(120), int64(140))
}

// Timer 和 After 的异同
//    Timer 和 After 均可以在指定时间后发送一个信号, 以达到定时任务的效果
//    Timer (或者 Ticker) 具备更丰富的操作手段, 包括 Stop (打断计时) 和 Reset (重新计时)
//    After 的功能比较简单
