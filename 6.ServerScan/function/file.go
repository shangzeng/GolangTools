package function


import(
	"ServerScan/model"
	"ServerScan/models"  
	"strings"
	"strconv"
	"os"
	"bufio"
	"fmt"
	"time"
	"sync"
	"net"

	// 根据插件进行注册 - 弱口令
	_ "ServerScan/weakpassword/sshweakpassword"
	_ "ServerScan/weakpassword/ftpweakpassword"
	_ "ServerScan/weakpassword/redisweakpassword"

	// 根据插件进行注册 - 主机漏洞扫描
	_ "ServerScan/vulnerability/ms17010"
)


var (
	AliveAddr []models.Service
	mutex     sync.Mutex
)


// 读取文件数据，返回
func ReadIpList(fileName string) (ipList []models.Service) {
	ipListFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Open ip List file err, %v", err)
		os.Exit(0)
	}

	defer ipListFile.Close()

	scanner := bufio.NewScanner(ipListFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		ipPort := strings.TrimSpace(line)
		t := strings.Split(ipPort, ":")
		ip := t[0]
		portProtocol := t[1]
		tmpPort := strings.Split(portProtocol, "|")
		port, _ := strconv.Atoi(tmpPort[0])
		protocol := strings.ToUpper(tmpPort[1])

		addr := models.Service{Ip: ip, Port: port, Protocol: protocol}
		ipList = append(ipList, addr)
	}

	return ipList
}


// 用于检查加载的IP端口是否存活，将存活的写入下一轮进行爆破
// 还没写完 嘤嘤嘤
func CheckAlive(requestlist []models.Service) (aliveIpList []models.Service){
	// 可以考虑加上进度条
	// 使用并发
	var wg sync.WaitGroup
	wg.Add(len(requestlist))
	for _, addr := range requestlist {
		go func(addr models.Service) {
			defer wg.Done()
			SaveAddr(check(addr))
		}(addr)
	}
	wg.Wait()
	return AliveAddr 
}



// 用于检查TCP/UDP 链接是否存在 
func check(ipAddr models.Service) (bool, models.Service) {
	// 默认不存活
	alive := false
	// 判断协议，选择链接方式
	if models.UdpProtocols[ipAddr.Protocol] {
		_, err := net.DialTimeout("udp", fmt.Sprintf("%v:%v", ipAddr.Ip, ipAddr.Port), models.Timeout)
		if err == nil {
			alive = true
		}
	} else {
		_, err := net.DialTimeout("tcp", fmt.Sprintf("%v:%v", ipAddr.Ip, ipAddr.Port), models.Timeout)
		if err == nil {
			alive = true
		}
	}
	return alive, ipAddr
}

func SaveAddr(alive bool, ipAddr models.Service) {
	if alive {
		mutex.Lock() // 用于保证保证在同一时间,只能有同一个函数对该变量的操作 资源保护锁
		AliveAddr = append(AliveAddr, ipAddr)
		mutex.Unlock()
	}
}

// 接收探测的存活端口的IP ，写入任务进行扫描
func RunTask(requestlist []models.ServiceWeakPassword ) {
	wg := &sync.WaitGroup{}
	taskChan := make(chan models.ServiceWeakPassword, models.Threads*2)
	for i := 0; i < models.Threads; i++ {
		go work(taskChan, wg)
	}

	// 生产者，不断地往taskChan channel发送数据，直到channel阻塞 （摘抄大佬解释）
	for _, request := range requestlist {
		wg.Add(1)
		taskChan <- request
	}

	close(taskChan)
	// 从大佬那里学的超时间方式
	waitTimeout(wg, models.Timeout*2)

}



// waitTimeout waits for the waitgroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}


// 测试
// 主要用于执行task中的任务
// 这里根据models中的map值的不同，进入不同的脚本进行扫描
func work(taskChan chan models.ServiceWeakPassword, wg *sync.WaitGroup) {
	for task := range taskChan {
		//fmt.Println(task)
		models.Start(task)
		wg.Done()
	}
}



// 将账号密码装载进入任务进行扫描
func AddUserAndPass (aliveIpList []models.Service, users []string, passwords []string) (aliveIpListtask []models.ServiceWeakPassword){
	for _, user := range users {
		for _, password := range passwords {
			for _, addr := range aliveIpList {
				service := models.ServiceWeakPassword{Server: addr, Username: user, Password: password}
				aliveIpListtask = append(aliveIpListtask, service)
			}
		}
	}
	return aliveIpListtask
}




//  装载用户名
func ReadUserDict(userDict string) (users []string) {
	file, err := os.Open(userDict)
	if err != nil {
		fmt.Println("Open user dict file err", err)
		os.Exit(0)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		user := strings.TrimSpace(scanner.Text())
		if user != "" {
			users = append(users, user)
		}
	}
	return users
}







// 装载密码
func ReadPasswordDict(passDict string) (password []string) {
	file, err := os.Open(passDict)
	if err != nil {
		fmt.Println("Open password dict file err", err)
		os.Exit(0)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		passwd := strings.TrimSpace(scanner.Text())
		if passwd != "" {
			password = append(password, passwd)
		}
	}
	password = append(password, "")
	return password
}





// 用于筛选出是弱口令爆破的服务
func Checkmodels(aliveIpList []models.Service) (weakpasswordaliveIpList,httpliveIpList,vulnerabillityaliveIpList,otheraliveIpList []models.Service){
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





// 用于装载和运行 主机漏洞扫描任务

func RunVulnerabillityTask(requestlist []models.Service ) {
	wg := &sync.WaitGroup{}
	taskChan := make(chan models.Service, models.Threads*2)
	for i := 0; i < models.Threads; i++ {
		go Vulnerabillitywork(taskChan, wg)
	}

	// 生产者，不断地往taskChan channel发送数据，直到channel阻塞 （摘抄大佬解释）
	for _, request := range requestlist {
		wg.Add(1)
		taskChan <- request
	}

	close(taskChan)
	// 从大佬那里学的超时间方式
	waitTimeout(wg, models.Timeout*2)

}


// 测试
// 主要用于执行task中的任务
// 这里根据models中的map值的不同，进入不同的脚本进行扫描
func Vulnerabillitywork(taskChan chan models.Service, wg *sync.WaitGroup) {
	for task := range taskChan {
		//fmt.Println(task)
        model.Start(task)
		wg.Done()
	}
}









