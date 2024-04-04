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
title: "上线时反复 OOM"
date: 2024-04-05T01:12:45+08:00
lastmod: 2024-04-05T01:12:45+08:00
categories: ["实战"]
tags: []
---

最近在 B 站实习，mentor 在上线服务时实例出现 OOM 问题。记录一下

## OOM

现象：出现 OOM 后，容器会自动重启，内存飙升，然后继续 OOM

{{< image src="/images/OOM 问题/现象.png">}}

因为线上出现问题，降低损耗是第一位。mentor 考虑到刚刚进行一次发布上线，怀疑是新代码内存泄露造成的 OOM，立刻进行**回滚**

{{< image src="/images/OOM 问题/线上问题.png">}}

但是回滚后，**问题依旧存在**，陷入重启、OOM、重启的无限循环

由于该服务主要是**消费上游稿件指纹识别结果的 MQ**，然后进行一系列侵权规则匹配的处理。所以 mentor 把消息 MQ 的速率进行降级，即消费 MQ 的速率进行限流，服务重启后逐渐恢复到正常的状态。在此之前，为了更好的排查分析问题，对重启后的容器生成内存的火焰图如下：

{{< image src="/images/OOM 问题/pprof.png">}}

通过火焰图很明显可以推断出消耗大量内存的正是 `loadVideoTortRuleMap` 方法：

```go
type Service struct {
	videoTortRuleMap           map[int64]*dmm.VideoTortRule                           
	videoTortRulePropertyMap   map[int64][]*dmm.VideoTortRuleProperty 
}

// 消费每一条 MQ
func (s *Service) Process() {
    ...
    if len(s.videoTortRuleMap) == 0 {
		s.loadVideoTortRuleMap(ctx)
    }
    ...
}

func (s *Service) loadVideoTortRuleMap(ctx context.Context) {
    // select 读数据库
    s.videoTortRuleMap = 数据库读取的数据 
    videoTortRulePropertyMap = 数据库读取的数据
}
```

其实很容易发现问题，当服务刚刚启动时，`if len(s.videoTortRuleMap) == 0` 判断为 true，会执行 `loadVideoTortRuleMap` 进行读取，把数据库加载到本地内存中。而由于服务器启动时，并发消费 MQ 将会导致 `if len(s.videoTortRuleMap) == 0` 多次判断为 true，并多次执行 `loadVideoTortRuleMap`，该方法会扫描把整个表都读取加载到内存，多次加载过程中，导致了 OOM（规则表有 5 k 行，规则资产表有 18 万行）

既然该代码存在很久，为什么现在才会暴露呢？

mentor 说其实之前上线时实例会偶发 OOM，但重启后都没事。这次由于上线时刚好在跑任务，消费速度加快，并发度高，才暴露了这个问题

## 解决方法

由于本质问题是对并发访问没有做到互斥，所以加个 sync.Once 来达到只执行一次即可

```go
type Service struct {
	videoTortRuleMap           map[int64]*dmm.VideoTortRule                           
	videoTortRulePropertyMap   map[int64][]*dmm.VideoTortRuleProperty 
    once                       sync.Once
}

// 消费每一条 MQ
func (s *Service) Process() {
    ...
    if len(s.videoTortRuleMap) == 0 {
		s.once.Do(s.loadVideoTortRuleMap(ctx))
    }
    ...
}
```

其实这个是个历史遗留问题，应该在服务启动的时候就执行该方法，而不是消费 MQ 时才**懒加载**

## 总结

其实在字节实习的时候也因为**资源的并发懒加载**踩过坑，所以**懒加载**其实是个坑，一般放在 `main` 函数进行饿汉式加载即可，才不会出现并发问题。如果真要懒加载，那么该**加载函数一定要保证是并发安全**的

{{< image src="/images/OOM 问题/线上问题.png">}}

## End
