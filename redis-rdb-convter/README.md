
redis rdb database file convter.

```
# usage:
$ rdbconvter -d path/to/name.rdb

# eg:
$ rdbconvter -d ../tinyredislib/data/dump.rdb
convert file: /path/to/tinyredislib/data/dump.rdb
--------------------------------------------------------------------------------
rdb raw in hex                            means                                   
--------------------------------------------------------------------------------
52 45 44 49 53 30 30 30 37                => [rdbLoadMagicString] magic string and rdb version: REDIS0007
--------------------------------------------------------------------------------
FA                                        => [rdbLoadType] rdb optype: FA, 250
09                                        => [rdbLoadLen] 6 bit len: [00001001], 9
72 65 64 69 73 2D 76 65 72                => [rdbGenericLoadStringObject]: redis-ver
05                                        => [rdbLoadLen] 6 bit len: [00000101], 5
33 2E 32 2E 33                            => [rdbGenericLoadStringObject]: 3.2.3
                                          [opcodeHandlerAux] key: redis-ver, val: 3.2.3
FA                                        => [rdbLoadType] rdb optype: FA, 250
0A                                        => [rdbLoadLen] 6 bit len: [00001010], 10
72 65 64 69 73 2D 62 69 74 73             => [rdbGenericLoadStringObject]: redis-bits
C0                                        => [rdbLoadLen] 6 bit encoding type: [11000000], 0
40                                        => [rdbLoadIntegerObject] encode int8: [01000000], 64
                                          [opcodeHandlerAux] key: redis-bits, val: 64
FA                                        => [rdbLoadType] rdb optype: FA, 250
05                                        => [rdbLoadLen] 6 bit len: [00000101], 5
63 74 69 6D 65                            => [rdbGenericLoadStringObject]: ctime
C2                                        => [rdbLoadLen] 6 bit encoding type: [11000010], 2
7C 8D 7F 60                               => [rdbLoadIntegerObject] encode int32: [01111100 10001101 01111111 01100000], 1618972028
                                          [opcodeHandlerAux] key: ctime, val: 1618972028
FA                                        => [rdbLoadType] rdb optype: FA, 250
08                                        => [rdbLoadLen] 6 bit len: [00001000], 8
75 73 65 64 2D 6D 65 6D                   => [rdbGenericLoadStringObject]: used-mem
C2                                        => [rdbLoadLen] 6 bit encoding type: [11000010], 2
E8 96 0B 00                               => [rdbLoadIntegerObject] encode int32: [11101000 10010110 00001011 00000000], 759528
                                          [opcodeHandlerAux] key: used-mem, val: 759528
FE                                        => [rdbLoadType] rdb optype: FE, 254
00                                        => [rdbLoadLen] 6 bit len: [00000000], 0
                                          [opcodeHandlerSelectDB] select db: 0
FB                                        => [rdbLoadType] rdb optype: FB, 251
01                                        => [rdbLoadLen] 6 bit len: [00000001], 1
00                                        => [rdbLoadLen] 6 bit len: [00000000], 0
                                          [opcodeHandlerResizeDB] db_size:1, expire_size: 0
00                                        => [rdbLoadType] rdb optype: 0, 0
03                                        => [rdbLoadLen] 6 bit len: [00000011], 3
61 61 61                                  => [rdbGenericLoadStringObject]: aaa
C0                                        => [rdbLoadLen] 6 bit encoding type: [11000000], 0
6F                                        => [rdbLoadIntegerObject] encode int8: [01101111], 111
                                          [rdbTypeCommonHandler] key: aaa, valType: 0, val: 111
FE                                        => [rdbLoadType] rdb optype: FE, 254
07                                        => [rdbLoadLen] 6 bit len: [00000111], 7
                                          [opcodeHandlerSelectDB] select db: 7
FB                                        => [rdbLoadType] rdb optype: FB, 251
01                                        => [rdbLoadLen] 6 bit len: [00000001], 1
00                                        => [rdbLoadLen] 6 bit len: [00000000], 0
                                          [opcodeHandlerResizeDB] db_size:1, expire_size: 0
00                                        => [rdbLoadType] rdb optype: 0, 0
03                                        => [rdbLoadLen] 6 bit len: [00000011], 3
61 61 61                                  => [rdbGenericLoadStringObject]: aaa
C1                                        => [rdbLoadLen] 6 bit encoding type: [11000001], 1
DE 00                                     => [rdbLoadIntegerObject] encode int16: [11011110 00000000], 222
                                          [rdbTypeCommonHandler] key: aaa, valType: 0, val: 222
FF                                        => [rdbLoadType] rdb optype: FF, 255
                                          [opcodeHandlerEOF] end of data.
--------------------------------------------------------------------------------
59 F6 28 AB B6 AC BC 30                   => [rdbLoadCRC64Checksum] rdbver:7, CRC 64 checksum: [01011001 11110110 00101000 10101011 10110110 10101100 10111100 00110000]
--------------------------------------------------------------------------------
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