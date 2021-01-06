// assembler - 主程序，驱动整个编译过程

package tinyassembler

import (
	"fmt"
	"os"
	"strconv"
)

// CommandT - 指令类型
type CommandT int8

const (
	_ CommandT = iota
	// ACommand A指令：0 v v v  v  v  v  v   v  v  v  v   v  v  v  v
	ACommand
	// CCommand C指令：1 1 1 a  c1 c2 c3 c4  c5 c6 d1 d2  d3 j1 j2 j3
	CCommand
	// LCommand 伪指令：(XXX)，XXX为符号
	LCommand
)

// MnemonicDestTable - CCommand中dest助记符
var MnemonicDestTable = map[string]string{
	"null": "000",
	"M":    "001",
	"D":    "010",
	"MD":   "011",
	"A":    "100",
	"AM":   "101",
	"AD":   "110",
	"AMD":  "111",
}

// MnemonicCompTable - CCommand中comp助记符
var MnemonicCompTable = map[string]string{
	"0":  "0101010",
	"1":  "0111111",
	"-1": "0111010",
	"D":  "0001100",
	"A":  "0110000", "M": "1110000",
	"!D": "0001101",
	"!A": "0110001", "!M": "1110001",
	"-D": "0001111",
	"-A": "0110011", "-M": "1110011",
	"D+1": "0011111",
	"A+1": "0110111", "M+1": "1110111",
	"D-1": "0001110",
	"A-1": "0110010", "M-1": "1110010",
	"D+A": "0000010", "D+M": "1000010",
	"D-A": "0010011", "D-M": "1010011",
	"A-D": "0000111", "M-D": "1000111",
	"D&A": "0000000", "D&M": "1000000",
	"D|A": "0010101", "D|M": "1010101",
}

// MnemonicJumpTable - CCommand中Jump助记符
var MnemonicJumpTable = map[string]string{
	"null": "000",
	"JGT":  "001",
	"JEQ":  "010",
	"JGE":  "011",
	"JLT":  "100",
	"JNE":  "101",
	"JLE":  "110",
	"JMP":  "111",
}

// PredefinedSymbols - 预定义符号
var PredefinedSymbols = map[string]int{
	"SP":     0,
	"LCL":    1,
	"ARG":    2,
	"THIS":   3,
	"THAT":   4,
	"R0":     0,
	"R1":     1,
	"R2":     2,
	"R3":     3,
	"R4":     4,
	"R5":     5,
	"R6":     6,
	"R7":     7,
	"R8":     8,
	"R9":     9,
	"R10":    10,
	"R11":    11,
	"R12":    12,
	"R13":    13,
	"R14":    14,
	"R15":    15,
	"SCREEN": 16384,
	"KBD":    24576,
}

// NextSymbolAddress - 符号地址从16开始分配
var NextSymbolAddress = 16

// Run -
func Run(source string, target string) {
	// init symbol table
	st := NewTSymbolTable()
	for symbol, address := range PredefinedSymbols {
		st.AddEntry(symbol, address)
	}

	// init code
	code := NewTCode()

	// init source/target file
	sf, err := os.Open(source)
	if err != nil {
		fmt.Println("open source file err:", err)
		return
	}
	defer sf.Close()
	tf, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("create target file err:", err)
		return
	}
	defer tf.Close()

	// first parser
	commandAddressCounter := -1
	parser := NewTParser(sf)
	for {
		hasMore := parser.HasMoreCommands()
		if !hasMore {
			break
		}

		parser.Advance()

		ct := parser.CommandType()
		if ct == ACommand || ct == CCommand { // A指令和C指令会被加载至指令内存
			commandAddressCounter++
		}
		if ct == LCommand { // 伪指令符号地址值为下一条指令地址，不会被加载至指令内存
			st.AddEntry(parser.Symbol(), commandAddressCounter+1)
		}
	}

	// second parser
	_, err = sf.Seek(0, 0)
	if err != nil {
		fmt.Println("second parser err:", err)
		return
	}
	parser = NewTParser(sf)
	for {
		hasMore := parser.HasMoreCommands()
		if !hasMore {
			break
		}

		parser.Advance()

		ct := parser.CommandType()
		if ct == LCommand {
			continue
		}
		if ct == ACommand {
			symbol := parser.Symbol()
			if symbol[0] >= '0' && symbol[0] <= '9' {
				value, _ := strconv.Atoi(symbol)
				tf.WriteString(fmt.Sprintf("%016b\n", value)) // 将A指令翻译为二进制指令
				continue
			}

			if !st.Contains(symbol) { // 分配新符号的RAM地址
				st.AddEntry(symbol, NextSymbolAddress)
				NextSymbolAddress++
			}
			tf.WriteString(fmt.Sprintf("%016b\n", st.GetAddress(symbol))) // 将A指令翻译为二进制指令
			continue
		}
		if ct == CCommand {
			tf.WriteString("111" + code.Comp(parser.Comp()) + code.Dest(parser.Dest()) + code.Jump(parser.Jump()) + "\n") // 将C指令翻译为二进制指令
			continue
		}
	}
}
