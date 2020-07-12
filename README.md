# Gin 学习笔记

## 目录

- [模块代理](#模块代理)
- [安装配置](#安装配置)
  - 设置 VS Code
- [快速开始](#快速开始)
- [Gin 特性](#Gin-特性)
- [通过 jsoniter 编译 JSON](#通过-jsoniter-编译-JSON)
- [测试用例](#测试用例)
- [官方示例](#官方示例)
  - [HTTP 基本方法](#HTTP-基本方法)
  - [URL 路径参数](#URL-路径参数)
  - [解析 URL 编码的 POST 表单](#解析-URL-编码的-POST-表单)
  - [POST 数据是映射格式](#POST-数据是映射格式)
  - [上传文件](#上传文件)
    - [单个文件](#单个文件)
    - [多个文件](#多个文件)
  - [路由组](#路由组)
  - [自定义中间件](#自定义中间件)
  - [自定义从错误中恢复的行为](#自定义从错误中恢复的行为)
  - [写入日志文件](#写入日志文件)
  - [自定义日志格式](#自定义日志格式)
  - [控制日志输出颜色](#控制日志输出颜色)

## 模块代理

国内网络基本上无法访问官方网址，所以设置代理是最好的办法。

```shell
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=off
go env -w GO111MODULE=on
```

## 安装配置

1. 下载安装 Gin

```shell
go get -u github.com/gin-gonic/gin
```

2. 初始化工作区

```shell
go mod init your_project
go mod edit -require github.com/gin-gonic/gin@latest
go mod vendor
```

3. 设置 VS Code

我个人写 Go 其实习惯于 Goland，但 VS Code 仍是我的最爱，只是看代码的话还是 VS Code 爽。首先必须安装 Microsoft 官方的 Go 插件，其次快捷键 `Ctrl + Shift + P` 进入命令模式，键入 `go:install/update tools` ，将全部插件选中，然后点击确定开始安装。

以下是 VS Code 与 Go 插件相关的全局设置，我一般就是这样设置的：

```json
{
  "go.useLanguageServer": true,
  "go.autocompleteUnimportedPackages": true,
  "go.gocodePackageLookupMode": "go",
  "go.gotoSymbol.includeImports": true,
  "go.useCodeSnippetsOnFunctionSuggest": true,
  "go.useCodeSnippetsOnFunctionSuggestWithoutType": true,
  "go.inferGopath": true,
  "go.docsTool": "gogetdoc",
  "go.formatTool": "goimports",
  "workbench.startupEditor": "newUntitledFile",
  "go.languageServerExperimentalFeatures": {
    "format": true,
    "autoComplete": true,
    "rename": true,
    "goToDefinition": true,
    "hover": true,
    "signatureHelp": true,
    "goToTypeDefinition": true,
    "goToImplementation": true,
    "documentSymbols": true,
    "workspaceSymbols": true,
    "findReferences": true,
    "diagnostics": true,
    "documentLink": true
  },
  // 只检查新加入代码
  "go.lintFlags": ["--enable-all", "--new"],
}
```

4. 通过以下方式导入模块

```go
import "github.com/gin-gonic/gin"
```

## 快速开始

1. 创建 `example.go` 文件

```go
package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.Run()
}
```

2. 命令行运行

```shell
# localhost:8080/ping
go run example.go
```

3. 访问 `localhost:8080/ping`

```json
{
    "message": "pong"
}
```

## Gin 特性

- 零分配路由
- 从路由到数据写入，仍是最快的 HTTP 框架
- 完备的单元测试套件
- 可经受生产环境挑战
- 稳定的 API，新的版本发布不影响之前的代码

## 通过 jsoniter 编译 JSON

默认标准库处理 JSON 数据，可以通过配置修改为更快的 [`jsoniter`](https://github.com/json-iterator/go)

```shell
go build -tags=jsoniter .
```

## 测试用例

HTTP 测试推荐 `net/http/httptest` 标准库，以 `example.go` 为例：

```go
package main

import "github.com/gin-gonic/gin"

func pingRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})
	return r
}
```

可编写如下测试用例：

```go
package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPingRoute(t *testing.T) {
	router := pingRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "pong", w.Body.String())
}
```

## 官方示例

### HTTP 基本方法

```go
func main() {
	router := gin.Default()

	router.GET("get", get)
	router.GET("post", post)
	router.GET("put", put)
	router.GET("delete", delete)
	router.GET("patch", patch)
	router.GET("head", head)
	router.GET("options", options)

	// 默认 8080 端口
	router.Run()
	// router.Run(":80")
}
```

### URL 路径参数

```go
func main() {
	router := gin.Default()

	router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "hello %s", name)
	})

	router.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + "is" + action
		c.String(http.StatusOK, message)
	})

	router.POST("user/:name/*action", func(c *gin.Context) {
		_ = c.FullPath() == "/user/:name/*action"
	})
	router.Run()
}
```

```shell
$ curl localhost:8080/user/rustlekarl
hello rustlekarl
```

### 解析字符串查询参数

```go
func main() {
	router := gin.Default()
	// 示例解析：/welcome?firstname=Jane&lastname=Doe
	router.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
    // 等价于 c.Request.URL.Query().Get("lastname")
    lastname := c.Query("lastname")

		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})
	router.Run(":8080")
}
```

```shell
$ curl --location --request GET 'localhost:8080/welcome?firstname=Jane&lastname=Doe'
Hello Jane Doe
```

### 解析 URL 编码的 POST 表单

```go
func main() {
	router := gin.Default()

	router.POST("/form_post", func(c *gin.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "anonymous")

		c.JSON(200, gin.H{
			"status":  "posted",
			"message": message,
			"nick":    nick,
		})
	})
	router.Run(":8080")
}
```

```shell
$ curl --location --request POST 'localhost:8080/form_post' --header 'Content-Type: application/x-www-form-urlencoded' --data-urlencode 'message=hello' --data-urlencode 'nick=rustlekarl'
{"message":"hello","nick":"rustlekarl","status":"posted"}
```

### 同时带有字符串查询参数和表单数据

```go
func main() {
	router := gin.Default()

	router.POST("/post", func(c *gin.Context) {

		id := c.Query("id")
		page := c.DefaultQuery("page", "0")
		name := c.PostForm("name")
		message := c.PostForm("message")
		c.String(http.StatusOK, "id: %s; page: %s; name: %s; message: %s", id, page, name, message)
	})
	router.Run(":8080")
}
```

```shell
$ curl --location --request POST 'localhost:8080/post?id=1&page=2' --header 'Content-Type: application/x-www-form-urlencoded' --data-urlencode 'message=hello' --data-urlencode 'name=rustlekarl'
id: 1; page: 2; name: rustlekarl; message: hello
```

### POST 数据是映射格式

```go
func main() {
	router := gin.Default()

	router.POST("/post", func(c *gin.Context) {

		ids := c.QueryMap("ids")
		names := c.PostFormMap("names")

		// c.String(http.StatusOK, "ids: %v; names: %v", ids, names)
		c.String(http.StatusOK, "ids[a]: %v; names[first]: %v", ids["a"], names["first"])
	})
	router.Run(":8080")
}
```

```shell
$ curl --location --request POST 'localhost:8080/post?ids[a]=1234&ids[b]=hello' --header 'Content-Type: application/x-www-form-urlencoded' --data-urlencode 'names[first]=rustle' --data-urlencode 'names[second]=karl'
ids[a]: 1234; names[first]: rustle
```

### 上传文件

#### 单个文件

```go
func main() {
	router := gin.Default()
	// 设置上传文件的最大内存限制
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		file, _ := c.FormFile("file")
		log.Println(file.Filename)

		dst := "upload/file"
		c.SaveUploadedFile(file, dst)

		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})
	router.Run(":8080")
}
```

```shell
$ curl --location --request POST 'http://localhost:8080/upload' 
--form 'file=@D:/LICENSE' --header "Content-Type: multipart/form-data"        
'LICENSE' uploaded!
```

#### 多个文件

```go
func main() {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20  // 8 MiB
	router.POST("/upload", func(c *gin.Context) {
		// Multipart form
		form, _ := c.MultipartForm()
		files := form.File["upload[]"]

		for _, file := range files {
			log.Println(file.Filename)
			c.SaveUploadedFile(file, file.Filename)
		}
		c.String(http.StatusOK, fmt.Sprintf("%d files uploaded!", len(files)))
	})
	router.Run(":8080")
}
```

```shell
$ curl --location --request POST 'http://localhost:8080/upload' --form 'upload[]=@D:/LICENSE' --form 'upload[]=@D:/README.rst'
2 files uploaded!
```

### 路由组

```go
func main() {
	router := gin.Default()

	// Simple group: v1
	v1 := router.Group("/v1")
	{
		v1.POST("/login", loginEndpoint)
		v1.POST("/submit", submitEndpoint)
		v1.POST("/read", readEndpoint)
	}

	// Simple group: v2
	v2 := router.Group("/v2")
	{
		v2.POST("/login", loginEndpoint)
		v2.POST("/submit", submitEndpoint)
		v2.POST("/read", readEndpoint)
	}

	router.Run(":8080")
}
```

### 自定义中间件

```go
func main() {
  // 创建一个不带任何中间件的路由
	r := gin.New()

	// 全局中间件
	// 日志中间件将日志写入 gin.DefaultWriter 即使设置了 GIN_MODE=release
	// 默认 gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

  // 恢复中间件从任何错误中恢复状态，返回 500 状态码
	r.Use(gin.Recovery())

	// Per route middleware, you can add as many as you desire.
	r.GET("/benchmark", MyBenchLogger(), benchEndpoint)

	// Authorization group
	// authorized := r.Group("/", AuthRequired())
	// exactly the same as:
	authorized := r.Group("/")
	// per group middleware! in this case we use the custom created
	// AuthRequired() middleware just in the "authorized" group.
	authorized.Use(AuthRequired())
	{
		authorized.POST("/login", loginEndpoint)
		authorized.POST("/submit", submitEndpoint)
		authorized.POST("/read", readEndpoint)

		// nested group
		testing := authorized.Group("testing")
		testing.GET("/analytics", analyticsEndpoint)
	}

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
```

### 自定义从错误中恢复的行为

```go
func main() {
	// Creates a router without any middleware by default
	r := gin.New()

	// Global middleware
	// Logger middleware will write the logs to gin.DefaultWriter even if you set with GIN_MODE=release.
	// By default gin.DefaultWriter = os.Stdout
	r.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	r.Use(gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	r.GET("/panic", func(c *gin.Context) {
		// panic with a string -- the custom middleware could save this to a database or report it to the user
		panic("foo")
	})

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "ohai")
	})

	// Listen and serve on 0.0.0.0:8080
	r.Run(":8080")
}
```

### 写入日志文件

```go
func main() {
    // Disable Console Color, you don't need console color when writing the logs to file.
    gin.DisableConsoleColor()

    // Logging to a file.
    f, _ := os.Create("gin.log")
    gin.DefaultWriter = io.MultiWriter(f)

    // Use the following code if you need to write the logs to file and console at the same time.
    // gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

    router := gin.Default()
    router.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })

    router.Run(":8080")
}
```

### 自定义日志格式

```go
func main() {
	router := gin.New()

	// LoggerWithFormatter middleware will write the logs to gin.DefaultWriter
	// By default gin.DefaultWriter = os.Stdout
	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// your custom format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
				param.ClientIP,
				param.TimeStamp.Format(time.RFC1123),
				param.Method,
				param.Path,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Request.UserAgent(),
				param.ErrorMessage,
		)
	}))
	router.Use(gin.Recovery())

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.Run(":8080")
}
```

### 控制日志输出颜色

```go
func main() {
    // Disable log's color
    gin.DisableConsoleColor()
    
    // Creates a gin router with default middleware:
    // logger and recovery (crash-free) middleware
    router := gin.Default()
    
    router.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })
    
    router.Run(":8080")
}
```

```go
func main() {
    // Force log's color
    gin.ForceConsoleColor()
    
    // Creates a gin router with default middleware:
    // logger and recovery (crash-free) middleware
    router := gin.Default()
    
    router.GET("/ping", func(c *gin.Context) {
        c.String(200, "pong")
    })
    
    router.Run(":8080")
}
```

### URL 路径参数

```go

```

### URL 路径参数

```go

```

### URL 路径参数

```go

```

