package routes

import (
	"net/http"
	"study-gin/core/server"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 获取用户信息页面
//
// gin 框架使用了 Golang 内置的 `html/template` 包进行 HTML 模板渲染
func GetUser(ctx *gin.Context) {
	// 通过 gin.Context 的 HTML 方法渲染一个 HTML 模板
	// 传入响应状态, 模板名称和模板参数, 其中:
	//   - 模板参数是一个 `map[string]any` 集合, 该集合在模板中使用 `.` 引用, 即 `title` 表示为 `.title`
	ctx.HTML(http.StatusOK, "user.tmpl", gin.H{
		"title": "User",
		"user": &User{
			Name:     "Alvin",
			Gender:   "M",
			Birthday: time.Date(1981, 3, 17, 0, 0, 0, 0, time.UTC),
		},
		"list": []string{"A", "B", "C"},
	})
}

// 获取用户编辑页面
func GetUserEditor(ctx *gin.Context) {
	// 渲染用户编辑页面
	ctx.HTML(http.StatusOK, "user_editor.tmpl", gin.H{
		"title": "User Editor",
		"user": &UserForm{
			Gender:     "M",
			BirthYear:  1960,
			BirthMonth: 1,
			BirthDay:   1,
		},
	})
}

// 提交用户信息
func PostUser(ctx *gin.Context) {
	var user UserForm

	// 提交表单转为 `UserForm` 类型结构体
	err := ctx.ShouldBind(&user)
	if err != nil {
		// 表单验证失败, 渲染编辑页面并展示错误信息
		ctx.HTML(http.StatusOK, "user_editor.tmpl", gin.H{
			"title": "User Editor",
			"user":  user,
			"errs":  server.MappedValidatorErrors(err.(validator.ValidationErrors), &user),
		})
		return
	}

	// 提交成功, 渲染用户展示页面
	ctx.HTML(http.StatusOK, "user.tmpl", gin.H{
		"title": "User",
		"user":  user.toUser(),
		"list":  []string{"A", "B", "C"},
	})
}
