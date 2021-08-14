

代码是写给人看的，理论上讲，我们不期望把代码规范定的太死，任何人都可以按照自己的习惯进行编码，但是当这种现象如果放大到一个项目、团队、组织中来看，就会导致后人维护代码成本过高等诸多问题。所以，我们制定了本开发指南，以期能在一个项目、团队、组织中达成统一的共识，形成一个统一的风格。

## 编码规范

### 总规范

todo::...

参考 Golang 官方文档：

+ [Effective Go](https://golang.org/doc/effective_go)
+ [Effective Go 中文版](https://github.com/bingohuang/effective-go-zh-en)

### 代码格式化

使用 goimports 工具进行代码格式化。

注：goimports 会额外修正 import 代码块

### 命名

变量、函数、结构体、结构体方法、文件的命名应尽量有意义，对于一些明显的拼写错误，发现后，应及时
修正。

### 模块化

+ 对于可以抽象为 1 个函数/方法的相同逻辑，抽象为 1 个函数/方法；
+ 对于不可以抽象为 1 个函数/方法的逻辑块，用 1 个空行与其他代码逻辑块分开；
+ 对于逻辑相似的 N 个函数/方法，建议考虑通过函数生成器的方式抽象为 1 个函数/方法。

### 返回值

除去约定的入口函数/方法，其他的函数/方法须有一个 error 返回。

### 错误的判断

优先判断错误逻辑，并提前返回，正常的代码放在后面，对正常代码保持最少的缩进。

注：这也是 Golang 的设计思想之一

### 接口请求

+ 风格统一为：小写字母 + 下划线；

示例：

```
POST {{host}}/demo/foobar
Content-Type: application/json

{
    "order_id": 27149
}
```

### 接口返回

+ 返回的 json 格式数据，风格统一为：小写字母 + 下划线;
+ 返回格式固定为：code、message、data 字段，code 为 0 表示正常，message 为 code 的描述信息，data 为接口数据；
+ 对于表示列表的字段，字段名建议加上 _list 后缀；
+ 对于表示对象的字段，字段名建议加上 _info 后缀。

示例：

```
{
    "code": 0,
    "message": "ok",
    "data": {
        "user_info": {
            "id": "5",
            "name": "旺仔"
        },
        "user_list": [
            {
                "id": "5",
                "name": "旺仔"
            },
            {
                "id": "7",
                "name": "旺旺"
            }
        ]
    }
}
```

## web app 目录结构及开发说明

```
.
├── ci.yml
├── Makefile
├── README.md
├── bin
├── log
├── data
├── doc
├── go.mod
├── go.sum
├── main.go
├── conf
├── routers
├── entities
├── controllers
├── models
├── library
└── script
```

代码的执行流程如下：

![请求执行流程]()

下面逐一对各个目录/文件进行说明。

### doc

todo::...

### main.go

todo::...

### conf

todo::...

### routers

todo::...

### entities

todo::...

### controllers

todo::...

### models/service

todo::...

### models/dao

todo::...

### library

todo::...

### script

todo::...

## web app lib 目录结构及开发说明

```
lib
├── app
│   ├── app.go
│   ├── app_error.go
│   ├── env.go
│   └── inits.go
├── config
│   ├── agollo.go
│   └── conf.go
├── middlewares
│   └── router_not_found.go
├── server
│   ├── Engine.go
│   ├── grpc_response.go
│   ├── http_response.go
│   ├── regexes.go
│   ├── regexes_test.go
│   ├── server_grpc.go
│   ├── server_http.go
│   └── webcontext.go
└── storage
    ├── mysql.go
    ├── redis.go
    └── sqlite.go
```

## 整体设计说明

业务层框架以及规范

规范化业务分层、逻辑拆分、接口约定、模块的目录/命名结构，以及部署同构这一系列，应用开发、部署过程中所面临的问题


### 名词解释

业务层:	在具体web框架之上进行的业务逻辑开发的部分。

业务层框架:	主要是把业务逻辑开发过程中一些共性的部分抽取出来建立的框架。这部分包含控制器组件，通用业务组件（参数处理，session处理，打日志等），配置组件，日志组件等等。

业务层规范:	侧重于具体业务/应用开发过程中需要处理的诸如，如何分层，层次之间接口如何约定，部署环境，以及业务/应用内部的目录，命名组织等。


### 设计目标


### 概览

```
app1 | app2 | ... | appn
------------------------
    gin | GRPC
```

#### 逻辑层次规范

```
controller
pageservice
dataservice
dao | orm
```

#### 编译产出规范

项目目录中的以下目录及文件:

+ ci.yml
+ Makefile
+ bin/appctl
+ bin/start.sh
+ bin/stop.sh
+ supervise

为打包、编译、服务启动及服务状态监控等脚本文件，默认提供，无需关心。 代码提交、编译通过后，下载 tar.gz 包，解压完成后，目录结构如下:

```
output
├── bin
│   ├── demo
│   ├── start.sh
│   ├── stop.sh
│   └── appctl
├── conf
├── log
├── data
└── supervise
```

