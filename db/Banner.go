package db

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"log"
	"fmt"
	"github.com/go-martini/martini"
	"time"
)

type Banner struct {
	Id bson.ObjectId `bson:"_id"`
	ImgPath string `bson: "ImgPath"`
	Url string `bson: "Url"`
	Status string `bons: "Status"`
	CreateTime time.Time `bson: "CreateTime"`
}

type bannerDaoImpl struct {
	d *mgo.Database
	logger *log.Logger
}

func init() {
	fmt.Println("加载 BannerDao")
	DaoList.PushBack(func (d *mgo.Database,c martini.Context, logger *log.Logger){
		c.MapTo(&bannerDaoImpl{d:d,logger:logger},(*BannerDao)(nil))
	})
}

type BannerDao interface {
}

func (bd *bannerDaoImpl) FindOne(id bson.ObjectId) (Banner,error){
	var banner Banner
	err := bd.d.C("Banner").FindId(id).One(&banner)
	return banner,err
}
func (bd *bannerDaoImpl) Save(banner *Banner) error{
	if(banner.Id.Valid()){	// 修改
		var old Banner
		err := bd.d.C("Banner").FindId(banner.Id).One(&old)
		if(err != nil){
			bd.logger.Fatalln("Banner update error, ", banner.Id , " not found", err)
			return err
		}
		banner.CreateTime = old.CreateTime
		err = bd.d.C("Banner").Update(bson.M{"_id":banner.Id}, banner)
		if(err != nil){
			bd.logger.Fatalln("Banner update error, ", banner,  err)
			return err
		}
	}else{	// 新增
		banner.CreateTime = time.Now()
		err := bd.d.C("Banner").Insert(banner);
		if(err != nil){
			bd.logger.Fatalln("Banner insert error.", banner, err)
			return err
		}
	}
	return nil
//	bd.d.C("Banner")
}
