package structmapper

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

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
//	   Name string `struct:"name"`
//	   Age  int    `struct:"age"`
//	}
//
// 表示 `User` 结构体的 `Name` 字段将被转换为 Map 的 `name` Key 值, `Age` 字段将被转换为 Map 的 `age` Key 值
type StructMapper struct {
	// 结构体字段注解
	tag string

	// 指定类型的转换函数
	mapperFns map[reflect.Type]MapperFn
}

// 创建实例
//
// 参数:
//   - `tag` 默认为空, 表示使用结构体的字段名称作为 Map 的 Key 值, 如果需要使用自定义的 tag, 则需要指定该参数
func New(tag string) *StructMapper {
	// 返回 StructMapper 结构体指针
	return &StructMapper{
		tag: tag,
		mapperFns: map[reflect.Type]MapperFn{
			// 设置默认的类型转换函数
			reflect.TypeFor[time.Time](): mapTime,
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
func (sm *StructMapper) AddMapperFn(t reflect.Type, mapperFn MapperFn) {
	sm.mapperFns[t] = mapperFn
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
func (sm *StructMapper) FindTag(f *reflect.StructField) string {
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

func (sm *StructMapper) fromSlice(v reflect.Value) ([]any, error) {
	// 验证参数, 如果参数 v 不是切片类型, 则返回错误
	if v.Kind() != reflect.Slice {
		return nil, fmt.Errorf("\"obj\" argument must be a slice")
	}

	// 获取切片的长度
	l := v.Len()

	// 创建结果切片, 长度为参数 v 表示的切片的长度
	r := make([]any, l)

	// 遍历切片的每一项, 将其设置到结果切片中
	for i := range l {
		// 获取切片的每一项, 并将对应值添加到结果切片中
		r[i] = v.Index(i).Interface()
	}

	return r, nil
}
func (sm *StructMapper) fromStruct(v reflect.Value) (map[string]any, error) {
	// 验证参数, 如果参数 v 不是结构体类型, 则返回错误
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("\"obj\" argument must be a struct")
	}

	// 获取结构体的类型
	t := v.Type()

	// 创建结果 map 对象
	m := make(map[string]any)

	// 遍历结构体的字段
	for f := range t.Fields() {

	}
}

func (sm *StructMapper) Encode(obj any) (map[string]any, error) {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
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

	return r, nil
}

// 将指定的结构体对象设置到指定的值反射对象中
//
// 参数:
//   - `v` 值反射对象
//   - `data` 要设置的结构体对象
//
// 返回:
//   - `error` 错误对象, 如果没有错误, 则返回 `nil`
func (sm *StructMapper) assignToStruct(v reflect.Value, data map[string]any) error {
	// 获取值对象的类型
	t := v.Type()

	// 遍历值对象的各字段
	for f := range t.Fields() {
		// 查找字段的 tag
		tag := sm.FindTag(&f)

		// 根据 tag 在 map 中查询目标值
		dv, ok := data[tag]
		if !ok {
			continue
		}

		// 根据字段类型查找预设的转换函数
		mapper, ok := sm.mapperFns[f.Type]
		if ok {
			// 对目标值进行转换
			var err error
			if dv, err = mapper(dv); err != nil {
				return err
			}
		}

		// 获取结构体字段的值反射对象
		fv := v.Field(f.Index[0])

		// 将目标值设置到结构体字段中
		if err := sm.assign(fv, dv); err != nil {
			return fmt.Errorf("error field \"%v\": %v", f.Name, err)
		}
	}
	return nil
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
func (sm *StructMapper) Decode(data any, target any) error {
	// 获取目标对象的值反射对象
	rv := reflect.ValueOf(target)

	// 确认目标对象必须为指针类型, 否则返回错误
	if rv.Kind() != reflect.Pointer {
		return fmt.Errorf("\"target\" argument must be a struct or slice pointer")
	}

	// 获取指针指向的实例值反射对象
	rv = rv.Elem()

	var err error

	// 根据实例值反射对象的类型, 选择不同的方式设置值
	switch rv.Kind() {
	case reflect.Slice: // 当对象是切片类型时, 将 Map 中的切片内容设置到目标对象中
		err = sm.assignToSlice(rv, data.([]any))
	case reflect.Struct: // 当对象是结构体类型时, 将 Map 中的内容设置到目标对象中
		err = sm.assignToStruct(rv, data.(map[string]any))
	default: // 其它类型, 返回错误
		err = fmt.Errorf("\"target\" argument must be a struct pointer")
	}
	return err
}

// 将指定的值设置到指定的值反射对象中
//
// 参数:
//   - `v` 值反射对象
//   - `val` 要设置的值
//
// 返回:
//   - `error` 返回错误对象, 如果没有错误, 则返回 nil
func (sm *StructMapper) assign(v reflect.Value, val any) (err error) {
	// 如果发生 panic 错误, 则进行恢复, 不中断程序
	defer func() {
		if r := recover(); r != nil {
			// 设置错误对象
			var ok bool
			if err, ok = r.(error); !ok {
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
			// 设置默认值
			v.Set(reflect.ValueOf(val).Convert(v.Type()))
		}
	}
	return
}

// 将指定的切片对象设置到指定的值反射对象中
//
// 参数:
//   - `v` 值反射对象
//   - `s` 要设置的切片对象
func (sm *StructMapper) assignToSlice(v reflect.Value, s []any) error {
	// 创建一个新的切片对象, 类型与值反射对象的类型一致, 长度和容量与所给切片一致
	sv := reflect.MakeSlice(v.Type(), len(s), len(s))

	// 将切片的每一项设置到值反射对象的每一项中
	for i := range sv.Len() {
		// 将
		if err := sm.assign(sv.Index(i), s[i]); err != nil {
			return err
		}
	}

	// 赋值切片
	v.Set(sv)
	return nil
}
