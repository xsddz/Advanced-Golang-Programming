package main

import (
	"Advanced-Golang-Programming/tinyecho/common"
	"flag"
	"fmt"
	"log"
	"net"
	"sync"
)

var (
	network string
	host    string
	port    string

	addresses string

	clientConnMap sync.Map
)

func init() {
	flag.StringVar(&network, "proto", "tcp", "network protocol")
	flag.StringVar(&host, "host", "127.0.0.1", "host")
	flag.StringVar(&port, "port", "8897", "port")
	flag.Parse()

	addresses = fmt.Sprintf("%s:%s", host, port)
}

func main() {
	// Listen on <network> port <port> on IP addresses <host> of the local system.
	listerer, err := net.Listen(network, addresses)
	if err != nil {
		log.Fatal("Listen error:", err)
	}
	defer listerer.Close()
	fmt.Printf("Server running on %s %s ...\n", network, addresses)

	for {
		// Wait for a connection.
		conn, err := listerer.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}

		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		fmt.Printf("Accept connection from <%v>.\n", conn.RemoteAddr().String())
		go handleClientConn(conn)
	}
}

func handleClientConn(c net.Conn) {
	saveClientConn(c)
	defer cleanClientConn(c)

	// Set client info
	clientAddr := c.RemoteAddr().String()

	for {
		// Receive message from client
		playload, err := common.ReadMessagePlayload(c)
		if err != nil {
			fmt.Printf("Borken connetion from <%v>: %v, %v\n", clientAddr, string(playload), err)
			break
		}
		recvDetail := fmt.Sprintf("<%v># %v", clientAddr, string(playload))
		fmt.Printf("Receive message: %v\n", recvDetail)

		// Send message to all clients except this
		clientConnMap.Range(func(key, value interface{}) bool {
			sc := value.(net.Conn)
			if sc != c {
				n, err := common.WriteMessage(sc, recvDetail)
				fmt.Printf("Send message to <%v>: %v,%v,%v\n", sc.RemoteAddr().String(), recvDetail, n, err)
			}

			return true
		})
	}
}

func saveClientConn(c net.Conn) {
	clientConnMap.Store(c, c)
}

func cleanClientConn(c net.Conn) {
	clientConnMap.Delete(c)
	c.Close()
}
