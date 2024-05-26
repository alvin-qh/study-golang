package pool

import (
	"context"
	"fmt"
	"study/basic/builtin/slice/utils"
	"study/basic/testing/assertion"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试用结构体实例
type Value struct {
	Val int
}

func (v *Value) String() string {
	return fmt.Sprintf("%d", v.Val)
}

// 测试创建池对象
func TestSizedPool_New(t *testing.T) {
	// 创建一个具有 10 个元素容量的池
	pool := NewSizedPool(10, func() *Value {
		return &Value{}
	})

	assert.Equal(t, 10, pool.Size())
	assert.Equal(t, 10, pool.MaxSize())
}

// 测试从池中尝试获取一个元素
func TestSizedPool_TryGet(t *testing.T) {
	lastId := int32(0)

	pool := NewSizedPool(10, func() *Value {
		id := atomic.AddInt32(&lastId, 1)
		if id > 10 {
			assert.Fail(t, "")
		}
		return &Value{Val: int(id)}
	})

	rs := make([]int, 0, 10)

	var mux sync.Mutex
	var wg sync.WaitGroup

	// 启动 20 个 goroutine
	for i := 0; i < 20; i++ {
		wg.Add(1)

		// 启动 goroutine, 并尝试从池中获取元素
		go func() {
			defer wg.Done()

			// 尝试从池中获取元素
			elem, ok := pool.TryGet()
			if ok {
				// 使用完毕后归还池
				defer elem.Release()

				// 如果成功从池中获取到元素, 则将获取到的元素实例存储到结果切片中
				mux.Lock()
				rs = append(rs, elem.Get().Val)
				mux.Unlock()

				// 模拟使用实例 50ms
				time.Sleep(50 * time.Millisecond)
			}
		}()
	}

	wg.Wait()

	// 确定所有池元素归还后, 池大小回复初始值
	assert.Equal(t, 10, pool.Size())

	// 确认从池中获取的共 10 个元素
	assert.ElementsMatch(t, utils.Range(1, 11, 1), rs)
}

// 测试从池中获取元素
//
// 通过 `NewSizedPool.Get` 方法可以从池中获取元素, 如果池此时为空, 则 `Get` 方法会阻塞直到池中有元素可用
func TestSizedPool_Get(t *testing.T) {
	lastId := int32(0)

	// 创建池实例, 容量为 10
	pool := NewSizedPool(10, func() *Value {
		id := atomic.AddInt32(&lastId, 1)
		if id > 10 {
			assert.Fail(t, "")
		}
		return &Value{Val: int(id)}
	})

	// 创建切片, 存储从池中获取的元素
	rs := make([]*SizedPoolElem[*Value], 0, pool.MaxSize())

	start := time.Now()

	// 将池中的元素全部取出, 放入切片中
	for pool.Size() > 0 {
		elem, err := pool.Get(context.Background())
		assert.Nil(t, err)

		rs = append(rs, elem)
	}

	// 上面的操作未发生任何阻塞, 因为池中现存元素
	assertion.Between(t, time.Since(start).Milliseconds(), int64(0), int64(10))

	// 共从池中获取 10 个元素
	assert.Len(t, rs, pool.MaxSize())

	// 目前池为空, 继续获取元素会导致阻塞或超时
	_, ok := pool.TryGet()
	assert.False(t, ok)

	// 启动 goroutine, 并以 10ms 间隔将将池元素归还池
	go func() {
		// 将之前取出的元素按照 10ms 的间隔逐一归还
		for _, elem := range rs {
			time.Sleep(10 * time.Millisecond)
			elem.Release()
		}
	}()

	start = time.Now()

	var wg sync.WaitGroup

	// 启动 10 个 goroutine
	for i := 0; i < pool.MaxSize(); i++ {
		wg.Add(1)

		// 启动 goroutine, 从池中获取一个元素
		go func() {
			defer wg.Done()

			// 创建可超时上下文实例, 由于池的元素整体经过 pool.MaxSize() * 10ms 后才能释放完毕,
			// 从池中获取元素的最长等待时间也也应为 pool.MaxSize() * 10ms, 增加 20 ms 作为其它运行损耗
			ctx, cancel := context.WithTimeout(
				context.Background(),
				time.Duration(pool.MaxSize()*10+20)*time.Millisecond,
			)
			defer cancel()

			// 从池中获取元素
			elem, err := pool.Get(ctx)

			assert.Nil(t, err)
			// 确认取出的元素包含在上次取出元素集合中
			assert.Contains(t, rs, elem)
		}()
	}

	// 等待所有 goroutine 执行完毕
	wg.Wait()

	// 确认创建的池元素总数
	assert.Equal(t, int32(pool.MaxSize()), lastId)

	// 确认第二次获取全部池元素消耗的时间
	assertion.Between(t, time.Since(start).Milliseconds(), int64(pool.MaxSize()*10), int64(pool.MaxSize()*10+20))
}
