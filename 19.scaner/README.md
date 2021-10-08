## 原理

简单模仿 `crawlergo_x_XRAY` 使用golang写的一个 demo, 功能原理一致




## crawlergo + xray

1. xray 设置被动扫描

```
./xray_darwin_amd64 webscan --listen 0.0.0.0:7777 --html-output xray-testphpaabbcc.html
```

2. 启动 crawlergo 扫描

```golang
go run start.go
```


## 测试效果



![image](https://github.com/shangzeng/GolangTools/blob/master/19.scaner/WeChat93e0ff90a19f177ecf95f2c86787d8f5.png)
