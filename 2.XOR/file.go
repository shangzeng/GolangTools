package main

import (
	"fmt"
	"os"
	"io/ioutil"
	"regexp"
	"flag"

)

var (
	d string
)


//读取文件
func read1(filename string) []byte{
    tet,_ := ioutil.ReadFile(filename)
    return tet
}

//亦或数字返回
func xornum( data []byte, number int) []byte {
	var vals []byte
	for i := 0; i < len(data) ;  i++ {
		hhhh := data[i] ^ byte(number)
		vals = append(vals,hhhh)
	}
	return vals
}


// 命令交互
func init() {
	flag.StringVar(&d, "d", "", "-d 你要亦或的文件名")
}


//用于写入 解密 文件
func write(name int,data []byte) {
	names := fmt.Sprintf("XOR_%d",name)
	f,err := os.Create(names)
	defer f.Close()
	if err !=nil {
		fmt.Println(err.Error())
	} else {
		_,err=f.Write(xornum(data,name))
		}
}


func main () {
	flag.Parse()
	if d == ""  {
	fmt.Println("\033[33m \033[20m[缺少搜索语句！] ：-d <options> \033[0m")
	os.Exit(0)
	}

	data := read1(d)
	for  j := 0 ; j < 500 ; j++ { 
	var validID = regexp.MustCompile(`</script>`)
	panduan := validID.MatchString(string(xornum(data,j))) //true
	if panduan {
		write(j,data)
		break
    	}
	}
}








