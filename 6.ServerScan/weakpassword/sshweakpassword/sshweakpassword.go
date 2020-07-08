package sshweakpassword


import(
	"fmt"
	"net"
	"time"
	"ServerScan/models"
	"golang.org/x/crypto/ssh"
)



func init() {
    p_1 := plugin_1{}
    models.Regist("ssh 弱口令爆破插件", p_1)
}


type plugin_1 struct {
}


// 告诉插件我没生病
func (this plugin_1) Flag() string {
    return "SSH"
}


// 告诉插件我是干啥的
func (this plugin_1) Start(IP models.ServiceWeakPassword) {
    //fmt.Println(IP)
    result := false
    config := &ssh.ClientConfig{
		User: IP.Username,
		Auth: []ssh.AuthMethod{
			ssh.Password(IP.Password),
		},
		Timeout: 3 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", IP.Server.Ip, IP.Server.Port), config)
	if err == nil {
		defer client.Close()
		session, err := client.NewSession()
		errRet := session.Run("echo xsec")
		if err == nil && errRet == nil {
			defer session.Close()
			result = true
		}
	}
	if result { fmt.Println(IP)}

}