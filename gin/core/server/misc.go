package server

import (
	"io"

	"github.com/gin-gonic/gin"
)

func DisableGinLogger() {
	gin.SetMode(gin.ReleaseMode)

	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard
}
