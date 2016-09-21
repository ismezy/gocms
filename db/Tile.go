package db

import (
	"gopkg.in/mgo.v2"
	"log"
	"fmt"
	"github.com/go-martini/martini"
)

type Tile struct {
	PluginName	string			`bson:"_id"`
	Items 		[]TileItem		`bson:"Item"`
}

type TileItem struct{
	Index		int			`bson:"Index" json:"Index"`
	Title		string		`bson:"Title" json:"Title"`
	Memo		string		`bson:"Memo" json:"Memo"`
	ImgPath		string		`bson:"ImgPath" json:"ImgPath"`
	Url			string		`bson:"Url" json:"Url"`
}
func init() {
	fmt.Println("加载 TileDao")
	DaoList.PushBack(func (d *mgo.Database,c martini.Context, logger *log.Logger){
		c.MapTo(&tileDao{d:d,logger:logger},(*TileDao)(nil))
	})
}
type TileDao interface {
	Upload(item TileItem) error
	Insert(item TileItem) (TileItem,error)
	Get() (Tile,error)
}

type tileDao struct {
	d *mgo.Database
	logger *log.Logger
}

func (td *tileDao) Upload(item TileItem) error {
	tile,err := td.Get()
	tile.Items[item.Index-1] = item
	err = td.d.C("Plugin").UpdateId("Tile",tile)
	return err
}
func (td *tileDao) Insert(item TileItem) (TileItem,error) {
	tile,err := td.Get()
	size := len(tile.Items)
	item.Index = size + 1
	tile.Items = append(tile.Items,item)
	err = td.d.C("Plugin").UpdateId("Tile",tile)
	if(err != nil){
		td.logger.Println("insert TileItem error",item, err)
	}
	return item,err
}
func (td *tileDao) Get() (Tile,error){
	tile := Tile{}
	err := td.d.C("Plugin").FindId("Tile").One(&tile)
	if(err != nil){
		tile = Tile{PluginName:"Tile",Items:[]TileItem{}}
		err = td.d.C("Plugin").Insert(tile)
	}
	return tile,err
}
