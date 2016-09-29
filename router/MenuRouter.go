package router

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessionauth"
	dbw "../db"
	"fmt"
	"github.com/martini-contrib/binding"
)
func init() {
	fmt.Println("加载 MenuRouter")
	routers.PushFront(MenuRouter)
}
func MenuRouter(m *martini.ClassicMartini){
	m.Get("/admin/menu",sessionauth.LoginRequired, func(r render.Render, menuDao dbw.MenuDao){
		menus,err := menuDao.All()
		if(err != nil){
			r.JSON(405, JsonMsg{"message":"get menu error"})
			return
		}
		r.Header().Add("X-Total-Count",fmt.Sprintf("%d",len(menus)))
		r.JSON(200,menus)
	})
	m.Post("/admin/menu",sessionauth.LoginRequired, binding.Bind(dbw.Menu{}), func(r render.Render, menu dbw.Menu, menuDao dbw.MenuDao) {
		fmt.Println("----------",menu)
		menu,err := menuDao.Save(menu)
		if(err != nil){
			r.JSON(405, JsonMsg{"message":"add menu error"})
			return
		}
		r.JSON(200,menu)
	})
	m.Put("/admin/menu",sessionauth.LoginRequired, binding.Bind(dbw.Menu{}), func(r render.Render, menu dbw.Menu, menuDao dbw.MenuDao) {
		fmt.Println("+++++++",menu)
		menu,err := menuDao.Save(menu)
		if(err != nil){
			r.JSON(405, JsonMsg{"message":"update menu error"})
			return
		}
		r.JSON(200,menu)
	})
	m.Delete("/admin/menu/:id",sessionauth.LoginRequired, func(r render.Render, menuDao dbw.MenuDao, param martini.Params) {
		id,ok := param["id"]
		if(!ok){
			r.JSON(405,JsonMsg{"message":"id not found"})
			return
		}
		err := menuDao.Remove(id)
		if(err != nil){
			r.JSON(405,JsonMsg{"message":"remove error"})
			return
		}
		r.JSON(200,JsonMsg{"message":"success"})
	})
}
