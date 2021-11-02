package types

import (
	"fmt"
	"reflect"
)

// 获取一个类型的完全限定名，由 <pkgPath>.<name>[<kind>] 组成
// pkgPath 表示类型所在包的路径，系统内置类型为 ""
// name 为类型名称，array, ptr, interface, slice, map 为空
// kind 为 Kind 类型枚举，通过 String 函数获取其字符串表达
func GetFullTypeName(t reflect.Type) string {
	return fmt.Sprintf("%v.%v[%v]", t.PkgPath(), t.Name(), t.Kind().String())
}
