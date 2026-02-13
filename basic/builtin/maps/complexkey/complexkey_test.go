package complexkey_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// 定义结构体作为 Map 的 Key
type Key struct {
	Id   int
	Name string
}

// 定义结构体作为 Map 的 Value
type Value struct {
	Gender   rune
	Birthday string
	Address  string
}

// 测试结构体作为 Map 的 Key
//
// 由于 Go 语言支持结构体的比较和散列, 所以结构体可以直接作为 Map 的 Key
func TestMap_ComplexKey(t *testing.T) {
	// 使用结构体作为 Map 的 Key
	m := map[Key]*Value{}

	m[Key{1, "Alvin"}] = &Value{
		Gender:   'M',
		Birthday: "1981-03",
		Address:  "ShanXi, Xi'an",
	}

	assert.Equal(t, Value{
		Gender:   'M',
		Birthday: "1981-03",
		Address:  "ShanXi, Xi'an",
	}, *(m[Key{1, "Alvin"}]))
}
