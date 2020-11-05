---
date: 2020-07-12T19:15:24+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Gin 源码解析"  # 文章标题
url:  "posts/gin/abc/source"  # 设置网页链接，默认使用文件名
tags: [ "gin", "go" ]  # 自定义标签
series: [ "Gin 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 文章分类

# 章节
weight: 20 # 文章在章节中的排序优先级，正序排序
chapter: false  # 将页面设置为章节

index: true  # 文章是否可以被索引
draft: false  # 草稿
---

## 路由详解

Gin 用的是定制版本的 [httprouter](https://github.com/julienschmidt/httprouter)，其路由的原理是大量运用公共前缀的树结构，它基本上是一个紧凑的 [Trie tree](https://baike.sogou.com/v66237892.htm) 或者是 [Radix Tree](https://baike.sogou.com/v73626121.htm)。具有公共前缀的节点也共享一个公共父节点。

### Radix Tree

基数树（Radix Tree）又称为 PAT 位树（Patricia Trie or Crit Bit Tree），是一种更节省空间的前缀树（Trie Tree）。

对于基数树的每个节点，如果该节点是唯一的子树的话，就和父节点合并。下图为一个基数树示例：

![](../imgs/radix_tree.png)

`Radix Tree` 可以被认为是一棵简洁版的前缀树。我们注册路由的过程就是构造前缀树的过程，具有公共前缀的节点也共享一个公共父节点。假设我们现在注册有以下路由信息：

```go
r := gin.Default()

r.GET("/", func1)
r.GET("/search/", func2)
r.GET("/support/", func3)
r.GET("/blog/", func4)
r.GET("/blog/:post/", func5)
r.GET("/about-us/", func6)
r.GET("/about-us/team/", func7)
r.GET("/contact/", func8)
```

那么我们会得到一个 `GET` 方法对应的路由树，具体结构如下：

```
Priority   Path             Handle
9          \                *<1>
3          ├s               nil
2          |├earch\         *<2>
1          |└upport\        *<3>
2          ├blog\           *<4>
1          |    └:post      nil
1          |         └\     *<5>
2          ├about-us\       *<6>
1          |        └team\  *<7>
1          └contact\        *<8>
```

上面最右边那一列每个 `*<数字>` 表示 Handle 处理函数的内存地址。从根节点遍历到叶子节点我们就能得到完整的路由表。

例如：`blog/:post` 其中 `:post` 只是实际文章名称的占位符。与 `hash-maps` 不同，这种树结构还允许我们使用像 `:post` 参数这种动态部分，因为我们实际上是根据路由模式进行匹配，而不仅仅是比较哈希值。

由于 URL 路径具有层次结构，并且只使用有限的一组字符，所以很可能有许多常见的前缀。这使我们可以很容易地将路由简化为更小的问题。此外，**路由器为每个请求方法管理一棵单独的树**。一方面，它比在每个节点中都保存一个 `method-> handle` Map 更加节省空间，它还使我们甚至可以在开始在前缀树中查找之前大大减少路由问题。

为了获得更好的可伸缩性，每个树级别上的子节点都按 `Priority` 排序，其中优先级（最左列）就是在子节点中注册的句柄的数量。这样做有两个好处:

1. 首先优先匹配被大多数路由路径包含的节点。这样可以让尽可能多的路由快速被定位。
2. 类似于成本补偿。最长的路径可以被优先匹配，补偿体现在最长的路径需要花费更长的时间来定位，如果最长路径的节点能被优先匹配（即每次拿子节点都命中），那么路由匹配所花的时间不一定比短路径的路由长。

### 路由树节点

路由树是由一个个节点构成的，Gin 框架路由树的节点由 `node` 结构体表示，它有以下字段：

```go
// tree.go
type node struct {
   // 节点路径，比如上面的 s，earch，和 upport
	path      string
	// 和 children 字段对应, 保存的是分裂的分支的第一个字符
	// 例如 search 和 support, 那么 s 节点的 indices 对应的 "eu"
	// 代表有两个分支, 分支的首字母分别是 e 和 u
	indices   string
	// 儿子节点
	children  []*node
	// 处理函数链条（切片）
	handlers  HandlersChain
	// 优先级，子节点、子子节点等注册的 handler 数量
	priority  uint32
	// 节点类型，包括 static, root, param, catchAll
	// static: 静态节点（默认），比如上面的 s，earch 等节点
	// root: 树的根节点
	// catchAll: 有 * 匹配的节点
	// param: 参数节点
	nType     nodeType
	// 路径上最大参数个数
	maxParams uint8
	// 节点是否是参数节点，比如上面的 :post
	wildChild bool
	// 完整路径
	fullPath  string
}
```

### 请求方法树

在 Gin 的路由中，每一个 `HTTP Method`（GET、POST、PUT、DELETE）都对应了一棵  `radix tree`，注册路由的时候会调用下面的 `addRoute` 函数：

```go
// gin.go
func (engine *Engine) addRoute(method, path string, handlers HandlersChain) {
   // 获取请求方法对应的树
	root := engine.trees.get(method)
	if root == nil {
		root = new(node)
		root.fullPath = "/"
		engine.trees = append(engine.trees, methodTree{method: method, root: root})
	}
	root.addRoute(path, handlers)
}
```

暂时还没研究深入的必要性：

```go
https://www.liwenzhou.com/posts/Go/read_gin_sourcecode/
```
