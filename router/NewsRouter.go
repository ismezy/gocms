package router

import (
	"fmt"
	"github.com/go-martini/martini"
	"log"
	"gopkg.in/mgo.v2/bson"
	dbw "../db"
	"net/http"
	"github.com/martini-contrib/binding"
	"strconv"
	"strings"
	"github.com/martini-contrib/sessionauth"
)

func init() {
	fmt.Println("加载 NewsRouter")
	routers.PushFront(NewsRouter)
}
func NewsRouter(m *martini.ClassicMartini){
	m.Get("/admin/news/list/:page",sessionauth.LoginRequired, func(r LayoutWrapper, req *http.Request ,params martini.Params, newsDao dbw.NewsDao, newsTypeDao dbw.NewsTypeDao) {
		// 分页数据
		page,ok := params["page"]
		ipage := 1
		if(ok){
			ipage,_ = strconv.Atoi(page)
		}

		// 栏目字典
		newsTypeMap := newsTypeDao.Map()
		columns:=newsTypeDao.Parents()
		ret := map[string] interface{} {"typemap":newsTypeMap,"columns":columns,"column":"","subcolumn":""}
		// 处理筛选
		col,cok := req.URL.Query()["column"]
		scol,sok := req.URL.Query()["subcolumn"]
		filter := bson.M{}
		column := ""
		if(sok && scol[0] != ""){
			ret["subcolumn"] = scol[0]
			filter["SubType"] = scol[0]
			st,_ := newsTypeDao.FindOne(bson.ObjectIdHex(scol[0]))
			ret["column"] = st.ParentId
			column = st.ParentId
		}
		if(cok){
			column = col[0]
		}
		if(column != ""){
			ret["column"] = column
			subcolumns := newsTypeDao.FindByParent(column)
			ret["subcolumns"] = subcolumns
			filter["Type"] = column
		}else{
			subcolumns := newsTypeDao.FindByParent(columns[0].Id_.Hex())
			ret["subcolumns"] = subcolumns
		}

		newspage := newsDao.FindPage(ipage,30,filter,"-CreateTime")
		ret["newspage"] = newspage
		r.HTML(200,"admin/newslist", ret, "admin")
	})
	m.Get("/admin/news",sessionauth.LoginRequired, func(r LayoutWrapper, logger *log.Logger,params martini.Params, newsTypeDao dbw.NewsTypeDao,
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
	m.Post("/admin/news",sessionauth.LoginRequired,binding.Bind(dbw.News{}),func(r LayoutWrapper, newsDao dbw.NewsDao,newsTypeDao dbw.NewsTypeDao, newsTagDao dbw.NewsTagDao,
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
	m.Get("/admin/news/preview",sessionauth.LoginRequired, func(r LayoutWrapper,newsDao dbw.NewsDao, newsTagDao dbw.NewsTagDao, newsTypeDao dbw.NewsTypeDao,req *http.Request,res http.ResponseWriter) {
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

