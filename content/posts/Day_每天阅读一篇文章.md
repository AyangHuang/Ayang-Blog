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
title: "每三天阅读一篇文章"
date: 2023-12-13T00:00:00+08:00
lastmod: 2023-12-13T00:00:00+08:00
categories: ["Day"]
tags: []
---

从 2023.11.12 起，3 天分享一篇文章，文章大部分来源微信公众号。因为我关注了好多编程公众号，收藏了好多文章，但是好多都没仔细看过。之前在学校的时候，会在吃饭的时候慢慢看，但是现在在公司，吃饭跟同事吃，都没看哈哈。所以攒了好多哈哈。

**目的**有：

1. 逼迫自己每天学习一点知识，细水长流    
2. 看的时候也复习下涉及的知识  
3. 文章不同于书籍，会较口语化，也会更加注重实战

PS：从微信公众号，点击“复制链接”，都是短链的方式，不知道这种**短链的映射时间是否永久**~~~

## 推荐的订阅

架构师合集：<a href="https://mp.weixin.qq.com/mp/appmsgalbum?__biz=MzkzMDI0ODg4NQ==&action=getalbum&album_id=2247053463681564673" target="_blank">https://mp.weixin.qq.com/mp/appmsgalbum?__biz=MzkzMDI0ODg4NQ==&action=getalbum&album_id=2247053463681564673</a>   

## 2024

### 1 月

#### 01.01

#### 01.04

## 2023

### 12 月

#### 12.31

#### 12.28

#### 12.25 什么是云原生

URL：<a href="https://mp.weixin.qq.com/s/fsO4Zh1spalmrIwxtn8fvw" target="_blank">https://mp.weixin.qq.com/s/fsO4Zh1spalmrIwxtn8fvw</a>

总结：该文章通过一张图概括了软件架构和流程的演变。概括了云原生的 4 点：

1. 开发流程：DevOps，开发、测试、部署、运维紧密相连  
2. 应用架构：从单体向微服务转变，每一个微服务很小  
3. 部署和打包：打包成镜像并部署在容器中（区别于直接部署在物理机上），例如 docker  
4. 基础设施：部署在云基础设施上，不管是公有云还是私有云，可以动态扩容和缩容  

{{< image src="/images/每天一篇文章/什么是云原生.png" width=100% height=100% caption="什么是云原生" >}}

#### 12.22 职场新人-如何快速变得专业

URL 1：<a href="https://mp.weixin.qq.com/s/TFCrRQT3k5P2Pz3YOWfMKw" target="_blank">https://mp.weixin.qq.com/s/TFCrRQT3k5P2Pz3YOWfMKw</a>   
URL 2：<a href="https://mp.weixin.qq.com/s/LIpXt9aCHGqrV1s0jddMnw" target="_blank">https://mp.weixin.qq.com/s/LIpXt9aCHGqrV1s0jddMnw</a>   

总结：首先讲了新人的定义：**不熟练、不系统、不严谨、不开放**，然后讲解了如何通过**快速变得熟练、能够系统化思考、以严谨的态度、开放的心态去展开工作**快速变得专业

* 快速变得熟练：快速了解并学会使用你日常工作需要的工具
  * 研发类工具（git，idea 插件）
  * 运维类工具（日志查询平台、系统监控平台）
  * 泛文档类工具
* 能够系统化思考
  * 提升思考全面性  
    团队的文档模板（例如技术文档中通常会要求你考虑一些工程方面的内容例如“变更风险评估”、“上下游影响面分析”、“安全评估”、“容量评估”、“切流设计”等等）
  * 提升内容逻辑性  
    目前存在的问题->产生问题的原因->已有的方法能否复用解决，新的解决方法->会不会产生新的影响
* 以严谨的态度
  * 技术方案考虑的是长期性、可维护性
* 以开放的心态
  * 勇于承担事情，要主动揽不重复的、自己不熟悉（能够促进自己学习）的活
  * 不要害怕犯错

最近在字节实习了快三个月，这两篇文章看下来，感触颇多。能够总结下这两篇文章，真的很厉害。反思下自己，或多或少上面的问题都有存在着。例如  
1. 能够系统化思考：
   1. 越底层的接口应该考虑长期性、通用性
   2. 当问 mentor 说 api 层为什么没有 swagger 或者内部的 bam 时，回复是旧项目，框架不支持。但是仅仅到这里就结束了，没有想为什么不支持以及如何去支持
2. 以开放的心态：
   1. 自己只等着 mentor 发放需求，没有主动揽活
   2. 害怕犯错，很多时候都是一对一联系，而不是群里同步联系
   
希望下次实习或者正式工作，对于以上四点，能够有所成长把~

#### 12.19 大厂二面重点

URL：<a href="https://mp.weixin.qq.com/s/8IzKAWNqqaCE9xL6nqFClQ" target="_blank">https://mp.weixin.qq.com/s/8IzKAWNqqaCE9xL6nqFClQ</a>  

总结：讲解了大厂二面的重点：不能在于面经层面，而应该是**重视实践，问题拆解，重视细节**

#### 12.16 架构师：提升系统稳定性

URL 1：<a href="https://mp.weixin.qq.com/s/2sbvHavhqQw5QQilHdSY2w" target="_blank">https://mp.weixin.qq.com/s/2sbvHavhqQw5QQilHdSY2w</a>   
URL 2：<a href="https://mp.weixin.qq.com/s/3l0XvssIbYVhhtJa0qCyeA" target="_blank">https://mp.weixin.qq.com/s/3l0XvssIbYVhhtJa0qCyeA</a>  

#### 12.13 架构师：提升系统读性能

#### 12.10 架构师：提升系统写性能

#### 12.07 幂等性

URL：<a href="https://mp.weixin.qq.com/s/QufXfnJj5kPX8K3M5gICqw" target="_blank">https://mp.weixin.qq.com/s/QufXfnJj5kPX8K3M5gICqw</a>   

总结：讲解了支付幂等场景的解决方案。（但总感觉缺少很多细节）

下面是我总结了幂等的常见解决方案：

{{< image src="/images/每天一篇文章/请求幂等.png" width=100% height=100% caption="请求幂等" >}}

#### 12.04 视频直播如何工作

URL：<a href="https://mp.weixin.qq.com/s/YqeJvwBtEf0UKr3EzrJ2iQ" target="_blank">https://mp.weixin.qq.com/s/YqeJvwBtEf0UKr3EzrJ2iQ</a>    


#### 12.01 SQL 索引优化

URL：<a href="https://mp.weixin.qq.com/s/sAxUb9ho6eYxjrwNkRENJA" target="_blank">https://mp.weixin.qq.com/s/sAxUb9ho6eYxjrwNkRENJA</a>    

总结：

### 11 月

#### 11.28 百亿数据索引和更新

URL：<a href="https://mp.weixin.qq.com/s/dOH7bQc2CsP2ZjiC76ZEpQ" target="_blank">https://mp.weixin.qq.com/s/dOH7bQc2CsP2ZjiC76ZEpQ</a>    

总结：讲解了使用 MySQL 同步到 ES 为了更好的快速、复杂查询的方案

{{< image src="/images/每天一篇文章/百亿数据索引和更新.png" width=100% height=100% caption="百亿数据索引和更新" >}}

但是，文章只讲解了大概思路，但是具体实现没讲，受限于篇幅的原因，就没讲把：

1. ES 索引全量更新指的什么？如何更新？指的是从宽表重新构建整个新的索引？（问了作者，是的，就是重新建立整个 ES 索引） 
2. MySQL 宽表是从订阅 MySQL binlog，写 MQ，消费 MQ 写入宽表数据的，那如何保证宽表和 MySQL 原表的一致性？

#### 11.25 跨域问题

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

#### 11.22 MySQL 大数据量快速插入

URL 1：<a href="https://mp.weixin.qq.com/s/mwXOAsmh2RQaTTR-57VXsw" target="_blank">https://mp.weixin.qq.com/s/mwXOAsmh2RQaTTR-57VXsw</a>  
URL 2：<a href="https://mp.weixin.qq.com/s/DieIcM-BNHkL9aT6JU7Epw" target="_blank">https://mp.weixin.qq.com/s/DieIcM-BNHkL9aT6JU7Epw</a>   
URL 3：<a href="https://mp.weixin.qq.com/s/mWL7eacJyYg-6evQZmBVWQ" target="_blank">https://mp.weixin.qq.com/s/mWL7eacJyYg-6evQZmBVWQ</a>    

总结：主要讲解了 MySQL 10 亿数据量如何快速插入的问题

{{< image src="/images/每天一篇文章/MySQL 10 亿数据插入.png" width=100% height=100% caption="MySQL 10 亿数据插入" >}}

{{< image src="/images/每天一篇文章/2 kw 数据？.png" width=100% height=100% caption="2 kw 数据" >}}

#### 11.19 Go GC 调优

URL：<a href="https://mp.weixin.qq.com/s/xDb4PFmd3eDx3pbxgwtd7A" target="_blank">https://mp.weixin.qq.com/s/xDb4PFmd3eDx3pbxgwtd7A</a>

总结：主要讲解了由于响应时长毛刺问题引出 GC 调优。本质问题：GC 阈值设置太低，比上一次 GC 后活跃内存增加 100% 即 GC，导致 GC 频繁，进而导致 STW 时间较长，产生请求响应时间毛刺。最后通过设置百分比为 160% 和设置最大阈值为 1600 MB 降低 GC 频率

#### 11.16 网盘/短链系统设计

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
