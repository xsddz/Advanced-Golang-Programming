package common

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

// version 1 data protocol format:
//     <version><messageLen><message><checkSum>
//     - version: 2 byte, max version number is 99
//     - messageLen: 4 byte, max message length value is 4096-2-4=4090
//     - message: messageLen byte
//     - checkSum: for feature
//
//     eg:
//         010013hello, world!
//         010018你好，世界。

const (
	// protocolVersion1 protocol version 1 number
	protocolVersion1 = 1
	// v1ProtoDataMaxLen protocol version 1 data max length
	v1ProtoDataMaxLen = 4096
	// v1ProtoMessageMaxLen protocol version 1 message max length
	v1ProtoMessageMaxLen = 4090
)

var (
	errMessageVersion = errors.New("message version not support")
	errMessageToLarge = errors.New("to large message")
)

// MakeV1ProtoData MakeV1ProtoData
func MakeV1ProtoData(playload string) []byte {
	return []byte(fmt.Sprintf("%02d", protocolVersion1) + fmt.Sprintf("%04d", len(playload)) + playload)
}

// WriteData write data to io writer
func WriteData(w io.Writer, playload string) (n int, err error) {
	if len(playload) > v1ProtoMessageMaxLen {
		err = errMessageToLarge
		return
	}

	n, err = w.Write(MakeV1ProtoData(playload))

	return
}

// ReadMessage read message from io reader
func ReadMessage(r io.Reader) (buf []byte, err error) {
	// Read message version
	version, err := readProtoVersion(r)
	if err != nil {
		return
	}

	switch version {
	case protocolVersion1:
		// Read message
		buf, err = readV1Message(r)
	default:
		err = errMessageVersion
	}

	return
}

func readProtoVersion(r io.Reader) (version int64, err error) {
	buf := make([]byte, 2)
	_, err = r.Read(buf)
	if err != nil {
		return 0, err
	}

	version, err = strconv.ParseInt(string(buf), 10, 64)
	return
}

func readV1Message(r io.Reader) (buf []byte, err error) {
	buf = make([]byte, 4)
	_, err = r.Read(buf)
	if err != nil {
		return []byte{}, err
	}

	messageLen, err := strconv.ParseInt(string(buf), 10, 64)

	buf = make([]byte, messageLen)
	n, err := r.Read(buf)

	return buf[:n], err
}
