package routes

import (
	"fmt"
	"net/http"

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
		var er ErrorResult
		er.Code = InputErrorCode

		er.FromValidator(err.(validator.ValidationErrors), &agreeing)
		ctx.AbortWithStatusJSON(http.StatusBadRequest, &er)
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
