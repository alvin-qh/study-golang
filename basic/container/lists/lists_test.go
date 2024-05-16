package lists

import (
	"container/list"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试创建链表实例
func TestList_New(t *testing.T) {
	// 通过 New 创建一个空列表
	l := list.New()

	// 此时链表为空
	assert.Equal(t, 0, l.Len())
}

// 测试在列表末尾添加一个元素
func TestList_PushBack(t *testing.T) {
	l := list.New()

	// 在列表末尾添加一个元素
	elem := l.PushBack(1)
	assert.Equal(t, elem.Value, 1)

	elem = l.PushBack(2)
	assert.Equal(t, elem.Value, 2)

	// 此时链表长度为 2
	assert.Equal(t, 2, l.Len())
	assert.Equal(t, []any{1, 2}, ToSlice[any](l))
}

// 测试在列表开头添加一个元素
func TestList_PushFront(t *testing.T) {
	l := list.New()

	// 在列表开头添加一个元素
	elem := l.PushFront(1)
	assert.Equal(t, elem.Value, 1)

	elem = l.PushFront(2)
	assert.Equal(t, elem.Value, 2)

	// 此时链表长度为 2
	assert.Equal(t, 2, l.Len())
	assert.Equal(t, []any{2, 1}, ToSlice[any](l))
}

// 测试重新初始化链表
//
// 该方法可以用于清空链表
func TestList_Init(t *testing.T) {
	l := list.New()

	l.PushBack(1)
	assert.Equal(t, 1, l.Len())

	// 重新初始化列表 (清空)
	l.Init()
	assert.Equal(t, 0, l.Len())
}

// 测试获取链表首元素
func TestList_First(t *testing.T) {
	l := FromSlice([]any{"Hello", 1})

	elem := l.Front()
	assert.Equal(t, "Hello", elem.Value)
}

// 测试获取链表末尾元素
func TestList_Back(t *testing.T) {
	l := FromSlice([]any{"Hello", 1})

	elem := l.Back()
	assert.Equal(t, 1, elem.Value)
}

// 测试获取链表指定节点的前一个节点
func TestList_Prev(t *testing.T) {
	l := FromSlice([]any{"Hello", 1})

	// 获取链表最后一个节点
	elem := l.Back()

	// 获取节点前一个节点
	elem = elem.Prev()
	assert.Equal(t, "Hello", elem.Value)
}

// 测试获取链表指定节点的后一个节点
func TestList_Next(t *testing.T) {
	l := FromSlice([]any{"Hello", 1})

	// 获取链表第一个节点
	elem := l.Front()

	// 获取节点后一个节点
	elem = elem.Next()
	assert.Equal(t, 1, elem.Value)
}

// 测试在指定节点前插入新节点
func TestList_InsertBefore(t *testing.T) {
	l := FromSlice([]any{"Hello", 1})

	// 找到链表倒数第 2 个节点
	elem := l.Back().Prev()
	assert.Equal(t, "Hello", elem.Value)

	// 在找到的节点前插入, 返回插入的新节点
	elem = l.InsertAfter("OK", elem)
	assert.Equal(t, "OK", elem.Value)
	assert.Equal(t, []any{"Hello", "OK", 1}, ToSlice[any](l))
}

// 测试在指定节点后插入新节点
func TestList_InsertAfter(t *testing.T) {
	l := FromSlice([]any{"Hello", 1})

	// 找到链表第 2 个节点
	elem := l.Front().Next()
	assert.Equal(t, 1, elem.Value)

	// 在找到的节点后插入, 返回插入的新节点
	elem = l.InsertAfter("Bye", elem)
	assert.Equal(t, "Bye", elem.Value)
	assert.Equal(t, []any{"Hello", 1, "Bye"}, ToSlice[any](l))
}

// 测试将节点移动到指定节点前
func TestList_MoveBefore(t *testing.T) {
	l := FromSlice([]any{1, 2, 3, 4})

	// 找到链表第 2 个节点
	elem := l.Front().Next()
	assert.Equal(t, 2, elem.Value)

	// 将节点移动到链表末尾元素之前
	l.MoveBefore(elem, l.Back())
	assert.Equal(t, []any{1, 3, 2, 4}, ToSlice[any](l))

	// 将节点移动到首节点之前
	l.MoveBefore(elem, l.Front())
	assert.Equal(t, []any{2, 1, 3, 4}, ToSlice[any](l))
}

// 测试将节点移动到指定节点之后
func TestList_MoveAfter(t *testing.T) {
	l := FromSlice([]any{1, 2, 3, 4})

	// 找到链表第 2 个节点
	elem := l.Front().Next()
	assert.Equal(t, 2, elem.Value)

	// 将节点移动到链表末尾
	l.MoveAfter(elem, l.Back())
	assert.Equal(t, []any{1, 3, 4, 2}, ToSlice[any](l))

	// 将节点移动到首节点之后
	l.MoveAfter(elem, l.Front())
	assert.Equal(t, []any{1, 2, 3, 4}, ToSlice[any](l))
}

// 测试将指定节点移动到链表末尾
func TestList_MoveToBack(t *testing.T) {
	l := FromSlice([]any{1, 2, 3, 4})

	// 找到链表第 2 个节点
	elem := l.Front().Next()
	assert.Equal(t, 2, elem.Value)

	// 将节点移动到链表末尾
	l.MoveToBack(elem)
	assert.Equal(t, []any{1, 3, 4, 2}, ToSlice[any](l))
}

// 测试将指定节点移动到链表起始
func TestList_MoveToFront(t *testing.T) {
	l := FromSlice([]any{1, 2, 3, 4})

	// 找到链表第 2 个节点
	elem := l.Front().Next()
	assert.Equal(t, 2, elem.Value)

	// 将节点移动到链表首
	l.MoveToFront(elem)
	assert.Equal(t, []any{2, 1, 3, 4}, ToSlice[any](l))
}

// 测试删除链表指定节点
func TestList_Remove(t *testing.T) {
	l := FromSlice([]any{"Hello", 1})

	// 找到链表第 2 个节点
	elem := l.Front().Next()
	assert.Equal(t, 1, elem.Value)

	// 删除找到的节点, 返回节点的 Value
	value := l.Remove(elem)
	assert.Equal(t, 1, value)
	assert.Equal(t, []any{"Hello"}, ToSlice[any](l))
}

// 测试在链表末尾添加另一个链表的所有元素
func TestList_PushBackList(t *testing.T) {
	l := FromSlice([]any{"Hello", 1})

	// 反转原链表, 得到新链表
	lr := Reverse(l)

	// 在链表的末尾添加新链表
	l.PushBackList(lr)
	assert.Equal(t, []any{"Hello", 1, 1, "Hello"}, ToSlice[any](l))
}

// 测试在链表前端添加另一个链表的所有元素
func TestList_PushFrontList(t *testing.T) {
	l := FromSlice([]any{"Hello", 1})

	// 反转原链表, 得到新链表
	lr := Reverse(l)

	// 在链表的前端添加新链表
	l.PushFrontList(lr)
	assert.Equal(t, []any{1, "Hello", "Hello", 1}, ToSlice[any](l))
}
