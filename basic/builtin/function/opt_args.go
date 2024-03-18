package function

import "study-golang/basic/builtin/types"

// 用于 `CreateUser` 函数的可选参数结构体
type option struct {
	Id     int
	Name   string
	Gender rune
}

// 用于设置 `CreateUser` 可选参数的回调类型
type UserOption = func(*option)

// 用于设置 `Id` 属性的回调函数
//
// 参数:
//   - `id`: 可选参数值
//
// 返回 `UserOption` 回调函数
func WithUserId(id int) UserOption {
	return func(o *option) {
		o.Id = id
	}
}

// 用于设置 User Name 属性
func WithUserName(name string) UserOption {
	return func(o *option) {
		o.Name = name
	}
}

// 用于设置 User Gender 属性
func WithUserGender(gender rune) UserOption {
	return func(o *option) {
		o.Gender = gender
	}
}

// 创建一个 User 对象
func CreateUser(opts ...UserOption) *types.User {
	def := option{
		Id:     1,
		Name:   "Alvin",
		Gender: 'M',
	}

	// 根据传入的 opts 参数值, 设置 arg 中的各属性, 未设置的属性保持其默认值
	for _, opt := range opts {
		opt(&def)
	}
	// 使用设置后的 arg 参数
	return &types.User{Id: def.Id, Name: def.Name, Gender: def.Gender}
}
