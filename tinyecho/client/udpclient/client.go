package udpclient

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"tinyecho/common"
)

const (
	// Protocol protocol
	Protocol = "udp"
)

// Run run a udp client on local address, and connect to server address
func Run(localAddress string, serverAddress string) {
	localAddr, err := net.ResolveUDPAddr(Protocol, localAddress)
	if err != nil {
		log.Fatal("[ERROR] address:", Protocol, err)
	}
	serverAddr, err := net.ResolveUDPAddr(Protocol, serverAddress)
	if err != nil {
		log.Fatal("[ERROR] address:", Protocol, err)
	}
	udpConn, err := net.ListenUDP(Protocol, localAddr)
	if err != nil {
		log.Fatal("[ERROR] listen:", err)
	}

	// send login message
	udpConn.WriteToUDP(common.V1ProtoMsgMaker("Login"), serverAddr)
	// wait for udp connection and read welcome message
	data := make([]byte, 4096)
	n, _, err := udpConn.ReadFromUDP(data)
	if err != nil {
		log.Fatal("[ERROR] receive udp connection and read message:", err)
	}
	// parse welcome message and get name
	msg, err := common.ReadMessage(bytes.NewReader(data[:n]))
	welcomeMsg := string(msg)
	idname := welcomeMsg[strings.Index(welcomeMsg, "<")+1 : strings.Index(welcomeMsg, ">")]
	fmt.Printf("%v\n\n", welcomeMsg)

	var wg sync.WaitGroup
	wg.Add(2)
	// read message loop
	go func() {
		defer wg.Done()
		data := make([]byte, 4096)
		for {
			n, _, err := udpConn.ReadFromUDP(data)
			if err != nil {
				log.Println("[ERROR] receive udp connection and read message:", err)
				continue
			}
			msg, err := common.ReadMessage(bytes.NewReader(data[:n]))
			fmt.Printf("\033[1E%v\n", string(msg))
		}
	}()
	// send message loop
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Printf("<%v>... ", common.YellowString(idname))

			text, err := reader.ReadString('\n')
			if err != nil {
				log.Fatal("[ERROR] read input error:", err)
			}
			text = strings.TrimRight(text, "\n")
			if len(text) == 0 {
				continue
			}

			_, err = udpConn.WriteToUDP(common.V1ProtoMsgMaker(text), serverAddr)
			if err != nil {
				log.Fatal("[ERROR] write connetion error:", err)
			}
		}
	}()
	wg.Wait()
}
