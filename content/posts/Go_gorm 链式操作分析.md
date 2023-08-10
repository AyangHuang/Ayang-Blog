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
title: "gorm 链式操作源码分析"
date: 2023-08-10T01:12:45+08:00
lastmod: 2023-08-10T01:12:45+08:00
categories: ["Go"]
tags: []
---

对 gorm 的链式操作一直有疑惑，最近做项目，翻了下源码看了具体的实现，做下记录。

当然推荐先看官方文档了解下链式操作罗：  
[https://gorm.io/zh_CN/docs/method_chaining.html](https://gorm.io/zh_CN/docs/method_chaining.html)

看完是不是有很多疑问？例如：   

{{< image src="/images/gorm 链式操作分析/why.png" width=100% height=100% caption="why" >}}  

## getInstance() 源码分析

链式方法在每次调用时，内部会先调用 `getInstance()` 函数，该函数的实现如下：

```go
func (db *DB) getInstance() *DB {
	if db.clone > 0 {
        // 创建新的 DB 实例，由于 Go 变量零值的原因，新的 DB.clone = 0 
		tx := &DB{Config: db.Config, Error: db.Error}

		if db.clone == 1 {
			// 创建一个全新的 statement
			tx.Statement = &Statement{
				DB:       tx,
				ConnPool: db.Statement.ConnPool,
				Context:  db.Statement.Context,
				Clauses:  map[string]clause.Clause{},
				Vars:     make([]interface{}, 0, 8),
			}
		} else {
            // 复用原来的的 statement，即 SQL 语句会复用
			tx.Statement = db.Statement.clone()
			tx.Statement.DB = tx
		}

		return tx
	}

    // clone = 0，直接返回原 db
	return db
}
```

再查看 clone 值的使用：

{{< image src="/images/gorm 链式操作分析/clone.png" width=100% height=100% caption="clone" >}}  

追进去，发现只有两个函数三处地方对 clone 修改：

```go
func Open(dialector Dialector, opts ...Option) (db *DB, err error) {
    ...
    // 第一个函数，第一处
	db = &DB{Config: config, clone: 1}
    ...
}

func (db *DB) Session(config *Session) *DB {
	var (
		txConfig = *db.Config
		tx       = &DB{
			Config:    &txConfig,
			Statement: db.Statement,
			Error:     db.Error,
            // 第二个函数，第一处
			clone:     1,
		}
	)

    ...
    // 第二个函数，第二处
    // 只有 db.Session(&gorm.config{NewDB:true}) 才不会执行这里的逻辑
    //（因为 NewDB 的零值为 false），调用即 Session 后的新 DB 的 clone 值默认设置为 2
	if !config.NewDB {
		tx.clone = 2
	}
    ...
	return tx
}
```

其实还有第三个函数会对 clone 值做出修改，即 `getInstance()` 函数，因为创建一个新的 DB 后的 clone 为零值，默认为 0 。

根据以上源码分析，可以得到下图：  

1. 箭头表示调用 `getInstance` 的 `clone` 的变化；
  
2. 椭圆旁边的文字表示 clone 强制改变的方式。

{{< image src="/images/gorm 链式操作分析/clone 状态转变.png" width=100% height=100% caption="clone 状态转变" >}} 

注意：  

1. 一旦 `clone` 为 1 或 2，在首次调用 `getInstance()` 后，`clone` 都为变成 0；而 clone 为 0，调用 `getInstance()` 时会陷入死循环，即无限为 0。  

2. 只有两个函数能修改，将 `clone` 从 0 变 1 或 2：`WithContext()` 和 `Session()`（`WithContext()` 内部调用 `Session()`）。`Session()` 可以根据传入的 `config.NewDB` 的值来决定把 `clone` 变成 1 还是 2。

## End
