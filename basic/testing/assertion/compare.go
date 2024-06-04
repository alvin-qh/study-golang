package assertion

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	c "golang.org/x/exp/constraints"
)

// 断言 min <= val < max
func Between[V c.Ordered](t *testing.T, val V, min, max V) {
	if val < min || val >= max {
		assert.Failf(t, "value is not between min and max", "value: %v, min: %v, max: %v", val, min, max)
	}
}

// 断言切片中所有的元素值都和指定值一致
func All[V comparable](t *testing.T, val []V, target V) {
	for _, v := range val {
		assert.Equal(t, target, v)
	}
}

// 断言一个时长是否符合预期
func DurationMatch(t *testing.T, expect time.Duration, actual time.Duration) {
	var max time.Duration
	if expect == 0 {
		max = 1 + time.Duration(math.Ceil(float64(1)/float64(15*time.Millisecond))*15*float64(time.Millisecond))
	} else {
		max = expect + time.Duration(math.Ceil(float64(expect)/float64(15*time.Millisecond))*15*float64(time.Millisecond))
	}

	fmt.Printf("actual: %v, min: %v, max: %v\n", actual, expect, max)
	if actual < expect || actual > max {
		assert.Failf(t, "duration not match", "expect %v, actual %v", expect, actual)
	}
}
