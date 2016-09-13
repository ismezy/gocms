package router

import (
	"fmt"
	"github.com/go-martini/martini"
	"github.com/martini-contrib/sessionauth"
	"container/list"
	dbw "../db"
	"net/http"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
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
	m.Get("/admin", sessionauth.LoginRequired, func(r LayoutWrapper){
		r.HTML(200,"admin/index","admin", "admin")
	})
	m.Get("/login",func(wr LayoutWrapper,r render.Render,user sessionauth.User, req *http.Request){
		if(user.IsAuthenticated()){
			r.Redirect("/admin")
		}
		wr.HTML(200,"login","")
	})
	m.Post("/login",func(r render.Render, userDao dbw.LoginUserDao, session sessions.Session, req *http.Request){
		userName := req.FormValue("user")
		password := req.FormValue("password")
		u,err := userDao.FindByLoginId(userName);
		fmt.Println("21111111")
		if(err != nil){
			r.HTML(200,"login","error")
		}
		if(password == u.Password) {
			next, ok := req.URL.Query()["next"]
			url := "/admin"
			sessionauth.AuthenticateSession(session,&u)
			if (ok) {
				url = next[0]
			}
			fmt.Println("goto:", url)
			r.Redirect(url)
		}else{
			r.HTML(200,"login","error")
		}
	})
	m.Get("/logout",sessionauth.LoginRequired,func(r render.Render, user sessionauth.User){
		user.Logout()
		r.Redirect("/")
	})
}

func GetRouters() *list.List {
	return routers;
}