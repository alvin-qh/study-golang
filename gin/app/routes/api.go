package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 获取用户 API, 返回 `User` 类型 JSON 对象
func ApiGetUser(ctx *gin.Context) {
	name := ctx.DefaultQuery("name", "Alvin")

	// 返回 JSON 结果
	// `JSON` 表示 JSON 中非 ASCII 编码的文本会使用类似 `\uXXXX` 的编码格式
	ctx.JSON(200, &User{
		Name:     name,
		Gender:   GenderM,
		Birthday: time.Date(1981, 3, 17, 0, 0, 0, 0, time.UTC),
	})
}

// 创建用户 API, 返回 `User` 类型 JSON 对象
func ApiPostUser(ctx *gin.Context) {
	var form UserForm
	// 将请求 body 转为 `UserForm` 类型对象
	if err := ctx.ShouldBindJSON(&form); err != nil {
		// 返回错误, 包含验证错误结果
		ctx.AbortWithStatusJSON(http.StatusBadRequest, NewInputError(err.(validator.ValidationErrors), &form))
		return
	}

	// 返回 JSON 结果
	// `PureJSON` 表示 JSON 为纯文本, 不使用 `\uXXXX` 这类的编码
	ctx.PureJSON(http.StatusOK, form.toUser())
}
