package router

import (
	"github.com/go-martini/martini"
	"fmt"
	dbw "../db"
	"github.com/martini-contrib/sessionauth"
	"github.com/martini-contrib/binding"
	"net/http"
	"io"
	"os"
	"log"
	"path/filepath"
	"github.com/satori/go.uuid"

	"gopkg.in/mgo.v2/bson"
	"github.com/martini-contrib/render"
)

func init() {
	fmt.Println("加载 BannerRouter")
	routers.PushBack(BannerRouter)
}

func BannerRouter(m *martini.ClassicMartini)  {
	m.Get("/admin/banner",sessionauth.LoginRequired, func(r LayoutWrapper,bannerDao dbw.BannerDao){
		banners,_ := bannerDao.FindAll()
		model := map[string]interface{}{"banners": banners}
		r.HTML(200,"admin/banner", model,"admin")
	})
	m.Post("/admin/banner", sessionauth.LoginRequired, binding.Bind(dbw.Banner{}), func(r LayoutWrapper, logger *log.Logger, req *http.Request, bannerDao dbw.BannerDao, banner dbw.Banner) {
		file,fileHeader,ferr := req.FormFile("ImgFile")
		banners,_ := bannerDao.FindAll()
		model := map[string]interface{}{"banners": banners}
		newFile := false
		defer file.Close()
		method,mok := req.URL.Query()["method"]
		id,iok := req.URL.Query()["id"]
		if(mok && method[0] == "add"){
			banner.Id = ""
		}else if(mok && method[0] == "upload" && iok){
			banner.Id = bson.ObjectIdHex(id[0])
		}else{
			model["error"] = "无效的请求"
			r.HTML(200,"admin/banner", model,"admin")
			return
		}
		if(ferr == nil){
			targetPath := "/banner/" + uuid.NewV4().String() +  filepath.Ext(fileHeader.Filename)
			os.MkdirAll(uploadPath + "/banner",0660)
			target,err := os.Create(uploadPath + targetPath)
			if(err != nil){
				logger.Println("upload error,create file failed", err, targetPath)
				model["error"] = "创建文件失败"
				r.HTML(200,"admin/banner", model,"admin")
				return
			}
			defer target.Close()
			_,err = io.Copy(target, file)
			if(err != nil){
				logger.Println("upload error,copy file failed", err, targetPath)
				model["error"] = "复制文件失败"
				r.HTML(200,"admin/banner", model,"admin")
				return
			}
			banner.ImgPath = targetPath
			newFile = true
		}
		err := bannerDao.Save(banner,newFile)
		if(err != nil){
			model["error"] = "保存数据失败"
		}
		banners,_ = bannerDao.FindAll()
		model["banners"] = banners
		r.HTML(200,"admin/banner", model,"admin")
	})
	m.Get("/admin/banner/:id/remove", func(params martini.Params,bannerDao dbw.BannerDao,r render.Render) {
		id,ok := params["id"]
		model := map[string]interface{}{}
		if(ok){
			bannerDao.RemoveId(bson.ObjectIdHex(id))
			banners, _ := bannerDao.FindAll()
			model["banners"] = banners
		}
		r.Redirect("/admin/banner")
	})
}
