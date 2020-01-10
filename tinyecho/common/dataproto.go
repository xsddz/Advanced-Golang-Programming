package common

import (
	"errors"
	"fmt"
	"io"
	"strconv"
)

// dataProto data protocol format
// a data transmit on network shoule format to string based on this proto,
// and convert to []byte. a version 1 data on network should like:
//    010013hello, world!
//    010018你好，世界。
type dataProto struct {
	version    [2]byte // max 99
	messageLen [4]byte // max 9999
	message    messageProto
	// checksum   int64
}

// messageProto message protocol format
type messageProto struct {
	user  []byte
	sep   byte
	words []byte
}

const (
	// messageMaxLen message max length
	messageMaxLen = 9999
)

var (
	errMessageVersion = errors.New("wrong message version")
	errMessageToLarge = errors.New("to large message")
)

func makeV1Message(playload string) []byte {
	return []byte(fmt.Sprintf("%02d", 1) + fmt.Sprintf("%04d", len(playload)) + playload)
}

// WriteMessage write messge to io writer
func WriteMessage(w io.Writer, playload string) (n int, err error) {
	if len(playload) > messageMaxLen {
		err = errMessageToLarge
		return
	}

	n, err = w.Write(makeV1Message(playload))

	return
}

// ReadMessagePlayload read message from io reader
func ReadMessagePlayload(r io.Reader) (buf []byte, err error) {
	// Read message version
	version, err := readMessageVersion(r)
	if err != nil {
		return
	}

	switch version {
	case messageVersion1:
		// Read message
		buf, err = readV1MessagePlayload(r)
	default:
		err = errMessageVersion
	}

	return
}

func readMessageVersion(r io.Reader) (version int64, err error) {
	buf := make([]byte, 2)
	_, err = r.Read(buf)
	if err != nil {
		return 0, err
	}

	version, err = strconv.ParseInt(string(buf), 10, 64)
	return
}

func readV1MessagePlayload(r io.Reader) (buf []byte, err error) {
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
