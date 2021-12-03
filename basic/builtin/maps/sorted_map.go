package maps

import "sort"

// 利用 map 和 slice 组合一个 key 有序的 map 类型

// 用于比较的函数类型
type ComparatorFunc = func(a, b interface{}) int

// 用于迭代的函数
type RangeFunc = func(key interface{}, val interface{})

// 定义结构体
type SortedMap struct {
	m          map[interface{}]interface{} // map 对象
	s          []interface{}               // 保存有序 key 集合的 slice 对象
	sorted     bool                        // s 字段是否已排序
	comparator ComparatorFunc              // 排序使用的比较函数
}

// 创建 SortedMap 对象
func NewSortedMap(comparator ComparatorFunc) *SortedMap {
	sm := &SortedMap{}
	sm.init(comparator) // 初始化对象
	return sm
}

// 初始化 SortedMap 对象
func (sm *SortedMap) init(comparator ComparatorFunc) {
	sm.m = make(map[interface{}]interface{}) // 创建 map 对象
	sm.s = make([]interface{}, 0)            // 创建 slice 对象
	sm.sorted = true                         // 一开始 key 数量为 0，所以默认已排序
	sm.comparator = comparator               // 保存比较函数
}

// 获取存储 key 的个数
func (sm *SortedMap) Len() int {
	return len(sm.s)
}

// 存储一对 key/value
func (sm *SortedMap) Put(key interface{}, value interface{}) {
	if _, ok := sm.m[key]; !ok { // 判断 key 是否已存在，若 key 不存在，则额外将 key 在 slice 中存储一份
		sm.s = append(sm.s, key)
		sm.sorted = false // 由于加入了新的 key，所以之前的排序失效
	}
	sm.m[key] = value // 设置 value
}

// 根据 key 获取 value
func (sm *SortedMap) Get(key interface{}) (value interface{}, ok bool) {
	value, ok = sm.m[key] // 从 map 集合中获取 value
	return
}

// 对 slice 中存储的 key 进行排序
func (sm *SortedMap) sortKey() {
	// 如果已经排序，则跳过此步骤
	if !sm.sorted {
		// 对 key slice 进行排序
		sort.Slice(sm.s, func(i, j int) bool { return sm.comparator(sm.s[i], sm.s[j]) >= 0 })
		sm.sorted = true // 设置已排序的标志
	}
}

// 删除一个 key
func (sm *SortedMap) Remove(key interface{}) {
	if _, ok := sm.m[key]; ok { // 判断 key 是否存在，若存在，则执行删除操作
		delete(sm.m, key) // 从 map 中删除 key

		sm.sortKey() // 对 key slice 进行排序

		// 二分查找法，在 key slice 中找到所需删除 key 的下标
		i := sort.Search(len(sm.s), func(n int) bool { return sm.comparator(sm.s[n], key) <= 0 })
		// 重建 slice，忽略要删除的下标
		sm.s = append(sm.s[:i], sm.s[i+1:]...)
	}
}

// 获取所有的 key
func (sm *SortedMap) Keys() []interface{} {
	if !sm.sorted {
		sm.sortKey()
	}
	return sm.s
}

// 获取所有的 value
func (sm *SortedMap) Values() []interface{} {
	vs := make([]interface{}, len(sm.s))
	for i, v := range sm.Keys() {
		vs[i] = sm.m[v]
	}
	return vs
}

// 迭代所有的 key/value
func (sm *SortedMap) Range(r RangeFunc) {
	for _, k := range sm.Keys() { // 变量有序的 key 集合
		r(k, sm.m[k]) // 回调迭代函数
	}
}

// 清空当前对象
func (sm *SortedMap) Clear() {
	sm.init(sm.comparator) // 重新初始化
}
