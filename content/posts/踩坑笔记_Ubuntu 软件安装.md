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
title: "Ubuntu 软件安装"
date: 2022-12-16T21:16:23+08:00
lastmod: 2022-12-16T21:16:23+08:00
categories: ["踩坑笔记"]
tags: ["Ubuntu"]
---

从今天开始，正式步入 Linux 系统。起因是之前的 dio 系统调用 Windows 不支持，就下定决心后面有时间一定把开发环境转移到 Linux。因为疫情原因提前回家放寒假罗，回家第一天当然是安装 Ubuntu 双系统，然后配置开发环境罗。在这里做个记录，把 Ubuntu 下载软件的方式记录下来。下次如果重装系统就不会耗费太长时间了。

## vim 

由于ubuntu预安装的是tiny版本（键位会错乱），所以会导致我们在使用上的产生上述的不便。但是，我们安装了vim的full版本之后，键盘的所有键在vi下就很正常了。

```bash
sudo apt remove vim-common
sudo apt install vim
```

## netstat 

`sudo apt install net-tools`

## vscode

### vscoe

1. 添加微软公钥  
   
   `wget -q https://packages.microsoft.com/keys/microsoft.asc -O- | sudo apt-key add -`

2. 安装依赖文件   
   
    `sudo apt install software-properties-common apt-transport-https wget`

3. 将vscode的apt源添加到到本地  
   
    `sudo add-apt-repository "deb [arch=amd64] https://packages.microsoft.com/repos/vscode stable main"`

4. 安装vscode  
     
    ```bush
    sudo apt update 
    sudo apt install code
    ```


### vscode plugin

* markdown pre

## edge

直接到官网下载 .deb 安装包 

`sudo dpkg -i  edge.deb`

百度百科：deb 格式是 Debian 系统（包含Debian 和 Ubuntu）专属安装包格式，配合 APT 软件管理系统，成为了当前在 Linux 下非常流行的一种安装包。

## typora

太早下载了，忘记记录了。

## ToolBox 

到官网下载 .tar.gz 文件

```bash
sudo tar -C /usr/local/ -zxvf jetbrains-toolbox-1.27.1.13673.tar.gz
cd /usr/local/jetbrains-toolbox-1.27.1.13673 
./jetbrains-toolbox # 打开后加入菜单快捷页面
```

然后直接下载 jetbrains 全家桶

## go


[go 语言中文网](https://studygolang.com/dl) 找到 .tar.gz 下载地址

```bash
cd /home/ayang/Downloads
wget https://studygolang.com/dl/golang/go1.19.4.linux-amd64.tar.gz
tar -C /usr/local/ -zxvf go1.19.4.linux-amd64.tar.gz
```

环境变量配置

```bash
mv /etc/profile /etc/profile.cp
vim /etc/profile
export PATH=$PATH:/usr/local/go/bin
source /etc/profile
```

改镜像源和自定义 gopath
```bash
ayang@Ubuntu22:~$ go env -w GOPROXY="https://goproxy.cn"
ayang@Ubuntu22:~$ go env -w GOPATH="/home/ayang/gopath"
```

## java 

官网下载 .tar.gz 文件。

https://www.oracle.com/java/technologies/downloads/archive/

```bash
ayang@Ubuntu22:~/Downloads$ sudo mkdir /usr/local/java
ayang@Ubuntu22:~/Downloads$ sudo tar -C /usr/local/java -zxvf jdk-8u341-linux-x64.tar.gz
ayang@Ubuntu22:~/Downloads$ mv /etc/profile etc/profile.cp
ayang@Ubuntu22:~/Downloads$ sudo vim /etc/profile
export PATH=$PATH:/usr/local/java/jdk1.8.0_341/bin
ayang@Ubuntu22:~/Downloads$ source /etc/profile

```

## xmind2020

下载 .deb 安装包

https://dl2.xmind.cn/XMind-2020-for-Linux-amd-64bit-10.3.1-202101132117.deb


```bash
ayang@Ubuntu22:~/Downloads$ wget https://dl2.xmind.cn/XMind-2020-for-Linux-amd-64bit-10.3.1-202101132117.deb
ayang@Ubuntu22:~/Downloads$ sudo mkdir /usr/local/xmind
ayang@Ubuntu22:~/Downloads$ sudo dpkg -i Xmind-for-Linux-amd64bit-22.10.0920.deb
```

链接: https://pan.baidu.com/s/1ftNSdsLf3aZTLu53Hvwtow?pwd=idc5 提取码: idc5  

解压 XMind_2020_10.3.1_Linux_补丁.7z 后，将app.asar文件覆盖 /opt/XMind/resources/app.asar

## git 

```bash
ayang@Ubuntu22:~$ sudo apt install git
ayang@Ubuntu22:~$ git config --global user.name "ayang-linux"
ayang@Ubuntu22:~$ git config --global user.email ayang@Ubuntu22:~$ ssh-keygen
```

生成 ssh 公钥和私钥默认存储在 ~/.ssh/ 目录下

```bash
ayang@Ubuntu22:~$ ssh-keygen
```

## site

```bash
cd ~/Documents
mkdir site & cd site
git init 
git clone git@github.com:AyangHuang/Ayang-Blog.git
git submodule init
git submodule update
```

## hugo

```bash
sudo apt install hugo
```

## clash

https://github.com/Fndroid/clash_for_windows_pkg/releases

```bash
ayang@Ubuntu22:~/Downloads$ sudo mkdir /opt/clash
[sudo] password for ayang: 
ayang@Ubuntu22:~/Downloads$ tar -C /opt/clash -zxvf Clash.for.Windows-0.20.10-x64-linux.tar.gz 
ayang@Ubuntu22:~/Downloads$ cd /opt/clash
ayang@Ubuntu22:~/Downloads$ ./cfw
```

## 截图软件

windows 下肯定是 snipaste，可惜没开发 Linux 版本。

`sudo apt install flameshot`

设置 f1 快捷建

setting -> keyboard -> keyboard shortcuts -> custom shortcuts -> add 

command 填 `flameshot gui`，命令行启动 flameshot。

## MySQL5.7

```bash
wget https://downloads.mysql.com/archives/get/p/23/file/mysql-server_5.7.39-1ubuntu18.04_amd64.deb-bundle.tar
sudo dpkg -i mysql-*.deb
# 报错了缺少依赖，然后我就执行下面的，还是报错了
sudo apt install libaio1 && sudo apt install libmecab2
# 然后我跟着报错指示，我执行了下面不知道是啥的指令
sudo apt --fix-broken install
# 然后成了，显示输入密码界面
```

## redis 

`sudo apt install redis-server`  

`sudo apt install redis-tools`

## End