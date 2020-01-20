package main

import (
	"Advanced-Golang-Programming/tinyecho/client/tcpclient"
	"Advanced-Golang-Programming/tinyecho/client/udpclient"
	"flag"
	"fmt"
)

var (
	protocol      string
	serverAddress string
	localAddress  string
)

func init() {
	flag.StringVar(&protocol, "proto", "tcp", "network protocol")
	flag.StringVar(&serverAddress, "saddr", "127.0.0.1:8897", "server address")
	flag.StringVar(&localAddress, "laddr", "127.0.0.1:55310", "local address")
	flag.Parse()
}

func main() {
	switch protocol {
	case "udp":
		udpclient.Run(localAddress, serverAddress)
	case "tcp":
		tcpclient.Run(serverAddress)
	default:
		fmt.Println("Server protocol error: ", protocol)
	}
}
