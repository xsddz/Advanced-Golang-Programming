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

	idname string
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

	// Parse welcome message and get client idname
	msg, err := common.ReadMessage(conn)
	welcomeMsg := string(msg)
	idname = welcomeMsg[strings.Index(welcomeMsg, "<")+1 : strings.Index(welcomeMsg, ">")]
	fmt.Printf("%v\n\n", welcomeMsg)

	// Receive message loop
	go func() {
		for {
			msg, err := common.ReadMessage(conn)
			fmt.Printf("\033[1E%v\n", string(msg)) // display in next line
			if err != nil {
				log.Fatal("Read connection error:", err)
			}
		}
	}()

	// Send message loop
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("%v... ", common.YellowString(idname))

		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("Read input error:", err)
		}
		text = strings.TrimRight(text, "\n")
		if len(text) == 0 {
			continue
		}

		_, err = common.WriteData(conn, text)
		// fmt.Printf("Send message: %v,%v,%v\n", text, n, err)
		if err != nil {
			log.Fatal("Write connetion error:", err)
		}
	}
}
