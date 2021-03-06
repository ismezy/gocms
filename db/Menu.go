package db

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"log"
	"fmt"
	"github.com/go-martini/martini"
)

type Menu struct {
	Id bson.ObjectId `bson:"_id" json:"Id"`
	Title string `bson:"Title" json:"Title"`
	Path string `bson:"Path" json:"Path"`
	Code string `bson:"Code" json:"Code"`
	Index int `bson:"Index" json:"Index"`
	SubMenu []Menu `bson:"SubMenu"`
}

type MenuDao interface {
	Remove(id string) error
	All() ([]Menu,error)
	FindOne(id string) (Menu,error)
	Save(menu Menu) (Menu,error)
}

type menuDaoImpl struct {
	d *mgo.Database
	logger *log.Logger
}

func init() {
	fmt.Println("加载 MenuDao")
	DaoList.PushBack(func (d *mgo.Database,c martini.Context, logger *log.Logger){
		c.MapTo(&menuDaoImpl{d:d,logger:logger},(*MenuDao)(nil))
	})
}

func (md *menuDaoImpl) All() ([]Menu,error){
	var menu []Menu
	err := md.d.C("Menu").Find(nil).Sort("Index").Sort("Index").All(&menu)
	if(err != nil){
		md.logger.Println("FindAll Menu Error ", err)
		return menu,err
	}
	return menu,nil
}

func (md *menuDaoImpl) FindOne(id string) (Menu,error){
	var menu Menu
	err := md.d.C("Menu").FindId(bson.ObjectIdHex(id)).One(&menu)
	if(err != nil){
		md.logger.Println("FindOne Menu Error ", err)
		return menu,err
	}
	return menu,nil
}
func (md *menuDaoImpl) Save(menu Menu) (Menu,error) {
	if(menu.Id.Valid()){	// 更新
		err := md.d.C("Menu").UpdateId(menu.Id,menu)
		if(err != nil){
			md.logger.Println("Update Menu Error ", err)
			return menu,err
		}
	}else{		// 新增
		menu.Id = bson.NewObjectId()
		err := md.d.C("Menu").Insert(menu)
		if(err != nil){
			md.logger.Println("Insert Menu Update Error ", err)
			return menu,err
		}
	}
	return menu,nil
}
func (md *menuDaoImpl) Remove(id string) error{
	err := md.d.C("Menu").RemoveId(bson.ObjectIdHex(id))
	if(err != nil){
		md.logger.Println("Remove MenuId Error.", id, err)
		return err
	}
	return nil
}
