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
title: "通过 go 汇编了解函数调用栈帧"
date: 2022-11-06T01:17:04+08:00
lastmod: 2022-11-06T23:59:04+08:00
categories: ["底层"]
tags: ["go 汇编"]
---

因为大二上准备跟大三先修《操作体统》，所以暑假的时候花一周时间过了一遍王爽老师的《汇编语言》，学习了 8086 实模式下的汇编语言相关知识。虽然已经忘得差不多，不过基础框架还依稀记得。真的，再次感觉到基础知识的重要性：学习其他知识会快很多。  

本文章通过对 go 的两个简单函数进行反汇编，讲解 go 函数调用的过程，顺便提一下尾递归优化。

## 前置知识

Go 使用了 plan9 汇编，所以有必要了解一下。（tip：必须看，或者边看后面的汇编时边查看）  

* <a href="https://github.com/cch123/asmshare/blob/master/layout.md" target="_blank">https://github.com/cch123/asmshare/blob/master/layout.md</a>
* <a href="https://github.com/cch123/golang-notes/blob/master/assembly.md" target="_blank">https://github.com/cch123/golang-notes/blob/master/assembly.md</a>
* <a href="https://www.bilibili.com/video/BV1WZ4y1p7JT/" target="_blank">https://www.bilibili.com/video/BV1WZ4y1p7JT/</a>

go 的函数调用栈帧大致如下：   

{{< image src="/images/通过 go 汇编了解函数调用栈帧/go函数调用栈帧.jpg" width=60% height=60% caption="go函数调用栈帧" >}}

**生成汇编的命令**：  `go tool compile main.go`    
该命令可以编译 Go 文件生成汇编代码，-N 参数表示禁止编译优化， -l 表示禁止内联，-S 表示打印汇编。

## 详解函数调用栈过程

```go
package main

func swap(a, b int) int {   
	a, b = b, a             
	c := a + b             
	return c               
}

func main() {               
	a := 1                 
	b := 2                 
	c := swap(a, b)         
	print(c)                
}
```

```shell
# 记得开启禁止编译优化和禁止内联优化
go tool compile -S -N -l demo2.go
```

```shell
"".main STEXT size=95 args=0x0 locals=0x30 funcid=0x0 align=0x0
        0x0000 00000 (demo2.go:9)       TEXT    "".main(SB), ABIInternal, $48-0 # 48 表示函数栈帧大小占 48Byte，在编译时确认
        # (接上面）go 比较特殊，不是在函数执行过程中逐步扩展栈帧大小，而是直接确认栈帧大小，然后一口气分配（把 SP 指针移到需要的位置）
        0x0000 00000 (demo2.go:9)       CMPQ    SP, 16(R14)
        0x0004 00004 (demo2.go:9)       PCDATA  $0, $-2
        0x0004 00004 (demo2.go:9)       JLS     88
        0x0006 00006 (demo2.go:9)       PCDATA  $0, $-1
        0x0006 00006 (demo2.go:9)       SUBQ    $48, SP #（1）SP-48，即把 SP 向低地址移动48字节作为 main 栈的栈顶
        0x000a 00010 (demo2.go:9)       MOVQ    BP, 40(SP) #（2）把 BP 的值保存到 SP+40，即 main 栈帧的栈底
        0x000f 00015 (demo2.go:9)       LEAQ    40(SP), BP #（3）把 SP+40 的地址保存到 BP 作为 main 栈的栈基
        0x0014 00020 (demo2.go:9)       FUNCDATA        $0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x0014 00020 (demo2.go:9)       FUNCDATA        $1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x0014 00020 (demo2.go:10)      MOVQ    $1, "".a+32(SP) #（4）赋值，把 1 赋值给 SP+32
        0x001d 00029 (demo2.go:11)      MOVQ    $2, "".b+24(SP) #（5）赋值，把 2 赋值给 SP+24
        0x0026 00038 (demo2.go:12)      MOVQ    "".a+32(SP), AX #（6) 把 a 的值放入 AX 寄存器中
        0x002b 00043 (demo2.go:12)      MOVL    $2, BX #（7）把 b 的值放入 BX 寄存器中
        0x0030 00048 (demo2.go:12)      PCDATA  $1, $0
        0x0030 00048 (demo2.go:12)      CALL    "".swap(SB) #（8）把下一条指令的地址入栈（隐含SP-1），并跳到 swap 函数的地址执行代码
        0x0035 00053 (demo2.go:12)      MOVQ    AX, "".c+16(SP) #（21）通过寄存器把 swap 的返回值存储在 main.SP+16 
        0x003a 00058 (demo2.go:13)      CALL    runtime.printlock(SB)
        0x003f 00063 (demo2.go:13)      MOVQ    "".c+16(SP), AX
        0x0044 00068 (demo2.go:13)      CALL    runtime.printint(SB)
        0x0049 00073 (demo2.go:13)      CALL    runtime.printunlock(SB)
        0x004e 00078 (demo2.go:14)      MOVQ    40(SP), BP
        0x0053 00083 (demo2.go:14)      ADDQ    $48, SP
        0x0057 00087 (demo2.go:14)      RET
        0x0058 00088 (demo2.go:14)      NOP
        0x0058 00088 (demo2.go:9)       PCDATA  $1, $-1
        0x0058 00088 (demo2.go:9)       PCDATA  $0, $-2
        0x0058 00088 (demo2.go:9)       CALL    runtime.morestack_noctxt(SB)
        0x005d 00093 (demo2.go:9)       PCDATA  $0, $-1
        0x005d 00093 (demo2.go:9)       JMP     0
```

```shell
"".swap STEXT nosplit size=90 args=0x10 locals=0x20 funcid=0x0 align=0x0
        0x0000 00000 (demo2.go:3)       TEXT    "".swap(SB), NOSPLIT|ABIInternal, $32-16 # swap 函数栈帧大小为 32Byte
        0x0000 00000 (demo2.go:3)       SUBQ    $32, SP #（9）SP-32，即把 SP 向低地址移动32字节作为 swap 栈的栈顶
        0x0004 00004 (demo2.go:3)       MOVQ    BP, 24(SP) #（10）把 BP 的值保存到 SP+24，即 swap 栈帧的栈底
        0x0009 00009 (demo2.go:3)       LEAQ    24(SP), BP #（11）把 SP+24 的地址保存到 BP 作为 swap 栈的栈基
        0x000e 00014 (demo2.go:3)       FUNCDATA        $0, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x000e 00014 (demo2.go:3)       FUNCDATA        $1, gclocals·33cdeccccebe80329f1fdbee7f5874cb(SB)
        0x000e 00014 (demo2.go:3)       FUNCDATA        $5, "".swap.arginfo1(SB)
        0x000e 00014 (demo2.go:3)       MOVQ    AX, "".a+40(SP) #（12）通过寄存器传参 a，注意是 main 栈的内存空间
        0x0013 00019 (demo2.go:3)       MOVQ    BX, "".b+48(SP) #（13）同上
        0x0018 00024 (demo2.go:3)       MOVQ    $0, "".~r0(SP) #（14）把 swap 的栈顶赋值为 0（不懂原因，也用不到）
        0x0020 00032 (demo2.go:4)       MOVQ    "".a+40(SP), CX #（15）15 都是交换a，b的值
        0x0025 00037 (demo2.go:4)       MOVQ    CX, ""..autotmp_4+16(SP)  #（15）把 CX 赋值给SP+16
        0x002a 00042 (demo2.go:4)       MOVQ    "".b+48(SP), CX #（15）
        0x002f 00047 (demo2.go:4)       MOVQ    CX, "".a+40(SP) #（15）
        0x0034 00052 (demo2.go:4)       MOVQ    ""..autotmp_4+16(SP), CX #（15）
        0x0039 00057 (demo2.go:4)       MOVQ    CX, "".b+48(SP) #（15）
        0x003e 00062 (demo2.go:5)       MOVQ    "".a+40(SP), DX #（16）16 都是对a，b相加，并赋值给 c 和 寄存器 AX
        0x0043 00067 (demo2.go:5)       LEAQ    (DX)(CX*1), AX #（16）
        0x0047 00071 (demo2.go:5)       MOVQ    AX, "".c+8(SP) #（16）
        0x004c 00076 (demo2.go:6)       MOVQ    AX, "".~r0(SP) #（17）把 返会值 c 值存入 栈顶（不知道意义何在，没用到）
        0x0050 00080 (demo2.go:6)       MOVQ    24(SP), BP #（18）把之前存储 swap 栈底的 main.BP 重新存入 BP
        0x0055 00085 (demo2.go:6)       ADDQ    $32, SP #（19）SP-32，释放 swap 函数栈帧
        0x0059 00089 (demo2.go:6)       RET  #（20）弹出栈顶（隐含SP+1）（当时保存了CALL指令的下一条指令的地址，并跳到该地址继续执行）
```

{{< image src="/images/通过 go 汇编了解函数调用栈帧/汇编执行过程.png" width=100% height=100% caption="汇编执行过程" >}}

## 尾递归优化

尾递归学习了：
* <a href="https://zhuanlan.zhihu.com/p/36587160" target="_blank">https://zhuanlan.zhihu.com/p/36587160</a>

下面是直接 copy 作为我的笔记：（建议直接看作者原文）  

**普通递归**  

```c
function fact(n) {
    if (n <= 0) {
        return 1;
    } else {
        return n * fact(n - 1);
    }
}

函数递归调用展开：
6 * fact(5)
6 * (5 * fact(4))
6 * (5 * (4 * fact(3))))
// two thousand years later...
6 * (5 * (4 * (3 * (2 * (1 * 1)))))) // <= 最终的展开

展开后回溯计算：
6 * (5 * (4 * (3 * (2 * 1)))))
6 * (5 * (4 * (3 * 2))))
6 * (5 * (4 * 6)))
// two thousand years later...
720 // <= 最终的结果
```

**尾递归优化**

```c
function fact(n, r) {
    if (n <= 0) {
        return 1 * r;
    } else {
        return fact(n - 1, r * n);
    }
}

fact(6, 1) // 1 是 fact(0) 的值，我们需要手动写一下
fact(5, 6)
fact(4, 30)
fact(3, 120)
fact(2, 360)
fact(1, 720)
720 // <= 最终的结果
```

**尾递归定义**  
> 若函数在尾位置调用自身（或是一个尾调用本身的其他函数等等），则称这种情况为尾递归。尾递归也是递归的一种特殊情形。尾递归是一种特殊的尾调用，即在尾部直接调用自身的递归函数。对尾递归的优化也是关注尾调用的主要原因。尾调用不一定是递归调用，但是尾递归特别有用，也比较容易实现。  
>
> 特点：  
> 尾递归在普通尾调用的基础上，多出了2个特征:
> 1. 在尾部调用的是函数自身 (Self-called)；  
> 2. 可通过优化，使得计算仅占用常量栈空间 (Stack Space)。

**函数栈的作用**   
栈的意义其实非常简单，五个字——保持入口环境

**尾递归为什么可以优化**  
尾递归，可以把函数栈的入口环境其实是无意义的，所有可以优化掉。

```c
function fact(n, r) {
    if (n <= 0) {
        return 1 * r;
    } else {
        return fact(n - 1, r * n); // <= 这里的入口环境没有必要保留。
    }
}
```
当里面这个 fact(n - 1, r * n) 返回的时候，外面的 fact(n, r) 就马上要返回了，所以保存栈是没有任何意义的，既然没意义我们毫无疑问就要优化掉。

**尾递归是编译时的优化**  
例如 go 目前并没有对尾递归进行优化，所有调用尾递归还是会 Stack Overflow，而 C 则不会。

## End
