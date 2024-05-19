package core

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrInvalidPtr = errors.New("ptr must an address")
)

// 获取一个类型的完全限定名
//
// 完全限定名由 `<pkgPath>.<name>[<kind>]` 组成, 其中:
//   - `pkgPath` 表示类型所在包的路径, 系统内置类型为 `""`
//   - `Name` 为类型名称, `array`, `ptr`, `interface`, `slice`, `map` 为空
//   - `kind` 为 `reflect.Kind` 类型枚举, 通过 `String` 函数获取其字符串表达
func GetFullTypeName(t reflect.Type) string {
	return fmt.Sprintf("%v.%v[%v]", t.PkgPath(), t.Name(), t.Kind().String())
}

// 通过反射设置变量值
//
// 要对一个对象的值进行设置, 需要操作该对象的地址, 否则会报告"非地址类型异常", 即:
//   - `ptr`: 要设置的变量的指针
//   - `newVal`: 要设置的新值
func SetValueByReflect(ptr interface{}, newVal interface{}) (err error) {
	// ptr 参数转为 Value 类型
	tv := reflect.ValueOf(ptr)
	if tv.Kind() != reflect.Ptr { // 判断 ptr 参数是否为指针类型
		return ErrInvalidPtr
	}

	// 通过 Elem() 函数解引指针类型, 获取 ptr 指向的变量, 并设置新的值
	tv.Elem().Set(reflect.ValueOf(newVal))
	return nil
}
