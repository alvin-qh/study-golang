package utils

import (
	"fmt"
	"reflect"
	"strings"
	"time"
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
		return AnyToString(elems[0])
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
			v = AnyToString(elem)
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
		return AnyToString(elems)
	}

	// 处理数组只有 0 和 1 项的情况
	switch t.Len() {
	case 0:
		return ""
	case 1:
		return AnyToString(t.Index(0).Interface())
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
			v = AnyToString(elem)
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

var (
	// 注册自定义类型的 Map
	registerTypes = map[reflect.Type]func(any) (any, error){
		// 为 time.Time 类型注册默认的处理函数
		reflect.TypeOf(time.Time{}): func(v any) (rv any, err error) {
			// 根据值的不同类型进行不同的转换
			switch _v := v.(type) {
			case string:
				// 将字符串表示的时间转为 time.Time 类型
				rv, err = time.Parse(time.RFC3339, _v)
				if err != nil {
					rv, err = time.Parse(time.TimeOnly, _v)
					if err != nil {
						rv, err = time.Parse(time.DateOnly, _v)
					}
				}
				if err != nil {
					return
				}
			case float64:
				// 将 float64 表示的时间转为 time.Time 类型
				rv = time.Unix(0, int64(v.(float64))*int64(time.Millisecond))
			case int64:
				// 将 int64 表示的时间转为 time.Time 类型
				rv = time.Unix(0, v.(int64)*int64(time.Millisecond))
			default:
				// 其它类型, 不做任何转换
				rv = v
			}
			return
		},
	}
)

// 为指定类型注册转换函数
func RegisterTypeConverter(t reflect.Type, converter func(any) (any, error)) {
	registerTypes[t] = converter
}

// 判断一个反射值对象是否表示结构体指针
//
// 参数:
//   - `v` (`*reflect.Value`): 反射值对象指针
//
// 返回:
//   - `bool` 返回 `v` 参数是否表示结构体指针
func checkIfStructPointer(v *reflect.Value) bool {
	return v.Kind() == reflect.Pointer && v.Elem().Kind() == reflect.Struct
}

type fieldDefine struct {
	isStruct bool
	value    reflect.Value
	fieldMap map[string]*fieldDefine
}

// 将结构体反射值转为字段 + 字段反射值
func structFieldsToMap(v reflect.Value) map[string]*fieldDefine {
	// 获取结构体类型
	t := v.Type()

	// 定义递归闭包函数, 对当前结构体 (以及其结构体字段) 逐级进行处理
	var fieldsToMap func(v reflect.Value, t reflect.Type) map[string]*fieldDefine
	fieldsToMap = func(v reflect.Value, t reflect.Type) map[string]*fieldDefine {
		// 创建结构体字段 map 集合
		fm := make(map[string]*fieldDefine)

		// 遍历结构体的所有字段反射值
		for i := 0; i < v.NumField(); i++ {
			// 获取结构体字段的反射值
			fv := v.Field(i)
			if !fv.IsValid() || !fv.CanSet() {
				continue
			}

			// 获取结构体字段类型
			ft := t.Field(i)

			// 获取字段上定义的 tag
			tag := ft.Tag.Get("json")
			if len(tag) == 0 {
				tag = strings.ToLower(ft.Name)
			} else {
				tag = strings.TrimSpace(strings.Split(tag, ",")[0])
			}

			// 查询字段类型是否已被注册
			_, reg := registerTypes[ft.Type]

			// 对于未注册的类型, 且类型为结构体, 则递归处理结构体内的内容
			if !reg && fv.Kind() == reflect.Struct {
				fm[tag] = &fieldDefine{
					isStruct: true,
					fieldMap: fieldsToMap(fv, fv.Type()),
				}
			} else {
				// 将结构体字段反射值和 tag 对应
				fm[tag] = &fieldDefine{
					value: fv,
				}
			}
		}

		return fm
	}

	return fieldsToMap(v, t)
}

// 将一个值转为反射值
func toReflectValue(v any) reflect.Value {
	// 获取 value 值的反射值
	rv := reflect.ValueOf(v)

	// 如果该值为 interface{} 类型, 则进一步获取其原始类型
	if rv.Kind() == reflect.Interface {
		rv = rv.Elem()
	}
	return rv
}

// 将 map 集合内容填充到结构体对象
//
// 参数:
//   - `data` (`map[string]any`): map 集合对象
//   - `target` (`any`): 结构体对象指针
//
// 返回
//   - `error`: 错误对象
func MapToStruct(data map[string]any, target any) error {
	// 获取 target 参数的反射值
	tv := reflect.ValueOf(target)

	// 判断 target 参数的类型是否为结构体指针
	if !checkIfStructPointer(&tv) {
		return fmt.Errorf("target not a struct")
	}

	// 获取 target 参数指针指向的实际类型
	tv = tv.Elem()

	var mapToStruct func(map[string]any, map[string]*fieldDefine) error
	mapToStruct = func(data map[string]any, m map[string]*fieldDefine) (err error) {
		// 遍历保存 key/value 的 map 集合
		for k, v := range data {
			// 根据 key 值查询结构体对应的字段反射值
			fv, ok := m[strings.ToLower(k)]
			if !ok {
				continue
			}

			if fv.isStruct {
				// 进一步递归字段为结构体的内容
				if err = mapToStruct(v.(map[string]any), fv.fieldMap); err != nil {
					return err
				}
			} else {
				// 获取结构体字段类型
				ft := fv.value.Type()

				// 查找注册的转换函数
				fn, ok := registerTypes[ft]
				if ok {
					if v, err = fn(v); err != nil {
						return err
					}
				}

				vv := toReflectValue(v)

				// 将 value 设置到结构体字段中
				fv.value.Set(vv.Convert(ft))
			}
		}
		return
	}

	return mapToStruct(data, structFieldsToMap(tv))
}
