package maptostruct

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// 定义转换函数类型, 用于将一个值转为另一个类型值
//
// 参数:
//   - `v` (`any`): 原始类型值
//
// 返回：
//   - `r` (`any`): 转换后的值
//   - `err` (`error`): 转换错误对象
type MapperFn func(v any) (r any, err error)

// 结构体转换类型
type MapToStruct struct {
	tagKey  string                    // 结构体字段注解的 Key
	mappers map[reflect.Type]MapperFn // 指定类型的转换函数
}

// 创建 `MapToStruct` 对象
//
// 参数:
//   - `tagKey` (`string`): 结构体字段注解 Key
//
// 返回:
//   - `*MapToStruct`: `MapToStruct` 结构体指针
func New(tagKey string) *MapToStruct {
	return &MapToStruct{
		tagKey: tagKey,
		mappers: map[reflect.Type]MapperFn{
			reflect.TypeOf(time.Time{}): mapTime, // 设置默认的类型转换函数
		},
	}
}

// 为指定类型添加转换函数
//
// 参数:
//   - `t` (`reflect.Type`): 类型对象, 对于此类型的值进行转换
//   - `mapper` (`MapperFn`): 转换函数
func (m *MapToStruct) AddMapper(t reflect.Type, mapper MapperFn) {
	m.mappers[t] = mapper
}

// 从结构体字段获取注解的 tag
//
// 参数:
//   - `f` (`*reflect.StructField`): 结构体字段反射值
//
// 返回:
//   - `string`: 找到的 tag 值
func (m *MapToStruct) findTag(f *reflect.StructField) string {
	if len(m.tagKey) > 0 {
		// 查找预设 tag key 对应的 tag
		if tag, ok := f.Tag.Lookup(m.tagKey); ok {
			// 获取 tag 的第一部分并返回
			return strings.Split(tag, ",")[0]
		}
	}

	// 如果没有找到 tag, 则将字段名称转为 tag, 即字段名称的首字母小写
	name := []rune(f.Name)
	return fmt.Sprintf("%v%v", strings.ToLower(string(name[0])), string(f.Name[1:]))
}

// 将指定的值设置到指定的值反射对象中
//
// 参数:
//   - `v` (`reflect.Value`): 被设置值的值反射对象
//   - `val` (`any`): 要设置的值
//
// 返回:
//   - `err` (`error`): 错误对象
func (m *MapToStruct) assign(v reflect.Value, val any) (err error) {
	// 如果发生 panic 错误, 则进行恢复, 不中断程序
	defer func() {
		if e, ok := recover().(error); ok {
			// 设置错误对象
			err = e
		}
	}()

	// 如果两边类型一致, 则直接设置值
	if reflect.TypeOf(val) == v.Type() {
		v.Set(reflect.ValueOf(val).Convert(v.Type()))
	} else {
		// 根据值反射对象的类型选择不同的方式设置值
		switch v.Kind() {
		case reflect.Struct:
			// 对于值反射对象为结构体类型的, 要求 val 参数必须为 map 集合
			if err = m.assignToStruct(v, val.(map[string]any)); err != nil {
				return
			}
		case reflect.Slice:
			// 对于值反射对象为切片类型的, 要求 val 参数必须也为切片类型
			if err = m.assignToSlice(v, val.([]any)); err != nil {
				return
			}
		default:
			// 设置值
			v.Set(reflect.ValueOf(val).Convert(v.Type()))
		}
	}
	return
}

// 将指定的切片对象设置到指定的值反射对象中
//
// 参数:
//   - `v` (`reflect.Value`): 被设置值的值反射对象
//   - `s` (`[]any`): 要设置的切片类型对象
//
// 返回:
//   - `err` (`error`): 错误对象
func (m *MapToStruct) assignToSlice(v reflect.Value, s []any) error {
	// 为值反射对象设置切片长度
	sv := reflect.MakeSlice(v.Type(), len(s), len(s))

	// 将切片的每一项设置到值反射对象的每一项中
	for i := 0; i < sv.Len(); i++ {
		if err := m.assign(sv.Index(i), s[i]); err != nil {
			return err
		}
	}

	// 赋值切片
	v.Set(sv)
	return nil
}

// 将指定的结构体对象设置到指定的值反射对象中
//
// 参数:
//   - `v` (`reflect.Value`): 被设置值的值反射对象
//   - `data` (`map[string]any`): 要设置的结构体类型对象
//
// 返回:
//   - `err` (`error`): 错误对象
func (m *MapToStruct) assignToStruct(v reflect.Value, data map[string]any) error {
	// 获取值对象的类型
	t := v.Type()

	// 遍历值对象的各字段
	for i := 0; i < v.NumField(); i++ {
		// 获取指定字段的值对象
		fv := v.Field(i)

		// 如果字段值对象无效或者不允许设置值, 则忽略当前字段
		if !fv.IsValid() || !fv.CanSet() {
			continue
		}

		// 获取指定字段的字段类型
		ft := t.Field(i)
		// 查找字段的 tag
		tag := m.findTag(&ft)

		// 根据 tag 在 map 中查询目标值
		dv, ok := data[tag]
		if !ok {
			continue
		}

		// 根据字段类型查找预设的转换函数
		mapper, ok := m.mappers[ft.Type]
		if ok {
			var err error
			// 对目标值进行转换
			if dv, err = mapper(dv); err != nil {
				return err
			}
		}

		// 将目标值设置到结构体字段中
		if err := m.assign(fv, dv); err != nil {
			return fmt.Errorf("error field \"%v\": %v", ft.Name, err)
		}
	}
	return nil
}

// 将 map 对象内容解码到结构体对象中
//
// 参数:
//   - `data` (`map[string]any`): map 对象
//   - `target` (`any`): 结构体对象指针
//
// 返回:
//   - `err` (`error`): 错误对象
func (m *MapToStruct) Decode(data any, target any) error {
	tv := reflect.ValueOf(target)
	if tv.Kind() != reflect.Pointer {
		return fmt.Errorf("\"target\" argument must be a struct or slice pointer")
	}
	tv = tv.Elem()

	switch tv.Kind() {
	case reflect.Slice:
		return m.assignToSlice(tv, data.([]any))
	case reflect.Struct:
		return m.assignToStruct(tv, data.(map[string]any))
	default:
		return fmt.Errorf("\"target\" argument must be a struct pointer")
	}
}
