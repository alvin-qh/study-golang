package serialize

import (
	"bytes"
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 定义 user 结构体, 用于测试序列化, 反序列化
type XMLUser struct {
	XMLName xml.Name `xml:"user"`    // 用于定义 XML 根节点名称
	Id      int64    `xml:"id,attr"` // "attr" 表示在 XML 中, "id" 字段的值在 根节点属性上表示, 而不是使用 XML 节点
	Name    string   `xml:"name"`
	Email   string   `xml:"email,omitempty"` // "omitempty" 表示如果为空, 则不出现在 XML 中
	Phone   []string `xml:"phones>tel"`      // "phones>tel" 表示在 XML 中, 切片类型字段位于 "phones" 节点下, 每一项是一个 "tel" 节点
}

// 为结构体实例添加手机号码
func (u *XMLUser) AddPhone(phone string) {
	u.Phone = append(u.Phone, phone)
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
func TestXML_Marshal(t *testing.T) {
	u := XMLUser{
		Id:    1,
		Name:  "Alvin",
		Email: "alvin@fake.com",
		Phone: []string{
			"13999912345",
			"13000056789",
		},
	}

	// 结构体实例转为 XML
	data, err := xml.Marshal(&u)
	assert.Nil(t, err)

	// 将返回的字节串编码为字符串
	s := string(data)
	assert.Equal(t, XML, s)
}

// 将 XML 反序列化为结构体实例
func TestXML_Unmarshal(t *testing.T) {
	u := XMLUser{}

	// 将 bytes 反序列化, 填充到结构体对象中
	xml.Unmarshal([]byte(XML), &u)

	assert.Equal(t, int64(1), u.Id)
	assert.Equal(t, "Alvin", u.Name)
	assert.Equal(t, "alvin@fake.com", u.Email)
	assert.Equal(t, []string{"13999912345", "13000056789"}, u.Phone)
}

// 将结构体实例序列化为 XML, 并写入 `io.Write` 接口实例中
func TestXML_Encode(t *testing.T) {
	pu := XMLUser{
		Id:    1,
		Name:  "Alvin",
		Email: "alvin@fake.com",
		Phone: []string{
			"13999912345",
			"13000056789",
		},
	}

	// 产生一个字节缓冲区
	buf := bytes.NewBuffer(make([]byte, 0))

	// 通过一个 io.Writer 对象, 产生一个 XML Encoder 对象
	enc := xml.NewEncoder(buf)

	// 设置 XML 格式
	enc.Indent("", "  ")

	// 将结构体序列化到 io.Writer 中
	err := enc.Encode(pu)
	assert.Nil(t, err)

	enc.Flush()
	assert.Equal(t, XML_FORMATTED, buf.String())
}

// 从 `io.Reader` 接口实例中读取 XML, 并反序列化为结构体实例
func TestXML_Decode(t *testing.T) {
	buf := bytes.NewBuffer([]byte(XML_FORMATTED))

	// 通过一个 io.Reader 对象, 产生一个 XML Decoder 对象
	dec := xml.NewDecoder(buf)

	u := XMLUser{}

	// 从 io.Reader 中读取 XML, 并反序列化到结构体中
	err := dec.Decode(&u)
	assert.Nil(t, err)

	assert.Equal(t, int64(1), u.Id)
	assert.Equal(t, "Alvin", u.Name)
	assert.Equal(t, "alvin@fake.com", u.Email)
	assert.Equal(t, []string{"13999912345", "13000056789"}, u.Phone)
}
