Time: 20220531:other:日常记录...:[日常笔记] SpringBootActuator未授权与钓鱼邮件
--------





## 背景

鸽了好久了，想写点啥但是也没啥指得写的最近，只好谢谢日常记录水水文字吧 == 

修图摄影也鸽了很久，这么长时间就修了几张图，待久了啥也不想干了属于是。

![](/static/img/IMG_5018.PNG)



## 关于 Spring Boot Actuator
`Actuator` 是 `Spring Boot` 提供的服务监控和管理中间件。当 `Spring Boot `应用程序运行时，它会自动将多个端点注册到路由进程中。而由于对这些端点的错误配置，就有可能导致一些系统信息泄露、XXE、甚至是 RCE 等安全问题。`Actuator`提供了13个接口:

```bash
GET /actuator 显示能够所有打开的接口
GET /auditevents    显示应用暴露的审计事件 (比如认证进入、订单失败)
GET /beans  描述应用程序上下文里全部的 Bean，以及它们的关系
GET /conditions 就是 1.0 的 /autoconfig ，提供一份自动配置生效的条件情况，记录哪些自动配置条件通过了，哪些没通过
GET /configprops    描述配置属性(包含默认值)如何注入Bean
GET /env    获取全部环境属性
GET /env/{name} 根据名称获取特定的环境属性值
GET /flyway 提供一份 Flyway 数据库迁移信息
GET /liquidbase 显示Liquibase 数据库迁移的纤细信息
GET /health 报告应用程序的健康指标，这些值由 HealthIndicator 的实现类提供
GET /heapdump   dump 一份应用的 JVM 堆信息
GET /httptrace  显示HTTP足迹，最近100个HTTP request/repsponse
GET /info   获取应用程序的定制信息，这些信息由info打头的属性提供
GET /logfile    返回log file中的内容(如果 logging.file 或者 logging.path 被设置)
GET /loggers    显示和修改配置的loggers
GET /metrics    报告各种应用程序度量信息，比如内存用量和HTTP请求计数
GET /metrics/{name} 报告指定名称的应用程序度量值
GET /scheduledtasks 展示应用中的定时任务信息
GET /sessions   如果我们使用了 Spring Session 展示应用中的 HTTP sessions 信息
POST    /shutdown   关闭应用程序，要求endpoints.shutdown.enabled设置为true
GET /mappings   描述全部的 URI路径，以及它们和控制器(包含Actuator端点)的映射关系
GET /threaddump 获取线程活动的快照
```


`Spring Boot` 也主要分成 `1.X` 和 `2.X` 两个版本, 上文是`1.x`版本。 在 `2.X` 的 `Actuator`  也有所不同：



```bash
# 只默认开放这两个端口
/actuator           // 显示能够所有打开的接口
/actuator/health    // 主要用于 查看系统的运行状态
/actuator/info
# 其他的需要手动开启
management.endpoints.web.exposure.include=*    // 全部打开
management.endpoints.web.exposure.exclude=beans,trace  //打开部分 
management.endpoints.web.base-path=/manage   // Actuator 默认所有的监控点路径都在/actuator/* 当然如果有需要这个路径也支持定制,设置完重启后，再次访问地址就会变成/manage/*
```

所有的接口都是可以关闭的（默认开启）， 在`application.properties` 文件中进行设置
```bash
management.endpoint.shutdown.enabled=true
```

### heapdump 信息泄漏

在`/heapdump` 接口会泄漏内存信息，我们可以通过分析内存在获取敏感信息 ， 这里我使用的是`MAC` 中的`MAT` 工具来进行分析


* 选择`File`  , `Open Heap Dump`  , 打开后点击左上角的`OQL`图标，搜索敏感信息命令如下， 点击感叹号图标执行：
```bash
select * from java.util.LinkedHashMap$Entry x WHERE (toString(x.key).contains("password"))
```

也可以使用[heapdump_tool](https://github.com/wyzxxz/heapdump_tool), [JDumpSpider-1.0-SNAPSHOT-full](https://github.com/whwlsfb/JDumpSpider)等工具 记性分析，但是会有遗漏(数据库信息，堆栈中的用户信息，加密解密秘钥信息)，比如下面这条

```bash
jdbc:mysql://rds.prod.na.v2.naiot.com:3306/api_test?useUnicode=true&characterEncoding=UTF-8§{user=api_test, password=Midea1234%}
```

### access 阿里云数据泄漏

主要还是搜索关键词 `application.aliyun.accesskey`  , `application.aliyun.accesssecret` 等等 ，然后使用工具测试。


## 钓鱼流程指南

BOSS要求做一下简单的钓鱼攻击，并提出流程。目标环境有以下特点：

1. 内外网系统隔离，有策略管理
2. 存在邮件安全网关和邮件沙箱 
3. 邮箱、OA 等系统全部存在内网，VPN等使用PIN码进行验证

这样一来，获取权限就比较困难，目标降级为获取敏感信息。采取的方式就很简单了：钓鱼邮件+手机二维码进行钓鱼。


### 使用smtp2go进行钓鱼邮件分发

这里主要是用于发件人信息伪造，经过测试发现其实伪造发件人是个很鸡肋的行为： 能点击邮件的憨憨也基本不会在乎发件人，在乎发件人的也不会点击邮件，麻。

在使用自己注册[smtp2go](https://www.smtp2go.com/)的邮箱后，建立密码,需要手机号码（支持中国手机号），登陆，增加自己的域名 ， 选择 `Settings`  , `STMPUsers`  , `ADD STMP Users`

### 使用gophish进行搭建

* 官方文档： https://docs.getgophish.com/user-guide/template-reference
* 下载地址： https://github.com/gophish/gophish

Gopdish 是一个功能强大的开源网络钓鱼框架，可以有一个强大的api，任何操作都能可以通过api来进行，他可以自动化的进行钓鱼，并且能够统计钓鱼的结果，例如邮件钓鱼的送达率，点开率，触发率，成功率等等，然后汇总成一个报表，让安全人员能够清晰明了的了解到本次钓鱼测试结果。默认账号admin密码在运行终端生成,安装命令如下：

```bash
// 使用ubuntu 系统
wget https://github.com/gophish/gophish/releases/download/v0.11.0/gophish-v0.11.0-linux-64bit.zip
apt install unzip
unzip gophish-v0.11.0-linux-64bit.zip -d gophish
cd gophish

```

在启动之前需要在`config.json`中将管理端口的`127.0.0.1:3333`修改为`0.0.0.0:333`,这样可以在外网进行访问(钓鱼端口和协议也是在这里设置)。账号密码在命令行中随机生成。如果不需要https服务也可以在`config.json`关闭。
```bash
chmod 777 gophish
vi config.json      // 修改配置文件
tmux                // 新建窗口
./gophish
```


### 使用gophish进行邮件伪造发送

在`gophish`的流程如下：

```bash

设置邮件发送服务器 -> 设置钓鱼网站页面 -> 设置邮件模版 -> 添加钓鱼用户组 -> 实施攻击

```

* 登陆后，在`Sending Profiles` 设置发送邮件的邮箱设置，找qq或者网易邮箱都可以进行设置
* 在 `Landing Pages`配置钓鱼网站模版，这里可以使用`import Site` 爬虫爬取模版，但是有一些不准确。下面的选项为抓取密码。注意，在加载钓鱼模版 的时候，有一个加载连接的选项需要勾选，否则使用`gophish`很难接收到数据。再就是需要在页面中存在连接，`gophish`会替换成自己的网址，没有的话就无法返回数据。 最后有一个跳转选项，也可以让中招者在输入敏感数据后跳转到真正的页面上去（总之就是能全勾选就全勾上）
 在钓鱼邮件建立这里还需要注意一点： 接收的账号密码格式也是需要一定程度的修改，否则就抓不到获取的密码，最好还是from表单进行提交的数据，比如最常见的格式如下：
 ```html
<form action="" method="POST">
    <input name="username" type="text" placeholder="username" />
    <input name="password" type="password" placeholder="password" />
    <input type="submit" value="Submit" />
</form>
 ```
上述的请求格式是可以进行转换的，但是后续抓包可以大致看到数据是经过POST传入的, 总之, 要是自动的没生效，那就手动修改传参就行了 ：
 ```http
POST /?rid=HmS249p HTTP/1.1

Host: XX.XX.XX.XX
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:100.0) Gecko/20100101 Firefox/100.0
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8
Accept-Language: zh-CN,zh;q=0.8,zh-TW;q=0.7,zh-HK;q=0.5,en-US;q=0.3,en;q=0.2
Accept-Encoding: gzip, deflate
Content-Type: application/x-www-form-urlencoded
Content-Length: 42
Origin: http://XX.XX.XX.XX
Connection: close
Referer: http://XX.XX.XX.XX/?rid=HmS249p
Cookie: __aegis_uid=1652952007292-7968
Upgrade-Insecure-Requests: 1

username=testaaa&password=qqqqqq&bianhao=1
 ```

* 在 `Email Templates`配置钓鱼网站模版，这里可以使用`import Email` 加载 `eml` 格式的邮件信息，这里需要注意邮件编码格式问题，可以先发送几个实验。这里需要勾选`Add Tracking Image` 这个选项后就可以根据图像知道这封邮件有没有被打开了(内外网隔离无效)。这里要注意 ，当我们想要使用钓鱼功能的时候，就需要在邮件模板的HTML中假如如下的内容来获取钓鱼网址
 ```bash
{{.URL}}
 ```

* `campaigns` 新建一个钓鱼攻击活动，需要输入上文编写的邮件模版名字,钓鱼页面，钓鱼网址，发送日期，配置文件，使用的组。这里需要注意：还是需要在`longing page` 页面进行设置，否则无法正常发送页面的, (重点）URL:
URL 是用来替换选定钓鱼邮件模板中超链接的值 ，这里的URL需要填写当前运行gophish脚本主机的ip。因为启动gophish后，gophish默认监听了3333和80端口，其中3333端口是后台管理系统，而80端口就是用来部署钓鱼页面的。
当URL填写了http://主机IP/，并成功创建了当前的钓鱼事件后。gophish会在主机的80端口部署当前钓鱼事件所选定的钓鱼页面，并在发送的钓鱼邮件里，将其中所有的超链接都替换成部署在80端口的钓鱼页面的url。所以，这里的URL填写我本地当前运行gophish主机的IP对应的url，即`http://192.168.141.1` , 再就是倒是第三行的 `Send Emails By ` 这里可以进行 时间延后发送，这样发送的邮件就是邮件数量除以这段时间的发送，避免了被封的危险。最后的结果在展示面板进行观看。


### 使用EwoMail搭建邮件服务器

要是想整点真实的，可以考虑自己购买域名+国外服务搭建邮箱（国内的服务器不行，25端口不给开）进行测试 ， 这是我的搭建环境：

* 购买服务器需要开启25端口（阿里腾讯等云服务器厂商一般都不给开）
* 需要纯净centos7/8服务器
* centos7/8服务器需要关闭默认防火墙

**命令如下**：

```bash
systemctl status firewalld      // 查看防火墙设置
systemctl stop firewalld        // 关闭防火墙
systemctl disable firewalld     // 关闭防火墙
#mail.symail.club 换成 mail.你的域名
hostnamectl set-hostname mail.polarnightsec.com
#修改hosts
vi /etc/hosts
#把symail.club 全换成你的域名,在下加入这一行, esc+:wq 进行保存
127.0.0.1 mail.polarnightsec.com polarnightsec.com smtp.polarnightsec.com imap.polarnightsec.com
# 安装git
yum -y install git
# git下载EwoMail
git clone https://github.com/gyxuehu/EwoMail.git
cd EwoMail/install
#下面的symail.club请换成你注册的域名
#要是外网服务器 ，设置en参数解析
#sh ./start.sh polarnightsec.com en
sh ./start.sh polarnightsec.com
```
**域名解析**

参考`EwoMail`官方配置解析方式: `v=spf1 ip4:45.32.212.223 -all`

**搭建进入**

默认账号`admin`密码`ewomail123`

```bash
管理端：http://45.32.212.223:8010/Center/Index/login
邮件端：http://45.32.212.223:8000/
```



