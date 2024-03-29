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
title: "cron job：阅读文章"
date: 2024-02-13T00:00:00+08:00
lastmod: 2024-02-13T00:00:00+08:00
categories: ["Day"]
tags: []
---

从 2023.11.12 起，每 3 天阅读分享一篇文章，文章大部分来源微信公众号。因为我关注了好多编程公众号，收藏了好多文章，但是好多都没仔细看过。之前在学校的时候，会在吃饭的时候慢慢看，但是现在在公司，吃饭跟同事吃，都没看哈哈。所以攒了好多哈哈。

**目的**有：

1. push 自己周期性学习知识，细水长流    
2. 看的时候也复习下涉及的知识  
3. 文章不同于书籍，会较口语化，也会更加注重实战  
4. 留下记录，定时回顾复习  

PS：从微信公众号，点击“复制链接”，都是短链的方式，不知道这种**短链的映射时间是否永久**~~~

## 推荐的订阅

架构师合集：<a href="https://mp.weixin.qq.com/mp/appmsgalbum?__biz=MzkzMDI0ODg4NQ==&action=getalbum&album_id=2247053463681564673" target="_blank">https://mp.weixin.qq.com/mp/appmsgalbum?__biz=MzkzMDI0ODg4NQ==&action=getalbum&album_id=2247053463681564673</a>   

## 2024

### 2 月

#### 02.28

#### 02.25

#### 02.22

#### 02.19

#### 02.13 热点库存扣减方案

URL：<a href="https://mp.weixin.qq.com/s/4XQAdXEVIjXnicauPE-yMw" target="_blank">https://mp.weixin.qq.com/s/4XQAdXEVIjXnicauPE-yMw</a>

总结：直接看下图把

{{< image src="/images/每天一篇文章/热点库存扣减.png" width=100% height=100% caption="热点库存扣减" >}} 

当然，这里没考虑 Redis 有主从的情况，如果主奔溃，发生主从替换，由于 Redis 异步复制的特性，会造成 Incr 丢失，造成库存超卖的问题  
当然，这种情况不用主从把，如果发生奔溃，重新开启一个 Redis 重新申请库存即可，奔溃掉的 Redis 的剩余库存直接丢失，虽然会造成库存少卖现象，但毕竟奔溃是很少发生的

之前笔者做的抽奖项目里的库存扣减就有考虑主从替换导致的 Redis 库存回滚问题，具体做法是：

```go
curStock = redis.incr($stock)
if curStock > sumStock {
   return 库存不足
} else {
  // 另一个 Redis，对每个库存上一个锁，过期时间为活动结束。这里上了一层保障，只要两个 Redis 没有全部奔溃，就不会造成库存回滚
   if redis.setnx(curStock, 过期时间为活动结束) {
      return 扣减库存成功
   } else {
      return 上锁失败
   }
}
```

#### 02.10 异地多活

URL：<a href="https://mp.weixin.qq.com/s/T6mMDdtTfBuIiEowCpqu6Q" target="_blank">https://mp.weixin.qq.com/s/T6mMDdtTfBuIiEowCpqu6Q</a>

总结：开拓下视野把

一个好的软件架构，应该遵循高性能、高可用、易扩展 3 大原则，其中提升高可用的核心是「冗余」，备份、主从副本、同城灾备、同城双活、两地三中心、异地双活，异地多活都是在做冗余。同城灾备分为「冷备」和「热备」，冷备只备份数据，不提供服务，热备实时同步数据，并做好随时切换的准备，即已经部署相应的服务，只是正常时不对外暴露提供服务。同城双活比灾备的优势在于，主机房可以接入「写」流量，两个机房都可以接入「读写」流量，提高可用性的同时，还提升了系统性能。异地双活，两地机房都可以接入「读写」流量，所以对于异地双活或多活，必须做一个良好的业务路由隔离，防止写入时对同一条数据产生并发而导致不一致的结果，这应该就是最难的点。当然对于库存等强一致性要求的数据，肯定不适合异地双活或多活的，这类服务依旧只能采用**只写主机房**，读从机房的方案

#### 02.07 异构存储

URL 1：<a href="https://mp.weixin.qq.com/s/Se_WkDhAls4PA8DQ5jqmvQ" target="_blank">https://mp.weixin.qq.com/s/Se_WkDhAls4PA8DQ5jqmvQ</a>    
URL 2：<a href="https://mp.weixin.qq.com/s/4x7WoOMmhLrVfFpgKKrgZQ" target="_blank">https://mp.weixin.qq.com/s/4x7WoOMmhLrVfFpgKKrgZQ</a>   
URL 3：<a href="https://mp.weixin.qq.com/s/vn0JTD7Rq_8PNdBIbk3IhA" target="_blank">https://mp.weixin.qq.com/s/vn0JTD7Rq_8PNdBIbk3IhA</a>    

总结：开拓视野把，实际上就我实习来说，非常多的场景都需要异构数据，不同数据库提供了不同场景下的能力

{{< image src="/images/每天一篇文章/异构存储.png" width=100% height=100% caption="异构存储" >}} 

#### 02.04 淘宝购物车扩容

URL：<a href="https://mp.weixin.qq.com/s/CAUZIIOxsr6kHZcgHvtrcg" target="_blank">https://mp.weixin.qq.com/s/CAUZIIOxsr6kHZcgHvtrcg</a>

#### 02.01 Raft 简学

URL：<a href="https://learn.lianglianglee.com/%e4%b8%93%e6%a0%8f/%e5%88%86%e5%b8%83%e5%bc%8f%e4%b8%ad%e9%97%b4%e4%bb%b6%e5%ae%9e%e8%b7%b5%e4%b9%8b%e8%b7%af%ef%bc%88%e5%ae%8c%ef%bc%89/09%20%e5%88%86%e5%b8%83%e5%bc%8f%e4%b8%80%e8%87%b4%e6%80%a7%e7%ae%97%e6%b3%95%20Raft%20%e5%92%8c%20Etcd%20%e5%8e%9f%e7%90%86%e8%a7%a3%e6%9e%90.md" target="_blank">https://learn.lianglianglee.com/%e4%b8%93%e6%a0%8f/%e5%88%86%e5%b8%83%e5%bc%8f%e4%b8%ad%e9%97%b4%e4%bb%b6%e5%ae%9e%e8%b7%b5%e4%b9%8b%e8%b7%af%ef%bc%88%e5%ae%8c%ef%bc%89/09%20%e5%88%86%e5%b8%83%e5%bc%8f%e4%b8%80%e8%87%b4%e6%80%a7%e7%ae%97%e6%b3%95%20Raft%20%e5%92%8c%20Etcd%20%e5%8e%9f%e7%90%86%e8%a7%a3%e6%9e%90.md</a>

总结：讲解了 Raft 三个子问题，清晰地列举步骤一、二等。虽然每个步骤没有往深处讲，但真的非常清晰。建议搭配 B 站 UP 主“戌米的论文笔记”的讲解 Raft 的视频食用更佳。

{{< image src="/images/每天一篇文章/Raft1.png" width=100% height=100% caption="Leader 选举" >}}

{{< image src="/images/每天一篇文章/Raft2.png" width=100% height=100% caption="日志复制" >}}

{{< image src="/images/每天一篇文章/Raft3.png" width=100% height=100% caption="日志存储格式" >}}

## 2023

### 12 月

#### 12.31 ETCD 分布式公平锁

URL：<a href="https://learn.lianglianglee.com/%E4%B8%93%E6%A0%8F/%E5%88%86%E5%B8%83%E5%BC%8F%E4%B8%AD%E9%97%B4%E4%BB%B6%E5%AE%9E%E8%B7%B5%E4%B9%8B%E8%B7%AF%EF%BC%88%E5%AE%8C%EF%BC%89/10%20%E5%9F%BA%E4%BA%8E%20Etcd%20%E7%9A%84%E5%88%86%E5%B8%83%E5%BC%8F%E9%94%81%E5%AE%9E%E7%8E%B0%E5%8E%9F%E7%90%86%E5%8F%8A%E6%96%B9%E6%A1%88.md" target="_blank">https://learn.lianglianglee.com/%E4%B8%93%E6%A0%8F/%E5%88%86%E5%B8%83%E5%BC%8F%E4%B8%AD%E9%97%B4%E4%BB%B6%E5%AE%9E%E8%B7%B5%E4%B9%8B%E8%B7%AF%EF%BC%88%E5%AE%8C%EF%BC%89/10%20%E5%9F%BA%E4%BA%8E%20Etcd%20%E7%9A%84%E5%88%86%E5%B8%83%E5%BC%8F%E9%94%81%E5%AE%9E%E7%8E%B0%E5%8E%9F%E7%90%86%E5%8F%8A%E6%96%B9%E6%A1%88.md</a>

总结：ETCD 实现分布公平锁主要依赖于其 prefix 机制和 revision 机制

#### 12.28 分布式锁 

URL：<a href="https://mp.weixin.qq.com/s/yZC6VJGxt1ANZkn0SljZBg" target="_blank">https://mp.weixin.qq.com/s/yZC6VJGxt1ANZkn0SljZBg</a>

总结：很全面的分布式锁各方面的讲解

{{< image src="/images/每天一篇文章/分布式锁.png" width=100% height=100% caption="Redis 分布式锁" >}}

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
