package router

import (
	"github.com/martini-contrib/render"
	"github.com/go-martini/martini"
	dbw "../db"
	"html/template"
	"net/http"
)

type LayoutWrapper interface {
	HTML(status int, name string, v interface{}, layout ...string)
}

type wrapper struct {
	r render.Render
}

func (w *wrapper) HTML(status int, name string, v interface{}, layout...string){
	opt := render.HTMLOptions{}
	if(len(layout) > 0){
		opt.Layout = "layout/" + layout[0] + "/layout";
	}else {
		opt.Layout = "layout/default/layout";
	}
	w.r.HTML(status,name,v,opt)
}

func Wapper(r render.Render,c martini.Context,menuDao dbw.MenuDao,req *http.Request){
	menus,_ := menuDao.All()
	menuHtml := ""
	for _,menu := range menus{
		if(req.URL.Path == menu.Path){
			menuHtml += `<li class="active"><a href="` + menu.Path+ `">` + menu.Title + `</a></li>`
		}else{
			menuHtml += `<li><a href="` + menu.Path+ `">` + menu.Title + `</a></li>`
		}
	}
	r.Template().Funcs(template.FuncMap{"MainMenu":func() template.HTML{
		return template.HTML(menuHtml)
	}})
	c.MapTo(&wrapper{r:r},(*LayoutWrapper)(nil))
}
