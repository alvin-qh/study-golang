package lists

import "container/list"

// 将切片转化为列表, 返回列表指针
func FromSlice[T any](slice []T) *list.List {
	l := list.New()
	for _, v := range slice {
		l.PushBack(v)
	}
	return l
}

// 将列表转化为切片
func ToSlice[T any](l *list.List) []T {
	// 生成一个 cap 为列表长度的空切片
	slice := make([]T, 0, l.Len())

	// 遍历列表, 将列表元素依次加入切片中
	for iter := l.Front(); iter != nil; iter = iter.Next() {
		v := iter.Value.(T)
		slice = append(slice, v)
	}
	return slice
}

// 将列表元素反转, 返回新列表
func Reverse(l *list.List) *list.List {
	rl := list.New()
	for iter := l.Back(); iter != nil; iter = iter.Prev() {
		rl.PushBack(iter.Value)
	}
	return rl
}
