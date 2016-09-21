package db

import (
	"gopkg.in/mgo.v2/bson"
	"time"
	"gopkg.in/mgo.v2"
	"fmt"
	"github.com/go-martini/martini"
	"log"
)

type News struct {
	Id bson.ObjectId `bson:"_id"`
	Title string `bson:"Title" form:"Title"`
	Content string `bson:"Content" form:"Content"`
	Status string `bson:"Status" form:"Status"`
	CreateTime time.Time `bson:"CreateTime"`
	PublishTime time.Time `bson:"PublishTime"`
	Type string `bson:"Type" form:"Type"`
	SubType string `bson:"SubType" form:"SubType"`
	Tags string `bson:"Tags" form:"Tags"`
}

func (n News) InOneMonth() bool{
	newTime := n.PublishTime
	newTime.AddDate(0,1,0)
	return newTime.Before(time.Now())
}

func init() {
	fmt.Println("加载 NewsDao")
	DaoList.PushBack(func (d *mgo.Database,c martini.Context, logger *log.Logger){
		c.MapTo(&newsDao{d:d,logger:logger},(*NewsDao)(nil))
	})
}
type NewsDao interface {
	FindOne(id string)  (News,error)
	Save(news News) (News,error)
	FindByTypeId(id string)([]News,error)
	FindPage(page int,pageSize int,filter map[string]interface{},sort...string) Page
	FindPageByType(page int,pageSize int,typeId string,sort...string) Page
}
type newsDao struct {
	d *mgo.Database
	logger *log.Logger
}

func (nd *newsDao) FindOne(id string)  (News,error) {
	news := News{}
	err := nd.d.C("News").FindId(bson.ObjectIdHex(id)).One(&news)
	if(err != nil){
		nd.logger.Println("NewsDao FindId error ", err,id)
		return news,err
	}
	return news,nil
}

func (nd *newsDao) Save(news News) (News,error){
	if(news.Status == "发布"){
		news.PublishTime = time.Now()
	}
	if(news.Id.Hex() == ""){
		news.CreateTime = time.Now()
		news.Id = bson.NewObjectId()
		err := nd.d.C("News").Insert(news)
		if(err != nil){
			nd.logger.Println("insert news error,", err, news)
			return news,err
		}
	}else{
		// 获取更新前的数据
		old,err := nd.FindOne(news.Id.Hex())
		if(err != nil){
			nd.logger.Println("update news error", err, news)
			return news,err
		}
		// 不能从界面转递的数据从表中获取
		news.CreateTime = old.CreateTime
		err = nd.d.C("News").Update(bson.M{"_id": news.Id},news)
		if(err != nil){
			nd.logger.Println("update news error,", err, news)
			return news,err
		}
	}
	return news,nil
}

func (nd *newsDao) FindByTypeId(id string)([]News,error){
	var news []News
	err := nd.d.C("News").Find(bson.M{"$or":[]bson.M{{"Type":id},{"SubType":id}}}).Sort("-CreateTime").All(&news)
	if(err != nil){
		log.Println("NewsDao findByTypeId error", id, err)
	}
	return news,err
}

func (nd *newsDao) FindPage(page int,pageSize int,filter map[string]interface{},sort...string) Page{
	fmt.Println(filter)
	var news []News
	rows,_ := nd.d.C("News").Find(filter).Count()
	nd.d.C("News").Find(filter).Sort(sort...).Skip(pageSize * (page - 1)).Limit(pageSize).All(&news)
	newspage := Page{List:news, TotalRows:rows, Page:page, PageSize: pageSize}
	return newspage
}
func (nd *newsDao) FindPageByType(page int,pageSize int,typeId string,sort...string) Page{
	filter := bson.M{"Status":"发布","$or":[]bson.M{{"Type":typeId},{"SubType":typeId}}}
	return nd.FindPage(page,pageSize,filter,sort...)
}
