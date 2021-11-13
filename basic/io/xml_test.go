package io

import (
	"basic/io/user"
	"bytes"
	"encoding/xml"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 将 结构体 编码为 json，或者从 json 中还原 结构体 的字段值
// 要求：结构体必须为公开，且结构体字段为公开，私有的结构体字段会被忽略
// 如果需要对 json 字段做进一步描述，需要在结构体中使用 tag 进行标注，参考 user.User 结构体

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
