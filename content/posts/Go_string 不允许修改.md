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
title: "string 不允许修改"
date: 2023-02-24T00:31:29+08:00
lastmod: 2023-02-24T00:31:29+08:00
categories: ["Go"]
tags: []
---

这是重新学习 Go string 类型时做的笔记，同时解决了我长期对于 string 类型不可修改的困惑。

## string struct

```go
// str 变量是一个指向 stringStruct 的指针
str := "string"

type stringStruct struct {
	str unsafe.Pointer // 指向底层的内存空间的起始位置
	len int  // 表示内存空间的大小
}
```

## 不允许修改有两层意思  

1. 不允许 `str[1] = 'o'`，编译器编译不通过。即 **string 类型的变量在编译层面不允许被修改**。  

2. 在编译时**字符串字面量**分配于 SRODATA，该内存只能读取不能修改。即**存储字面量 string 的内存不允许修改**（当然如果 string 指向的底层是堆、栈等可读可写内存，是可以通过unsafe指针方式强制修改的）。  
   ```asm
   # func main() {
   #	s := "123456"
   #	s = s + "7"
   #	print(s)
   #}
   
   # 只有字符串字面量在编译时就分配好内存，且位于 RODATA 只读内存段
   go.string."123456" SRODATA dupok size=6
           0x0000 31 32 33 34 35 36                                123456
   go.string."7" SRODATA dupok size=1
           0x0000 37                                               7
   ```

## 不允许修改的原因 

1. go 实现中，string struct 不包含字符串实际内存空间，只有一个指向内存的指针。这样做的好处是 string 变得非常轻量，可以很方便地进行传递而不用担心内存拷贝。  
  
2. 保证对底层字符串的并发安全。

## 举例不是字面量时如何修改 string 底层内存空间  

```go
package main

import "unsafe"

func main() {
	str := "123456"
	
	// 直接报错，因为底层内存在编译时就分配好，在 SRODATA 只读段
	//change(s)

	// 此时在运行时动态分配在堆区或栈区
	str = str + "7" 

	change(str)

	print(str) // 0234567
}

func change(str string) string {
	// 直接变成切片指针，然后再解引用变成切片，然后直接修改
	slices := *(*([]byte))(unsafe.Pointer(&str))
	slices[0] = '0'
	return *((*string)(unsafe.Pointer(&slices)))
}
```

## End
