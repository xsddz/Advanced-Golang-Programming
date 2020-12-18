// parser - 语法分析器，对输入文件进行语法分析

package tinyassembler

// Paser -
type Paser interface {
	New(file string)
	HasMoreCommands() bool
	Advance()
	CommandType() CommandT
	Symbol() string
	Dest() string
	Comp() string
	Jump() string
}
