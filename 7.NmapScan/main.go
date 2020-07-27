package main

import(

	"fmt"
	"sync"
	_"time"
	"NmapScan/functions"
	"NmapScan/plugins"
	//pb "gopkg.in/cheggaaa/pb.v1"


	// 弱口令库调用
	_ "NmapScan/weakpassword/sshweakpassword"
	_ "NmapScan/weakpassword/ftpweakpassword"
	_ "NmapScan/weakpassword/redisweakpassword"


	// 主机漏洞调用
	_ "NmapScan/vulnerability/ms17010"

	
)



func init() {

	fmt.Println(functions.Banner)
	// 加载配置文件内容 读取json文件，准备扫描
	functions.IplistFile, functions.Nmapports , functions.Nmapresults ,functions.Nmapthreads ,functions.Nmaptimeout, functions.Nmapresultfile = functions.Configing("config.json")
	// 读取iplist
	functions.NmapIplist, functions.NmapIpnumber, _=  functions.ReadLines(functions.IplistFile)
	// 读取用户名/密码字典
	functions.Usernamelist = functions.ReadUserDict("iplist/username.txt")
	functions.Passwordlist = functions.ReadPasswordDict("iplist/password.txt")

	// 生成excel表格等待结果写入
	//functions.ExcleCreate(functions.Nmapresultfile)

}



func main() {
	//使用nmap进行扫描
	fmt.Printf("需要nmap扫描数量：  %d  \n",functions.NmapIpnumber)
	informations := make(chan functions.NmapRequests, functions.Nmapthreads)
	var wg sync.WaitGroup

	for i := 0; i < cap(informations); i++ {
		go functions.Worker(informations, &wg)
	}

	for _,v := range functions.NmapIplist {
		wg.Add(1)
		informations <-  functions.NmapRequests{Host: v,Ports: functions.Nmapports}
	}
	wg.Wait()
	close(informations)






	//根据获得的数据进行扫描
	requestlist := functions.ReadIpList("iplist/nmapresult.txt")
	// 根据服务进行分类 【 增加扫描类型记得添加插件 ！】
	weakpasswordaliveIpList,httpliveIpList,vulnerabillityaliveIpList,otheraliveIpList := functions.Checkmodels(requestlist)
	// 显示详情
	fmt.Printf("弱口令服务如下：\n")
	fmt.Println(weakpasswordaliveIpList)
	fmt.Printf("HTTP服务如下：\n")
	fmt.Println(httpliveIpList)
	fmt.Printf("可能存在服务器漏洞的服务如下：\n")
	fmt.Println(vulnerabillityaliveIpList)
	fmt.Printf("其他服务如下：\n")
	fmt.Println(otheraliveIpList)


	// 进行分类 - 弱口令任务
	Task := functions.AddUserAndPass(weakpasswordaliveIpList,functions.Usernamelist,functions.Passwordlist)
	RunTaskWeakPassword(Task)
	//fmt.Println(Task)
	// 使用插件
	//fn := plugins.ScanFuncMap["SSH"]
	//fn()
}







// 并发，执行弱口令任务
func RunTaskWeakPassword(Task []functions.ServiceWeakPassword) {
	fmt.Printf("弱口令爆破数量：  %d  \n",len(Task))
	informations := make(chan functions.ServiceWeakPassword, 10)
	var wg sync.WaitGroup

	for i := 0; i < cap(informations); i++ {
		go WeakPasswordWorker(informations, &wg)
	}

	for _,v := range Task {
		wg.Add(1)
		informations <-  v
	}


	wg.Wait()
	close(informations)

}



func WeakPasswordWorker(informations chan functions.ServiceWeakPassword, wg *sync.WaitGroup) {
	for v := range informations {
		//fmt.Println(v.Server.Protocol)
		plugins.Start(v)
		wg.Done()
	}	
}











