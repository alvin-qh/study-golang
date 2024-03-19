package sortedmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试 SortedMap 对象
func TestSortedMap(t *testing.T) {
	// 测试基本的 key/value 存储

	sm := New(func(a, b interface{}) int { return b.(int) - a.(int) }) // 新建对象
	assert.Equal(t, 0, sm.Len())

	sm.Put(100, "A")
	sm.Put(100, "B") // 使用相同的 key 再次存储, 因为 key 相同, 所以会覆盖掉之前的 key/value

	keys := sm.Keys() // 目前只存储了 1 个 key
	assert.Len(t, keys, 1)

	v, ok := sm.Get(100) // 利用 key 获取 value
	assert.True(t, ok)
	assert.Equal(t, "B", v)

	sm.Put(10, "A")
	sm.Put(1000, "C")

	keys = sm.Keys() // 目前存储了 3 个 key
	assert.Len(t, keys, 3)
	assert.Equal(t, []interface{}{10, 100, 1000}, keys) // 获取到的 key 是有序的

	sm.Put(1, "D")
	vs := sm.Values()
	assert.Len(t, vs, 4)                                   // 目前存储了 4 个 key
	assert.Equal(t, []interface{}{"D", "A", "B", "C"}, vs) // values 的属性依据 keys 的顺序
}

// 测试迭代内容
func TestSortedMapRange(t *testing.T) {
	sm := New(func(a, b interface{}) int { return b.(int) - a.(int) })

	sm.Put(100, "B")
	sm.Put(1000, "C")
	sm.Put(1, "D")
	sm.Put(10, "A")

	n := 0
	expKeys := []interface{}{1, 10, 100, 1000}     // 期待的有序 keys 集合
	expValues := []interface{}{"D", "A", "B", "C"} // 期待的 values 集合

	// 按 key 的顺序依次迭代
	sm.Range(func(key interface{}, value interface{}) {
		assert.Equal(t, expKeys[n], key)
		assert.Equal(t, expValues[n], value)
		n++
	})
}

// 测试删除 key 和清空集合
func TestSortedMapRemoveAndClear(t *testing.T) {
	sm := New(func(a, b interface{}) int { return b.(int) - a.(int) })

	sm.Put(100, "B")
	sm.Put(1000, "C")
	sm.Put(1, "D")
	sm.Put(10, "A")

	sm.Remove(100)                                         // 删除一个 key
	assert.Equal(t, []interface{}{1, 10, 1000}, sm.Keys()) // 剩余的 key 依旧有序

	sm.Remove(10) // 继续删除 key
	assert.Equal(t, []interface{}{1, 1000}, sm.Keys())

	assert.Equal(t, []interface{}{"D", "C"}, sm.Values()) // 剩余的 values 依旧保持和 keys 的有序对应

	sm.Clear() // 清空, 返回初始状态

	assert.Len(t, sm.m, 0)
	assert.Len(t, sm.s, 0)

	keys := sm.Keys()
	assert.Len(t, keys, 0)
	assert.True(t, sm.sorted)
}
