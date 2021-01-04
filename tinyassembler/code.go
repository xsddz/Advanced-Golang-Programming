// code - 编码，将所有汇编命令助记符翻译为对应的二进制码

package tinyassembler

// Code -
type Code interface {
	// 返回dest助记符对应的二进制码，3bits
	Dest(string) string
	// 返回comp助记符对应的二进制码，7bits
	Comp(string) string
	// 返回jump助记符对应的二进制码，3bits
	Jump(string) string
}

// TCode -
type TCode struct{}

// NewTCode -
func NewTCode() *TCode {
	return &TCode{}
}

func (c *TCode) unpackCode(mnemonic string, mnemonicTable map[string]string) string {
	if code, ok := mnemonicTable[mnemonic]; ok {
		return code
	}
	return ""
}

// Dest -
func (c *TCode) Dest(dest string) string {
	return c.unpackCode(dest, MnemonicDest)
}

// Comp -
func (c *TCode) Comp(comp string) string {
	return c.unpackCode(comp, MnemonicComp)
}

// Jump -
func (c *TCode) Jump(jump string) string {
	return c.unpackCode(jump, MnemonicJump)
}
