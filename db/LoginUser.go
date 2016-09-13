package db

import (
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/mgo.v2"
	"log"
	"fmt"
	"github.com/go-martini/martini"
	"errors"
)

// 登录用户缓存
var userCache map[string] *LoginUser = map[string] *LoginUser {}

type LoginUser struct {
	Id 				bson.ObjectId	`bson:"_id"`
	Name 			string			`bson:"Name"`
	LoginId			string			`bson:"LoginId"`
	Password		string			`bson:"Password"`
	Authenticated	bool
}

func (u *LoginUser) IsAuthenticated() bool{
	return u.Authenticated
}

func (u *LoginUser) Login() {
	u.Authenticated = true
	userCache[u.Id.Hex()] = u
	fmt.Println("login",u)
}

func (u *LoginUser) Logout() {
	u.Authenticated = false
	delete(userCache, u.Id.Hex())
}

func (u *LoginUser) UniqueId() interface{} {
	return u.Id.Hex()
}

func (u *LoginUser) GetById(id interface{}) error {
	userId := id.(string)
	user,ok := userCache[userId]
	if(!ok){
		return errors.New("user not login")
	}

	u.Id = user.Id
	u.Authenticated = user.Authenticated
	u.LoginId = user.LoginId
	u.Name = user.Name
	return nil
}

type LoginUserDao interface {
	FindByLoginId(loginId string) (LoginUser,error)
}
type loginUserDaoImpl struct {
	d *mgo.Database
	logger *log.Logger
}
func init() {
	fmt.Println("加载 LoginUserDao")
	DaoList.PushBack(func (d *mgo.Database,c martini.Context, logger *log.Logger){
		c.MapTo(&loginUserDaoImpl{d:d,logger:logger},(*LoginUserDao)(nil))
	})
}
func (ld *loginUserDaoImpl) FindByLoginId(loginId string) (LoginUser,error){
	var users [] LoginUser
	err := ld.d.C("LoginUser").Find(bson.M{"LoginId":loginId}).All(&users)
	if err != nil{
		ld.logger.Println("find user by login id error", err, loginId)
		return LoginUser{},err
	}
	if len(users) < 1 {
		ld.logger.Println("user not found", loginId)
		return LoginUser{},errors.New("user not found")
	}
	return users[0],nil
}