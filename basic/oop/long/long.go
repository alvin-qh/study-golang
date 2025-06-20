package long

import (
	errs "basic/oop/errors"
	"strconv"
)

// 定义 int64 为 Long 类型, 相当于给 int64 类型一个别名
//
// 可以为类型定义方法, 格式为:
//
//	func (receiver_type) func_name([parameter_list]) [return_types] {...}
//
// 其中类型的方法会传入类型实例或类型实例的指针, 称为 Receiver
//   - 类型 Long 包含全部 Receiver 为 Long 的方法
//   - 类型 *Long 包含全部 Receiver 为 Long + *Long 的方法
type Long int64

// 实现 typedef.Comparable 接口, 比较两个对象大小
func (i Long) Compare(other interface{}) int {
	val, ok := other.(Long)
	if !ok {
		panic(errs.ErrInvalidType)
	}
	return int(i - val)
}

// ToString 方法, 将值转为字符串
func (i Long) String() string {
	return strconv.FormatInt(int64(i), 10)
}
