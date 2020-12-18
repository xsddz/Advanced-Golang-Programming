// parser - 语法分析器，对输入文件进行语法分析

package tinyassembler

import "os"

// Paser -
type Paser interface {
	New(file string)
	Init(f os.File)
	HasMoreCommands() bool
	Advance()
	CommandType() CommandT
	Symbol() string
	Comp() string
	Dest() string
	Jump() string
}
