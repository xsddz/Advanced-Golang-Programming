package httpclientcoder

import (
	"io"

	"golang.org/x/tools/imports"
)

type Linker struct {
	data string
}

func NewLinker() (*Linker, error) {
	return &Linker{}, nil
}

func (l *Linker) Content(cnt string) {
	l.data += cnt
}

func (l *Linker) Format() error {
	formatContent, err := imports.Process("", []byte(l.data), nil)
	if err != nil {
		return err
	}

	l.data = string(formatContent)
	return nil
}

func (l *Linker) Write(w io.Writer) error {
	err := l.Format()
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(l.data))
	return err
}
