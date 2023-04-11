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

由于 http.Client 源码太长，我就简单模拟了一下超时处理的场景。

## 问题代码

```go
package main

import "time"

func main() {
	// 这里是同步等待执行完毕，deal 有超时控制
	result := deal()
	print("执行完毕，result 为：", result)
	// 循环，防止 main goroutine 退出
	for {}
}

func deal() string {
	timeout := time.After(time.Second * 2)
	done := make(chan string)

	result := ""

	// 开启一个协程异步执行任务
	go func() {
		// 模拟处理逻辑花费时间超长（如需要 MySQL 请求等）
        do.....
		time.Sleep(time.Second * 3)
		ans := "ayang"
		// 问题所在（bug）：如果超时了，那么下面会走 timeout case
		// 而如果一段时间后上面的逻辑处理完毕，执行下面代码时发现 done 已经没有协程在等待
        // 将会永久阻塞，造成 goroutine 永远不会回收，也就是协程泄露
		done <- ans
	}()

    // select 等待子协程处理完毕或超时直接返回
	// 问：既然异步了，为什么这里还要阻塞等待执行完毕？不是还是又回到同步么？
	// 其实逻辑上 deal 函数还是同步的，开启一个协程只不过是为了能够实现超时功能
	select {
	case <-timeout:
		result = "超时"
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
	print("执行完毕，result 为：", result)
	for {}
}

func deal() string {
	timeout := time.After(time.Second * 2)
	done := make(chan string)

	result := ""

	go func() {
		time.Sleep(time.Second * 3)
		ans := "ayang"
		done <- ans
	}()

	select {
	case <-timeout:
		result = "超时"
	case result = <-done:
	}

    // 解决办法：关闭 chan。
    //（1）会唤醒所有阻塞等待在该 chan 的 done，并返回零值；
    //（2）如果有从已关闭的 chan 中取值，不会阻塞，也会返回零值
	// close(done)
    return result
}	
```

## End
