package routers

import (
	"beego/www/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.HomeController{})

	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/home",
			beego.NSRouter("/", &controllers.HomeController{})),
	)
	beego.AddNamespace(ns)
}
