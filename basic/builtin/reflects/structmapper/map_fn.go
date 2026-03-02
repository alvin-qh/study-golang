package structmapper

import "time"

// 将所给类型转换为 `time.Time` 类型
//
// 该函数作为默认的时间转换函数, 传入 `MapToStruct` 的 `mappers` 字段中
func mapTime(obj any) (r any, err error) {
	// 根据值的不同类型进行不同的转换
	switch val := obj.(type) {
	case string:
		// 将字符串表示的时间转为 time.Time 类型
		r, err = time.Parse(time.RFC3339, val)
		if err != nil {
			r, err = time.Parse(time.TimeOnly, val)
			if err != nil {
				r, err = time.Parse(time.DateOnly, val)
			}
		}

		if err != nil {
			return
		}
	case float64:
		// 将 float64 表示的时间转为 time.Time 类型
		r = time.Unix(0, int64(val)*int64(time.Millisecond))
	case int64:
		// 将 int64 表示的时间转为 time.Time 类型
		r = time.Unix(0, val*int64(time.Millisecond))
	default:
		// 其它类型, 不做任何转换
		r = val
	}
	return
}
