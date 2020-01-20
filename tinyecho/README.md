an echo service

TOTO
+ [√]TCP协议支持
+ [√]通信消息格式定义
+ [√]服务端消息分发
+ [√]客户端交互输入
+ []客户端控制台输入、输出优化
+ []服务端随机唯一名字生成器
+ []客户端网页版
+ [√]UDP协议支持
+ []配置文件及平滑加载

## 编译

- GOOS：目标平台的操作系统（darwin、freebsd、linux、windows） 
- GOARCH：目标平台的体系架构（386、amd64、arm） 
- 交叉编译不支持 CGO 所以要禁用它

Mac 下编译 Linux 和 Windows 64位可执行程序

```
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
```

Linux 下编译 Mac 和 Windows 64位可执行程序

```
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build main.go
```