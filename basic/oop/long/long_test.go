package long

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompare(t *testing.T) {
	l1 := Long(100)
	assert.True(t, l1.Compare(l1) == 0)

	l2 := Long(200)
	assert.True(t, l1.Compare(l2) < 0)

	l3 := Long(50)
	assert.True(t, l1.Compare(l3) > 0)

	l4 := l1
	assert.True(t, l1.Compare(l4) == 0)
}

func TestToString(t *testing.T) {
	var l Long
	assert.Equal(t, "0", l.String())

	l = Long(100)
	assert.Equal(t, "100", l.String())
}
