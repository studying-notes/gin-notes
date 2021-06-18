---
date: 2021-01-25T13:34:55+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "HTTP Status Code 详解"  # 文章标题
url:  "posts/gin/abc/status_code"  # 设置网页永久链接
tags: [ "gin", "http"]  # 标签
series: [ "Gin 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 文章分类

# 章节
weight: 20 # 排序优先级
chapter: false  # 设置为章节

index: true  # 是否可以被索引
toc: true  # 是否自动生成目录
draft: false  # 草稿
---

200 OK (from disk cache) 是浏览器没有跟服务器确认， 就是直接用浏览器缓存。

304 是浏览器和服务器确认了一次缓存有效性，再用的缓存。
