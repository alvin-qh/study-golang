// GO 语言集合操作：
//

package builtin

import (
	"container/list"
	"fmt"
)

// 定义 IntArray 类型，为整数数组的别名
type IntArray []int

// 创建一个新 整型 数组
// 	len: 数组的初始长度
//	cap: 数组的最大长度
func NewInts(len int, cap int) IntArray {
	return make(IntArray, len, cap) // make([]int, len, cap)
}

func (i *IntArray) Append(n int) {
	*i = append(*i, n)
}

func (i *IntArray) Remove(n int) {
	*i = append((*i)[:n], (*i)[n+1:]...)
}

func (i *IntArray) Clear() {
	*i = make(IntArray, 0)
}

func (i IntArray) Size() int {
	return len(i)
}

type Any interface{}

func ListAssign(items ...Any) *list.List {
	lst := list.New()
	for _, v := range items {
		lst.PushBack(v)
	}
	return lst
}

func ListAt(list *list.List, index int) Any {
	if index >= 0 && index < list.Len() {
		for iter, i := list.Front(), 0; iter != nil; iter, i = iter.Next(), i+1 {
			if i == index {
				return iter.Value
			}
		}
	}
	panic(fmt.Errorf("out of range [%d, %d]", 0, list.Len()))
}

func ListToSlice(list *list.List) []Any {
	if list == nil {
		panic(fmt.Errorf("invalid list"))
	}

	array := make([]Any, list.Len())
	for iter, i := list.Front(), 0; iter != nil; iter, i = iter.Next(), i+1 {
		array[i] = iter.Value
	}
	return array
}
