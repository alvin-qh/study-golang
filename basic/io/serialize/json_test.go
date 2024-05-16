package serialize

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 定义 user 结构体, 用于测试序列化, 反序列化
type JSONUser struct {
	Ignore bool     `json:"-" ` // `-` 表示该字段忽略序列化
	Id     int64    `json:"id"`
	Name   string   `json:"name"`
	Email  string   `json:"email,omitempty"` // "omitempty" 表示如果为空
	Phone  []string `json:"phone,omitempty"`
}

// 为结构体实例添加手机号码
func (u *JSONUser) AddPhone(phone string) {
	u.Phone = append(u.Phone, phone)
}

const (
	JSON_FROM_MAP  = `{"email":"alvin@fake.com","id":1,"name":"Alvin","phone":["13999912345","13000056789"]}`
	JSON           = `{"id":1,"name":"Alvin","email":"alvin@fake.com","phone":["13999912345","13000056789"]}`
	JSON_FORMATTED = `{
  "id": 1,
  "name": "Alvin",
  "email": "alvin@fake.com",
  "phone": [
    "13999912345",
    "13000056789"
  ]
}`
	JSON_OMITEMPTY = `{"id":1,"name":"Alvin","phone":["13999912345","13000056789"]}`
)

// 测试 JSON 序列号
func TestJSON_Marshal(t *testing.T) {
	// 将 Map 实例序列化为 JSON 字符串
	t.Run("marshal map to json", func(t *testing.T) {
		// 产生一个 map 对象, key 为 string, value 任意
		m := map[string]interface{}{
			"id":    1,
			"name":  "Alvin",
			"email": "alvin@fake.com",
			"phone": []string{
				"13999912345",
				"13000056789",
			},
		}

		// json.Marshal 函数的参数为 Map 实例指针, 返回包含 JSON 字符串的字节串
		data, err := json.Marshal(&m)
		assert.Nil(t, err)

		// 将返回的字符串编码为 string 对象
		s := string(data)
		assert.Equal(t, JSON_FROM_MAP, s)
	})

	// 将结构体实例序列化为 JSON 字符串
	//
	// 使用 `json.Marshal` 函数, 传入结构体实例指针, 要求:
	//   - 结构体的所有被序列化字段名称必须以大写字母开头 (即公开属性)
	//   - 可以通过字段 tag 设置结构体字段转换后的字段名, 如: `json:name`; `omitempty` 表示如果字段为空,
	//     则不出现在 JSON 结果中, 例如 `json:name,omitempty`
	t.Run("marshal struct to json", func(t *testing.T) {
		// 产生结构体变量, 指针类型
		u := JSONUser{
			Id:    1,
			Name:  "Alvin",
			Email: "alvin@fake.com",
			Phone: []string{
				"13999912345",
				"13000056789",
			},
		}

		// 将结构体对象转为 json 字符串, 参数为 结构体变量地址
		data, err := json.Marshal(&u)
		assert.Nil(t, err)
		assert.Equal(t, JSON, string(data))
	})

	// 设置 JSON 序列化结果的格式
	t.Run("marshal with indent", func(t *testing.T) {
		// 产生结构体变量, 指针类型
		u := JSONUser{
			Id:    1,
			Name:  "Alvin",
			Email: "alvin@fake.com",
			Phone: []string{
				"13999912345",
				"13000056789",
			},
		}

		// 带 json 格式的转换, 可以设置 json 结构的前缀字符串和缩进字符
		data, err := json.MarshalIndent(&u, "", "  ")
		assert.Nil(t, err)
		assert.Equal(t, JSON_FORMATTED, string(data))
	})

	// 如果在结构体字段上标记 `omitempty`, 则当该字段为 `nil` 时不会包含在 JSON 序列化结果中
	t.Run("struct omitempty field", func(t *testing.T) {
		// 产生结构体变量, 指针类型
		// 结构体 Email 字段的 tag 标记为 omitempty, 所以如果为空, 则不出现在结果 json 中
		u := JSONUser{
			Id:   1,
			Name: "Alvin",
			Phone: []string{
				"13999912345",
				"13000056789",
			},
		}

		data, err := json.Marshal(&u)
		assert.Nil(t, err)
		assert.Equal(t, JSON_OMITEMPTY, string(data))
	})
}

// 测试 JSON 反序列化
func TestJSON_Unmarshal(t *testing.T) {
	// 将 JSON 反序列化到 Map 实例中
	t.Run("unmarshal to map", func(t *testing.T) {
		// 产生一个空的 map 对象
		m := make(map[string]interface{})

		// 进行反序列化操作, 将 json 字符串转为 map 对象
		// json.Unmarshal 函数接受一个字符串和 map 对象指针, 将 json 对象反序列化后填充到 map 对象中
		// 返回 error 对象, 表示转换过程中是否出现错误
		err := json.Unmarshal([]byte(JSON_FORMATTED), &m)
		assert.Nil(t, err)

		// JSON 中的数值型字段会反序列化为 float64 类型
		assert.Equal(t, float64(1), m["id"].(float64))
		assert.Equal(t, "Alvin", m["name"].(string))
		assert.Equal(t, "alvin@fake.com", m["email"].(string))
		assert.ElementsMatch(t, []string{"13999912345", "13000056789"}, m["phone"])
	})

	// 将 JSON 反序列化到结构体实例中
	t.Run("unmarshal to struct", func(t *testing.T) {
		// 定义 User 类型对象, 无需进行初始化
		u := JSONUser{}

		// 进行反序列化操作, 将 JSON 转为结构体实例
		err := json.Unmarshal([]byte(JSON_FORMATTED), &u)
		assert.Nil(t, err)

		// 确认反序列化正确
		assert.Equal(t, JSONUser{
			Id:    1,
			Name:  "Alvin",
			Email: "alvin@fake.com",
			Phone: []string{
				"13999912345",
				"13000056789",
			},
		}, u)
	})
}
