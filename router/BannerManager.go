package router

import (
	"github.com/go-martini/martini"
	"fmt"
)

func init() {
	fmt.Println("加载 BannerManagerRouter")
	routers.PushBack(BannerManagerRouter)
}

func BannerManagerRouter(m *martini.ClassicMartini)  {
	m.Get("/admin/banner", func(r LayoutWrapper) {
		r.HTML(200,"admin/banner", "","admin")
	})
}
