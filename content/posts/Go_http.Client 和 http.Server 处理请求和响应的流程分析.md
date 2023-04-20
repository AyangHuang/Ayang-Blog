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
title: "http.Client 和 http.Server 处理请求和响应的流程分析"
date: 2023-04-11T00:11:45+08:00
lastmod: 2023-04-11T00:11:45+08:00
categories: ["Go"]
tags: []
---

参考： https://mp.weixin.qq.com/s/zFG6_o0IKjXh4RxKmPTt4g

这篇文章写得非常棒，看完再去看源码就轻松很多了。下面是初略看了源码后的**流程总结**。

## 服务端

### 处理请求流程

{{< image src="/images/http.Client 和 http.Server 处理请求和响应的流程分析/服务端.jpg" width=100% height=100% caption="客户端" >}}

* **一条 TCP 连接一个协程处理**，而且是串行化，一个 HTTP 请求就是一个“事务”，必须读取 Request，发送 Response 后才能处理下一个 Request 请求报文。    
* 这样仅仅是用到了 TCP 的**半双工**，即虽然是可读可写，但是读时不写，写时不读。

其实 HTTP 是有 **pipeline 机制**的，即客户端可以连续发送多个 Request 而不需要等待 Response，但一般不会这么实现，Go 的 HTTP 也没有这么实现。下面是源码的注释：

```go
// HTTP cannot have multiple simultaneous active requests.[*]
// Until the server replies to this request, it can't read another,
// so we might as well run the handler in this goroutine.
// [*] Not strictly true: HTTP pipelining. We could let them all process
// in parallel even if their responses need to be serialized.
// But we're not going to implement HTTP pipelining because it
// was never deployed in the wild and the answer is HTTP/2.

// HTTP 不能有多个同时活动的请求。[*]
// 在服务器回复这个请求之前，它不能读取另一个，
// 所以我们不妨在这个 goroutine 中运行处理程序。
// [*] 不完全正确：HTTP 流水线。 我们可以让他们全部处理
// 并行，即使它们的响应需要序列化。
// 但是我们不打算实现 HTTP 流水线，因为它
// 从未在野外部署过，答案是 HTTP/2。
```

## 客户端

### 发送请求报文流程

客户端比较复杂一点。

* 客户端使用了**连接池**，复用 TCP 连接，避免重复建立连接。  
* 充分利用 goroutine 机制，使用了 chan 进行**协程间通信**和**超时控制**等；

{{< image src="/images/http.Client 和 http.Server 处理请求和响应的流程分析/客户端.jpg" width=100% height=100% caption="客户端" >}}

**为什么服务端的 read 和 write 不是同时进行的，而客户端可以呢？**

```go
// Write the request concurrently with waiting for a response,
// in case the server decides to reply before reading our full
// request body.

// 在等待响应的同时写入请求，
// 如果服务器决定在阅读我们的完整内容之前回复
// 请求正文。
```

## 实战：利用连接池配置实现长连接

* 客户端和服务端都配置了永不关闭 TCP 连接  

* 客户端连接池配置：    
  * 同一个服务端连接池最多保留 2 个**空闲连接**   
  * 同一时间，只允许和一个服务端保持 3 个连接  

```go
func ServeHTTP(w http.ResponseWriter, req *http.Request) {
	fmt.Println("服务端收到 request")
	time.Sleep(time.Second * 10)
	w.WriteHeader(200)
	w.Write([]byte("客户端收到 response"))
}

func Client() {
	// 客户端
	client := &http.Client{
		Transport: &http.Transport{
			IdleConnTimeout:     0, // 连接永不超时
			MaxIdleConnsPerHost: 2, // 控制每个服务端保持的最大空闲（保持活动）连接数
			MaxConnsPerHost:3,// 和同一个服务端同时连接的最大数目
		},
	}

	// 发送三个请求
	for i := 0; i < 4; i++ {
		go func() {
			resq, err := client.Get("http://localhost:9999/")

			if err == nil {
				var s []byte = make([]byte, 30)
				resq.Body.Read(s)
				fmt.Println()
				fmt.Println(string(s))
			}
			resq.Body.Close()
		}()
	}
}

func Server() {
	// 服务端
	server := &http.Server{
		Addr:        "localhost:9999",
		IdleTimeout: 0, // 连接永不超时，如果服务端设置超时时间 30 分钟，30 分钟内没有新的报文，服务端会主动关闭连接。
		Handler:     http.DefaultServeMux,
	}
	http.DefaultServeMux.HandleFunc("/", ServeHTTP)
	server.ListenAndServe()
}

func main() {
	go Server()
	time.Sleep(time.Second)
	go Client()
	select {}
}
```

**实验结果：**

{{< image src="/images/http.Client 和 http.Server 处理请求和响应的流程分析/实验结果.png" width=100% height=100% caption="实验结果" >}}

## End
