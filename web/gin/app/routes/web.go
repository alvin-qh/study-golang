package routes

import (
	"fmt"
	"net/http"
	"time"

	"study/web/gin/core/server"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 获取用户信息页面
//
// gin 框架使用了 Golang 内置的 `html/template` 包进行 HTML 模板渲染
func WebGetUsers(ctx *gin.Context) {
	// 通过 gin.Context 的 HTML 方法渲染一个 HTML 模板
	// 传入响应状态, 模板名称和模板参数, 其中:
	//   - 模板参数是一个 `map[string]any` 集合, 该集合在模板中使用 `.` 引用, 即 `title` 表示为 `.title`
	ctx.HTML(http.StatusOK, "users.tmpl", gin.H{
		"title": "User List",
		"users": []User{
			{
				Id:       "001",
				Name:     "Alvin",
				Gender:   GenderM,
				Birthday: time.Date(1981, 3, 17, 0, 0, 0, 0, time.UTC),
			},
			{
				Id:       "002",
				Name:     "Emma",
				Gender:   GenderF,
				Birthday: time.Date(1985, 3, 29, 0, 0, 0, 0, time.UTC),
			},
		},
	})
}

// 获取用户编辑页面
func WebNewUsers(ctx *gin.Context) {
	// 渲染用户编辑页面
	ctx.HTML(http.StatusOK, "users_new.tmpl", gin.H{
		"title": "New User",
		"user": &UserForm{
			Gender:     GenderM,
			BirthYear:  1960,
			BirthMonth: 1,
			BirthDay:   1,
		},
	})
}

// 提交用户信息
func WebPostUsers(ctx *gin.Context) {
	var user UserForm

	// 提交表单转为 `UserForm` 类型结构体
	err := ctx.ShouldBind(&user)
	if err != nil {
		// 表单验证失败, 渲染编辑页面并展示错误信息
		ctx.HTML(http.StatusBadRequest, "users_new.tmpl", gin.H{
			"title": "New User",
			"user":  user,
			"errs":  server.MappedValidatorErrors(err.(validator.ValidationErrors), &user, "form"),
		})
		return
	}

	// 提交成功, 渲染用户展示页面
	ctx.HTML(http.StatusOK, "users_detail.tmpl", gin.H{
		"title": "User",
		"user":  user.toUser("003"),
	})
}

// 获取用户编辑页面
func WebEditUsers(ctx *gin.Context) {
	var user *UserForm

	// 获取路径中的 :id 参数
	id := ctx.Param("id")

	// 根据 id 值获取不同的 UserForm 对象
	switch id {
	case "001":
		user = &UserForm{
			Name:       "Alvin",
			Gender:     GenderM,
			BirthYear:  1981,
			BirthMonth: 3,
			BirthDay:   17,
		}
	case "002":
		user = &UserForm{
			Name:       "Emma",
			Gender:     GenderF,
			BirthYear:  1985,
			BirthMonth: 3,
			BirthDay:   29,
		}
	default:
		// id 参数错误, 返回 404 页面
		ctx.AbortWithError(http.StatusNotFound, fmt.Errorf("user not exit by id \"%v\"", id))
		return
	}

	// 渲染 `users_edit.tmpl` 页面
	ctx.HTML(http.StatusOK, "users_edit.tmpl", gin.H{
		"title": "User Edit",
		"id":    id,
		"user":  user,
	})
}

// 修改用户信息
func WebPutUsers(ctx *gin.Context) {
	// 获取待修改用户的 id 参数
	id := ctx.Param("id")

	var user UserForm

	// 提交表单转为 `UserForm` 类型结构体
	err := ctx.ShouldBind(&user)
	if err != nil {
		// 表单验证失败, 渲染编辑页面并展示错误信息
		ctx.HTML(http.StatusBadRequest, "users_new.tmpl", gin.H{
			"title": "User Editor",
			"user":  user,
			"errs":  server.MappedValidatorErrors(err.(validator.ValidationErrors), &user, "form"),
		})
		return
	}

	// 提交成功, 渲染用户展示页面
	ctx.HTML(http.StatusOK, "users_detail.tmpl", gin.H{
		"title": "User",
		"user":  user.toUser(id),
	})
}
