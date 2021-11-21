package gob

import (
	"bytes"
	"encoding/gob"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 测试通过 gob 序列化和反序列化对象
func TestGobEncodeAndDecode(t *testing.T) {
	// 用于序列化的 Buffer 对象
	buf := bytes.NewBuffer(make([]byte, 0))

	// 序列化对象

	s := "Hello World"

	enc := gob.NewEncoder(buf) // 生成 Encoder 对象

	err := enc.Encode(len(s)) // 编码一个字符串长度
	assert.NoError(t, err)

	err = enc.Encode(s) // 编码一个字符串
	assert.NoError(t, err)

	prod, err := NewProduct("Apple", 100, 2.56)
	assert.NoError(t, err)

	err = prod.Serialize(buf) // 编码一个结构体
	assert.NoError(t, err)

	// 用于反序列化的 Reader 对象
	data := buf.Bytes()
	reader := bytes.NewReader(data)

	// 反序列化对象

	dec := gob.NewDecoder(reader) // 生成 Decoder 对象

	var n int
	err = dec.Decode(&n) // 反序列化字符串长度
	assert.NoError(t, err)
	assert.Equal(t, 11, n)

	s = ""
	err = dec.Decode(&s) // 反序列化字符串
	assert.NoError(t, err)
	assert.Equal(t, "Hello World", s)

	res, err := DeserializeProduct(reader) // 反序列化结构体
	assert.NoError(t, err)
	assert.Equal(t, *prod, *res)
}
