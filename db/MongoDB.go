package db

import (
	"github.com/go-martini/martini"
	"gopkg.in/mgo.v2"
	"log"
	"container/list"
)

type Mongo interface {
	open() (*mgo.Session,error)
	close()
}

type mongo struct {
	connectString string
	db string
	session *mgo.Session
}

var DaoList = list.New()

func Mongoer(connectString string, db string) martini.Handler{
	mongo := mongo{connectString:connectString,db:db}
	return func(c martini.Context, logger *log.Logger) {
		session,err := mongo.open()
		defer session.Close()
		if(err != nil){
			logger.Fatalln("mongo数据库连接失败",err)
			return
		};
		d := session.DB(mongo.db)
		c.Map(d)
		c.Next()
	}
}
/**
 * 连接数据库
 */
func (m *mongo) open() (*mgo.Session,error){
	if(m.session == nil){
		session, err := mgo.Dial(m.connectString);
		if(err != nil){
			return nil,err
		}
		m.session = session
	}
	m.session.SetMode(mgo.Monotonic, true)
	return m.session.Clone(),nil;
}

func (m *mongo) close(){
	m.session.Close()
}

