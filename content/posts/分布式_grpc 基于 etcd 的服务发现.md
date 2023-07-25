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
linkToMarkdown: true
# 上面一般不用动
title: "grpc 基于 etcd 的服务发现"
date: 2023-07-25T01:12:45+08:00
lastmod: 2023-07-25T01:12:45+08:00
categories: ["分布式"]
tags: []
---

学习下 grpc 基于 etcd 的服务发现，做下学习笔记。

参考：

* [https://etcd.io/docs/v3.5/dev-guide/grpc_naming/](https://etcd.io/docs/v3.5/dev-guide/grpc_naming/)  

* [https://mp.weixin.qq.com/s/x-vC1gz7-x6ELjU-VYOTmA](https://mp.weixin.qq.com/s/x-vC1gz7-x6ELjU-VYOTmA)

* [https://mp.weixin.qq.com/s/iptZLaGFLd1rDedclQUMFg](https://mp.weixin.qq.com/s/iptZLaGFLd1rDedclQUMFg)

以下大部分理论内容摘抄自上面文章。

## 基本概念

{{< image src="/images/grpc  基于 etcd 服务发现/整体逻辑.png" width=100% height=100% caption="整体逻辑" >}}****

### 客户端的服务发现

通常情况下客户端需要知道服务端的 **IP + 端口号**才能建立连接，但服务端的IP和端口号并不是那么容易记忆。还有更重要的，在云部署的环境中，服务端的IP和端口可能随时会发生变化。

所以我们可以给某一个服务起一个名字，客户端通过名字创建与服务端的连接，客户端底层使用服务发现系统，解析这个名字来获取真正的IP和端口，并在服务端的IP和端口发生变化时，重新建立连接。这样的系统通常也会被叫做 **name-system（名字服务）**。（名字服务的概念很重要，下面会有用到）

**注册中心应用场景：**   
1. 其实注册中心相对于多加一层，解除 rpc 客户端对服务端的 ip:port 的依赖。如果服务端 ip:port 变动，只需要修改存储在注册中心的 key:value 的对应值；  
2. 一般注册中心都有类似于心跳的功能，能够检查 rpc 服务端是否正常运行；
3. 一个 rpc 服务有多个服务器实例，客户端可以通过相同的名字拉取到所有的地址；  


### 服务端的服务注册

如果 gRPC 服务端的地址是静态的，可以在客户端服务发现时直接解析为静态的地址

如果 gRPC 服务端的地址是动态的，可以有两种选择

* 自注册：当gRPC的服务启动后，向一个集中的**注册中心**进行注册（就是今天学习的 etcd）
* 平台的服务发现：使用k8s平台时，平台会感知gPRC实例的变化（貌似看过：本质还是依赖于 etcd 的注册中心？？？以后学 k8s 再学习把）

### 为什么是 etcd

以下来自 GPT：

Etcd 是一个开源的分布式键值存储系统，用于在分布式系统中存储配置数据、元数据和小规模的持久化数据。

主要特点包括：

1. 分布式：Etcd 被设计为在多个节点上运行，以保证高可用性和容错性。它使用 **Raft 一致性算法**来确保数据的**强一致性**，并在节点之间自动进行数据复制和同步。

2. 键值存储：Etcd 提供了一个简单的键值对存储模型，其中每个键和其对应的值都是字符串。这使得它非常适合用于存储配置数据和小规模的元数据。
   ```shell
    ayang@Ubuntu22:~$ etcdctl put ayang good # put key value
    OK
    ayang@Ubuntu22:~$ etcdctl get ayang # get key
    ayang
    good
    ayang@Ubuntu22:~$ etcdctl get --prefix a # 按前缀获取，挺重要的把我觉得
    ayang
    good
   ```

3. watch 机制：就是实时监听某一个 key，当 key 发生变化时，监听的服务可以快速发现。
   ```shell
    ayang@Ubuntu22:~$ etcdctl watch ayang
      # 会阻塞监听 ayang key 的变化，当变化时会即时接收到新的值**
   ```

4. 租约机制。注意：多个 key 可以共用同一个租约，可以给租约续期。
   ```shell
    ayang@Ubuntu22:~$ etcdctl lease grant 100 # 创建一个 100 秒过期的租约
    lease 694d8985ff795807 granted with TTL(100s) # 694... 为租约 id
    ayang@Ubuntu22:~$ etcdctl put ayang good --lease=694d8985ff795807 # 创建一个带租约的键值对
   ```

### RPC 的负载均衡

负载均衡的两种种方式：（文章里有第三种，我觉得实际用较少，就不写进来作为笔记记忆了）

1. 集中式负载均衡（Proxy Model）  
   例如 Nginx
    
    {{< image src="/images/grpc  基于 etcd 服务发现/集中式负载均衡.png" width=100% height=100% caption="集中式负载均衡" >}}  

	**缺点**：  
    * 请求转发的方式，转发耗时；   
    * 所有服务调用流量都经过 LB，LB 容易成为瓶颈；  
    * 存在单点故障，即 LB 宕机，影响大。  

2. 客户端负载均衡（Balancing-aware Client）

    {{< image src="/images/grpc  基于 etcd 服务发现/客户端负载均衡.png" width=100% height=100% caption="客户端负载均衡" >}}  

    **优点**：直接发送到服务端，不用经过 LB 转发，速度快；    
    **缺单**：如果有多个语言的客户端，每个客户端内部都要开发负载均衡的代码，代码量大。

**RPC 一般采用客户端内部的负载均衡的方式**。

注意：负载均衡发生的前提条件是有**多台完全相同的 RPC 服务端**。不然压根就不需要=负载均衡。

## 简单的整体逻辑

先以 grpc 来整理下整体逻辑。

1. 服务注册：  
   grpc 服务端以 key：value = **serverName：ip:port** 的键值对存入注册中心 etcd 中；

2. 服务发现：  
   grpc 客户端指定 serverName 为 **key**从注册中心拉去 value，即获得该 grpc 服务的地址；

## 测试

```shell
grpc-etcd
├── go.mod
├── grpc-client
│   └── main.go
├── grpc-server
│   ├── etcd.go
│   └── main.go
└── pb
    ├── hello_grpc.pb.go
    ├── hello.pb.go
    └── hello.proto
```

### proto 文件

```proto
// grpc-etcd/pb/hello/proto

syntax = "proto3"; // 版本声明，使用Protocol Buffers v3版本

option go_package = "grpc-etcd/pb";  // 指定生成的Go代码在你项目中的导入路径

package hello; // 包名，方便其他 proto 文件引入

// 定义服务，到时候是 Greeter.SayHello。注册到服务发现中心的是 ip:port
service Greeter1 {
    // SayHello 方法
    rpc SayHello1 (HelloRequest) returns (HelloResponse) {}
}


service Greeter2 {
    // SayHello 方法
    rpc SayHello2 (HelloRequest) returns (HelloResponse) {}
}

// 请求消息
message HelloRequest {
    string name = 1;
}

// 响应消息
message HelloResponse {
    string reply = 1;
}
```

### grpc 服务端

注意：两个不同的 grpc 服务端。

```go
// grpc-etcd/grpc-server/mian.go

package main

import (
	"context"
	"google.golang.org/grpc"
	"grpc-etcd/pb"
	"log"
	"net"
)

// 服务1
type server1 struct {
	pb.UnimplementedGreeter1Server
}

func (server1) SayHello1(context.Context, *pb.HelloRequest) (*pb.HelloResponse, error) {
	resp := new(pb.HelloResponse)
	resp.Reply = "server1:hello"
	return resp, nil
}

// 服务2
type server2 struct {
	pb.UnimplementedGreeter2Server
}

func (server2) SayHello2(context.Context, *pb.HelloRequest) (*pb.HelloResponse, error) {
	resp := new(pb.HelloResponse)
	resp.Reply = "server2:hello"
	return resp, nil
}

const (
	ServerAddr1 = "127.0.0.1:8080"
	ServerAddr2 = "127.0.0.1:8081"
	ServerName1 = "ayang/server1"
	ServerName2 = "ayang/server2"
)

func main() {
	var err error
	// 1. 创建两个 tcp 连接
	conn1, err := net.Listen("tcp", ServerAddr1)
	conn2, err := net.Listen("tcp", ServerAddr2)

	if err != nil {
		log.Fatal(err.Error())
		return
	}

	// 2. 创建两个 grpc 服务器
	s1 := grpc.NewServer()
	s2 := grpc.NewServer()

	// 3. 注册到 grpc 服务器中
	pb.RegisterGreeter1Server(s1, &server1{})
	pb.RegisterGreeter2Server(s2, &server2{})

	// 4. 注册到 etcd 中
	go registerEndPointToEtcd(context.TODO(), ServerAddr1, ServerName1)
	go registerEndPointToEtcd(context.TODO(), ServerAddr2, ServerName2)

	go func() {
		err = s1.Serve(conn1)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
	}()

	go func() {
		err = s2.Serve(conn2)
		if err != nil {
			log.Fatal(err.Error())
			return
		}
	}()

	<-make(chan struct{})

}
```

创建 etcd 客户端，向 etcd 中注册。（其实就是把 serverName:serverAddr 以 key:value 加入 etcd 中，并按时续期。）

```go
// grpc-etcd/grpc-server/etcd.go

package main

import (
	"context"
	"fmt"
	eclient "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/naming/endpoints"
	"time"
)

const (
    // etcd 服务器的地址
	EtcdAddr = "http://localhost:2379"
)

func registerEndPointToEtcd(ctx context.Context, serverAddr, serverName string) {
	// 创建 etcd 客户端
	etcdClient, _ := eclient.NewFromURL(EtcdAddr)
	etcdManager, _ := endpoints.NewManager(etcdClient, serverName) 

	// 创建一个租约，每隔 10s 需要向 etcd 汇报一次心跳，证明当前节点仍然存活
	var ttl int64 = 10
	lease, _ := etcdClient.Grant(ctx, ttl)

	// 添加注册节点到 etcd 中，并且携带上租约 id  
    // 以 serverName/serverAddr 为 key，serverAddr 为 value
    // serverName/serverAddr 中的 serverAddr 可以自定义，只要能够区分同一个 grpc 服务器功能的不同机器即可
	_ = etcdManager.AddEndpoint(ctx, fmt.Sprintf("%s/%s", serverName, serverAddr), endpoints.Endpoint{Addr: serverAddr}, eclient.WithLease(lease.ID))

	// 每隔 5 s进行一次延续租约的动作
	for {
		select {
		case <-time.After(5 * time.Second):
			// 续约操作
			resp, _ := etcdClient.KeepAliveOnce(ctx, lease.ID)
			fmt.Printf("keep alive resp: %+v\n", resp)
		case <-ctx.Done():
			return
		}
	}
}
```

etcd 存储的键值对如下：

```bash
ayang@Ubuntu22:~$ etcdctl get --prefix ""
ayang/server1/127.0.0.1:8080
{"Op":0,"Addr":"127.0.0.1:8080","Metadata":null}
ayang/server2/127.0.0.1:8081
{"Op":0,"Addr":"127.0.0.1:8081","Metadata":null}
```

### grpc 客户端测试

```go
// grpc-etcd/grpc-client/etcd.go

package main

import (
	"context"
	"fmt"
	eresolver "go.etcd.io/etcd/client/v3/naming/resolver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-etcd/pb"
	"log"
	// etcd
	eclient "go.etcd.io/etcd/client/v3"
)

const (
	serverNamePreResolve = "etcd:///ayang/server1"
	EtcdAddr             = "http://localhost:2379"
)

func main() {
	var err error
	// 创建 etcd 客户端
	etcdClient, err := eclient.NewFromURL(EtcdAddr)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	// 创建 etcd 实现的 grpc 服务注册发现模块 resolver
	etcdResolverBuilder, err := eresolver.NewBuilder(etcdClient)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	// 创建 grpc 连接代理
	conn, err := grpc.Dial(
		// 服务名称
		serverNamePreResolve,
		// 注入 etcd resolver
		grpc.WithResolvers(etcdResolverBuilder),
		// 声明使用的负载均衡策略为 roundrobin，轮询。（测试 target 时去除该注释）
		// grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln(err.Error())
		return
	}

	for i := 0; i < 4; i++ {
		greeter1 := pb.NewGreeter1Client(conn)
		resp, err := greeter1.SayHello1(context.Background(), &pb.HelloRequest{
			Name: "ayang",
		})

		if err != nil {
			log.Fatalln(err.Error())
			return
		}

		fmt.Printf("%d  %s\n", i, resp.Reply)
	}

	defer conn.Close()
}
```

客户端打印结果如下：

```go
0  server1:hello
1  server1:hello
2  server1:hello
3  server1:hello
```

### 测试：grpc.Dail() 中的 target 参数

如果没有加 etcd 注册中心，采用直连的方式，就是传入 grpc 服务端的 ip:port。**直连的方式不需要注入名字服务解析器**。

如果注入了 etcd 提供的名字服务解析器，经过测试，我认为 target 参数是**服务器名字前缀**。（注意：这里只是我根据测试结果得出的结论）

测试如下：

**测试一：**

修改客户端代码       
1serverNamePreResolve = "etcd:///ayang/server"  

结果：多运行几次，有 1/2 的概率会发生错误 `2023/07/25 23:40:45 rpc error: code = Unimplemented desc = unknown service hello.Greeter1`  

**测试二：** 

修改客户端代码     
（1）serverNamePreResolve = "etcd:///ayang/server"    
（2）同时注释掉负载均衡的注释，即开启轮询的负载均衡  
   
结果：在第一次调用或第二次调用会发生错误

```bash
0  server1:hello
2023/07/25 23:43:28 rpc error: code = Unimplemented desc = unknown service hello.Greeter1
```

**测试结果分析**

以 ayang/server 为前缀拉去到 ayang/server1 和 ayang/server2 的 ip:prot，并认为是同一个 grpc 服务的不同服务器实例。实际是完全不同的服务器，所以客户端在调用发送到 grpc 服务端时，服务端发现没有该服务方法，返回 error。

也解决了我的一个困惑，**grpc.Dail 只对应一个 grpc 服务**（当然一个 grpc 服务可以有多个服务器实例），内部有连接池，封装了 **rpc 通信、多服务器负载均衡**的过程。如果需要调用多个 grpc 服务，即需要需要多次调用 grpc.Dail。

注意：

## 具体的整体逻辑

最后整理下整体的逻辑。（注意跟上面的整体逻辑进行对比，增加了一些细节）

1. 服务注册：  
   grpc 服务端以 key：value = **serverName/serverN：ip:port** 的键值对存入注册中心 etcd 中。serverName/serverN 中的 serverN 可以是任意的，只要识别出不同的 grpc 服务实例即可。

2. 服务发现：  
   grpc 客户端指定 serverName 为 **key 前缀**从注册中心拉去 value，即获得该 grpc 服务的多个实例地址（如有的话）；

{{< image src="/images/grpc  基于 etcd 服务发现/完整整体逻辑.png" width=100% height=100% caption="完整整体逻辑" >}}