package structure_test

import (
	"fmt"
	"reflect"
	"study/basic/builtin/reflects/structure"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 表示用户的结构体
type User struct {
	Id     int    `primaryKey:"true" null:"false"`
	Name   string `default:"Alvin"`
	Gender rune
}

// 将结构体转为字符串
func (u *User) String(n int) string {
	return fmt.Sprintf(`user{Id: %v, Name: "%v", Gender: %v, Num: %v}`, u.Id, u.Name, u.Gender, n)
}

// 测试获取结构体字段标签
func TestStructure_GetFieldTags(t *testing.T) {
	u := User{}

	s, err := structure.New(&u)
	assert.Nil(t, err)

	// 获取结构体字段标签值
	tag, err := s.GetFieldTags("Id", "primaryKey")
	assert.Nil(t, err)
	assert.Equal(t, tag, "true")

	tag, err = s.GetFieldTags("Id", "null")
	assert.Nil(t, err)
	assert.Equal(t, tag, "false")
}

// 测试获取结构体类型
func TestStructure_Kind(t *testing.T) {
	u := User{}

	s, err := structure.New(&u)
	assert.Nil(t, err)

	// 获取结构体类型
	k := s.Kind()
	assert.Equal(t, reflect.Struct, k)
}

// 测试获取结构体名称
func TestStructure_Name(t *testing.T) {
	u := User{}

	s, err := structure.New(&u)
	assert.Nil(t, err)

	// 获取结构体名称
	n := s.Name()
	assert.Equal(t, "User", n)
}

// 测试获取结构体包路径
func TestStructure_PackagePath(t *testing.T) {
	u := User{}

	s, err := structure.New(&u)
	assert.Nil(t, err)

	// 获取结构体名称
	path := s.PackagePath()
	assert.Equal(t, "study/basic/builtin/reflects/structure_test", path)
}

// 测试获取结构体字段实例
func TestStructure_FindField(t *testing.T) {
	u := User{}

	s, err := structure.New(&u)
	assert.Nil(t, err)

	// 获取结构体字段实例
	f, err := s.FindField("Id")
	assert.Nil(t, err)
	assert.Equal(t, f.Name, "Id")
	assert.Equal(t, reflect.Int, f.Type.Kind())
}

// 测试获取结构体字段值
func TestStructure_GetFieldValue(t *testing.T) {
	u := User{}

	s, err := structure.New(&u)
	assert.Nil(t, err)

	// 获取结构体字段实例
	u.Id = 100
	v, err := s.GetFieldValue("Id")
	assert.Nil(t, err)
	assert.Equal(t, 100, v)
}

// 测试设置结构体字段值
func TestStructure_SetFieldValue(t *testing.T) {
	u := User{}

	s, err := structure.New(&u)
	assert.Nil(t, err)

	// 获取结构体字段实例
	ov, err := s.SetFieldValue("Id", 100)
	assert.Nil(t, err)
	assert.Equal(t, 0, ov)
	assert.Equal(t, 100, u.Id)
}

// 测试获取结构体所有字段名称
func TestStructure_AllFieldNames(t *testing.T) {
	u := User{}

	s, err := structure.New(&u)
	assert.Nil(t, err)

	// 获取结构体字段实例
	ns := s.AllFieldNames()
	assert.Equal(t, []string{"Id", "Name", "Gender"}, ns)
}

// 测试根据方法名称调用方法
func TestStructure_CallMethodByName(t *testing.T) {
	u := User{
		Id:     100,
		Name:   "Alvin",
		Gender: 'M',
	}

	s, err := structure.New(&u)
	assert.Nil(t, err)

	u.Name = "Alvin"

	// 调用 `User` 结构体的 `String` 方法
	r, err := s.CallMethodByName("String", 100)
	assert.Nil(t, err)

	assert.Equal(t, r[0], "user{Id: 100, Name: \"Alvin\", Gender: 77, Num: 100}")
}
