package router

import (
	"fmt"
	"github.com/go-martini/martini"
	dbw "../db"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"net/http"
)

func init() {
	fmt.Println("加载 columnRouter")
	routers.PushBack(columnRouter)
}

func columnRouter(m *martini.ClassicMartini){
	m.Get("/column/:id", func(r LayoutWrapper,req *http.Request, columnDao dbw.NewsTypeDao, newsDao dbw.NewsDao,params martini.Params) {
		model := map[string]interface{} {}
		id,ok := params["id"]
		page,pok := req.URL.Query()["page"]
		ipage := 1
		if(!ok){
			return
		}
		if(pok){
			ipage,_ = strconv.Atoi(page[0])
		}
		newsType,err := columnDao.FindOne(bson.ObjectIdHex(id))
		if(err != nil){
			return
		}
		if(newsType.ParentId == ""){
			subTypes := columnDao.FindByParent(id)
			model["subCols"] = subTypes
		}else{
			subTypes := columnDao.FindByParent(newsType.ParentId)
			model["subCols"] = subTypes
			model["subCol"] = newsType
			newsType,_ = columnDao.FindOne(bson.ObjectIdHex(newsType.ParentId))
		}
		model["col"] = newsType
		model["id"] = id
		newspage := newsDao.FindPageByType(ipage, 20, id, "-CreateTime")
		model["news"] = newspage
		r.HTML(200, "column", model,"column")
	})
}
