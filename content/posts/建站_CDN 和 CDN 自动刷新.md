---
# 主页简介
# summary: ""
# 文章副标题
# subtitle: ""
# 作者信息
# author: ""
# authorLink: ""
# authorEmail: ""
# description: ""
# keywords: ""
# license: ""
# images: []
# 文章的特色图片
# featuredImage: ""
# 用在主页预览的文章特色图片
# featuredImagePreview: ""
# password:加密页面内容的密码，详见 主题文档 - 内容加密
# message:  加密提示信息，详见 主题文档 - 内容加密
# linkToMarkdown: true
# 上面一般不用动
title: "CDN 和 CDN 自动刷新"
date: 2022-11-10T01:12:27+08:00
lastmod: 2022-11-10T01:12:27+08:00
categories: ["建站"]
tags: ["CDN"]
---

我的云服务是在广州，如果西藏或者北京的朋友来访问需要通过几千公里的网线传输，速度会很慢。解决方法是 **CDN 加速**。该文章讲解 CDN 的简单工作方式和部署。


## CDN 加速

> CDN 本质是**缓存**。把源服务的静态资源 copy 缓存在全国各地的CDN服务器中，客户端访问时选择最近的CDN服务器请求资源即可，不用请求远处的源服务器。

{{< image src="/images/搭建网站流程/cdn3.png" width=100% height=100% caption="CDN加速" >}}  

**为我们的网站设置CDN**  
很简单，只需要在 CDN 服务厂商购买CDN服务，然后填写网站的相关信息即可。最后记得把域名 DNS 解析改成 **CDN 服务厂商的域名**，且暂停域名解析到 服务器 IP。  
注意：这样做也就意味着以后访问 ayang.ink 域名都会先访问 CDN 服务器，如果没有缓存或者是动态请求（eg：post）再去请求源服务器。

{{< image src="/images/搭建网站流程/cdn1.png" width=100% height=100% caption="CDN加速" >}}  


同样可以通过`ping ayang.ink`来查看是否成功。如成功，访问的是CDN服务器，而不是源服务器。

```bash
ping ayang.ink

正在 Ping fwp74c6f.slt.sched.tdnsv8.com [113.105.165.82] 具有 32 字节的数据:
来自 113.105.165.82 的回复: 字节=32 时间=24ms TTL=52
来自 113.105.165.82 的回复: 字节=32 时间=23ms TTL=52
来自 113.105.165.82 的回复: 字节=32 时间=21ms TTL=52
来自 113.105.165.82 的回复: 字节=32 时间=22ms TTL=52
```
## 缓存更新问题

CDN 并不都是好处连连，当网站的内容（存储在源服务器中）更改时，在 CDN 中的缓存就是脏数据。如果用户此时访问网址，得到的还是 CDN 缓存的脏数据。所以需要我们**主动更新** CDN 中的缓存。

腾讯云服务器提供了缓存管理的三种方式：（以下截取腾讯 CDN 文档） 

>**URL 刷新**：删除 CDN 所有节点上对应资源的缓存。  
>**刷新目录**：有两种模式。选择 “刷新变更资源” 模式，当用户访问匹配目录下资源时，会回源获取资源的 Last-Modify 信息，若与当前缓存资源一致，则直接返回已缓存资源，若不一致，回源拉取资源并重新缓存；选择 “刷新全部资源” 时，当用户访问匹配目录下资源时，直接回源拉取新资源返回给用户，并重新缓存新资源。  
>**URL 预热**：可将指定资源主动从源站加载至 CDN 加速节点并缓存。当用户首次请求资源时，可直接从 CDN 加速节点获取缓存的资源，无需再次回源。

我们可以通过腾讯云 CDN 服务页面直接输入 URL 直接进行 CDN 缓存管理：

{{< image src="/images/CDN 和 CDN 自动刷新/CDN 手动更新缓存.png" width=100% height=100% caption="CDN 手动更新缓存" >}} 

## 自动化刷新 CDN

但是每一次修改文章后都需要手动修改，效率低下。其实腾讯云为上面三种方式都提供了相应的 API 调用。<a href="https://cloud.tencent.com/document/product/228/37870" target="_blank">API 文档</a>  

至此，我们可以把自动刷新 CDN 集合到自动化部署网站中。对网站内容的修改或增加，`git push` 后，利用 GitHUb WebHook 提供的信息和 CDN 提供的 API，我们可以做到自动化刷新 CDN 缓存。此外，对于网站新增加内容，我们进行缓存预热。

```json
  // 截取 GitHUb WebHook 发送 post 请求的部分信息
  "head_commit": {
    // commit -m "update post"，commit 备注的信息
    "message": "update post",
    "added": [

    ],
    // git push 更改的文件
    "modified": [
      "content/posts/底层_通过 go 汇编了解函数调用栈帧.md"
    ]
  }
```

tip：对 json 数据的处理，可用第三方库 <a href="github.com/tidwall/gjson" target="_blank">github.com/tidwall/gjson</a> 简化开发。

```bash
# 项目目录如下
├─cdn
│  ├─.idea
│  ├─cdn
│  └─deploy
```

```go
// 对腾讯提供的 API 的进一步封装
// cdn/cdn.go
package cdn

import (
	"fmt"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"os"
	"sync"
)

const (
	POST     = "https://ayang.ink/posts/"
	CATEGORY = "https://ayang.ink/categories/"
	TAG      = "https://ayang.ink/tags/"
)

var (
	// 文章跟新时，需要更改其他页面
	updatePath = []string{POST, CATEGORY, TAG}
)

type Cdn interface {
	Add(urlEncode bool, urls ...string)
	Update(urlEncode bool, urls ...string)
}

type Client struct {
	once   sync.Once
	client *cdn.Client
}

func (c *Client) initClient() {
	if c == nil {
		panic("未初始化 Client")
	}
	// 实例化一个认证对象，入参需要传入腾讯云账户secretId，secretKey,此处还需注意密钥对的保密
	// 密钥可前往https://console.cloud.tencent.com/cam/capi网站进行获取
	secretKey := os.Getenv("CDN_SECRET_KEY")
	credential := common.NewCredential(
		"AKID4xjHckecMu9D4g4FEzWh4E9qNglbDbqs",
		secretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cdn.tencentcloudapi.com"
	//实例化要请求产品的client对象,clientProfile是可选的
	client, _ := cdn.NewClient(credential, "", cpf)
	c.client = client
}

// Add 网站 post 文章
func (c *Client) Add(urlEncode bool, urls ...string) {
	c.once.Do(c.initClient)
	// 缓存预热
	c.pushUrlsCache(urls, urlEncode)
	// 刷新目录，包括 TAG 等
	c.purgePathCache(updatePath, urlEncode)
}

// Update 更该文章
func (c *Client) Update(urlEncode bool, urls ...string) {
	c.once.Do(c.initClient)
	// 只刷新这篇文章
	c.purgePathCache(urls, urlEncode)
}

// purgePathCache 刷新目录
func (c *Client) purgePathCache(urls []string, urlEncode bool) {
	urlsPoint := strToStrPoint(urls)
	boolPoint := new(bool)
	if urlEncode {
		*boolPoint = true
	} else {
		*boolPoint = false
	}
	str := new(string)
	*str = "delete"
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := cdn.NewPurgePathCacheRequest()
	request.Paths = append(request.Paths, urlsPoint...)
	request.UrlEncode = boolPoint
	request.FlushType = str

	// 返回的resp是一个PurgePathCacheResponse的实例，与请求对象对应
	response, err := c.client.PurgePathCache(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		panic(err)
	}
	// 输出json格式的字符串回包
	fmt.Printf("%s", response.ToJsonString())
}

// pushUrlsCache 预热 url
func (c *Client) pushUrlsCache(urls []string, urlEncode bool) {
	urlsPoint := strToStrPoint(urls)
	boolPoint := new(bool)
	if urlEncode {
		*boolPoint = true
	} else {
		*boolPoint = false
	}

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := cdn.NewPushUrlsCacheRequest()

	request.Urls = append(request.Urls, urlsPoint...)
	request.UrlEncode = boolPoint

	// 返回的resp是一个PurgePathCacheResponse的实例，与请求对象对应
	response, err := c.client.PushUrlsCache(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		fmt.Printf("An API error has returned: %s", err)
		return
	}
	if err != nil {
		panic(err)
	}
	// 输出json格式的字符串回包
	fmt.Printf("%s", response.ToJsonString())
}

// strToStrPoint 不知道为啥腾讯云全是指针
func strToStrPoint(str []string) (strPoint []*string) {
	strPoint = make([]*string, len(str))
	for i, _ := range str {
		strPoint[i] = &str[i]
	}
	return
}
```

```go
// 自动化部署，集成刷新、预热 CDN 缓存
// deploy/deploy.go
package deploy

import (
	"autodeploy/cdn"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"github.com/tidwall/gjson"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const (
	UPDATE = "update"
	ADD    = "add"
	MD     = ".md"
	HOST   = "https://ayang.ink"
)

var (
	CDN cdn.Cdn = &cdn.Client{}
)

func AutoDeploy(w http.ResponseWriter, req *http.Request) {
	body, _ := ioutil.ReadAll(req.Body)
	// 解密错误，why
	//if !verify(body, req) {
	//	w.WriteHeader(400)
	//	return
	//} else {
	//	w.WriteHeader(200)
	//}
	w.WriteHeader(200)
	command := "/home/ayang/site/Ayang-Blog/autodeploy/autodeploy.sh"
	cmd := exec.Command("/bin/bash", command)
	// 执行 shell 脚本文件
	if err := cmd.Run(); err != nil {
		// 缓存刷新
		cdnFunc(body)
	}
}

// false，why
func verify(body []byte, req *http.Request) bool {
	if req.Method == http.MethodPost {
		// 获取环境变量token
		token := os.Getenv("SECRET_TOKEN_GITHUB_AUTO_DEPLOY")
		// 获取加密
		head := req.Header.Get("X-Hub-Signature-256")
		// hmac是Hash-based Message Authentication Code的简写，就是指哈希消息认证码
		h := hmac.New(sha256.New, []byte(token))
		h.Write(body)
		signature := "sha256=" + hex.EncodeToString(h.Sum(nil))
		if hmac.Equal([]byte(signature), []byte(head)) {
			return true
		}
	}
	return false
}

// cdn 的刷新、预热等操作
func cdnFunc(body []byte) {
	message := gjson.GetBytes(body, "head_commit.message").String()
	if strings.Contains(message, UPDATE) {
		update(body)
	} else if strings.Contains(message, ADD) {
		add(body)
	}
}

func update(body []byte) {
	modified := gjson.GetBytes(body, "head_commit.modified").Array()
	var strs []string
	for _, m := range modified {
		if str := judgeMD(m.String()); str != "" {
			str = HOST + str
			strs = append(strs, str)
		}
	}
	if len(strs) > 0 {
		CDN.Update(true, strs...)
	}
}

func add(body []byte) {
	modified := gjson.GetBytes(body, "head_commit.added").Array()
	var strs []string
	for _, m := range modified {
		if str := judgeMD(m.String()); str != "" {
			str = HOST + str
			strs = append(strs, str)
		}
	}
	if len(strs) > 0 {
		CDN.Add(true, strs...)
	}
}

// judgeMD 判断后缀是否未 “.MD”，并返回 /目录/
func judgeMD(str string) string {
	if strings.HasSuffix(str, MD) {
		return cut(str) + "/"
	} else {
		return ""
	}
}

// 截取有效目录，去掉“.MD”后缀，把“ ”替换成“-”，并把目录中大写改成小写，返回 /目录
func cut(str string) string {
	start := strings.LastIndex(str, "/")
	str = str[start : len(str)-3]
	str = strings.ReplaceAll(str, " ", "-")
    str = strings.ToLower(str)
	return str
}
```

```go
package main

import (
	"autodeploy/deploy"
	"net/http"
)

func main() {
	http.HandleFunc("/autodeploy", deploy.AutoDeploy)
	_ = http.ListenAndServe(":1314", nil)
}
```

## end