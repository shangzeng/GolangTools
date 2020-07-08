package main

import(
	"ServerScan/models"   // 用于存放类型 变量 常量
	"ServerScan/function" // 加载常用方法 函数
	"fmt"
)

/*
待解决：

1. 密码验证成功就一次进行 【这个最好要可以设置，毕竟弱口令扫描的目的就是扫描出所有的弱口令】
2. 注册插件还需要修改识别协议的函数
3. 缺少进度条

*/



func init() {
	models.Iplist       = "iplist.txt"         // 扫描列表格式为   IP：端口｜服务
	models.Passwordlist = "password.txt"       // 弱口令爆破密码
	models.Usernamelist = "username.list"      // 弱口令爆破用户名
	models.Threads      = 10                   // 线程数量（数量过大，并发扫描结果就不准确，很难受）
	models.Httplever    = "1"                  // 设定http漏洞扫描等级 分为 1 2 3 
	//models.Timeout      = 3                  // 每个扫描模块的超时时间    [超时由于时间的类型暂时不会所以在models 是直接设置的]
}





func main() {
	// 装载用户名,密码
	// models.Service 装载全部数据 并判断端口存活
	users       := function.ReadUserDict(models.Usernamelist)
	passwords   := function.ReadPasswordDict(models.Passwordlist)
	requestlist := function.ReadIpList(models.Iplist)
	aliveIpList := function.CheckAlive(requestlist)
	//fmt.Println(aliveIpList)

	// 根据服务进行分类 【 增加扫描类型记得添加插件 ！】
	weakpasswordaliveIpList,httpliveIpList,vulnerabillityaliveIpList,otheraliveIpList := function.Checkmodels(aliveIpList)



	// 进行分类 - 弱口令任务
	function.RunTask(function.AddUserAndPass(weakpasswordaliveIpList,users,passwords))
	//fmt.Println(weakpasswordaliveIpList)



	// 进行分类 - 主机漏洞任务
	function.RunVulnerabillityTask(vulnerabillityaliveIpList)
	fmt.Println(vulnerabillityaliveIpList)



	// 进行分类 - HTTP 常见漏洞扫描任务
	// 获取HTTPheader信息等等
	fmt.Println(httpliveIpList)




	// 进行分类 - 其他端口服务任务
	fmt.Println(otheraliveIpList)






}