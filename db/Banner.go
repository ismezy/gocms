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
	Id bson.ObjectId `bson:"_id" form:"Id"`
	ImgPath string `bson:"ImgPath"`
	Url string `bson:"Url" form:"Url"`
	Status string `bson:"Status"`
	CreateTime time.Time `bson:"CreateTime"`
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
	FindTop(topN int)([]Banner,error)
	FindAll()([]Banner,error)
	FindOne(id bson.ObjectId) (Banner,error)
	Save(banner Banner,newFile bool) error
	RemoveId(id bson.ObjectId) error
}

func (bd *bannerDaoImpl) FindAll()([]Banner,error){
	var banners [] Banner
	err := bd.d.C("Banner").Find(nil).Sort("-CreateTime").All(&banners)
	if(err != nil){
		bd.logger.Println("Banner FindAll error", err)
	}
	return banners,err
}

func (bd *bannerDaoImpl) FindTop(topN int)([]Banner,error){
	var banners [] Banner
	err := bd.d.C("Banner").Find(nil).Sort("-CreateTime").Limit(topN).All(&banners)
	if(err != nil){
		bd.logger.Println("Banner FindTop error", err, topN)
	}
	return banners,err
}

func (bd *bannerDaoImpl) FindOne(id bson.ObjectId) (Banner,error){
	var banner Banner
	err := bd.d.C("Banner").FindId(id).One(&banner)
	return banner,err
}
func (bd *bannerDaoImpl) Save(banner Banner, newFile bool) error{
	if(banner.Id != ""){	// 修改
		var old Banner
		err := bd.d.C("Banner").FindId(banner.Id).One(&old)
		if(err != nil){
			bd.logger.Println("Banner update error, ", banner.Id , " not found", err)
			return err
		}
		banner.CreateTime = old.CreateTime
		if(!newFile){
			banner.ImgPath = old.ImgPath
		}
//		banner.Status = old.Status
		err = bd.d.C("Banner").Update(bson.M{"_id":banner.Id}, banner)
		if(err != nil){
			bd.logger.Println("Banner update error, ", banner,  err)
			return err
		}
	}else{	// 新增
		banner.Id = bson.NewObjectId()
		banner.CreateTime = time.Now()
		err := bd.d.C("Banner").Insert(banner);
//		banner.Status = "启用"
		if(err != nil){
			bd.logger.Println("Banner insert error.", banner, err)
			return err
		}
	}
	return nil
//	bd.d.C("Banner")
}
func (bd *bannerDaoImpl) RemoveId(id bson.ObjectId) error{
	err := bd.d.C("Banner").RemoveId(id)
	if(err != nil){
		bd.logger.Println("Banner removeid error",id,err)
	}
	return err
}
