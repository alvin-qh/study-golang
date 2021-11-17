package optargs

import "basic/builtin/types"

// CreateUser 函数的参数结构体
type UserOptions struct {
	Id     int
	Name   string
	Gender rune
}

// CreateUser 函数参数项接口
type UserOption interface {
	// 用于设置 UserOptions 某项的函数
	apply(*UserOptions)
}

// 用于设置 UserOptions 结构体某个属性值
type UserOptionsSetter struct {
	// 设置 UserOptions 属性值的函数
	set func(*UserOptions)
}

// 实现 UserOption 接口
func (setter *UserOptionsSetter) apply(opts *UserOptions) {
	setter.set(opts)
}

// 产生一个 UserOptionsSetter 结构体对象
func newFuncOption(setter func(*UserOptions)) *UserOptionsSetter {
	return &UserOptionsSetter{set: setter}
}

// 用于设置 User Id 属性
func WithUserId(id int) UserOption {
	return newFuncOption(func(opt *UserOptions) { opt.Id = id })
}

// 用于设置 User Name 属性
func WithUserName(name string) UserOption {
	return newFuncOption(func(opt *UserOptions) { opt.Name = name })
}

// 用于设置 User Gender 属性
func WithUserGender(gender rune) UserOption {
	return newFuncOption(func(opt *UserOptions) { opt.Gender = gender })
}

// 产生默认的 UserOption 对象
func defaultUserOption() *UserOptions {
	return &UserOptions{Id: 1, Name: "Alvin", Gender: 'M'}
}

// 创建一个 User 对象
func CreateUser(opts ...UserOption) *types.User {
	// 先将 arg 设置为默认的 UserOptions 对象
	arg := defaultUserOption()

	// 根据传入的 opts 参数值，设置 arg 中的各属性，未设置的属性保持其默认值
	for _, opt := range opts {
		opt.apply(arg)
	}
	// 使用设置后的 arg 参数
	return &types.User{Id: arg.Id, Name: arg.Name, Gender: arg.Gender}
}
