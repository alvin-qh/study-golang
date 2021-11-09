package user

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddPhoneToUser(t *testing.T) {
	pu := New(1, "Alvin", "alvin@fake.com", []string{"13999912345", "13000056789"})
	assert.IsType(t, &User{}, pu)

	pu.AddPhone("13456789111")
	assert.Equal(t, []string{"13999912345", "13000056789", "13456789111"}, pu.Phone)
}
