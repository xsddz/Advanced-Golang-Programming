package main

import (
	"Advanced-Golang-Programming/tinyecho/common"
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var (
	network string
	host    string
	port    string

	addresses string
)

func init() {
	flag.StringVar(&network, "proto", "tcp", "network protocol")
	flag.StringVar(&host, "host", "127.0.0.1", "host")
	flag.StringVar(&port, "port", "8897", "port")
	flag.Parse()

	addresses = fmt.Sprintf("%s:%s", host, port)
}

func main() {
	conn, err := net.Dial(network, addresses)
	if err != nil {
		log.Fatal("Dialing:", err)
	}
	defer conn.Close()
	fmt.Printf("Connected to <%v> by <%v>: \n", conn.RemoteAddr().String(), conn.LocalAddr().String())

	// Receive message
	go func() {
		playload, err := common.ReadMessagePlayload(conn)
		fmt.Printf("Receive message: %v\n", string(playload))
		if err != nil {
			log.Fatal("Borken connetion.")
		}
	}()

	// Send message
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		text = strings.TrimRight(text, "\n")

		n, err := common.WriteMessage(conn, text)
		fmt.Printf("Send message: %v,%v\n", n, err)
		if err != nil {
			log.Fatal("Borken connetion.")
		}
	}

}
