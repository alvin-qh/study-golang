package server

import (
	"html/template"
	"io"
	"time"

	"study-gin/core/conf"

	"github.com/gin-gonic/gin"
)

// 禁用 gin 框架内置的日志
func DisableGinLogger() {
	gin.SetMode(gin.ReleaseMode)

	gin.DefaultWriter = io.Discard
	gin.DefaultErrorWriter = io.Discard
}

// 设置模板
//
// 参数:
//   - `engine` (`*gin.Engine`): gin 框架 http 对象
func SetupTemplate(engine *gin.Engine) {
	if conf.Config.Server.Template.Enable {
		// 读取模板文件, 由于本例中为多模板, 所以无需执行
		// engine.LoadHTMLGlob(fmt.Sprintf("%v/*", conf.Config.Server.Template.TemplatesPath))

		// 设置静态文件路径和 publicPath
		engine.Static(conf.Config.Server.Template.StaticBaseURI, conf.Config.Server.Template.StaticPath)

		// 设置模板中的自定义函数
		engine.SetFuncMap(template.FuncMap{
			"date": func(t *time.Time) string {
				return t.Format(time.DateOnly)
			},
		})

		// 多模板渲染, 加载多模板渲染对象
		engine.HTMLRender = LoadTemplates()
	}
}
