package long

import (
	"basic/oop/typedef"
	"strconv"
)

type Long int64

func (i Long) Compare(other interface{}) int {
	val, ok := other.(Long)
	if !ok {
		panic(typedef.ErrType)
	}
	return int(i - val)
}

func (i Long) ToString() string {
	return strconv.FormatInt(int64(i), 10)
}
