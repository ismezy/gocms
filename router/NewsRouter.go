package router

import (
	"fmt"
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2"
	"log"
	"gopkg.in/mgo.v2/bson"
	dbw "../db"
	"net/http"
	"github.com/martini-contrib/binding"
	"strconv"
	"strings"
)

func init() {
	fmt.Println("加载 NewsRouter")
	routers.PushFront(NewsRouter)
}
func NewsRouter(m *martini.ClassicMartini){
	m.Get("/admin/news/list/:page", func(r LayoutWrapper, d *mgo.Database, logger *log.Logger,params martini.Params, newsDao dbw.NewsDao, newsTypeDao dbw.NewsTypeDao) {
		// 分页数据
		page,ok := params["page"]
		ipage := 1
		if(ok){
			ipage,_ = strconv.Atoi(page)
		}
		/*
		var news []dbw.News
		rows,_ := d.C("News").Find(nil).Count()
		d.C("News").Find(nil).Sort("-CreateTime").Skip(30 * (ipage - 1)).Limit(30).All(&news)
		var pageCount int = rows / 30 + ((rows % 30 + 1) / (rows % 30 + 1))
		newspage := map[string] interface {} {"list":news, "totalRows":rows, "page":ipage, "pageCount": pageCount}
		*/
		newspage := newsDao.FindPage(ipage,30,nil,"-CreateTime")

		// 栏目字典
		newsTypeMap := newsTypeDao.Map()

		ret := map[string] interface{} {"newspage": newspage,"typemap":newsTypeMap}

		fmt.Println(ret)
		r.HTML(200,"admin/newslist", ret, "admin")
	})
	m.Get("/admin/news", func(r LayoutWrapper, logger *log.Logger,params martini.Params, newsTypeDao dbw.NewsTypeDao,
		newsTagDao dbw.NewsTagDao,newsDao dbw.NewsDao,req *http.Request) {
		types := newsTypeDao.Parents()
		subTypes :=[]dbw.NewsType{}
		tags,_ := newsTagDao.FindAll()
		model := map[string] interface{} {"types":types,"tags":tags}
		id,ok := req.URL.Query()["id"]
		if(ok){
			news,_:= newsDao.FindOne(id[0])
			fmt.Println(news)
			subTypes = newsTypeDao.FindByParent(news.Type)
			model["news"] = news
		}else if(len(types) > 0){
			subTypes = newsTypeDao.FindByParent(types[0].Id_.Hex())
			model["news"] = dbw.News{}
		}
		model["subTypes"] = subTypes
		r.HTML(200,"admin/news", model, "admin")
	})
	m.Post("/admin/news",binding.Bind(dbw.News{}),func(r LayoutWrapper, newsDao dbw.NewsDao,newsTypeDao dbw.NewsTypeDao, newsTagDao dbw.NewsTagDao,
		logger *log.Logger, req *http.Request, res http.ResponseWriter, news dbw.News){
		fmt.Println(news)
		if(news.Status == "预览"){
			typeMap := newsTypeDao.Map()
			model := map[string]interface{}{"news":news, "type":typeMap}
			r.HTML(200,"admin/newspreview",model,"admin")
			return
		}
		tags := strings.Split(news.Tags, ",")
		err := newsTagDao.AddAll(tags)
		if(err != nil){
			http.Redirect(res,req,"/admin/news/list/1?msg=error",302)
		}

		id := req.FormValue("Id");
		if(id != ""){
			news.Id = bson.ObjectIdHex(id)
		}
		_,err = newsDao.Save(news)
		if(err != nil){
			http.Redirect(res,req,"/admin/news/list/1?msg=error",302)
		}else{
			http.Redirect(res,req,"/admin/news/list/1?msg=success",302)
		}
	})
	m.Get("/admin/news/preview", func(r LayoutWrapper,newsDao dbw.NewsDao, newsTagDao dbw.NewsTagDao, newsTypeDao dbw.NewsTypeDao,req *http.Request,res http.ResponseWriter) {
		id,ok := req.URL.Query()["id"]
		if(!ok){
			http.Redirect(res,req,"/admin/news/list/1?msg=error",302)
			return
		}
		news,err := newsDao.FindOne(id[0])
		if(err != nil){
			http.Redirect(res,req,"/admin/news/list/1?msg=error",302)
			return
		}
		typeMap := newsTypeDao.Map()
		model := map[string]interface{}{"news":news, "type":typeMap}
		r.HTML(200, "admin/newspreview", model,"admin")
	})
	m.Get("/news/:id",func(r LayoutWrapper, newsDao dbw.NewsDao, columnDao dbw.NewsTypeDao, params martini.Params) {
		id,ok:= params["id"]
		if(!ok){
			return
		}
		news,err := newsDao.FindOne(id)
		if(err != nil){
			return
		}
		newsType ,_ := columnDao.FindOne(bson.ObjectIdHex(news.Type))
		subType,_ := columnDao.FindOne(bson.ObjectIdHex(news.SubType))
		subTypes := columnDao.FindByParent(news.Type)
		model := map[string] interface{}{
			"col": newsType,
			"subCol":subType,
			"subCols": subTypes,
			"news": news,
		}
		r.HTML(200,"news", model, "column")
	})
}

