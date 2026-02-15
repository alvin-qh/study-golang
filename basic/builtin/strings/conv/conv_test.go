package conv_test

import (
	"strconv"
	"study/basic/builtin/strings/conv"
	"study/basic/builtin/strings/utils"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/assert"
)

// 测试将值类型转为字符串
func TestStrConv_Format(t *testing.T) {
	// 将整数以所给进制转为字符串
	// 需要指定结果的进制 (`base` 参数), 可以为 `2`, `4`, `8`, `10`, `16`, `32`, `64` 等
	t.Run("strconv.FormatInt", func(t *testing.T) {
		// 以 10 进制转换为字符串
		s := strconv.FormatInt(100, 10)
		assert.Equal(t, "100", s)

		// 以 16 进制转换为字符串
		s = strconv.FormatInt(100, 16)
		assert.Equal(t, "64", s)

		// 以 8 进制转换为字符串
		s = strconv.FormatInt(100, 8)
		assert.Equal(t, "144", s)

		// 以 2 进制转换为字符串
		s = strconv.FormatInt(100, 2)
		assert.Equal(t, "1100100", s)
	})

	// 测试将无符号整数以所给进制转为字符串
	// 需要指定结果的进制 (`base` 参数), 可以位 `2`, `4`, `8`, `10`, `16`, `32`, `64` 等
	t.Run("strconv.FormatUint", func(t *testing.T) {
		// 以 10 进制转换为字符串
		s := strconv.FormatUint(100, 10)
		assert.Equal(t, "100", s)

		// 以 16 进制转换为字符串
		s = strconv.FormatUint(100, 16)
		assert.Equal(t, "64", s)

		// 以 8 进制转换为字符串
		s = strconv.FormatUint(100, 8)
		assert.Equal(t, "144", s)

		// 以 2 进制转换为字符串
		s = strconv.FormatUint(100, 2)
		assert.Equal(t, "1100100", s)
	})

	// 将浮点数以所给进制转为字符串
	// 需要指定结果的格式 (`fmt` 参数), 可以为:
	//   - 'b' (`-dddp±dddd`, 以 2 为低的指数)
	//   - 'e' (`-d.dddde±dd`, 以 10 为底的指数)
	//   - 'E' (`-d.ddddE±dd`, 同上, `E` 字母大写)
	//   - 'f' (`-ddd.dddd`，无指数)
	//   - 'g' (根据小数位数自动选择 'e' 或 'f')
	//   - 'G' (根据小数位数自动选择 'E' 或 'f')
	//   - 'x' (`-0xd.ddddp±dddd`, 小数部分为 16 进制, 指数部分以 2 为底, 'p' 字母小写)
	//   - 'X' (`-0Xd.ddddP±dddd`, 小数部分为 16 进制, 指数部分以 2 为底, 'P' 字母大写)
	// 需要指定保留小数位数 (`prec` 参数)
	// 需要指定二进制位数, 可以为 `32` 和 `64`
	t.Run("strconv.FormatFloat", func(t *testing.T) {
		// 将 64 位浮点数转为小数字符串, 保留最多小数位
		s := strconv.FormatFloat(100.002, 'f', -1, 64)
		assert.Equal(t, "100.002", s)

		// 将 64 位浮点数转为 p 指数 (以 2 为底) 的字符串, 保留最多小数位
		s = strconv.FormatFloat(100.002, 'b', -1, 64)
		assert.Equal(t, "7037015155254755p-46", s)

		// 将 64 位浮点数转为 e 指数 (以 10 为底) 的字符串, 保留最多小数位
		s = strconv.FormatFloat(100.002, 'e', -1, 64)
		assert.Equal(t, "1.00002e+02", s)

		// 将 32 位浮点数转小数字符串, 保留 2 位小数位
		s = strconv.FormatFloat(100.002, 'f', 2, 64)
		assert.Equal(t, "100.00", s)
	})

	// 将布尔值转为字符串
	t.Run("strconv.FormatBool", func(t *testing.T) {
		// 将布尔值转为字符串
		s := strconv.FormatBool(true)
		assert.Equal(t, "true", s)

		// 将布尔值转为字符串
		s = strconv.FormatBool(false)
		assert.Equal(t, "false", s)
	})

	// 将复数值转为字符串
	// 需要指定结果的格式 (`fmt` 参数), 可以为
	//   - 'b' (`-dddp±dddd`, 以 2 为低的指数)
	//   - 'e' (`-d.dddde±dd`, 以 10 为底的指
	//   - 'E' (`-d.ddddE±dd`, 同上, `E` 字母
	//   - 'f' (`-ddd.dddd`，无指数)
	//   - 'g' (根据小数位数自动选择 'e' 或 'f
	//   - 'G' (根据小数位数自动选择 'E' 或 'f
	//   - 'x' (`-0xd.ddddp±dddd`, 小数部分为
	//   - 'X' (`-0Xd.ddddP±dddd`, 小数部分为
	//
	// 需要指定保留小数位数 (`prec` 参数)
	// 需要指定二进制位数, 可以为 `64` 和 `128`
	t.Run("strconv.FormatComplex", func(t *testing.T) {
		c := 100.002 + 20i

		// 将复数转为小数格式字符串, 保留最多小数位
		s := strconv.FormatComplex(c, 'f', -1, 64)
		assert.Equal(t, "(100.002+20i)", s)

		// 将复数转为 p 指数格式 (以 2 为底) 字符串, 保留最多小数位
		s = strconv.FormatComplex(c, 'b', -1, 64)
		assert.Equal(t, "(13107462p-17+10485760p-19i)", s)

		// 将复数转为 e 指数格式 (以 10 为底) 字符串, 保留最多小数位
		s = strconv.FormatComplex(c, 'e', -1, 64)
		assert.Equal(t, "(1.00002e+02+2e+01i)", s)

		// 将复数转为小数格式字符串, 保留 2 位小数位
		s = strconv.FormatComplex(c, 'f', 2, 64)
		assert.Equal(t, "(100.00+20.00i)", s)
	})
}

// 测试将字符串转为值类型
func TestStrConv_Parse(t *testing.T) {
	// 将字符串按照要求的进制及二进制位数转为整数
	// 需要指定结果整数的进制 (`base` 参数) 及二进制位数 (`bitSize` 参数)
	// `bitSize` 参数可以为 `0`, `8`, `16`, `32`, `64`, `0` 表示自动判断
	t.Run("strconv.ParseInt", func(t *testing.T) {
		// 将字符串转为 10 进制 64 位整数
		s, err := strconv.ParseInt("100", 10, 64)
		assert.Nil(t, err)
		assert.Equal(t, int64(100), s)

		// 将字符串转为 16 进制 64 位整数
		s, err = strconv.ParseInt("64", 16, 64)
		assert.Nil(t, err)
		assert.Equal(t, int64(100), s)

		// 将字符串转为 8 进制 32 位整数
		s, err = strconv.ParseInt("144", 8, 32)
		assert.Nil(t, err)
		assert.Equal(t, int64(100), s)

		// 将字符串转为 2 进制 64 位整数
		s, err = strconv.ParseInt("1100100", 2, 64)
		assert.Nil(t, err)
		assert.Equal(t, int64(100), s)

		// 测试转换错误, 字符串中不包含非数字字符
		_, err = strconv.ParseInt("abcd", 10, 64)
		assert.EqualError(t, err, "strconv.ParseInt: parsing \"abcd\": invalid syntax")
	})

	// 将字符串按照要求的进制及二进制位数转为无符号整数
	// 需要指定结果整数的进制 (`base` 参数) 及位数 (`bitSize` 参数)
	// `bitSize` 参数可以为 `0`, `8`, `16`, `32`, `64`, `0` 表示自动判断
	t.Run("strconv.ParseUint", func(t *testing.T) {
		// 将字符串转为 10 进制 64 位无符号整数
		s, err := strconv.ParseUint("100", 10, 64)
		assert.Nil(t, err)
		assert.Equal(t, uint64(100), s)

		// 将字符串转为 16 进制 64 位无符号整数
		s, err = strconv.ParseUint("64", 16, 64)
		assert.Nil(t, err)
		assert.Equal(t, uint64(100), s)

		// 将字符串转为 8 进制 32 位无符号整数
		s, err = strconv.ParseUint("144", 8, 32)
		assert.Nil(t, err)
		assert.Equal(t, uint64(100), s)

		// 将字符串转为 2 进制 64 位无符号整数
		s, err = strconv.ParseUint("1100100", 2, 64)
		assert.Nil(t, err)
		assert.Equal(t, uint64(100), s)

		// 测试转换错误, 字符串中不包含非数字字符
		_, err = strconv.ParseUint("abcd", 10, 64)
		assert.EqualError(t, err, "strconv.ParseUint: parsing \"abcd\": invalid syntax")
	})

	// 将字符串按照要求的进制及二进制位数转为浮点数
	//
	// `bitSize` 参数可以为 `0`, `32`, `64`, `0` 表示自动判断
	t.Run("strconv.ParseFloat", func(t *testing.T) {
		// 将小数格式的字符串转为 64 位浮点数
		s, err := strconv.ParseFloat("100.002", 64)
		assert.Nil(t, err)
		assert.Equal(t, float64(100.002), s)

		// 将指数格式 (以 10 为底) 的字符串转为 64 位浮点数
		s, err = strconv.ParseFloat("1.00002e+02", 32)
		assert.Nil(t, err)
		assert.Equal(t, float32(100.002), float32(s))

		// 测试转换错误, 字符串中不包含非数字字符
		_, err = strconv.ParseFloat("abcd", 64)
		assert.EqualError(t, err, "strconv.ParseFloat: parsing \"abcd\": invalid syntax")
	})

	// 将字符串转为布尔值
	t.Run("strconv.ParseBool", func(t *testing.T) {
		// 将字符串转为布尔值
		s, err := strconv.ParseBool("true")
		assert.Nil(t, err)
		assert.Equal(t, true, s)

		// 将字符串转为布尔值
		s, err = strconv.ParseBool("false")
		assert.Nil(t, err)
		assert.Equal(t, false, s)

		// 测试转换错误, 字符串中不包含表示布尔值的内容
		_, err = strconv.ParseBool("abcd")
		assert.EqualError(t, err, "strconv.ParseBool: parsing \"abcd\": invalid syntax")
	})

	// 将字符串转为复数
	// 需要指定复数的二进制位数 (`bitSize` 参数), 可以为 `0`, `64` 和 `128`, `0` 表示自动判断
	t.Run("strconv.ParseComplex", func(t *testing.T) {
		// 将小数格式字符串转为复数
		s, err := strconv.ParseComplex("100.002+20i", 128)
		assert.Nil(t, err)
		assert.Equal(t, 100.002+20i, s)

		// 将指数格式 (以 10 为底) 的字符串转为复数
		s, err = strconv.ParseComplex("1.00002e+02+2e+01i", 128)
		assert.Nil(t, err)
		assert.Equal(t, 100.002+20i, s)

		// 测试转换错误, 字符串中不包含表示复数值的内容
		_, err = strconv.ParseComplex("abcd", 64)
		assert.EqualError(t, err, "strconv.ParseComplex: parsing \"abcd\": invalid syntax")

	})

}

// 测试将整数转为字符串
//
// 相当于 `strconv.FormatInt(int64(100), 10)` 函数
func TestStrConv_Itoa(t *testing.T) {
	s := strconv.Itoa(100)
	assert.Equal(t, "100", s)
}

// 测试将字符串转为整数
//
// 相当于 `strconv.ParseInt("100", 10, 0)` 函数, 并且返回值位 `int` 类型
func TestStrConv_Atoi(t *testing.T) {
	n, err := strconv.Atoi("100")

	assert.Nil(t, err)
	assert.Equal(t, 100, n)

	_, err = strconv.Atoi("abcd")
	assert.EqualError(t, err, "strconv.Atoi: parsing \"abcd\": invalid syntax")
}

// 测试将值类型转为字符串后, 追加到现有的字节切片中
func TestStrConv_Append(t *testing.T) {
	// 将一个整数转为字符串, 并追加到一个字节串之后
	t.Run("strconv.AppendInt", func(t *testing.T) {
		s := []byte("Hello ")

		bs := strconv.AppendInt(s, 100, 10)
		assert.Equal(t, []byte("Hello 100"), bs)
	})

	// 将一个布尔值转为字符串, 并追加到一个字节串之后
	t.Run("strconv.AppendBool", func(t *testing.T) {
		s := []byte("Hello ")

		bs := strconv.AppendBool(s, true)
		assert.Equal(t, []byte("Hello true"), bs)
	})

	// 将一个浮点数转为字符串, 并追加到一个字节串之后
	t.Run("strconv.AppendFloat", func(t *testing.T) {
		s := []byte("Hello ")

		bs := strconv.AppendFloat(s, 100.002, 'f', -1, 64)
		assert.Equal(t, []byte("Hello 100.002"), bs)
	})
}

// 测试判断一个字符是否为可打印字符
func TestStrConv_IsPrint(t *testing.T) {
	// 可打印字符返回 true
	assert.True(t, strconv.IsPrint('a'))
	assert.True(t, strconv.IsPrint('\u0020'))

	// 非可打印字符返回 false
	assert.False(t, strconv.IsPrint('\n'))
	assert.False(t, strconv.IsPrint(127))
}

// 测试判断一个字符是否为图形字符
//
// 图形字符包括类别字母, 标记, 数字, 符号, 标点, 空格等
func TestStrConv_IsGraphic(t *testing.T) {
	assert.True(t, strconv.IsGraphic('a'))
	assert.True(t, strconv.IsGraphic('☺'))
	assert.False(t, strconv.IsGraphic('\n'))
	assert.False(t, strconv.IsGraphic(127))
	assert.True(t, strconv.IsGraphic('\u0020'))
}

// 测试字符
func TestStrConv_Quote(t *testing.T) {
	// 将字符串转为“双引号”包围的字符串字面量
	// 返回的字符串会被一对双引号包围, 并对其中的特殊字符进行转义
	t.Run("strconv.Quote", func(t *testing.T) {
		s := strconv.Quote(`Fran & Freddie's Diner	"☺", Ok`)

		assert.Equal(t, `"Fran & Freddie's Diner\t\"☺\", Ok"`, s)
		assert.Equal(t, "\"Fran & Freddie's Diner\\t\\\"☺\\\", Ok\"", s)
	})

	// 将特殊字符转为转义字符
	// 注意, 结果中会包含字符最外层的单引号
	t.Run("strconv.QuoteRune", func(t *testing.T) {
		c := strconv.QuoteRune('	')
		assert.Equal(t, `'\t'`, c)
	})

	// 将字符转为 Unicode 表示
	// 返回一个字符串, 内容为一个单引号包围的 Unicode 字符
	// 对于非可见字符, 会返回其转义字符, 例如 `\n`, `\t` 等
	t.Run("strconv.QuoteRuneToGraphic", func(t *testing.T) {
		s := strconv.QuoteRuneToGraphic('	')
		assert.Equal(t, `'\t'`, s)

		s = strconv.QuoteRuneToGraphic('\u263a')
		assert.Equal(t, `'☺'`, s)

		s = strconv.QuoteRuneToGraphic('\u000a')
		assert.Equal(t, `'\n'`, s)

		s = strconv.QuoteRuneToGraphic('中')
		assert.Equal(t, `'中'`, s)
	})

	// 去掉字符串外围的引号
	t.Run("strconv.Unquote", func(t *testing.T) {
		_, err := strconv.Unquote("无法处理不包含引号的字符串")
		assert.EqualError(t, err, "invalid syntax")

		sq := strconv.Quote(`Fran & Freddie's Diner	☺`)

		s, err := strconv.Unquote(sq)
		assert.Nil(t, err)
		assert.Equal(t, "Fran & Freddie's Diner\t☺", s)
	})
}

// 测试获取字符串数据指针
//
// `unsafe.StringData` 函数获取字符串数据指针, 即字符串结构中存储的字节数据地址
func TestUnsafe_StringData(t *testing.T) {
	// 获取字符串数据地址
	bs := unsafe.StringData("hello world")
	assert.Equal(t, byte('h'), *bs)

	// 通过移动指针访问字符串中其它的字节数据
	bs = utils.PtrAdd(bs, 6)
	assert.Equal(t, byte('w'), *bs)
}

// 测试基于字节指针产生字符串变量
//
// `unsafe.String` 函数可以零拷贝方式, 基于一个指向连续字节数据的指针产生一个字符串变量
func TestUnsafe_String(t *testing.T) {
	bs := []byte("hello world")

	// 获取字节切片的数据指针, 转为字符串
	s := unsafe.String(unsafe.SliceData(bs), len(bs))
	assert.Equal(t, "hello world", s)

	// 获取字节切片的数据指针, 并将指定长度的数据转为字符串
	s = unsafe.String(unsafe.SliceData(bs), len(bs)-6)
	assert.Equal(t, "hello", s)

	// 获取字节切片的数据指针, 将指针移动 6 字节后, 将剩余部分转为字符串
	s = unsafe.String(utils.PtrAdd(unsafe.SliceData(bs), 6), len(bs)-6)
	assert.Equal(t, "world", s)

	n := int64(7163384699739271026)

	// 将指向整数的指针转为字节指针, 并基于该指针将数据转为字符串
	s = unsafe.String((*byte)(unsafe.Pointer(&n)), unsafe.Sizeof(n))
	assert.Equal(t, "romantic", s)
}

// 测试基于字节切片产生字符串变量 (过时方式)
func TestUnsafe_BytesToStringLegacy(t *testing.T) {
	bs := []byte("hello world")

	s := conv.BytesToStringLegacy(bs)
	assert.Equal(t, "hello world", s)
}

// 测试基于字符串产生字节切片变量 (过时方式)
func TestUnsafe_StringToBytesLegacy(t *testing.T) {
	s := "hello world"

	bs := conv.StringToBytesLegacy(s)
	assert.Equal(t, []byte("hello world"), bs)
}
