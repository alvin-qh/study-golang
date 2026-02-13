package optargs_test

import (
	"study/basic/builtin/functions/optargs"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试函数可选参数
//
// 可以通过不定参数 + 回调函数模拟可选参数
func TestFunction_OptionalArgs(t *testing.T) {
	u := optargs.CreateUser(
		optargs.WithUserId(2),
		optargs.WithUserName("Emma"),
		optargs.WithUserGender('F'),
	)

	assert.Equal(t, 2, u.Id)
	assert.Equal(t, "Emma", u.Name)
	assert.Equal(t, 'F', u.Gender)
}
