package orderedmap

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试 OrderedMap 对象
func TestOrderedMap_New(t *testing.T) {
	// 测试基本的 key/value 存储

	sm := New[int, string]() // 新建对象
	assert.Equal(t, 0, sm.Len())

	sm.Put(100, "A")

	// 使用相同的 key 再次存储, 因为 key 相同, 所以会覆盖掉之前的 key/value
	sm.Put(100, "B")

	// 目前只存储了 1 个 key
	keys := sm.Keys()
	assert.Len(t, keys, 1)

	// 利用 key 获取 value
	v, ok := sm.Get(100)
	assert.True(t, ok)
	assert.Equal(t, "B", v)

	sm.Put(10, "A")
	sm.Put(1000, "C")

	// 目前存储了 3 个 key
	keys = sm.Keys()
	assert.Len(t, keys, 3)

	// 获取到的 key 是有序的
	assert.Equal(t, []int{10, 100, 1000}, keys)

	sm.Put(1, "D")
	vs := sm.Values()

	// 目前存储了 4 个 key
	assert.Len(t, vs, 4)
	// values 的属性依据 keys 的顺序
	assert.Equal(t, []string{"D", "A", "B", "C"}, vs)
}

// 测试迭代内容
func TestOrderedMap_Do(t *testing.T) {
	sm := New[int, string]()

	sm.Put(100, "B")
	sm.Put(1000, "C")
	sm.Put(1, "D")
	sm.Put(10, "A")

	n := 0

	// 期待的有序 keys 集合
	expKeys := []int{1, 10, 100, 1000}

	// 期待的 values 集合
	expValues := []string{"D", "A", "B", "C"}

	// 按 key 的顺序依次迭代
	sm.Do(func(key int, value string) {
		assert.Equal(t, expKeys[n], key)
		assert.Equal(t, expValues[n], value)
		n++
	})
}

// 测试删除 key 和清空集合
func TestOrderedMap_Remove(t *testing.T) {
	sm := New[int, string]()

	sm.Put(100, "B")
	sm.Put(1000, "C")
	sm.Put(1, "D")
	sm.Put(10, "A")

	// 删除一个 key
	sm.Remove(100)
	// 剩余的 key 依旧有序
	assert.Equal(t, []int{1, 10, 1000}, sm.Keys())

	// 继续删除 key
	sm.Remove(10)
	assert.Equal(t, []int{1, 1000}, sm.Keys())


    assert.Equal(t, []string{"D", "C"}, sm.Values()) // 剩余的 values 依旧保持和 keys 的有序对应

	// 清空, 返回初始状态
	sm.Init()

	assert.Len(t, sm.m, 0)
	assert.Len(t, sm.s, 0)

	keys := sm.Keys()
	assert.Len(t, keys, 0)
	assert.True(t, sm.sorted)
}

// 测试迭代函数
func TestOrderedMap_Iterator(t *testing.T) {
	sm := New[int, string]()

	sm.Put(100, "B")
	sm.Put(1000, "C")
	sm.Put(1, "D")
	sm.Put(10, "A")

	ks := make([]int, 0, 4)
	vs := make([]string, 0, 4)

	it := sm.Iterate()
	for {
		k, v, ok := it()
		if !ok {
			break
		}
		ks = append(ks, k)
		vs = append(vs, v)
	}

	assert.Equal(t, []int{1, 10, 100, 1000}, ks)
	assert.Equal(t, []string{"D", "A", "B", "C"}, vs)
}
