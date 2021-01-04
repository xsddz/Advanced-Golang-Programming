// SymbolTable - 符号表，在符号标签和数字地址之间建立关联

package tinyassembler

// SymbolTable -
type SymbolTable interface {
	// 创建空的符号表
	// New()

	// 将(symbol, address)配对加入符号表
	AddEntry(symbol string, address int)

	// 符号表是否包含指定的symbol
	Contains(symbol string) bool

	// 返回与symbol关联的地址
	GetAddress(symbol string) int
}

// TSymbolTable -
type TSymbolTable struct {
	savedTable map[string]int
}

// NewTSymbolTable -
func NewTSymbolTable() *TSymbolTable {
	return &TSymbolTable{
		savedTable: make(map[string]int),
	}
}

// AddEntry -
func (t *TSymbolTable) AddEntry(symbol string, address int) {
	t.savedTable[symbol] = address
}

// Contains -
func (t *TSymbolTable) Contains(symbol string) bool {
	if _, exist := t.savedTable[symbol]; exist {
		return true
	}
	return false
}

// GetAddress -
func (t *TSymbolTable) GetAddress(symbol string) int {
	return t.savedTable[symbol]
}
