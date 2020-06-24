package main

import (

    "encoding/json"
    "strings"
    "flag"
    "fmt"
    "os"
    "strconv"
    "net/http"
    "io/ioutil"
    "encoding/base64"
    "github.com/gookit/color"
    "github.com/360EntSecGroup-Skylar/excelize"

)


var (

	banner   = `

	   ___          ___            ____    _____      
	 / ___\       / ___\          /\  _ \ /\  __ \    
	/\ \__/  ___ /\ \__/   __     \ \ \L\_\ \ \/\ \   
	\ \  __\/ __ \ \  __\/ __ \    \ \ \L_L\ \ \ \ \  
	 \ \ \_/\ \L\ \ \ \_/\ \L\.\_   \ \ \/, \ \ \_\ \ 
	  \ \_\\ \____/\ \_\\ \__/.\_\   \ \____/\ \_____\
	   \/_/ \/___/  \/_/ \/__/\/_/    \/___/  \/_____/  

				fofa_GO version:1.0.2 @shangzeng
	`


	email   	string
	key     	string
	qouer   	string
	attribute   string
	excle	    string
	history 	bool 
	write       bool
	qouersize   int
	numbers	    int

)


type configuration struct {
    Email   	string
    Key     	string
    Attribute   string
}


type Fofaresponse struct {
	Mode    string     `json:"mode"`
	Error   bool       `json:"error"`
	Query   string     `json:"query"`
	Page    int        `json:"page"`
	Size    int        `json:"size"`
	Results [][]string `json:"results"`
}


func usage() {
    flag.PrintDefaults()
}

func init() {
	fmt.Println(banner)
	// 加载配置文件内容
	file, err := os.Open("config.json")
	if err != nil {
		fmt.Println("配置文件读取错误 ！")
		os.Exit(0)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	conf := configuration{}
	err1 := decoder.Decode(&conf)
	if err1 != nil {
		fmt.Println("编码错误 ！请检查json文件格式 ！")
		os.Exit(0)
	}
	email = conf.Email
	key   = conf.Key
	attribute = conf.Attribute


	// 加载命令行交互
	flag.StringVar(&qouer, "q", "", "FOFA 查询语句:注意转译例如：-q domain=\\\"jljcxy.com\\\" ")
	flag.BoolVar(&history, "h", false, "-h 获取历史数据，默认为false")
	flag.StringVar(&excle, "e", "", "生成详细信息的excel文件，例如： -e xxx资产梳理报告.xlsx ")
	flag.IntVar(&qouersize, "s", 10000, "显示/下载数量，默认为搜索总量，例如： -s 50")
	flag.Usage = usage
	flag.Parse()
	if qouer == "" {
		flag.Usage()
		os.Exit(0)
	}


	if excle != "" {
		f := excelize.NewFile()
		f.SetCellValue("Sheet1", "A1", "IP")
		f.SetCellValue("Sheet1", "B1", "子域名")
		f.SetCellValue("Sheet1", "C1", "端口")
		f.SetCellValue("Sheet1", "D1", "协议")
		f.SetCellValue("Sheet1", "E1", "国家")
		f.SetCellValue("Sheet1", "F1", "省")
		f.SetCellValue("Sheet1", "G1", "城市")
		f.SetCellValue("Sheet1", "H1", "信息")
		f.SetCellValue("Sheet1", "I1", "主域名")
		//设置宽度
		f.SetColWidth("Sheet1", "A", "B", 30)
		f.SetColWidth("Sheet1", "C", "G", 15)
		f.SetColWidth("Sheet1", "H", "H", 100)
		f.SetColWidth("Sheet1", "I", "I", 15)	
		if err := f.SaveAs(excle); err != nil {
			fmt.Println("生成表格错误 ！")
		}			
		write = true
		
	}	

	// 这里写的有问题
	if qouersize != 10000 && qouersize <= 10000 {
		numbers = qouersize
	} else if qouersize > 10000{
		color.Red.Print("[+]")
		fmt.Println(" 数量错误  ：最多只能显示10000条 ！")
		os.Exit(0)
	} else {
		numbers = qouersize
	}
}

func zhuanhuan(zhuan []string) string {
	str2 := strings.Join(zhuan, " ")
	return str2
}

func fofa_requests(url string) []byte {
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




// 写入 excle 
func excle360(IP string, domain string, port string, xieyi string, country string,city string, region string, information string, Fdomain string, listnumber int, exclename string ) {

	f, err1 := excelize.OpenFile(exclename)
    if err1 != nil {
    	fmt.Println("1")
    }

	list1 := fmt.Sprintf("A%d", listnumber+2)
	list2 := fmt.Sprintf("B%d", listnumber+2)
	list3 := fmt.Sprintf("C%d", listnumber+2)
	list4 := fmt.Sprintf("D%d", listnumber+2)
	list5 := fmt.Sprintf("E%d", listnumber+2)
	list6 := fmt.Sprintf("F%d", listnumber+2)
	list7 := fmt.Sprintf("G%d", listnumber+2)
	list8 := fmt.Sprintf("H%d", listnumber+2)
	list9 := fmt.Sprintf("I%d", listnumber+2)

	f.SetCellValue("Sheet1",list1,IP)
	f.SetCellValue("Sheet1",list2,domain)
	f.SetCellValue("Sheet1",list3,port)
	f.SetCellValue("Sheet1",list4,xieyi)
	f.SetCellValue("Sheet1",list5,country)
	f.SetCellValue("Sheet1",list6,city)
	f.SetCellValue("Sheet1",list7,region)
	f.SetCellValue("Sheet1",list8,information)
	f.SetCellValue("Sheet1",list9,Fdomain)

	if err := f.SaveAs(exclename); err != nil {
        fmt.Println(err)
    }
}

func main() {

	// 请求并输出结果
	base64decode := base64.StdEncoding.EncodeToString([]byte(qouer))
	historys := strconv.FormatBool(history)
	url := fmt.Sprintf("https://fofa.so/api/v1/search/all?email=%s&key=%s&fields=%s&page=1&qbase64=%s&full=%s",email,key,attribute,base64decode,historys)
	//fmt.Println(url)
	results := fofa_requests(url)	


	// 处理 json 信息
	var fofajson Fofaresponse
	err := json.Unmarshal(results, &fofajson)
	if err != nil {
		fmt.Println("json 处理错误 ！")
		os.Exit(0)
	}


	// 显示搜索内容
	color.Green.Print("[+]")
	fmt.Println(" 搜索语句 ：", fofajson.Query)

	if fofajson.Size == 0 {
		color.Red.Print("[+]")
		fmt.Println(" 搜索总量 ：", fofajson.Size)
		os.Exit(0)





	} else if fofajson.Size > 0 && fofajson.Size < 100 {
		color.Yellow.Print("[+]")
		fmt.Println(" 搜索总量 ：", fofajson.Size)

		// 命令行显示详细信息 
		fmt.Printf("|%-16s|%-3s|%-26s|%10s\n","ip地址","端口","域名","网站标题")
		for number,resultslist := range fofajson.Results{
			if number > numbers-1 {
				break
			}  
			data, err := json.Marshal(resultslist)
			if err != nil {
				panic(err)
			}
			var fofalistjson []string
			errs := json.Unmarshal(data, &fofalistjson)
			if errs != nil {
				panic(err)
			}

			fmt.Printf("|%-18s|%-5s|%-28s|%10s\n",zhuanhuan(fofalistjson[0:1]),zhuanhuan(fofalistjson[2:3]),zhuanhuan(fofalistjson[1:2]), zhuanhuan(fofalistjson[7:8]))
			if write == true {
				excle360(zhuanhuan(fofalistjson[0:1]),zhuanhuan(fofalistjson[1:2]),zhuanhuan(fofalistjson[2:3]),zhuanhuan(fofalistjson[3:4]),zhuanhuan(fofalistjson[4:5]),zhuanhuan(fofalistjson[5:6]),zhuanhuan(fofalistjson[6:7]),zhuanhuan(fofalistjson[7:8]),zhuanhuan(fofalistjson[8:9]),number,excle)
			}
		}





	} else if fofajson.Size > 100 && fofajson.Size < 10000{
		var urlsize string
		color.Green.Print("[+]")
		fmt.Println(" 搜索总量 ：", fofajson.Size)
		if numbers != 10000 {
			color.Green.Print("[+]")
			fmt.Println(" 显示数量 ：", numbers)
			urlsize = fmt.Sprintf("https://fofa.so/api/v1/search/all?email=%s&key=%s&fields=%s&page=1&size=%d&qbase64=%s&full=%s",email,key,attribute,numbers,base64decode,historys)
		} else {
			urlsize = fmt.Sprintf("https://fofa.so/api/v1/search/all?email=%s&key=%s&fields=%s&page=1&size=%d&qbase64=%s&full=%s",email,key,attribute,fofajson.Size,base64decode,historys)
		}
		resultsize := fofa_requests(urlsize)
		//fmt.Println(urlsize)


		var fofajsonsize Fofaresponse
		err1 := json.Unmarshal(resultsize, &fofajsonsize)
		if err1 != nil {
			fmt.Println("json1 处理错误 ！")
			os.Exit(0)
		} 	

		//fmt.Println(fofajsonsize.Results)



		if !fofajsonsize.Error {
			// 命令行显示详细信息 
			fmt.Printf("|%-16s|%-3s|%-26s|%10s\n","ip地址","端口","域名","网站标题")
			for number,resultslist := range fofajsonsize.Results{
				if number > numbers-1 {
					break
				} 
				data, err := json.Marshal(resultslist)
				if err != nil {
					panic(err)
				}
				var fofajsonsize []string
				errs := json.Unmarshal(data, &fofajsonsize)
				if errs != nil {
					panic(err)
				}

				fmt.Printf("|%-18s|%-5s|%-28s|%10s\n",zhuanhuan(fofajsonsize[0:1]),zhuanhuan(fofajsonsize[2:3]),zhuanhuan(fofajsonsize[1:2]), zhuanhuan(fofajsonsize[7:8]))
				if write == true {
					excle360(zhuanhuan(fofajsonsize[0:1]),zhuanhuan(fofajsonsize[1:2]),zhuanhuan(fofajsonsize[2:3]),zhuanhuan(fofajsonsize[3:4]),zhuanhuan(fofajsonsize[4:5]),zhuanhuan(fofajsonsize[5:6]),zhuanhuan(fofajsonsize[6:7]),zhuanhuan(fofajsonsize[7:8]),zhuanhuan(fofajsonsize[8:9]),number,excle)
				}
			}
			


		}else {
			color.Red.Print("[+]")
			fmt.Println(" 搜索错误  ：fofa币不足或者非高级会员 ！")
			color.Yellow.Print("[+]")
			fmt.Println(" 搜索建议  ：使用 -s 参数指定数量 ")
		}




	} else if fofajson.Size > 10000 {
		var urlsize string
		color.Green.Print("[+]")
		fmt.Println(" 搜索总量 ：", fofajson.Size)
		if numbers != 10000 {
			color.Green.Print("[+]")
			fmt.Println(" 显示数量 ：", numbers)
			urlsize = fmt.Sprintf("https://fofa.so/api/v1/search/all?email=%s&key=%s&fields=%s&page=1&size=%d&qbase64=%s&full=%s",email,key,attribute,numbers,base64decode,historys)
			//fmt.Println(urlsize)
		} else {
			color.Red.Print("[+]")
			fmt.Println("在数量一万以上时，未制定参数 -s 暂时无法查询 ！")
			os.Exit(0)
		}
		resultsize := fofa_requests(urlsize)

		var fofajsonsize Fofaresponse
		err1 := json.Unmarshal(resultsize, &fofajsonsize)
		if err1 != nil {
			fmt.Println("json2 处理错误 ！")
			os.Exit(0)
		} 	


		if !fofajsonsize.Error {
			// 命令行显示详细信息 
			fmt.Printf("|%-16s|%-3s|%-26s|%10s\n","ip地址","端口","域名","网站标题")
			for number,resultslist := range fofajsonsize.Results{
				if number > numbers-1 {
					break
				} 
				data, err := json.Marshal(resultslist)
				if err != nil {
					panic(err)
				}
				var fofajsonsize []string
				errs := json.Unmarshal(data, &fofajsonsize)
				if errs != nil {
					panic(err)
				}

				fmt.Printf("|%-18s|%-5s|%-28s|%10s\n",zhuanhuan(fofajsonsize[0:1]),zhuanhuan(fofajsonsize[2:3]),zhuanhuan(fofajsonsize[1:2]), zhuanhuan(fofajsonsize[7:8]))
				if write == true {
					excle360(zhuanhuan(fofajsonsize[0:1]),zhuanhuan(fofajsonsize[1:2]),zhuanhuan(fofajsonsize[2:3]),zhuanhuan(fofajsonsize[3:4]),zhuanhuan(fofajsonsize[4:5]),zhuanhuan(fofajsonsize[5:6]),zhuanhuan(fofajsonsize[6:7]),zhuanhuan(fofajsonsize[7:8]),zhuanhuan(fofajsonsize[8:9]),number,excle)
				}
			}


		}else {
			color.Red.Print("[+]")
			fmt.Println(" 搜索错误  ：fofa币不足或者非高级会员 ！")
			color.Yellow.Print("[+]")
			fmt.Println(" 搜索建议  ：使用 -s 参数指定数量 ")
		}




	}

}














