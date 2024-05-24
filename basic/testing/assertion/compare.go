package assertion

import (
	"testing"

	"github.com/stretchr/testify/assert"
	c "golang.org/x/exp/constraints"
)

// 断言 min <= val < max
func Between[V c.Ordered](t *testing.T, val V, min, max V) {
	if val < min || val >= max {
		assert.Fail(t, "value is not between min and max", "value: %v, min: %v, max: %v", val, min, max)
	}
}

// 断言切片中所有的元素值都和指定值一致
func All[V comparable](t *testing.T, val []V, target V) {
    for _, v := range val {
        assert.Equal(t, target, v)
    }
}
