package controllers

import (
	"github.com/astaxie/beego"
)

type HomeController struct {
	beego.Controller
}

// @router / [get]
func (_this *HomeController) Get() {
	_this.Ctx.WriteString("Hello")
}
