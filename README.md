# redka-debugger
debug tool for redka

## 一、项目搭建记录
```bash
# 项目结构搭建
mkdir redka-debugger
cd redka-debugger/
go mod init redka-debugger
# Go包 ldap 安装
go get -u modernc.org/sqlite
go get -u github.com/nalgeon/redka
touch main.go
```

### 二、设置Go包源加速和安装Go包
```bash
# 启用 Go Mod
set GO111MODULE=on
# 使用国内源
set GOPROXY=https://proxy.golang.com.cn,direct
# 安装Go包
go mod tidy -v
```

### (三)、Windows 编译&运行
```bash
# 编译
go build -o redka-debugger.exe main.go
# 运行
redka-debugger.exe -h
# example
redka-debugger.exe -d "d:\proj\SBTProject\SBTOperationTools\logbus-parser\build\record.db" -k "*"
redka-debugger.exe -d "d:\proj\SBTProject\SBTOperationTools\logbus-parser\build\record.db" -v "load-offset"
```
