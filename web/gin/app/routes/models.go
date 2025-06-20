package routes

import (
	"time"
	"web/gin/core/server"

	"github.com/go-playground/validator/v10"
)

// 定义表示性别的类型
type Gender string

// 定义性别常量
const (
	GenderF Gender = "F"
	GenderM Gender = "M"
)

// 定义表示拥护的结构体
type User struct {
	Id       string    `json:"id,omitempty"`
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
// 参数:
//   - `id` (`string`): 用户 Id
//
// 返回:
//   - `User` 结构体变量指针
func (u *UserForm) toUser(id string) *User {
	return &User{
		Id:     id,
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

// 错误代码
const (
	OkCode         = "ok"
	InputErrorCode = "input_error"
)

// 响应结果负载结构体
type ResponseData struct {
	Code    string `json:"code"`
	Payload any    `json:"payload"`
}

// 表示字段错误的结构体
type ErrorField struct {
	Name  string `json:"name"`
	Error any    `json:"error"`
}

// 错误结构体
type ErrorResult struct {
	Error       string       `json:"error,omitempty"`
	ErrorFields []ErrorField `json:"errorFields,omitempty"`
}

// 创建结果负载
func NewResponseData(payload any) *ResponseData {
	return &ResponseData{
		Code:    OkCode,
		Payload: payload,
	}
}

// 创建输入错误对象
func NewInputError(err error, target any) *ResponseData {
	er := new(ErrorResult)

	if e, ok := err.(validator.ValidationErrors); ok {
		fs := make([]ErrorField, 0)
		for k, v := range server.MappedValidatorErrors(e, target, "json") {
			fs = append(fs, ErrorField{
				Name:  k,
				Error: v,
			})
		}
		er.ErrorFields = fs
	} else {
		er.Error = err.Error()
	}

	return &ResponseData{
		Code:    InputErrorCode,
		Payload: er,
	}
}
