# 操作数据库

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
