package long

import (
	"strconv"
	"study/basic/oop/err"
)

// 定义 int64 为 Long 类型
// 相当于给 int64 类型一个别名
type Long int64

// 实现 typedef.Comparable 接口, 比较两个对象大小
func (i Long) Compare(other interface{}) int {
	val, ok := other.(Long)
	if !ok {
		panic(err.ErrType)
	}
	return int(i - val)
}

// ToString 方法, 将值转为 字符串
func (i Long) String() string {
	return strconv.FormatInt(int64(i), 10)
}
