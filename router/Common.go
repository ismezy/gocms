package router

import (
	"container/list"
	"fmt"
	"github.com/go-martini/martini"
	"net/http"
	"os"
	"io"
	"path/filepath"
)

var routers *list.List = list.New()
var uploadPath = "d:\\temp\\"

func init() {
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
		ext := filepath.Ext(file)
		res.Header().Add("Content-Type","image/" + ext[1:len(ext)])
//		res.Header().Add("attachment;filename", filepath.])
		fmt.Println("download :%s",path)
		f, err := os.Open(file)
		defer f.Close()
		if(err == nil){
			io.Copy(res,f)
		}
	})
}
