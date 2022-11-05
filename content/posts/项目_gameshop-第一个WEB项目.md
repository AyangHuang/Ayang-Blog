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
title: "gameshop - 第一个 WEB 项目"
date: 2022-10-31T01:30:49+08:00
lastmod: 2022-10-31T01:30:49+08:00
categories: ["项目"]
tags: ["JavaWeb"]
---

gameshop 项目是大一暑期实训个人独立完成（除 html、css 外）的小项目，也是我的的第一个 WEB 项目。15 天的时间学习到了很多很多，可以说是除了参与学校网络创新实验室一二轮考核外短时间内收获知识最丰富的一段时间（汗，因为都有 ddl 的 push），收获了 WEB 的整体工作流程，了解并实践 JavaWeb 的开发。包括了 jquery，ajax，session，cookie，JavaEE 中的 Servlet，Filter 等和数据库连接池，ThreadLocal 等等。这篇文章主要是回顾整理下项目的简单实现（不涉及代码细节）。  

项目用到的技术栈：JavaWeb、jQuery、ajax、MySQL。

实训结束后简单学习了 Linux，已经部署到云服务器上，请访问 <a href="https://ayang.ink/gameshop" target="_blank">ayang.ink/gameshop</a>  
Github 源码地址：<a href="https://github.com/AyangHuang/gameshop-javaweb" target="_blank">GitHub</a>  

## 前期准备

因为问过去年的学长得知暑期实训是做 WEB 项目，所以在期末前 3 周我是看了两本书《网络是怎么连接起来的》、《图解HTTP》（推荐这俩本书，计网的简单入门），粗略地、全局地了解整个网络中数据的传输过程。我个人觉得**理论知识是特别重要的**，是实践的前提。但是实际实践过程中，还是会碰壁然后重新看书，在实践过程中也更加了解了 session、cookie（尽管对现在的我来说，so easy）。

## 项目整体概括

>1. ~~项目使用前后端分离开发方式~~（并不是），利用**ajax**发送异步请求，以**json**交互数据
>2. 明文密码使用**MD5**的方式进行传输，服务端采用**MD5加盐**的方式进行二次加密
>3. 运用**cookie和session**进行用户**登录状态管理**，运用**ServletContext**存储全局用户登录信息，**HttpSessionListener**监听**session**并利用**ServletContext**管理用户登录信息：（1）实现**7天自动登录**；（2）实时监视**网站登录总人数**；（3）不允许一个账号同时在不同浏览器登录，实现一方登录，另一方**强制下线**功能。
>4. 运用**session**实现**购物车**
>5. 运用**filter**进行过滤请求，防止未登录用户访问需要登录的页面和获取登录后的信息等非法访问
>6. 数据库封装**BaseDao**类：（1）运用apache的**dbutils**工具类；（2）运用阿里的**Druid连接池**；（3）数据库**可选择是否开启事务**，配合**ThreadLocal+filter**进行一个请求线程内的事务管理。
>7. 前端使用**jquery+ajax**进行**实时交互渲染**。例：游戏**搜索功能**的实现，动态搜索并显示。
>8. **分页功能**的实现和其**可变性**（可通过参数简单设置分页细节）

## 1. 前后端分离？

项目整体是 HTML + CSS 编写静态网页，然后通过 ajax 向服务端发送异步请求获取动态数据进行网页的进一步渲染。我以为这样（只要不用 jsp，而是用 ajax）就是**前后端分离**，直到一两周前，和同学交流的时候才知道**前后端分离是前端和后端分成两个项目**（简单说，就是两个互不相干的文件夹并行开发）。具体再深入的我也没有了解，等我实践过再来补把。

## 2.搜索游戏功能的实现

{{< image src="/images/gameshop-第一个WEB项目/search.png" width=100% height=100% caption="搜索功能" >}}

1. 先定义一个监听事件监听输入框文字的改变，只要改变就发送 ajax 向服务器请求
2. 运用 MySQL 进行模糊搜索，并返回 json 数据
3. 每次收到响应后，先删除下拉框的全部内容，然后根据返回的数据重新把下拉框拼接上去。

## 3. 密码传输和存储：MD5 + 盐

> **MD5 消息摘要算法**，一种被广泛使用的密码散列函数，可以根据文本产生出一个128位（16字节）的散列值（hash value），用于确保信息传输完整一致。其典型应用就是对一段信息产生信息摘要，防止被篡改。
>
> 1. **抗修改性**：对原数据进行任何改动，哪怕只修改1个字节，所得到的MD5值都有很大区别。
>
> 2. **不是加密算法**，理论是不能解密的，但是彩虹表，字典等可破解常见密码。解决方案：通过**加salt**。
>
> 3. 加**随机salt**：**MD5（MD5（密码）+随机字符串（salt））**。

具体业务逻辑如下：  

{{< image src="/images/gameshop-第一个WEB项目/login.png" width=100% height=100% caption="注册流程" >}}

{{< image src="/images/gameshop-第一个WEB项目/sign_up.png" width=100% height=100% caption="登录流程" >}}

注意实际生产的时候，salt 和 最后密码是放在两个不同的数据库中，减少被同时入侵导致被破解的可能性。

## 4. 登录状态管理

### cookie 和 session

> **cookie**：数据存储在浏览器中，由服务器在发送 HTTP Respond 报文时通过 `Set-Cookie` 响应头发送给浏览器。  
> **session**：数据存储在服务器中（准确来说应该是由 WEB 程序在内存中维护，其实就是一个数据结构，比如是 map 之类）。在第一次访问服务器时，WEB 程序会生成一个 session 数据结构，并把标识该 session 的 sessionid 通过 session 发送回浏览器。下次访问时，cookie 带上 sessionid。WEB 程序通过 sessionid 找到对应的 session 数据结构，可以进程存储或删除数据。注意：session 是会话级别，浏览器关闭时会清楚 sessionid，session 在 WEB 程序同时也有超时机制，超过一定时间没访问，会删除该 session。  

{{< image src="/images/gameshop-第一个WEB项目/request.png" width=100% height=100% caption="请求头" >}}

{{< image src="/images/gameshop-第一个WEB项目/respond.png" width=100% height=100% caption="响应头" >}}

### ServletContext

> ServletContext：一个 WEB 工程，有一个 ServletContext 实例对象。ServletContext 对象是一个域对象，也就是有一个 map 属性，可以存储数据。全局需要用到的数据可以存储在这里。 

session 是通过 sessionid 来标识的，并没有关联用户。所以我另外通过在 ServletContext 中定义一个全局 ConcurrentHashMap（线程安全问题，不能用 map）来作为用户登录表记录处于登录状态的账户，key 是 `usr_id`，value 是 sessionid。然后就可以实现用户登录管理了，如记录网站实时访问人数、不允许一个账号在两处浏览器登录，强制下线一方等。

### (1) 第一次登录

{{< image src="/images/gameshop-第一个WEB项目/first_login.png" width=100% height=100% caption="第一次登录" >}}

### (2) 自动登录例子 

每次打开页面会发送ajax请求进行自动登录验证，验证成功则自动登录

{{< image src="/images/gameshop-第一个WEB项目/auto_login.png" width=100% height=100% caption="自动登录" >}}

### (3) 监视网站实时登录人数

1. **ServletContext**储存用户登录表，key 是 `user_id`，value 是 **sessionid**

2. session 超时**自动销毁**或**强制销毁**（重复登录），**HttpSessionListener**监听session销毁并从用户登陆表移除该键值对（不监听销毁，会导致 session 被销毁，但是用户登录表却还指向这里，会造成野指针）。

## 下面摆烂了======================

## 5. ThreadLocal + Filter 进行事务管理

（后序会写 ThreadLocal 的源码刨析，实训时看过一点，好像不难？）

1. **ThreadLocal**是解决线程安全问题一个很好的思路，它通过为每个线程提供一个独立的变量副本解决了变量并发访问的冲突问题。

   该类提供了线程局部 (thread-local) 变量。这些变量不同于它们的普通对应物，因为访问某个变量(通过其 get 或 set 方法)的每个线程都有自己的局部变量，它独立于变量的初始化副本。

2. **Tomacat**（多线程） 对于每一个请求都会开启一个**新的线程**来处理（运用线程池），故可用 ThreadLocal 存储 Connection 来进行事务管理，使一个请求里的**Connection是同一个**。

3. filter 类似 go 中间件功能，在发生 error 时，可利用 filter 回滚事务

## 6. Filter 进行请求过滤

1. 前端会验证是否存在sessionid ，存在则发送请求，不存在则跳到登录页面

2. 规定请求包含用户信息的请求都携带**/user/**，后端 Filter 对 /user/ 进行session验证，判断是否处于登录状态，是则放行。

3. 后端的过滤，能有效避免跳过前端验证直接手动发送请求恶意获取信息的方式。

## 7.分页功能实现

可根据需要修改两个参数即可灵活改变分页参数。

## 8.session 实现购物车
## End
