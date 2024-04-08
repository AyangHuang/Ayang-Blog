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
title: "数据库加一个 field"
date: 2024-03-25T01:12:45+08:00
lastmod: 2024-03-25T15:58:26+08:00
categories: ["数据库"]
---

最近在实习遇到一个需要在数据库表 `property` 加一个字段存储的需求，mentor 告诉我不用直接在 `property` 表加字段，让我去看 `property_extra` 表，于是我恍然大悟，“加一个字段”还有这种骚操作 

## 表增加一个字段的三种方法

### ALTER TABLE

`ALTER TABLE table_name ADD COLUMN new_column INT;`

Alter Table 操作执行过程：

1. 用新的结构创建一张空表
2. 从旧表查出所有数据插入到新表中
3. 然后删除旧表

整个 DDL 语句执行过程中，会上 **MDL 写锁**，MDL 是避免 DML 和 DDL 并发执行的锁，所以在整个 DDL 语句执行过程中，无法进行执行 DML 增删改查。当然，MySQL 5.6 增加了 oline DDL 的机制，DDL 执行时，会降级成 MDL 读锁，DDL 语句也是上读锁，所以可以并发执行（当然对于 DDL 中加入新字段时数据 copy 是逃不掉的）

online DDL 具体可看：http://mysql.taobao.org/monthly/2021/03/06/

我理解在 DDL 时，DML 的 insert 和 update 的相关数据会写入一个 log 文件，只要 copy 后，对 log 文件在旧表跑一遍，就可以了

注意：MySQL 8.0 对于**新增列**支持 INSTANT DDL，不需要 copy 数据

### extra 字段

字节小说团队实习的时候，对于一张表加字段，是直接在加在该表的 `extra` 即可，即 `extra` 其实存的是 json 格式的 string 类型，只需要在业务层的 model 层将 extra 字段转化成 `map[string]interface{}`，增加枚举即可

优点：非常方便，几乎不需要改动   
缺点：（1）只查询该新增字段时，需要查出整个 `extra` 字段，性能较低  
     （2）业务增长，`extra` 会无限扩大，小说书籍元数据表 `extra` 里有两百多个字段

### table_extra 表

在原表的基础上增加 extra 表，那么就有两张表了

extra 表定义如下


二级索引为基础表的 id  
其实我理解如果查询比较多，写比较少的情况，可以直接把 extra 表主键设置为基础表的 id，这样可以杜绝回表

Add 和 Update 业务代码如下：



Select 业务代码如下：


优点：业务层改动较少  
缺点：查询需要多一个 extra 表的查询

## End
