package maptostruct

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

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

func Decode(data map[string]any, target any) error {
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
