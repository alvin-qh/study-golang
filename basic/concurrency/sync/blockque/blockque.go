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
	lst *list.List         // 存储数据的单项链表
	sem semaphore.Weighted // 用于控制队列长度的信号量对象
	mux sync.RWMutex       // 保护链表的读写互斥锁
}

// 创建 BlockQueue 结构体实例
// `size` 参数指定了队列的长度, 即队列中最多可以存储多少个元素
func New[T any](size int64) *BlockQueue[T] {
	// 创建一个 BlockQueue 结构体实例, 并返回该实例的指针
	return &BlockQueue[T]{
		lst: list.New(),
		sem: *semaphore.NewWeighted(size),
	}
}

// 获取队列中链表的长度
func (bq *BlockQueue[T]) Len() int {
	return bq.lst.Len()
}

// 获取队列中现存全部元素的切片
func (bq *BlockQueue[T]) List() []T {
	// 锁定互斥量
	bq.mux.RLock()

	// 在函数返回前解锁互斥量
	defer bq.mux.RUnlock()

	// 获取队列中元素的数量, 并创建一个对应长度的切片用于存储这些元素
	len := bq.Len()

	s := make([]T, len)
	elem := bq.lst.Front()

	// 遍历链表, 将每个元素的值存储到切片中
	for i := range len {
		s[i] = elem.Value.(T)
		elem = elem.Next()
	}
	return s
}

// 获取队列是否为空
func (bq *BlockQueue[T]) Empty() bool {
	return bq.Len() == 0
}

// 将元素加入队列
//
// `ctx` 参数用于控制入队操作的上下文, 该参数可以用于控制入队操作的截止时间等
// `val` 参数表示要加入队列的元素
//
// 如果队列已满, 则该方法会阻塞, 直到队列空出至少一个元素后方能入队成功
func (bq *BlockQueue[T]) Offer(ctx context.Context, val T) bool {
	// 尝试占用一个信号量值, 如果队列已满, 则该方法会阻塞, 直到队列空出至少一个元素后方能入队成功
	if err := bq.sem.Acquire(ctx, 1); err != nil {
		return false
	}

	// 锁定互斥量, 在函数返回前解锁互斥量
	bq.mux.Lock()
	defer bq.mux.Unlock()

	// 将元素加入链表的尾部
	bq.lst.PushBack(val)

	return true
}

// 尝试将元素加入队列
//
// `val` 参数表示要加入队列的元素
//
// 如果队列已满, 则加入元素失败, 返回 `false`
func (bq *BlockQueue[T]) TryOffer(val T) bool {
	// 尝试占用一个信号量值, 如果队列已满, 则该方法会返回 `false`
	if ok := bq.sem.TryAcquire(1); !ok {
		return false
	}

	// 锁定互斥量, 在函数返回前解锁互斥量
	bq.mux.Lock()
	defer bq.mux.Unlock()

	// 将元素加入链表的尾部
	bq.lst.PushBack(val)
	return true
}

// 从队列的头部弹出一个元素
//
// `defValue` 参数表示如果队列为空时返回的默认值
//
// 从队列中弹出表示, 获取队列头部元素并将其删除
//
// 如果队列为空, 则返回 `defValue` 参数表示的默认值及 `false` 值
func (bq *BlockQueue[T]) Poll(defVal T) (T, bool) {
	// 锁定互斥量, 在函数返回前解锁互斥量
	bq.mux.Lock()
	defer bq.mux.Unlock()

	// 获取队列头部元素
	elem := bq.lst.Front()
	if elem == nil {
		return defVal, false
	}

	// 删除队列头部元素
	bq.lst.Remove(elem)

	// 释放一个信号量值
	bq.sem.Release(1)

	// 返回队列头部元素的值
	return elem.Value.(T), true
}

// 从队列的头部删除一个元素
//
// 如果队列为空, 则返回 `false` 值
func (bq *BlockQueue[T]) Remove() bool {
	// 锁定互斥量, 在函数返回前解锁互斥量
	bq.mux.Lock()
	defer bq.mux.Unlock()

	// 获取队列头部元素
	elem := bq.lst.Front()
	if elem == nil {
		return false
	}

	// 删除队列头部元素
	bq.lst.Remove(elem)

	// 释放一个信号量值
	bq.sem.Release(1)
	return true
}

// 获取队列的头部元素, 但不从队列中删除该元素
//
// `defValue` 参数表示如果队列为空时返回的默认值
//
// 如果队列为空, 则返回 `defValue` 参数表示的默认值及 `false` 值
func (bq *BlockQueue[T]) Peek(defVal T) (T, bool) {
	// 锁定互斥量, 在函数返回前解锁互斥量
	bq.mux.RLock()
	defer bq.mux.RUnlock()

	// 获取队列头部元素, 但不从队列中删除该元素
	elem := bq.lst.Front()
	if elem == nil {
		return defVal, false
	}

	// 返回队列头部元素的值
	return elem.Value.(T), true
}
