package router

type UeditorConfig struct {
	ImageActionName string `json:"imageActionName"`
	ImagePathFormat string `json:"imagePathFormat"`
	ImageFieldName string `json:"imageFieldName"`
	ImageUrlPrefix string `json:"imageUrlPrefix"`
	ImageCompressEnable bool `json:"imageCompressEnable"`
	ImageCompressBorder int `json:"imageCompressBorder"`
	ImageMaxSize int `json:"imageMaxSize"`
	ImageInsertAlign string `json: "imageInsertAlign"`
	ImageAllowFiles []string `json:"imageAllowFiles"`
	ScrawlActionName string `json:"scrawlActionName"`
	ScrawlFieldName string `json:"scrawlFieldName"`
	ScrawlPathFormat string `json:"scrawlPathFormat"`
	ScrawlMaxSize int `json:"scrawlMaxSize"`
	ScrawlUrlPrefix string `json:"scrawlUrlPrefix"`
}

var ueditorConfig UeditorConfig = UeditorConfig{}
func init() {
	ueditorConfig.ImageUrlPrefix,ueditorConfig.ScrawlUrlPrefix = "",""
	ueditorConfig.ImageActionName,ueditorConfig.ScrawlActionName = "img.upload","draw.upload"
	ueditorConfig.ImageFieldName,ueditorConfig.ScrawlFieldName = "file","file"
	ueditorConfig.ImagePathFormat,ueditorConfig.ScrawlPathFormat = "/ueditor/upload/{filename}","/ueditor/upload/{filename}"
	ueditorConfig.ImageCompressEnable = true
	ueditorConfig.ImageCompressBorder = 1600
	ueditorConfig.ImageMaxSize,ueditorConfig.ScrawlMaxSize = 5 * 1024 * 1024, 1 * 1024 * 1024
	ueditorConfig.ImageAllowFiles = []string{".jpg",".png",".jpeg",".gif"}
	ueditorConfig.ImageInsertAlign = "none"
}

func GetDefaultUeditorConfig()  *UeditorConfig{
	return &ueditorConfig
}