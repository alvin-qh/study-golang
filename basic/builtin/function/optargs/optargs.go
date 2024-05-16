package optargs

// 在 Go 中使用参数默认值
//
// 在 Go 语言中, 并没有提供参数默认值语法, 但可以通过结构体进行模拟
//
// 定义结构体, 用于测试
type User struct {
	Id     int
	Name   string
	Gender rune
}

// 用于 `CreateUser` 函数的可选参数结构体
type option struct {
	id     int
	name   string
	gender rune
}

// 用于设置 `CreateUser` 可选参数的回调类型
type UserOption = func(*option)

// 创建一个 User 对象
func CreateUser(opts ...UserOption) *User {
	// 设置可选参数默认值
	def := option{
		id:     1,
		name:   "Alvin",
		gender: 'M',
	}

	// 根据传入的 opts 参数值, 设置 arg 中的各属性, 未设置的属性保持其默认值
	for _, opt := range opts {
		opt(&def)
	}

	// 使用设置后的 arg 参数
	return &User{
		Id:     def.id,
		Name:   def.name,
		Gender: def.gender,
	}
}

// 用于设置名为 `id` 的可选参数
func WithUserId(id int) UserOption {
	return func(o *option) {
		o.id = id
	}
}

// 用于设置名为 `Name` 的可选参数
func WithUserName(name string) UserOption {
	return func(o *option) {
		o.name = name
	}
}

// 用于设置名为 `Gender` 的可选参数
func WithUserGender(gender rune) UserOption {
	return func(o *option) {
		o.gender = gender
	}
}
