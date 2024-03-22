package sets

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试 Set 集合的创建和添加元素
func TestSetCreateAndAdd(t *testing.T) {
	// 初始化并添加元素
	s1 := New[int](100)

	// 批量添加元素
	s1.Add(1, 2, 3, 4, 2)
	assert.Equal(t, 4, s1.Len()) // 实际添加了 4 个元素, 重复的 2 只存在 1 份

	// 判断集合是否包含指定值
	ok := s1.Contains(1)
	assert.True(t, ok)

	ok = s1.Contains(2, 3, 4) // 多值判断
	assert.True(t, ok)

	ok = s1.Contains(3, 4, 5) // 5 不在集合中, 返回 false
	assert.False(t, ok)
}

// 集合相等判断
func TestSetCompare(t *testing.T) {
	s1 := New[int](10)
	s1.Add(1, 2, 3, 4)

	// 产生一个元素相同的集合
	s2 := New[int](10)
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

// 集合子集判断
func TestSetSubset(t *testing.T) {
	s1 := New[int](10)
	s1.Add(1, 2, 3, 4)

	// 产生一个包含 `s1` 集合部分元素的集合
	s2 := New[int](10)
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
