---
date: 2020-07-12T19:15:24+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "分布式限流器"  # 文章标题
url:  "posts/gin/project/limiter"  # 设置网页链接，默认使用文件名
tags: [ "gin", "go" ]  # 自定义标签
series: [ "Gin 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 文章分类

# 章节
weight: 20 # 文章在章节中的排序优先级，正序排序
chapter: false  # 将页面设置为章节

index: true  # 文章是否可以被索引
draft: false  # 草稿
---

## 简介

限流会导致用户在短时间内（这个时间段是毫秒级的）系统不可用，一般我们衡量系统处理能力的指标是每秒的 QP 或者 TPS，假设系统每秒的流量阈值是 1000，理论上一秒内有第 1001 个请求进来时，那么这个请求就会被限流。

## 限流方案

### 计数器

```go
package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

// 普通计数限流器

type CountLimiter struct {
	counter  int64         // 计数器
	max      int64         // 最大数量
	interval time.Duration // 间隔时间
	last     time.Time     // 上一次时间
}

func NewCountLimiter(max int64, interval time.Duration) *CountLimiter {
	return &CountLimiter{
		counter:  0,
		max:      max,
		interval: interval,
		last:     time.Now(),
	}
}

func (c *CountLimiter) Allow() bool {
	current := time.Now()
	// 超过时间计数清零
	if current.After(c.last.Add(c.interval)) {
		atomic.StoreInt64(&c.counter, 1)
		c.last = current
		return true
	}
	// 取出一个
	atomic.AddInt64(&c.counter, 1)
	// 判断是否超过限流个数
	if c.counter <= c.max {
		return true
	}
	return false
}

func main() {
	limiter := NewCountLimiter(3, time.Second)
	for i := 0; i < 10; i++ {
		for limiter.Allow() {
			fmt.Println(i)
		}
		time.Sleep(time.Second)
	}
}
```

### 漏桶算法

我们把水比作是请求，漏桶比作是系统处理能力极限，水先进入到漏桶里，漏桶里的水按一定速率流出，当流出的速率小于流入的速率时，由于漏桶容量有限，后续进入的水直接溢出（拒绝请求），以此实现限流。

![](../imgs/leaky_bucket.jpg)

```go

```

```go

```

### 令牌桶算法

系统会维护一个令牌（token）桶，以一个恒定的速度往桶里放入令牌（token），这时如果有请求进来想要被处理，则需要先从桶里获取一个令牌（token），当桶里没有令牌（token）可取时，则该请求将被拒绝服务。令牌桶算法通过控制桶的容量、发放令牌的速率，来达到对请求的限制。

![](../imgs/token_bucket.jpg)

```go

```


```go

```


```go

```

