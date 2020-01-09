package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
)

// a message version 1 should like
//     <version num><message len><message><check sum>
//     version num: two byte
//     message len: four byte
//     message: <message len> byte
//     check sum: for feature
//     eg:
//         010018你好，世界。
const messageVersion1 = 1
const messageMaxLen = 4096

func readConn(c net.Conn) (buf []byte, e error) {
	// Read message version
	buf = make([]byte, 2)
	_, err := c.Read(buf)
	if err != nil {
		return []byte{}, err
	}
	version, _ := strconv.ParseInt(string(buf), 10, 32)

	switch version {
	case messageVersion1:
		// Read message length
		buf = make([]byte, 4)
		_, err := c.Read(buf)
		if err != nil {
			return []byte{}, err
		}
		messageLen, _ := strconv.ParseInt(string(buf), 10, 32)

		// Read message
		buf = make([]byte, messageLen)
		_, err = c.Read(buf)
	}

	return
}

func writeConn(c net.Conn, s string) (n int, e error) {
	if len(s) > messageMaxLen {
		e = errors.New("to large message")
		return
	}

	// Bulid version 1 message
	s = fmt.Sprintf("%02d", 1) + fmt.Sprintf("%04d", len(s)) + s

	// Write to client
	n, e = c.Write([]byte(s))

	return
}

func main() {
	conn, err := net.Dial("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("Dialing:", err)
	}
	defer conn.Close()
	fmt.Printf("Connected to %v: \n", conn.RemoteAddr().String())

	// Send message
	n, err := writeConn(conn, "你好")
	fmt.Printf("Send message: %v,%v\n", n, err)
	if err != nil {
		log.Fatal("Borken connetion.")
	}

	// Receive message
	msg, err := readConn(conn)
	fmt.Printf("Receive message: %v\n", string(msg))
	if err != nil {
		log.Fatal("Borken connetion.")
	}
}
