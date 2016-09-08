package db

import (
	"gopkg.in/mgo.v2/bson"
	"fmt"
	"gopkg.in/mgo.v2"
	"github.com/go-martini/martini"
	"log"
)

type NewsTag struct {
	Id bson.ObjectId `bson:"_id"`
	Tag string `bson:"Tag"`
}

type NewsTagDao interface{
	Add(tag string) error
	AddAll(tags []string) error
	FindAll() ([]NewsTag,error)
}

type newsTagDao struct{
	d *mgo.Database
	logger *log.Logger
}

func init() {
	fmt.Println("加载 NewsTagDao")
	DaoList.PushBack(func (d *mgo.Database,c martini.Context, logger *log.Logger){
		c.MapTo(&newsTagDao{d:d,logger:logger},(*NewsTagDao)(nil))
	})
}
func (nd *newsTagDao) FindAll() ([]NewsTag,error){
	var tags []NewsTag
	err := nd.d.C("NewsTag").Find(nil).All(&tags)
	if(err != nil){
		nd.logger.Fatalln("FindAll NewsTags error", err)
		return tags,err
	}
	return tags,nil
}
func (nd *newsTagDao) Add(tag string) error{
	var old NewsTag
	err := nd.d.C("NewsTag").Find(bson.M{"Tag":tag}).One(&old)
	if(err == nil){
		return nil
	}
	newsTag := NewsTag{Id:bson.NewObjectId(), Tag:tag}
	err = nd.d.C("NewsTag").Insert(newsTag)
	if(err != nil){
		nd.logger.Println("Add tag ",tag,"error.",err)
		return err
	}
	return nil
}
func (nd *newsTagDao) AddAll(tags []string) error  {
	for _,tag := range tags{
		err := nd.Add(tag)
		if(err != nil){
			nd.logger.Println("Add tags ", tags, " error.",err)
		}
	}
	return nil
}
