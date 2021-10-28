package container

import (
	"basic/container/lists"
	"container/list"
	"testing"

	"github.com/stretchr/testify/assert"
)

// go 语言的 list 实际上是双向链表，适合一些需要插入或删除中间节点的集合操作
// list 也能作为队列或栈来使用
func TestList(t *testing.T) {
	// 通过 New 创建一个空列表
	lst := list.New()
	assert.Equal(t, 0, lst.Len()) // 列表长度目前为 0

	// 添加元素
	// 返回值是一个 Element 指针，表示链表的节点
	elem := lst.PushBack(1) // 在列表末尾添加一个元素
	assert.Equal(t, []interface{}{1}, lists.ToSlice(lst))
	assert.Equal(t, elem.Value, 1) // 返回的 Element 即为刚添加元素的节点

	elem = lst.PushFront("Hello") // 在列表开头添加一个元素
	assert.Equal(t, []interface{}{"Hello", 1}, lists.ToSlice(lst))
	assert.Equal(t, elem.Value, "Hello")

	// 插入元素
	// 插入元素表示将新的元素添加在已有 Element 之前（或之后）
	elem = lst.Back().Prev() // 获取列表末尾元素的前一个元素节点
	assert.Equal(t, "Hello", elem.Value)

	lst.InsertAfter("OK", elem) // 在节点前插入
	assert.Equal(t, []interface{}{"Hello", "OK", 1}, lists.ToSlice(lst))

	elem = lst.Front().Next()
	assert.Equal(t, "OK", elem.Value)

	lst.InsertAfter("Bye", elem)
	assert.Equal(t, []interface{}{"Hello", "OK", "Bye", 1}, lists.ToSlice(lst))

	// 删除元素
	// 删除元素依赖被删除元素的节点对象，所以要先找到这个节点
	elem = lst.Front().Next().Next()
	assert.Equal(t, "Bye", elem.Value)

	value := lst.Remove(elem) // 删除节点，返回节点的 Value
	assert.Equal(t, "Bye", value)
	assert.Equal(t, []interface{}{"Hello", "OK", 1}, lists.ToSlice(lst))

	// 连接两个列表
	lst.PushBackList(lists.Reverse(lst)) // 在列表后连接列表
	assert.Equal(t, []interface{}{"Hello", "OK", 1, 1, "OK", "Hello"}, lists.ToSlice(lst))

	lst.PushFrontList(lists.Reverse(lst)) // 在列表前连接列表
	assert.Equal(t, []interface{}{"Hello", "OK", 1, 1, "OK", "Hello", "Hello", "OK", 1, 1, "OK", "Hello"}, lists.ToSlice(lst))

	// 重新初始化列表（清空）
	lst.Init()
	assert.Equal(t, 0, lst.Len())
}
