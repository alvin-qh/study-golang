package structure

import (
	"errors"
	"reflect"
)

// 错误定义
var (
	ErrInvalidType       = errors.New("invalid type, need struct")
	ErrInvalidFieldName  = errors.New("no field found by given name")
	ErrInvalidMethodName = errors.New("no method found by given name")
	ErrInvalidTagName    = errors.New("no tag found by given name")
)

// 定义结构体反射类型
type Structure struct {
	rt     reflect.Type  // 结构体类型
	rv     reflect.Value // 结构体值
	fields []string      // 字段名称
}

// 创建结构体反射类型实例
//
// `t` 参数为一个实例或实例指针, 表示该实例的类型
func New(t any) (*Structure, error) {
	// 根据 `t` 参数获取其类型
	rt := reflect.TypeOf(t)
	if rt.Kind() == reflect.Pointer {
		// 如果 `typ` 参数是指针类型, 则获取其指向的类型
		rt = rt.Elem()
	}

	// 如果 `typ` 参数不是结构体类型, 则返回错误
	if rt.Kind() != reflect.Struct {
		return nil, ErrInvalidType
	}

	// 根据 `t` 参数获取其结构体值
	rv := reflect.ValueOf(t)
	if rv.Kind() == reflect.Pointer {
		// 如果 `t` 参数是指针类型, 则获取其指向的值
		rv = rv.Elem()
	}

	return &Structure{
		rt, // 设置结构体类型实例
		rv, // 设置结构体反射值实例
		nil,
	}, nil
}

// 获取结构体类型
func (s *Structure) Kind() reflect.Kind { return s.rt.Kind() }

// 获取结构体类型名称
func (s *Structure) Name() string { return s.rt.Name() }

// 获取结构体的包路径
func (s *Structure) PackagePath() string { return s.rt.PkgPath() }

// 根据所给的字段名称, 获取结构体字段对象
func (t *Structure) FindField(field string) (reflect.StructField, error) {
	f, ok := t.rt.FieldByName(field)
	if !ok {
		return f, ErrInvalidFieldName
	}
	return f, nil
}

// 根据所给的字段名称, 获取结构体字段值
//
// 通过结构体相关 `reflect.Value` 对象的 `FieldByName` 方法可以获取到结构体自动的反射值实例
//
//	rv := reflect.ValueOf(&user{}).Elem()
//	val := rv.FieldByName("Name").Interface()
func (t *Structure) GetFieldValue(field string) (any, error) {
	fv := t.rv.FieldByName(field)
	if !fv.IsValid() {
		return nil, ErrInvalidFieldName
	}

	return fv.Interface(), nil
}

// 根据所给的字段名称, 设置结构体字段值
//
// 返回指定字段之前的值
//
// 通过结构体相关 `reflect.Value` 对象的 `FieldByName` 方法可以获取到结构体自动的反射值实例
//
//	rv := reflect.ValueOf(&user{}).Elem()
//	rv.FieldByName("Name").Set(reflect.ValueOf("Tom"))
func (t *Structure) SetFieldValue(field string, vals any) (any, error) {
	fv := t.rv.FieldByName(field)
	if !fv.IsValid() {
		return nil, ErrInvalidFieldName
	}

	old := fv.Interface()
	fv.Set(reflect.ValueOf(vals))

	return old, nil
}

// 获取结构体所有
//
// 如果获取了结构体的反射类型 (`reflect.Type` 实例), 则可通过其 `NumField` 方法获取到结构体字段数量,
// 通过 `Field` 方法获取第 n 个字段类型实例, 并通过字段类型实例的 `Name` 属性获取字段名称
//
//	rt := reflect.TypeOf(&user{}).Elem()
//	n := rt.NumField()
//	name := rt.Field(0).Name
func (t *Structure) AllFieldNames() []string {
	if t.fields != nil {
		return t.fields
	}

	// 获取结构体字段数量
	n := t.rt.NumField()

	names := make([]string, 0, n)
	// 逐个获取结构体字段名称
	for i := range n {
		// 获取第 i 个字段类型实例, 并通过字段类型实例的 `Name` 属性获取字段名称
		f := t.rt.Field(i)
		names = append(names, f.Name)
	}

	t.fields = names
	return names
}

// 获取结构体指定字段上注解的标签值
//
// 通过结构体的反射类型 (`reflect.Type` 实例), 可以通过 `Tag` 属性获取到结构体字段上注解的标签值
//
//	rt := reflect.TypeOf(&user{}).Elem()
//	tag := rt.Field(0).Tag.Get("json")
//
// 除了 `Get` 方法外, `Lookup` 方法可以返回标签值以及标签是否存在
//
//	tagVal, ok := rType.Field(0).Tag.Lookup("json")
func (t *Structure) GetFieldTags(fields string, tag string) (string, error) {
	// 根据字段名获取字段类型实例
	f, ok := t.rt.FieldByName(fields)
	if !ok {
		return "", ErrInvalidFieldName
	}

	// 查找字段标签
	tagVal, ok := f.Tag.Lookup(tag)
	if !ok {
		return "", ErrInvalidTagName
	}

	return tagVal, nil
}

// 根据所给的方法名执行该方法, 并返回执行结果
//
// 通过 `reflect.Value` 实例的 `MethodByName` 方法即可获取指定名称的结构体方法实例, 调用 `Call` 方法即可执行该方法
//
// 注意, 有些方法是定义在结构体实例上的, 有些方法是定义在结构体实例指针上的, 因此获取结构体方法时要覆盖这两种情况
//
//	rv := reflect.ValueOf(&user{}).Elem()
//	res := rVal.MethodByName("GetName").Call()
//	res := rVal.Addr().MethodByName("GetName").Call()
func (t *Structure) CallMethodByName(method string, args ...any) ([]any, error) {
	m := t.rv.MethodByName(method)
	if !m.IsValid() {
		m = t.rv.Addr().MethodByName(method)
		if !m.IsValid() {
			return nil, ErrInvalidMethodName
		}
	}

	avs := make([]reflect.Value, 0, len(args))
	for _, arg := range args {
		avs = append(avs, reflect.ValueOf(arg))
	}

	res := m.Call(avs)

	retVals := make([]any, 0, len(res))
	for _, rv := range res {
		retVals = append(retVals, rv.Interface())
	}
	return retVals, nil
}
