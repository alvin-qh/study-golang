package structmapper_test

import (
	"study/basic/builtin/reflects/structmapper"
	"study/basic/builtin/reflects/structure"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试用的类型, 表示性别
type Gender string

const (
	// 定义表示男女性别的常量
	GenderM = Gender("M")
	GenderF = Gender("F")
)

// 测试用结构体, 以 `struct` 作为字段标签
type User struct {
	Id       int       `struct:"id"`
	Name     string    `struct:"name"`
	Gender   Gender    `struct:"gender"`
	Birthday time.Time `struct:"birthday"`
	Titles   []string  `struct:"titles"`
	Locale   struct {
		Country  string `struct:"country"`
		Language string `struct:"language"`
	} `struct:"locale"`
}

// 测试查找结构体字段标签
func TestStructMapper_findTag(t *testing.T) {
	// 生成结构体反射工具类型对象
	s, err := structure.New(&User{})
	assert.Nil(t, err)

	// 找到 user 结构体对象的 Id 字段
	f, err := s.FindField("Id")
	assert.Nil(t, err)

	// 实例化 MapToStruct 对象, 以 json 为 tag key
	m := structmapper.New("struct")
	assert.Equal(t, "id", m.FindTag(&f))

	// 实例化 MapToStruct 对象, 设置不存在的 tag key
	m = structmapper.New("unknown")
	assert.Equal(t, "id", m.FindTag(&f))
}

// 测试 `Decode` 方法错误参数
//
// 测试当 `target` 参数不正确时, `Decode` 方法返回预期的错误信息
func TestStructMapper_DecodeWrongTarget(t *testing.T) {
	m := structmapper.New("struct")

	v := 123

	// 传递错误的 `target` 参数, 返回预期错误
	err := m.Decode(map[string]any{}, &v)
	assert.EqualError(t, err, "\"target\" argument must be a struct pointer")
}

// 测试将 Map 的 Value 值填入结构体字段中
func TestStructMapper_DecodeStruct(t *testing.T) {
	data := map[string]any{
		"id":       100,
		"name":     "Alvin",
		"gender":   "M",
		"birthday": "1981-03-17",
		"titles":   []string{"Manager", "Engineer"},
		"locale": map[string]any{
			"country":  "China",
			"language": "Chinese",
		},
	}

	// 创建 StructMapper 实例
	m := structmapper.New("struct")

	// 创建目标结构体, 获取其地址
	u := new(User)

	// 将 map 解码到 user 结构体中
	err := m.Decode(data, u)
	assert.Nil(t, err)

	// 确认结构体内容符合预期
	assert.Equal(t, 100, u.Id)
	assert.Equal(t, "Alvin", u.Name)
	assert.Equal(t, GenderM, u.Gender)
	assert.Equal(t, time.Date(1981, 3, 17, 0, 0, 0, 0, time.UTC), u.Birthday)
	assert.Equal(t, []string{"Manager", "Engineer"}, u.Titles)
	assert.Equal(t, "China", u.Locale.Country)
	assert.Equal(t, "Chinese", u.Locale.Language)
}

// 测试将切片解码到结构体切片中
func TestStructMapper_DecodeSlice(t *testing.T) {
	s := []any{
		map[string]any{
			"id":       100,
			"name":     "Alvin",
			"gender":   "M",
			"birthday": "1981-03-17",
			"titles":   []string{"Manager", "Engineer"},
			"locale": map[string]any{
				"country":  "China",
				"language": "Chinese",
			},
		},
		map[string]any{
			"id":       101,
			"name":     "Emma",
			"gender":   "F",
			"birthday": "1985-03-29",
			"titles":   []string{"Manager", "Engineer"},
			"locale": map[string]any{
				"country":  "China",
				"language": "Chinese",
			},
		},
	}

	// 实例化 MapToStruct 对象
	m := structmapper.New("struct")

	// 创建目标切片变量, 获取其地址
	var u []User

	err := m.Decode(s, &u)
	assert.Nil(t, err)
	assert.Len(t, u, 2)

	// 确认切片内容符合预期
	assert.Equal(t, 100, u[0].Id)
	assert.Equal(t, "Alvin", u[0].Name)
	assert.Equal(t, GenderM, u[0].Gender)
	assert.Equal(t, time.Date(1981, 3, 17, 0, 0, 0, 0, time.UTC), u[0].Birthday)
	assert.Equal(t, []string{"Manager", "Engineer"}, u[0].Titles)
	assert.Equal(t, "China", u[0].Locale.Country)
	assert.Equal(t, "Chinese", u[0].Locale.Language)

	assert.Equal(t, 101, u[1].Id)
	assert.Equal(t, "Emma", u[1].Name)
	assert.Equal(t, GenderF, u[1].Gender)
	assert.Equal(t, time.Date(1985, 3, 29, 0, 0, 0, 0, time.UTC), u[1].Birthday)
	assert.Equal(t, []string{"Manager", "Engineer"}, u[1].Titles)
	assert.Equal(t, "China", u[1].Locale.Country)
	assert.Equal(t, "Chinese", u[1].Locale.Language)
}
