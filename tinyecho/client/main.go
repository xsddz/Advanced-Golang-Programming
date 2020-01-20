package main

import (
	"Advanced-Golang-Programming/tinyecho/common"
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	protocol      string
	serverAddress string
	localAddress  string

	idname string
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
		udpClient()
	case "tcp":
		tcpClient()
	default:
		fmt.Println("Server protocol error: ", protocol)
	}
}

func udpClient() {
	localUDPAddr, err := net.ResolveUDPAddr(protocol, localAddress)
	if err != nil {
		log.Fatal("Address error:", protocol, err)
	}
	serverUDPAddr, err := net.ResolveUDPAddr(protocol, serverAddress)
	if err != nil {
		log.Fatal("Address error:", protocol, err)
	}
	udpConn, err := net.ListenUDP(protocol, localUDPAddr)
	if err != nil {
		log.Fatal("Listen error:", err)
	}

	// 登录获取用户信息
	udpConn.WriteToUDP(common.MakeV1ProtoData("Login"), serverUDPAddr)
	data := make([]byte, 4096)
	n, _, err := udpConn.ReadFromUDP(data)
	if err != nil {
		fmt.Printf("error during read: %s", err)
	}
	msg, err := common.ReadMessage(bytes.NewReader(data[:n]))
	welcomeMsg := string(msg)
	idname = welcomeMsg[strings.Index(welcomeMsg, "<")+1 : strings.Index(welcomeMsg, ">")]
	fmt.Printf("%v\n\n", welcomeMsg)

	var wg sync.WaitGroup
	wg.Add(2)
	// 进入获取信息循环
	go func() {
		defer wg.Done()
		data := make([]byte, 4096)
		for {
			n, _, err := udpConn.ReadFromUDP(data)
			if err != nil {
				fmt.Printf("error during read: %s", err)
			}
			msg, err := common.ReadMessage(bytes.NewReader(data[:n]))
			fmt.Println("\033[1E\n", string(msg), err)
		}
	}()
	// 进入发送信息循环
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Printf("<%v>... ", common.YellowString(idname))

			text, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal("Read input error:", err)
			}
			text = strings.TrimRight(text, "\n")
			if len(text) == 0 {
				continue
			}

			_, err = udpConn.WriteToUDP(common.MakeV1ProtoData(text), serverUDPAddr)
			if err != nil {
				log.Fatal("Write connetion error:", err)
			}
		}
	}()
	wg.Wait()
}

func tcpClient() {
	conn, err := net.Dial(protocol, serverAddress)
	if err != nil {
		log.Fatal("Dialing:", err)
	}
	defer conn.Close()
	fmt.Printf("Connected to <%v> by <%v>: \n", conn.RemoteAddr().String(), conn.LocalAddr().String())

	common.WriteData(conn, "")
	// Parse welcome message and get client idname
	msg, err := common.ReadMessage(conn)
	welcomeMsg := string(msg)
	idname = welcomeMsg[strings.Index(welcomeMsg, "<")+1 : strings.Index(welcomeMsg, ">")]
	fmt.Printf("%v\n\n", welcomeMsg)

	var wg sync.WaitGroup
	wg.Add(2)
	// Receive message loop
	go func() {
		defer wg.Done()
		for {
			msg, err := common.ReadMessage(conn)
			fmt.Printf("\033[1E%v\n", string(msg)) // display in next line
			if err != nil {
				log.Fatal("Read connection error:", err)
			}
		}
	}()
	// Send message loop
	go func() {
		defer wg.Done()
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
	}()
	wg.Wait()
}
