package builtin

import (
	"container/list"
	"fmt"
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

func TestIntsFromArray(t *testing.T) {
	array := [...]int{1, 2, 3, 4, 5}

	ints := Ints(array[:])
	assert.Equal(t, len(ints), len(array))

	ints = Ints(array[0:])
	assert.Equal(t, len(ints), len(array))

	ints = Ints(array[1:2])
	assert.Equal(t, len(ints), 1)
	assert.Equal(t, ints[0], 2)

	ints = Ints(array[1:3])
	assert.Equal(t, len(ints), 2)
	assert.Equal(t, ints[0], 2)
	assert.Equal(t, ints[1], 3)

	ints = Ints(array[:3])
	assert.Equal(t, len(ints), 3)
	assert.Equal(t, ints[0], 1)
	assert.Equal(t, ints[1], 2)
	assert.Equal(t, ints[2], 3)
}

func TestSliceGrowUp(t *testing.T) {
	var ints []int
	assert.Equal(t, len(ints), 0)
	assert.Equal(t, cap(ints), 0)

	c := 1
	for i := 0; i < 20; i++ {
		if len(ints) > 0 && len(ints) == cap(ints) {
			c = len(ints) * 2
		}

		ints = append(ints, i)
		assert.Equal(t, len(ints), i+1)
		assert.Equal(t, cap(ints), c)
	}
}

func TestSliceShare(t *testing.T) {
	a := []int{1, 2, 3}
	b := a
	assert.Equal(t, &a, &b)

	a[1] = 100
	assert.Equal(t, &a, &b)
	assert.Equal(t, b[1], 100)

	a = append(a, 200)
	assert.NotEqual(t, &a, &b)
}

func TestNewInts(t *testing.T) {
	ints := NewInts(10, 100)
	assert.Equal(t, len(ints), 10)
	assert.Equal(t, cap(ints), 100)
}

func TestInts_Append(t *testing.T) {
	ints := Ints{1, 2, 3, 4}
	ints.Append(5)
	assert.Equal(t, ints[4], 5)
	assert.Equal(t, len(ints), 5)
}

func TestInts_Remove(t *testing.T) {
	ints := Ints{1, 2, 3, 4}
	ints.Remove(2)
	assert.Equal(t, ints[2], 4)
	assert.Equal(t, len(ints), 3)

	ints.Remove(0)
	assert.Equal(t, ints[1], 4)
	assert.Equal(t, len(ints), 2)

	ints.Remove(len(ints) - 1)
	assert.Equal(t, ints[0], 2)
	assert.Equal(t, len(ints), 1)
}

func TestInts_Clear(t *testing.T) {
	ints := Ints{1, 2, 3, 4}
	ints.Clear()

	assert.Equal(t, len(ints), 0)
}

func TestInts_Size(t *testing.T) {
	ints := Ints{1, 2, 3, 4}
	assert.Equal(t, ints.Size(), 4)

	ints.Clear()
	assert.Equal(t, ints.Size(), 0)
}

func TestList(t *testing.T) {
	lst := list.New()
	lst.PushBack(1)
	assert.Equal(t, lst.Len(), 1)

	lst.PushBack("Hello")
	assert.Equal(t, lst.Len(), 2)

	type Any interface{}

	array := make([]Any, lst.Len())
	for iter, i := lst.Front(), 0; iter != nil; iter, i = iter.Next(), i+1 {
		array[i] = iter.Value
	}
	assert.Equal(t, len(array), 2)
}

func TestListAt(t *testing.T) {
	lst := ListAssign(1, 2, 3, "Hello")
	assert.Equal(t, ListAt(lst, 0), 1)
	assert.Equal(t, ListAt(lst, 1), 2)
	assert.Equal(t, ListAt(lst, 3), "Hello")

	fn := func() {
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}

	defer fn()
	ListAt(lst, 4)
}

func TestListToSlice(t *testing.T) {
	lst := ListAssign(1, 2, 3, "Hello")
	array := ListToSlice(lst)
	assert.Equal(t, len(array), 4)
	assert.ElementsMatch(t, array, []Any{1, 2, 3, "Hello"})
}
