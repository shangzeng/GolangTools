package ftpweakpassword


import(
	"fmt"
	"time"
	"ServerScan/models"
	"github.com/jlaffaye/ftp"
)



func init() {
    p_1 := plugin_1{}
    /* 与容器进行连接 */
    models.Regist("ftp弱口令扫描 插件", p_1)
}


type plugin_1 struct {
}


// 告诉插件我没生病
func (this plugin_1) Flag() string {
    return "FTP"
}


// 告诉插件我是干啥的
func (this plugin_1) Start(IP models.ServiceWeakPassword) {
    //fmt.Println(IP)
    result := false
    conn, err := ftp.DialTimeout(fmt.Sprintf("%v:%v", IP.Server.Ip, IP.Server.Port), 3 * time.Second)
    if err == nil {
    	err = conn.Login(IP.Username, IP.Password)
    	if err == nil {
    		defer conn.Logout()
    		result = true
    	}
    }
    if result { fmt.Println(IP)}
}