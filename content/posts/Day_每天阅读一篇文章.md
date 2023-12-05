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
linkToMarkdown: false
# 上面一般不用动
title: "每天阅读一篇文章"
date: 2023-11-23T00:00:00+08:00
lastmod: 2023-11-23T00:00:00+08:00
categories: ["Day"]
tags: []
---

从 2023.11.12 起，每天观看一篇文章（周五六日除外，复盘），文章大部分来源微信公众号。因为我关注了好多编程公众号，收藏了好多文章，但是好多都没仔细看过。之前在学校的时候，会在吃饭的时候慢慢看，但是现在在公司，吃饭跟同事吃，都没看哈哈。所以攒了好多哈哈。

**目的**有：

1. 逼迫自己每天学习一点知识，细水长流    
2. 看的时候也复习下涉及的知识  
3. 文章不同于书籍，会较口语化，也会更加注重实战

## 2023

### 11 月

### 12 月

#### 12.06

#### 12.05

#### 12.04

#### 11.30

#### 11.29

#### 11.28

#### 11.27

#### 11.24

#### 11.23 视频直播如何工作

URL：<a href="https://mp.weixin.qq.com/s/YqeJvwBtEf0UKr3EzrJ2iQ" target="_blank">https://mp.weixin.qq.com/s/YqeJvwBtEf0UKr3EzrJ2iQ</a>    


#### 11.22 SQL 索引优化

URL：<a href="https://mp.weixin.qq.com/s/sAxUb9ho6eYxjrwNkRENJA" target="_blank">https://mp.weixin.qq.com/s/sAxUb9ho6eYxjrwNkRENJA</a>    

总结：

#### 11.21 百亿数据索引和更新

URL：<a href="https://mp.weixin.qq.com/s/dOH7bQc2CsP2ZjiC76ZEpQ" target="_blank">https://mp.weixin.qq.com/s/dOH7bQc2CsP2ZjiC76ZEpQ</a>    

总结：讲解了使用 MySQL 同步到 ES 为了更好的快速、复杂查询的方案

{{< image src="/images/每天一篇文章/百亿数据索引和更新.png" width=100% height=100% caption="百亿数据索引和更新" >}}

但是，文章只讲解了大概思路，但是具体实现没讲，受限于篇幅的原因，就没讲把：

1. ES 索引全量更新指的什么？如何更新？  
2. MySQL 宽表是从订阅 MySQL binlog，写 MQ，消费 MQ 写入宽表数据的，那如何保证宽表和 MySQL 原表的一致性？

#### 11.20 跨域问题

URL 1：<a href="https://mp.weixin.qq.com/s/PSViiWNSZiq2Tu04D11Lmw" target="_blank">https://mp.weixin.qq.com/s/PSViiWNSZiq2Tu04D11Lmw</a>    
URL 2：<a href="https://blog.csdn.net/weixin_42318691/article/details/121187785" target="_blank">https://blog.csdn.net/weixin_42318691/article/details/121187785</a>    

总结：主要讲解了跨域问题的本质：浏览器的**同源策略**（三元组：协议、域名和端口），同时给出了跨域问题的几种解决方案，服务端的话主要关注：CORS（跨域资源共享, CORS 是Cross-Origin Resource Sharing）和 nginx 域名代理。第 2 篇写得详细一点，还讲解了简单请求和预请求的区别，预请求可以防止服务端执行业务代码

自己也验证了下：

在本地起一个 web 服务器：

```go
package main

import (
	"log"
	"net/http"
)

func main() {
	http.Handle("/test", http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {
		// 尝试注释和不注释的区别
        //resp.Header().Set("Access-Control-Allow-Origin", "*")  // 这是允许访问所有域
		log.Println("访问", req)
		_, _ = resp.Write([]byte("hello"))
	}))
	_ = http.ListenAndServe("127.0.0.1:8080", nil)
}
```

随便在一个网站（例如本网站）的控制台里运行：

```js
const url = 'http://127.0.0.1:8080/test'; 

fetch(url, {
method: 'GET',
})
.then(response => {
console.log("成功"); 
})
.catch(error => {
});
```

可知：浏览器对于跨域的请求都会自动带上 origin 请求头，**同源策略是在响应报文返回时浏览器拦截的**，除非服务端加上 `Access-Control-Allow-Origin", "*"` 代表允许跨域请求，浏览器才会放行

{{< image src="/images/每天一篇文章/cors.png" width=100% height=100% caption="cors" >}}

#### 11.16 MySQL 大数据量快速插入

URL 1：<a href="https://mp.weixin.qq.com/s/mwXOAsmh2RQaTTR-57VXsw" target="_blank">https://mp.weixin.qq.com/s/mwXOAsmh2RQaTTR-57VXsw</a>  
URL 2：<a href="https://mp.weixin.qq.com/s/DieIcM-BNHkL9aT6JU7Epw" target="_blank">https://mp.weixin.qq.com/s/DieIcM-BNHkL9aT6JU7Epw</a>   
URL 3：<a href="https://mp.weixin.qq.com/s/mWL7eacJyYg-6evQZmBVWQ" target="_blank">https://mp.weixin.qq.com/s/mWL7eacJyYg-6evQZmBVWQ</a>    

总结：主要讲解了 MySQL 10 亿数据量如何快速插入的问题

{{< image src="/images/每天一篇文章/MySQL 10 亿数据插入.png" width=100% height=100% caption="MySQL 10 亿数据插入" >}}

{{< image src="/images/每天一篇文章/2 kw 数据？.png" width=100% height=100% caption="2 kw 数据" >}}

#### 11.15 Go GC 调优

URL：<a href="https://mp.weixin.qq.com/s/xDb4PFmd3eDx3pbxgwtd7A" target="_blank">https://mp.weixin.qq.com/s/xDb4PFmd3eDx3pbxgwtd7A</a>

总结：主要讲解了由于响应时长毛刺问题引出 GC 调优。本质问题：GC 阈值设置太低，比上一次 GC 后活跃内存增加 100% 即 GC，导致 GC 频繁，进而导致 STW 时间较长，产生请求响应时间毛刺。最后通过设置百分比为 160% 和设置最大阈值为 1600 MB 降低 GC 频率

#### 11.14 网盘/短链系统设计

URL：<a href="https://mp.weixin.qq.com/s/Siz3YHxsobRIbZC1JYKMhw" target="_blank">https://mp.weixin.qq.com/s/Siz3YHxsobRIbZC1JYKMhw</a>

总结：主要讲解了网盘系统的简要设计。主要学到大文件的元数据（关系型数据库）和文件内容（对象存储服务器）分开存储，大文件分 block，可实现并行传输、断点续传

URL：<a href="https://mp.weixin.qq.com/s/ySA-RbJcC5iVIKqL2GZMaA" target="_blank">https://mp.weixin.qq.com/s/ySA-RbJcC5iVIKqL2GZMaA</a>

总结：主要讲解短链系统的简要设计。没我设计得好哈哈，下面图是我设计的哈哈

{{< image src="/images/每天一篇文章/短链系统.png" width=100% height=100% caption="短链系统" >}}

#### 11.13 微服务优/缺点

URL：<a href="https://mp.weixin.qq.com/s/P-DseTGlD2IA1KQVrKWlPg" target="_blank">https://mp.weixin.qq.com/s/P-DseTGlD2IA1KQVrKWlPg</a>

总结：主要讲解了微服务的优缺点和适用场景

{{< image src="/images/每天一篇文章/微服务的优缺点.png" width=100% height=100% caption="微服务优缺点" >}}

## End
