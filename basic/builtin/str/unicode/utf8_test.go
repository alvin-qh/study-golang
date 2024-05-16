package unicode

import (
	"slices"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

// 测试获取一个 UTF-8 编码字符的字节长度
func TestUtf8_RuneLen(t *testing.T) {
	// 获取中文字符的字节长度, 为 3 字节
	c := '中'
	assert.Equal(t, 3, utf8.RuneLen(c))
}

// 测试将一个 UTF-8 编码字符编码到字节切片中以及从字节切片中解码一个 UTF-8 字符
func TestUtf8_EncodeRune_DecodeRune(t *testing.T) {
	c := '中'

	// 传教字节切片, 长度为字符的字节长度
	bs := make([]byte, utf8.RuneLen(c))

	// 将字符编码到字节切片中, 返回编码后的字节长度
	size := utf8.EncodeRune(bs, c)
	assert.Equal(t, utf8.RuneLen(c), size)
	assert.Equal(t, []byte{0xe4, 0xb8, 0xad}, bs)

	// 从字节切片中解码一个 UTF-8 字符, 返回字符和解码后的字节长度
	r, size := utf8.DecodeRune(bs)
	assert.Equal(t, c, r)
	assert.Equal(t, utf8.RuneLen(c), size)
}

// 测试计算字符串中 UTF-8 字符的个数
func TestUtf8_RuneCountInString(t *testing.T) {
	s := "Hello,大家好"
	// 确认字符串的字节长度, 即 6 个 ASCII 字符 (占 6 字节) + 3 个 UTF-8 字符 (占 9 字节)
	assert.Len(t, s, 15)

	// 计算字符串的字符个数
	count := utf8.RuneCountInString(s)
	assert.Equal(t, 9, count)

	// 字符串字符个数和字符串转为 []rune 的长度是一致的
	assert.Equal(t, len([]rune(s)), count)
}

// 测试从字符串中解码第一个 UTF-8 字符
func TestUtf8_DecodeRuneInString(t *testing.T) {
	s := "Hello,大家好"
	n := 0

	// 保存解码结果的字符切片
	cs := make([]rune, 0)

	// 循环, 直到完成解码
	for n < len(s) {
		// 从字符串 n 开始的位置解码一个 UTF-8 字符, 返回字符和解码后的字节长度
		r, size := utf8.DecodeRuneInString(s[n:])
		if r == utf8.RuneError {
			break
		}
		cs = append(cs, r)
		n += size
	}

	// 确认解码了所有的字符
	assert.Equal(t, []rune(s), cs)
}

// 测试从字符串中解码最后一个 UTF-8 字符
func TestUtf8_DecodeLastRuneInString(t *testing.T) {
	s := "Hello,大家好"
	n := len(s)

	// 保存解码结果的字符切片
	cs := make([]rune, 0)

	// 循环, 直到完成解码
	for n > 0 {
		// 从字符串 n 开始的位置解码一个 UTF-8 字符, 返回字符和解码后的字节长度
		r, size := utf8.DecodeLastRuneInString(s[:n])
		if r == utf8.RuneError {
			break
		}
		cs = append(cs, r)
		n -= size
	}

	rs := []rune(s)
	slices.Reverse(rs)

	// 确认解码了所有的字符
	assert.Equal(t, rs, cs)
}
