package serial

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	JSON_FROM_MAP = `{"email":"alvin@fake.com","id":1,"name":"Alvin","phone":["13999912345","13000056789"]}`

	JSON = `{"id":1,"name":"Alvin","email":"alvin@fake.com","phone":["13999912345","13000056789"]}`

	JSON_FORMATTED = `{
  "id": 1,
  "name": "Alvin",
  "email": "alvin@fake.com",
  "phone": [
    "13999912345",
    "13000056789"
  ]
}`
)

// 将 map 对象转为 json，即将 map 的 key 作为 json 的字段名，value 作为 json 的字段值

// 将 结构体 编码为 json，或者从 json 中还原 结构体 的字段值
// 要求：结构体必须为公开，且结构体字段为公开，私有的结构体字段会被忽略
// 如果需要对 json 字段做进一步描述，需要在结构体中使用 tag 进行标注，参考 user.User 结构体

// 序列化：从 map 对象生成 json 字符串
func TestJsonMarshalFromMap(t *testing.T) {
	// 产生一个 map 对象，key 为 string，value 任意
	m := map[string]interface{}{
		"id":    1,
		"name":  "Alvin",
		"email": "alvin@fake.com",
		"phone": []string{
			"13999912345",
			"13000056789",
		},
	}

	var s string

	// 将 map 对象转为 json 对象
	// json.Marshal() 函数的参数为 map 对象的指针，返回 ([]byte, error)，前者为 json 字符串（未编码），后者为是否错误
	if data, err := json.Marshal(&m); err == nil {
		s = string(data) // 将返回的字符串编码为 string 对象
	}
	assert.Equal(t, JSON_FROM_MAP, s)
}

// 序列化：从结构体对象生成 json 字符串
// 同样是使用 json.Marshal 函数，参数为结构体对象指针
// 结构体要求：所有需要被序列化为 json 的字段名称必须以 大写字母 开头（即公开属性）
// 可以通过 tag 设置结构体字段转换为 json 的字段名，如：`json:name`，另外，omitempty 表示如果该字段为空，则不出现在 json 结果中，例如 `json:name,omitempty`
func TestJsonMarshalFromStruct(t *testing.T) {
	u := NewUser(1, "Alvin", "alvin@fake.com", []string{"13999912345", "13000056789"}) // 产生结构体变量，指针类型

	data, err := json.Marshal(u) // 将结构体对象转为 json 字符串，参数为 结构体变量地址
	assert.NoError(t, err)       // 确认转换成功
	assert.Equal(t, JSON, string(data))

	data, err = json.MarshalIndent(u, "", "  ") // 带 json 格式的转换，可以设置 json 结构的前缀字符串和缩进字符
	assert.NoError(t, err)
	assert.Equal(t, JSON_FORMATTED, string(data))

	// 结构体 Email 字段的 tag 标记为 omitempty，所以如果为空，则不出现在结果 json 中
	u.Email = ""
	data, err = json.Marshal(u)
	assert.NoError(t, err)
	assert.Equal(t, `{"id":1,"name":"Alvin","phone":["13999912345","13000056789"]}`, string(data))
}

// 反序列化：将 json 字符串转换为 map 对象
// 转换的 map 为 map[string]interface{} 类型的，所以要处理
func TestJsonUnMarshalToMap(t *testing.T) {
	// 定义 json 字符串并转为 byte 序列
	s := []byte(JSON_FORMATTED)

	// 产生一个空的 map 对象
	m := make(map[string]interface{}, 10)

	// 进行反序列化操作，将 json 字符串转为 map 对象
	// json.Unmarshal 函数接受一个字符串和 map 对象指针，将 json 对象反序列化后填充到 map 对象中
	err := json.Unmarshal(s, &m) // 返回 error 对象，表示转换过程中是否出现错误
	assert.NoError(t, err)

	phones, err := ConvertInterfaceToStringSlice(m["phone"])
	assert.NoError(t, err)

	assert.Equal(t, 1, int(m["id"].(float64)))
	assert.Equal(t, "Alvin", m["name"].(string))
	assert.Equal(t, "alvin@fake.com", m["email"].(string))
	assert.Equal(t, []string{"13999912345", "13000056789"}, phones)
}

// 反序列化：将 json 字符串
func TestJsonUnMarshalToStruct(t *testing.T) {
	// 定义 json 字符串并转为 byte 序列
	s := []byte(JSON_FORMATTED)

	// 定义 User 类型对象，无需进行初始化
	u := User{}

	// 将 json 字符串反序列化为 User 对象
	// 仍用 json.Unmarshal 函数，第一个参数为 json 字符串，第二个参数为 User 对象指针
	// 将 json 字符串按 User 对象字段名（或 tag 标记）填充到 User 类型对象中
	err := json.Unmarshal(s, &u)

	assert.NoError(t, err)
	assert.Equal(t, User{Id: 1, Name: "Alvin", Email: "alvin@fake.com", Phone: []string{"13999912345", "13000056789"}}, u)
}
