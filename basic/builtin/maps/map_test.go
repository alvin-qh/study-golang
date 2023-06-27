package maps

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 创建 map 对象
func TestCreateMap(t *testing.T) {
	// 定义一个 key 为 string 类型, value 为 int 类型的 map 变量
	var m map[string]int
	assert.Nil(t, m)           // 变量此时为 nil
	assert.Equal(t, 0, len(m)) // nil map 的长度为 0

	// 定义 map 并初始化
	m = map[string]int{
		"a": 100,
		"b": 200,
	}
	assert.Equal(t, 2, len(m))   // map 中包含 2 个 key
	assert.Equal(t, 100, m["a"]) // 根据 key 获取 value
	assert.Equal(t, 200, m["b"]) // 根据 key 获取 value

	// 通过 make 函数初始化 map, 第二个参数为 map 的初始容积, 默认为 0
	m = make(map[string]int, 100)
	assert.Equal(t, 0, len(m)) // 此时 map 长度为 0
}

// 从 map 中添加或删除 key
func TestAddAndRemoveMapKey(t *testing.T) {
	m := map[string]int{} // 声明一个空 map
	m["a"] = 100          // 设置 key "a"

	v, exist := m["a"] // 获取 key "a" 是否存在, 以及其值
	assert.True(t, exist)
	assert.Equal(t, 100, v)

	delete(m, "a") // 删除 key "a"

	_, exist = m["a"] // 判断 key "a" 是否存在
	assert.False(t, exist)
}

// 遍历 map
func TestThroughMap(t *testing.T) {
	m := map[string]interface{}{
		"a": 100,
		"b": "B",
		"c": []int{1, 2, 3},
	}

	// 遍历数组
	ks := make([]string, 0, len(m))
	vs := make([]interface{}, 0, len(m))

	for k, v := range m { // 遍历 key/value
		ks = append(ks, k)
		vs = append(vs, v)
	}

	assert.ElementsMatch(t, []string{"a", "b", "c"}, ks)
	assert.ElementsMatch(t, []interface{}{100, "B", []int{1, 2, 3}}, vs)

	// 遍历所有 key
	ks = make([]string, 0, len(m))

	for k := range m {
		ks = append(ks, k)
	}
	assert.ElementsMatch(t, []string{"a", "b", "c"}, ks)

	vs = make([]interface{}, 0, len(m))

	// 遍历所有的 value
	for _, v := range m {
		vs = append(vs, v)
	}
	assert.ElementsMatch(t, []interface{}{100, "B", []int{1, 2, 3}}, vs)
}

// 测试复合类型作为 `map` 的 key
// 复合类型即结构体 `struct`, 由于 go 语言支持接口体的比较和散列, 所以结构体可以直接作为 `map` 的 key
func TestComplexMapKey(t *testing.T) {
	// 定义结构体作为 map key
	type Key struct {
		id   int
		name string
	}

	// 定义结构体作为 map value
	type Value struct {
		gender   rune
		birthday string
		address  string
	}

	// 使用复杂类型作为 map key
	m := map[Key]*Value{}
	m[Key{1, "Alvin"}] = &Value{gender: 'M', birthday: "1981-03", address: "ShanXi, Xi'an"}

	assert.Equal(t, Value{gender: 'M', birthday: "1981-03", address: "ShanXi, Xi'an"}, *(m[Key{1, "Alvin"}]))
}

// 测试同步 `sync.Map`
// 同步 `Map` 用于异步场合, 当多个任务同时访问一个 map 时, 必须使用锁, 否则会导致错误
// 如果需要降低锁对性能的影响, 则需要使用 `sync.Map` 进行操作
func TestSyncMap(t *testing.T) {
	// 定义同步 map 对象
	sm := sync.Map{}

	// 定义一个等待组, 设置 2 个等待任务
	wg := sync.WaitGroup{}
	wg.Add(2)

	// 开启任务 1
	go func() {
		defer wg.Done() // 表示任务结束后, 等待数 减1

		for n := 0; n < 1000; n++ {
			sm.Store(n, n+1) // 向 同步map 中添加元素
		}
	}()

	// 开启任务 2
	go func() {
		defer wg.Done()

		for n := 0; n < 1000; n++ {
			sm.Store(fmt.Sprintf("%d", n), fmt.Sprintf("%d", n+1)) // 向 同步map 中添加元素
		}
	}()

	wg.Wait() // 等待组的数值为 0 时返回

	// 通过 key 读取 value
	v, ok := sm.Load("999")
	assert.True(t, ok) // ok 表示 key 是否存在
	assert.Equal(t, "1000", v)

	v, ok = sm.Load(999)
	assert.True(t, ok)
	assert.Equal(t, 1000, v)

	// 通过 key 删除 key
	sm.Delete("999")
	_, ok = sm.Load("999")
	assert.False(t, ok)

	// 通过 key 读取并同时删除
	v, ok = sm.LoadAndDelete(999)
	assert.True(t, ok) // ok 表示 key 是否存在
	assert.Equal(t, 1000, v)

	_, ok = sm.Load(999)
	assert.False(t, ok)

	// 通过 key 读取 Value 否则存储新 Value
	_, ok = sm.LoadOrStore(999, 1000)
	assert.False(t, ok) // ok 表示存储新 Value 前 key 是否存在

	v, ok = sm.Load(999)
	assert.True(t, ok)
	assert.Equal(t, 1000, v)

	// 遍历 key/value
	ks := make([]interface{}, 0, 1000)
	vs := make([]interface{}, 0, 1000)

	sm.Range(func(k, v interface{}) bool { // 遍历需要通过传递一个函数参数完成
		ks = append(ks, k)
		vs = append(vs, v)
		return true // 返回遍历是否结束, 任意一个迭代返回 false, 则整个遍历结束
	})
	assert.Contains(t, ks, 999)
	assert.Contains(t, ks, "0")

	assert.Contains(t, vs, 1000)
	assert.Contains(t, vs, "1")
}
