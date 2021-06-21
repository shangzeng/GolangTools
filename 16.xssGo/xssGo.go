package main

import(
	"github.com/gin-gonic/gin"
	"xssGo/view"
	"net/http"
	"time"
	_"fmt"
	"io"
	"os"
)



func main() {

	// 启动web页面
	server := &http.Server{
		Addr:         "0.0.0.0:8099",        // 地址
		Handler:      RunAdmin(),            // 路由
		ReadTimeout:  5 * time.Second,       // 超时
		WriteTimeout: 10 * time.Second,
	}

	server.ListenAndServe()                  // 监听开启服务 
}





// 路由
func RunAdmin() http.Handler {
	gin.DisableConsoleColor()                // 禁止控制台颜色
	f, _ := os.Create("./logs/hfish.log")    // 创建日志
	gin.DefaultWriter = io.MultiWriter(f)
	// 引入gin  创建路由
	r := gin.Default()
	// 使用 Recovery 中间件
	r.Use(gin.Recovery())
	// 引入html资源 定义html模板路径
	r.LoadHTMLGlob("admin/*")
	// 引入静态资源 定义静态资源模板路径
	r.Static("/static", "./static")
	// 加载路由 view 是个库 = = 
	view.LoadUrl(r)
	return r
}