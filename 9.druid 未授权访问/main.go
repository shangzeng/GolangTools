package main


import (
	"fmt"
	"time"
	"sync"
	"bufio"
	"os"
	"net"
	"regexp"
	"net/http"
	"crypto/tls"
	"io/ioutil"
	"golang.org/x/text/transform"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/net/html/charset"

)


var (
	ch 			chan bool
	wg          sync.WaitGroup
	filename    string
	savename    string
	path        string
	threads     int
	timeout     int
	titleRe = regexp.MustCompile(`>(.*?)\s?</title>`)

)


func init() {
	threads  = 3                        // 线程数量
	timeout  = 2                         // 超时
	path     = "/druid/index.html"       // 请求路径
	filename = "butian.txt"              // 读取路径
	savename = "save.txt"                // 存储结果路径
}


func main() {
	ch = make(chan bool, threads)
	urllist , num, err := readLines(filename)
	if err != nil {
		fmt.Println("读取路径数据错误 ！" , err)
		os.Exit(0)
	}
	fmt.Println("共扫描：",num)


	for _, v := range urllist {
		urlpath := v + path
		ch <- true
		wg.Add(1)
		go requestworker(urlpath)
	}
	wg.Wait()
}


func requestworker(url string) {
	defer func() {
		<-ch
		wg.Done()
	}()

	var req *http.Request
	var err error

	req, err = http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("请求错误 ！")
	}

	client := createHTTPClient()
	time.Sleep(1e9) 
	resp, err1 := client.Do(req)
	if err1 != nil {
		fmt.Println("网站超时错误 ！",url,err1)
		return
	}

	defer resp.Body.Close()

	reader := bufio.NewReader(resp.Body)
	e := determineEncoding(reader)
	utf8Reader := transform.NewReader(reader, e.NewDecoder())
	bodyss, err := ioutil.ReadAll(utf8Reader)
	if err != nil {
		bodyss = []byte("")
	}

	respBody := string(bodyss)
	m := titleRe.FindStringSubmatch(respBody)
	if len(m) >= 2 {
		if string(m[1]) == "Druid Stat Index" {
			tracefile(url+"\n",savename)
		}
	}
}


func tracefile(str_content,savename string)  {
    fd,_:=os.OpenFile(savename,os.O_RDWR|os.O_CREATE|os.O_APPEND,0644)
    fd_content:=str_content
    buf:=[]byte(fd_content)
    fd.Write(buf)
    fd.Close()
}


func determineEncoding(r *bufio.Reader) encoding.Encoding {
	b, err := r.Peek(1024)
	if err != nil {
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(b, "")
	return e
}


func createHTTPClient() *http.Client {
	// 不校验证书
	tr := &http.Transport{
		Dial: (&net.Dialer{
			Timeout:   time.Duration(timeout) * time.Second,
			Deadline:  time.Now().Add(time.Duration(timeout) * time.Second),
			KeepAlive: time.Duration(timeout) * time.Second,
		}).Dial,
		TLSHandshakeTimeout: time.Duration(timeout) * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true}}

	// 主要是这里构造超时
	client := &http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: tr,
	}
	return client
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