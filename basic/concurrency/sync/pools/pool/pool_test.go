package pool_test

import (
	"study/basic/concurrency/sync/pools/pool"
	"study/basic/testing/assertion"
	"sync"
	"testing"
	"time"
)

// 实例状态
//
// 池中无法满足, 新创建的实例状态为 NEW,
// 回到池中的状态为 OLD
const (
	NEW = iota
	OLD
)

// 用于测试的结构体
type Object struct {
	// 实例状态
	state int
}

// 获取实例状态
func (o *Object) State() int {
	return o.state
}

// 修改实例状态
func (o *Object) ChangeState(state int) {
	o.state = state
}

// 测试通过池元素实例简化池元素返回池的操作
func TestPool_PoolElem(t *testing.T) {
	// 创建池, 定义实例创建函数
	pool := pool.New(func() *Object {
		return &Object{
			state: NEW,
		}
	})

	// 定义等待组对象, 用于等待全部任务完成
	var wg sync.WaitGroup

	// 创建 10 个 goroutine, 并在
	for range 10 {
		// 创建 goroutine, 每个 goroutine 中利用池获取 10 次 `Object` 类型实例,
		// 在获取实例时, 一部分实例会新建, 一部分实例会从池中直接获取, 通过实例状态加以区分
		wg.Go(func() {
			for range 10 {
				// 从池中获取一个实例
				elem := pool.Get()

				// 模拟实例操作时长
				time.Sleep(10 * time.Millisecond)

				// 判断实例的状态, 如果是新建状态, 则将其转为已存在状态
				if elem.Get().State() == NEW {
					elem.Get().ChangeState(OLD)
				}

				// 将用完的实例返回池
				elem.Release()

				// 模拟其它流程执行时长
				time.Sleep(20 * time.Millisecond)
			}

			// 任务执行完毕
		})
	}

	// 等待所有 goroutine 执行完毕
	wg.Wait()

	// 收集池中所有标记为 OLD 的实例, 查看共产生了多少个实例
	objs := make([]*Object, 0, 100)
	for {
		elem := pool.Get()
		if elem.Get().State() != OLD {
			break
		}
		objs = append(objs, elem.Get())
	}

	// 确认产生的实例总数不超过 20, 实际使用实例次数为 100
	assertion.Between(t, len(objs), 1, 20)
}
