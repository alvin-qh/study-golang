package oop

import (
	"basic/oop/errors"
	"basic/oop/size"
	"basic/oop/size3d"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试 typedef.Comparable 接口进行对象比较
func TestInterface(t *testing.T) {
	// 比较两个对象的大小
	s1, s2 := size.New(10, 20), size.New(11, 21)
	s11 := s1

	assert.True(t, Eq(s1, s1))

	assert.True(t, Eq(s1, s11))
	assert.True(t, Eq(s11, s1))

	assert.True(t, Ne(s1, s2))
	assert.True(t, Ne(s2, s1))

	assert.True(t, Gt(s2, s1))
	assert.True(t, Lt(s1, s2))

	assert.True(t, Ge(s2, s1))
	assert.True(t, Ge(s1, s1))
	assert.True(t, Ge(s1, s11))

	assert.True(t, Le(s1, s2))
	assert.True(t, Le(s1, s1))
	assert.True(t, Le(s1, s11))

	// 测试比较不同类型对象时, 出现的 panic 异常
	defer func() {
		err := recover().(error)
		assert.ErrorIs(t, err, errors.ErrType)
	}()

	ss1 := size3d.New(10, 20, 30)
	Eq(ss1, s1)
	assert.Fail(t, "Cannot run here")
}
