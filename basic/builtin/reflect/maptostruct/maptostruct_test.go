package maptostruct

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// 测试用的类型, 表示性别
type Gender string

// 定义表示男女性别的常量
const (
	GenderM = Gender("M")
	GenderF = Gender("F")
)

// 定义测试用的结构体, 表示一个用户, 用 `json` 作为 tag key
type user struct {
	Id       int       `json:"id"`
	Name     string    `json:"name"`
	Gender   Gender    `json:"gender"`
	Birthday time.Time `json:"birthday"`
	Titles   []string  `json:"titles"`
	Locale   struct {
		Country  string `json:"country"`
		Language string `json:"language"`
	} `json:"locale"`
}

// 测试查找结构体字段的 tag
func TestFindTag(t *testing.T) {
	// 找到 user 结构体对象的 Id 字段
	f, _ := reflect.TypeOf(new(user)).Elem().FieldByName("Id")

	// 实例化 MapToStruct 对象, 以 json 为 tag key
	mts := New("json")
	// 确认字段标签
	assert.Equal(t, "id", mts.findTag(&f))

	// 实例化 MapToStruct 对象, 设置不存在的 tag key
	mts = New("unknown")
	// 确认字段标签
	assert.Equal(t, "id", mts.findTag(&f))
}

// 测试当 target 参数不正确时, Decode 方法返回预期的错误信息
func TestDecodeByInvalidTarget(t *testing.T) {
	mts := New("json")

	v := 123
	m := map[string]any{}

	// 传递错误的 target 参数 v
	err := mts.Decode(m, &v)
	// 确认返回的错误信息符合预期
	assert.EqualError(t, err, "\"target\" argument must be a struct pointer")
}

// 测试将 map 解码到结构体对象中
func TestDecodeStruct(t *testing.T) {
	m := map[string]any{
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

	// 实例化 MapToStruct 对象
	mts := New("json")

	// 创建目标结构体, 获取其地址
	u := new(user)

	// 将 map 解码到 user 结构体中
	err := mts.Decode(m, u)
	assert.NoError(t, err)

	// 确认结构体内容符合预期
	assert.Equal(t, 100, u.Id)
	assert.Equal(t, "Alvin", u.Name)
	assert.Equal(t, GenderM, u.Gender)
	assert.Equal(t, time.Date(1981, 3, 17, 0, 0, 0, 0, time.UTC), u.Birthday)
	assert.Equal(t, []string{"Manager", "Engineer"}, u.Titles)
	assert.Equal(t, "China", u.Locale.Country)
	assert.Equal(t, "Chinese", u.Locale.Language)
}

// 测试将 map 切片解码到结构体切片中
func TestDecodeSlice(t *testing.T) {
	m := []any{
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
	mts := New("json")

	// 创建目标切片变量, 获取其地址
	var us []user

	err := mts.Decode(m, &us)
	assert.NoError(t, err)
	assert.Len(t, us, 2)

	// 确认切片内容符合预期
	assert.Equal(t, 100, us[0].Id)
	assert.Equal(t, "Alvin", us[0].Name)
	assert.Equal(t, GenderM, us[0].Gender)
	assert.Equal(t, time.Date(1981, 3, 17, 0, 0, 0, 0, time.UTC), us[0].Birthday)
	assert.Equal(t, []string{"Manager", "Engineer"}, us[0].Titles)
	assert.Equal(t, "China", us[0].Locale.Country)
	assert.Equal(t, "Chinese", us[0].Locale.Language)

	assert.Equal(t, 101, us[1].Id)
	assert.Equal(t, "Emma", us[1].Name)
	assert.Equal(t, GenderF, us[1].Gender)
	assert.Equal(t, time.Date(1985, 3, 29, 0, 0, 0, 0, time.UTC), us[1].Birthday)
	assert.Equal(t, []string{"Manager", "Engineer"}, us[1].Titles)
	assert.Equal(t, "China", us[1].Locale.Country)
	assert.Equal(t, "Chinese", us[1].Locale.Language)
}
