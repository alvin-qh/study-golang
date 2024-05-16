package sync

import (
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试同步 Map 的基本操作
//
// 同步 Map, 即 `sync.Map` 类型一般用于并发场景, 当多个任务同时访问一个 Map 时, 必须使用锁,
// 否则会导致错误如果需要降低锁对性能的影响, 则需要使用 `sync.Map` 进行操作
func TestSyncMap_Store_Load(t *testing.T) {
	// 定义同步 Map 对象
	sm := sync.Map{}

	// 将键值对存入 Map 中
	sm.Store("A", 100)

	// 根据 Key 读取 Value, 返回 Value 值以及 Key 是否存在
	v, ok := sm.Load("A")
	assert.True(t, ok)
	assert.Equal(t, 100, v)
}

// 测试遍历同步 Map 中的键值对
func TestSyncMap_Range(t *testing.T) {
	sm := sync.Map{}

	for i := 0; i < 5; i++ {
		sm.Store(string([]rune{rune(65 + i)}), i)
	}

	// 定义保持 Key 和 Value 的切片
	ks := make([]string, 0)
	vs := make([]int, 0)

	// 遍历 Map 中的所有键值对
	sm.Range(func(key, value any) bool {
		ks = append(ks, key.(string))
		vs = append(vs, value.(int))

		return true
	})

	// 确认键值对
	assert.ElementsMatch(t, []string{"A", "B", "C", "D", "E"}, ks)
	assert.ElementsMatch(t, []int{0, 1, 2, 3, 4}, vs)
}

// 测试在多个协助中使用同步 Map
func TestSyncMap_InGoroutine(t *testing.T) {
	sm := sync.Map{}

	// 定义一个等待组, 设置 2 个等待任务
	wg := sync.WaitGroup{}
	wg.Add(2)

	// 执行任务 1
	go func() {
		// 表示任务结束后, 等待数减1
		defer wg.Done()
		for n := 0; n < 1000; n++ {
			// 向同步 Map 中添加元素
			sm.Store(n, n+1)
		}
	}()

	// 执行任务 2
	go func() {
		// 表示任务结束后, 等待数减1
		defer wg.Done()
		for n := 0; n < 1000; n++ {
			// 向同步 Map 中添加元素
			sm.Store(strconv.Itoa(n), strconv.Itoa(n+1))
		}
	}()

	// 等待组的数值为 0 时返回
	wg.Wait()

	// 确认两个协程都完成操作
	v, ok := sm.Load("999")
	assert.True(t, ok)
	assert.Equal(t, "1000", v)

	v, ok = sm.Load(999)
	assert.True(t, ok)
	assert.Equal(t, 1000, v)

}
