// code - 编码，提供所有汇编命令对应的二进制码

package tinyassembler

// Code -
type Code interface {
	Dest(string) string
	Comp(string) string
	Jump(string) string
}
