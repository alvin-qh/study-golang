package complexkey

// 定义结构体作为 Map 的 Key
type Key struct {
	Id   int
	Name string
}

// 定义结构体作为 Map 的 Value
type Value struct {
	Gender   rune
	Birthday string
	Address  string
}
