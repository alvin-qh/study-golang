package structmapper

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

// 将所给类型转换为 `time.Time` 类型
//
// 该函数作为默认的时间转换函数, 传入 `MapToStruct` 的 `mappers` 字段中
func mapTime(v any) (r any, err error) {
	// 根据值的不同类型进行不同的转换
	switch _v := v.(type) {
	case string:
		// 将字符串表示的时间转为 time.Time 类型
		r, err = time.Parse(time.RFC3339, _v)
		if err != nil {
			r, err = time.Parse(time.TimeOnly, _v)
			if err != nil {
				r, err = time.Parse(time.DateOnly, _v)
			}
		}

		if err != nil {
			return
		}
	case float64:
		// 将 float64 表示的时间转为 time.Time 类型
		r = time.Unix(0, int64(v.(float64))*int64(time.Millisecond))
	case int64:
		// 将 int64 表示的时间转为 time.Time 类型
		r = time.Unix(0, v.(int64)*int64(time.Millisecond))
	default:
		// 其它类型, 不做任何转换
		r = v
	}
	return
}

// 定义转换函数类型
//
// 该函数用于将一个值转为另一个类型值
type MapperFn func(v any) (r any, err error)

// 结构体转换类型
//
// 该类型可以将结构体实例和 Map 实例相互转换
//
// 需要指定结构体字段的 tag, 并依据 tag 指定的名称作为 Map 的 Key 值, 例如:
//
//	type User struct {
//	  Name string `struct:"name"`
//	  Age  int    `struct:"age"`
//	}
//
// 表示 `User` 结构体的 `Name` 字段将被转换为 Map 的 `name` Key 值, `Age` 字段将被转换为 Map 的 `age` Key 值
type StructMapper struct {
	// 结构体字段注解
	tag string

	// 指定类型的转换函数
	mappers map[reflect.Type]MapperFn
}

// 创建实例
//
// 该函数传入标签名称, 返回 `MapToStruct` 实例
func New(tag string) *StructMapper {
	return &StructMapper{
		tag: tag,
		mappers: map[reflect.Type]MapperFn{
			// 设置默认的类型转换函数
			reflect.TypeOf(time.Time{}): mapTime,
		},
	}
}

// 为指定类型添加转换函数
//
// 该转换函数用于将特定类型字段和 Map 的 Value 值进行转换
//
// `t` 参数用于指定要进行转换的字段类型, 对于没有指定的类型, 将使用 Go 语言的默认规则进行转换
func (sm *StructMapper) AddMapper(t reflect.Type, mapper MapperFn) {
	sm.mappers[t] = mapper
}

// 从结构体字段获取注解的 tag
func (sm *StructMapper) findTag(f *reflect.StructField) string {
	if len(sm.tag) > 0 {
		// 查找预设 tag key 对应的 tag
		if tag, ok := f.Tag.Lookup(sm.tag); ok {
			// 获取 tag 的第一部分并返回
			return strings.Split(tag, ",")[0]
		}
	}

	// 如果没有找到 tag, 则将字段名称转为 tag, 即字段名称的首字母小写
	name := []rune(f.Name)
	return fmt.Sprintf("%v%v", strings.ToLower(string(name[0])), string(f.Name[1:]))
}

// 将指定的值设置到指定的值反射对象中
func (sm *StructMapper) assign(v reflect.Value, val any) (err error) {
	// 如果发生 panic 错误, 则进行恢复, 不中断程序
	defer func() {
		if r := recover(); r != nil {
			// 设置错误对象
			if e, ok := r.(error); ok {
				err = e
			} else {
				err = fmt.Errorf("unknown error: %v", r)
			}
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
			if err = sm.assignToStruct(v, val.(map[string]any)); err != nil {
				return
			}
		case reflect.Slice:
			// 对于值反射对象为切片类型的, 要求 val 参数必须也为切片类型
			if err = sm.assignToSlice(v, val.([]any)); err != nil {
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
func (sm *StructMapper) assignToSlice(v reflect.Value, s []any) (err error) {
	// 为值反射对象设置切片长度
	sv := reflect.MakeSlice(v.Type(), len(s), len(s))

	// 将切片的每一项设置到值反射对象的每一项中
	for i := 0; i < sv.Len(); i++ {
		if err = sm.assign(sv.Index(i), s[i]); err != nil {
			return
		}
	}

	// 赋值切片
	v.Set(sv)
	return
}

// 将指定的结构体对象设置到指定的值反射对象中
func (sm *StructMapper) assignToStruct(v reflect.Value, data map[string]any) (err error) {
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
		tag := sm.findTag(&ft)

		// 根据 tag 在 map 中查询目标值
		dv, ok := data[tag]
		if !ok {
			continue
		}

		// 根据字段类型查找预设的转换函数
		mapper, ok := sm.mappers[ft.Type]
		if ok {
			// 对目标值进行转换
			if dv, err = mapper(dv); err != nil {
				return
			}
		}

		// 将目标值设置到结构体字段中
		if err = sm.assign(fv, dv); err != nil {
			err = fmt.Errorf("error field \"%v\": %v", ft.Name, err)
		}
	}
	return
}

// 将 Map 对象内容解码到结构体对象中
func (m *StructMapper) Decode(data any, target any) (err error) {
	tv := reflect.ValueOf(target)
	if tv.Kind() != reflect.Pointer {
		err = fmt.Errorf("\"target\" argument must be a struct or slice pointer")
		return
	}

	tv = tv.Elem()

	switch tv.Kind() {
	case reflect.Slice:
		err = m.assignToSlice(tv, data.([]any))
	case reflect.Struct:
		err = m.assignToStruct(tv, data.(map[string]any))
	default:
		err = fmt.Errorf("\"target\" argument must be a struct pointer")
	}
	return
}
