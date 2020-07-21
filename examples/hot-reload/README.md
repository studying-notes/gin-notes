# 用 Air 实现 Gin框架实时重新加载

```shell
go get -u github.com/cosmtrek/air
```

## 用法

1. 进入项目目录

```shell
cd /path/to/project
```

2. 在当前目录创建一个新的配置文件 `.air.conf`

3. 复制 `air.conf.example` 中的内容到这个文件，根据需要修改

```shell
# [Air](https://github.com/cosmtrek/air) TOML 格式的配置文件

# 工作目录
# `tmp_dir` 目录必须在 `root` 目录下
root = "."
tmp_dir = "tmp"

[build]
# 编译命令
cmd = "go build -o ./tmp/main ."
# 由 `cmd` 命令得到的二进制文件名
bin = "tmp/main"
# 自定义的二进制，可以添加额外的编译标识
full_bin = "APP_ENV=dev APP_USER=air ./tmp/main"
# 监听以下文件扩展名的文件
include_ext = ["go", "tpl", "tmpl", "html"]
# 忽略这些文件扩展名或目录
exclude_dir = ["assets", "tmp", "vendor", "frontend/node_modules"]
# 监听以下指定目录的文件
include_dir = []
# 排除以下文件
exclude_file = []
# 设置触发构建的延迟时间
delay = 1000 # ms
# 发生构建错误时，停止运行旧的二进制文件
stop_on_error = true
# 日志文件名，放置在 `tmp_dir` 中
log = "air_errors.log"

[log]
# 显示日志时间
time = true

[color]
# 自定义每个部分显示的颜色
main = "magenta"
watcher = "cyan"
build = "yellow"
runner = "green"

[misc]
# 退出时删除 tmp 目录
clean_on_exit = true
```

4. 输入 `air` 命令运行

```shell
air
```
