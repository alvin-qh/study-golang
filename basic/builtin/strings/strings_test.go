package strings_test

import (
	"io"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

// 测试比较字符串
func TestStrings_Compare(t *testing.T) {
	s := "abc"

	// 测试比较两个字符串
	t.Run("strings.Compare", func(t *testing.T) {
		// 两个字符串比较, 相等则返回 0
		assert.Equal(t, 0, strings.Compare(s, "abc"))

		// 第一个字符串大于第二个字符串, 返回 1
		assert.Equal(t, 1, strings.Compare(s, "Abc"))

		// 第一个字符串小于第二个字符串, 返回 -1
		assert.Equal(t, -1, strings.Compare(s, "bbc"))
	})

	// 测试对两个字符串进行忽略大小写的比较, 返回是否相等
	t.Run("strings.EqualFold", func(t *testing.T) {
		assert.True(t, strings.EqualFold(s, "ABC"))
	})
}

// 测试指定字符串是否包含子字符串
func TestStrings_Contains(t *testing.T) {
	s := "Hello,大家好"

	assert.True(t, strings.Contains(s, ",大"))
	assert.False(t, strings.Contains(s, "好大"))
}

// 测试查找子字符串在字符串中出现的位置
func TestStrings_Index(t *testing.T) {
	s := "Hello,大家好"

	// 测试查找子字符串在字符串中第一次出现的位置
	t.Run("strings.Index", func(t *testing.T) {
		n := strings.Index(s, "家好")
		assert.Equal(t, 9, n)
	})

	// 测试从字符串末尾开始查找
	t.Run("strings.LastIndex", func(t *testing.T) {
		n := strings.LastIndex(s, "家好")
		assert.Equal(t, 9, n)
	})

	// 测试所给字符串中的任意字符在指定字符串中第一次出现的位置
	t.Run("strings.IndexAny", func(t *testing.T) {
		n := strings.IndexAny(s, "o好")
		assert.Equal(t, 4, n)
	})

	// 测试所给字符串中的任意字符在指定字符串中第一次出现的位置 (反向查找)
	t.Run("strings.LastIndexAny", func(t *testing.T) {
		n := strings.LastIndexAny(s, "o好")
		assert.Equal(t, 12, n)
	})

	// 测试在字符串中查找指定字节第一次出现的位置
	t.Run("strings.IndexByte", func(t *testing.T) {
		n := strings.IndexByte(s, 'l')
		assert.Equal(t, 2, n)
	})

	// 测试在字符串中查找指定字节第一次出现的位置 (反向查找)
	t.Run("strings.LastIndexByte", func(t *testing.T) {
		n := strings.LastIndexByte(s, 'l')
		assert.Equal(t, 3, n)
	})

	// 测试通过回调函数查找字符串指定字符第一次出现的位置
	t.Run("strings.IndexRune", func(t *testing.T) {
		// 获取 ',' 字符第一次出现的位置
		n := strings.IndexFunc(s, func(r rune) bool { return r == ',' })
		assert.Equal(t, 5, n)
	})

	// 测试通过回调函数查找字符串指定字符第一次出现的位置
	t.Run("strings.LastIndexRune", func(t *testing.T) {
		// 从字符串末尾开始, 获取 ',' 字符第一次出现的位置
		n := strings.LastIndexFunc(s, func(r rune) bool { return r == ',' })
		assert.Equal(t, 5, n)
	})
}

// 测试字符串是否以指定的子字符串开始
func TestStrings_HasPrefix(t *testing.T) {
	s := "Hello,大家好"

	b := strings.HasPrefix(s, "Hello")
	assert.True(t, b)
}

// 测试字符串是否以指定的子字符串结束
func TestStrings_HasSuffix(t *testing.T) {
	s := "Hello,大家好"

	b := strings.HasSuffix(s, "家好")
	assert.True(t, b)
}

// 测试子字符串在原字符串出现的次数
func TestStrings_Count(t *testing.T) {
	s := "Hello,大家好"

	// 统计字符串中 "l" 字符串出现的次数
	n := strings.Count(s, "l")
	assert.Equal(t, 2, n)

	// 统计字符串中空字串出现的次数, 比字符串长度多 1
	n = strings.Count(s, "")
	assert.Equal(t, utf8.RuneCountInString(s)+1, n)
}

// 测试替换字符串中的一部分
func TestStrings_Replace(t *testing.T) {
	s := "Hello,大家好"

	// 将子字符串替换为指定的字符串, 共替换 1 次
	r := strings.Replace(s, "l", "L", 1)
	assert.Equal(t, "HeLlo,大家好", r)

	// 将子字符串替换为指定的字符串, 共替换 2 次
	r = strings.Replace(s, "l", "L", 2)
	assert.Equal(t, "HeLLo,大家好", r)

	// 将子字符串替换为指定的字符串, 共替换任意次
	r = strings.Replace(s, "l", "L", -1)
	assert.Equal(t, "HeLLo,大家好", r)

	// 替换所有的指定子字符串, 相当于 `strings.Replace(s, "l", "L", -1)`
	r = strings.ReplaceAll(s, "l", "L")
	assert.Equal(t, "HeLLo,大家好", r)
}

// 测试删除字符串前后的指定字符
func TestStrings_Trim(t *testing.T) {
	s := "Hello,大家好"

	// 测试删除字符串前后的指定字符
	t.Run("strings.Trim", func(t *testing.T) {
		// 删除字符串前后的指定字符, 字符集合中的内容会被全部删除
		r := strings.Trim(s, "He好")
		assert.Equal(t, "llo,大家", r)
	})

	// 测试删除字符串前后的空白字符
	// 相当于 `strings.Trim(s, " \r\n\t")`
	t.Run("strings.TrimSpace", func(t *testing.T) {
		s := " \rAAAA\t\n"

		r := strings.TrimSpace(s)
		assert.Equal(t, "AAAA", r)
	})

	// 测试删除字符串起始位置的指定字符
	t.Run("strings.TrimLeft", func(t *testing.T) {
		r := strings.TrimLeft(s, "He好")
		assert.Equal(t, "llo,大家好", r)
	})

	// 删除字符串开始的指定部分
	// 注意, 和 `strings.TrimLeft` 不同, `strings.TrimPrefix` 是删除整个指定的字符串, 而不是其中的某个字符
	t.Run("strings.TrimPrefix", func(t *testing.T) {
		r := strings.TrimPrefix(s, "Hello,")
		assert.Equal(t, "大家好", r)
	})

	// 测试删除字符串结束位置的指定字符
	t.Run("strings.TrimRight", func(t *testing.T) {
		r := strings.TrimRight(s, "He好")
		assert.Equal(t, "Hello,大家", r)
	})

	// 删除字符串结束的指定部分
	// 注意, 和 `strings.TrimRight` 不同, `strings.TrimSuffix` 是删除整个指定的字符串, 而不是其中的某个字符
	t.Run("strings.TrimSuffix", func(t *testing.T) {
		r := strings.TrimSuffix(s, ",大家好")
		assert.Equal(t, "Hello", r)
	})

	// 测试根据回调函数删除字符串开始位置的指定字符
	t.Run("strings.TrimLeftFunc", func(t *testing.T) {
		r := strings.TrimLeftFunc(s, func(r rune) bool {
			// 删除字节长度为 1 的字符, 即 ASCII 编码字符
			return utf8.RuneLen(r) == 1
		})
		assert.Equal(t, "大家好", r)
	})

	// 测试根据回调函数删除字符串结束位置的指定字符
	t.Run("strings.TrimRightFunc", func(t *testing.T) {
		r := strings.TrimRightFunc(s, func(r rune) bool {
			// 删除字节长度大于 1 的字符, 即汉字字符
			return utf8.RuneLen(r) > 1
		})
		assert.Equal(t, "Hello,", r)
	})
}

// 测试字符串分割
func TestStrings_Split(t *testing.T) {
	s := "Hello,大家好"

	// 将字符串分隔成若干子字符串
	t.Run("strings.Split", func(t *testing.T) {
		r := strings.Split(s, ",")
		assert.Equal(t, []string{"Hello", "大家好"}, r)
	})

	// 指定分割结果的最大数量
	t.Run("strings.SplitN", func(t *testing.T) {
		// 指定分割结果的数量, 至多将字符串分割为 1 个部分, 相当于不做分割
		r := strings.SplitN(s, ",", 1)
		assert.Equal(t, []string{"Hello,大家好"}, r)

		// 指定分割结果的数量, 至多将字符串分割为 2 个部分
		r = strings.SplitN(s, ",", 2)
		assert.Equal(t, []string{"Hello", "大家好"}, r)

		// 指定分割结果的数量, 至多将字符串分割为 2 个部分
		r = strings.SplitN(s, ",", 2)
		assert.Equal(t, []string{"Hello", "大家好"}, r)

		// 分割为任意部分, 相当于 strings.Split(s)
		r = strings.SplitN(s, ",", -1)
		assert.Equal(t, []string{"Hello", "大家好"}, r)
	})

	// 在分割结果中包含用于分割的字符串本身
	// 用于分隔的字符串和前一个分割结果合并在一起
	t.Run("strings.SplitAfter", func(t *testing.T) {
		// 分割结果中包含用于分割的字符串本身
		r := strings.SplitAfter(s, ",")
		assert.Equal(t, []string{"Hello,", "大家好"}, r)
	})

	// 在分割结果中包含用于分割的字符串本身, 并指定分割结果的最大数量
	// 用于分隔的字符串和前一个分割结果合并在一起
	t.Run("strings.SplitAfterN", func(t *testing.T) {
		// 分割结果中至多包含 1 个子字符串
		r := strings.SplitAfterN(s, ",", 1)
		assert.Equal(t, []string{"Hello,大家好"}, r)

		// 分割结果中至多包含 2 个子字符串
		r = strings.SplitAfterN(s, ",", 2)
		assert.Equal(t, []string{"Hello,", "大家好"}, r)

		// 分割为任意部分, 相当于 strings.SplitAfter(s)
		r = strings.SplitAfterN(s, ",", -1)
		assert.Equal(t, []string{"Hello,", "大家好"}, r)
	})
}

// 测试字符串连接
//
// 字符串连接, 即将若干个子字符串连接成一个完整的字符串
func TestStrings_Join(t *testing.T) {
	s := "Hello"

	// 字符串连接, 将字符串数组 (切片) 通过连接符进行连接
	r := strings.Join([]string{s, "World"}, " ")
	assert.Equal(t, "Hello World", r)
}

// 测试字符串重复
//
// 即将指定的字符串重复指定次数后, 形成新字符串
func TestStrings_Repeat(t *testing.T) {
	s := "Hello"

	// 重复指定次数的字符串
	r := strings.Repeat(s, 3)
	assert.Equal(t, "HelloHelloHello", r)
}

// 测试将字符串中的字母转为小写或大写
func TestStrings_ToLower_ToUpper(t *testing.T) {
	s := "Hello,大家好"

	// 将字符串转换为小写
	r := strings.ToLower(s)
	assert.Equal(t, "hello,大家好", r)

	// 将字符串转换为大写
	r = strings.ToUpper(s)
	assert.Equal(t, "HELLO,大家好", r)
}

// 测试字符串切割
func TestStrings_Cut(t *testing.T) {
	s := "a,b,c,d,e"

	// 测试根据分隔符将字符串分割为两部分
	t.Run("strings.Cut", func(t *testing.T) {
		// 根据分隔符切分字符串
		// 返回切分后的前一部分, 后一部分以及是否找到分隔符
		before, after, found := strings.Cut(s, ",")

		assert.True(t, found)
		assert.Equal(t, "a", before)
		assert.Equal(t, "b,c,d,e", after)
	})

	// 测试根据所给的前缀将字符串分为两部分
	t.Run("strings.CutPrefix", func(t *testing.T) {
		// 根据前缀切分字符串
		// 返回切分后的后一部分以及是否找到前缀
		after, found := strings.CutPrefix(s, "a,b,")

		assert.True(t, found)
		assert.Equal(t, "c,d,e", after)
	})

	// 测试根据所给的后缀将字符串分为两部分
	t.Run("strings.CutPrefix", func(t *testing.T) {
		// 根据后缀切分字符串
		// 返回切分后的前一部分以及是否找到后缀
		before, found := strings.CutSuffix(s, "c,d,e")

		assert.True(t, found)
		assert.Equal(t, "a,b,", before)
	})
}

// 测试字符串读取器, 用于从字符串中读取不同数据
func TestStrings_Reader(t *testing.T) {
	// 获取 `Reader` 长度
	t.Run("length", func(t *testing.T) {
		// 通过字符串创建一个读取器实例
		r := strings.NewReader("abcdefghijklmnopqrstuvwxyz")

		// 获取可读取字节数
		assert.Equal(t, 26, r.Len())
		assert.Equal(t, int64(26), r.Size())
	})

	// 依次读取指定字节内容
	t.Run("read", func(t *testing.T) {
		// 通过字符串创建一个读取器实例
		r := strings.NewReader("abcdefghijklmnopqrstuvwxyz")

		data := make([]byte, 10)

		// 从字符串当前位置读取 5 字节
		// 读取结束后, 当前读取位置移动 5 字节
		n, err := r.Read(data)
		assert.Nil(t, err)
		assert.Equal(t, 10, n)
		assert.Equal(t, "abcdefghij", string(data))

		// 继续读取 5 字节
		n, err = r.Read(data[:5])
		assert.Nil(t, err)
		assert.Equal(t, 5, n)
		assert.Equal(t, "klmno", string(data[:5]))
	})

	// 读取字节
	t.Run("read byte", func(t *testing.T) {
		// 通过字符串创建一个读取器实例
		r := strings.NewReader("abcdefghijklmnopqrstuvwxyz")

		// 读取一个字节
		b, err := r.ReadByte()
		assert.Nil(t, err)
		assert.Equal(t, 'a', rune(b))

		// 撤销一个字节读取, 即当前读取位置前移一个字节
		err = r.UnreadByte()
		assert.Nil(t, err)

		// 再次读取一个字节, 为之前已读取过的字节
		b, err = r.ReadByte()
		assert.Nil(t, err)
		assert.Equal(t, 'a', rune(b))
	})

	// 读取字符
	t.Run("read rune", func(t *testing.T) {
		// 通过字符串创建一个读取器实例
		r := strings.NewReader("abcdefghijklmnopqrstuvwxyz")

		// 继续读取一个字符
		c, n, err := r.ReadRune()
		assert.Nil(t, err)
		assert.Equal(t, 1, n)
		assert.Equal(t, 'a', c)

		// 撤销一个字符读取, 即当前读取位置前移一个字节
		err = r.UnreadByte()
		assert.Nil(t, err)

		// 再次读取一个字符, 为之前已读取过的字符
		c, n, err = r.ReadRune()
		assert.Nil(t, err)
		assert.Equal(t, 1, n)
		assert.Equal(t, 'a', c)
	})

	// 从指定位置读取字节
	t.Run("read at", func(t *testing.T) {
		// 通过字符串创建一个读取器实例
		r := strings.NewReader("abcdefghijklmnopqrstuvwxyz")

		data := make([]byte, 10)

		// 从第 10 个字节开始读取 10 字节
		// 该方法和 `Reader` 当前的读取位置无关
		n, err := r.ReadAt(data, 10)
		assert.Nil(t, err)
		assert.Equal(t, 10, n)
		assert.Equal(t, "klmnopqrst", string(data))
	})

	// 测试移动读取位置
	t.Run("seek", func(t *testing.T) {
		// 通过字符串创建一个读取器实例
		r := strings.NewReader("abcdefghijklmnopqrstuvwxyz")

		data := make([]byte, 10)

		// 读取 10 字节, 将读取位置移动到 10 字节位置
		n, err := r.Read(data)
		assert.Nil(t, err)
		assert.Equal(t, 10, n)
		assert.Equal(t, "abcdefghij", string(data))

		// 将读取位置从当前位置向前移动 3 字节, 将读取位置移动到 7 字节位置
		pos, err := r.Seek(-3, io.SeekCurrent)
		assert.Nil(t, err)
		assert.Equal(t, int64(7), pos)

		// 从第 7 字节开始, 继续读取 10 字节
		n, err = r.Read(data)
		assert.Nil(t, err)
		assert.Equal(t, 10, n)
		assert.Equal(t, "hijklmnopq", string(data))

		r.Reset("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	})

	// 重新设置字符串内容
	t.Run("reset", func(t *testing.T) {
		// 通过字符串创建一个读取器实例
		r := strings.NewReader("abcdefghijklmnopqrstuvwxyz")

		// 读取 Reader 中的全部内容
		data, err := io.ReadAll(r)
		assert.Nil(t, err)
		assert.Equal(t, "abcdefghijklmnopqrstuvwxyz", string(data))

		// 重设 Reader 中的字符串
		r.Reset("ABCDEFGHIJKLMNOPQRSTUVWXYZ")

		// 再次读取 Reader 中的全部内容
		data, err = io.ReadAll(r)
		assert.Nil(t, err)
		assert.Equal(t, "ABCDEFGHIJKLMNOPQRSTUVWXYZ", string(data))
	})
}
