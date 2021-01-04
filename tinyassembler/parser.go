// parser - 语法分析器，对输入文件进行语法分析

package tinyassembler

import (
	"io"
	"os"
)

// Parser -
type Parser interface {
	// 打开输入文件或输出流，为语法解析作准备
	// New(file string)

	// 输入中还有更多命令吗
	HasMoreCommands() bool

	// 从输入中读取下一条命令，将其当作“当前命令”，
	// 仅当HasMoreCommands()为真时，才能调用该方法，
	// 最初始的时候，没有“当前命令”
	Advance()

	// 返回当前命令的类型：
	// - ACommand 当@xxx中的xxx是符号或者十进制数字时
	// - CCommand 用于dest=comp;jump
	// - LCommand 伪指令，当(xxx)中的xxx是符号时
	CommandType() CommandT

	// 返回形如@xxx或(xxx)的当前命令的符号或十进制值，
	// 仅当CommandType()是ACommand或LCommand时才能调用
	Symbol() string

	// 返回当前C指令的dest助记符（共有8种形式），
	// 仅当CommandType()是CCommand时才能调用
	Dest() string
	// 返回当前C指令的comp助记符（共有28种形式），
	// 仅当CommandType()是CCommand时才能调用
	Comp() string
	// 返回当前C指令的jump助记符（共有8种形式），
	// 仅当CommandType()是CCommand时才能调用
	Jump() string
}

// TParser -
type TParser struct {
	file           *os.File
	currentCommand string
	commandBuf     string
}

// NewTParser -
func NewTParser(file *os.File) *TParser {
	return &TParser{
		file: file,
	}
}

// HasMoreCommands -
func (p *TParser) HasMoreCommands() bool {
	for {
		line, err := readline(p.file)
		if err != nil && err != io.EOF {
			return false
		}

		// 注释行
		if len(line) >= 2 && string(line[:2]) == "//" {
			if err == io.EOF {
				return false
			}
			continue
		}

		// 空白行
		if string(line) == "" {
			if err == io.EOF {
				return false
			}

			continue
		}

		p.commandBuf = string(line)
		break
	}

	return true
}

// readline - 从打开的文件中读取一行并返回
func readline(f *os.File) (line []byte, err error) {
	var data [1]byte
	for {
		_, err = f.Read(data[:])
		if err != nil { // incluce io.EOF
			break
		}
		if data[0] == '\n' {
			break
		}
		line = append(line, data[:]...)
	}
	return line, err
}

// Advance -
func (p *TParser) Advance() {
	p.currentCommand = p.commandBuf
}

// CommandType -
func (p *TParser) CommandType() CommandT {
	return ACommand
}

// Symbol -
func (p *TParser) Symbol() string {
	return ""
}

// Dest -
func (p *TParser) Dest() string {
	return ""
}

// Comp -
func (p *TParser) Comp() string {
	return ""
}

// Jump -
func (p *TParser) Jump() string {
	return ""
}
