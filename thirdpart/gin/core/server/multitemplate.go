package server

import (
	"path/filepath"
	"study/thirdpart/gin/core/conf"

	"github.com/gin-contrib/multitemplate"
)

// 获取多模板渲染对象
//
// 所谓多模板渲染, 即通过 go 1.6 之后的 `{{block "name" pipline}}` 模板语法, 可以组合多个模板为一个
//
// 多模板渲染解决了 HTML 每个模板中都需要重复定义相同内容的问题, 可以通过定义一个类似 `layout` 的模板, 并在模板中定义 `block` 占位符,
// 在其它模板中定义用于取代 `block` 占位符的内容
//
// 参考: <https://github.com/gin-contrib/multitemplate>
//
// `multitemplate` 库的使用方法即: 定义一个 `render.HTMLRender` 对象 (`multitemplate.Renderer` 类型从 `render.HTMLRender` 继承),
// 并设置每个渲染都是由那些模板文件组合而成, 例如: `index.tmpl` 是由 [`base.tmpl` 和 `index.tmpl`] 组成的
//
// 在本例中, `base.tmpl` 模板文件定义了页面的 "layout", 其它模板文件都需要和该文件组合
//
// 返回:
//   - `multitemplate.Renderer`, 渲染对象
func LoadTemplates() multitemplate.Renderer {
	// 获取 base.tmpl 文件的路径
	base := filepath.Join(conf.Config.Server.Template.TemplatesPath, "base.tmpl")

	// 查找其它的模板文件
	htmls, err := filepath.Glob(
		filepath.Join(conf.Config.Server.Template.TemplatesPath, "*"),
	)
	if err != nil {
		panic(err)
	}

	// 生成渲染对象
	r := multitemplate.NewRenderer()

	// 将其它模板文件和 base.tmpl 文件文件进行组合
	for _, html := range htmls {
		if html == base {
			continue
		}
		// 添加模板组合以及模板中可以使用的函数
		r.AddFromFilesFuncs(filepath.Base(html), Engine.FuncMap, base, html)
	}

	return r
}
