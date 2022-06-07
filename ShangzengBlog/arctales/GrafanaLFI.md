Time: 20211222:golang: 第一次在没有获取到漏洞细节，反推出漏洞存在点，感觉还是很有成就感的...:Grafana 任意文件读取漏洞 CVE-2021-43798
--------


# 背景


payload如下：

```golang
public/plugins/grafana-clock-panel/../../../../../../etc/passwd
```

![](/static/img/WechatIMG2464.png)



# 下载/安装

进入 **pkg/api/api.go**文件夹，寻找路由，跟进 **getPluginAssets** ,换 **Goland** 方便跟进函数，函数**getPluginAssets** 在 **plugins.go** 中，这里**filepath.Join** 直接获取了插件的名字**plugin.PluginDir**和路径**requestedFile**进行了拼接（但是在之前有判断插件是否存在，因此得是存在的插件），然后直接使用 **os.Open** 打开了插件，这时候导致了问题的存在：


* **filepath.Join**   用于拼接URL路径   filepath.Join("a","a").  => "a/a"
* **filepath.Clean**   清理路径中的多余字符 , 也就是这里出现了问题，导致 **../** 进行了拼接

![](/static/img/WechatIMG2465.png)

## 关于 filepath.Clean

在过滤时候，  **/../../../**可以进行处理

![](/static/img/WechatIMG2466.png)

但是 **../../../**  并不会被处理

![](/static/img/WechatIMG2467.png)

使用 `os.Open` 打开文件，使用ServeContent


* 关于**ServeContent** https://pkg.go.dev/net/http#ServeContent

![](/static/img/WechatIMG2468.png)


## 查看修复方式

在最新版本已经修复了，修复方式就是在数组前加上 **/**

```golang
// https://github.com/grafana/grafana/blob/main/pkg/api/plugins.go
requestedFile := filepath.Clean(filepath.Join("/", web.Params(c.Req)["*"]))
rel, err := filepath.Rel("/", requestedFile)
```

# 其他 TIPS

可以通过在前面加上//../,让nginx不处理后的锚点


```golang
/public/plugins/alertlist/#/../..%2f..%2f..%2f..%2f..%2f..%2f..%2f..%2f/etc/hosts
```
