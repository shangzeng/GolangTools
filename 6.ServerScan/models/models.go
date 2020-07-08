package models


import (
	//"ServerScan/weakpassword"
	"time"
)


var(
	
	Iplist          string
	Passwordlist    string
	Usernamelist    string
	Httplever       string
	Threads         int
	Timeout = 3 * time.Second


	// 在扫描前的端口判断是否存活，根据服务判断是TCP还是UDP
	// 注意： 如果要增加漏洞模式需要在这里进行注册
	PortNames = map[int]string{
		21:    "FTP",
		22:    "SSH",
		161:   "SNMP",
		445:   "SMB",
		1433:  "MSSQL",
		3306:  "MYSQL",
		5432:  "POSTGRESQL",
		6379:  "REDIS",
		9200:  "ELASTICSEARCH",
		27017: "MONGODB",
	}

	UdpProtocols = map[string]bool{
		"SNMP": true,
	}
)



type Service struct {
	Ip       string   // IP
	Port     int      // 端口
	Protocol string   // 协议名称
}

type ServiceWeakPassword struct {
	Server       Service   
	Username     string    // 用户名
	Password     string    // 密码	
}



