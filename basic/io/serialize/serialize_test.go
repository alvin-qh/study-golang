package serialize

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"study-golang/basic/generic"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 定义 user 结构体, 用于测试序列化, 反序列化
type user struct {
	XMLName xml.Name `json:"-" xml:"user"`     // 用于定义 XML 根节点名称, json 中忽略 (json:"-" 表示该字段不出现在 json 中)
	Id      int64    `json:"id" xml:"id,attr"` // "attr" 表示在 XML 中, "id" 字段的值在 根节点属性上表示, 而不是使用 XML 节点
	Name    string   `json:"name" xml:"name"`
	Email   string   `json:"email,omitempty" xml:"email,omitempty"` // "omitempty" 表示如果为空, 则不出现在 json 或 XML 中
	Phone   []string `json:"phone,omitempty" xml:"phones>tel"`      // "phones>tel" 表示在 XML 中, 切片类型字段位于 "phones" 节点下, 每一项是一个 "tel" 节点
}

// 为结构体实例添加手机号码
func (u *user) AddPhone(phone string) {
	u.Phone = append(u.Phone, phone)
}

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

// 将 Map 类型实例转为 JSON
//
// 将 Map 类型实例转为 JSON, 即将 Map 的 Key 作为 JSON 的字段名, Map 的 Value 作为 JSON 的字段值
func TestJsonMarshalFromMap(t *testing.T) {
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
	dat, err := json.Marshal(&m)
	assert.NoError(t, err)

	// 将返回的字符串编码为 string 对象
	s := string(dat)
	assert.Equal(t, JSON_FROM_MAP, s)
}

// 将 JSON 转换为 Map 类型实例
//
// 转换得到的 Map 类型为 `map[string]interface{}` 类型
func TestJsonUnMarshalToMap(t *testing.T) {
	// 定义 json 字符串并转为 byte 序列
	s := []byte(JSON_FORMATTED)

	// 产生一个空的 map 对象
	m := make(map[string]interface{}, 10)

	// 进行反序列化操作, 将 json 字符串转为 map 对象
	// json.Unmarshal 函数接受一个字符串和 map 对象指针, 将 json 对象反序列化后填充到 map 对象中
	err := json.Unmarshal(s, &m) // 返回 error 对象, 表示转换过程中是否出现错误
	assert.NoError(t, err)

	phones, err := generic.InterfaceToSlice[string](m["phone"])
	assert.NoError(t, err)

	assert.Equal(t, 1, int(m["id"].(float64)))
	assert.Equal(t, "Alvin", m["name"].(string))
	assert.Equal(t, "alvin@fake.com", m["email"].(string))
	assert.Equal(t, []string{"13999912345", "13000056789"}, phones)
}

// 从结构体实例生成 JSON 字符串
//
// 使用 `json.Marshal` 函数, 传入结构体实例指针, 要求:
//   - 结构体的所有被序列化字段名称必须以大写字母开头 (即公开属性)
//   - 可以通过字段 tag 设置结构体字段转换后的字段名, 如: `json:name`; `omitempty` 表示如果字段为空,则不出现在 JSON 结果中, 例如 `json:name,omitempty`
func TestJsonMarshalFromStruct(t *testing.T) {
	// 产生结构体变量, 指针类型
	u := &user{
		Id:    1,
		Name:  "Alvin",
		Email: "alvin@fake.com",
		Phone: []string{
			"13999912345",
			"13000056789",
		},
	}

	// 将结构体对象转为 json 字符串, 参数为 结构体变量地址
	data, err := json.Marshal(u)
	assert.NoError(t, err) // 确认转换成功
	assert.Equal(t, JSON, string(data))

	// 带 json 格式的转换, 可以设置 json 结构的前缀字符串和缩进字符
	data, err = json.MarshalIndent(u, "", "  ")
	assert.NoError(t, err)
	assert.Equal(t, JSON_FORMATTED, string(data))

	// 结构体 Email 字段的 tag 标记为 omitempty, 所以如果为空, 则不出现在结果 json 中
	u.Email = ""
	data, err = json.Marshal(u)
	assert.NoError(t, err)
	assert.Equal(t, `{"id":1,"name":"Alvin","phone":["13999912345","13000056789"]}`, string(data))
}

// 将 JSON 转为结构体实例
//
// 通过 `json.Unmarshal 函数, 将一个 JSON 字符串 (参数1) 反序列化到结构体实例 (参数2) 中
//
// 函数会将 JSON 字符串按结构体字段名 (或 tag 标记) 填充到结构体实例中
func TestJsonUnMarshalToStruct(t *testing.T) {
	// 定义 json 字符串并转为 byte 序列
	s := []byte(JSON_FORMATTED)

	// 定义 User 类型对象, 无需进行初始化
	u := user{}

	// 进行反序列化操作, 将 JSON 转为结构体实例
	err := json.Unmarshal(s, &u)
	assert.NoError(t, err)

	// 确认反序列化正确
	expected := user{
		Id:    1,
		Name:  "Alvin",
		Email: "alvin@fake.com",
		Phone: []string{
			"13999912345",
			"13000056789",
		},
	}
	assert.Equal(t, expected, u)
}

const (
	XML = `<user id="1">` +
		`<name>Alvin</name>` +
		`<email>alvin@fake.com</email>` +
		`<phones><tel>13999912345</tel><tel>13000056789</tel></phones>` +
		`</user>`

	XML_FORMATTED = `<user id="1">
  <name>Alvin</name>
  <email>alvin@fake.com</email>
  <phones>
    <tel>13999912345</tel>
    <tel>13000056789</tel>
  </phones>
</user>`
)

// 将结构体实例转为 XML
//
// 通过 `xml.Marshal` 函数将结构体实例转为 XML, 要求:
//   - 结构体的所有被序列化字段名称必须以大写字母开头 (即公开属性)
//   - 可以通过字段 tag 设置结构体字段转换后的字段名, 如: `xml:name`; `omitempty` 表示如果字段为空,则不出现在 XML 结果中, 例如 `xml:name,omitempty`
func TestXmlMarshal(t *testing.T) {
	pu := &user{
		Id:    1,
		Name:  "Alvin",
		Email: "alvin@fake.com",
		Phone: []string{
			"13999912345",
			"13000056789",
		},
	}

	// 结构体实例转为 XML
	data, err := xml.Marshal(pu)
	assert.NoError(t, err)

	// 将返回的字节串编码为字符串
	s := string(data)
	assert.Equal(t, XML, s)
}

// 将 XML 反序列化为结构体实例
func TestXmlUnmarshal(t *testing.T) {
	u := user{}

	// 将 bytes 反序列化, 填充到结构体对象中
	xml.Unmarshal([]byte(XML), &u)

	assert.Equal(t, int64(1), u.Id)
	assert.Equal(t, "Alvin", u.Name)
	assert.Equal(t, "alvin@fake.com", u.Email)
	assert.Equal(t, []string{"13999912345", "13000056789"}, u.Phone)
}

// 将结构体实例序列化为 XML, 并写入 `io.Write` 接口实例中
func TestXmlEncoder(t *testing.T) {
	pu := &user{
		Id:    1,
		Name:  "Alvin",
		Email: "alvin@fake.com",
		Phone: []string{
			"13999912345",
			"13000056789",
		},
	}

	buf := bytes.NewBuffer(make([]byte, 0, 1024))

	// 通过一个 io.Writer 对象, 产生一个 XML Encoder 对象
	enc := xml.NewEncoder(buf)

	// 设置 XML 格式
	enc.Indent("", "  ")

	// 将结构体序列化到 io.Writer 中
	err := enc.Encode(pu)
	assert.NoError(t, err)

	enc.Flush()

	s := buf.String()
	assert.Equal(t, XML_FORMATTED, s)
}

// 从 `io.Reader` 接口实例中读取 XML, 并反序列化为结构体实例
func TestXmlDecoder(t *testing.T) {
	buf := bytes.NewBuffer([]byte(XML_FORMATTED))

	// 通过一个 io.Reader 对象, 产生一个 XML Decoder 对象
	dec := xml.NewDecoder(buf)

	u := user{}

	// 从 io.Reader 中读取 XML, 并反序列化到结构体中
	err := dec.Decode(&u)
	assert.NoError(t, err)

	assert.Equal(t, int64(1), u.Id)
	assert.Equal(t, "Alvin", u.Name)
	assert.Equal(t, "alvin@fake.com", u.Email)
	assert.Equal(t, []string{"13999912345", "13000056789"}, u.Phone)
}
