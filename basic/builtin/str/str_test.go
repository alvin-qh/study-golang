package str

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

// 测试字符类型
//
// # Go 语言的字符串是一组 `rune` 类型组成的串, `rune` 类型表示一个 UTF-8 类型字符
//
// 通过如下方法可以生成字符串对象
//
//   - 通过字节串生成字符串, 要求字节串为 UTF-8 编码
//   - 通过 `rune` 类型串生成字符串
//
// 例如:
//
//	bs := []byte("Hello World") // 将字符串转为字节串
//	s := string(bs) // 将字节串转为字符串
//
//	rs := []rune{'a', 'b', 'c'} // 将 rune 串转为字符串
//	s = string(rs) // 将字符串转为 rune 串
//
// 另外, 可以通过 `unicode/utf8` 包提供的函数库, 对 `rune` 类型数据进行操作, 包括:
//
//	c := '好'
//	n := utf8.RuneLen(s) // 计算一个字符编码后的字节长度
//	bs := make([]byte, n)
//	utf8.EncodeRune(bs, c) // 将字符 c 编码为字节, 存储到 bs 数组中
//
//	c, n = utf8.DecodeRune(bs) // 从字节串中解码第一个 rune 字符
//
//	bs = []byte("你好")
//	c, n = utf8.DecodeLastRune(bs) // 从字节串中解码最后一个 rune 字符
//
// 进一步的, 可以通过 `unicode/utf8` 包提供的函数库, 对字符串进行操作, 包括:
//
//	c, n := utf8.DecodeRuneInString("Hello") // 从字符串中解码最后一个 rune 字符
//	c, n := utf8.DecodeLastRuneInString("Hello") // 从字符串中解码最后一个 rune 字符
func TestRune(t *testing.T) {
	// 定义字符类型
	r := 'H'
	assert.Equal(t, "int32", reflect.TypeOf(r).Name()) // 字符类型的值实际上是 int32 类型
	assert.Equal(t, int32(72), r)                      // 字符类型的值是 utf8 编码的 72

	// 将 rune 数组生成字符串
	rs := []rune{'A', 'B', 'C'}
	s := string(rs)
	assert.Equal(t, "ABC", s)

	// 将字符串和 utf-8 字节进行转换
	bs := []byte("Hello,大家好") // 字符串转为 byte 数组 (utf8编码)
	s = string(bs)            // byte 数组转换为 字符串
	assert.Equal(t, "Hello,大家好", s)
	assert.Equal(t, []byte(s), bs)

	// 字符和 byte 类型转换
	// 一个 rune 类型可能会转换为 1~4 个 byte 类型 (utf8编码)
	// 实际情况中, 需要从 byte 数组中解码指定的字符, 或将指定的字符编码为 byte 数组
	r = '好'
	bs = make([]byte, utf8.RuneLen(r)) // 获取 rune 类型值转为 byte 数组所需的空间大小, 依据此大小生成 byte 数组作为缓存
	size := utf8.EncodeRune(bs, r)     // 将 rune 编码为 byte 数组, 返回编码长度, 该长度和 utf8.RuneLen 返回的值一致
	assert.Equal(t, 3, size)           // 汉字编码为 byte 数组需 3 个字节
	assert.Equal(t, size, utf8.RuneLen(r))

	r, size = utf8.DecodeRune(bs) // 从 byte 数组解码出第一个字符, 并返回解码了多少个字节
	assert.Equal(t, 3, size)      // 解码了 3 个字节
	assert.Equal(t, '好', r)       // 解码的第一个字符

	s = "好"
	r, size = utf8.DecodeRuneInString(s) // 从字符串解码出第一个字符, 并返回解码了多少个字节
	assert.Equal(t, 3, size)             // 解码了 3 个字节
	assert.Equal(t, '好', r)              // 解码的第一个字符

	r, size = utf8.DecodeLastRune(bs) // 从 byte 数组解码出最后一个字符, 并返回解码了多少个字节
	assert.Equal(t, '好', r)
	assert.Equal(t, 3, size)

	r, size = utf8.DecodeLastRuneInString(s) // 从 byte 数组解码出最后一个字符, 并返回解码了多少个字节
	assert.Equal(t, '好', r)
	assert.Equal(t, 3, size)
}

// 测试字符串长度
//
// 字符串在内存中默认是以 UTF-8 编码存储的字节串
//
// `rune` 类型表示一个完整的 UTF-8 字符, 可能占据 1~4 个字节
//
// 获取字符串长度, 包括如下两种方式:
//   - `len(string)`: 返回的是字符串的字节数量, 对于非 ASCII 编码的字符串, 不能作为字符串的真实长度
//   - `len([]rune(string))` 为字符串 UTF-8 字符个数, 即字符串的真正长度
func TestStringLength(t *testing.T) {
	s := "Hello, 大家好"

	// 将字符串转化为字符数组后求长度
	size := len([]rune(s))
	assert.Equal(t, 10, size)

	// 利用 utf8 包计算字符串长度
	size = utf8.RuneCountInString(s)
	assert.Equal(t, 10, size)

	// 直接对字符串使用 len, 得出字符串的 字节总数, 对于包含非 ASCII 编码的字符, 结果不正确
	size = len(s)
	assert.Equal(t, 7+9, size) // 得出结果为 16, 即 7 个 ASCII 字符 (占 7 字节) + 3 个 UTF8 字符 (占 9 字节)
}

// 测试字符串转换
//
// 其它类型和字符串类型之间的转换, 主要是通过 `strconv` 包来完成
func TestConvertToString(t *testing.T) {
	n := int64(-100)

	// 将整数转为字符串, 以十进制表示
	s := strconv.FormatInt(n, 10)
	assert.Equal(t, "-100", s)

	// 将字符串转为 64 位整数, 按十进制处理
	n, err := strconv.ParseInt(s, 10, 64)
	assert.NoError(t, err)
	assert.Equal(t, int64(-100), n)

	// 将整数转为字符串, 以八进制表示
	s = strconv.FormatInt(n, 8)
	assert.Equal(t, "-144", s)

	// 将字符串转为 64 位整数, 按八进制处理
	n, err = strconv.ParseInt(s, 8, 64)
	assert.NoError(t, err)
	assert.Equal(t, int64(-100), n)

	// 将无符号整数转为字符串, 以十进制表示
	s = strconv.FormatUint(uint64(n), 10)
	assert.Equal(t, "18446744073709551516", s)

	// 将字符串转为 64 位整数, 按十进制处理
	un, err := strconv.ParseUint(s, 10, 64)
	assert.NoError(t, err)
	assert.Equal(t, int64(-100), int64(un))

	f := 123.456

	// 将浮点数转为字符串, 以小数形式表示 ('f'), 保留所有小数位 (-1), 按 32 位浮点数处理 (32)
	s = strconv.FormatFloat(f, 'f', -1, 32)
	assert.Equal(t, "123.456", s)

	// 将字符串转为浮点数, 按 64 位浮点数处理
	f, err = strconv.ParseFloat(s, 64)
	assert.NoError(t, err)
	assert.Equal(t, float64(123.456), f)

	// 将浮点数转为字符串, 以科学计数法形式表示 ('e'), 保留所有小数位 (-1), 按 32 位浮点数处理 (64)
	s = strconv.FormatFloat(f, 'e', -1, 64)
	assert.Equal(t, "1.23456e+02", s)

	// 将字符串 (科学计数法形式) 转为浮点数, 按 64 位浮点数处理
	f, err = strconv.ParseFloat(s, 64)
	assert.NoError(t, err)
	assert.Equal(t, float64(123.456), f)

	// 将浮点数转为字符串, 以小数形式表示 ('f'), 保留 1 位小数位 (1), 按 32 位浮点数处理 (32)
	s = strconv.FormatFloat(f, 'f', 1, 32)
	assert.Equal(t, "123.5", s)

	// 将字符串转为浮点数, 按 64 位浮点数处理
	f, err = strconv.ParseFloat(s, 64)
	assert.NoError(t, err)
	assert.Equal(t, float64(123.5), f)

	// 将浮点数转为字符串, 以小数形式表示 ('f'), 保留所有小数位 (-1), 按 64 位浮点数处理 (64)
	s = strconv.FormatFloat(f, 'f', -1, 64) // 64 表示按 float64 长度进行转换
	assert.Equal(t, "123.5", s)

	b := true

	// 将布尔类型值转为字符串
	s = strconv.FormatBool(b)
	assert.Equal(t, "true", s)

	// 将字符串转化为布尔类型值
	b, err = strconv.ParseBool(s)
	assert.NoError(t, err)
	assert.Equal(t, true, b)

	// 复数类型转换
	c := complex(100, 20)

	// 将复数转为字符串, 'f', -1, 128 含义和 strconv.FormatFloat 函数类似
	s = strconv.FormatComplex(c, 'f', -1, 128)
	assert.Equal(t, "(100+20i)", s)

	// 将字符串转化为 128 位复数类型
	c, err = strconv.ParseComplex(s, 128)
	assert.NoError(t, err)
	assert.Equal(t, (100 + 20i), c)
}

// 测试字符串和字符
//
// `rune` 表示一个"字符"而非"字节", 所以要正确的从字符串获取指定位置的字符, 需要将字符串类型转为 `[]rune` 类型来处理
func TestRuneOfString(t *testing.T) {
	s := "Hello, 大家好"
	assert.Equal(t, rune(s[1]), 'e')   // 一个 rune 类型表示一个字符, 用单引号 ' 包围
	assert.Equal(t, string(s[1]), "e") // rune 类型可以转为 string 类型

	// 将字符串转为 rune 数组, 相当于字符数组
	rs := []rune(s)
	assert.Equal(t, rs[1], int32(s[1]))    // 字符串下标返回一个字节 (uint8), rune下标返回一个 utf8 字符 (int32)
	assert.NotEqual(t, rs[8], int32(s[8])) // 第 8 个字符是中文字符, 所以 string 和 rune 下标相同时, 值不再相同
	assert.Equal(t, '家', rs[8])
}

// 测试比较字符串
//
// 字符串比较是通过 `strings` 包下的函数来完成, 包括:
//   - `strings.Compare`: 比较两个字符串, 返回整数表示比较结果
//   - `strings.EqualFold`: 忽略大小写比较两个字符串, 返回整数表示比较结果
func TestStringCompare(t *testing.T) {
	s := "abc"

	assert.Equal(t, 0, strings.Compare(s, "abc"))  // 两个字符串比较, 相等则返回 0
	assert.Equal(t, 1, strings.Compare(s, "Abc"))  // 第一个字符串大于第二个字符串, 返回 1
	assert.Equal(t, -1, strings.Compare(s, "bbc")) // 第一个字符串小于第二个字符串, 返回 -1

	assert.True(t, strings.EqualFold(s, "ABC")) // 对两个字符串进行忽略大小写的比较, 返回是否相等
}

// 测试子字符串处理
//
// 所谓"子字符串", 即字符串的一部分, 通过 `strings` 包下的函数来完成, 包括:
//   - `strings.Contains`, 是否包含子字符串
//   - `strings.Index`, `strings.LastIndex`, `strings.IndexAny`, `strings.LastIndexAny`, 查找子字符串在字符串中出现的位置
//   - `strings.IndexByte`, `strings.LastIndexByte`, 查找一个字节在字符串中出现的位置
//   - `strings.IndexFunc`, `strings.LastIndexFunc`, 通过一个回调函数来查找字符串内容, 并返回查找位置
//   - `strings.HasPrefix`, `strings.HasSuffix`, 判断字符串是否包含指定前缀或后缀
//   - `strings.Count`, 统计字符串中指定的子字符出现的次数
//   - `strings.Replace`, `strings.ReplaceAll`, 替换子字符串
//   - `strings.Trim`, `strings.TrimSpace`, `strings.TrimLeft`, `strings.TrimRight`,
//     `strings.TrimLeftFunc`, `strings.TrimRightFunc`, `strings.TrimPrefix`, `strings.TrimSuffix`, 字符串截断
//   - `strings.Split`, `strings.SplitN`, `strings.SplitAfter`, `strings.SplitAfterN`, 字符串切分
func TestSubString(t *testing.T) {
	s := "Hello,大家好"

	// 查看字符串是否包括所给的子字符串
	assert.True(t, strings.Contains(s, ",大"))
	assert.False(t, strings.Contains(s, "好大"))

	// 查找子字符串在字符串中出现的位置
	n := strings.Index(s, "家好") // 查找子字符串在字符串中第一次出现的位置
	assert.Equal(t, 9, n)       // 在位置 10 找到子字符串

	n = strings.LastIndex(s, "家好") // 从字符串末尾开始查找
	assert.Equal(t, 9, n)          // 在位置 10 找到 '家' 字符

	// 查找一组字符中任意字符在字符串中首次出现的位置
	n = strings.IndexAny(s, "o好") // 查找所给字符中任意字符在字符串中第一次出现的位置
	assert.Equal(t, 4, n)         // 在位置 4 找到子字符串

	n = strings.LastIndexAny(s, "o好") // 从字符串末尾开始查找
	assert.Equal(t, 12, n)            // 在位置 13 找到 'o' 字符

	// 查找一个 byte 值在字符串中第一次出现的位置
	n = strings.IndexByte(s, 'l') // 在字符串中查找一个 byte 第一次出现的位置
	assert.Equal(t, 2, n)         // 在位置 2 找到 byte 'o'

	n = strings.LastIndexByte(s, 'l') // 在字符串中查找一个 byte 第一次出现的位置
	assert.Equal(t, 3, n)             // 在位置 3 找到 byte 'o'

	// 通过逐字符回调指定函数, 查找符合结果的内容第一次出现的位置
	n = strings.IndexFunc(s, func(r rune) bool {
		return r == ',' // 判断当前字符是否 ',' 字符
	})
	assert.Equal(t, 5, n)

	n = strings.LastIndexFunc(s, func(r rune) bool {
		return r == ',' // 判断当前字符是否 ',' 字符
	})
	assert.Equal(t, 5, n)

	// 判断字符串是否以指定的子字符串开头 (或结束)
	b := strings.HasPrefix(s, "Hello") // 字符串是否以指定的子字符串开头
	assert.True(t, b)

	b = strings.HasSuffix(s, "家好") // 字符串是否以指定的子字符串结束
	assert.True(t, b)

	// 计算子字符串出现的次数
	n = strings.Count(s, "l")
	assert.Equal(t, 2, n) // 指定的子字符串在源字符串中出现了 2 次

	n = strings.Count(s, "")                        // 计算空字符串出现的次数
	assert.Equal(t, utf8.RuneCountInString(s), n-1) // 空字符串出现的此时相当于字符串长度+1

	// 子字符串替换
	sr := strings.Replace(s, "l", "L", 1) // 将子字符串替换为指定的字符串, 共替换 1 次
	assert.Equal(t, "HeLlo,大家好", sr)

	sr = strings.Replace(s, "l", "L", 2) // 将子字符串替换为指定的字符串, 共替换 2 次
	assert.Equal(t, "HeLLo,大家好", sr)

	sr = strings.Replace(s, "l", "L", -1) // 将子字符串替换为指定的字符串, 共替换 任意 次
	assert.Equal(t, "HeLLo,大家好", sr)

	sr = strings.ReplaceAll(s, "l", "L") // 替换所有的指定子字符串, 相当于 strings.Replace(s, "l", "L", -1)
	assert.Equal(t, "HeLLo,大家好", sr)

	// 去除字符串前后的指定内容
	sr = strings.Trim(s, "He好") // 删除字符串前后的指定字符, 字符集合中的内容会被全部删除
	assert.Equal(t, "llo,大家", sr)

	sr = strings.TrimSpace(" \r" + s + "\t\n") // 去除字符串前后的空白字符, 相当于 strings.Trim(s, " \r\n\t")
	assert.Equal(t, s, sr)

	sr = strings.TrimLeft(s, "He好") // 删除字符串开始位置的指定字符
	assert.Equal(t, "llo,大家好", sr)

	sr = strings.TrimLeftFunc(s, func(r rune) bool { // 根据函数的返回值决定是否要去掉指定字符
		return utf8.RuneLen(r) == 1 // 去除所有 byte 长度 大于 1 的字符
	})
	assert.Equal(t, "大家好", sr)

	sr = strings.TrimRight(s, "He好") // 删除字符串结束位置的指定字符
	assert.Equal(t, "Hello,大家", sr)

	sr = strings.TrimRightFunc(s, func(r rune) bool { // 根据函数的返回值决定是否要去掉指定字符
		return utf8.RuneLen(r) > 1 // 去除所有 byte 长度 等于 1 的字符
	})
	assert.Equal(t, "Hello,", sr)

	sr = strings.TrimPrefix(s, "Hel") // 删除字符串开始位置的指定子字符串, 需要匹配整个子字符串
	assert.Equal(t, "lo,大家好", sr)

	sr = strings.TrimSuffix(s, "家好") // 删除字符串结束位置的指定子字符串, 需要匹配整个子字符串
	assert.Equal(t, "Hello,大", sr)

	// 将字符串分隔成若干子字符串
	ss := strings.Split(s, ",")                   // 通过 ',' 将字符串分割
	assert.Equal(t, []string{"Hello", "大家好"}, ss) // 分割为两个子字符串

	ss = strings.SplitN(s, ",", 1)             // 指定分割结果的数量, 至多将字符串分割为 1 个部分
	assert.Equal(t, []string{"Hello,大家好"}, ss) // 分隔为 1 个子字符串, 相当于不做分割

	ss = strings.SplitN(s, ",", 2)                // 指定分割结果的数量, 至多将字符串分割为 2 个部分
	assert.Equal(t, []string{"Hello", "大家好"}, ss) // 分隔为 2 个子字符串

	ss = strings.SplitN(s, ",", -1)               // 分割为任意部分, 相当于 strings.Split(s)
	assert.Equal(t, []string{"Hello", "大家好"}, ss) // 分割为 2 个子字符串

	ss = strings.SplitAfter(s, ",")                // 分割结果中包含用于分割的字符串本身
	assert.Equal(t, []string{"Hello,", "大家好"}, ss) // 用于分隔的字符串和前一个分割结果合并在一起

	ss = strings.SplitAfterN(s, ",", 1)        // 分割结果中至多包含 1 个子字符串
	assert.Equal(t, []string{"Hello,大家好"}, ss) // 分隔为 1 个子字符串, 相当于不做分割

	ss = strings.SplitAfterN(s, ",", 2)            // 分割结果中至多包含 2 个子字符串
	assert.Equal(t, []string{"Hello,", "大家好"}, ss) // 分隔为 2 个子字符串

	ss = strings.SplitAfterN(s, ",", -1)           // 分割为任意部分, 相当于 strings.SplitAfter(s)
	assert.Equal(t, []string{"Hello,", "大家好"}, ss) // 分割为 2 个子字符串
}

// 字符串连接
//
// 字符串连接, 即将若干个子字符串连接成一个完整的字符串
//
// Go 语言中提供了三种字符串连接方式:
//   - `+`, 通过 `+` 运算符连接字符串, 得到新字符串
//   - `strings.Join`, 将多个字符串进行连接, 组成新字符串
//   - `strings.Repeat`, 将一个字符串重复指定次数, 组成新字符串
func TestStringConcat(t *testing.T) {
	s := "Hello"

	// 字符串连接, 通过 '+' 可以连接两个字符串
	sc := s + ", World"
	assert.Equal(t, "Hello, World", sc)

	// 字符串连接, 将字符串数组 (切片) 通过连接符进行连接
	sc = strings.Join([]string{s, "World"}, " ")
	assert.Equal(t, "Hello World", sc)

	// 重复指定字符串若干次
	sc = strings.Repeat(s, 2)
	assert.Equal(t, "HelloHello", sc)
}

// 通过内存缓冲区连接字符串
//
// 对于复杂的字符串拼接, 可通过 `bytes.Buffer` 类型进行, 使用其如下方法:
//   - `WriteString`, 写入字符串
//   - `WriteRune`, 写入字符
//   - `Write`, 写入字节
//
// 完成内容写入后, 可通过 `String` 方法获取结果字符串
func TestStringBuffer(t *testing.T) {
	buffer := bytes.NewBuffer([]byte{})

	// 写入字符串
	buffer.WriteString("Hello ")
	buffer.WriteString("World")

	assert.Equal(t, 11, buffer.Len())
	assert.Equal(t, "Hello World", buffer.String()) // 获取结果字符串

	// 写入字符
	buffer.WriteRune(' ')
	// 写入字节
	buffer.Write([]byte("ABC"))
	assert.Equal(t, "Hello World ABC", buffer.String())

	// 将其它类型转为字符串后写入
	buffer.WriteString(strconv.FormatInt(123, 10))
	assert.Equal(t, "Hello World ABC123", buffer.String())
}

// 字符串格式化
//
// 通过 `fmt.Sprint`, `fmt.Sprintf`, `fmt.Sprintln` 等函数, 可以将一组参数组成一个字符串;
// 其中, `fmt.Sprintf` 函数可以按所给的字符串格式以及参数, 产生格式化后的字符串
func TestStringFormat(t *testing.T) {
	// 将一组参数组成字符串
	s := fmt.Sprint("Hello", "World", 123)
	assert.Equal(t, "HelloWorld123", s)

	// 将一组参数组成字符串, 参数间用 空格 分隔, 末尾增加换行符
	s = fmt.Sprintln("Hello", "World", 123)
	assert.Equal(t, "Hello World 123\n", s)

	// 根据所给的字符串格式, 生成格式化后的字符串
	s = fmt.Sprintf("%s %s %.2f", "Hello", "World", 123.456)
	assert.Equal(t, "Hello World 123.46", s)
}
