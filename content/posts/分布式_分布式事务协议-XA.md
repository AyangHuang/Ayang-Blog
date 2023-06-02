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
title: "分布式事务协议-XA"
date: 2023-06-02T01:12:45+08:00
lastmod: 2023-06-02T15:58:26+08:00
categories: ["分布式"]
---


## XA 协议引出的问题

如果不懂 XA 协议，直接看下面的文章把：

[https://dtm.pub/practice/xa.html](https://dtm.pub/practice/xa.html)

下面是我这篇文章主要想说的：（以下是我测试得出的，可能有误）

问题：**如果 commit 阶段前一个分事务 commit 成功（无法回滚），而 另一个分事务在要 commit 时宕机了，此时该怎么办？**

先说答案：只要是 **prepare 成功，即表示该事务已经执行且 redo 日志已经刷新到磁盘**。所以我们只需要宕机重启，然后客户端不断进行 commit 尝试，该事务还是能够正常 commit。这就**与普通的事务有区别，普通事务在宕机后，会直接回滚**。而 prepare 后的 XA 事务则不会。

注意：未 prepare 的 xa 事务，宕机重启会自动回滚，无法继续 prepare 和 commit。

## 测试

* 初始状态
  
```bash
mysql> select * from grades;
+------------+-----------+-------+
| id         | course    | grade |
+------------+-----------+-------+
| 2125121024 | 数据库    |    99 |
+------------+-----------+-------+
1 row in set (0.00 sec)
```

* 执行 xa 事务到 prepare 阶段，待 commit

```bash
mysql> xa start 'c';
Query OK, 0 rows affected (0.00 sec)

mysql> delete from grades where grade=99;
Query OK, 1 row affected (0.00 sec)

mysql> xa end 'c';
Query OK, 0 rows affected (0.00 sec)

mysql> xa prepare 'c';
Query OK, 0 rows affected (0.01 sec)
```

* 另一个客户端加锁查询

很明显，由于 xa 事务的原因，该查询超时了

```bash
mysql> select * from grades for update;
ERROR 1205 (HY000): Lock wait timeout exceeded; try restarting transaction
```

* kill 掉 mysql 服务端，模拟宕机场景。然后重启 mysql 服务端。

```bash
ayang@Ubuntu22:~$ ps -aux | grep mysql
mysql     420584  0.0  0.6 1624532 187852 ?      Sl   12:14   0:02 /usr/sbin/mysqld --daemonize --pid-file=/var/run/mysqld/mysqld.pid
ayang@Ubuntu22:~$ sudo kill -9 420584
ayang@Ubuntu22:~$ sudo systemctl start mysql
```

* 开启客户端并加锁查询

可以看到依旧超时了，说明 xa 事务还存在，即使宕机重启了。

```bash
mysql> select * from grades for update;
ERROR 1205 (HY000): Lock wait timeout exceeded; try restarting transaction
```

* 提交 xa 事务，再次查询

提交成功，查询立刻返回。说明 prepare，宕机重启后，依旧能够 commit

```bash
mysql> xa commit 'c';
Query OK, 0 rows affected (0.00 sec)

mysql> select * from grades for update;
Empty set (0.00 sec)
```

## End
