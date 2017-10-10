package basic

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestPointer(t *testing.T) {
	assert.Equal(t, Int, *PInt, "expect *PInt == 10")

	assert.Equal(t, &PInt, PPInt, "expect &PInt == PPInt")
	assert.Equal(t, Int, **PPInt, "expect **PInt == 10")
	assert.Equal(t, Int, **PPInt, "expect **PInt == 10")
}

func TestVarPointer(t *testing.T) {
	a, b := 1, 1

	VarPointer(a, &b)

	assert.Equal(t, 1, a, "expect a == 1")
	assert.Equal(t, 100, b, "expect b == 100")
}

func TestPointerWithArray(t *testing.T) {
	assert.Equal(t, Ints[0], (*PInts)[0], `expect Ints[0] == (*PInts)[0]`)
}

func TestArrayPointer(t *testing.T) {
	a := []int{1, 2, 3}
	b := []int{1, 2, 3}

	ArrayPointer(a, &b)

	assert.Equal(t,100, a[0], "expect a[0] == 100")
	assert.Equal(t,100, b[0], "expect b[0] == 100")
}

func TestExchangeByPointer(t *testing.T) {
	a, b := 1, 2
	ExchangeByPointer(&a, &b)

	assert.Equal(t, 2, a, `expect a == 2`)
	assert.Equal(t, 1, b, `expect b == 1`)
}

func TestStructPointer(t *testing.T) {
	p1 := Person{"Alvin", 36, Male}
	p2 := Person{"Emma", 30, Female}

	StructPointer(p1, &p2)

	assert.Equal(t, 36, p1.Age, `expect p1.Age == 36`)
	assert.Equal(t, 20, p2.Age, `expect p2.Age == 20`)
}

