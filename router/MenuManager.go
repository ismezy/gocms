package router

import "github.com/go-martini/martini"

func menuManagerRouter(m *martini.ClassicMartini){
	m.Get("/admin/menu", func(r LayoutWrapper){
		r.HTML(200,"admin/menu","","admin")
	})
}
