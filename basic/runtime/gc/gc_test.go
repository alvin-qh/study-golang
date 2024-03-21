package gc

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 用于占用内存的结构体
type memory struct {
	data []byte
}

// 用于分配内存
func allocate(size uint32) *memory {
	return &memory{make([]byte, size)}
}

func (m *memory) Clear() {
	m.data = nil
}

// 测试 GC 内存回收
//
// 本例中通过 `runtime.ReadMemStats` 函数来收集实时内存状态, 以判断 GC 执行的情况
func TestGc(t *testing.T) {
	runtime.GC() // 手动唤起 GC 释放堆

	// 内存状态结构体
	var ms runtime.MemStats

	// 获取内存状态
	runtime.ReadMemStats(&ms)
	fmt.Printf("Alloc: %d(KB) HeapIdle: %d(KB) HeapReleased: %d(KB)\n", ms.Alloc/1024, ms.HeapIdle/1024, ms.HeapReleased/1024)

	// 分配内存, 用于测试 GC 回收内存
	m := allocate(1024 * 100000)

	finalized := false

	// 调用 GC 时, 会同时调用设置在对象上的 finalizer 函数, 释放资源
	runtime.SetFinalizer(m, func(m *memory) { finalized = true })

	// 获取内存状态
	runtime.ReadMemStats(&ms)
	fmt.Printf("Alloc: %d(KB) HeapIdle: %d(KB) HeapReleased: %d(KB)\n", ms.Alloc/1024, ms.HeapIdle/1024, ms.HeapReleased/1024)

	// 释放内存
	m.Clear()

	// 手动唤起 GC 释放堆
	runtime.GC()

	// 获取内存状态
	runtime.ReadMemStats(&ms)
	fmt.Printf("Alloc: %d(KB) HeapIdle: %d(KB) HeapReleased: %d(KB)\n", ms.Alloc/1024, ms.HeapIdle/1024, ms.HeapReleased/1024)

	assert.True(t, finalized)
}
