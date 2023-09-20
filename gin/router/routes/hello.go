package routes

import "github.com/gin-gonic/gin"

func HelloGet(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"result":  "OK",
		"content": "Hello World",
	})
    ctx.Next()
}
