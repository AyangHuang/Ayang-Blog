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
title: "database.sql 库-策略模式"
date: 2022-12-17T00:17:04+08:00
lastmod: 2022-12-17T00:17:04+08:00
categories: ["Go"]
tags: ["设计模式"]
---

## 策略模式简介

> **定义**：它定义了算法家族，分别封装起来，让他们之间可以相互替换，此模式让算法的改变，不会影响到使用算法的用户。  

{{< image src="/images/database.sql 库-策略模式/策略模式.png" width=100% height=100% caption="go函数调用栈帧" >}}

```java
// 策略接口
public interface Strategy {
    public void Algorithm();
}

// 具体算法1
class Algorithm1 implements Strategy {
    @Override
    public void Algorithm() {
        System.out.println("使用策略1");
    }
}

// 具体算法2
class Algorithm2 implements Strategy {
    @Override
    public void Algorithm() {
        System.out.println("使用策略2");
    }
}

// 提供算法服务
class Serve {
    // 具体策略
    private Strategy strategy;
    // 设置具体算法策略，可用简单工厂模式
    public void setStrategy(int algo) {
        switch (algo) {
            case 1:
                this.strategy = new Algorithm1();
                break;
            case 2:
                this.strategy = new Algorithm2();
                break;
            default:
                this.strategy = new Algorithm1();
        }
    }

    public void useAlgo() {
        strategy.Algorithm();
    }

    public Serve(int algo) {
        setStrategy(algo);
    }
}
```

## database.sql 

> 这里是 1.19.4

```go
package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql" // 匿名映入
)

var db *sql.DB

func main() {
	dsn := "root:123456@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=True"
	db, _ = sql.Open("mysql", dsn)
	_ = db.Ping()
}
```

* 策略接口   
  
  首先可以明确 database/sql 是类 SQL 数据库都可以调用的，其中 driver 文件下的 driver.go 定义了很多接口，相当于策略接口。  

* Serve 类  
  
  database/sql/sql.go 里面的 DB 结构体，相当于上面的 Serve 类。

* Strategy 类

  database/sql/sql.go 里面有 `map[string]driver.Driver`，相当与 Strategy 策略类。

  ```go
  // sql 包变量
  var (
	driversMu sync.RWMutex
	drivers   = make(map[string]driver.Driver)
  )
  ```

* Algorithm 类

  github.com/go-sql-driver/mysql 相当于具体的算法类。里面实现了 database/sql driver.go 的所有接口。


**接下来我们追下上面代码的流程：**  

1. `_ "github.com/go-sql-driver/mysql"`   
   
   mysql 通过匿名引入的方法，引入时会调用包里的  `init()` 函数如下：
   
   ```go
    // github.com/go-sql-driver/mysql driver.go
    func init() {
        // 注意：调用的是 sql 包
	    sql.Register("mysql", &MySQLDriver{})
    }
    
    // database/sql sql.go
    func Register(name string, driver driver.Driver) {
	    driversMu.Lock()
	    defer driversMu.Unlock()
	    drivers[name] = driver
    }
   ```
   
   实际这里就是所谓的注册 MySQL 的驱动，就是策略模式里面的**设置具体算法**的步骤。

2. `db, _ = sql.Open("mysql", dsn)`
   
   这里 open 函数并不会真正地与数据库建立连接，而是检查 database/sql 里面的包变量 drivers 有没有 mysql 这个驱动，如果我们上面有通过匿名引入的方式注册 mysql 驱动，那么就不会报错。然后初始化 DB 对象并返回。
   
   ```go
   // database/sql sql.go
    func Open(driverName, dataSourceName string) (*DB, error) {
	    driveri, ok := drivers[driverName]
    	if !ok {
    		return nil, "err"
    	}
    	return OpenDB(dsnConnector{dsn: dataSourceName, driver: driveri}), nil
    }
    
    func OpenDB(c driver.Connector) *DB {
    	db := &DB{
    		connector:    c,
    		openerCh:     make(chan struct{}, connectionRequestQueueSize),
    		lastPut:      make(map[*driverConn]string),
    		connRequests: make(map[uint64]chan connRequest),
    		stop:         cancel,
    	}
    	return db
    }
   ```
3. `db.Ping()`  
   
   Ping 验证与数据库的连接是否仍处于活动状态，如有必要，建立连接。
   `db.Ping()` 以及如果调用 `db.Query("select * from demo")` 其实都是内部真正执行的是 `db.driver.ping()` 或者 `db.driver.query()`。如果所有代码都写好了，想更改数据库，只需要更改引入的数据库类型，也就是更改注册的数据库驱动的类型即可。实现了**解耦**。
   
## End