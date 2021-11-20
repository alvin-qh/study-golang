package gc

// 内存分配结构体
type Memory struct {
	buf []byte
}

// 分配内存
func (m *Memory) Alloc(n int) {
	m.buf = make([]byte, n, n*2)
}

// 获取所分配内存大小
func (m *Memory) Size() int {
	return len(m.buf)
}

// 释放已分配的内存
func (m *Memory) Clear() {
	m.buf = nil
}
