package structure

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 表示用户的结构体
type user struct {
	Id     int    `primaryKey:"true" null:"false"`
	Name   string `default:"Alvin"`
	Gender rune
}

// 将结构体转为字符串
func (u *user) String(n int) string {
	return fmt.Sprintf(`user{Id: %v, Name: "%v", Gender: %v, Num: %v}`, u.Id, u.Name, u.Gender, n)
}

// 测试 `Structure` 类型
//
// 通过 `Structure` 类型实例可以获取结构体类型的反射信息
func TestStructure(t *testing.T) {
	u := user{}

	stu, err := New(&u)
	assert.NoError(t, err)

	// 确认类型为结构体类型
	assert.Equal(t, stu.Kind(), reflect.Struct)
	// 确认类型名称
	assert.Equal(t, stu.Name(), "user")
	// 确认类型所在包名称
	assert.Equal(t, stu.PackagePath(), "study-golang/basic/builtin/reflect/structure")

	// 获取结构体指定名称的字段类型
	f, err := stu.FindField("Id")
	assert.NoError(t, err)
	assert.Equal(t, f.Name, "Id")

	// 根据字段名称设置结构体字段值
	old, err := stu.SetFieldValue("Id", 100)
	assert.NoError(t, err)
	assert.Equal(t, old, 0)
	assert.Equal(t, u.Id, 100)

	// 根据字段名获取字段值
	val, err := stu.GetFieldValue("Id")
	assert.NoError(t, err)
	assert.Equal(t, val, 100)

	// 获取结构体所有字段值
	names := stu.AllFieldNames()
	assert.Equal(t, 3, len(names))

	// 获取结构体字段标签值
	tag, err := stu.GetFieldTags("Id", "primaryKey")
	assert.NoError(t, err)
	assert.Equal(t, tag, "true")

	u.Name = "Alvin"
	// 根据方法名调用结构体方法
	res, err := stu.CallMethodByName("String", 100)
	assert.NoError(t, err)
	assert.Equal(t, res[0], "user{Id: 100, Name: \"Alvin\", Gender: 0, Num: 100}")
}
