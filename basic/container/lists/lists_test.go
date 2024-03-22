package lists

import (
	"container/list"
	"testing"

	"github.com/stretchr/testify/assert"
)

func _S(a ...any) []any {
	return a
}

// 测试在链表两端添加节点
func TestLinkedListPush(t *testing.T) {
	// 通过 New 创建一个空列表
	lst := list.New()

	// 在列表末尾添加一个元素
	// 返回值是一个 Element 指针, 表示链表的节点
	elem := lst.PushBack(1)
	assert.Equal(t, _S(1), ToSlice[any](lst))
	assert.Equal(t, elem.Value, 1)

	// 在列表开头添加一个元素
	elem = lst.PushFront("Hello")
	assert.Equal(t, _S("Hello", 1), ToSlice[any](lst))
	assert.Equal(t, elem.Value, "Hello")
}

// 测试在链表指定节点前后插入新节点
func TestLinkedListInsert(t *testing.T) {
	lst := FromSlice([]any{"Hello", 1})

	// 找到链表倒数第 2 个节点
	elem := lst.Back().Prev()
	assert.Equal(t, "Hello", elem.Value)

	// 在找到的节点前插入, 返回插入的新节点
	elem = lst.InsertAfter("OK", elem)
	assert.Equal(t, "OK", elem.Value)
	assert.Equal(t, _S("Hello", "OK", 1), ToSlice[any](lst))

	// 找到链表第 2 个节点
	elem = lst.Front().Next()
	assert.Equal(t, "OK", elem.Value)

	// 在找到的节点后插入, 返回插入的新节点
	elem = lst.InsertAfter("Bye", elem)
	assert.Equal(t, "Bye", elem.Value)
	assert.Equal(t, _S("Hello", "OK", "Bye", 1), ToSlice[any](lst))
}

// 测试删除链表指定节点
func TestLinkedListRemove(t *testing.T) {
	lst := FromSlice([]any{"Hello", 1})

	// 找到链表第 2 个节点
	elem := lst.Front().Next()
	assert.Equal(t, 1, elem.Value)

	// 删除找到的节点, 返回节点的 Value
	value := lst.Remove(elem)
	assert.Equal(t, 1, value)
	assert.Equal(t, _S("Hello"), ToSlice[any](lst))
}

// 测试在链表末尾或前端链接另一个链表
func TestLinkedListConcat(t *testing.T) {
	lst := FromSlice([]any{"Hello", 1})

	// 在列表后连接列表
	lst.PushBackList(Reverse(lst))
	assert.Equal(t, _S("Hello", 1, 1, "Hello"), ToSlice[any](lst))

	lst = FromSlice([]any{"Hello", 1})

	// 在列表前连接列表
	lst.PushFrontList(Reverse(lst))
	assert.Equal(t, _S(1, "Hello", "Hello", 1), ToSlice[any](lst))
}

// 测试重新初始化链表
func TestLinkedListClean(t *testing.T) {
	lst := FromSlice([]any{"Hello", 1})
	assert.Equal(t, 2, lst.Len())

	// 重新初始化列表 (清空)
	lst.Init()
	assert.Equal(t, 0, lst.Len())
}
