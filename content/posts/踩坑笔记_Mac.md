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
title: "Mac 使用"
date: 2024-03-20T01:12:45+08:00
lastmod: 2024-03-20T01:12:45+08:00
categories: ["踩坑笔记"]
tags: ["Hugo", "Markdown"]
---

新的实习，发现好多 Mac 的操作都忘记了，又得重新查一遍。记录一下把，下次入职会更快熟练

## Mac

### 配置环境

首先有两种 shell（命令解释器），bash 和 zsh，macOS 默认是 zsh

如何确定是哪一种 bash： `echo $0` 

两者加载的环境配置文件不同（包括顺序）：（挺乱的，所以这里很简）

```shell
# bash
/etc/profile   # 全局配置，每一个用户登录
~/.bash_profile # 用户配置
```

```shell
# zsh 强大很多，有很多配置和一定的加载顺序，这里只列了两个
/etc/zshenv # 全局配置，每一个用户登录
~/.zshenv # 用户配置
```

**最省心的方法：**   
配置全部写在 /etc/profile，然后在 /etc/zshenv 文件里加上 `source /etc/profile`

### Mac 快捷键

#### 自定义快捷键方法

例如打开命令行居然没有快捷键，可以自己设一个，`Command + T`

ChatGPT：

1. 打开 "系统偏好设置"。  
2. 选择 "键盘"。  
3. 切换到 "快捷键" 选项卡。  
4. 选择 "应用程序快捷键"。  
5. 点击左下角的 "+" 号来添加新的快捷键。在 "应用程序" 中选择 "终端"。  
6. 在 "菜单标题" 中输入 "新建窗口" 或 "新建标签"。  
7. 输入你想要的快捷键。  
8. 点击 "添加"。  

#### 系统定义快捷键  

* 进入和退出全面屏 `Command + Control + F`  
* 锁屏 `Control + Command + Q`  
* 复制黏贴文本 `Command + C` `Command + V`  
* 多屏幕多窗口  
  * 固定窗口设置 https://blog.csdn.net/genius_yym/article/details/81508654  
  * 向左/右切换一个屏幕 `control + >` `control + >`   
  * 切换到某一个屏幕 `control + n`  
  * 不同应用的窗口切换 `command +tab`  

## End
