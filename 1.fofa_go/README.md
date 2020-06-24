# FoFa Go

FoFa Go 是使用 golang 编写的一个调用fofa API 调用工具


## 安装

安装go语言环境  

```
sudo apt-get install golang 	  # (Linux)
brew install go            		  # (Mac)
```

安装第三方库excelize

```
go get github.com/360EntSecGroup-Skylar/excelize
```


编辑 config.json 填写API与email

```
email    =  "xxxxxxxxxxxxxxxxx"
api_key  =  "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
```

执行

```
go run fofa_go_1.1.go -h
```


## 编译/跨平台

在填写API与email后，我们可以编译达到跨平台的目的,其中参数设置可以看[这里](https://golang.org/doc/install/source#environment/)

```
go build fofa_go_1.1.go
```

编译常见linux使用的版本：

```
GOOS="linux" GOARCH="amd64" go build -ldflags "-w -s" fofa_go_1.1.go
```

编译常见windows使用的版本：

```
GOOS="windows" GOARCH="amd64" go build -ldflags "-w -s" fofa_go_1.1.go
```

## 使用说明

自行查看帮助，注意转译





## 待完成

1. win编译后输出字体颜色缺失
2. 目前只能下载前一万条（因为没有币子做实验下载，呜呜呜）
3. 老八一样的代码需要优化
4. 与xray进行联动进行主动扫描


## 更多

[shangzeng's blog ](https://www.shangzeng.club)
















