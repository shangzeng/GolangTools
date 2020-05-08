package main

import (
    "fmt"
    "flag"
    "os"
    "math/big"
    "bytes"
    "crypto/rand"
    "strings"
)

var (
	banner = `
 _    _ _____  ______      ______  _____  
\ \  / / ___ \(_____ \    / _____)/ ___ \ 
 \ \/ / |   | |_____) )  | /  ___| |   | |
  )  (| |   | (_____ (   | | (___) |   | |
 / /\ \ |___| |     | |  | \____/| |___| |
/_/  \_\_____/      |_|   \_____/ \_____/ 

                         -h for help @shangzeng
	`
	strput string
	whilestr string
	help   bool
	number   int
	second_char = ""
)


//用于编码
func XorEncodeStr(msg, key string) string {
	pwd := ""
	for i := 0; i < len(msg); i++ {
		pwd += (string((key[i]) ^ (msg[i])))
	}
	return pwd
}


//用于生成随机数
func CreateRandomString(len int) string  {
    var container string
    b := bytes.NewBufferString(whilestr)
    length := b.Len()
    bigInt := big.NewInt(int64(length))
    for i := 0;i < len ;i++  {
        randomInt,_ := rand.Int(rand.Reader,bigInt)
        container += string(whilestr[randomInt.Int64()])
    }
    return container
}


//命令行交互
func init() {
	flag.StringVar(&strput, "d", "", "-d 你要进行异或编码的字符串")
	flag.StringVar(&whilestr, "c", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890", "-c 异或范围 默认数字字母大小写")
	flag.BoolVar(&help, "h", false, "帮助")
	flag.IntVar(&number, "n", 1, "-n 输出数量，默认为1")
	flag.Usage = usage
}
func usage() {
    flag.PrintDefaults()
}


//判断字符串内容是否在范围内
func isstrin(isstr string) bool{
	for i := 0 ; i < len(isstr) ; i++ {
		if strings.Index(whilestr, string(isstr[i])) == -1 { return false }
	}
	return true
}


func main() {
	//输入控制 判断
	fmt.Println(banner)
	flag.Parse()
	if help { flag.Usage()}
	if strput == "" {os.Exit(0)}
	for i := 0 ; i < number ; i++ {
		//进行异或判断 ： 写死一个循环，只有达成才会跳出
		for {
			//生成要随之亦或的字符串
			first_char := CreateRandomString(len(strput))
			//根据亦或原理出来第二个 second_char 有可能是特殊字符 
			second_char := XorEncodeStr(first_char,strput)
			//使用函数 判断字符串内容
			if isstrin(second_char) { 
				fmt.Printf("\033[33m \033[38m[+]:\033[0m%s 的异或运算结果为：  %s   %s \n",strput,second_char,first_char)
				break
			}
		}
	}
}
























