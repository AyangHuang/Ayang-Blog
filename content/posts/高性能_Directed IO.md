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
title: "Directed IO"
date: 2022-11-27T00:56:45+08:00
lastmod: 2022-11-27T00:56:45+08:00
categories: ["高性能"]
tags: ["DIO"]
---

学习一下 IO 的两种方式。

## 参考资料

* 理论知识：《现代操作系统 原理与实现》    
* <a href="https://mp.weixin.qq.com/s/gW_3JD52rtRdEqXvyg-lJQ" target="_blank">https://mp.weixin.qq.com/s/gW_3JD52rtRdEqXvyg-lJQ</a>

## block    

磁盘**一次物理读写**的基本单位是**扇区**（一般 512Byte)。扇区的空间比较小且数目众多，在寻址时比较困难，操作系统的虚拟文件文件系统就将**多个扇区**组合在一起，形成一个更大的单位，就是**块**（block）。虚拟文件系统通过**块**（一般为 4KB）作为读取等操作数据的**基本单位**（即文件系统读写的最小粒度为**块**）。   

{{< image src="/images/Directed-IO/block.png" width=40% height=40% caption="block" >}}

**总结**  

1. 扇区是对硬盘而言，是**物理层**的；块是对虚拟文件系统而言，是**逻辑层**的。   
2. 磁盘可以看成是一个由 block 组成的数组。

## 页缓存（Page Cache）  

### 页缓存

Page Cache（页缓存）是**位于内核地址空间**中，以内存页为单位，**缓存数据块**。  
当一个文件被读取时，文件系统首先会先检查其内容是否已经保存在页缓存中。如果文件数据已经保存在页缓存中，则文件系统直接从页缓冲中读取数据返回给应用程序；否则，文件系统会在页缓存中创建新的内存页，并从存储设备中读取相关的数据，然后将其保存在创建的内存页中。之后，文件系统同样会再次检查内存页中找到并读取相应的页数据，返回给应用程序。

### 预读机制

操作系统为基于 Page Cache 的读缓存机制提供预读机制（PAGE_READAHEAD），一个例子是：  
用户线程仅仅请求读取磁盘上文件 A 的 offset 为 0-3KB 范围内的数据，由于文件系统的的基本读写单位为 block（4KB），于是操作系统至少会读 0-4KB 的内容，这恰好可以在一个 page 中装下。
但是操作系统出于局部性原理会选择将磁盘块 offset [4KB,8KB)、[8KB,12KB) 以及 [12KB,16KB) 都加载到内存，于是额外在内存中申请了 3 个 page。

### 延迟写入

大多数现代操作系统将**写入**在内存中缓冲 5-30s，即只是修改内存中的数据，并且把该内存页设置为**脏页**。系统中存在定期任务（表现形式为内核线程），周期性地将文件系统中文件脏数据块写回磁盘（即**异步 Write Back 机制**）。将内存写回磁盘时间延长，则可以通过批处理写入磁盘，减少 I/O 次数，提高性能。  
如果对于数据要求较高的，可以利用系统调用`fsync(int fd)`：将 fd 代表的文件的脏数据和脏元数据全部刷新至磁盘中。

### 直接 I/O 和 缓存 I/O

根据**是否利用操作系统的页缓存**（page cache），可以把文件 I/O 分为直接 I/O 与缓存 I/O。  

{{< image src="/images/Directed-IO/DIO.png" width=80% height=80% caption="缓存 IO 和 DIO 对比" >}}


* **缓存 I/O**（标准I/O）：读操作时，数据先从磁盘 copy 到**内核页缓存**中，然后再从内核页缓存中拷贝给用户程序，写操作时，数据从用户程序拷贝给内核缓存，再由内核决定什么时候写入数据到磁盘。（缓存I/O又被称作**标准I/O**，大多数文件系统的默认I/O操作都是缓存I/O。） 
* **直接 I/O**（Direct I/O，DIO）：直接IO就是应用程序直接访问磁盘数据，而不经过内核缓冲区，也就是绕过内核缓冲区，自己管理IO缓存区。  
    应用程序可以在打开文件时，通过附带的`O_DIRECT`标志，提示文件系统不要使用页缓存。

**缓存I/O**  

* 优点：在一定程度上分离了内核空间和用户空间，保护系统本身的运行安全；可以减少 I/O的次数，从而提高性能。
* 缺点：数据多次拷贝，性能降低（一些应用程序（如数据库）会自己实现缓存机制对数据进行缓存和管理，此时，操作系统的缓存是冗余的）；数据缓存在内核空间中，一段时间再写入磁盘，断电丢失。  

**直接I/O**  

* 优点：数据写的时候直接写回磁盘，确保掉电不丢失；减少内核 page cache 的内存使用，业务层自己控制内存，更加灵活。

## DIO 实践

**直接 I/O 对齐操作**    
DIO 模式由于自己维护缓冲区（即直接从磁盘 copy 到自定义的缓冲区），需要程序自己**保证对齐规则**，否则 IO 会报错。   
为什么对齐的是扇区大小？我猜的原因：因为物理磁盘读写的最小粒度是扇区。

* 用于传递数据的缓冲区，其大小必须和扇区大小（一般 512Byte）（即 IO 的大小是扇区大小的倍数）。
* 数据传输的开始点，即文件的读写偏移位置（offset）必须是扇区大小的倍数。
* 用于传递数据的缓冲区的地址（内存边界）必须与扇区大小对齐（eg：对齐 512 扇区，则虚拟地址的后 9 位 必须全部是 0）。

**代码**  
```go
package dio

import (
	"errors"
	"unsafe"
)

var (
	sectorSize = 512
)

func AlterSectorSize(size int) {
	sectorSize = size
}

func align(buf []byte) int {
	//return int(uintptr(unsafe.Pointer(&buf[0])) % uintptr(sectorSize))
	// 上下结果相同，位运算会快一点
	return int(uintptr(unsafe.Pointer(&buf[0])) & uintptr(sectorSize-1))
}

// NewDioBuf bufSize 必须是 sectorSize 的整数倍
func NewDioBuf(bufSize int) ([]byte, error) {
	if bufSize%sectorSize != 0 {
		panic("缓冲区的大小必须和扇区大小对齐")
	}
	buf := make([]byte, bufSize+sectorSize)
	offset := sectorSize - align(buf)
	if offset != 0 {
		buf = buf[offset : offset+bufSize]
	} else {
		buf = buf[:bufSize]
	}
	// 再判断一次
	if judgeAlign(buf) {
		return buf, nil
	}
	return buf, errors.New("err")
}

func judgeAlign(buf []byte) bool {
	if len(buf)%sectorSize == 0 {
		if align(buf) == 0 {
			return true
		}
	}
	return false
}
```

## End
