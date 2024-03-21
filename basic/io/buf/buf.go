package buf

import (
	"bufio"
	"io"
)

// 从输入中读取所有的行
//
// 返回行数组
func ReadLines(r io.Reader) ([]string, error) {
	br, ok := r.(*bufio.Reader)
	if !ok {
		br = bufio.NewReader(r)
	}

	lines := make([]string, 0, 100)

	for {
		line, err := br.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		lines = append(lines, line[:len(line)-1])
	}

	return lines, nil
}
