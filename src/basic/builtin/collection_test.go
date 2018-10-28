package builtin

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateIntArray(t *testing.T) {
	array1 := [10]int{}
	assert.Equal(t, len(array1), 10)

	for i := 0; i < len(array1); i++ {
		array1[i] = i + 1
	}
	assert.Equal(t, array1[9], 10)

	for i, v := range array1 {
		assert.Equal(t, i, v-1)
	}

	array2 := [...]int{1, 2, 3}
	assert.Equal(t, len(array2), 3)

	array3 := [9][9]int{}
	assert.Equal(t, len(array3), 9)
	assert.Equal(t, len(array3[0]), 9)

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			array3[i][j] = (i + 1) * (j + 1)
		}
	}
	assert.Equal(t, array3[2][6], 21)
}

type Any interface{}

func TestAnyArray(t *testing.T) {
	array := [...]Any{"Hello", 1, false}
	assert.Equal(t, len(array), 3)
	assert.Equal(t, array[0], "Hello")
	assert.Equal(t, array[2], false)
}

func TestPointerOfArray(t *testing.T) {
	array := [...]Any{"Hello", 1, false}
	pa := &array // var pa *[3]Any = &array
	assert.Equal(t, (*pa)[0], "Hello")
}
