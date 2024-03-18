package types

import "reflect"

type ValueType string

const (
	NAN  ValueType = ""
	I32  ValueType = "INT32"
	I64  ValueType = "INT64"
	F32  ValueType = "FLOAT32"
	F64  ValueType = "FLOAT64"
	BYTE ValueType = "BYTE"
)

type VT interface {
	~int | ~int32 | ~int64 | ~float32 | ~float64 | ~byte
}

type Value[T VT] struct {
	value T
}

func NewValue[T VT](value T) Value[T] {
	return Value[T]{value}
}

func ExplainValue(v interface{}) (ValueType, interface{}) {
	if reflect.ValueOf(v).Kind() == reflect.Ptr {
		switch n := v.(type) {
		case *Value[int32]:
			return I32, n.value
		case *Value[int]:
			return I32, n.value
		case *Value[int64]:
			return I64, n.value
		case *Value[float32]:
			return F32, n.value
		case *Value[float64]:
			return F64, n.value
		case *Value[byte]:
			return BYTE, n.value
		}
	} else {
		switch n := v.(type) {
		case Value[int32]:
			return I32, n.value
		case Value[int]:
			return I32, n.value
		case Value[int64]:
			return I64, n.value
		case Value[float32]:
			return F32, n.value
		case Value[float64]:
			return F64, n.value
		case Value[byte]:
			return BYTE, n.value
		}
	}
	return NAN, nil
}
