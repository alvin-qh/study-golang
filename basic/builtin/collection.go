package builtin

import (
	"container/list"
	"fmt"
)

type Ints []int

func NewInts(len int, cap int) Ints {
	return make(Ints, len, cap) // make([]int, len, cap)
}

func (i *Ints) Append(n int) {
	*i = append(*i, n)
}

func (i *Ints) Remove(n int) {
	*i = append((*i)[:n], (*i)[n+1:]...)
}

func (i *Ints) Clear() {
	*i = make(Ints, 0)
}

func (i Ints) Size() int {
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
