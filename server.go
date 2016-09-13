package main

import (
	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/martini-contrib/sessions"
	"github.com/martini-contrib/sessionauth"
	r "./router"
	dbw "./db"
	//"container/list"
	"container/list"
	"html/template"
	"fmt"
	"strings"
	"regexp"
)

/* 模板自定义方法 start */
// 按次数循环
func templateLoop(n int,start int) []int {
	ret := make([]int,n)
	for i := 0; i < n; i++{
		ret[i] = i + start;
	}
	fmt.Println(ret)
	return ret
}
// int减法
func templateSub(i,j int) int{
	return i - j;
}
// int加法
func templateAdd(i,j int) int{
	return i + j;
}
// 是否为空
func templateIsNil(p interface{}) bool{
	return p == nil
}
// HTML处理
func templateHTML(html string) template.HTML{
	return template.HTML(html)
}
func templateToLower(str string) string{
	return strings.ToLower(str)
}
func templateHtml2Txt(html string) string {
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	text := re.ReplaceAllString(html , " ")
	return text
}
func templateLimitText(text string, limit int,suffix string) string{
	r := []rune(text)
	if(len(r) > limit){
		fmt.Println(r[0:limit])
		return string(r[0:limit]) + suffix
	}else{
		return text
	}
}
/* end 自定义方法 */

func main() {
	var m = martini.Classic()
	m.Use(martini.Static("assets"))
	store := sessions.NewCookieStore([]byte("secret123"))
	m.Use(sessions.Sessions("my_session", store))
	m.Use(sessionauth.SessionUser(func () sessionauth.User{
		return &dbw.LoginUser{}
	}))
	m.Use(render.Renderer(render.Options{
		Directory: "templates", // Specify what path to load the templates from.
		Layout: "layout/default/layout", // Specify a layout template. Layouts can call {{ yield }} to render the current template.
		Extensions: []string{".tmpl", ".html"}, // Specify extensions to load for templates.
		//Delims: render.Delims{"{[{", "}]}"}, // Sets delimiters to the specified strings.
		Charset: "UTF-8", // Sets encoding for json and html content-types. Default is "UTF-8".
		Funcs:[]template.FuncMap{{"loop":templateLoop,"sub":templateSub,"add":templateAdd,"toLower":templateToLower,
			"split":strings.Split,"isNil":templateIsNil,"HTML":templateHTML, "MainMenu":func()string{return ""},
			"Html2Txt":templateHtml2Txt,"textLimit":templateLimitText}},
//		IndentJSON: true, // Output human readable JSON
//		IndentXML: true, // Output human readable XML
	}))

	// mongo数据库
	m.Use(dbw.Mongoer("localhost","cms"))

	var routers *list.List = r.GetRouters();
	// 初始化router
	for rl := routers.Front(); rl != nil; rl = rl.Next() {
		fun,ok := rl.Value.(func(m *martini.ClassicMartini))
		if(ok){
			fun(m)
		}
	}
	// 初始化dao
	for el := dbw.DaoList.Front(); el != nil; el = el.Next() {
		m.Use(el.Value)
	}
	// 封装简化HTML模板
	m.Use(r.Wapper)
	m.Run()
}
