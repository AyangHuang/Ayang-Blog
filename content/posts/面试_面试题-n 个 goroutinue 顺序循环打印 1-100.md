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
title: "面试题-n 个 goroutine 循环打印 1-100"
date: 2023-07-16T01:12:45+08:00
lastmod: 2023-07-16T01:12:45+08:00
categories: ["面试"]
---

第一次面试遇到的问题。  

## 问题介绍

n 个 goroutine 顺序输出 1-100   

eg：   
2 个 goroutine：    
goroutine 0 ： 1  
goroutine 1 ： 2  
goroutine 0 ： 3  

## 我的不优雅做法

虽然结果是正确的，但实际不是面试官想要的方式。只要 sleep 稍微少一点，就会因为调度的问题导致结果不正确。

  ```go
package main

import (
	"fmt"
	"time"
)

func run(i int, c chan int) {
	for {
		n, is := <- c
		if is {
			fmt.Printf("goroutine %d : %d\n", i, n)
		}
	}
}

func main() {
	n := 0
	fmt.Scan(&n)
	s := make([]chan int, n)
	for i := 0; i < n; i++ {
		s[i] = make(chan int, 1)
		go run(i, s[i])
	}
	for i := 0; i < 100; i++ {
		s[i%n] <- i+1
		time.Sleep(time.Microsecond*800)
	}
}
  ```

## 面后反思：正确优雅做法

上面用 chan 来控制顺序，但可以直接用 chan 作为 n 个 goroutine 通信，完成了就通知下一个，即下一个 chan 是依赖上一个 chan 完成后来发送消息通知的。

```go
package main

import (
	"fmt"
	"time"
)

var n int

func run(i int, c []chan int) {
	for {
		var count int
		count, ok := <- c[i]
		if !ok {
			break
		}
		fmt.Printf("goroutine %d : %d\n", i, count)
		if count == 100 {
			for i := range c {
				close(c[i])
			}
			break
		}
		// 通知下一个 G
		c[(i+1)%n] <- count+1
	}
}

func main() {
	fmt.Scan(&n)
	s := make([]chan int, n)
	for i := 0; i < n; i++ {
		s[i] = make(chan int, 1)
		go run(i, s)
	}
	s[0] <- 1
	time.Sleep(time.Second)
}
```

## End
