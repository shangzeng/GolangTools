package main

import(
	"nmapmasscan/nmap"
	"nmapmasscan/masscan"
	"nmapmasscan/functions"
	"github.com/gookit/color"
	"strings"
	"strconv"
	"context"
	_"reflect"
	"regexp"
	"fmt"
	"log"
	"os"
	"sync"
	"time"
	pb "gopkg.in/cheggaaa/pb.v1"
)


/*

-----------------------------------------------------------------------------------
																				  
注意: 主要用于安服人员用于内网资产收集（别忘记开白名单） 默认环境 masscan + nmap           
 
使用逻辑： nmap直接判断端口信息/masscan 扫描出端口结果 -> nmap判断端口服务 -> 生成数据

-----------------------------------------------------------------------------------

参考资料:  1. https://github.com/Ullaakut/nmap
		  2. https://github.com/dean2021/go-masscan



*/


var(

	banner   = `

                                                                       █                                    ██████ 
 ██████████   ██████    ██████  ██████  █████   ██████   ███████      ░█     ███████  ██████████   ██████  ░██░░░██
░░██░░██░░██ ░░░░░░██  ██░░░░  ██░░░░  ██░░░██ ░░░░░░██ ░░██░░░██  █████████░░██░░░██░░██░░██░░██ ░░░░░░██ ░██  ░██
 ░██ ░██ ░██  ███████ ░░█████ ░░█████ ░██  ░░   ███████  ░██  ░██ ░░░░░█░░░  ░██  ░██ ░██ ░██ ░██  ███████ ░██████ 
 ░██ ░██ ░██ ██░░░░██  ░░░░░██ ░░░░░██░██   ██ ██░░░░██  ░██  ░██     ░█     ░██  ░██ ░██ ░██ ░██ ██░░░░██ ░██░░░  
 ███ ░██ ░██░░████████ ██████  ██████ ░░█████ ░░████████ ███  ░██     ░      ███  ░██ ███ ░██ ░██░░████████░██     
░░░  ░░  ░░  ░░░░░░░░ ░░░░░░  ░░░░░░   ░░░░░   ░░░░░░░░ ░░░   ░░            ░░░   ░░ ░░░  ░░  ░░  ░░░░░░░░ ░░      


				                                                               version:1.0.0 @shangzeng
	`



	bar                *pb.ProgressBar
	iplist             []string
	justString         string
	masscanports       string
	masscanthreads     string
	masscaniplist      string
	nmapiplist         string
	masscantimeout     int
	number             int
	nmapthreads        int
	nmaptimeout        int

)

// 设置nmap请求路径
type NmapRequests struct {

	Host		   string
	Port           string
	Protocol       string

}



func init() {
	// 存放要扫描的IP清单
	iplist, number, _=  functions.ReadLines("iplist/iplist.txt")
	justString     = strings.Join(iplist,",")           // 扫描IP


	// 设置总体扫描模式
	// masscan + nmap 
	// nmap 

    // masscan 设置
	masscanports   = "21,22,23,25,53,80,81,110,111,123,135,137,139,161,389,443,445,465,500,515,520,523,548,623,636,873,902,1080,1099,1433,1521,1604,1645,1701,1883,1900,2049,2181,2375,2379,2425,3128,3306,3389,4730,5060,5222,5351,5353,5432,5555,5601,5672,5683,5900,5938,5984,6000,6379,7001,7077,8080,8081,8443,8545,8686,9000,9042,9092,9100,9200,9418,9999,11211,27017,37777,50000,50070,61616"       					    // 扫描端口
	masscanthreads = "200"    					        // 扫描rates
	masscantimeout = 5                                  // 设置masscan扫描超时
	masscaniplist  = "iplist/ipalive.txt"               // masscan扫描存活端口

	// nmap 设置
	nmapthreads    = 5                                  // 设置线程
	nmaptimeout    = 5                                  // 设置nmap扫描超时
	nmapiplist     = "iplist/nmapscan.txt"              // 存储 nmap 端口扫描结果


	// 输出扫描.xlsx 结果
	//results        = "iplist/result.xlsx"               // 输出最终结果文件名

}



func main() {
	// 准备
	fmt.Println(banner)
	// 扫描后数据写入 iplist/ipalive.txt
	masscanscan(masscaniplist,masscantimeout)
	// 根据 masscan 的 ipalive.txt 端口判断  
	color.Green.Print("[+]")
	fmt.Println(" masscan扫描完成，文件存储在",masscaniplist)
	color.Green.Print("[+]")
	fmt.Println(" 开始nmap扫描... ")
	nmapscan(masscaniplist) 
	color.Green.Print("[+]")
	fmt.Println(" masscan+nmap 扫描完成，文件存储在",nmapiplist)                  
}



func nmapscan(masscaniplist string) {
	// 读取写入
	iplisttest, numbertest, _:=  functions.ReadLines(masscaniplist)
	fmt.Printf("扫描数量 %d ",numbertest)

	// 使用并发进行nmap端口识别
	informations := make(chan NmapRequests, nmapthreads)
	var wg sync.WaitGroup
	startTime := time.Now()

	// 显示进度条
	bar = pb.New(len(iplisttest))
	bar.ShowSpeed = false
	bar.ShowTimeLeft = false
	bar.Start()


	for i := 0; i < cap(informations); i++ {
		go worker(informations, &wg)
	}

	for _,v := range iplisttest {
		wg.Add(1)
		host,port := regxhostport(v)
		informations <-  NmapRequests{Host: host,Port: port}
	}


	wg.Wait()
	close(informations)
	finishMessage := fmt.Sprintf(" nmap指纹识别共耗时 : %v\n\n", time.Since(startTime))
	bar.FinishPrint(fmt.Sprintf(finishMessage))
}



func masscanscan(masscaniplist string, masscantimeout int) {
	// 设置超时   cancel  用于释放资源 https://golang.org/pkg/context/
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(masscantimeout)*time.Minute)
	defer cancel()
	// 新建一个扫描类
	scanner, err := masscan.NewScanner(
		masscan.WithTargets(justString),
		masscan.WithPorts(masscanports),
		masscan.WithContext(ctx),
		masscan.MasscanThreads(masscanthreads),    
	)
	// 开始扫描
	if err != nil {
		log.Fatalf("unable to create masscan scanner: %v", err)
	}

	result, _, err := scanner.Run()
	if err != nil {
		log.Fatalf("unable to run masscan scan: %v", err)
	}
	// 输出结果
	for _, host := range result.Hosts {

		if len(host.Ports) == 0 || len(host.Addresses) == 0 {
			continue
		}

		for _, port := range host.Ports {
			portint := strconv.Itoa(int(port.ID))
			portstring := string(portint)
			address := fmt.Sprintf("%q  %s \n",host.Addresses[0],portstring)
			tracefile(address,masscaniplist)
		}
	}
	color.Green.Print("[+]")
	fmt.Printf(" Masscan 扫描完成: %d 主机，共耗时 %.2f 秒\n", len(result.Hosts), result.Stats.Finished.Elapsed)
}


// 执行任务
func worker(informations chan NmapRequests, wg *sync.WaitGroup) {
	for v := range informations {
		// 使用 nmap 扫描端口，将结果写入文件 
		nmapscans(v.Host,v.Port)
		// 超时测试 time.Sleep(time.Duration(2)*time.Second)
		// 进度条显示
		if bar != nil {
			bar.Increment()
		}
		wg.Done()
	}	
}

// 使用 nmap 进行端口扫描
func nmapscans(host,port string) {
	// 使用nmap进行端口验证
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(nmaptimeout)*time.Minute)
	defer cancel()
	// 新建一个扫描类
	scanner, err := nmap.NewScanner(
		nmap.WithTargets(host),
		nmap.WithPorts(port),
		nmap.WithContext(ctx),
		nmap.WithSkipHostDiscovery(),
	)
	// 开始扫描
	if err != nil {
		log.Fatalf("unable to create nmap scanner: %v", err)
	}	
	result, _, err := scanner.Run()
	if err != nil {
		log.Fatalf("unable to run nmap scan: %v", err)
	}

	for _, host := range result.Hosts {
		if len(host.Ports) == 0 || len(host.Addresses) == 0 {
			continue
		}

		for _, port := range host.Ports {
			portint := strconv.Itoa(int(port.ID))
			portstring := string(portint)
			//  host.Addresses[0] 主机名字 portstring 端口 port.State 是否开放  port.Service.Name 服务名字
			address := fmt.Sprintf("%q  %s  %s [%s] \n",host.Addresses[0], portstring, port.State, port.Service.Name)
			resulthost, resultport := regxhostport(address)

			// 进行格式整合(这里暂时不采用nmap判断的端口是否开启)
			resultaddress := fmt.Sprintf("%s:%s|%s \n",resulthost, resultport, port.Service.Name)

			//  写入txt文件
			tracefile(resultaddress,nmapiplist)
		}
	}
}


// 用于写入数据
func tracefile(str_content,name string)  {

    fd,_:=os.OpenFile(name,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)

    fd_content:=str_content
    buf:=[]byte(fd_content)
    fd.Write(buf)
    fd.Close()
}



//  正则匹配主机与端口
func regxhostport(v string) (host,port string) {
	var nmaphost = regexp.MustCompile(`"(.*)"`)
	hoststring := nmaphost.FindStringSubmatch(v)

	var nmapport = regexp.MustCompile(`\ (\d{1,5})\ `)
	portstring := nmapport.FindStringSubmatch(v)


	host = fmt.Sprintf("%s", hoststring[1])
	port = fmt.Sprintf("%s", portstring[1])
	return host,port
	
}


















