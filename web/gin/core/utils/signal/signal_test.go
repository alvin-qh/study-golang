package signal

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试 `WaitInterruptSignal` 函数
//
// 确认是否可以正确的接收到 `SIGINT` 信号
func TestWaitInterruptSignal(t *testing.T) {
	t.Skip("One process cannot have two signal notify")

	var until time.Duration = 0

	// 用于等待协程结束的通道
	ch := make(chan struct{})

	// 协程函数, 在协程内等待进程结束信号
	go func() {
		start := time.Now()

		// 等待进程结束信号
		WaitInterruptSignal()

		// 计算等待时间
		until = time.Since(start)

		// 发送表示协程结束的对象
		ch <- struct{}{}
	}()

	// 主线程等待 2 秒
	time.Sleep(2 * time.Second)

	// 获取当前进程对象
	process, err := os.FindProcess(syscall.Getpid())
	assert.NoError(t, err)

	// 向当前进程发送 `SIGINT` 信号
	err = process.Signal(os.Interrupt)
	assert.NoError(t, err)

	// 等待协程结束
	<-ch

	// 确认协程结束时间在 2s 左右
	assert.GreaterOrEqual(t, until, 2*time.Second)
	assert.Less(t, until, 2100*time.Millisecond)
}
