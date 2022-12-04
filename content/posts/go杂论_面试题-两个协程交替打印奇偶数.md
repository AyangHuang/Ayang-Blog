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
title: "面试题-两个协程交替打印奇偶数（内含三种方法）"
date: [":git", ":fileModTime"]
date: 2022-12-05T19:52:23+08:00
lastmod: 2022-12-05T19:52:23+08:00
tags: ["Go"]
---

最近翻了下牛客网，发现都在喊秋招缩招，一片红海之类的。go 岗位数量更是少。难呀！不管如何，还是得迎难而上！！！  

随手翻到一道面试题：两个协程交替打印奇偶数，从 1 到 100。这我不得试一试？

## 方法一

先上代码。

```go
// 包变量
var count int32 = 1

// 打印奇数
func printOne() {
	for {
        // 死循环无限取
		num := atomic.LoadInt32(&count)
		if num > 100 {
			return
		}
        // 取的不是想要就不理，继续无限取
		if num%2 == 1 {
			println("奇数：", num)
			atomic.AddInt32(&count, 1)
		}
	}
}

// 打印偶数
func printTwo() {
	for {
        // 死循环无限取
		num := atomic.LoadInt32(&count)
		if num > 100 {
			return
		}
        // 取的不是想要就不理，继续无限取
		if num%2 == 0 {
			println("偶数：", num)
			atomic.AddInt32(&count, 1)
		}
	}
}

func main() {
	go printOne()
	go printTwo()
	var c chan struct{}
    // 阻塞
	<-c
}
```

第一个方法利用**原子操作保证并发安全**。两个协程 while 无限取，判断是否取到想要的数，如果是想要的，则 ++。  
eg：奇数协程取到 1，则打印并 ++ 变成 2，继续运行，发现还是 2，continue....直到偶数协程取到 2，然后 ++ 变成 3，偶数线程陷入无效自旋，奇数线程运行...实现两个协程**交替运行**的效果。

## 方法二

还是先上代码。

```go
func print(c1, c2 chan int, str string) {
	for {
		num, closed := <-c1
        // 判断 chanel 是否已被另一个协程关闭，如关闭则说明到 100 了，直接返回
		if closed {
			if num <= 100 {
				println(str, num)
                // 让另一个线程运行
				c2 <- num + 1
			} else {
				close(c1)
				close(c2)
				return
			}
		} else {
			return
		}
	}

}

func main() {
	c1 := make(chan int, 1)
	c2 := make(chan int, 1)
    // 注意：上面的 c1 必须是有缓冲的，不然会阻塞在这里
	c1 <- 1
	go print(c1, c2, "奇数：")
	go print(c2, c1, "偶数：")
	time.Sleep(time.Second)
}
```

第二个方法呢是利用两个协程各自绑定一个 chanel，两个 chanel 相互依赖，利用 **chanel 的收发阻塞来阻塞和运行协程**。  
写文章写到这里，发现其本质（上一行的加粗字体），那么可以更改更简单的代码。

```go
// 包变量
var count int32 = 1

func print(c1, c2 chan struct{}, str string) {
	for {
		if _, closed := <-c1; closed {
			if count <= 100 {
				println(str, count)
				count++
				c2 <- struct{}{}
			} else {
				close(c1)
				close(c2)
				return
			}
		} else {
			return
		}
	}

}
func main() {
	c1 := make(chan struct{})
	c2 := make(chan struct{})
	go print(c1, c2, "奇数：")
	go print(c2, c1, "偶数：")
	go func() {
		c1 <- struct{}{}
	}()
	time.Sleep(time.Second)
}
```

额，好像没区别哈哈。只是改成了无缓冲 chanel，更纯粹一点就是利用了 chanel 实现两个协程交互运行，count 直接用包变量，跟 chanel 解耦了。好吧，还是没区别哈哈，强行解释。既然写都写了，就不删了，凑下字数。

## 方法三

法三用了条件变量。

```go
var condition int = 1 // 1 是奇数，2 是偶数
var count int = 1

func printOne(c1, c2 *sync.Cond) {
	for {
		c1.L.Lock()
		// 这里可以换成 if，因为就俩协程而已，多协程就必须用 while
		for condition != 1 {
			c1.Wait()
		}
		c1.L.Unlock()

		if count <= 100 {
			println("奇数：", count)
			count++
		}

		c2.L.Lock()
		condition = 2
		c2.Signal()
		c2.L.Unlock()

		if count > 100 {
			println("我结束了")
			return
		}
	}
}

func printTwo(c1, c2 *sync.Cond) {
	for {
		c1.L.Lock()
		// 这里可以换成 if，因为就俩协程而已，多协程就不行
		for condition != 2 {
			c1.Wait()
		}
		c1.L.Unlock()

		if count <= 100 {
			println("偶数数：", count)
			count++
		}

		c2.L.Lock()
		condition = 1
		c2.Signal()
		c2.L.Unlock()

		if count > 100 {
			println("我结束了")
			return
		}
	}
}

func main() {
	m := new(sync.Mutex)
	c1, c2 := sync.NewCond(m), sync.NewCond(m)
	// 这里每一个条件变量配合一个锁也可以
	// m2 := new(sync.Mutex)
	// c1, c2 := sync.NewCond(m), sync.NewCond(m2)
	go printOne(c1, c2)
	go printTwo(c2, c1)
	time.Sleep(time.Second)
}
```

## 总结

其实要实现两个协程交替运行，本质是要实现线程的相互通信。

1. 法一直接用 count 作为通信的载体，通知取消自旋
2. 法二利用 go 的特性，用 chanel 作为通信的载体（chanel 内部其实也是用自旋实现阻塞的效果）
3. 法三运用条件变量，用 condition 来通信实现真正的线程休眠唤醒。

如果你有其他方法欢迎交流哈！

## End
