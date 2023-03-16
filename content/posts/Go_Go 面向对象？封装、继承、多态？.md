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
title: "Go 面向对象？封装、继承、多态？"
date: 2023-03-17T00:59:45+08:00
lastmod: 2023-03-17T00:59:45+08:00
categories: ["Go"]
tags: ["面向对象"]
---

系统学习并整理了下 Go 的面向对象的思想，分享下。

参考：

* https://mp.weixin.qq.com/s/r5gIVtyBWtD7UQncK_JPNQ  
* 《Go 语言学习笔记》--作者：雨痕


## Go 的面向对象/过程

**面向对象**？（我的理解）  
函数不能独立存在，必须依托于类。面向对象把世界中的所有东西都抽象成类（对象），类具有属性（字段）和方法（函数）。如 Java Main.main。当然，面向对象还有三大特性：封装，继承，多态。

**面向过程**？（我的理解）  
函数可以单独存在，比如 Go 的包函数，可以不通过 类.Method **直接调用**。逻辑指令被划分为一个个函数调用，每个函数拥有自己的功能。  
缺点：重用性低，维护麻烦。（函数一多就特乱。）

但是 C++ 也支持直接调用函数，并且人家可是面向对象的语言呢。所以上面的理解应该半错把。  
面向对象应该更加注重的是那三个属性把：封装，继承，多态。（看下面一次讲解把）

**Go 是一门面向对象的语言么？**

是又不是。    
* Go 支持面向对象编程，但却不是通过提供经典的类、对象以及**类型层次**来实现的。
* Go 同时支持面向过程的编程，可以直接调用进行函数调用，而不需要经过类的调用。

## 封装

> 封装就是把数据以及操作数据的方法“打包”到一个抽象数据类型中，这个类型封装隐藏了实现的细节，所有数据仅能通过导出的方法来访问和操作。

**Go 的 Struct 可以起到抽象数据类型的效果**，可以拥有属性和方法。但是 Go 没有提供任何方法来隐藏属性，如 Java 可以通过 private 来表示字段私有，外部就不能通过 class.field 访问了。

注意：Go 的 Struct 在声明时可以通过大小写来标识包外部能否调用。

## 继承

###  组合大于继承

> 子类自动拥有了父类的属性和方法。  
> 具有明显的层级关系，被继承的称为父类，继承的称为子类。子类对象实例可以赋值给父类变量。

**Go 的匿名嵌套可以起到继承的效果**，但 Go 的匿名嵌套却没有层级关系，两个 Struct 完全是不同的 Struct 类型，根本没有任何关系，当然也不能赋值。以经典 OO 理论话术去理解就是两个 Struct 的关系不是 is-a，而是 has-a 的关系。

其实这种继承更应该被称为**组合**。Go 更愿意将模块分成互相独立的小单元，分别处理不同方面的需求，最后以**匿名嵌入**的方式组合到一起，共同实现对外接口。也就是**组合大于继承**的思想。

**组合没有父子依赖**，不会破坏封装。且整体和局部松耦合，可任意增加来实现扩展。各单元持有单一职责，互不关联，**自由灵活组合**，实现和维护更加简单。

### 匿名嵌套 

匿名嵌套在**编译时**会根据嵌套类型生成**包装方法**，**包装方法实际是调用嵌套类型的原始方法**。

{{< image src="/images/Go 面向Go_对象？封装、继承、多态如何实现？/包装方法.png" width=100% height=100% caption="包装方法" >}}

TIP：图片来源 B 站**幼麟实验室**，主讲 Go 语言的各种底层实现，特别推荐！

### 类似与继承的“重写”

匿名嵌套有**同名遮蔽问题**，编译器编译时会自动选择**深度最浅**的作为 Struct 方法集中的方法，类似与继承的重写。

```go
type A struct {
	B  
}

type B struct {
}

func (*A) M() {
	print("A")
}

func (B) M() {
	print("B")
}

func (B) M2() {
	print("B M2")
}

func main() {
	a := A{}
	a.M()   // A // 相当于重写 // 因为 A 也有 M 的同名方法，编译器生成的 A 的方法集中的 M 方法 是 A 的 M 方法
	a.B.M() // B
	a.M2()  // 相当于继承
}
```

## 多态

**Go 通过接口来实现多态。**

Go的接口类型本质就是一组方法集合(行为集合)，**一个类型如果实现了某个接口类型中的所有方法**，那么就可以作为动态类型赋值给接口类型（注意：定义一个接口变量，该变量本质是个 Struct 类型的**变量**哦）。

Go 的接口是特别重要的东西，通过学习 Go 接口的底层实现可以学到很多东西，例如**动态语言**的实现，**方法动态派发**实现等。这块下次再聊把。

## 拓：匿名嵌套的多种玩法

* struct 匿名嵌套 struct（上面已经展示过了）  
* interface 匿名嵌套 interface
* struct 匿名嵌套 interface  

### interface 匿名嵌套 interface

接口可嵌入其他匿名接口，**相当于将其声明的方法集导入**。

当然，注意只有实现了两个接口的全部的方法，才算实现大接口哈。

**Go 标准库中经典用法如下：**

```go
type Reader interface {
    Read(p []byte) (n int, err error)
}

type Writer interface {
    Write(p []byte) (n int, err error)
}

type ReadWriter interface {
   Reader
   Writer
}
```

### struct 匿名嵌套 interface

**编译器自动为 struct 的方法集加上 interface 的所有方法**。    
（后面是我猜的）如我们通过 struct.M() 调用 interface 的 M 方法，编译器实际为 struct 生成包装方法 `struct.interface.M()`。

注意：struct 中的 interface 要记得赋值哈，不然调用时会显示 interface nil panic。

```go
type I interface {
	M()
}

type A struct {
	I
}

type B struct {
}

func (B) M() {
	print("B")
}

func main() {
	var a A = A{I: B{}}
	a.M() // B

	// 当然 A 也是 I 接口类型
	var i I = A{I: B{}}
	i.M() // A
}
```

我们同时验证以下匿名嵌套的同名覆盖问题：

```go
type I interface {
	M()
}

type A struct {
	I
}

type B struct {
}

// 多加这个：验证匿名方法同名覆盖
func (A) M() {
	print("A")
}

func (B) M() {
	print("B")
}

func main() {
	var a A = A{I: B{}}
	a.M() // A
}
```

**Go 标准库中经典用法如下：**

context 包中：

```go
type valueCtx struct {
  Context  // 匿名接口
  key, val interface{}
}

// 创建 valueCtx
func WithValue(parent Context, key, val interface{}) Context {
  return &valueCtx{parent, key, val}
}

// 实际重写了 Value() 接口，其他父 context 的方法依旧可以调用
func (c *valueCtx) Value(key interface{}) interface{} {
  if c.key == key {
    return c.val
  }
  return c.Context.Value(key)
}
```