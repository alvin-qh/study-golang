package str

import (
	"io"
	"strings"
	"unicode/utf8"
)

// 从 Reader 中读取指定数量的 byte, 写入 to 参数中
func ReadBytes(r *strings.Reader, to []byte, len int) (int, error) {
	for i := 0; i < len; i++ {
		b, err := r.ReadByte() // io.ByteReader 接口函数, 读取一个字节, 读取指针后移 1
		if err != nil {
			if err == io.EOF {
				return i, nil
			}
			return i, err
		}
		to[i] = b
	}
	return len, nil
}

// 从 Reader 中读取指定数量的字符, 转为 bytes 后写入 to 参数中
func ReadRune(r *strings.Reader, to []byte, len int) (int, error) {
	total := 0
	for i := 0; i < len; i++ {
		c, _, err := r.ReadRune() // io.RuneReader 接口函数, 读取一个 rune 字符
		if err != nil {
			if err == io.EOF {
				return total, nil
			}
			return total, err
		}

		for _, b := range RuneToBytes(c) { // 将 rune 转为 bytes
			to[total] = b // 存入 data 切片的后续位置
			total++
		}
	}
	return total, nil
}

// Rune 转 []byte
// 这里认为 Rune 为 UTF-8 编码, 一个 Rune 可以转为 1~4 个 bytes
func RuneToBytes(r rune) []byte {
	data := make([]byte, 4)       // 可接受最长 utf8 编码的 bytes
	n := utf8.EncodeRune(data, r) // 将字符编码为 utf8 bytes
	return data[:n]               // 返回有效长度的 []byte 切片
}
