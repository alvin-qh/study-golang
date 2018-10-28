package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReturn(t *testing.T) {
	one := ReturnOne()
	assert.Equal(t, one, "First")

	one, two := ReturnTwo()
	assert.Equal(t, one, "First")
	assert.Equal(t, two, "Second")
}

func TestReturnFunc(t *testing.T) {
	fun := ReturnFunc() // var fun func(a int, b int) int = ReturnFunc()
	assert.Equal(t, fun(10, 20), 30)

	pFun := &fun
	assert.Equal(t, (*pFun)(10, 20), 30)

	sFun := fun
	assert.Equal(t, sFun(10, 20), 30)
}

func TestArgumentAsFunction(t *testing.T) {
	f := func(a interface{}, b interface{}) interface{} {
		return a.(int) + b.(int)
	}

	result := ArgumentAsFunction(10, 20, f)
	assert.Equal(t, result, 30)
}

func TestNamedReturnValue(t *testing.T) {
	s := NamedReturnValue("Hello")
	assert.Equal(t, s, "Hello")
}
