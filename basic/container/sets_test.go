package container

import (
	"basic/container/sets"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试 Set 集合
func TestSet(t *testing.T) {
	// 初始化并添加元素
	s1 := sets.New(100)          // 初始化
	s1.Add(1, 2, 3, 4, 2)        // 批量添加元素
	assert.Equal(t, 4, s1.Len()) // 实际添加了 4 个元素，重复的 2 只存在 1 份

	// 判断集合是否包含指定值
	ok := s1.Contains(1)
	assert.True(t, ok)

	ok = s1.Contains(2, 3, 4) // 多值判断
	assert.True(t, ok)

	ok = s1.Contains(3, 4, 5) // 5 不在集合中，返回 false
	assert.False(t, ok)

	// 集合相等判断
	s2 := sets.New(10) // 产生一个元素相同的集合
	s2.Add(1, 2, 3, 4)

	ok = s2.Equal(s1) // 两个集合元素是否相同
	assert.True(t, ok)

	s2.Add("Hello")   // 在其中一个集合中添加新元素
	ok = s2.Equal(s1) // 此时两个集合不再相同
	assert.False(t, ok)

	s2.Remove("Hello") // 移除之前添加的新元素
	ok = s2.Equal(s1)  // 此时两个集合恢复相同
	assert.True(t, ok)

	// 判断是否为子集
	ok = s2.IsSubset(s1) // 判断两个相同的集合是否互为子集
	assert.True(t, ok)
	ok = s1.IsSubset(s2)
	assert.True(t, ok)

	s1.Remove(2) // 删除集合元素，此时两个集合不再互为子集
	ok = s2.IsSubset(s1)
	assert.False(t, ok)
	ok = s1.IsSubset(s2)
	assert.True(t, ok)
}
