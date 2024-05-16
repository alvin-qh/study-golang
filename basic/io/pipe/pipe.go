package pipe

import (
	"io"
	"os"
)

// 定义管道对象
type Pipe struct {
	r *os.File
	w *os.File
}

// 创建实例
func New() (*Pipe, error) {
	// 创建管道实例, 返回一对读写文件接口
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	return &Pipe{
		r,
		w,
	}, nil
}

// 获取管道读接口
func (p *Pipe) Reader() *os.File {
	return p.r
}

// 获取管道写接口
func (p *Pipe) Writer() *os.File {
	return p.w
}

// 关闭管道
//
// 同时关闭管道读写接口
func (p *Pipe) Close() error {
	if err := p.CloseWriter(); err != nil {
		return err
	}
	if err := p.CloseReader(); err != nil {
		return err
	}
	return nil
}

// 关闭管道写接口
func (p *Pipe) CloseWriter() error {
	if err := p.w.Close(); err != nil {
		return err
	}
	return nil
}

// 关闭管道读接口
func (p *Pipe) CloseReader() error {
	if err := p.r.Close(); err != nil {
		return err
	}
	return nil
}

// 写入内容
func (p *Pipe) Write(data []byte) (int, error) {
	return p.w.Write(data)
}

// 从指定的 `io.Reader` 实例中读取内容并写入管道
func (p *Pipe) WriteFrom(r io.Reader) error {
	buf := make([]byte, 1024)
	for {
		// 从 `io.Reader` 实例中读取内容
		n, err := r.Read(buf)
		if err != nil {
			// 如果读取到结尾, 则结束读取
			if err == io.EOF {
				break
			}
			// 如果发送其它错误, 则返回错误
			return err
		}

		// 将读取的内容写入管道
		if _, err = p.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
}

// 读取内容
func (p *Pipe) Read(buf []byte) (int, error) {
	return p.r.Read(buf)
}

// 将管道内容读出, 并写入指定的 `io.Writer` 实例
func (p *Pipe) ReadTo(w io.Writer) error {
	buf := make([]byte, 1024)
	for {
		// 从管道中读取内容
		n, err := p.Read(buf)
		if err != nil {
			// 如果读取到结尾, 则结束读取
			if err == io.EOF {
				break
			}
			// 如果发送其它错误, 则返回错误
			return err
		}
		// 将读取内容写入指定的 `io.Writer` 实例
		w.Write(buf[:n])
	}
	return nil
}
