package router

import (
	"fmt"
	"github.com/go-martini/martini"
	"container/list"
	dbw "../db"
)
func init() {
	fmt.Println("加载index router")
	routers.PushBack(indexRouters)
}

func indexRouters(m *martini.ClassicMartini){
	m.Get("/", func(r LayoutWrapper,menuDao dbw.MenuDao){
		menus,_ := menuDao.All()
		model := map[string] interface{}{
			"menus":menus,
		}

	 	r.HTML(200,"index", model)
	})
	m.Get("/admin",func(r LayoutWrapper){
		r.HTML(200,"admin/index","admin", "admin")
	})
}

func GetRouters() *list.List {
	return routers;
}