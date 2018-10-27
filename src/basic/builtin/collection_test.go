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

	type Any interface{}
	array3 := [...]Any{"Hello", 1, false}
	assert.Equal(t, len(array3), 3)
	assert.Equal(t, array3[0], "Hello")
	assert.Equal(t, array3[2], false)

	array4 := [9][9]int{}
	assert.Equal(t, len(array4), 9)
	assert.Equal(t, len(array4[0]), 9)

	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			array4[i][j] = (i + 1) * (j + 1)
		}
	}
	assert.Equal(t, array4[2][6], 21)

	pa := &array4
	println((*pa)[0][0])
}
