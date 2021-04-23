
redis rdb database file convter.

```
# usage:
rdbconvter -d path/to/name.rdb

# eg:
rdbconvter -d ../tinyredislib/data/dump.rdb
```

## 编译说明

Mac 下编译 Linux 和 Windows 64位可执行程序

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o rdbconvter main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o rdbconvter main.go
```

Linux 下编译 Mac 和 Windows 64位可执行程序

```
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o rdbconvter main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o rdbconvter main.go
```

编译选项说明：
- GOOS：目标平台的操作系统（darwin、freebsd、linux、windows） 
- GOARCH：目标平台的体系架构（386、amd64、arm） 
- 交叉编译不支持 CGO 所以要禁用它