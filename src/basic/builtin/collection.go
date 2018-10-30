package builtin

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
