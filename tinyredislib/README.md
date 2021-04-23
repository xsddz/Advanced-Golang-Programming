
记录阅读 redis 源码过程中的思考，分析 redis 的设计思想（design Philosophy）。

目录说明：

```
tinyredislib
├── README.md
├── ae             // io多路复用模块
├── data           // 测试数据
├── networking     // 网络io模块
└── rdb            // rdb模块
```

以上各个模块，在分析、模拟实现之后，可以用于开发一些 redis 的小工具，比如

+ 基于 rdb 模块开发的 rdbconvter

