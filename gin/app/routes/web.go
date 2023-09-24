package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type Gender string

type User struct {
	Name     string
	Gender   Gender
	Birthday time.Time
}

func RenderHTML(ctx *gin.Context) {
	ctx.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "HTML Render",
		"user": &User{
			Name:     "Alvin",
			Gender:   "M",
			Birthday: time.Date(1981, 3, 17, 0, 0, 0, 0, time.UTC),
		},
		"list": []string{"A", "B", "C"},
	})
}
