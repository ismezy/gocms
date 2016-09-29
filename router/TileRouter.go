package router

import (
	"fmt"
	"github.com/go-martini/martini"
	dbw "../db"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/binding"
	"strconv"
)

func init() {
	fmt.Println("加载TileRouter")
	routers.PushFront(TileRouter)
}

func TileRouter(m *martini.ClassicMartini){
	m.Get("/admin/tile/item", sessionauth.LoginRequired, func(r render.Render, tileDao dbw.TileDao) {
		tile,err:= tileDao.Get()
		if(err != nil){
			r.JSON(405,map[string]string{"message":"获取tile出错"})
		}else{
			r.Header()["X-Total-Count"] = []string{fmt.Sprintf("%d",len(tile.Items))}
			r.JSON(200,tile.Items)
		}
	})
	m.Post("/admin/tile/item", sessionauth.LoginRequired, binding.Bind(dbw.TileItem{}),func(r render.Render, tileDao dbw.TileDao, tile dbw.TileItem){
		t,err := tileDao.Insert(tile)
		if(err != nil) {
			r.JSON(405, map[string]string{"message":"新增磁贴出错"})
		}else{
			r.JSON(200,t)
		}
	})
	m.Put("/admin/tile/item", sessionauth.LoginRequired, binding.Bind(dbw.TileItem{}),func(r render.Render, tileDao dbw.TileDao, tile dbw.TileItem){
		err := tileDao.Upload(tile)
		if(err != nil){
			r.JSON(405, map[string]string{"message":"修改磁贴出错"})
		}else{
			r.JSON(200,tile)
		}
	})
	m.Delete("/admin/tile/item/:index", sessionauth.LoginRequired, func(r render.Render, tileDao dbw.TileDao, param martini.Params) {
		index,ok := param["index"]
		if(!ok){
			r.JSON(400,JsonMsg{"message":"缺少参数"})
			return
		}
		idx,perr := strconv.Atoi(index)
		if(perr != nil){
			r.JSON(500,JsonMsg{"message":"无效的参数"})
			fmt.Println(perr)
		}
		err := tileDao.RemoveItem(idx)
		if(err != nil){
			r.JSON(500, JsonMsg{"message":"删除失败"})
			return
		}
		r.JSON(200, JsonMsg{"message":"删除成功"})
	})
}
