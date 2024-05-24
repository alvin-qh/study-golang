package pool

import (
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

// 测试同步池
//
// "池" (`sync.Pool`) 是一个 goroutine 安全的实例容器, 可以从池中获取对象, 使用完毕后归还到池中
//
// 如果池空了 (新池或所有实例都被取出), 则会自动调用池的 `New` 方法创建实例
//
// 池的主要作用为管理大实例, 大实例的创建和销毁会消耗较多的计算资源, 通过池可以显著降低实例的创建次数,
// 提高系统效率
func TestPool_GetAndPut(t *testing.T) {
	// 创建池, 定义实例创建函数
	pool := sync.Pool{
		New: func() interface{} {
			return &Object{
				state: NEW,
			}
		},
	}

	var wg sync.WaitGroup

	// 创建 10 个 goroutine
	for i := 0; i < 10; i++ {
		wg.Add(1)

		// 创建 goroutine, 每个 goroutine 中利用池获取 10 次 `Object` 类型实例,
		// 在获取实例时, 一部分实例会新建, 一部分实例会从池中直接获取, 通过实例状态加以区分
		go func() {
			for j := 0; j < 10; j++ {
				// 从池中获取一个实例
				obj := pool.Get().(*Object)

				// 模拟实例操作时长
				time.Sleep(10 * time.Millisecond)

				// 判断实例的状态, 如果是新建状态, 则将其转为已存在状态
				if obj.State() == NEW {
					obj.ChangeState(OLD)
				}

				// 将用完的实例返回池
				pool.Put(obj)

				// 模拟其它流程执行时长
				time.Sleep(20 * time.Millisecond)
			}

			// 任务执行完毕
			wg.Done()
		}()
	}

	// 等待所有 goroutine 执行完毕
	wg.Wait()

	// 收集池中所有标记为 OLD 的实例, 查看共产生了多少个实例
	objs := make([]*Object, 0, 100)
	for {
		obj := pool.Get().(*Object)
		if obj.State() != OLD {
			break
		}
		objs = append(objs, obj)
	}

	// 确认产生的实例总数不超过 20, 实际使用实例次数为 100
	assertion.Between(t, len(objs), 1, 20)
}
