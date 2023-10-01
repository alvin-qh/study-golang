package routes

import (
	"study-gin/core/server"
	"time"

	"github.com/go-playground/validator/v10"
)

// 定义表示性别的类型
type Gender string

// 定义表示拥护的结构体
type User struct {
	Name     string    `json:"name"`
	Gender   Gender    `json:"gender"`
	Birthday time.Time `json:"birthday"`
}

// 定义表示用户表单的结构体
//
// Tag 中的描述如下:
//   - `form:"name"` 表示结构体字段和表单字段的对应关系
//   - `binding:"required,min=3,max=20"` 表示绑定表单时对表单的验证方式, 也可以写为 `validate:"required,..."`
type UserForm struct {
	Name       string `form:"name" json:"name" binding:"required,min=3,max=20"`
	Gender     Gender `form:"gender" json:"gender" binding:"required,oneof=F M"`
	BirthYear  int    `form:"birth_year" json:"birthYear" binding:"required,min=1960,max=9999"`
	BirthMonth int    `form:"birth_month" json:"birthMonth" binding:"required,min=1,max=12"`
	BirthDay   int    `form:"birth_day" json:"birthDay" binding:"required,min=1,max=31"`
}

// 将当前 `UserForm` 结构体变量转为 `User` 结构体变量
//
// 返回:
//   - `User` 结构体变量指针
func (u *UserForm) toUser() *User {
	return &User{
		Name:   u.Name,
		Gender: u.Gender,
		Birthday: time.Date(
			u.BirthYear,
			time.Month(u.BirthMonth),
			u.BirthDay,
			0, 0, 0, 0, time.UTC,
		),
	}
}

const (
	InputErrorCode = "input_error"
)

// 定义打招呼的请求结构体
type Agreeing struct {
	Name   string `json:"name" binding:"required"`
	Gender Gender `json:"gender" binding:"required,oneof=F M"`
	Age    int    `json:"age" binding:"min=10,max=99"`
}

// 打招呼的响应结构体
type Answer struct {
	Ok     bool   `json:"ok"`
	Answer string `json:"answer"`
}

// 表示字段错误的结构体
type ErrorField struct {
	Name  string `json:"name"`
	Error any    `json:"error"`
}

// 错误结构体
type ErrorResult struct {
	Code        string       `json:"code"`
	ErrorFields []ErrorField `json:"errorFields"`
}

func (e* ErrorResult) FromValidator(errs validator.ValidationErrors, target any) {
	fs := make([]ErrorField, 0)
	for k, v := range server.MappedValidatorErrors(errs, target, "json") {
		fs = append(fs, ErrorField{
			Name:  k,
			Error: v,
		})
	}
	e.ErrorFields = fs
}
