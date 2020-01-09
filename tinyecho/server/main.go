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

func handleClientConn(c net.Conn) {
	defer c.Close()

	// Set client info
	clientAddr := c.RemoteAddr().String()

	for {
		// Receive message from client
		msg, err := readConn(c)
		fmt.Printf("Receive message from <%v>: %v\n", clientAddr, string(msg))
		if err != nil {
			fmt.Printf("Borken connetion from <%v>.\n", clientAddr)
			break
		}

		// Send message to client
		n, err := writeConn(c, clientAddr+": "+string(msg))
		fmt.Printf("Send message to <%v>: %v,%v\n", clientAddr, n, err)
	}
}

func main() {
	// Listen on TCP port 2000 on all available unicast and
	// anycast IP addresses of the local system.
	listerer, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Listen TCP error:", err)
	}
	defer listerer.Close()

	for {
		// Wait for a connection.
		conn, err := listerer.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}

		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		fmt.Printf("Recvive connection from %v.\n", conn.RemoteAddr())
		go handleClientConn(conn)
	}
}
