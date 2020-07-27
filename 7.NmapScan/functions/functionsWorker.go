package functions


import(
	"fmt"
	"log"
	"sync"
	"time"
	"strconv"
	"context"
	"NmapScan/nmap"
	"github.com/360EntSecGroup-Skylar/excelize"
)






func Worker(informations chan NmapRequests, wg *sync.WaitGroup) {
	for v := range informations {
		// 使用 nmap 扫描端口，将结果写入文件 
		NmapScans(v.Host,v.Ports)
		//time.Sleep(time.Duration(2)*time.Second)
		// 进度条显示
		wg.Done()
	}	
}





func NmapScans(host,port string) {
	//fmt.Println(host,port)
	// 使用nmap进行端口验证
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(Nmaptimeout)*time.Minute)
	defer cancel()
	// 新建一个扫描类
	scanner, err := nmap.NewScanner(
		nmap.WithTargets(host),
		nmap.WithPorts(port),
		nmap.WithContext(ctx),
		nmap.WithSYNScan(),
		nmap.WithOpenOnly(),
		nmap.WithSkipHostDiscovery(),
		nmap.WithDisabledDNSResolution(),
	)

	// 开始扫描
	if err != nil {
		log.Fatalf("unable to create nmap scanner: %v", err)
	}	
	result, _, err := scanner.Run()
	if err != nil {
		log.Fatalf("unable to run nmap scan: %v", err)
	}

	// 输出结果
	for _, host := range result.Hosts {

		if len(host.Ports) == 0 || len(host.Addresses) == 0 {
			continue
		}

		for _, port := range host.Ports {
			//fmt.Printf("\tPort %d/%s %s %s\n", port.ID, port.Protocol, port.State, port.Service.Name)

			// 端口转换
			portint := strconv.Itoa(int(port.ID)) 
			portstring := string(portint)
			// 主机转换
			address := fmt.Sprintf("%q",host.Addresses[0])
			// 服务名字
			servername  := fmt.Sprintf("%s",port.Service.Name)
			// 端口状态转换
			//isopen      := fmt.Sprintf("%s", port.State)
			//excle360(address , portstring, isopen, servername, number)
			all := fmt.Sprintf("%s:%s|%s\n",address , portstring, servername)
			tracefile(all,"iplist/nmapresult.txt")
		}
	}
	//fmt.Printf("Nmap done: %d hosts up scanned in %.2f seconds\n", len(result.Hosts), result.Stats.Finished.Elapsed)
}









// 生成excel表格等待结果写入
func ExcleCreate(filename string) {
	f := excelize.NewFile()
	f.SetCellValue("Sheet1", "A1", "IP")
	f.SetCellValue("Sheet1", "B1", "端口")
	f.SetCellValue("Sheet1", "C1", "状态")
	f.SetCellValue("Sheet1", "D1", "协议")
	//设置宽度
	f.SetColWidth("Sheet1", "A", "B", 30)
	f.SetColWidth("Sheet1", "C", "G", 15)
	if err := f.SaveAs(filename); err != nil {
		fmt.Println("生成表格错误 ！")
	}			
}





// 写入 excle 
func excle360(IP string, port string, xieyi string, ipopen string, listnumber int ) {

	f, err1 := excelize.OpenFile(Nmapresultfile)
    if err1 != nil {
    	fmt.Println("1")
    }

	list1 := fmt.Sprintf("A%d", listnumber+2)
	list2 := fmt.Sprintf("B%d", listnumber+2)
	list3 := fmt.Sprintf("C%d", listnumber+2)
	list4 := fmt.Sprintf("D%d", listnumber+2)

	f.SetCellValue("Sheet1",list1,IP)
	f.SetCellValue("Sheet1",list2,port)
	f.SetCellValue("Sheet1",list3,xieyi)
	f.SetCellValue("Sheet1",list4,ipopen)


	if err := f.SaveAs(Nmapresultfile); err != nil {
        fmt.Println(err)
    }
}




// 用于筛选服务类型
func Checkmodels(aliveIpList []Service) (weakpasswordaliveIpList,httpliveIpList,vulnerabillityaliveIpList,otheraliveIpList []Service){
	for _, aliveIpListmodel := range aliveIpList {
		if aliveIpListmodel.Protocol == "SSH" || aliveIpListmodel.Protocol == "FTP" || aliveIpListmodel.Protocol == "REDIS"{
			weakpasswordaliveIpList = append(weakpasswordaliveIpList, aliveIpListmodel)
		} else if aliveIpListmodel.Protocol == "HTTP" {
			httpliveIpList =  append(httpliveIpList, aliveIpListmodel)
		} else if aliveIpListmodel.Protocol == "MICROSOFI-DS"{
			vulnerabillityaliveIpList =  append(vulnerabillityaliveIpList, aliveIpListmodel)
		} else {
			otheraliveIpList = append(otheraliveIpList, aliveIpListmodel)
		}
	}
	return weakpasswordaliveIpList, httpliveIpList, vulnerabillityaliveIpList, otheraliveIpList
}















