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
title: "面试题-用 Go 实现非阻塞互斥锁"
date: 2023-07-16T00:12:45+08:00
lastmod: 2023-07-16T00:12:45+08:00
categories: ["面试"]
---

刷牛客看到的问题，尝试解决下并记录下来。

## 利用 chan + select 实现

```go
type Mutex struct {
	c chan struct{} 
}

func NewMutex() *Mutex {
	m := &Mutex{
		c: make(chan struct{}, 1),
	}

	m.c <- struct{}{}
	return m
}

func (m *Mutex) Lock() bool {
	select {
	case <- m.c:
		return true
	default:
		return false
	}
}

func (m *Mutex) UnLock() {
	select {
	case <- m.c:
		panic("没有上锁，却调用解锁")
	default:
	}
	m.c <- struct{}{}
}


var count = 0

// Test
func main() {
	mutex := NewMutex()
	m := make(map[int]int)
	
	for i := 0; i < 100; i++ {
		go func() {
			for {
				if mutex.Lock() {
					m[count] = count
					print(count)
					count++
					time.Sleep(time.Second)
					mutex.UnLock()
				} 
			}
			
		}()
	}

	var c  chan struct{}
	<- c
}
```

这里有两种校验锁是否满足并发安全的方法：

1. 使用 map 作为校验并发安全性，因为 map 内部有并发安全的校验机制，如果发生并发读写，会 panic  
2. 为了观看方便，print 打印出来

## End
