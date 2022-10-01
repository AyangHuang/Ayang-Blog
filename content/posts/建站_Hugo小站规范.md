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
title: "Hugo小站规范"
date: 2022-09-29T01:12:45+08:00
lastmod: 2022-10-01T15:58:26+08:00
categories: ["建站"]
tags: ["Hugo", "Markdown"]
---

做事得有规范和计划，所以趁着小站刚刚搭建，先制定一下该网站的目录规范和md文件书写规范等，也把该page文件作为以后page文件的模板。

## 规范

### 文件规范  

1. page命名规范  

    ~~每一个page文件命名为“category_category_pageTitle_lastMod.md”。  
    且category排序按“字母顺序>中文拼音顺序”。  
    eg: 该page命名为“建站_Hugo小站规范_20220929”。~~   

    20221001更改为：  
    每一个page文件命名为“category_category_pageTitle”。  
    且category排序按“字母顺序>中文拼音顺序”。  
    eg: 该page命名为“建站_Hugo小站规范”。  

    更改原因：  
    修改文件名称即导致URI改变，个人觉得不太好，遂改变。  

2. image存储规范  

    1. image统一存储在“site/static/images/”中，经hugo转化后存储在“public/images/”中。 

    2. 一个page对应一个文件夹，该文件夹命名为“pageTitle”  
    
    eg：该page中有一张命名为“图片.png”的image，该image存储路径为“site/static/images/Hugo小站规范/图片.png”，引入路径为“/images/Hugo小站规范/图片.png”。

### 目录规范  

pageTitle是一级目录，且**一级目录后是该page的简介**。  

主体文章从二级目录开始书写（theme.conf 设置从二级目录开始解析)  

```toml
  # 目录设置
  # 推荐：文章的标题为一级目录，目录从二级目录开始
  [markup.tableOfContents]
    startLevel = 2
    endLevel = 6
```

### 图片引入规范

图片用HTML模式引入，可设置**居中显示**和**大小**。

```HTML
<div align=center> <img src="/images/Hugo小站规范/图片.png" style="width:50%; height:50%"/> </div>
```

或者用Hugo shortcode（去掉\，防止被解析）
```Go HTML Template
\{\{< image src="/images/Hugo小站规范/展示写作.png" width=50% height=50% caption="我是下面的文字" >\}\}
```

### Markdown规范

相同的Markdown文件在不同软件渲染出来的可能会有所差异，为了兼容大部分软件且使得Hugo解析成HTML利于观看，遂制定如下规范。

1. 换行规范  

    记住换行时一定要（兼容）：**在本行尾按两下space再按enter键进行换行**。  

2. 空行规范  

    注意：在vscode编辑的时候无论有多少个空行(只要这一行只有回车或者space没有其他的字符就算空行)，**在渲染之后，只隔着一行**。也就是说无论如何只能空一行（当然在Typora即写即渲染可以空很多行）。  

3. 段落规范  
   
    规范：只要不在同一段文字中，用空行进行分段切割。

    注意：在Markdown语言中，**唯一决定两行文字是否是段落的，就在于这两行文字之间是否有空行**。如果这两行文字之间,有空行了，就代表，这两行文字为两个段落，如果这两行文字之间，没有空行，仅仅换行，就代表这两行文字是属于同一个段落。即使是在一行文字中的末尾，添加了两个空格之后换行，这两个行文字依旧是一个段落。  

    段落的作用：**增大间隙，体现层次感，增强观感**。

    无段落和有段落对比如下：

    {{< image src="/images/Hugo小站规范/无段落.png" caption="无段落">}}
    {{< image src="/images/Hugo小站规范/有段落.png"
    caption="有段落">}}

4. 缩进规范  

    只要是利用**数字**和**小圆点***进行分割的，下方文字统一缩进。且缩进统一用Tab（四个space）。  

## 写作工具  

我的note写作工具现在是**vscode+本地网页实时渲染**。  

实时渲染记得设置

```bash
hugo server --disableFastRender
```

{{< image src="/images/Hugo小站规范/展示写作.png">}}


## Markdown源码  

该page的md源码可点击左下角的**阅读原始文档**下载，然后在vscode即可清晰观看到**该page Markdown的规范**。

## End