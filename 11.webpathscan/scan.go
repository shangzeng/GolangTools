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
	titleRe = regexp.MustCompile(`>(.*?)\s?</title>`)


)


func init() {
	threads  = 100                         // 线程数量
	timeout  = 5                           // 超时
	path     = "spring-Eureka.txt"         // 请求路径(修改成加载路径脚本)
	filename = "url.txt"                   // 读取网址路径
	savename = "save.txt"                  // 存储结果路径
}


func main() {
	startTime := time.Now()
	ch = make(chan bool, threads)

	// 显示网址数量
	urllist , num, err := readLines(filename)
	if err != nil {
		fmt.Println("读取路径数据错误(url) ！" , err)
		os.Exit(0)
	}
	fmt.Println("  共扫描网址：",num)


	// 显示path数量
	pathlist, pathnum, errs := readLines(path)
	if errs != nil {
		fmt.Println("读取路径数据错误(path) ！" , err)
		os.Exit(0)
	}
	fmt.Println("  共扫描路径：",pathnum)


	// 网址路径组合扫描
 	urlall := urlweb(urllist,pathlist)
 	fmt.Println("  共扫描数量：",len(urlall))


	// 显示进度条
	bar = pb.New(len(urlall))
	bar.ShowSpeed = false
	bar.ShowTimeLeft = false
	bar.Start()


	for _, v := range urlall {
		ch <- true
		wg.Add(1)
		go requestworker(v)
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

	req, err = http.NewRequest("GET", url, nil)
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
	resp.Header.Add("User-Agent", "Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1)")
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

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		fmt.Println(url,m)
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




func urlweb(urls []string, paths []string) []string {
	var allurls []string
	var allurl    string
	for _,url := range urls{
		for _,path := range paths {
			allurl = url + path
			allurls = append(allurls,allurl)
		}
	}
	return allurls
}