package conv

import (
	"fmt"
	"strconv"
)

// 将任意类型变量转为对应的字符串
//
// 参数:
//   - `val` (`any`): 表示一个任意类型的任意参数
//
// 返回:
//   - 参数 `val` 转换得到的字符串
func AnyToString(val any) string {
	// 判断参数类型, 并根据不同类型进行对应的转换操作
	switch sval := val.(type) {
	case string:
		return sval
	case int:
		return strconv.FormatInt(int64(sval), 10)
	case int8:
		return strconv.FormatInt(int64(sval), 10)
	case int16:
		return strconv.FormatInt(int64(sval), 10)
	case int32:
		return strconv.FormatInt(int64(sval), 10)
	case int64:
		return strconv.FormatInt(sval, 10)
	case uint:
		return strconv.FormatUint(uint64(sval), 10)
	case uint8:
		return strconv.FormatUint(uint64(sval), 10)
	case uint16:
		return strconv.FormatUint(uint64(sval), 10)
	case uint32:
		return strconv.FormatUint(uint64(sval), 10)
	case uint64:
		return strconv.FormatUint(sval, 10)
	case float32:
		return strconv.FormatFloat(float64(sval), 'g', 5, 32)
	case float64:
		return strconv.FormatFloat(float64(sval), 'g', 11, 64)
	case bool:
		return strconv.FormatBool(sval)
	default:
		// 如果变量类型实现了 `fmt.Stringer` 接口, 则调用接口的 `String` 方法进行转换
		if sval, ok := val.(fmt.Stringer); ok {
			return sval.String()
		}

		// 使用字符串格式化方式进行转换
		return fmt.Sprintf("%v", sval)
	}
}
