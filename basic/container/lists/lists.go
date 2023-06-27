package lists

import "container/list"

// 将列表转化为切片
func ToSlice(l *list.List) []interface{} {
	// 生成一个 cap 为列表长度的空切片
	slice := make([]interface{}, 0, l.Len())

	// 遍历列表, 将列表元素依次加入切片中
	for iter := l.Front(); iter != nil; iter = iter.Next() {
		slice = append(slice, iter.Value)
	}
	return slice
}

// 将列表元素反转, 得到新列表
func Reverse(l *list.List) *list.List {
	rl := list.New()
	for iter := l.Back(); iter != nil; iter = iter.Prev() {
		rl.PushBack(iter.Value)
	}
	return rl
}
