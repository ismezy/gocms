package router

import (
	"fmt"
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2"
	dbw "../db"
	"log"
	"gopkg.in/mgo.v2/bson"
	"github.com/martini-contrib/render"
	"net/http"
)

func init() {
	fmt.Println("加载typeManagerRouter")
	routers.PushFront(typeManagerRouter)
}

func typeManagerRouter(m *martini.ClassicMartini){
	// 类别页面
	m.Get("/admin/news/type",func(r LayoutWrapper, d *mgo.Database, logger *log.Logger ){
		ret := make(map[string] interface{})
		subMap := make(map[string] interface{})
		ret["subTypes"] = subMap
		var newsTypes []dbw.NewsType
		// 获取大类
		err := d.C("NewsType").Find(bson.M{"$or":[]bson.M{{"ParentId":nil}, {"ParentId":""}}}).All(&newsTypes)
		if(err != nil){
			logger.Fatalln(err)
		}
		ret["rootTypes"] = newsTypes;
		// 获取小类
		for _,t := range newsTypes{
			var subTypes []dbw.NewsType
			err := d.C("NewsType").Find(bson.M{"ParentId":t.Id_.Hex()}).All(&subTypes)
			if(err != nil){
				logger.Fatalln(err)
			}
			subMap[t.Id_.Hex()] = subTypes
		}
		fmt.Println(ret)
		r.HTML(200,"admin/newstype",ret,"admin")
	})
	m.Get("/admin/news/type/:id/subtypes", func(r render.Render,d *mgo.Database, logger *log.Logger,params martini.Params) {
		id,ok := params["id"]
		if(!ok){
			r.Text(405,"id cannot be empty")
			return
		}
		var list []dbw.NewsType
		d.C("NewsType").Find(bson.M{"ParentId":id}).All(&list)
		r.JSON(200,list)
	})
	// 保存类别
	m.Post("/admin/news/type",func(r render.Render,d *mgo.Database, logger *log.Logger,params martini.Params,req *http.Request ){
		title := req.FormValue("title")
		id := req.FormValue("id")
		parentId := req.FormValue("parentId")
		// 如果标题为空返回错误
		if(title == ""){
			r.Error(405)
			return
		}
		newType := dbw.NewsType{}
		newType.Title = title
		if(parentId != ""){
			newType.ParentId = parentId
		}
		var err error
		if(id != ""){
			newType.Id_ = bson.ObjectIdHex(id)
			err = d.C("NewsType").UpdateId(newType.Id_, newType)
			logger.Println("update newstype", newType)
		}else{
			newType.Id_ = bson.NewObjectId()
			err = d.C("NewsType").Insert(newType)
			logger.Println("insert newstype", newType)
		}
		if(err != nil){
			fmt.Println(err)
			r.JSON(405,"save newstype faild")
			return
		}
		r.JSON(200,newType)
	})
	// 删除类别
	m.Delete("/admin/news/type/:id", func(r render.Render,d *mgo.Database, logger *log.Logger,params martini.Params,req *http.Request) {
		id,ok := params["id"]
		if(!ok){
			r.JSON(405,"id cannot be empty")
		}
//		newType := db.NewsType{}
//		d.C("NewsType").FindId(bson.ObjectIdHex(id)).One(&newType)
//		if(newType.ParentId != nil && newType.ParentId != ""){
		d.C("NewsType").Remove(bson.M{"ParentId":id})
//		}
		err := d.C("NewsType").RemoveId(bson.ObjectIdHex(id))
		if(err != nil){
			logger.Fatalln("remove NewsType failed, id is %s. Reason:", id, err)
			r.JSON(405,"remove NewsType by id[" + id + "] failed")
		}
		r.JSON(200,"success")
	})
}
