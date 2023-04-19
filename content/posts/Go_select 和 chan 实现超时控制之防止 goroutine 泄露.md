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
title: "select 和 chan 实现超时控制之防止 goroutine 泄露"
date: 2023-04-12T01:22:45+08:00
lastmod: 2023-04-12T01:22:45+08:00
categories: ["Go"]
tags: []
---

前一篇文章说到 http.Client 巧妙利用 chan 进行**协程间通信**和**超时控制**，其实这里面有很多细节，包括需要注意避免**协程泄露**等。  

由于 http.Client 源码太长，我就简单模拟了一下**协程通信**和**超时处理**的场景。

## 问题代码

```go
package main

import "time"

func main() {
	// 这里是同步等待执行完毕，deal 有超时控制
	result := deal()
	print("执行完毕，result 为：", result, "\n")
	time.Sleep(time.Second * 5)
	print("父协程结束\n")
}

func deal() string {
	timeout := time.NewTimer(time.Second * 2)
	// 结束后立刻关闭计时器，防止高并发情况下很多个计时器运行，造成性能损耗
	defer timeout.Stop()

	done := make(chan string)

	result := ""

	// 开启一个协程异步执行任务
	go func() {
		// 模拟处理逻辑花费时间超长（如需要 MySQL 请求等）
		// do.....
		time.Sleep(time.Second * 3)
		ans := "ayang"
		// 问题所在（bug）：如果超时了，那么下面会走 timeout case
		// 而如果一段时间后上面的逻辑处理完毕，执行下面代码时发现 done 已经没有协程在等待
		// 将会永久阻塞，造成 goroutine 永远不会回收，也就是协程泄露
		done <- ans
		print("子协程正常结束\n")
	}()

	// select 等待子协程处理完毕或超时直接返回
	// 问：既然异步了，为什么这里还要阻塞等待执行完毕？不是还是又回到同步么？
	// 确实是同步的，不过同步的只有最多 timeout 的时间，过了这个时间会立刻返回（那个子协程此时就是异步执行了）
	// 如果不采取子协程异步发送，那么可能会一直阻塞在该方法很久
	// 也就是说实现超时 立刻返回 功能必须是 异步 的

	select {
	case <-timeout.C:
		result = "nil"
	case result = <-done:
	}

	return result
}
```

## 解决方法

解决方法也是比较简单，直接看里面注释把。

```go
package main

import "time"

func main() {
	result := deal()
	print("执行完毕，result 为：", result, "\n")
	time.Sleep(time.Second * 5)
	print("父协程结束\n")
}

func deal() string {
	timeout := time.NewTimer(time.Second * 2)
	defer timeout.Stop()

	done := make(chan string)

	result := ""

	go func() {
		time.Sleep(time.Second * 3)
		ans := "ayang"

		// 解决办法：非阻塞写入。
		// 因为如果超时了，那么外面已经没有协程阻塞等待接受了，
		select {
		case done <- ans:
		default:
		}
		print("子协程正常结束\n")
	}()

	select {
	case <-timeout.C:
		result = "nil"
	case result = <-done:
	}

	return result
}
```

## End
