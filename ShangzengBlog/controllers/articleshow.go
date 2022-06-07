package controllers

import (
	"bufio"
	"github.com/astaxie/beego"
	"github.com/russross/blackfriday"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"fmt"
)



type MainController struct {
	beego.Controller
}

// 自我介绍页面 AboutMe
func (c *MainController) AboutMe() {
	AddViews()
	fileread, _ := ioutil.ReadFile("arctales/about.html")
	unsafe := blackfriday.MarkdownCommon(fileread)
	c.Data["Content"] = string(unsafe)
	c.Data["Number"] = lognumber()
	c.TplName = "page.html"
}

// 友链 Friends
func (c *MainController) Friends() {
	AddViews()
	fileread, _ := ioutil.ReadFile("arctales/friends.html")
	unsafe := blackfriday.MarkdownCommon(fileread)
	c.Data["Content"] = string(unsafe)
	c.Data["Number"] = lognumber()
	c.TplName = "page.html"
}

// 写类型，用于存储文章名 文章时间 文章标签
type ArctlsInfo struct {
	Name 		string				// 文件名字
	ArcName     string              // 文章标题
	Time 		string				// 文章编写时间
	Tag  		string				// 文章标签
	Information string				// 文章介绍信息
}

// 获取信息
// 在这里最好还要根据时间把文章进行排名
func GetArctlsInfo() []ArctlsInfo {
	var ArctlsInfos []ArctlsInfo
	for _, filename := range findarctlename() {
		fileread, _ := ioutil.ReadFile("arctales/"+filename+".md")
		title := strings.Split(string(fileread), "--------")[0]
		titleline := strings.Split(title, ":")
		ArctlsInfos= append(ArctlsInfos,ArctlsInfo{Name:filename,Time:titleline[1],Tag:titleline[2],Information:titleline[3],ArcName:titleline[4]})
		fmt.Println(ArctlsInfos)
	}
	// 文章按照时间顺序进行排序
	sort.SliceStable(ArctlsInfos, func(i, j int) bool {
		return ArctlsInfos[i].Time > ArctlsInfos[j].Time
	})
	return ArctlsInfos
}

//尝试显示首页的展示
func (c *MainController) Index() {
	var (
		page         int      // 页面设置
		//number       int      // 浏览量设置
	)
	// GetInt 是beego中的函数，在默认无的情况下为 1 ， 这里主要用于限制页面刷领
	if page, _ = c.GetInt("page"); page < 1 {
		page = 1
	}

	AddViews()
	c.Data["CurrentPage"] = page
	c.Data["Len"] = len(GetArctlsInfo())
	c.Data["Title"] = "Shangzeng's blog"
	c.Data["Number"] = lognumber()
	c.Data["Data"] = GetArctlsInfo()
	c.TplName = "index.html"
}





// 拆解成为两种： 一个函数记录浏览量并写入 配置文件  一个函数读取
// 查看日志浏览量
func lognumber() int {
	file, _ := os.Open("view")
	fileScanner := bufio.NewScanner(file)
	var lines []string
	for fileScanner.Scan() {
		lines = append(lines, fileScanner.Text())
	}
	intVar, _ := strconv.Atoi(lines[0])
	//保存浏览量并输出
	return intVar
}
// 增加日志量
func AddViews() {
	add := lognumber() + 1
	tracefile(strconv.Itoa(add))
}
func tracefile(str_content string)  {
	fd,_:=os.OpenFile("view",os.O_WRONLY|os.O_TRUNC|os.O_APPEND,0644)
	fd_content:=str_content
	buf:=[]byte(fd_content)
	fd.Write(buf)
	fd.Close()
}

// 遍历获取 arctales 文件夹中以 md 结尾的文件名 ，准备输出到首页上
func findarctlename() (files []string){
	file, err := ioutil.ReadDir("arctales")
	if err != nil {
		beego.Info(err)
	}
	for _, file := range file {
		//判断是否是md 文件
		if strings.Index(file.Name(), ".md") > -1 {
			//fmt.Println(strings.Trim(file.Name(),".md"))
			files = append(files,strings.Trim(file.Name(),".md"))
		}
	}
	return files
}

// 寻找文件是否正确
func findarctleconnect(name string) (bool) {
	file, err := ioutil.ReadDir("arctales")
	if err != nil {
		beego.Info(err)
	}
	for _, file := range file {
		if file.Name() == name {
			return true
		}
	}
	return false
}

// 并返回 md 内容
func readarctleconnect(name string) (file string) {
	fileread, _ := ioutil.ReadFile("arctales/"+name)
	title := strings.Split(string(fileread), "--------")[1]
	unsafe := blackfriday.MarkdownCommon([]byte(title))
	// 内容处理下 ， "---" 之前的不要了
	//fmt.Println(string(unsafe))
	return string(unsafe)
}

// 实现动态解析分析 ArticleInfo
func (c *MainController) ArticleInfo() {
	AddViews()
	slug := c.Ctx.Input.Param(":slug")+".md"
	// 根据 slug 读取 相关的 md 文件代码， 如果不存在返回 404
	// 先判断文件是是否存在
	if findarctleconnect(slug) {
		c.Data["Content"] = readarctleconnect(slug)
		c.Data["Title"] = check(slug).ArcName
		c.Data["Time"] = check(slug).Time
		c.Data["Tag"] = check(slug).Tag
		c.Data["Information"] = check(slug).Information
		c.Data["Number"] = lognumber()
		c.TplName = "article.html"
		return
	}
	c.TplName = "404.html"
}

// 循环判断GetArctlsInfo() 中内容返回合适的
func check(filename string ) ArctlsInfo {
	for _, name := range GetArctlsInfo() {
		if name.Name+".md" == filename {
			return name
		}
	}
	return ArctlsInfo{}
}

// 文章分类展示 Categories
func (c *MainController) Categories() {
	AddViews()
	tag := c.Ctx.Input.Param(":link")
	// 搜索包含关键词的所有文章
	//frintln(checkTag(tag))
	c.Data["Tag"] = tag
	c.Data["Data"] = checkTag(tag)
	c.Data["Number"] = lognumber()
	c.TplName = "category.html"
}

// 循环判断GetArctlsInfo() 中内容返回合适的 tag
func checkTag(tag string ) []ArctlsInfo {
	var TagArctlsInfo []ArctlsInfo
	for _, name := range GetArctlsInfo() {
		if name.Tag == tag {
			//return name
			TagArctlsInfo = append(TagArctlsInfo ,name)
		}
	}
	return TagArctlsInfo
}

