package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"yawebapp/library/infra/helper"
	"yawebapp/library/infra/tools/httpclientcoder"
)

const (
	VERSION = "1.0.0"
)

func init() {
	flag.Usage = usage
	flag.Parse()
}

func usage() {
	fmt.Fprintf(os.Stderr, "usage: http-client-coder [path/to/proto/file ...]\n")
	flag.PrintDefaults()
	fmt.Fprintf(os.Stderr, "\n")
}

func main() {
	if flag.NArg() == 0 {
		usage()
		return
	}

	failedTPL := "coder file: " + helper.ColorYellowBold + "%s" + helper.ColorGreenBold + " -> " + helper.ColorRedBold + "%s" + helper.ColorReset + "\n"
	successTPL := "coder file: " + helper.ColorYellowBold + "%s" + helper.ColorGreenBold + " -> " + helper.ColorYellowBold + "%s" + helper.ColorReset + "\n"
	for _, sourceFile := range flag.Args() {
		if !helper.IsExist(sourceFile) {
			fmt.Printf(failedTPL, sourceFile, "not exist")
			continue
		}

		toFile := filepath.Dir(sourceFile) + "/" + strings.ReplaceAll(filepath.Base(sourceFile), ".", "_") + ".go"
		if !helper.IsDir(filepath.Dir(toFile)) {
			helper.MakeDirP(filepath.Dir(toFile))
		}

		fmt.Printf(successTPL, sourceFile, toFile)
		processFile(toFile, sourceFile)
	}
}

func processFile(toFile string, sourceFile string) {
	fr, err := os.Open(sourceFile)
	if err != nil {
		panic(err)
	}
	p, err := httpclientcoder.NewPaser(fr)
	if err != nil {
		panic(err)
	}

	c, err := httpclientcoder.NewCoder()
	if err != nil {
		panic(err)
	}

	l, err := httpclientcoder.NewLinker()
	if err != nil {
		panic(err)
	}

	l.Content(c.Head(&httpclientcoder.HeadCmd{CoderVer: VERSION, ProtoFile: filepath.Base(sourceFile)}))
	for p.HasMoreCommands() {
		p.Advance()

		switch p.CommandType() {
		case httpclientcoder.CMD_NONE:
			// nothing need todo now
		case httpclientcoder.CMD_COMMENT_INLINE:
			// nothing need todo now
		case httpclientcoder.CMD_COMMENT:
			// nothing need todo now
		case httpclientcoder.CMD_COMMENT_CONTENT:
			// nothing need todo now
		case httpclientcoder.CMD_COMMENT_END:
			// nothing need todo now
		case httpclientcoder.CMD_CONF:
			// nothing need todo now
		case httpclientcoder.CMD_CLIENT:
			l.Content(c.Client(p.Cmd()))
		case httpclientcoder.CMD_CLIENT_ENDPOINT:
			// nothing need todo now
		case httpclientcoder.CMD_CLIENT_METHOD:
			l.Content(c.Method(p.Cmd()))
		case httpclientcoder.CMD_ENTITY:
			l.Content(c.Entity(p.Cmd()))
		case httpclientcoder.CMD_ENTITY_PROP:
			l.Content(c.EntityProp(p.Cmd()))
		case httpclientcoder.CMD_END_BLOCK:
			l.Content(c.EndBlock(p.Cmd()))
		default:
			panic(fmt.Errorf("error command:\n%v", p.CommandRaw()))
		}
	}

	os.Remove(toFile)
	fileWriter, err := os.OpenFile(toFile, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	err = l.Write(io.MultiWriter(fileWriter))
	if err != nil {
		panic(err)
	}
}
