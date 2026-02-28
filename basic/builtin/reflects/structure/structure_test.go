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
func (u *User) String() string {
	return fmt.Sprintf(`user{Id: %v, Name: "%v", Gender: %v}`, u.Id, u.Name, u.Gender)
}

// 测试获取结构体类型
func TestStructure_Kind(t *testing.T) {
	// 创建结构体对象
	var u User

	// 基于结构体对象创建 structure.Structure 对象, 确认创建成功
	s, err := structure.New(&u)
	assert.Nil(t, err)

	// 获取结构体类型
	k := s.Kind()
	assert.Equal(t, reflect.Struct, k)
}

// 测试获取结构体名称
func TestStructure_Name(t *testing.T) {
	// 创建结构体对象
	var u User

	// 基于结构体对象创建 structure.Structure 对象, 确认创建成功
	s, err := structure.New(&u)
	assert.Nil(t, err)

	// 获取结构体名称
	n := s.Name()
	assert.Equal(t, "User", n)
}

// 测试获取结构体包路径
func TestStructure_PackagePath(t *testing.T) {
	// 创建结构体对象
	var u User

	// 基于结构体对象创建 structure.Structure 对象, 确认创建成功
	s, err := structure.New(&u)
	assert.Nil(t, err)

	// 获取结构体名称
	path := s.PackagePath()
	assert.Equal(t, "study/basic/builtin/reflects/structure_test", path)
}

// 测试获取结构体字段实例
func TestStructure_FindField(t *testing.T) {
	// 创建结构体对象
	var u User

	// 基于结构体对象创建 structure.Structure 对象, 确认创建成功
	s, err := structure.New(&u)
	assert.Nil(t, err)

	// 获取结构体字段实例, 确认获取的字段示例符合预期 (字段名, 字段类型)
	f, ok := s.FindField("Id")
	assert.True(t, ok)
	assert.Equal(t, f.Name, "Id")
	assert.Equal(t, reflect.Int, f.Type.Kind())
}

// 测试获取结构体字段值
func TestStructure_GetFieldValue(t *testing.T) {
	// 创建结构体对象
	u := User{
		Id:     100,
		Name:   "Alvin",
		Gender: 'M',
	}

	// 基于结构体对象创建 structure.Structure 对象, 确认创建成功
	s, err := structure.New(&u)
	assert.Nil(t, err)

	// 获取结构体字段实例
	v, err := s.GetFieldValue("Id")
	assert.Nil(t, err)
	assert.Equal(t, 100, v)
}

// 测试设置结构体字段值
func TestStructure_SetFieldValue(t *testing.T) {
	// 创建结构体对象
	var u User

	// 基于结构体对象创建 structure.Structure 对象, 确认创建成功
	s, err := structure.New(&u)
	assert.Nil(t, err)

	// 为结构体 Id 字段设置值, 确认设置成功以及原始值符合预期
	ov, err := s.SetFieldValue("Id", 100)
	assert.Nil(t, err)
	assert.Equal(t, 0, ov)

	// 为结构体 Name 字段设置值, 确认设置成功以及原始值符合预期
	ov, err = s.SetFieldValue("Name", "Alvin")
	assert.Nil(t, err)
	assert.Equal(t, "", ov)

	// 为结构体 Gender 字段设置值, 确认设置成功以及原始值符合预期
	ov, err = s.SetFieldValue("Gender", rune('M'))
	assert.Nil(t, err)
	assert.Equal(t, rune(0), ov)

	// 确认结构体字段值均被设置并符合预期
	assert.Equal(t, User{100, "Alvin", rune('M')}, u)
}

// 测试获取结构体所有字段名称
func TestStructure_AllFieldNames(t *testing.T) {
	// 创建结构体对象
	u := User{}

	// 基于结构体对象创建 structure.Structure 对象, 确认创建成功
	s, err := structure.New(&u)
	assert.Nil(t, err)

	// 获取结构体字段所有字段名称
	ns := s.AllFieldNames()
	assert.Equal(t, []string{"Id", "Name", "Gender"}, ns)
}

// 测试根据方法名称调用方法
func TestStructure_CallMethodByName(t *testing.T) {
	// 创建结构体实例
	u := User{
		Id:     100,
		Name:   "Alvin",
		Gender: 'M',
	}

	// 获取结构体实例的元数据对象
	s, err := structure.New(&u)
	assert.Nil(t, err)

	// 为结构体字段设置值
	u.Name = "Alvin"

	// 通过反射调用 `User` 结构体的 `String` 方法, 确认返回值符合预期
	r, err := s.CallMethodByName("String", 100)
	assert.Nil(t, err)
	assert.IsType(t, string(""), r)
	assert.Equal(t, r, "user{Id: 100, Name: \"Alvin\", Gender: 77, Num: 100}")
}

// 测试获取结构体字段标签
func TestStructure_GetFieldTags(t *testing.T) {
	// 创建结构体实例
	u := User{}

	// 基于结构体对象创建 structure.Structure 对象, 确认创建成功
	s, err := structure.New(&u)
	assert.Nil(t, err)

	// 根据 Id 字段和标签 Key 为 primaryKey 对应的标签值, 确认获取成功以及获取的标签值正确
	tag, err := s.GetFieldTags("Id", "primaryKey")
	assert.Nil(t, err)
	assert.Equal(t, tag, "true")

	// 根据 Id 字段和标签 Key 为 null 对应的标签值, 确认获取成功以及获取的标签值正确
	tag, err = s.GetFieldTags("Id", "null")
	assert.Nil(t, err)
	assert.Equal(t, tag, "false")

	// 根据 Name 字段和标签 Key 为 default 对应的标签值, 确认获取成功以及获取的标签值正确
	tag, err = s.GetFieldTags("Name", "default")
	assert.Nil(t, err)
	assert.Equal(t, tag, "Alvin")

	// 根据 Gender 字段和标签 Key 为 default 对应的标签值, 该标签不存在, 故确认返回的错误信息符合预期
	tag, err = s.GetFieldTags("Gender", "default")
	assert.Equal(t, "no tag found by given tag key", err.Error())
}
