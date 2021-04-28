package common

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

/*******************************************************************************
 * 自定义协议实现
 ******************************************************************************/

const (
	// protoVer protocol version
	// version 1 data protocol format:
	//     <version><payloadlen><payload><checkSum>
	//     - version: 2 byte, interpreted as ascii characters, so max version number is 99
	//     - payloadlen: 4 byte, interpreted as ascii characters, max payload length value is 9999
	//     - payload: payloadlen byte payload
	//     - checkSum: for feature
	//
	//     eg:
	//         []byte("010013hello, world!")
	//         []byte("010018你好，世界。")
	protoVer  = 1
	maxMsgLen = 9999
)

var (
	errMsgVer      = errors.New("message version not support")
	errMsgTooLarge = errors.New("message to large")
)

// WriteData write data to io writer
func WriteData(w io.Writer, payload string) (n int, err error) {
	if len(payload) > maxMsgLen {
		err = errMsgTooLarge
		return
	}

	return w.Write(V1ProtoMsgMaker(payload))
}

// ReadMessage read message from io reader
func ReadMessage(r io.Reader) (msg string, err error) {
	// Read message version
	ver, err := readProtoVer(r)
	if err != nil {
		return
	}

	return readProtoPayload(r, ver)
}

// V1ProtoMsgMaker -
func V1ProtoMsgMaker(payload string) []byte {
	return []byte(fmt.Sprintf("%02d", protoVer) + fmt.Sprintf("%04d", len(payload)) + payload)
}

func readProtoVer(r io.Reader) (ver uint64, err error) {
	buf := make([]byte, 2)
	_, err = r.Read(buf)
	if err != nil {
		return 0, err
	}

	return strconv.ParseUint(string(buf), 10, 64)
}

func readProtoPayload(r io.Reader, protover uint64) (msg string, err error) {
	msgLen, err := readProtoPayloadLen(r)
	if err != nil {
		return "", err
	}

	buf := make([]byte, msgLen)
	_, err = r.Read(buf)
	return string(buf), err
}

func readProtoPayloadLen(r io.Reader) (len uint64, err error) {
	buf := make([]byte, 4)
	_, err = r.Read(buf)
	if err != nil {
		return 0, err
	}

	return strconv.ParseUint(string(buf), 10, 64)
}
