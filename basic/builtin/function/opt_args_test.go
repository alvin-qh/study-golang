package function

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试所有参数为默认的情形
func TestOptionArgsNone(t *testing.T) {
	u := CreateUser()
	assert.Equal(t, 1, u.Id)
	assert.Equal(t, "Alvin", u.Name)
	assert.Equal(t, 'M', u.Gender)
}

// 测试仅设置 `UserId` 参数的情形
func TestOptionArgsId(t *testing.T) {
	u := CreateUser(WithUserId(2))
	assert.Equal(t, 2, u.Id)
	assert.Equal(t, "Alvin", u.Name)
	assert.Equal(t, 'M', u.Gender)
}

// 测试同时设置 `UserId`, `UserName` 参数的情形
func TestOptionArgsIdName(t *testing.T) {
	u := CreateUser(WithUserId(2), WithUserName("Emma"))
	assert.Equal(t, 2, u.Id)
	assert.Equal(t, "Emma", u.Name)
	assert.Equal(t, 'M', u.Gender)
}

// 测试同时设置 `UserId`, `UserGender` 参数的情形
func TestOptionArgsIdGender(t *testing.T) {
	u := CreateUser(WithUserId(2), WithUserGender('F'))
	assert.Equal(t, 2, u.Id)
	assert.Equal(t, "Alvin", u.Name)
	assert.Equal(t, 'F', u.Gender)
}
