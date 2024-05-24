package pool

import (
	"study/basic/builtin/slice/utils"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type Value struct {
	Val int
}

// 测试创建池对象
func TestSizedPool_New(t *testing.T) {
	// 创建一个具有 10 个元素容量的池
	pool := NewSizedPool(10, func() *Value {
		return &Value{}
	})

	assert.Equal(t, 10, pool.Size())
}

// 测试从池中尝试获取一个元素
func TestSizedPool_TryGet(t *testing.T) {
	lastId := int32(0)

	pool := NewSizedPool(10, func() *Value {
		return &Value{Val: int(atomic.AddInt32(&lastId, 1))}
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
				// 如果成功从池中获取到元素, 则将获取到的元素实例存储到结果切片中
				mux.Lock()
				rs = append(rs, elem.Get().Val)
				mux.Unlock()

				// 模拟使用实例 50ms
				time.Sleep(50 * time.Millisecond)

				// 使用完毕后归还池
				elem.Release()
			} else {
				// 如果获取元素失败, 则池一定为空
				assert.Equal(t, 0, pool.Size())
			}
		}()
	}

	wg.Wait()

	// 确定所有池元素归还后, 池大小回复初始值
	assert.Equal(t, 10, pool.Size())
	// 确认从池中获取的共 10 个元素
	assert.ElementsMatch(t, utils.Range(1, 11, 1), rs)
}
