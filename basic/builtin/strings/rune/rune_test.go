package rune_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试字符类型定义
//
// 在 Go 语言中, `rune` 类型表示一个 UTF-8 编码的字符
func TestRune_TypeDefine(t *testing.T) {
	// 定义 rune 类型变量
	r := 'H'

	// 字符类型的值实际上是 int32 类型
	assert.Equal(t, "int32", reflect.TypeOf(r).Name())

	// 字符类型的值是 utf8 编码的 72
	assert.Equal(t, int32(72), r)
}

// 测试字符串转化为 `[]rune` 切片
//
// 在 Go 语言中, 字符串可以转化为 `[]rune` 切片, 以便操作字符, 对于一个字符串变量 `s`, 有
//
//	s := "Hello, 大家好"
//	rs := []rune(s)
//
// 则
//
//   - `len(s)` 和 `len(rs)` 不一定相等, 前者是字符串的字节长度, 后者是字符串的 UTF-8 字符数量;
//   - `rs[n]` 表示获取 `n` 位置的字符, `s[n]` 则表示获取 `n` 位置的字节
func TestRune_Slice(t *testing.T) {
	// 定义字符串并确认其长度
	s := "Hello, 大家好"
	assert.Len(t, s, 16)

	// 将字符串转为字符切片并确认切片长度
	rs := []rune(s)
	assert.Len(t, rs, 10)
	assert.Equal(t, []rune{'H', 'e', 'l', 'l', 'o', ',', ' ', '大', '家', '好'}, rs)

	// 定义字符切片
	rs = []rune{'H', 'e', 'l', 'l', 'o', ',', ' ', '大', '家', '好'}
	assert.Len(t, s, 16)

	// 将字符切片转为字符串实例
	s = string(rs)
	assert.Len(t, rs, 10)
	assert.Equal(t, "Hello, 大家好", s)

	// 获取字符串下标为 8 的字符
	c1 := rs[8]
	// 获取字符串下标为 8 的字节
	c2 := s[8]
	assert.Equal(t, '家', c1)
	assert.Equal(t, int32(164), rune(c2))
}
