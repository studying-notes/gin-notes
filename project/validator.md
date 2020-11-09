---
date: 2020-08-12T19:15:24+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Gin 传入参数校验"  # 文章标题
url:  "posts/gin/project/validator"  # 设置网页链接，默认使用文件名
tags: [ "gin", "go", "validator" ]  # 自定义标签
series: [ "Gin 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 文章分类

# 章节
weight: 20 # 文章在章节中的排序优先级，正序排序
chapter: false  # 将页面设置为章节

index: true  # 文章是否可以被索引
draft: false  # 草稿
---

```shell
go get github.com/go-playground/validator/v10
```

Gin 框架默认设置 `TagName` 为 `binding`。

```go
func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		config := &validator.Config{TagName: "binding"}
		v.validate = validator.New(config)
	})
}
```

## 官网

```shell
https://github.com/go-playground/validator
```

## 标签详解

### 示例

```go
`binding:"required,min=2,max=100" example:"我的博客"`       
```

### 限制字符串长度

2 <= length < 100

```go
`binding:"required,min=2,max=100" example:"我的博客"`       
```

### 数值比较

- `eq`: equal，等于；
- `ne`: not equal，不等于；
- `gt`: great than，大于；
- `gte`: great than equal，大于等于；
- `lt`: less than，小于；
- `lte`: less than equal，小于等于；

### 嵌套结构验证

`dive` 一般用在 slice、array、map、嵌套的 struct 验证中，作为分隔符表示进入里面一层的验证规则。

```go
type Test struct {
	Array []string `validate:"required,gt=0,dive,max=2"`
    // gt=0 代表 field.len()>0

	Map map[string]string `validate:"required,gt=0,dive,keys,max=10,endkeys,required,max=2"`
	// dive 表示进入里面一层的验证，例如 a={"b":"c"} 中 dive 之前的 required 表示 a 是必填项，大于0，
	// dive 之后出现 keys 与 endkeys 之间代表验证 map 的 keys 的 tag 值：max=10，即长度不大于 10
	// 后面的是验证 value，required 必填并且最大长度是 2
}
```

### 必填项

`required` 表示字段必须有值，并且不为默认值，例如 `bool` 默认值为 `false`、`string` 默认值为 `””`、`int` 默认值为 `0`。如果有些字段是可填的，并且需要满足某些规则的，那么需要使用 `omitempty`。

```go
type Test struct {
	Id         string `form:"charger_id" validate:"omitempty,uuid4"`
    // Gin 验证 URL 中参数的 tag为 form
    // omitempty 表示变量可以不填，但是填的时候必须满足条件
	Password   string `form:"password" validate:"omitempty,min=5,max=128"`
}
```

### 自定义枚举值

oneof 自定义枚举值。

```go
type Test struct {
	Gender uint8  `json:"gender" binding:"oneof=0 1 2"`
}
```

### 字段间关系

- `eqfield=Field`: 必须等于 Field 的值；
- `nefield=Field`: 必须不等于 Field 的值；
- `gtfield=Field`: 必须大于 Field 的值；
- `gtefield=Field`: 必须大于等于 Field 的值；
- `ltfield=Field`: 必须小于 Field 的值；
- `ltefield=Field`: 必须小于等于 Field 的值；
- `eqcsfield=Other.Field`: 必须等于 struct Other 中 Field 的值；
- `necsfield=Other.Field`: 必须不等于 struct Other 中 Field 的值；
- `gtcsfield=Other.Field`: 必须大于 struct Other 中 Field 的值；
- `gtecsfield=Other.Field`: 必须大于等于 struct Other 中 Field 的值；
- `ltcsfield=Other.Field`: 必须小于 struct Other 中 Field 的值；
- `ltecsfield=Other.Field`: 必须小于等于 struct Other 中 Field 的值；
- `required_with=Field1 Field2`: 在 Field1 或者 Field2 存在时，必须；
- `required_with_all=Field1 Field2`: 在 Field1 与 Field2 都存在时，必须；
- `required_without=Field1 Field2`: 在 Field1 或者 Field2 不存在时，必须；
- `required_without_all=Field1 Field2`: 在 Field1 与 Field2 都存在时，必须；

## Gin 增加的标签

### Content-Type

```go
var (
	JSON          = jsonBinding{}
	XML           = xmlBinding{}
	Form          = formBinding{}
	Query         = queryBinding{}
	FormPost      = formPostBinding{}
	FormMultipart = formMultipartBinding{}
	ProtoBuf      = protobufBinding{}
	MsgPack       = msgpackBinding{}
)
```

常用的就 form、query、json 和 xml；经过测试 url 中的传参（application/x-www-form-urlencoded）和 body 中传参都可以使用 form 标签

```go
type Login struct {
	User     string `form:"user" json:"user" xml:"user"  binding:"required"`
	Password string `form:"password" json:"password" xml:"password" binding:"required"`
}
```

### 字段的默认值

```go
type Login struct {
	User     string `form:"user,default=admin"  binding:"required"`
	Password string `form:"password"  binding:"required"`
}
```
