package server

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	transZh "github.com/go-playground/validator/v10/translations/zh"
)

var (
	translate ut.Translator // 保存验证器翻译器的全局变量
)

// 初始化验证器
func SetupValidator() {
	var found bool
	// 创建翻译器对象, 以中文为主语言, 英文为备选语言
	// 翻译器对象保存为全局变量
	translate, found = ut.New(en.New(), zh.New()).GetTranslator("zh")
	if !found {
		panic(fmt.Errorf("cannot find validate transaction by \"zh\""))
	}

	// 为翻译器注册默认的验证器对象
	validate := binding.Validator.Engine().(*validator.Validate)
	if err := transZh.RegisterDefaultTranslations(validate, translate); err != nil {
		panic(err)
	}
}

// 将表单校验错误对象转为 `map` 集合对象
//
// 参数:
//   - `errs` (`v.ValidationErrors`): 校验失败的错误对象集合
//
// 返回:
//   - Key 为表单项名称, Value 为错误描述的 `map` 集合对象
func MappedValidatorErrors(errs validator.ValidationErrors, target any, tag string) map[string]any {
	errMap := make(map[string]any)

	// 获取表单结构体对象类型
	t := reflect.ValueOf(target)
	realType := t.Type().Elem()

	// 遍历错误对象
	for _, e := range errs {
		// 获取发生错误的原始表单结构体字段
		field, _ := realType.FieldByName(e.StructField())
		// 通过标注获取表单项名称, 并翻译错误信息内容
		errMap[field.Tag.Get(tag)] = e.Translate(translate)
	}

	return errMap
}
