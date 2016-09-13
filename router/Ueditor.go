package router

import (
	"fmt"
	"github.com/go-martini/martini"
//	"gopkg.in/mgo.v2"
//	"log"
	"net/http"
	"github.com/martini-contrib/render"
	"log"
	"os"
	"io"
	"github.com/satori/go.uuid"
	"strings"
	"path/filepath"
	"encoding/base64"
	"github.com/martini-contrib/sessionauth"
)

func init() {
	fmt.Println("加载 ueditorRouter")
	routers.PushFront(ueditorRouter)
}
func ueditorRouter(m *martini.ClassicMartini){
	m.Get("/admin/ueditor",sessionauth.LoginRequired, func(r render.Render, req *http.Request, res http.ResponseWriter,params martini.Params, logger *log.Logger) {
		fmt.Println(req.URL.Query())
		action := req.URL.Query()["action"][0]
		switch action{
		case "config":		// 获取配置
			r.JSON(200,GetDefaultUeditorConfig())
			break
		default:
			r.JSON(200,map[string]interface{}{"msg":"no action"})
		}
	})
	m.Post("/admin/ueditor",sessionauth.LoginRequired,func(r render.Render, req *http.Request, res http.ResponseWriter,params martini.Params, logger *log.Logger){
		action := req.URL.Query()["action"][0]
		fmt.Println("action:",action)
		switch action{
		case "draw.upload":	// 涂鸦
			uploadBase64(req,r,logger)
			break
		case "img.upload":	// 图片上传
			uploadImg(req,r,logger)
			break;
		default:
			r.JSON(200,map[string]interface{}{"msg":"no action"})
		}
	})
	var imgs = "jpg,png,gif"
	m.Get("/ueditor/upload/:filename\\.:ext", func(params martini.Params, res http.ResponseWriter) {
		if(strings.Contains(imgs,params["ext"])){
			res.Header().Add("Content-Type","image/" + params["ext"])
		}
		path := uploadPath + params["filename"] + "." + params["ext"]
		fmt.Println("download :%s",path)
		file, err := os.Open(path)
		defer file.Close()
		if(err == nil){
			io.Copy(res,file)
		}
	})
}
func uploadBase64(req *http.Request, r render.Render, logger *log.Logger) {
//	req.ParseForm()
	file := req.FormValue("file")
	if(file != ""){
		bytes,err := base64.StdEncoding.DecodeString(file)
		if(err != nil){
			r.Text(405,"base64 decode error")
			logger.Fatal("ueditor draw upload error, base64 decode error", err)
			return
		}

		newFileName := uuid.NewV4().String() + ".jpg" // 在后台生成唯一文件
		writer,err:= os.Create(uploadPath + newFileName);
		if(err != nil){
			r.Text(405,"create file error")
			logger.Fatal("ueditor draw upload error, create file error", err)
			return
		}
		defer  writer.Close()
		len,_ := writer.Write(bytes)
		fmt.Println(newFileName)
		r.JSON(200,map[string]interface{}{"state":"SUCCESS","url":"/ueditor/upload/" + newFileName,"title":newFileName,"original":newFileName,"type":"jpg","size": len})
	}else {
		r.Text(405,"file not found")
	}
}

// 文件上传处理
func uploadImg(req *http.Request, r render.Render, logger *log.Logger){
	err := req.ParseMultipartForm(10 * 1024 * 1024)			// 10MB内存
	if(err != nil){
		r.Text(405,"parse body error")
	}
	files,ok:= req.MultipartForm.File["file"]
	if(!ok){
		r.Text(405,"not file field")
		return
	}
	if(len(files) < 1){
		r.Text(405,"file is empty")
		return
	}
	file := files[0]
	ext := filepath.Ext(file.Filename)
	reader,_:= file.Open()
	defer  reader.Close()
	newFileName := uuid.NewV4().String() + ext		// 在后台生成唯一文件
	writer,_:= os.Create(uploadPath + newFileName)
	defer writer.Close()
	size,_ := io.Copy(writer,reader)
	r.JSON(200,map[string]interface{}{"state":"SUCCESS","url":"/ueditor/upload/" + newFileName,"title":newFileName,"original":file.Filename,"type":strings.TrimLeft(ext,"."),"size":size})
}