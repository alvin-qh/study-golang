package profile

import (
	"bufio"
	"io"
)

const (
	TIME_LAYOUT_UTC = "2006-01-02T15:04:05.000Z" // 格式化时间，UTC 格式
)

// 记录 Profile 数据的接口
type ProfileRecorder interface {
	start() error // 开始记录
	stop()        // 停止记录
}

// Profile 结构体
type Profile struct {
	recorders []ProfileRecorder // 保存 ProfileRecorder 的 slice
}

// 创建一个 Profile 对象
func NewProfile() *Profile {
	return &Profile{
		recorders: make([]ProfileRecorder, 0),
	}
}

// 开始记录 Profile
func (p *Profile) Start() error {
	for _, r := range p.recorders {
		if err := r.start(); err != nil {
			return err
		}
	}
	return nil
}

// 结束记录 Profile
func (p *Profile) Stop() {
	for _, r := range p.recorders {
		r.stop()
	}
}

// 添加一个 ProfileRecorder 对象
func (p *Profile) Use(recorder ProfileRecorder) {
	p.recorders = append(p.recorders, recorder)
}

// 表示频率的类型
type Frequency int64

// 将 io.Writer 包装为带缓冲的 bufio.Writer
func wrapWriter(w io.Writer) *bufio.Writer {
	if bw, ok := w.(*bufio.Writer); ok {
		return bw
	}
	return bufio.NewWriter(w)
}
