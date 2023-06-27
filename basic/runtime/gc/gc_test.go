package gc

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试 GC
// 可以通过 runtime.ReadMemStats 函数来获取 Go 内存状态, 以判断 GC 执行的情况
func TestGc(t *testing.T) {
	runtime.GC() // 手动唤起 GC 释放堆

	var ms runtime.MemStats // 内存状态结构体

	runtime.ReadMemStats(&ms) // 获取内存状态
	fmt.Printf("Alloc: %d(KB) HeapIdle: %d(KB) HeapReleased: %d(KB)\n", ms.Alloc/1024, ms.HeapIdle/1024, ms.HeapReleased/1024)

	m := Memory{}

	finalized := false
	runtime.SetFinalizer(&m, func(m *Memory) { finalized = true }) // 调用 GC 时, 会同时调用设置在对象上的 finalizer 函数, 释放资源

	size := 1024 * 100000

	m.Alloc(size) // 分配内存
	assert.Equal(t, size, m.Size())

	runtime.ReadMemStats(&ms) // 获取内存状态
	fmt.Printf("Alloc: %d(KB) HeapIdle: %d(KB) HeapReleased: %d(KB)\n", ms.Alloc/1024, ms.HeapIdle/1024, ms.HeapReleased/1024)

	m.Clear()    // 释放内存
	runtime.GC() // 手动唤起 GC 释放堆

	runtime.ReadMemStats(&ms) // 获取内存状态
	fmt.Printf("Alloc: %d(KB) HeapIdle: %d(KB) HeapReleased: %d(KB)\n", ms.Alloc/1024, ms.HeapIdle/1024, ms.HeapReleased/1024)

	assert.True(t, finalized)
}
