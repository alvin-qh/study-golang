package utils

import (
	"fmt"
	"strconv"
)

func AnyToString(val any) string {
	switch sval := val.(type) {
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
		return strconv.FormatFloat(float64(sval), 'f', 5, 32)
    case float64:
		return strconv.FormatFloat(float64(sval), 'f', 11, 64)
	case bool:
		return strconv.FormatBool(sval)
	case int:
		return strconv.FormatInt(int64(sval), 10)
	case int:
		return strconv.FormatInt(int64(sval), 10)
	case int:
		return strconv.FormatInt(int64(sval), 10)
	case int:
		return strconv.FormatInt(int64(sval), 10)
	case int:
		return strconv.FormatInt(int64(sval), 10)
	case int:
		return strconv.FormatInt(int64(sval), 10)
	default:
		return fmt.Sprintf("%v", sval)
	}
}
