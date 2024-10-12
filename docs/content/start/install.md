---
title: "快速安装"
weight: 1
# bookFlatSection: false
# bookToc: true
# bookHidden: false
# bookCollapseSection: false
bookComments: true
---

## 快速安装

### 安装Go

首先需要安装Go(支持版本1.20+), 你可以在
{{< button href="https://golang.org/dl" >}} Go官网{{< /button >}}
或者{{< button href="https://golang.google.cn/dl">}}Go中国镜像站{{< /button >}}
(国内用户推荐)找到对应的安装包

### 下载和导入项目

本项目使用了`GO111MOUDULE`, 你需要使用Go模块集成管理(如果不会就用Goland帮你解决吧)

使用下面指令下载源码

```bash
 go get github.com/wdvxdr1123/ZeroBot
```

然后你就可以在项目中导入ZeroBot了!

```golang
import "github.com/wdvxdr1123/ZeroBot"
```
