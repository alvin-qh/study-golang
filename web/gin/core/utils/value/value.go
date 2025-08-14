package value

import (
	"reflect"
	"strings"

	"study/web/gin/core/utils/conv"
)

const (
	// 最大的整数值
	maxInt = int(^uint(0) >> 1)
)

// 如果 `val` 参数的值为 `nil`, 则返回 `def` 参数的值, 否则返回 `val` 参数的值
//
// 参数:
//   - `val` (`any`): 参数值, 如果该值不为 `nil`, 则函数返回该参数值
//   - `def` (`T`): 默认值, 如果 `val` 参数值为 `nil`, 则函数返回该参数值
//
// 返回:
//   - `T`: 如果 `val` 参数为 `nil`, 则为 `def` 参数的值; 否则为 `val` 参数的值
func Default[T any](val any, def T) T {
	if val == nil {
		return def
	}

	var ok bool
	v, ok := val.(T)
	if !ok {
		// 如果 val 参数类型和 def 不一致, 也返回 def 参数值
		return def
	}
	return v
}

// 将任意元素类型切片连接为字符串
//
// 参数:
//   - `elems` (`[]T`): 任意元素类型切片
//   - `sep` (`string`): 用于将数组元素连接的分隔符
//
// 参数:
//   - `result` (`string`): 数组元素连接而成的字符串
func Join[T any](elems []T, sep string) (result string) {
	// 处理数组只有 0 和 1 项的情况
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return conv.AnyToString(elems[0])
	}

	// 计算分隔符在结果字符串中会占用的长度
	var n int
	if len(sep) > 0 {
		// 如果分隔符最终占用长度超过最大整数值, 则抛出异常
		if len(sep) >= maxInt/(len(elems)-1) {
			panic("join output length overflow")
		}
		// 计算分隔符最终占用的长度
		n += len(sep) * (len(elems) - 1)
	}

	// 计算生成字符串的占用长度
	sElems := make([]string, len(elems))
	for i, elem := range elems {
		// 遍历数组元素
		v, ok := (any(elem)).(string)
		if !ok {
			// 如果数组元素不为字符串类型, 则转为字符串
			v = conv.AnyToString(elem)
		}

		// 如果最终生成的字符串长度超过最大整数值, 则抛出异常
		if len(v) > maxInt-n {
			panic("strings: Join output length overflow")
		}
		// 记录最终生成字符串的长度
		n += len(v)
		sElems[i] = v
	}

	// 按最终生成字符串的长度生成 Builder 对象
	var b strings.Builder
	b.Grow(n)

	// 写入数组第一个元素
	b.WriteString(sElems[0])

	// 写入分隔符和数组剩余元素
	for _, s := range sElems[1:] {
		b.WriteString(sep)
		b.WriteString(s)
	}
	// 返回生成的字符串
	return b.String()
}

// 将各种类型切片 (或数组) 连接为字符串
//
// 参数:
//   - `elems` (`any`): 任意类型数组, 可以为 `[]int`, `[]string`, `[]any` 等
//   - `sep` (`string`): 用于将数组元素连接的分隔符
//
// 参数:
//   - `result` (`string`): 数组元素连接而成的字符串
func JoinAny(elems any, sep string) string {
	// 获取 elems 参数类型
	t := reflect.ValueOf(elems)
	if t.Kind() != reflect.Array && t.Kind() != reflect.Slice {
		// 如果 elemes 参数不为切片或数组, 则返回其字符串形式
		return conv.AnyToString(elems)
	}

	// 处理数组只有 0 和 1 项的情况
	switch t.Len() {
	case 0:
		return ""
	case 1:
		return conv.AnyToString(t.Index(0).Interface())
	}

	// 计算分隔符在结果字符串中会占用的长度
	var n int
	if len(sep) > 0 {
		if len(sep) >= maxInt/(t.Len()-1) {
			panic("join output length overflow")
		}
		n += len(sep) * (t.Len() - 1)
	}

	// 计算生成字符串的占用长度
	sElems := make([]string, t.Len())
	for i := 0; i < t.Len(); i++ {
		// 遍历数组元素
		elem := t.Index(i).Interface()
		v, ok := elem.(string)
		if !ok {
			// 如果数组元素不为字符串类型, 则转为字符串
			v = conv.AnyToString(elem)
		}

		// 如果最终生成的字符串长度超过最大整数值, 则抛出异常
		if len(v) > maxInt-n {
			panic("strings: Join output length overflow")
		}
		// 记录最终生成字符串的长度
		n += len(v)
		sElems[i] = v
	}

	// 按最终生成字符串的长度生成 Builder 对象
	var b strings.Builder
	b.Grow(n)

	// 写入数组第一个元素
	b.WriteString(sElems[0])
	for _, s := range sElems[1:] {
		b.WriteString(sep)
		b.WriteString(s)
	}
	// 返回生成的字符串
	return b.String()
}
