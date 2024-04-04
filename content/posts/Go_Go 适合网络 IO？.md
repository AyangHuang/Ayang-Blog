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
title: "Go 适合网络 IO"
date: 2024-02-08T01:12:45+08:00
lastmod: 2024-02-08T01:12:45+08:00
categories: ["Go"]
tags: []
---

最近在看一些 Java 的内容，途中会对比着 Go 学一下。发现 Go 的协程机制真的很牛逼，协程的机制搭配 epoll 使得 Go 非常适合网络 IO。本篇文章会讲下八股文，然后对比测试下网络 IO 和磁盘 IO

PS：需要对协程和 GMP 的基础知识有深入理解，对 epoll 有了解

## Go 适合网络 IO

计算机程序场景分为：**CPU 密集型和 IO 密集型**。IO 又分为网络 IO 和磁盘 IO

**为什么 Go 适合网络 IO 呢？**，这得先讲讲网络 socket 句柄和文件句柄的不同：

* **socket 句柄**
  * 实现了 `.poll`，可以用 epoll 池来管理；
  * socket 句柄可读（socket buffer 对端网络发数据过来）可写（socket buffer 有空间可以写入数据）事件有意义
  * 并且 go 中把 epoll 的 socket 句柄设置为 **noblocking**。这样当 epoll 中没有 socket fd 可读可写时，也会直接返回，然后调度其他协程，这样就实现了网络 IO 的并发 

* **文件句柄**
  * 文件句柄一般没有实现 `.poll`
  * 文件 IO 的 read/write 都是阻塞的
  * 文件句柄可读可写事件则没有意义，因为文件句柄理论上是永远都是可读可写的，不会阻塞等待

再讲**磁盘 IO 导致线程阻塞**：

磁盘 IO 是系统调用，且读写是阻塞的，只有完成整个系统调用才会返回，所以会导致**卡线程**。Go 能做的就是在**执行阻塞系统调用前，解除当前执行线程 M 和当前 P 的绑定，然后执行当前 G 的阻塞系统调用**。并且调度线程可能会**创建新线程**（不一定会创建）来绑定这个 P，然后执行 runq 队列中的 G。当进行很多磁盘 IO 且这些 IO 都很慢时，会**导致创建出很多线程并处于阻塞状态**。所以，对于磁盘 IO，Go 协程的机制和线程机制没有太大的区别，都是耗费一个线程阻塞等待

**那这种协程的机制对比线程的机制会有什么不同？**

看一个场景：Web 服务器，当一个 HTTP 请求过来，线程模型一个线程处理一个请求，协程模型一个协程处理一个请求，当该请求依赖于另一个服务时，需要通过网络 IO 再请求另一个服务器。

线程模型：  
此时，线程模型的线程只能阻塞等待，即使线程模型内部封装了 epoll，其实也是只有那个 epoll 线程阻塞了，但是处理请求的线程依赖于该请求的结果，所以要么阻塞同步等待 epoll 线程唤醒该线程，要么注册一个回调函数，该回调函数的函数体是处理该请求的剩余代码，然后该线程继续处理下一个请求。待 epoll 收到请求响应后执行该回调函数

协程模型：  
协程模型其实就是由 Go 的语言内部直接封装了 epoll 和“回调函数”，只是底层的协程机制让我们在 Go 业务代码书写层面使用同步的方法达到异步回调的效果，因为 Go 协程封装了帮我们“回调”的逻辑，我们可以理解 Go 在应用层把一个线程（其实就是指令流），切分成一串一串指令流，并发交替执行

## 实验网络 IO 和磁盘 IO

实验环境：Ubuntu22，CPU 32 核，内存 32 G

磁盘 IO：

```go

func main() {
    // 奇怪，设置 P 为 1 时，阻塞线程会变得比 32 少，可能有限制阻塞线程和 P 的最大比例关系？？？
	runtime.GOMAXPROCS(4)
	printThreadInfo("初始化 main 时：", "ps", "-T", "-p", strconv.Itoa(os.Getpid()))
	wg := sync.WaitGroup{}
	wg.Add(32)
	for i := 1; i <= 32; i++ {
		i := i
		go func() {
            // 32 个文件，每个文件名字 1 2 3，依次类推，每个文件 300 多 MB
			file, err := os.Open(fmt.Sprintf("/home/ayang/Downloads/%d", i))
			if err != nil {
				panic(err.Error())
			}
			_, err = io.ReadAll(file)
			if err != nil {
				panic(err.Error())
			}
			wg.Done()
		}()
	}

	time.Sleep(time.Second * 1)
	printThreadInfo("文件读取中：", "ps", "-T", "-p", strconv.Itoa(os.Getpid()))
	wg.Wait()

	printThreadInfo("文件读取后：", "ps", "-T", "-p", strconv.Itoa(os.Getpid()))
}

func printThreadInfo(preStr string, name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	output, err := cmd.Output()
	if err != nil {
		panic(err.Error())
		return
	}
	fmt.Println(preStr)
	fmt.Println(string(output))
}
```

实验结果：

```txt
初始化 main 时：
    PID    SPID TTY          TIME CMD
 194880  194880 pts/1    00:00:00 main
 194880  194881 pts/1    00:00:00 main
 194880  194882 pts/1    00:00:00 main
 194880  194883 pts/1    00:00:00 main
 194880  194884 pts/1    00:00:00 main

文件读取中：
    PID    SPID TTY          TIME CMD
 194880  194880 pts/1    00:00:00 main
 194880  194881 pts/1    00:00:00 main
 194880  194882 pts/1    00:00:00 main
 194880  194883 pts/1    00:00:00 main
 194880  194884 pts/1    00:00:00 main
 194880  194889 pts/1    00:00:00 main
 194880  194890 pts/1    00:00:00 main
 194880  194891 pts/1    00:00:00 main
 194880  194892 pts/1    00:00:00 main
 194880  194893 pts/1    00:00:00 main
 194880  194894 pts/1    00:00:00 main
 194880  194895 pts/1    00:00:00 main
 194880  194896 pts/1    00:00:00 main
 194880  194897 pts/1    00:00:00 main
 194880  194898 pts/1    00:00:00 main
 194880  194899 pts/1    00:00:00 main
 194880  194900 pts/1    00:00:00 main
 194880  194901 pts/1    00:00:00 main
 194880  194902 pts/1    00:00:00 main
 194880  194903 pts/1    00:00:00 main
 194880  194904 pts/1    00:00:00 main
 194880  194905 pts/1    00:00:00 main
 194880  194906 pts/1    00:00:00 main
 194880  194907 pts/1    00:00:00 main

文件读取后：
    PID    SPID TTY          TIME CMD
 194880  194880 pts/1    00:00:00 main
 194880  194881 pts/1    00:00:00 main
 194880  194882 pts/1    00:00:01 main
 194880  194883 pts/1    00:00:04 main
 194880  194884 pts/1    00:00:03 main
 194880  194889 pts/1    00:00:06 main
 194880  194890 pts/1    00:00:00 main
 194880  194891 pts/1    00:00:00 main
 194880  194892 pts/1    00:00:01 main
 194880  194893 pts/1    00:00:00 main
 194880  194894 pts/1    00:00:01 main
 194880  194895 pts/1    00:00:04 main
 194880  194896 pts/1    00:00:01 main
 194880  194897 pts/1    00:00:00 main
 194880  194898 pts/1    00:00:02 main
 194880  194899 pts/1    00:00:00 main
 194880  194900 pts/1    00:00:00 main
 194880  194901 pts/1    00:00:01 main
 194880  194902 pts/1    00:00:00 main
 194880  194903 pts/1    00:00:01 main
 194880  194904 pts/1    00:00:01 main
 194880  194905 pts/1    00:00:01 main
 194880  194906 pts/1    00:00:02 main
 194880  194907 pts/1    00:00:05 main
 194880  194923 pts/1    00:00:00 main
 194880  194924 pts/1    00:00:02 main
 194880  194925 pts/1    00:00:00 main
 194880  194926 pts/1    00:00:00 main
 194880  194927 pts/1    00:00:02 main
 194880  194928 pts/1    00:00:00 main
 194880  194929 pts/1    00:00:02 main
 194880  194930 pts/1    00:00:01 main
 194880  194931 pts/1    00:00:02 main
 194880  194932 pts/1    00:00:00 main
 194880  194933 pts/1    00:00:00 main
```

网络 IO：

```go
func main() {
	runtime.GOMAXPROCS(4)
	printThreadInfo("初始化 main 时：", "ps", "-T", "-p", strconv.Itoa(os.Getpid()))
	for i := 1; i <= 32; i++ {
		i := i
		go func() {
			// 开启 32 个服务器
			err := http.ListenAndServe(fmt.Sprintf(":%d", i+65500), nil)
			if err != nil {
				panic(err.Error())
			}
		}()
	}
	time.Sleep(time.Second * 3)
	printThreadInfo("开启 32 个服务器后：", "ps", "-T", "-p", strconv.Itoa(os.Getpid()))
	var c chan struct{}
	<-c
}

func printThreadInfo(preStr string, name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	output, err := cmd.Output()
	if err != nil {
		panic(err.Error())
		return
	}
	fmt.Println(preStr)
	fmt.Println(string(output))
}
```

实验结果：

```text
初始化 main 时：
    PID    SPID TTY          TIME CMD
 223184  223184 pts/1    00:00:00 main
 223184  223185 pts/1    00:00:00 main
 223184  223186 pts/1    00:00:00 main
 223184  223187 pts/1    00:00:00 main
 223184  223188 pts/1    00:00:00 main
 223184  223189 pts/1    00:00:00 main

开启 32 个服务器后：
    PID    SPID TTY          TIME CMD
 223184  223184 pts/1    00:00:00 main
 223184  223185 pts/1    00:00:00 main
 223184  223186 pts/1    00:00:00 main
 223184  223187 pts/1    00:00:00 main
 223184  223188 pts/1    00:00:00 main
 223184  223189 pts/1    00:00:00 main
 223184  223191 pts/1    00:00:00 main
```

结论：磁盘 IO 导致线程阻塞，网络 IO 可以切换协程执行，体现为多网络 IO（请求）并发处理

## End
