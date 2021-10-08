package main

import(
	"encoding/json"
	"crypto/tls"
	_"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"net/url"
	"time"
	"sync"
	"fmt"
	"os"
)

type Result struct {
	ReqList       []Request `json:"req_list"`
	AllReqList    []Request `json:"all_req_list"`
	AllDomainList []string  `json:"all_domain_list"`
	SubDomainList []string  `json:"sub_domain_list"`
}

type Request struct {
	Url     string                 `json:"url"`
	Method  string                 `json:"method"`
	Headers map[string]string 	   `json:"headers"`
	Data    string                 `json:"data"`
	Source  string                 `json:"source"`
}

func main() {
	// 读取 url.txt 数据
	urls := []string{"https://insight.rakuten.com/","https://www.rakuten.com/","https://www.victoriassecret.com","https://de.victoriassecret.com"}
	// 代理端口
	proxyUrl := "http://127.0.0.1:7777"
	// 线程
	threads := 2


	// 暂时循环
	for _ ,url := range urls {
		Crawlergo(url,proxyUrl,threads)
	}
}




func Crawlergo(url,proxyUrl string, threads int) {
	cmd := exec.Command("./crawlergos","-c","/Users/songjidong/Desktop/Chromium.app/Contents/MacOS/Chromium","-t","20","-f","smart","--fuzz-path","--output-mode","json",url)
	output, _ := cmd.CombinedOutput()
	jsonout := strings.Split(string(output),"--[Mission Complete]--")[1]

    var Crawlerjson Result
    err := json.Unmarshal([]byte(jsonout), &Crawlerjson)
	if err != nil {
		fmt.Println("json 处理错误 ！")
		os.Exit(0)
	}

	//并发
	informations := make(chan Request, threads)
	var wg sync.WaitGroup


	for i := 0; i < cap(informations); i++ {
		go worker(proxyUrl, informations, &wg)
	}



	for _,resultslist := range Crawlerjson.ReqList{
		wg.Add(1)
		//ProxyURL(proxyUrl,resultslist)
		informations <- Request{Url:resultslist.Url,Method:resultslist.Method,Headers:resultslist.Headers,Data:resultslist.Data,Source:resultslist.Source}
	}

	wg.Wait()
	close(informations)
}


func worker( proxyUrl string, resultslist chan Request, wg *sync.WaitGroup) {

	/*
		1. 代理请求
		2. 跳过https不安全验证
		3. 自定义请求头 User-Agent
	*/ 

	time.Sleep(10 * time.Second)
	for v := range resultslist {

	request, _ := http.NewRequest(v.Method, v.Url, strings.NewReader(v.Data))

	for key, value := range v.Headers {
		request.Header.Set(key, value)
	}

	proxy, _ := url.Parse(proxyUrl)
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 5, //超时时间
	}

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("出错了", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println(v.Method+"请求网址："+v.Url)
	wg.Done()

}

}



