package db

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"fmt"
	"github.com/go-martini/martini"
	"log"
)

func init(){
	fmt.Println("加载 NewsTypeDao")
	DaoList.PushBack(func (d *mgo.Database,c martini.Context,logger *log.Logger){
		c.MapTo(&newsTypeDao{d:d,logger:logger},(*NewsTypeDao)(nil))
	})
}

type NewsType struct {
	Id_ bson.ObjectId `bson:"_id"`
	Title string `bson:"Title"`
	ParentId string `bson:"ParentId"`
}

type NewsTypeDao interface{
	FindAll() ([]NewsType,error)
	Map() map[string]string
	Parents() []NewsType
	SubMap() map[string] []NewsType
	FindByParent(parentId string) []NewsType
	FindOne(id bson.ObjectId) (NewsType,error)
}

type newsTypeDao struct {
	d *mgo.Database
	logger *log.Logger
}

// 查找所有类型
func (nd *newsTypeDao) FindAll() ([]NewsType,error){
	var ret []NewsType
	err := nd.d.C("NewsType").Find(nil).All(&ret)
	return ret,err
}
// 获取所有类型的Id、Title的Map
func (nd *newsTypeDao) Map() map[string]string{
	newsType,_ := nd.FindAll()
	ret := map[string]string{}
	for _,t := range newsType{
		ret[t.Id_.Hex()] = t.Title
	}
	return ret
}
// 返回所有大类
func (nd *newsTypeDao) Parents() []NewsType{
	return nd.FindByParent("")
}
// 返回大类Id与子类的Map
func (nd *newsTypeDao) SubMap() map[string] []NewsType{
	submap := map[string] []NewsType{}
	parents:=nd.Parents()
	for _,parent:=  range parents{
		subTypes := []NewsType {}
		nd.d.C("NewsType").Find(bson.M{"ParentId":parent.Id_.Hex()}).All(&subTypes)
		submap[parent.Id_.Hex()] = subTypes
	}
	return submap
}
// 根据父类查找子类
func (nd *newsTypeDao) FindByParent(parentId string) []NewsType{
	parents := []NewsType{}
	err := nd.d.C("NewsType").Find(bson.M{"ParentId":parentId}).All(&parents)
	if(err != nil){
		nd.logger.Println("NewsTypeDao findByParentId error", err)
	}
	return  parents
}

func (nd *newsTypeDao) FindOne(id bson.ObjectId) (NewsType,error){
	var newsType NewsType
	err := nd.d.C("NewsType").FindId(id).One(&newsType)
	if(err != nil){
		log.Println("NewsTypeDao FindOne error", id, err)
	}
	return newsType,err
}