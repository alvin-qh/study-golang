package sets

// 定义 Set 集合结构体
type Set struct {
	m   map[interface{}]struct{}
	cap int
}

// 创建并初始化 Set 集合对象
//
//	cap: Set 集合初始容积
func New(cap int) *Set {
	s := Set{m: make(map[interface{}]struct{}, cap), cap: cap}
	return &s
}

// 向 Set 集合中添加元素
func (s *Set) Add(values ...interface{}) {
	for _, v := range values { // 遍历参数, 获取参数值 (range 返回的第 2 个值, 第 1 个值为下标)
		s.m[v] = struct{}{} // 以 某个参数值 为 key, 空结构为 value, 设置 map (相当于只给 map 设置了 key)
	}
}

// 从 Set 集合中删除指定的元素
func (s *Set) Remove(values ...interface{}) {
	for _, v := range values { // 遍历参数, 从 map 中删除参数所表示的 key
		delete(s.m, v)
	}
}

// 判断 元素 是否在 Set 集合中存在
func (s *Set) Contains(values ...interface{}) bool {
	for _, v := range values { // 遍历参数, 从 map 中查找 参数所表示的 key 是否存在
		if _, ok := s.m[v]; !ok {
			return false
		}
	}
	return true
}

// 获取 Set 集合元素个数
func (s *Set) Len() int {
	return len(s.m)
}

// 情况 Set 集合中的元素
func (s *Set) Clear(i, j int) {
	s.m = make(map[interface{}]struct{}, s.cap)
}

// 判断两个 Set 集合是否相同 (包含相同的元素)
func (s *Set) Equal(other *Set) bool {
	if s.Len() != other.Len() { // 两个集合元素个数是否相同
		return false // 元素个数不同不能相等
	}

	for v := range s.m { // 遍历 map, 判断 参数 所表示的 key 是否存在
		if !other.Contains(v) {
			return false // 某个参数的 key 不存在, 则返回 false
		}
	}
	return true
}

// 判断当前 Set 集合是否另一个集合的 子集
func (s *Set) IsSubset(other *Set) bool {
	if s.Len() > other.Len() {
		return false // 当前集合元素数必须不能大于另一个集合, 否则不能成为 子集
	}

	for v := range s.m {
		if !other.Contains(v) { // 另一个集合是否包含 当前集合的 所有元素
			return false
		}
	}
	return true
}
