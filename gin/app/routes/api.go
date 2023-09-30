package routes

import (
	"fmt"
	"net/http"
	"study-gin/core/server"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func GetHello(ctx *gin.Context) {
	name := ctx.DefaultQuery("name", "Alvin")
	ctx.JSON(200, gin.H{
		"name": name,
		"text": fmt.Sprintf("Hello %v", name),
	})
}

func PostHello(ctx *gin.Context) {
	var agreeing Agreeing
	if err := ctx.ShouldBindJSON(&agreeing); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &ErrorResult{
			Code:   InputErrorCode,
			Errors: server.MappedValidatorErrors(err.(validator.ValidationErrors), &agreeing, "json"),
		})
		return
	}

	var title string
	if agreeing.Gender == "M" {
		title = "先生"
	} else {
		if agreeing.Age < 40 {
			title = "小姐"
		} else {
			title = "女士"
		}
	}

	ctx.PureJSON(http.StatusOK, &Answer{
		Ok:     true,
		Answer: fmt.Sprintf("你好, %v%v", agreeing.Name, title),
	})
}
