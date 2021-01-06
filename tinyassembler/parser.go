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
	file *os.File

	currentCommand    string
	currentCommandLen int

	commandBuf string
}

// NewTParser -
func NewTParser(file *os.File) *TParser {
	return &TParser{
		file: file,
	}
}

// readline - 从打开的文件中读取一行文本
func readline(f *os.File) (line []byte, err error) {
	var data [1]byte
	for {
		_, err = f.Read(data[:])
		if err != nil { // incluce io.EOF
			break
		}
		if data[0] == '\n' {
			// remove char \r before char \n
			lineLen := len(line)
			if lineLen > 0 && line[lineLen-1] == '\r' {
				line = line[:lineLen-1]
			}
			break
		}
		line = append(line, data[0])
	}
	return line, err
}

// HasMoreCommands -
func (p *TParser) HasMoreCommands() bool {
	// 直到遇到一行指令或文件结束为止
	for {
		// 从文件中读取一行进行
		line, err := readline(p.file)
		if err != nil && err != io.EOF { // 读取出错
			return false
		}

		// 解析读取行
		charCount := 0
		mayComment := false
		for _, c := range line {
			// 去除注释：
			// - 跳过注释行
			// - 忽略行尾注释
			if c == '/' {
				if mayComment {
					charCount = charCount - 1
					break
				}
				mayComment = true
			} else {
				mayComment = false
			}

			// 去除空格：
			// - 忽略行内空白符
			if c == ' ' {
				continue
			}

			line[charCount] = c
			charCount = charCount + 1
		}

		// 跳过空白行
		if charCount == 0 {
			if err == io.EOF {
				return false
			}
			continue
		}

		// 缓存解析出来的汇编指令
		p.commandBuf = string(line[:charCount])
		break
	}
	return true
}

// Advance -
func (p *TParser) Advance() {
	p.currentCommand = p.commandBuf
	p.currentCommandLen = len(p.currentCommand)
}

// CommandType -
func (p *TParser) CommandType() CommandT {
	if p.currentCommand[0] == '@' {
		return ACommand
	} else if p.currentCommand[0] == '(' && p.currentCommand[p.currentCommandLen-1] == ')' {
		return LCommand
	} else {
		return CCommand
	}
}

// Symbol -
func (p *TParser) Symbol() string {
	// ACommand
	if p.currentCommand[0] == '@' {
		return p.currentCommand[1:]
	}
	// LCommand
	return p.currentCommand[1 : p.currentCommandLen-1]
}

// Dest -
func (p *TParser) Dest() string {
	destEndIndex := -1
	for i, c := range p.currentCommand {
		if c == '=' {
			destEndIndex = i - 1
		}
	}
	if destEndIndex == -1 {
		return "null"
	}
	return p.currentCommand[:destEndIndex+1]
}

// Comp -
func (p *TParser) Comp() string {
	compBeginIndex := 0
	compEndIndex := p.currentCommandLen - 1
	for i, c := range p.currentCommand {
		if c == '=' {
			compBeginIndex = i + 1
		}
		if c == ';' {
			compEndIndex = i - 1
		}
	}
	return p.currentCommand[compBeginIndex : compEndIndex+1]
}

// Jump -
func (p *TParser) Jump() string {
	jumpBeginIndex := -1
	for i, c := range p.currentCommand {
		if c == ';' {
			jumpBeginIndex = i + 1
		}
	}
	if jumpBeginIndex == -1 {
		return "null"
	}
	return p.currentCommand[jumpBeginIndex:]
}
