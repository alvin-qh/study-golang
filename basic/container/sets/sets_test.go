package sets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试创建 Set 实例
func TestSet_New(t *testing.T) {
	// 初始化并添加元素
	s := New[int]()
	assert.Equal(t, 0, s.Len())
}

// 测试添加元素
//
// 如果向集合中添加重复元素, 则相同的元素只包含一次
func TestSet_Add(t *testing.T) {
	s := New[int]()

	s.Add(1, 2, 3, 4, 2)
	assert.Equal(t, 4, s.Len())
	assert.ElementsMatch(t, []int{1, 2, 3, 4}, s.Slice())
}

// 测试判断元素是否包含在集合中
func TestSet_Contains(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3, 4)

	// 判断元素是否全部包含在集合中
	assert.True(t, s.Contains(1))
	assert.True(t, s.Contains(2, 3, 4))
	assert.False(t, s.Contains(3, 4, 5))
}

// 测试删除元素
func TestSet_Remove(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3, 4)

	// 删除一个元素
	s.Remove(1)
	assert.Equal(t, 3, s.Len())
	assert.False(t, s.Contains(1))

	// 删除元素
	s.Remove(2, 3)
	assert.Equal(t, 1, s.Len())
	assert.False(t, s.Contains(2))
	assert.False(t, s.Contains(3))
}

// 测试集合相等判断
func TestSet_Equal(t *testing.T) {
	s1 := New[int]()
	s1.Add(1, 2, 3, 4)

	// 产生一个元素相同的集合
	s2 := New[int]()
	s2.Add(1, 2, 3, 4)

	// 此时两个集合相等
	assert.True(t, s2.Equal(s1))

	// 在其中一个集合中添加新元素, 此时两个集合不再相同
	s2.Add(5)
	assert.False(t, s2.Equal(s1))

	// 移除之前添加的新元素
	s2.Remove(5)
	assert.True(t, s2.Equal(s1))
}

// 测试判断集合是否为指定集合的子集
func TestSet_IsSubset(t *testing.T) {
	s1 := New[int]()
	s1.Add(1, 2, 3, 4)

	// 产生一个包含 `s1` 集合部分元素的集合
	s2 := New[int]()
	s2.Add(2, 3, 4)

	// 此时 `s2` 为 `s1` 的子集, 但 `s1` 不是 `s2` 的子集
	assert.True(t, s2.IsSubset(s1))
	assert.False(t, s1.IsSubset(s2))

	// 删除 `s1` 集合的部分元素
	s1.Remove(1, 2)

	// 此时 `s2` 不再是 `s1` 的子集, 但 `s1` 成为了 `s2` 的子集
	assert.False(t, s2.IsSubset(s1))
	assert.True(t, s1.IsSubset(s2))
}

// 测试通过回调遍历集合
func TestSet_Do(t *testing.T) {
	s := New[int]()
	s.Add(1, 2, 3, 4)

	vs := make([]int, 0, s.Len())

	// 通过回调遍历集合
	s.Do(func(v int) bool {
		vs = append(vs, v)
		return true
	})

	assert.ElementsMatch(t, []int{1, 2, 3, 4}, vs)
}
