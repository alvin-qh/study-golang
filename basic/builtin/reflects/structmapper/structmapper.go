package structmapper

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Converter interface {
	From(src any) (any, error)
	To(value any) (any, error)
}

// 结构体转换类型
//
// 该类型可以将结构体实例和 Map 实例相互转换
//
// 需要指定结构体字段的 tag, 并依据 tag 指定的名称作为 Map 的 Key 值, 例如:
//
//	type User struct {
//	   Name string `struct:"name"`
//	   Age  int    `struct:"age"`
//	}
//
// 表示 `User` 结构体的 `Name` 字段将被转换为 Map 的 `name` Key 值, `Age` 字段将被转换为 Map 的 `age` Key 值
type StructMapper struct {
	// 结构体字段注解
	tag string

	// 指定类型的转换函数
	convs map[reflect.Type]Converter
}

// 创建实例
//
// 参数:
//   - `tag` 默认为空, 表示使用结构体的字段名称作为 Map 的 Key 值, 如果需要使用自定义的 tag, 则需要指定该参数
func New(tag string) *StructMapper {
	// 返回 StructMapper 结构体指针
	return &StructMapper{
		tag: tag,
		convs: map[reflect.Type]Converter{
			// 设置默认的类型转换函数
			reflect.TypeFor[time.Time](): timeConverter{},
		},
	}
}

// 为指定类型添加转换函数
//
// 该转换函数用于将特定类型字段和 Map 的 Value 值进行转换
//
// 参数:
//   - `t` 参数用于指定要进行转换的字段类型, 对于没有指定的类型, 将使用 Go 语言的默认规则进行转换
//   - `mapperFn` 参数用于指定转换函数, 该函数接收一个参数, 该参数为待转换的值, 返回一个值和错误对象
func (sm *StructMapper) AddMapperFn(t reflect.Type, conv Converter) {
	sm.convs[t] = conv
}

// 将字段名称转为 tag key
//
// 参数:
//   - `fieldName` 输入的字段名称
//
// 返回:
//   - `string` 返回的字段名称转为的 tag key
func fieldNameToTagKey(fieldName string) string {
	// 获取字段名称的第一个字符
	firstChar := []rune(fieldName)[0]

	// 判断字段名称的第一个字符是否大写
	if firstChar >= 'A' && firstChar <= 'Z' {
		// 如果大写, 则将字段名称的第一个字符转换为小写
		return strings.ToLower(string(firstChar)) + fieldName[1:]
	}

	// 返回字段名称
	return string(firstChar) + fieldName[1:]
}

// 从结构体字段获取注解的 tag
//
// 参数:
//   - `f` 输入的字段反射对象
//
// 返回:
//   - `string` 返回字段对应的 tag 值
func (sm *StructMapper) fieldName(f *reflect.StructField) string {
	// 判断是否设置了 tag key 值
	if len(sm.tag) > 0 {
		// 查找预设 tag key 对应的 tag 值
		if tag, ok := f.Tag.Lookup(sm.tag); ok {
			// 获取 tag 的第一部分并返回
			return strings.Split(tag, ",")[0]
		}
	}

	// 如果没有找到 tag key, 则将字段名称转为 tag key, 即字段名称的首字母小写并返回
	return fieldNameToTagKey(f.Name)
}

func (sm *StructMapper) from(src reflect.Value) (any, error) {

}

func (sm *StructMapper) fromSlice(src reflect.Value) []any {
	// 获取切片的长度
	vLen := src.Len()

	// 创建结果切片, 长度为参数 v 表示的切片的长度
	result := make([]any, vLen)

	// 遍历切片的每一项, 将其设置到结果切片中
	for i := range vLen {
		// 获取切片的每一项, 并将对应值添加到结果切片中
		result[i] = src.Index(i).Interface()
	}

	// 返回结果切片
	return result
}

func (sm *StructMapper) fromStruct(src reflect.Value) (map[string]any, error) {
	// 获取结构体的类型
	typ := src.Type()

	// 创建结果 map 对象
	result := make(map[string]any)

	// 遍历结构体的字段
	for field := range typ.Fields() {
		switch field.Type.Kind() {
		case reflect.Struct:
			// 如果字段类型为结构体, 则将字段值转为 map[string]any 对象
			if v, err := sm.fromStruct(f.Value); err != nil {
				return nil, err
			} else {
				// 将字段值转为 map[string]any 对象并设置到结果 map 中
				result[sm.fieldName(f)] = v
			}
		case reflect.Slice:
		}
	}
}

func (sm *StructMapper) Encode(src any) (map[string]any, error) {
	sv := reflect.ValueOf(src)
	if sv.Kind() == reflect.Pointer {
		sv = sv.Elem()
	}

	var r map[string]any

	switch v.Kind() {
	case reflect.Slice:
		return sm.fromSlice(v)
	case reflect.Struct:
		return sm.fromStruct(v)
	default:
		return nil, fmt.Errorf("\"obj\" argument must be a struct or slice")
	}
}
func (sm *StructMapper) from(src reflect.Value) (result any, err error) {
	if src.Kind() == reflect.Pointer {
		src = src.Elem()
	}

	switch src.Kind() {
	case reflect.Slice:
		result = sm.fromSlice(src)
	case reflect.Struct:
		result, err = sm.fromStruct(src)
	default:
		return nil, fmt.Errorf("\"obj\" argument must be a struct or slice")
	}

	return
}

// 将指定的结构体对象设置到指定的值反射对象中
//
// 参数:
//   - `v` 值反射对象
//   - `data` 要设置的结构体对象
//
// 返回:
//   - `error` 错误对象, 如果没有错误, 则返回 `nil`
func (sm *StructMapper) assignToStruct(target reflect.Value, value map[string]any) error {
	// 获取值对象的类型
	t := target.Type()

	// 遍历值对象的各字段
	for f := range t.Fields() {
		// 查找字段的 tag
		tag := sm.fieldName(&f)

		// 根据 tag 在 map 中查询目标值
		dv, ok := value[tag]
		if !ok {
			continue
		}

		// 获取结构体字段的值反射对象
		fv := target.Field(f.Index[0])

		// 将目标值设置到结构体字段中
		if err := sm.assign(fv, dv); err != nil {
			return fmt.Errorf("error field \"%v\": %v", f.Name, err)
		}
	}
	return nil
}

// 将指定的切片对象设置到指定的值反射对象中
//
// 参数:
//   - `v` 值反射对象
//   - `s` 要设置的切片对象
//
// 返回:
//   - `error` 错误对象, 如果没有错误则返回 `nil`
func (sm *StructMapper) assignToSlice(target reflect.Value, value []any) error {
	// 获取切片值长度
	vLen := len(value)

	// 创建一个新的切片对象, 类型与值反射对象的类型一致, 长度和容量与所给切片一致
	sv := reflect.MakeSlice(target.Type(), vLen, vLen)

	// 将切片的每一项设置到值反射对象的每一项中
	for i := range vLen {
		// 将
		if err := sm.assign(sv.Index(i), value[i]); err != nil {
			return err
		}
	}

	// 赋值切片
	target.Set(sv)
	return nil
}

// 将指定的值设置到指定的值反射对象中
//
// 参数:
//   - `v` 值反射对象
//   - `val` 要设置的值
//
// 返回:
//   - `error` 返回错误对象, 如果没有错误, 则返回 nil
func (sm *StructMapper) assign(target reflect.Value, value any) (err error) {
	// 如果发生 panic 错误, 则进行恢复, 不中断程序
	defer func() {
		// 从 panic 恢复, 并获取错误值
		if r := recover(); r != nil {
			// 确认错误值是否为 error 类型, 如果是 error 类型, 则返回错误,
			// 否则创建一个 error 对象并返回
			var ok bool
			if err, ok = r.(error); !ok {
				err = fmt.Errorf("unknown error: %v", r)
			}
		}
	}()

	// 确认目标对象是否为指针
	if target.Kind() == reflect.Pointer {
		// 确认目标对象是否为空, 如果为空则返回
		if target.IsNil() {
			if value != nil {
				return fmt.Errorf("invalid target, nil pointer")
			}
		}

		// 如果值为 nil, 则将值反射对象设置为 nil
		if value == nil {
			target.Set(reflect.Zero(target.Type()))
			return
		}

		// 获取目标对象的元素类型
		target = target.Elem()
	}

	// 根据字段类型查找预设的转换函数
	conv, ok := sm.convs[target.Type()]
	if ok {
		// 对目标值进行转换
		var err error
		if value, err = conv.To(value); err != nil {
			return err
		}
	}

	switch target.Kind() {
	case reflect.Slice:
		err = sm.assignToSlice(target, value.([]any))
	case reflect.Struct:
		err = sm.assignToStruct(target, value.(map[string]any))
	default:
		target.Set(reflect.ValueOf(value).Convert(target.Type()))
	}

	return
}

// 将 Map 对象内容解码到结构体对象中
//
// 该方法将 data 参数表示的 Map 对象键值对解码到 target 参数表示的结构体对象中
//
// 参数:
//   - `data` 输入的 Map 对象
//   - `target` 目标对象的指针
//
// 错误:
//   - `error` 错误对象, 如果没有错误, 则返回 `nil`
func (sm *StructMapper) Decode(target any, data map[string]any) error {
	// 获取目标对象的值反射对象
	rv := reflect.ValueOf(target)

	// 确认目标对象必须为指针类型, 否则返回错误
	if rv.Kind() != reflect.Pointer {
		return fmt.Errorf("\"target\" argument must be a struct or slice pointer")
	}

	// 将 data 值赋值给 rv 反射值
	return sm.assign(rv.Elem(), data)
}
