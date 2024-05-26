package blockque

import (
	"container/list"
	"context"
	"sync"

	"golang.org/x/sync/semaphore"
)

// 定义阻塞队列结构体
//
// 信号量的一个重要作用就是定义阻塞队列, 用于多个 goroutine 之间交换数据 (生产者消费者模式)
//
// 具体思路是: 定义队列长度 N, 并设置 N 个信号量值, 当入队一个元素时, 占用一个信号量值, 出队一个元素则释放
// 一个信号量值
//
// 当队列元素数量达到 N 时, 所有信号量值都被占用, 继续占用信号量会导致阻塞, 直到另一个并行程序出队了一个元素并释放信号量
type BlockQueue[T any] struct {
	l   *list.List
	sem semaphore.Weighted
	mux sync.RWMutex
}

// 创建结构体实例
func New[T any](size int64) *BlockQueue[T] {
	return &BlockQueue[T]{
		l:   list.New(),
		sem: *semaphore.NewWeighted(size),
	}
}

// 获取队列中链表的长度
func (bq *BlockQueue[T]) len() int {
	return bq.l.Len()
}

// 获取队列中现存全部元素的切片
func (bq *BlockQueue[T]) List() []T {
	bq.mux.RLock()
	defer bq.mux.RUnlock()

	len := bq.len()

	s := make([]T, len)
	elem := bq.l.Front()

	for i := 0; i < len; i++ {
		s[i] = elem.Value.(T)
		elem = elem.Next()
	}
	return s
}

// 获取队列长度
func (bq *BlockQueue[T]) Len() int {
	bq.mux.RLock()
	defer bq.mux.RUnlock()

	return bq.len()
}

// 获取队列是否为空
func (bq *BlockQueue[T]) Empty() bool {
	return bq.Len() == 0
}

// 将元素加入队列
//
// 如果队列已满, 则该方法会阻塞, 直到队列空出至少一个元素后方能入队成功
func (bq *BlockQueue[T]) Offer(ctx context.Context, val T) bool {
	if err := bq.sem.Acquire(ctx, 1); err != nil {
		return false
	}

	bq.mux.Lock()
	defer bq.mux.Unlock()

	bq.l.PushBack(val)

	return true
}

// 尝试将元素加入队列
//
// 如果队列已满, 则加入元素失败, 返回 `false`
func (bq *BlockQueue[T]) TryOffer(val T) bool {
	if ok := bq.sem.TryAcquire(1); !ok {
		return false
	}

	bq.mux.Lock()
	defer bq.mux.Unlock()

	bq.l.PushBack(val)
	return true
}

// 从队列的头部弹出一个元素
//
// 从队列中弹出表示, 获取队列头部元素并将其删除
//
// 如果队列为空, 则返回 `defValue` 参数表示的默认值及 `false` 值
func (bq *BlockQueue[T]) Poll(defVal T) (T, bool) {
	bq.mux.Lock()
	defer bq.mux.Unlock()

	elem := bq.l.Front()
	if elem == nil {
		return defVal, false
	}

	bq.l.Remove(elem)

	bq.sem.Release(1)
	return elem.Value.(T), true
}

// 从队列的头部删除一个元素
// 如果队列为空, 则返回 `false` 值
func (bq *BlockQueue[T]) Remove() bool {
	bq.mux.Lock()
	defer bq.mux.Unlock()

	elem := bq.l.Front()
	if elem == nil {
		return false
	}

	bq.l.Remove(elem)
	bq.sem.Release(1)

	return true
}

// 获取队列的头部元素, 但不从队列中删除该元素
func (bq *BlockQueue[T]) Peek(defVal T) (T, bool) {
	bq.mux.RLock()
	defer bq.mux.RUnlock()

	elem := bq.l.Front()
	if elem == nil {
		return defVal, false
	}

	return elem.Value.(T), true
}
