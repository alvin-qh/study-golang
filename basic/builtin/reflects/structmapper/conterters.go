package structmapper

import (
	"fmt"
	"time"
)

type timeConverter struct{}

func (timeConverter) From(value any) (result any, err error) {
	if tim, ok := value.(time.Time); ok {
		result = tim.Format(time.RFC3339)
	} else {
		err = fmt.Errorf("value is not time.Time")
	}
	return
}

func (timeConverter) To(value any) (result any, err error) {
	// 根据值的不同类型进行不同的转换
	switch val := value.(type) {
	case string:
		// 将字符串表示的时间转为 time.Time 类型
		result, err = time.Parse(time.RFC3339, val)
		if err != nil {
			result, err = time.Parse(time.TimeOnly, val)
			if err != nil {
				result, err = time.Parse(time.DateOnly, val)
			}
		}

		if err != nil {
			return
		}
	case float64:
		// 将 float64 表示的时间转为 time.Time 类型
		result = time.Unix(0, int64(val)*int64(time.Millisecond))
	case int64:
		// 将 int64 表示的时间转为 time.Time 类型
		result = time.Unix(0, val*int64(time.Millisecond))
	default:
		// 其它类型, 不做任何转换
		result = val
	}
	return
}
