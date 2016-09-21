package router

import (
	"fmt"
	"github.com/go-martini/martini"
	dbw "../db"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/binding"
)

func init() {
	fmt.Println("加载TileRouter")
	routers.PushFront(TileRouter)
}

func TileRouter(m *martini.ClassicMartini){
	m.Get("/admin/tile", sessionauth.LoginRequired, func(r render.Render, tileDao dbw.TileDao) {
		tile,err:= tileDao.Get()
		if(err != nil){
			r.JSON(405,map[string]string{"error":"获取tile出错"})
		}else{
			r.Header()["X-Total-Count"] = []string{fmt.Sprintf("%d",len(tile.Items))}
			r.JSON(200,tile.Items)
		}
	})
	m.Post("/admin/tile", sessionauth.LoginRequired, binding.Bind(dbw.TileItem{}),func(r render.Render, tileDao dbw.TileDao, tile dbw.TileItem){
		t,err := tileDao.Insert(tile)
		if(err != nil) {
			r.JSON(405, map[string]string{"error":"新增磁贴出错"})
		}else{
			r.JSON(200,t)
		}
	})
	m.Put("/admin/tile", sessionauth.LoginRequired, binding.Bind(dbw.TileItem{}),func(r render.Render, tileDao dbw.TileDao, tile dbw.TileItem){
		err := tileDao.Upload(tile)
		if(err != nil){
			r.JSON(405, map[string]string{"error":"修改磁贴出错"})
		}else{
			r.JSON(200,tile)
		}
	})
}
