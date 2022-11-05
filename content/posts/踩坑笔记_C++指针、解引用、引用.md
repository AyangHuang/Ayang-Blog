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
title: "C++ 指针、解引用、引用"
date: 2022-11-05T19:52:23+08:00
lastmod: 2022-11-05T19:52:23+08:00
categories: ["踩坑笔记"]
tags: ["指针"]
---

昨天同学来问我一个 bug，尴尬死，我在他面前调试了好久。后来还是直接看我大一时候的笔记。这篇文章算是对昨天的总结把。

## 过程

### 同学 bug 代码

同学的代码（为了易读，修改了一下），是关于前序建立一颗二叉树的。

```c++
#include <iostream>
#include<string.h>
#include<stack>
using namespace std;

struct Node {
    char data;
    struct Node* lchild, * rchild;
};

void creatBinTree(Node* root)  {
    char ch;
    cin >> ch;
    if (ch == '#') {
        root = NULL;
    } else {
        root = new(Node);
        root->data = ch;
        cout << ch << "的左子树为：";
        creatBinTree(root->lchild);
        cout << ch << "的右子树为：";
        creatBinTree(root->rchild);
    }
}

int main()
{
    Node* root = NULL;
    cout << "请输入根结点: ";
    creatBinTree(root);
    system("pause");
    return 0;
}
```

其实这种链表和二叉树的建立问题，我之前在大一的时候刚刚学数据结构，也是踩过类似的坑。所以我三下五除二，很快就找出问题所在：main 函数传入 root 指针变量到 `creatBinTree` 里面对 root 赋值修改，根本不会影响到外面的（main）的 root 变量。  

### 糊涂修改

于是我凭着感觉，一顿操作猛如虎，做出下面的修改：

```c++
#include <iostream>
#include<string.h>
#include<stack>
using namespace std;

struct Node {
    char data;
    struct Node* lchild, * rchild;
};

void creatBinTree(Node &root)  {
    char ch;
    cin >> ch;
    if (ch == '#') {
        root.data = '#'; // 由于不是指针变量，所以无法赋值为 NULL，直接赋值为 #
    } else {
        root = *new(Node);
        root.data = ch;
        cout << ch << "的左子树为：";
        creatBinTree(*root.lchild);
        cout << ch << "的右子树为：";
        creatBinTree(*root.rchild);
    }
}

int main()
{
    Node* root = new(Node);
    cout << "请输入根结点: ";
    creatBinTree(*root);
    system("pause");
    return 0;
}
```

可以说，这代码我改得稀巴烂，基本哪里波浪线爆红警告我就改哪里，缺乏思考。出现了很多处致命错误！

**第一个错误**    
因为分析出 creatBinTree 函数内的变量无法改变 main 函数 root 的值，所以我一个想法是把函数参数改成引用类型（&），实际一个引用类型参数实际就是一个指针（后面再详细讲），所以相当于错误没改。

**第二错误**  
`root = *new(Node);`，把 new 出来在堆区的全新变量解引用赋值给 root 变量，本质就是把一个 new 出来未初始化的变量 copy 给 root。what，神操作，啥也没用，甚至造成**内存泄漏**：堆区的变量没指针指向了，但没释放。  

**第三个错误**  

```c++
root = *new(Node);
root.data = ch;
cout << ch << "的左子树为：";
creatBinTree(*root.lchild);
```

最后一行，root 此时的 *lchild 指针变量并没有初始化，然后解引用再用按引用的方式传递参数，后面 creatBinTree 会对 *lchild 进行赋值，很明显访问了**野指针**。   

## 前序建立二叉树正确代码

前序建立二叉树也是递归遍历，递归回溯的很好的代码实践。

```c++
int main() {
    // 直接调用下面的函数即可
    TreeNode *root = builtBiTree();
}
```
```c++
// 法一：常规方法
class Solution {
    //可以改成二级指针
    void recursion(TreeNode* &cur){  
        char ch;
        cin >> ch;   
        if (ch=='#') {
            cur = NULL;
        } else {
            cur = new TreeNode;
            cur->val = ch;
            recursion(cur->lchild);
            recursion(cur->rchild);
        }
    }
    TreeNode* builtBiTree() {
        TreeNode* root = NULL;
        recursion(root)
        return root;
    }
};
```

```c++
// 法二：回溯的时候利用返回值进行左右孩子的赋值
class Solution {
    TreeNode* recursion() {
        char ch;
        cin >> ch;   
        if (ch=='#') {
            return  NULL;
        } else {
            TreeNode* cur = new TreeNode(ch, NULL, NULL);
            cur->lchild = recursion();
            cur->rchild = recursion();
            return cur;
        }
    }
    TreeNode* builtBiTree() {
        return recursion();
    }
};
```

## 总结复习

虽然很基础，但是却很重要。每次用指针时都应该仔细思考。

### 指针

常说的指针只是**普通变量**，存储了一个地址（该地址是进程地址空间的地址，是一个虚拟地址）。  

例如 `int* a` 表示的是 a 是一个指针变量，存储了一个地址，地址指向内存存储的是一个 int 类型。根据 int 类型的大小为 8字节（32位操作系统），我们就可以知道从该起始地址占用的内存空间。

### 解引用

* **作为 = 左值**  
    对指针变量存储的地址指向的变量赋值。  
    
    ```c++
    int *a = new int; // 让变量 a 存储堆区 new 出来 8 字节的 int 变量的地址
    *a = 2; // 对指针变量 a 存储的地址指向的堆区的int 变量赋值为 2
    ```

* **作为 = 的右值**  
    把指针变量存储的地址指向的变量 copy 给右值。

    ```c++
    // 接上面
    int b = *a; // 把 a 指针变量存储的地址指向的堆区的int变量的值copy给变量 b
    ```

### 引用类型

引用（&）本质是指针常量。  
```c++
int b  = 2;
int& a = b;
// 上面一条等价于下面一条
int* const a = &b;
```

注：我感觉就是一个语法糖，反正我不太喜欢，直接用指针就可以了。语义更加明确。

## End
