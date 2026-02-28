package structure

import (
	"errors"
	"reflect"
	"study/basic/builtin/slices"
)

// 错误定义
var (
	// 参数为无效指针错误
	ErrInvalidPtr = errors.New("need pointer of struct object")

	// 参数为无效结构体错误
	ErrInvalidType = errors.New("invalid type, need struct")

	// 找不到对应结构体字段错误
	ErrInvalidFieldName = errors.New("no field found by given field name")

	// 找不到对应结构体方法错误
	ErrInvalidMethodName = errors.New("no method found by given method name")

	// 找不到对应结构体标签错误
	ErrInvalidTagName = errors.New("no tag found by given tag key")
)

// 定义结构体反射类型
type Structure struct {
	// 结构体类型
	typ reflect.Type

	// 结构体值
	val reflect.Value

	// 结构体中所有字段名称
	fNames []string
}

// 创建结构体反射类型实例
//
// 参数:
//   - `structObj`: 任意结构体实例指针
//
// 返回:
//   - `*Structure`: `Structure` 类型实例指针
//   - `error`: 如果发生错误, 则返回错误对象
func New(structObj any) (*Structure, error) {
	// 获取 structObj 参数的反射值实例
	val := reflect.ValueOf(structObj)

	// 如果 structObj 参数不是指针类型, 则返回错误
	if val.Kind() != reflect.Pointer {
		return nil, ErrInvalidPtr
	}

	// 对反射值进行解引, 获取其指针指向的变量的反射值
	val = val.Elem()

	// 如果反射值表示的不是结构体类型, 则返回错误
	if val.Kind() != reflect.Struct {
		return nil, ErrInvalidType
	}

	// 获取反射值对应的反射类型对象
	typ := val.Type()
	if typ.Kind() != reflect.Struct {
		return nil, ErrInvalidType
	}

	// 创建结构体反射类型实例并返回
	return &Structure{
		typ, // 设置结构体类型实例
		val, // 设置结构体反射值实例
		nil,
	}, nil
}

// 获取结构体类型
//
// 返回:
//   - `reflect.Kind`: 结构体类型
func (s *Structure) Kind() reflect.Kind { return s.typ.Kind() }

// 获取结构体类型名称
//
// 返回:
//   - `string`: 结构体类型名称
func (s *Structure) Name() string { return s.typ.Name() }

// 获取结构体的包路径
//
// 返回:
//   - `string`: 包路径
func (s *Structure) PackagePath() string { return s.typ.PkgPath() }

// 根据所给的字段名称, 获取结构体字段对象
//
// 参数:
//   - `fieldName`: 字段名称
//
// 返回:
//   - `reflect.StructField`: 结构体字段反射对象
//   - `bool`: 如果字段存在, 则返回 `true`, 否则返回 `false`
func (s *Structure) FindField(fieldName string) (reflect.StructField, bool) {
	return s.typ.FieldByNameFunc(func(name string) bool { return name == fieldName })
}

// 根据所给的字段名称, 获取结构体字段值
//
// 参数:
//   - `fieldName`: 字段名称
//
// 错误:
//   - `error`: 如果发生错误, 则返回错误对象
//
// 返回:
//   - `any`: 字段值
func (s *Structure) GetFieldValue(fieldName string) (any, error) {
	// 通过结构体反射值对象获取结构体字段反射值对象
	fVal := s.val.FieldByName(fieldName)
	if !fVal.IsValid() {
		return nil, ErrInvalidFieldName
	}

	// 将字段值转为 any 类型后返回
	return fVal.Interface(), nil
}

// 根据所给的字段名称, 设置结构体字段值
//
// 参数:
//   - `fieldName`: 字段名称
//   - `value`: 要设置的字段值
//
// 返回:
//   - `any`: 字段原始值
//   - `error`: 如果发生错误, 则返回错误对象
func (s *Structure) SetFieldValue(fieldName string, value any) (any, error) {
	// 通过结构体反射值对象获取结构体字段反射值对象
	fVal := s.val.FieldByName(fieldName)
	if !fVal.IsValid() {
		return nil, ErrInvalidFieldName
	}

	// 获取字段原始值
	old := fVal.Interface()

	// 设置新字段值
	fVal.Set(reflect.ValueOf(value))

	// 返回字段原始值
	return old, nil
}

// 获取结构体所有字段名称
//
// 返回:
//   - `[]string`: 切片对象, 包括结构体下所有字段名称
func (s *Structure) AllFieldNames() []string {
	if s.fNames != nil {
		return s.fNames
	}

	// 初始化切片用于存储字段名称
	names := make([]string, 0, s.typ.NumField())

	// 遍历结构体字段, 将名称保存在切片中
	for field := range s.typ.Fields() {
		names = append(names, field.Name)
	}

	// 将名称切片进行缓存
	s.fNames = names

	// 返回字段名称切片
	return names
}

// 获取结构体指定字段上注解的标签值
//
// 参数:
//   - `fieldName`: 字段名称
//   - `tagKey`: 标签名称的 Key 值
//
// 返回:
//   - `string`: 标签值
//   - `error`: 如果发生错误, 则返回错误对象
func (s *Structure) GetFieldTags(fieldName string, tagKey string) (string, error) {
	// 根据字段名获取字段类型实例
	fType, ok := s.typ.FieldByName(fieldName)
	if !ok {
		return "", ErrInvalidFieldName
	}

	// 根据标签 Key 获取字段标签值
	tagVal, ok := fType.Tag.Lookup(tagKey)
	if !ok {
		return "", ErrInvalidTagName
	}

	// 返回字段标签值
	return tagVal, nil
}

// 根据所给的方法名执行该方法, 并返回执行结果
//
// 参数:
//   - `methodName`: 方法名称
//   - `args`: 方法参数列表
//
// 返回:
//   - `any`: 调用方法的返回值
//   - `error`: 如果发生错误, 则返回错误对象
func (s *Structure) CallMethodByName(methodName string, args ...any) ([]any, error) {
	// 通过方法名获取方法反射实例, 需要通过结构体反射值的地址进行获取
	m := s.val.Addr().MethodByName(methodName)
	if !m.IsValid() {
		return nil, ErrInvalidMethodName
	}

	// 准备函数调用参数列表, 参数列表为 `reflect.Value` 类型的切片
	avs := slices.Map(args, func(v any) reflect.Value {
		return reflect.ValueOf(v)
	})

	// 调用方法, 获取返回值切片, 内容为函数返回值的反射值类型
	res := m.Call(avs)

	// 将返回值切片转为实际的值并返回
	return slices.Map(res, func(v reflect.Value) any {
		return v.Interface()
	}), nil
}
