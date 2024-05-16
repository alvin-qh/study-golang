package channellock

import (
	"context"
	"sync"
	"sync/atomic"
)

// 定义 Channel 锁结构体
type Lock struct {
	ch     chan struct{} // 用来等待的 channel 对象
	mut    sync.Mutex    // 用来锁定 ch 字段的互斥量
	locked atomic.Value  // 用来记录是否锁定的原子量
}

// 产生一个新的锁对象
func New() *Lock {
	// 产生无缓冲的 channel 对象
	l := &Lock{ch: make(chan struct{})}
	// 初始化原子量
	l.locked.Store(false)
	return l
}

// 获取 channel 对象
func (l *Lock) channel() <-chan struct{} {
	// 通过互斥量进入临界区
	l.mut.Lock()
	defer l.mut.Unlock()

	// 返回结构体的 channel 对象
	return l.ch
}

// 锁定
func (l *Lock) Lock(ctx context.Context) bool {
	// 当原子量为 false 时, 设置其为 true, 并进入临界区
	for !l.locked.CompareAndSwap(false, true) {
		// 临界区

		// 获取 channel 对象
		ch := l.channel()

		select {
		case <-ch: // channel 被关闭, 表示被解锁
			// 继续循环, 直到完成锁定
		case <-ctx.Done(): // 超时时间到, 如果还未完成锁定, 则返回未锁定成功
			return false
		}
	}
	return true
}

// 交换结构体中的 channel 对象
func (l *Lock) swapChannel(new chan struct{}) chan struct{} {
	l.mut.Lock()
	defer l.mut.Unlock()

	old := l.ch
	// 设置新的 channel 对象
	l.ch = new

	// 返回旧的 channel 对象
	return old
}

// 解锁
func (l *Lock) Unlock() {
	// 判断锁对象是否锁定, 若未锁定则忽略操作, 否则进入临界区
	if l.locked.CompareAndSwap(true, false) {
		// 临界区

		defer func() {
			// 判断是否有错误发送
			if recover() != nil {
				l.ch = nil
			}
		}()

		// 设置新的 channel 对象, 返回原有 channel 对象
		ch := l.swapChannel(make(chan struct{}, 1))

		// 关闭 channel 对象, 表示解锁
		close(ch)
	}
}
