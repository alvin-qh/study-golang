package pool

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"golang.org/x/sync/semaphore"
)

// 定义一个池中实例的类型
type SizedPoolElem[T any] struct {
	elem T             // 内容实例
	pool *SizedPool[T] // 池实例指针
}

// 将元素返回池
func (p *SizedPoolElem[T]) Release() {
	p.pool.put(p)
}

// 获取池元素中存储的实例
func (pe *SizedPoolElem[T]) Get() T { return pe.elem }

func (pe *SizedPoolElem[T]) String() string {
	return fmt.Sprintf("%v", pe.elem)
}

// 定义固定大小的池类型
//
// Go 语言的 `sync.Pool` 类型并未对池的大小设置上限, 这会导致如果有一波处理峰值, 就有可能瞬间将池中的实例消耗光,
// 进而继续通过池创建大量的实例
//
// 所以限制池的大小, 可以在处理数据峰值的时候, 起到限流的作用
type SizedPool[T any] struct {
	pool     sync.Pool           // 池, 用于存储指定类型的实例
	weighted *semaphore.Weighted // 信号量, 用于限制池的最大尺寸
	size     int64               // 池的当前大小
	maxSize  int64               // 池的最大容积
}

// 创建实例
//
// 设置池的对最大容量, 并设置创建池中实例的函数
func NewSizedPool[T any](size int, creator func() T) *SizedPool[T] {
	pool := SizedPool[T]{
		weighted: semaphore.NewWeighted(int64(size)),
		size:     int64(size),
		maxSize:  int64(size),
	}

	// 设定池的 New 函数, 用于创建实例
	pool.pool.New = func() any {
		return &SizedPoolElem[T]{
			elem: creator(),
			pool: &pool,
		}
	}

	return &pool
}

// 尝试从池中获取一个池元素实例
func (p *SizedPool[T]) TryGet() (elem *SizedPoolElem[T], ok bool) {
	// 尝试消费一个信号量值
	if p.weighted.TryAcquire(1) {
		ok = true

		// 如果信号量消费成功, 则从池中获取一个实例
		elem = p.pool.Get().(*SizedPoolElem[T])

		// 减少池当前大小
		atomic.AddInt64(&p.size, -1)
	}
	return
}

// 尝试从池中获取一个池元素实例
//
// 这里可以通过 `Context` 实例限制池空后, 等待元素返回池的最长超时时间
func (p *SizedPool[T]) Get(ctx context.Context) (*SizedPoolElem[T], error) {
	if err := p.weighted.Acquire(ctx, 1); err != nil {
		return nil, err
	}

	elem := p.pool.Get().(*SizedPoolElem[T])
	atomic.AddInt64(&p.size, -1)
	return elem, nil
}

// 将元素返回池
func (p *SizedPool[T]) put(elem *SizedPoolElem[T]) {
	p.pool.Put(elem)
	atomic.AddInt64(&p.size, 1)

	p.weighted.Release(1)
}

// 获取池当前大小
func (p *SizedPool[T]) Size() int {
	return int(atomic.LoadInt64(&p.size))
}

// 获取池的最大容量
func (p *SizedPool[T]) MaxSize() int { return int(p.maxSize) }
