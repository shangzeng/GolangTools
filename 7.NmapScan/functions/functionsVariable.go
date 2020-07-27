package functions



var(
	Banner   = `

 ████     ██                                ████████                           
░██░██   ░██                       ██████  ██░░░░░░                            
░██░░██  ░██ ██████████   ██████  ░██░░░██░██         █████   ██████   ███████ 
░██ ░░██ ░██░░██░░██░░██ ░░░░░░██ ░██  ░██░█████████ ██░░░██ ░░░░░░██ ░░██░░░██
░██  ░░██░██ ░██ ░██ ░██  ███████ ░██████ ░░░░░░░░██░██  ░░   ███████  ░██  ░██
░██   ░░████ ░██ ░██ ░██ ██░░░░██ ░██░░░         ░██░██   ██ ██░░░░██  ░██  ░██
░██    ░░███ ███ ░██ ░██░░████████░██      ████████ ░░█████ ░░████████ ███  ░██
░░      ░░░ ░░░  ░░  ░░  ░░░░░░░░ ░░      ░░░░░░░░   ░░░░░   ░░░░░░░░ ░░░   ░░ 
	                                              
	                                                version:1.0.0 @shangzeng							
	` 

	NmapIplist          []string
	NmapIpnumber        int


    IplistFile   	    string
    Nmapports           string
    Nmapresults         string
    Nmapthreads         int
    Nmaptimeout         int
    Nmapresultfile      string

    Passwordlist    []string
    Usernamelist    []string



)



type configuration struct {

    Ipfile   	    string
    Nmapports       string
    Nmapresults     string
    Nmapresultfile  string
    Nmapthreads     int
    Nmaptimeout     int

}


// 设置nmap请求路径
type NmapRequests struct {

	Host		    string
	Ports           string
}




// Nmap 结果接收
type Service struct {
    Ip       string   // IP
    Port     int      // 端口
    Protocol string   // 协议名称
}

// 弱口令爆破任务
type ServiceWeakPassword struct {
    Server       Service   
    Username     string    // 用户名
    Password     string    // 密码    
}




























