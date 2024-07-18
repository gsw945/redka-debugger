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

redka-debugger.exe -d ".\data\record.db" -v "load-offset"
```

### examples

#### list keys
- command: `redka-debugger.exe -d ".\data\record.db" -k "*"`
- output:
   ```
   2024/07/18 10:41:53 reddka debugger
   key(ID=1, Key=load-offset, Type=String(1), Version=21764, ETime=<nil>, MTime=1721211182488)
   key(ID=2, Key=kafka-669789e8d791d51339cad7fa, Type=Hash(4), Version=13096, ETime=<nil>, MTime=1721211182488)
   ```

#### get value by key
- command: `redka-debugger.exe -d ".\data\record.db" -v "load-offset"`
- output:
   ```
   2024/07/18 10:42:01 reddka debugger
   key(ID=1, Key=load-offset, Type=String(1), Version=21764, ETime=<nil>, MTime=1721211182488)
   value: 33803
   ```

