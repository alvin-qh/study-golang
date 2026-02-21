package structure_test

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 定义结构体
//
// 结构体具备三个字段属性
type User struct {
	Id     int
	Name   string
	Gender rune
}

// 测试初始化结构体实例
func TestStructure_Initialize(t *testing.T) {
	// 定义结构体变量并进行初始化
	var u User = User{
		Id:     1,
		Name:   "Alvin",
		Gender: 'M',
	}

	assert.Equal(t, 1, u.Id)
	assert.Equal(t, "Alvin", u.Name)
	assert.Equal(t, 'M', u.Gender)
}

// 测试获取或设置结构体实例的字段值
func TestStructure_Properties(t *testing.T) {
	// 定义结构体并按默认值初始化字段值
	u := User{}

	// 查看变量类型
	assert.Equal(t, reflect.Struct, reflect.TypeFor[User]().Kind())
	assert.Equal(t, "User", reflect.TypeFor[User]().Name())

	// 为结构体字段赋值
	u.Id = 2
	u.Name = "Emma"
	u.Gender = 'F'

	// 确认结构体字段值已经被修改
	assert.Equal(t, 2, u.Id)
	assert.Equal(t, "Emma", u.Name)
	assert.Equal(t, 'F', u.Gender)
}

func TestStructure_Pointer(t *testing.T) {
	// 通过 new 操作符产生 User 类型的指针变量, 并在指针指向的内存空间中分配一个 User 结构体对象,
	// 但不对结构体字段进行初始化
	pu := new(User)

	// 记录原始指针
	puo := pu

	// 修改指针指向的结构体对象, 即将 pu 指针指向的结构体进行更换
	*pu = User{
		Id:     1,
		Name:   "Alvin",
		Gender: 'M',
	}

	// 查看指针指向的结构体字段值
	assert.Equal(t, puo, pu)
	assert.Equal(t, 1, pu.Id)
	assert.Equal(t, "Alvin", pu.Name)
	assert.Equal(t, 'M', pu.Gender)

	// 令 pu 变量存储一个新的结构体对象的地址
	pu = &User{
		Id:     2,
		Name:   "Emma",
		Gender: 'F',
	}

	// 确认此事 pu 和 puo 指向的结构体体对象已经不同
	assert.NotEqual(t, puo, pu)
	assert.Equal(t, 2, pu.Id)
	assert.Equal(t, "Emma", pu.Name)
	assert.Equal(t, 'F', pu.Gender)
}
