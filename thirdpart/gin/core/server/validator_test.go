package server

import (
	"testing"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

// 在所有测试前执行
func TestMain(m *testing.M) {
	// 初始化验证器对象
	SetupValidator()

	// 执行测试
	m.Run()
}

// 测试用结构体
type UserForm struct {
	Name       string `form:"name" binding:"required,min=3,max=20"`
	Gender     string `form:"gender" binding:"required,oneof=F M"`
	BirthYear  int    `form:"birth_year" binding:"required,min=1960,max=9999"`
	BirthMonth int    `form:"birth_month" binding:"required,min=1,max=12"`
	BirthDay   int    `form:"birth_day" binding:"required,min=1,max=31"`
}

// 测试验证器验证表单对象解析
func TestValidateErrorToString(t *testing.T) {
	// 获取 gin 默认的验证器对象
	validate := binding.Validator.Engine().(*validator.Validate)

	// 实例化表单对象
	user := UserForm{
		Name:       "Al",
		Gender:     "X",
		BirthYear:  0,
		BirthMonth: 0,
		BirthDay:   0,
	}

	// 对表单对象进行检验
	errors := validate.Struct(&user).(validator.ValidationErrors)
	assert.Len(t, errors, 5)

	// 确认验证器验证结果符合预期

	assert.Equal(t, "min", errors[0].ActualTag())
	assert.Equal(t, "3", errors[0].Param())

	assert.Equal(t, "oneof", errors[1].ActualTag())
	assert.Equal(t, "F M", errors[1].Param())

	assert.Equal(t, "required", errors[2].ActualTag())
	assert.Equal(t, "", errors[2].Param())
}

// 测试验证器验证表单对象解析
func TestMappedValidatorErrors(t *testing.T) {
	// 获取 gin 默认的验证器对象
	validate := binding.Validator.Engine().(*validator.Validate)

	// 实例化表单对象
	user := UserForm{
		Name:       "Al",
		Gender:     "X",
		BirthYear:  0,
		BirthMonth: 0,
		BirthDay:   0,
	}

	// 对表单对象进行检验, 并对验证结果进行转换
	errors := MappedValidatorErrors(
		validate.Struct(&user).(validator.ValidationErrors),
		&user,
		"form",
	)
	assert.Len(t, errors, 5)

	// 确认验证器验证结果符合预期

	assert.Equal(t, "Name长度必须至少为3个字符", errors["name"])
	assert.Equal(t, "Gender必须是[F M]中的一个", errors["gender"])
	assert.Equal(t, "BirthYear为必填字段", errors["birth_year"])
	assert.Equal(t, "BirthMonth为必填字段", errors["birth_month"])
	assert.Equal(t, "BirthDay为必填字段", errors["birth_day"])
}
