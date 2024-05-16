package unicode

import (
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

var (
	// 定义 Zh 类别, 包含全部汉字编码范围
	// 注意, `unicode.Range16` 和 `unicode.Range32` 中定义的编码范围要按照编码大小为顺序
	Zh = &unicode.RangeTable{
		R16: []unicode.Range16{
			{Lo: 0x2E80, Hi: 0x2EFF, Stride: 1}, // CJK 部首补充
			{Lo: 0x3000, Hi: 0x303F, Stride: 1}, // CJK 标点符号
			{Lo: 0x3400, Hi: 0x4DB5, Stride: 1}, // 标准 CJK 文字
			{Lo: 0x4E00, Hi: 0x9FA5, Stride: 1}, // 标准 CJK 文字
			{Lo: 0x9FA6, Hi: 0x9FBB, Stride: 1}, // 标准 CJK 文字
			{Lo: 0xF900, Hi: 0xFA2D, Stride: 1}, // 标准 CJK 文字
			{Lo: 0xFA30, Hi: 0xFA6A, Stride: 1}, // 标准 CJK 文字
			{Lo: 0xFA70, Hi: 0xFAD9, Stride: 1}, // 标准 CJK 文字
			{Lo: 0xFF00, Hi: 0xFFEF, Stride: 1}, // 全角 ASCII, 全角中英文标点, 半宽片假名, 半宽平假名, 半宽韩文字母
		},
		R32: []unicode.Range32{
			{Lo: 0x20000, Hi: 0x2A6D6, Stride: 1}, // 标准 CJK 文字
			{Lo: 0x2F800, Hi: 0x2FA1D, Stride: 1}, // 标准 CJK 文字
		},
	}
)

// 测试 Unicode 编码空间
//
// 通过 `unicode.RangeTable` 可以定义一个编码空间, 其中包含了一系列 Unicode16 编码范围以及一系列 Unicode32 编码范围,
// 通过 `unicode.In` 可以判断一个字符是否在指定的编码空间中
func TestUnicode_RangeTable(t *testing.T) {
	// 判读汉字, 偏旁, 标点符号, 生僻字在 Zh 类别中
	assert.True(t, unicode.In('中', Zh))
	assert.True(t, unicode.In('，', Zh))
	assert.True(t, unicode.In('扌', Zh))
	assert.True(t, unicode.In('壨', Zh))

	// 判读阿拉伯文字不属于 Zh 类别
	assert.False(t, unicode.In('م', Zh))
}

// 测试判断一个字符是否在一组指定的 Unicode 编码空间中
//
// Go 语言预定了一系列“类别”变量, 每个类别都是一个 `unicode.RangeTable` 实例,
// 表示一系列 Unicode 编码范围
func TestUnicode_In(t *testing.T) {
	// 判读字符是否在 L 或 C 类别中 (全部字母和控制字符)
	assert.True(t, unicode.In('a', unicode.L, unicode.C))
	assert.True(t, unicode.In('\n', unicode.L, unicode.C))

	// 判读字符是否在 N 类别中 (全部数字)
	assert.True(t, unicode.In('0', unicode.N))

	// 判读字符是否在 Han 或 Z 类别中 (全部汉字, 空格字符)
	assert.True(t, unicode.In('中', unicode.Han, unicode.Z))
	assert.True(t, unicode.In(' ', unicode.Han, unicode.Z))
}

// 测试判断一个字符是否在特定的 Unicode 编码空间中
func TestUnicode_Is(t *testing.T) {
	// 判读字符是否在 L 类别中 (全部字母)
	assert.True(t, unicode.Is(unicode.L, 'a'))

	// 判读字符是否在 N 类别中 (全部数字)
	assert.True(t, unicode.Is(unicode.N, '0'))

	// 判读字符是否在 Han 类别中 (全部汉字)
	assert.True(t, unicode.Is(unicode.Han, '中'))

	// 判读字符是否在 Z 类别中 (空格字符)
	assert.True(t, unicode.Is(unicode.Z, ' '))
}

// 测试判断一个字符是否为可打印字符
//
// 可打印字符是指可以在屏幕上显示的字符
func TestUnicode_IsPrint(t *testing.T) {
	assert.True(t, unicode.IsPrint('a'))
	assert.True(t, unicode.IsPrint('中'))
	assert.True(t, unicode.IsPrint(' '))

	// 换行符不是可打印字符
	assert.False(t, unicode.IsPrint('\n'))
}

// 测试判断一个字符是否为可见字符
//
// 可见字符是指可以显示并被识别的字符
func TestUnicode_IsGraphic(t *testing.T) {
	assert.True(t, unicode.IsGraphic('a'))
	assert.True(t, unicode.IsGraphic('中'))
	assert.True(t, unicode.IsGraphic(' '))

	// 换行符不是可显示字符
	assert.False(t, unicode.IsGraphic('\n'))
}

// 测试判断一个字符是否为空白字符
//
// 空白字符包括: 空格, 制表符, 回车符, 换行符
func TestUnicode_IsSpace(t *testing.T) {
	assert.True(t, unicode.IsSpace(' '))
	assert.True(t, unicode.IsSpace('\t'))
	assert.True(t, unicode.IsSpace('\n'))

	// 字母字符不是空白字符
	assert.False(t, unicode.IsSpace('a'))
}

// 测试判断一个字符是否为控制字符
//
// 控制字符包括: 换行符, 回车符等用于控制打印机或显示器格式的特殊字符
func TestUnicode_IsControl(t *testing.T) {
	assert.True(t, unicode.IsControl('\n'))
	assert.True(t, unicode.IsControl('\r'))

	// 字母字符不是控制字符
	assert.False(t, unicode.IsControl('a'))
}

// 测试判断一个字符是否为阿拉伯数字字符
func TestUnicode_IsDigit(t *testing.T) {
	assert.True(t, unicode.IsDigit('0'))
	assert.True(t, unicode.IsDigit('9'))

	// 字母字符不是数字
	assert.False(t, unicode.IsDigit('a'))
}

// 测试判断一个字符是否为表示数字的字符
//
// 与 `unicode.IsDigit` 函数不同, 所有表示数字的字符都会令 `unicode.IsNumber` 返回 `true`, 例如 `'Ⅷ'` 字符
func TestUnicode_IsNumber(t *testing.T) {
	assert.True(t, unicode.IsNumber('0'))
	assert.True(t, unicode.IsNumber('9'))
	assert.True(t, unicode.IsNumber('Ⅷ'))

	// 字母字符不是数字
	assert.False(t, unicode.IsNumber('a'))
}

// 测试判断一个字符是否为字母字符
func TestUnicode_IsLetter(t *testing.T) {
	assert.True(t, unicode.IsLetter('a'))
	assert.True(t, unicode.IsLetter('z'))

	// 数字字符不是字母
	assert.False(t, unicode.IsLetter('1'))
}

// 测试判断一个字符是否为小写字母
func TestUnicode_IsLower(t *testing.T) {
	assert.True(t, unicode.IsLower('a'))
	assert.True(t, unicode.IsLower('z'))

	// 'A' 字符不是小写字母
	assert.False(t, unicode.IsLower('A'))
}

// 测试判断一个字符是否为大写字母
func TestUnicode_IsUpper(t *testing.T) {
	assert.True(t, unicode.IsUpper('A'))
	assert.True(t, unicode.IsUpper('Z'))

	// 'a' 字符不是大写字母
	assert.False(t, unicode.IsUpper('a'))
}

// 判断一个字符是否为拉丁文标题大小写字符
//
// 标题大小写字符表参见: https://www.compart.com/en/unicode/category/Lt
func TestUnicode_IsTitle(t *testing.T) {
	assert.True(t, unicode.IsTitle('ǅ'))
	assert.True(t, unicode.IsTitle('ǈ'))
	assert.True(t, unicode.IsTitle('ῌ'))

	// 'A' 字符不是标题大小写字母
	assert.False(t, unicode.IsTitle('A'))
}

// 测试判断一个字符是否为标点符号
func TestUnicode_IsPunct(t *testing.T) {
	assert.True(t, unicode.IsPunct(','))
	assert.True(t, unicode.IsPunct('.'))

	// 字母字符不是标点符号
	assert.False(t, unicode.IsPunct('a'))
}

// 测试判断一个字符是否为符号字符
func TestUnicode_IsSymbol(t *testing.T) {
	assert.True(t, unicode.IsSymbol('©'))
	assert.True(t, unicode.IsSymbol('™'))

	// 字母字符不是标点符号
	assert.False(t, unicode.IsSymbol('a'))
}

// 测试判断一个字符是否为标记字符
//
// 标记字符一般用于标记一些特殊语言的音调, 例如 `"g̀"` 包括了两个字符 `g` 和 `̀̀ `, 其中 `̀̀ ` 表示音调标记字符
func TestUnicode_IsMark(t *testing.T) {
	// 两个字符 `g` 和 `̀̀ `, 其中 `̀̀ ` 表示音调标记字符
	s := []rune("g̀")
	assert.False(t, unicode.IsMark(s[0]))
	assert.True(t, unicode.IsMark(s[1]))

	// 其它符号不是标记字符
	assert.False(t, unicode.IsMark('℃'))
}

// 测试将一个字符转换为其它字符
//
// 转换策略包括: `unicode.LowerCase`, `unicode.UpperCase`, `unicode.TitleCase` 以及 `unicode.MaxCase`
func TestUnicode_To(t *testing.T) {
	// 转换为小写字母
	assert.Equal(t, 'a', unicode.To(unicode.LowerCase, 'A'))

	// 转换为大写字母
	assert.Equal(t, 'A', unicode.To(unicode.UpperCase, 'a'))

	// 转换为标题大小写字母
	assert.Equal(t, 'A', unicode.To(unicode.TitleCase, 'a'))
}

// 测试将一个字符转换为其它字符
//
// 通过 `unicode.ToLower`, `unicode.ToUpper`, `unicode.ToTitle` 进行转换
func TestUnicode_ToLower_ToUpper_ToTitle(t *testing.T) {
	// 转换为小写字母
	assert.Equal(t, 'a', unicode.ToLower('A'))

	// 转换为大写字母
	assert.Equal(t, 'A', unicode.ToUpper('a'))

	// 转换为标题大小写字母
	assert.Equal(t, 'A', unicode.ToTitle('a'))
}
