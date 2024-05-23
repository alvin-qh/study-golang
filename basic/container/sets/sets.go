package sets

type Nothing struct{}

var (
	Empty Nothing = struct{}{}
)

// 定义 Set 集合结构体
//
// Go 语言本身不提供 Set 集合, 需要用 Map 模拟, 基本思路为:
//   - 定义一个 `Map` 集合, 并将 Map 的 Value 类型定义为 `Nothing`
//   - 将 `Map` 集合的 Key 作为 Set 集合的元素
//
// `Nothing` (即 `struct{}`) 类型相当于一个空类型, 不占用实际的存储空间
type Set[T comparable] struct {
	m map[T]Nothing
}

// 创建并初始化 Set 集合对象
func New[T comparable]() *Set[T] {
	s := new(Set[T])
	s.Init()
	return s
}

// 初始化 Set 集合
func (s *Set[T]) Init() {
	s.m = make(map[T]Nothing)
}

// 向 Set 集合中添加元素
func (s *Set[T]) Add(values ...T) {
	// 遍历参数, 获取参数值 (range 返回的第 2 个值, 第 1 个值为下标)
	for _, v := range values {
		// 以 某个参数值 为 key, 空结构为 value, 设置 map (相当于只给 map 设置了 key)
		s.m[v] = Empty
	}
}

// 从集合中删除指定的元素
func (s *Set[T]) Remove(values ...T) {
	// 遍历参数, 从 map 中删除参数所表示的 key
	for _, v := range values {
		delete(s.m, v)
	}
}

// 判断元素是否在集合中存在
func (s *Set[T]) Contains(values ...T) bool {
	// 遍历参数, 从 map 中查找 参数所表示的 key 是否存在
	for _, v := range values {
		if _, ok := s.m[v]; !ok {
			return false
		}
	}
	return true
}

// 获取 Set 集合元素个数
func (s *Set[T]) Len() int {
	return len(s.m)
}

// 判断两个 Set 集合是否相同 (包含相同的元素)
func (s *Set[T]) Equal(other *Set[T]) bool {
	if s.Len() != other.Len() {
		// 元素个数不同不能相等
		return false
	}

	// 遍历 map, 判断 参数 所表示的 key 是否存在
	for v := range s.m {
		if !other.Contains(v) {
			// 某个参数的 key 不存在, 则返回 false
			return false
		}
	}
	return true
}

// 判断当前 Set 集合是否另一个集合的 子集
func (s *Set[T]) IsSubset(other *Set[T]) bool {
	if s.Len() > other.Len() {
		// 当前集合元素数必须不能大于另一个集合, 否则不能成为 子集
		return false
	}

	for v := range s.m {
        // 另一个集合是否包含 当前集合的 所有元素
		if !other.Contains(v) {
			return false
		}
	}
	return true
}

// 将 Set 转为切片
func (s *Set[T]) Slice() []T {
	rs := make([]T, 0, s.Len())

	for v := range s.m {
		rs = append(rs, v)
	}
	return rs
}

// 通过回调函数遍历所有元素
func (s *Set[T]) Do(fn func(v T) bool) {
	for v := range s.m {
		if !fn(v) {
			break
		}
	}
}
