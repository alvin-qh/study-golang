package xml

import (
	"basic/io/user"
	"bytes"
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 将 map 对象转为 xml，即将 map 的 key 作为 xml 的节点名，value 作为 xml 的节点值

// 将 结构体 编码为 xml，或者从 xml 中还原 结构体 的字段值
// 要求：结构体必须为公开，且结构体字段为公开，私有的结构体字段会被忽略
// 如果需要对 xml 节点做进一步描述，需要在结构体中使用 tag 进行标注，参考 user.User 结构体

// 序列化：从 map 对象生成 json 字符串
func TestXmlMarshal(t *testing.T) {
	pu := user.New(1, "Alvin", "alvin@fake.com", []string{"13999912345", "13000056789"})

	// 将 map 对象转为 json 对象
	// json.Marshal() 函数的参数为 map 对象的指针，返回 ([]byte, error)，前者为 json 字符串（未编码），后者为是否错误
	data, err := xml.Marshal(pu)
	assert.NoError(t, err)

	s := string(data) // 将返回的字符串编码为 string 对象
	assert.Equal(t, `<user id="1">`+
		`<name>Alvin</name>`+
		`<email>alvin@fake.com</email>`+
		`<phones><tel>13999912345</tel><tel>13000056789</tel></phones>`+
		`</user>`, s)
}

func TestXmlUnmarshal(t *testing.T) {
	s := `<user id="1">` +
		`<name>Alvin</name>` +
		`<email>alvin@fake.com</email>` +
		`<phones><tel>13999912345</tel><tel>13000056789</tel></phones>` +
		`</user>`
	u := user.User{}

	xml.Unmarshal([]byte(s), &u)

	assert.Equal(t, int64(1), u.Id)
	assert.Equal(t, "Alvin", u.Name)
	assert.Equal(t, "alvin@fake.com", u.Email)
	assert.Equal(t, []string{"13999912345", "13000056789"}, u.Phone)
}

func TestXmlEncoder(t *testing.T) {
	u := user.New(1, "Alvin", "alvin@fake.com", []string{"13999912345", "13000056789"})

	buf := bytes.NewBuffer(make([]byte, 0, 1024))
	enc := xml.NewEncoder(buf)
	enc.Indent("", "  ")

	err := enc.Encode(u)
	assert.NoError(t, err)
	enc.Flush()

	s := buf.String()
	assert.Equal(t, `<user id="1">
  <name>Alvin</name>
  <email>alvin@fake.com</email>
  <phones>
    <tel>13999912345</tel>
    <tel>13000056789</tel>
  </phones>
</user>`, s)
}

func TestXmlDecoder(t *testing.T) {
	buf := bytes.NewBuffer([]byte(`<user id="1">
  <name>Alvin</name>
  <email>alvin@fake.com</email>
  <phones>
    <tel>13999912345</tel>
    <tel>13000056789</tel>
  </phones>
</user>`))

	dec := xml.NewDecoder(buf)

	u := user.User{}

	err := dec.Decode(&u)
	assert.NoError(t, err)

	assert.Equal(t, int64(1), u.Id)
	assert.Equal(t, "Alvin", u.Name)
	assert.Equal(t, "alvin@fake.com", u.Email)
	assert.Equal(t, []string{"13999912345", "13000056789"}, u.Phone)
}
