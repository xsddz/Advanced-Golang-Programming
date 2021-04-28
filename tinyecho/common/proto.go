package common

import (
	"io"
)

// protoI 通信协议接口
type protoI interface {
	WriteData(w io.Writer, payload string) (n int, err error)
	ReadMessage(r io.Reader) (msg string, err error)
}
