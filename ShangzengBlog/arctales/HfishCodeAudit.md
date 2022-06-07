Time: 20210921:golang: 首发 「 SecIN」 这里做个记录：...:HFish初版审计学习
--------


# 背景


* 首发 https://sec-in.com/article/949

在之前学习golang的靶场之后，下来开始尝试一些真实环境的代码审计工作，于是我找到了HFsih并下载了最初的版本。

到这里可能有师傅问了：为什么不直接梭哈最新版本？当然是因为最初版本比较简单并且适合新人（并不是担心高版本找不到问题会很尴尬），哈哈

话不多说，我们开干

# 下载/安装

Hfish是使用golang编写的一款开源蜜罐，其中下载地址如下: https://github.com/hacklcx/HFish/releases?after=0.4

在安装golang环境后，我们还需要安装gin框架:

```goalng
go get -u github.com/gin-gonic/gin
```
设置go mode

```golang
go env -w GOBIN=$HOME/bin
go env -w GO111MODULE=on
go env -w GOPROXY=https://goproxy.cn,direct
```

接下来进入目录就可以运行了

```goalng
go run main.go run
```


# 代码审计/学习

首先看入口 **main.go** 中 **setting.Run()** 为运行入口，进入进入**HFish/utils/setting/setting.go**, 我们可以看到使用**conf.Get("xxx", "xxx")**的方式进行读取**config.ini**中的配置文件，以此判断蜜罐是否启动，如果开启使用**Start（）**进行运行

```golang
// 启动 Redis 钓鱼
redisStatus := conf.Get("redis", "status")

// 判断 Redis 钓鱼 是否开启
if redisStatus == "1" {
	redisAddr := conf.Get("redis", "addr")
	go redis.Start(redisAddr)
```

这里后面再来看蜜罐，先来看看管理端是是怎么设置的，读取**config.txt** 中设置的管理员地址，并在 **http.Server**中进行设置:

```golang
	// 启动 admin 管理后台
	adminbAddr := conf.Get("admin", "addr")
	serverAdmin := &http.Server{
		Addr:         adminbAddr,
		Handler:      RunAdmin(),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	serverAdmin.ListenAndServe()
}
```
看`RunAdmin() `中设置的路由，这里使用的gin框架，生成日志并写入静态资源

```goalng
func RunAdmin() http.Handler {
	gin.DisableConsoleColor()
	f, _ := os.Create("./logs/hfish.log")
	gin.DefaultWriter = io.MultiWriter(f)
	// 引入gin
	r := gin.Default()

	r.Use(gin.Recovery())
	// 引入html资源
	r.LoadHTMLGlob("admin/*")

	// 引入静态资源
	r.Static("/static", "./static")

	// 加载路由
	view.LoadUrl(r)

	return r
}
```

跟进**view.LoadUrl(r)**，查看加载的路由功能 ,进入**HFish/view/url.go**，访问**http://127.0.0.1/login**黑白结合进行查看，首先是登陆页面：直接夹在静态资源:


```goalng

func Html(c *gin.Context) {
	data := getSetting() //订阅通知等
	c.HTML(http.StatusOK, "setting.html", gin.H{
		"dataList": data,
	})
}
```

尝试登陆，这里是直接post提交 loginName和 loginPwd 为账号密码，我们在HFish/view/login/view.go中查看登陆逻辑 — 获取用户的登陆账号和密码，与config.ini重的账号密码进行比对，如果一直就在cookie中写入登陆的用户名:

```golang
func Login(c *gin.Context) {
	loginName := c.PostForm("loginName")
	loginPwd := c.PostForm("loginPwd")

	account := conf.Get("admin", "account")
	password := conf.Get("admin", "password")

	if loginName == account {
		if loginPwd == password {
			c.SetCookie("is_login", loginName, 60*60*24, "/", "*", false, true)
			c.JSON(http.StatusOK, error.ErrSuccessNull())
			return
		}
	}

	c.JSON(http.StatusOK, error.ErrLoginFail())
}

```

## 无验证码-容易爆破

其实这也不算是什么问题，但是这玩意万一有用呢？抱着这种想法，我写了脚本，结果居然大约爆破出来二十多个公网的蜜罐管理员账号密码，脚本如下：

* https://github.com/shangzeng/GolangTools/tree/master/12.HFishPassScan

## 绕过密码登陆管理员

继续查看学习管理员的路由设置，发现在登出操作中，就是清除cookie中的is_login 的用户名，从而达到目的


```golang
func Logout(c *gin.Context) {
	c.SetCookie("is_login", "", -1, "/", "*", false, true)
	c.Redirect(http.StatusFound, "/login")
}
```
而FHish怎么判断用户是否登陆的呢，这里使用login.Jump函数进行判断，HFish/view/login/view.go 的jump函数进行查看，这里出现了问题： 只是判断了cookie中的用户名为管理员用户名，就判断用户是已经登陆的状态，进入了管理员界面，根本没有用到密码:


```golang
func Jump(c *gin.Context) {
	account := conf.Get("admin", "account")
	loginCookie, _ := c.Cookie("is_login")
	if account != loginCookie {
		c.Redirect(http.StatusFound, "/login")
		c.Abort()
		return
	} else {
		c.Next()
	}
}
```
那么问题来了：如果我们知道管理员账号（一般都是admin），那么我们只要修改cookie中的内容，就可以进行登陆了。测试一下：当不存在的时候，登陆mail界面会跳转到登陆页面

![](https://sec-in.com/img/sin/M00/00/62/wKg0C2BLDgOAPP-IAAAlTBWatwM856.png)

但是is_login=admin时，就可以通过login.Jump函数验证，访问管理员界面

![](https://sec-in.com/img/sin/M00/00/62/wKg0C2BLDhKAOrcNAAAvdIxCAL0420.png)


也就是说：只要我们知道了管理员账号名，我们可以越过密码进行登录（划重点）。


## 前台存储XSS

继续查看，仪表盘功能和上钩列表功能主要是展示，没有什么输入操作，接下来查看邮件群发功能：接收邮件，并且在send.SendMail 中使用golang的gomail（https://gopkg.in/gomail.v2）进行发送邮件，使用的是sqlite数据库，但是这里只是用来一个查询并不可控。

![](https://sec-in.com/img/sin/M00/00/62/wKg0C2BLDkeABFNwAABLxhvHZZQ713.png)

接着查看配置功能,其中view.go中的GetSettingInfo 接收了ID参数用于查询邮件配置,但是这里使用占位符，不存在注入问题，

```golang
/*发送邮件*/
func SendEmailToUsers(c *gin.Context) {
	emails := c.PostForm("emails")
	title := c.PostForm("title")
	from := c.PostForm("from")
	content := c.PostForm("content")

	eArr := strings.Split(emails, ",")
	sql := `select status,info from hfish_setting where type = "mail"`
	isAlertStatus := dbUtil.Query(sql)
	info := isAlertStatus[0]["info"]
	config := strings.Split(info.(string), "&&")

	if from != "" {
		config[2] = from
	}

	send.SendMail(eArr, title, content, config)
	c.JSON(http.StatusOK, error.ErrSuccessNull())
}
```

接下来就是看API接口，主要的目的就是上报we蜜罐的信息,默认开启（这个重要）,我们可以看到上报的信息直接写入了sqlite数据库中，虽然用占位符不存在注入问题了，但是是是否会存在XSS呢？


```golang
// 获取钓鱼信息
func GetFishInfo(c *gin.Context) {
	id, _ := c.GetQuery("id")
	sql := `select info from hfish_info where id=?;`
	result := dbUtil.Query(sql, id)
	c.JSON(http.StatusOK, error.ErrSuccess(result))
}
```

我们看看蜜罐是怎么上传数据的 HFish/web/github/static/github.js 文件，sec_key直接写在js里面的可控:


```golang
function report() {
    var login_field = $("#login_field").val();
    var password = $("#password").val();

    $.ajax({
        type: "POST",
        url: "http://localhost:9001/api/v1/post/report",
        dataType: "json",
        data: {
            "name": "Github钓鱼",
            "info": login_field + "&&" + password,
            "sec_key": "9cbf8a4dcb8e30682b927f352d6559a0"
        },
        success: function (e) {
            if (e.code == "200") {
                window.location.href = "https://github.com";
            } else {
                console.log(e.msg)
            }
        },
        error: function (e) {
            console.log("fail")
        }
    });
}
```

查看fish. GetFishList而管理员界面也是直接读取数据，没有编码,很有可能存在XSS,使用钓鱼接口发送XSS测试

![](https://sec-in.com/img/sin/M00/00/62/wKg0C2BLDwuAPuGsAABNvVZvIIw979.png)

查看管理员界面，触发XSS http://127.0.0.1/fish

![](https://sec-in.com/img/sin/M00/00/62/wKg0C2BLDxuAddAtAAATlEEZEhk725.png)


# 最后

到这里代码就大致分析完毕了，还是有一些收获的哈哈（虽然是没人说的Nday）接下如果有时间可以继续学习蜜罐的编写或者查看下一个版本的HFih进行学习。最后感谢观看哈～

