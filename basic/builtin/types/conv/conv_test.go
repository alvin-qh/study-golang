package conv_test

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 定义结构体
type User struct {
	Id     int
	Name   string
	Gender rune
}

// 创建类型不同的对象
func makeDifferentTypeObject() any {
	var obj any

	// 根据 0~2 的随机数结果, 创建不同类型的对象
	switch rand.Intn(3) {
	case 0:
		obj = 100 // 创建整数类型对象
	case 1:
		obj = "Hello" // 创建字符串类型对象
	case 2:
		obj = User{ // 创建结构体类型对象
			Id:     1,
			Name:   "Alvin",
			Gender: 'M',
		}
	}
	return obj
}

// 测试利用 `switch` 语句进行类型转换
func TestConv_ConvertWithSwitch(t *testing.T) {
	// 创建不同类型的对象
	obj := makeDifferentTypeObject()

	// 通过 switch 语句进行类型转换
	// 每个分支用于判断 `v` 变量的一种类型, 如果类型匹配到具体分支, 则 `vv` 变量是该类型的值
	switch val := obj.(type) {
	case int:
		assert.Equal(t, 100, val)
	case string:
		assert.Equal(t, "Hello", val)
	case User:
		assert.Equal(t, User{
			Id:     1,
			Name:   "Alvin",
			Gender: 'M',
		}, val)
	default:
		assert.Fail(t, "unknown type")
	}
}
