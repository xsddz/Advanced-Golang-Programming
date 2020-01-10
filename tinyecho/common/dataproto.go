package common

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

// data protocol format:
//     <version><playload><checkSum>
//     - version: 2 byte, max version number is 99
//     - playload: n byte
//     - checkSum: for feature
//
// and a version 1 protocol playload shoule
//     <messageLen><message>
//     - messageLen: 4 byte, max message length is 9999
//     - message: messageLen byte
//
//     eg:
//         010013hello, world!
//         010018你好，世界。

const (
	// protocolVersion1 protocol version 1 number
	protocolVersion1 = 1
	// v1ProtoMessageMaxLen protocol version 1 message max length
	v1ProtoMessageMaxLen = 9999
)

var (
	errMessageVersion = errors.New("message version not support")
	errMessageToLarge = errors.New("to large message")
)

func makeV1ProtoData(playload string) []byte {
	return []byte(fmt.Sprintf("%02d", protocolVersion1) + fmt.Sprintf("%04d", len(playload)) + playload)
}

// WriteData write data to io writer
func WriteData(w io.Writer, playload string) (n int, err error) {
	if len(playload) > v1ProtoMessageMaxLen {
		err = errMessageToLarge
		return
	}

	n, err = w.Write(makeV1ProtoData(playload))

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
