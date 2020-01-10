package main

import (
	"Advanced-Golang-Programming/tinyecho/common"
	"crypto/md5"
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

var (
	network string
	host    string
	port    string

	addresses string

	clientConnMap sync.Map
)

type clientInfo struct {
	name      string
	addresses string
	conn      net.Conn
}

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
	fmt.Printf("Server running on %s %s ...\n\n", network, addresses)

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
	ci := registerClient(c)
	defer cleanClient(ci)

	// Send welcome message
	welcomeDetail := fmt.Sprintf("Welcome, your name is <%v>.", ci.name)
	n, err := common.WriteMessage(ci.conn, welcomeDetail)
	fmt.Printf("Send message to <%v>: %v,%v,%v\n", ci.addresses, welcomeDetail, n, err)

	for {
		// Receive message from client
		playload, err := common.ReadMessagePlayload(ci.conn)
		recvDetail := fmt.Sprintf("%v... %v", common.GreenString(ci.name), string(playload))
		if err != nil {
			fmt.Printf("\nRead connection error: %v,%v\n", err, recvDetail)
			break
		}
		fmt.Printf("\nReceive message: %v\n", recvDetail)

		// Send message to all clients except this
		clientConnMap.Range(func(key, value interface{}) bool {
			sci := value.(*clientInfo)
			if sci.conn != ci.conn {
				n, err := common.WriteMessage(sci.conn, recvDetail)
				fmt.Printf("Send message to <%v>: %v,%v,%v\n", sci.name, recvDetail, n, err)
			}

			return true
		})
	}
}

func registerClient(c net.Conn) *clientInfo {
	addresses := c.RemoteAddr().String()

	// TODO: generate uniq random name
	hash := md5.Sum([]byte(addresses))
	name := base64.StdEncoding.EncodeToString(hash[:])
	name = strings.TrimRight(name, "=")

	ci := &clientInfo{
		name:      name,
		addresses: addresses,
		conn:      c,
	}

	clientConnMap.Store(name, ci)

	return ci
}

func cleanClient(ci *clientInfo) {
	// clean from map
	clientConnMap.Delete(ci.name)

	// close connection
	ci.conn.Close()
}
