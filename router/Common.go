package router

import (
	"container/list"
	"fmt"
	"github.com/go-martini/martini"
	"net/http"
	"os"
	"io"
	"path/filepath"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/render"
	"github.com/satori/go.uuid"
	"log"
)

var routers *list.List = list.New()
var uploadPath = "/temp"
type JsonMsg map[string]string

func init() {
	up := os.Getenv("cms_upload_path")
	if(up != ""){
		uploadPath = up
	}

	fmt.Println("加载 CommonRouter")
	routers.PushBack(CommonRouter)
}

func CommonRouter (m *martini.ClassicMartini) {
	m.Get("/common/image", func(req *http.Request, res http.ResponseWriter) {
		path,ok := req.URL.Query()["path"]
		if(!ok){
			return
		}
		var file = uploadPath + path[0]
		http.ServeFile(res,req,file)
	})
	m.Post("/admin/image/upload", sessionauth.LoginRequired,func(req *http.Request,r render.Render, logger *log.Logger){
		file,fileHeader,err := req.FormFile("file")
		if(err != nil){
			logger.Println("upload image failed,file not found",err)
			r.JSON(405,map[string]string{"error":"upload failed"})
			return
		}
		defer file.Close()
		mkerr := os.MkdirAll(uploadPath + "/temp",0660)
		if(mkerr != nil){
			logger.Println("mkdir error", uploadPath + "/temp",mkerr)
		}
		path := "/temp/" + uuid.NewV4().String() +  filepath.Ext(fileHeader.Filename)
		target,terr := os.Create(uploadPath + path)
		if(terr != nil){
			logger.Println("upload image failed,create file failed",uploadPath + path,terr)
			r.JSON(405,map[string]string{"error":"upload failed"})
			return
		}
		defer target.Close()
		io.Copy(target,file)
		r.JSON(200,map[string]string{"path":path})
	})
}

