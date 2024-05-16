package orderedmap

import (
	"cmp"
	"slices"
)

// 利用 map 和 slice 组合一个 key 有序的 map 类型
//
// 定义结构体
type OrderedMap[K cmp.Ordered, V any] struct {
	m      map[K]V // map 对象
	s      []K     // 保存有序 key 集合的 slice 对象
	sorted bool    // s 字段是否已排序
}

// 创建 `OrderedMap` 对象
func New[K cmp.Ordered, V any]() *OrderedMap[K, V] {
	sm := &OrderedMap[K, V]{}

	// 初始化对象
	sm.Init()
	return sm
}

// 初始化 `OrderedMap` 对象
func (sm *OrderedMap[K, V]) Init() {
	// 创建 map 对象
	sm.m = make(map[K]V)
	// 创建 slice 对象
	sm.s = make([]K, 0)
	// 一开始 key 数量为 0, 所以默认已排序
	sm.sorted = true
}

// 获取存储 key 的个数
func (sm *OrderedMap[K, V]) Len() int {
	return len(sm.s)
}

// 存储一对 key/value
func (sm *OrderedMap[K, V]) Put(key K, value V) {
	// 判断 key 是否已存在, 若 key 不存在, 则额外将 key 在 slice 中存储一份
	if _, ok := sm.m[key]; !ok {
		sm.s = append(sm.s, key)
		// 由于加入了新的 key, 所以之前的排序失效
		sm.sorted = false
	}
	sm.m[key] = value // 设置 value
}

// 根据 key 获取 value
func (sm *OrderedMap[K, V]) Get(key K) (value V, ok bool) {
	// 从 map 集合中获取 value
	value, ok = sm.m[key]
	return
}

// 对 slice 中存储的 key 进行排序
func (sm *OrderedMap[K, V]) sortKey() {
	// 如果已经排序, 则跳过此步骤
	if !sm.sorted {
		// 对 key slice 进行排序
		slices.Sort(sm.s)

		// 设置已排序的标志
		sm.sorted = true
	}
}

// 删除一个 key
func (sm *OrderedMap[K, V]) Remove(key K) {
	// 判断 key 是否存在, 若存在, 则执行删除操作
	if _, ok := sm.m[key]; ok {
		// 从 map 中删除 key
		delete(sm.m, key)

		// 二分查找法, 在 key slice 中找到所需删除 key 的下标
		sm.sortKey()
		if i, ok := slices.BinarySearch(sm.s, key); ok {
			// 删除该元素
			sm.s = slices.Delete(sm.s, i, i+1)
		}
	}
}

// 获取所有的 key
func (sm *OrderedMap[K, V]) Keys() []K {
	sm.sortKey()
	return sm.s
}

// 获取所有的 value
func (sm *OrderedMap[K, V]) Values() []V {
	sm.sortKey()

	vs := make([]V, len(sm.s))
	for i, v := range sm.Keys() {
		vs[i] = sm.m[v]
	}
	return vs
}

// 迭代所有的 key/value
func (sm *OrderedMap[K, V]) Do(r func(key K, val V)) {
	sm.sortKey()
	for _, k := range sm.Keys() { // 变量有序的 key 集合
		r(k, sm.m[k]) // 回调迭代函数
	}
}

// 返回一个迭代函数对集合进行迭代
func (sm *OrderedMap[K, V]) Iterate() func() (K, V, bool) {
	sm.sortKey()

	i := 0

	// 返回迭代函数
	return func() (K, V, bool) {
		if i >= len(sm.s) {
			k := sm.s[len(sm.s)-1]
			return k, sm.m[k], false
		}

		k := sm.s[i]
		i++
		return k, sm.m[k], true
	}
}
