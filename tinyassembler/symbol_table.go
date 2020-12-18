// symbol table - 符号表，处理符号

package tinyassembler

// SymbolTable -
type SymbolTable interface {
	New()
	AddEntry(string, int)
	Contains(string) bool
	GetAddress(string) int
}
