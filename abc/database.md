---
date: 2020-07-12T19:15:24+08:00  # 创建日期
author: "Rustle Karl"  # 作者

# 文章
title: "Gin 操作数据库"  # 文章标题
url:  "posts/gin/doc/database"  # 设置网页链接，默认使用文件名
tags: [ "gin", "go" ]  # 自定义标签
series: [ "Gin 学习笔记"]  # 文章主题/文章系列
categories: [ "学习笔记"]  # 文章分类

# 章节
weight: 20 # 文章在章节中的排序优先级，正序排序
chapter: false  # 将页面设置为章节

index: true  # 文章是否可以被索引
draft: false  # 草稿
---

官方提供了一个操作数据库的通用接口，但是必须提供对应数据库的驱动。

`Open`

```go
func Open(driverName, dataSourceName string) (*DB, error)
```

`Prepare`

```go
func (db *DB) Prepare(query string) (*Stmt, error)
```

`Stmt` 是一个准备好了的声明语句，多线程安全。

`Exec`

```go
func (s *Stmt) Exec(args ...interface{}) (Result, error)
```

`Result`

```go
type Result interface {
        // 数据库返回一个数值
        LastInsertId() (int64, error)

        // 语句影响的行的序号
        RowsAffected() (int64, error)
}
```

`Query`

```go
func (s *Stmt) Query(args ...interface{}) (*Rows, error)
```

`Rows`

查询返回的结果，游标指向结果集的首行。

```go
// Next 方法准备下一个结果行供 Scan 方法读取
func (rs *Rows) Next() bool

// Scan 复制当前行的全部列到 dest
func (rs *Rows) Scan(dest ...interface{}) error
```

## SQLite3

官网：

```
https://www.sqlite.org/
```

一般常用的驱动为：

```
go get -u github.com/mattn/go-sqlite3
```

简单示例：

```go
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
)

func main() {
	os.Remove("./foo.db")

	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// 创建数据表
	sqlStmt := `
		CREATE TABLE IF NOT EXISTS userinfo
		(
			uid        INTEGER PRIMARY KEY AUTOINCREMENT,
			username   VARCHAR(64) NULL,
			departname VARCHAR(64) NULL,
			created    DATE        NULL
		);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	// 事务：插入
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	stmt, err := tx.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec("Zhu", "Yue", fmt.Sprintf("VOL.%03d", i))
		if err != nil {
			log.Fatal(err)
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}

	// 查询整个表
	rows, err := db.Query("select username, departname from userinfo")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var username string
		var departname string
		err = rows.Scan(&username, &departname)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(username, departname)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	// 查询指定条件，返回一行
	stmt, err = db.Prepare("select created from userinfo where uid = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()
	var created string
	err = stmt.QueryRow("3").Scan(&created)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(created)

	// 删除
	//_, err = db.Exec("delete from userinfo")
	//if err != nil {
	//	log.Fatal(err)
	//}

	// 插入
	_, err = db.Exec("INSERT INTO userinfo(username, departname, created) values('foo', 'bar', 'baz'), ('foo2', 'bar2', 'baz2')")
	if err != nil {
		log.Fatal(err)
	}
}
```


## MySQL/MariaDB

```shell
go get "github.com/go-sql-driver/mysql"
```

```go
import _ "github.com/go-sql-driver/mysql"
```

```go
connStr := "root:12345678@tcp(127.0.0.1:3306)/ginsql"

db, err := sql.Open("mysql", connStr)
    if err != nil {
        log.Fatal(err.Error())
        return
    }
```

操作数据库就是原生 SQL。

```go
sql := "sql"
_, err = db.Exec(sql)
```

查询数据：

```go
rows, err := db.Query(sql)
```

## xorm

一般而言，通过框架而不是原生 SQL 操作数据库是更好的选择。

```shell
go get github.com/go-xorm/xorm
```
