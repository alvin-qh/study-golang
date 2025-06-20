package maptostruct

import "time"

// 将所给类型转换为 `time.Time` 类型
func mapTime(v any) (r any, err error) {
	// 根据值的不同类型进行不同的转换
	switch _v := v.(type) {
	case string:
		// 将字符串表示的时间转为 time.Time 类型
		r, err = time.Parse(time.RFC3339, _v)
		if err != nil {
			r, err = time.Parse(time.TimeOnly, _v)
			if err != nil {
				r, err = time.Parse(time.DateOnly, _v)
			}
		}
		if err != nil {
			return
		}
	case float64:
		// 将 float64 表示的时间转为 time.Time 类型
		r = time.Unix(0, int64(v.(float64))*int64(time.Millisecond))
	case int64:
		// 将 int64 表示的时间转为 time.Time 类型
		r = time.Unix(0, v.(int64)*int64(time.Millisecond))
	default:
		// 其它类型, 不做任何转换
		r = v
	}
	return
}
