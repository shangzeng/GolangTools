package main

import (
	_ "ShangzengBlog/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

func main() {

	// 输出日志
	logs.Async()
	logs.SetLogger(logs.AdapterMultiFile, `{"filename": "log.log", "daily": true}`)

	// 方法函数
	beego.AddFuncMap("sub", sub)
	beego.AddFuncMap("add", add)
	beego.Run()
}

/*
后续增加：
 1. 文章目录
 2. 将静态资源打包进入二进制文件中
 3. 增加一个SRC漏洞挖掘页面


目前改正：
  1. 修复了标签问题
  2. DNS log 功能： 暂定采用 ceye 来进行检测
  3. 增加浏览量 (仅总量)
*/

func sub(in int) (out int) {
	out = in - 1
	return
}

func add(in int) (out int) {
	out = in + 1
	return
}