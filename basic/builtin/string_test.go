package builtin

import (
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

// 测试求字符串长度
// 字符串在内存中默认是以 UTF-8 编码存储的 byte 集合
// rune 类型表示一个完整的 UTF-8 字符，可能占据 1~4 个 byte
// 所以获取字符串长度，实际需要得到字符串中 rune 元素的数量
// len(string) 返回的是字符串的 byte 数量，不能作为字符串的真实长度
// len([]rune(string)) 则是字符串真实的字符个数
func TestStringLength(t *testing.T) {
	s := "Hello, 大家好"

	// 将字符串转化为字符数组后求长度
	size := len([]rune(s))
	assert.Equal(t, 10, size)

	// 利用 utf8 包计算字符串长度
	size = utf8.RuneCountInString(s)
	assert.Equal(t, 10, size)

	// 直接对字符串使用 len，得出字符串的 字节总数，对于包含非 ASCII 编码的字符，结果不正确
	size = len(s)
	assert.Equal(t, 7+9, size) // 得出结果为 16，即 7 个 ASCII 字符（占 7 字节）+ 3 个 UTF8 字符（占 9 字节）
}

// 测试字符串转换
// 其它类型和字符串类型之间的转换，主要是通过 strconv 包来完成
func TestConvertToString(t *testing.T) {
	// 整数转换
	n := int64(-100)
	s := strconv.FormatInt(n, 10) // 将 int64 转为 10 进制字符表示
	assert.Equal(t, "-100", s)

	n, ok := strconv.ParseInt(s, 10, 64) // 字符串按 10进制 形式转为 64位 整数
	assert.Nil(t, ok)
	assert.Equal(t, int64(-100), n)

	s = strconv.FormatInt(n, 8) // 将 n 转为 8 进制表示
	assert.Equal(t, "-144", s)

	n, ok = strconv.ParseInt(s, 8, 64) // 将字符串按 8 进制 形式转为 64位 整数
	assert.Nil(t, ok)
	assert.Equal(t, int64(-100), n)

	s = strconv.FormatUint(uint64(n), 10) // 将 n 作为 无符号整型 转为字符串表示
	assert.Equal(t, "18446744073709551516", s)

	un, ok := strconv.ParseUint(s, 10, 64) // 将字符串按 10进 制形式转为 64位 无符号整数
	assert.Nil(t, ok)
	assert.Equal(t, int64(-100), int64(un))

	// 浮点数转换
	f := 123.456
	s = strconv.FormatFloat(f, 'f', -1, 32) // 'f' 表示转为 浮点数值 形式, -1 表示保留所有小数位
	assert.Equal(t, "123.456", s)

	f, ok = strconv.ParseFloat(s, 64) // 将字符串转为 64 位浮点数
	assert.Nil(t, ok)
	assert.Equal(t, float64(123.456), f)

	s = strconv.FormatFloat(f, 'e', -1, 32) // 'e' 表示转为 科学计数法 形式
	assert.Equal(t, "1.23456e+02", s)

	f, ok = strconv.ParseFloat(s, 64) // 将 科学计数法 形式的字符串转为 64 位浮点数
	assert.Nil(t, ok)
	assert.Equal(t, float64(123.456), f)

	s = strconv.FormatFloat(f, 'f', 1, 32) // 1 表示保留 1位 小数位数，并进行四舍五入
	assert.Equal(t, "123.5", s)

	f, ok = strconv.ParseFloat(s, 64) // 将 科学计数法 形式的字符串转为 64 位浮点数
	assert.Nil(t, ok)
	assert.Equal(t, float64(123.5), f)

	s = strconv.FormatFloat(f, 'f', -1, 64) // 64 表示按 float64 长度进行转换
	assert.Equal(t, "123.5", s)

	// 布尔值转换
	b := true
	s = strconv.FormatBool(b) // 将布尔类型值转为字符串
	assert.Equal(t, "true", s)

	b, ok = strconv.ParseBool(s) // 将字符串转化为 布尔类型 值
	assert.Nil(t, ok)
	assert.Equal(t, true, b)

	// 复数类型转换，转换 complex128 到 string
	c := complex(100, 20)
	s = strconv.FormatComplex(c, 'f', -1, 128) // 'f', -1, 128 含义和 FormatFloat 类似
	assert.Equal(t, "(100+20i)", s)

	c, ok = strconv.ParseComplex(s, 128) // 将字符串转化为 128 位复数类型
	assert.Nil(t, ok)
	assert.Equal(t, (100 + 20i), c)
}

// rune 表示一个 ‘字符’ 而非 ‘byte’， 所以要正确的从字符串获取指定位置的字符，需要将字符串类型转为 []rune 来处理
func TestRuneOfString(t *testing.T) {
	s := "Hello, 大家好"
	assert.Equal(t, rune(s[1]), 'e')   // 一个 rune 类型表示一个字符，用单引号 ' 包围
	assert.Equal(t, string(s[1]), "e") // rune 类型可以转为 string 类型

	// 将字符串转为 rune 数组，相当于字符数组
	rs := []rune(s)
	assert.Equal(t, rs[1], int32(s[1]))    // 字符串下标返回一个字节 (uint8)，rune下标返回一个 utf8 字符 (int32)
	assert.NotEqual(t, rs[8], int32(s[8])) // 第 8 个字符是中文字符，所以 string 和 rune 下标相同时，值不再相同
	assert.Equal(t, '家', rs[8])
}

// 字符串比较是通过 strings 包来完成
// strings.EqualFold, strings.Compare
func TestStringCompare(t *testing.T) {
	s := "abc"

	assert.Equal(t, 0, strings.Compare(s, "abc"))  // 两个字符串比较，相等则返回 0
	assert.Equal(t, 1, strings.Compare(s, "Abc"))  // 第一个字符串大于第二个字符串，返回 1
	assert.Equal(t, -1, strings.Compare(s, "bbc")) // 第一个字符串小于第二个字符串，返回 -1

	assert.True(t, strings.EqualFold(s, "ABC")) // 对两个字符串进行忽略大小写的比较，返回是否相等
}

// 对子字符串的操作，是通过 strings 包下面的
// strings.Contains
func TestSubString(t *testing.T) {
	s := "Hello, 大家好"

	// 查看字符串是否包括所给的子字符串
	assert.True(t, strings.Contains(s, ", 大"))
	assert.False(t, strings.Contains(s, "好大"))

	// 查找子字符串在字符串中出现的位置
	n := strings.Index(s, "家好") // 查找子字符串在字符串中第一次出现的位置
	assert.Equal(t, 10, n)      // 在位置 10 找到子字符串

	n = strings.IndexAny(s, "o好") // 查找所给字符中任意字符在字符串中第一次出现的位置
	assert.Equal(t, 4, n)         // 在位置 2 找到子字符串
}

func TestString_Index(t *testing.T) {
	assert.Equal(t, String("abcde").Index("cde"), 2)
	assert.Equal(t, String("abcde").LastIndex("cde"), 2)
}

func TestString_Count(t *testing.T) {
	s := String("abababc")
	assert.Equal(t, s.Count("ab"), 3)
	assert.Equal(t, s.Count(""), s.Len()+1)
}

func TestRepeat(t *testing.T) {
	s := Repeat("abc", 3)
	assert.Equal(t, s.Len(), 9)
	assert.Equal(t, s.Count("abc"), 3)
}

func TestString_Replace(t *testing.T) {
	s := String("Hello")

	s = s.Replace("l", "L", 0) // do nothing replace
	assert.Equal(t, string(s), "Hello")

	s = s.Replace("l", "L", 1)
	assert.Equal(t, string(s), "HeLlo")

	s = String("Hello")
	s = s.Replace("l", "L", 2)
	assert.Equal(t, string(s), "HeLLo")

	s = String("Hello")
	s = s.Replace("ello", "ELLO", -1)
	assert.Equal(t, string(s), "HELLO")

	s = String("Hello")
	s.ReplaceSelf("ello", "ELLO", -1)
	assert.Equal(t, string(s), "HELLO")
}

func TestString_Trim(t *testing.T) {
	s := String("   Hello   ")
	assert.Equal(t, string(s.Trim()), "Hello")
}

func TestString_TrimWith(t *testing.T) {
	s := String("*^%Hello^&*")
	assert.Equal(t, string(s.TrimWith("*^%&")), "Hello")
}

func TestString_StartWith_And_EndWith(t *testing.T) {
	s := String("http://www.google.com")
	assert.True(t, s.StartWith("http://"))
	assert.True(t, s.EndWith(".com"))

	s = s.TrimStart("http://").TrimEnd(".com")
	assert.Equal(t, string(s), "www.google")
}

func TestString_Split(t *testing.T) {
	s := String("www.google.com")
	rs := s.Split(".")
	assert.Equal(t, len(rs), 3)
	assert.Equal(t, string(rs[0]), "www")
	assert.Equal(t, string(rs[1]), "google")
	assert.Equal(t, string(rs[2]), "com")
}

func TestJoin(t *testing.T) {
	parts := make([]string, 26)
	for i := 'A'; i <= 'Z'; i++ {
		parts[i-'A'] = string(i)
	}

	join := Join(",", parts...)
	assert.Equal(t, string(join), "A,B,C,D,E,F,G,H,I,J,K,L,M,N,O,P,Q,R,S,T,U,V,W,X,Y,Z")
}

func TestStringBuilder(t *testing.T) {
	builder := new(StringBuilder)

	builder.Append("Hello")
	builder.Append(" World")

	assert.Equal(t, builder.Size(), 11)
	assert.Equal(t, builder.ToString(), "Hello World")

	builder.AppendInt(123)
	assert.Equal(t, builder.Size(), 14)
	assert.Equal(t, builder.ToString(), "Hello World123")

	builder.AppendFloat(0.1234567)
	assert.Equal(t, builder.Size(), 23)
	assert.Equal(t, builder.ToString(), "Hello World1230.1234567")

	builder.Clear()
	assert.Equal(t, builder.Size(), 0)
	assert.Equal(t, builder.ToString(), "")
}
