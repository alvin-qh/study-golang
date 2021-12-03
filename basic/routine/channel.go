package routine

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
)

// 定义锁结构体
type Lock struct {
	ch     chan struct{} // 用来等待的 channel 对象
	mut    sync.Mutex    // 用来锁定 ch 字段的互斥量
	locked atomic.Value  // 用来记录是否锁定的原子量
}

// 产生一个新的锁对象
func NewLock() *Lock {
	l := &Lock{ch: make(chan struct{})} // 产生无缓冲的 channel 对象
	l.locked.Store(false)               // 初始化原子量
	return l
}

// 锁定
func (l *Lock) Lock(ctx context.Context) bool {
	for !l.locked.CompareAndSwap(false, true) { // 当原子量为 false 时，设置其为 true，并进入临界区

		// 临界区

		ch := l.channel() // 获取 channel 对象
		select {
		case <-ch: // channel 被关闭，表示被解锁，
			// 继续循环，直到完成锁定
		case <-ctx.Done(): // 超时时间到，如果还未完成锁定，则返回未锁定成功
			return false
		}
	}
	return true
}

// 获取 channel 对象
func (l *Lock) channel() <-chan struct{} {
	l.mut.Lock() // 通过互斥量进入临界区
	defer l.mut.Unlock()

	return l.ch // 返回结构体的 channel 对象
}

// 解锁
func (l *Lock) Unlock() {
	if l.locked.CompareAndSwap(true, false) { // 判断锁对象是否锁定，若未锁定则忽略操作，否则进入临界区
		// 临界区

		defer func() {
			if recover() != nil { // 判断是否有错误发送
				l.ch = nil
			}
		}()

		ch := l.swapChannel(make(chan struct{}, 1)) // 设置新的 channel 对象，返回原有 channel 对象
		close(ch)                                   // 关闭 channel 对象，表示解锁
	}
}

// 交换结构体中的 channel 对象
func (l *Lock) swapChannel(newOne chan struct{}) chan struct{} {
	l.mut.Lock()
	defer l.mut.Unlock()

	oldOne := l.ch
	l.ch = newOne // 设置新的 channel 对象
	return oldOne // 返回旧的 channel 对象
}

// 生成器错误集合
var (
	ErrChanClosed = errors.New("channel was closed")
)

// 定义生成器结构体
type Generator struct {
	ch chan interface{}
}

// 定义数据生成函数
type GeneratorFunc func(ch chan interface{}) interface{}

// 产生一个新的生成器对象
func NewGenerator(gen GeneratorFunc) *Generator {
	g := &Generator{ch: make(chan interface{})}
	runtime.SetFinalizer(g, func(g *Generator) { g.Close() }) // 析构函数，关闭生成器对象

	go func() {
		defer recover()
		gen(g.ch) // 在协程中异步生成数据
	}()

	return g
}

// 关闭生成器对象
func (g *Generator) Close() {
	if g.ch != nil {
		close(g.ch)
		g.ch = nil
	}
}

// 获取生成器下一个数据
func (g *Generator) Next() (interface{}, error) {
	if g.ch == nil {
		return nil, ErrChanClosed
	}

	// 从 channel 中接收下一个数据
	v, ok := <-g.ch
	if !ok {
		return nil, ErrChanClosed
	}
	return v, nil
}

func (g *Generator) Range() <-chan interface{} {
	if g.ch == nil {
		panic(ErrChanClosed)
	}
	return g.ch
}
