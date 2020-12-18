// symbol table - 符号表，处理符号

package tinyassembler

// SymbolTable -
type SymbolTable interface {
	New()                 // 创建空的符号表
	AddEntry(string, int) //
	Contains(string) bool
	GetAddress(string) int
}
