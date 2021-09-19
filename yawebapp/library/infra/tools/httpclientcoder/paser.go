package httpclientcoder

import (
	"bufio"
	"io"
	"regexp"
	"strings"
	"yawebapp/library/infra/helper"
)

var (
	regTrSpace = regexp.MustCompile(`(\s)+`)
	regTagItem = regexp.MustCompile(`([a-zA-Z0-9]+):"([^"]*)"`)
)

type Paser struct {
	f *bufio.Scanner

	nl string // next line
	cl string // current line

	cmd          Cmder
	lastBlockCmd Cmder

	clientCmd *ClientCmd
	entityCmd *EntityCmd
}

func NewPaser(r io.Reader) (*Paser, error) {
	scanner := bufio.NewScanner(r)
	return &Paser{f: scanner}, nil
}

func (p *Paser) nextLine() (string, error) {
	for p.f.Scan() {
		return p.f.Text(), nil
	}

	if err := p.f.Err(); err != nil {
		return "", err
	}

	return "", io.EOF
}

// HasMoreCommands 输入中还有更多命令吗？
func (p *Paser) HasMoreCommands() bool {
	line, err := p.nextLine()
	p.nl = line
	return err == nil
}

// Advance 从输入中读取下一条命令，将其当作“当前命令”，
// 仅当 hasMoreCommands() 为真时，才能调用该方法，
// 最初始的时候，没有“当前命令”
func (p *Paser) Advance() {
	p.cl = p.nl
	p.nl = ""
}

// CommandRaw 返回当前原始命令
func (p *Paser) CommandRaw() string {
	return p.cl
}

// CommandType 返回当前命令的类型
func (p *Paser) CommandType() CommandType {
	cmd := TrimSpace(p.cl)
	// empty
	if cmd == "" {
		return CMD_NONE
	}

	// cut cmd by space
	cmdArr := strings.Split(cmd, " ")

	// comment part
	if cmdArr[0] == "//" || strings.HasPrefix(cmdArr[0], "#") {
		return CMD_COMMENT_INLINE
	}
	if cmdArr[0] == "/**" {
		return CMD_COMMENT
	}
	if cmdArr[0] == "*" {
		return CMD_COMMENT_CONTENT
	}
	if cmdArr[0] == "*/" {
		return CMD_COMMENT_END
	}

	// conf
	if strings.ToLower(cmdArr[0]) == "conf" && len(cmdArr) == 2 {
		return CMD_CONF
	}

	// block end
	if cmdArr[0] == "}" {
		p.cmd = &EndBlockCmd{p.lastBlockCmd}
		return CMD_END_BLOCK
	}

	// client part
	if strings.ToLower(cmdArr[0]) == "client" && len(cmdArr) == 3 {
		p.clientCmd = &ClientCmd{cmdArr[1]}
		p.cmd = p.clientCmd
		p.lastBlockCmd = p.cmd
		return CMD_CLIENT
	}
	if strings.ToLower(cmdArr[0]) == "endpoint" && len(cmdArr) == 2 {
		return CMD_CLIENT_ENDPOINT
	}
	if helper.InArray(strings.ToLower(cmdArr[0]), []string{"get", "post", "put", "head", "options"}) {
		m := make(map[string]string)
		if len(cmdArr) == 5 {
			tags := strings.Split(strings.Trim(cmdArr[4], "`"), " ")
			for _, item := range tags {
				a := regTagItem.FindStringSubmatch(item)
				m[strings.ToLower(a[1])] = strings.ToLower(a[2])
			}
		}

		p.cmd = &ClientMethodCmd{
			BeginBlockCmd: p.clientCmd,
			Symbol:        strings.ToUpper(cmdArr[0]),
			Path:          strings.Trim(cmdArr[1], "\""),
			ReqSymbol:     p.clientCmd.Symbol + cmdArr[2][3:],
			RespSymbol:    p.clientCmd.Symbol + cmdArr[3][4:],
			Tag:           m,
		}
		return CMD_CLIENT_METHOD
	}

	// entity part
	if strings.ToLower(cmdArr[0]) == "entity" && len(cmdArr) == 3 {
		if strings.HasPrefix(cmdArr[1], "req") {
			p.entityCmd = &EntityCmd{Symbol: p.clientCmd.Symbol + cmdArr[1][3:], Type: "req"}
		} else {
			p.entityCmd = &EntityCmd{Symbol: p.clientCmd.Symbol + cmdArr[1][4:], Type: "resp"}
		}
		p.cmd = p.entityCmd
		p.lastBlockCmd = p.cmd
		return CMD_ENTITY
	}
	if len(cmdArr) == 2 {
		t := cmdArr[1]
		if strings.HasPrefix(t, "resp") {
			t = "Response" + p.clientCmd.Symbol + t[4:]
		}
		p.cmd = &EntityPropCmd{BeginBlockCmd: p.entityCmd, Symbol: cmdArr[0], Type: t}
		return CMD_ENTITY_PROP
	}

	return CMD_ERROR
}

func (p *Paser) Cmd() Cmder {
	return p.cmd
}
