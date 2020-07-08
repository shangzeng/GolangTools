package redisweakpassword


import(
	"fmt"
	"time"
    "ServerScan/models"
    "github.com/go-redis/redis"

)



func init() {
    p_1 := plugin_1{}
    /* 与容器进行连接 */
    models.Regist("redis弱口令扫描 插件", p_1)
}


type plugin_1 struct {
}


// 告诉插件我没生病
func (this plugin_1) Flag() string {
    return "REDIS"
}


// 告诉插件我是干啥的
func (this plugin_1) Start(IP models.ServiceWeakPassword) {
    //fmt.Println(IP)
    result :=  false
    opt := redis.Options{Addr: fmt.Sprintf("%v:%v", IP.Server.Ip, IP.Server.Port),
        Password: IP.Password, DB: 0, DialTimeout: 3 * time.Second}
    client := redis.NewClient(&opt)
    defer client.Close()
    _, err := client.Ping().Result()
    if err == nil {
        result = true
    }
    if result { fmt.Println(IP)}
}