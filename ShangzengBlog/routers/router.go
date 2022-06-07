package routers

import (
	"ShangzengBlog/controllers"
	"github.com/astaxie/beego"
	"html/template"
	"net/http"
)

func init() {
	// 实现了简单的首页功能
	beego.Router("/", &controllers.MainController{},"get:Index")
	beego.Router("page/about.html", &controllers.MainController{},"get:AboutMe")
	beego.Router("page/links.html", &controllers.MainController{},"get:Friends")
	// 实现了文章分类展示
	beego.Router("/categories/:link([\\w]+).html", &controllers.MainController{}, "get:Categories")
	// 实现动态文章的页面 article
	beego.Router("/article/:slug([\\w\\-]+).html", &controllers.MainController{},"get:ArticleInfo")
	// DNSlog 暂时不用了
	//beego.Router("/dnslog", &controllers.DNSController{},"get:Index")
	// 404
	beego.ErrorHandler("404", HttpNotFound)
	// rss 订阅内容生成
	beego.Router("/feed/", &controllers.FeedController{},"get:Index")
}

func HttpNotFound(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("views/404.html")
	if err != nil {
		panic(err)
	}
	err = t.Execute(w, "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}