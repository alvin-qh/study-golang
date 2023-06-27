package reflect

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	ErrInvalidPtr = errors.New("ptr must an address")
)

// 表示用户的结构体
type User struct {
	Id     int    `primaryKey:"true" null:"false"`
	Name   string `default:"Alvin"`
	Gender rune
}

// 用于测试反射的 User 对象行为
// 该函数的主体是 *User
func (u *User) String() string {
	return fmt.Sprintf("%v(%v)-%v", u.Name, u.Id, string(u.Gender))
}

// 用于测试反射的 User 对象行为
// 该函数的主体是 User
func (u User) AsString() string {
	return fmt.Sprintf("%v(%v)-%v", u.Name, u.Id, string(u.Gender))
}

// 获取一个类型的完全限定名，由 <pkgPath>.<name>[<kind>] 组成
// pkgPath 表示类型所在包的路径，系统内置类型为 ""
// name 为类型名称，array, ptr, interface, slice, map 为空
// kind 为 Kind 类型枚举，通过 String 函数获取其字符串表达
func GetFullTypeName(t reflect.Type) string {
	return fmt.Sprintf("%v.%v[%v]", t.PkgPath(), t.Name(), t.Kind().String())
}

// 通过反射设置值
// 要对一个对象的值进行设置，需要操作该对象的地址，否则会报告"非地址类型异常"

// 通过反射设置变量值
//
//	ptr: 要设置的变量的指针
//	newVal: 要设置的新值
func SetValueByReflect(ptr interface{}, newVal interface{}) (err error) {
	// ptr 参数转为 Value 类型
	tv := reflect.ValueOf(ptr)
	if tv.Kind() != reflect.Ptr { // 判断 ptr 参数是否为指针类型
		return ErrInvalidPtr
	}

	// 通过 Elem() 函数解引指针类型，获取 ptr 指向的变量，并设置新的值
	tv.Elem().Set(reflect.ValueOf(newVal))

	return nil
}

// 通过反射设置结构体字段值
//
//	ptr: 要设置字段的结构体变量指针
//	field: 要设置的结构体字段名
//	newVal: 要设置的字段新值
func SetStructFieldByReflect(ptr interface{}, field string, newVal interface{}) (err error) {
	// ptr 参数转为 Value 类型
	tv := reflect.ValueOf(ptr)
	if tv.Kind() != reflect.Ptr { // 判断 ptr 参数是否为指针类型
		return ErrInvalidPtr
	}

	// 通过 Elem() 函数解引指针类型，获取 ptr 指向的结构体对象
	// 通过名称反射字段，并设置新值
	tv.Elem().FieldByName(field).Set(reflect.ValueOf(newVal))
	return nil
}

// 测试函数反射的简单加法函数
func Add(a, b int) (r int) {
	r = a + b
	return
}
