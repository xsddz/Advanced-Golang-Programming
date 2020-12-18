// assembler - 主程序，驱动整个编译过程

package tinyassembler

// CommandT - 指令类型
type CommandT int8

const (
	_ CommandT = iota //
	// ACommand A指令：0 v v v  v  v  v  v   v  v  v  v   v  v  v  v
	ACommand
	// CCommand C指令：1 1 1 a  c1 c2 c3 c4  c5 c6 d1 d2  d3 j1 j2 j3
	CCommand
	// LCommand 伪指令：(Xxx)，Xxx为符号
	LCommand
)

// PredefinedSymbols -
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
