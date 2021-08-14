package tcpclient

import (
	"bufio"
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
	Protocol = "tcp"
)

// Run run a tcp client, and connect to server address
func Run(serverAddress string) {
	tcpConn, err := net.Dial(Protocol, serverAddress)
	if err != nil {
		log.Fatal("[ERROR] Dialing:", Protocol, err)
	}
	defer tcpConn.Close()
	fmt.Printf("Connected to <%v> by <%v>:\n", tcpConn.RemoteAddr().String(), tcpConn.LocalAddr().String())

	// send login message
	common.WriteData(tcpConn, "Login")
	// read welcome message
	msg, err := common.ReadMessage(tcpConn)
	// parse welcome message and get client idname
	welcomeMsg := string(msg)
	idname := welcomeMsg[strings.Index(welcomeMsg, "<")+1 : strings.Index(welcomeMsg, ">")]
	fmt.Printf("%v\n\n", welcomeMsg)

	var wg sync.WaitGroup
	wg.Add(2)
	// receive message loop
	go func() {
		defer wg.Done()
		for {
			msg, err := common.ReadMessage(tcpConn)
			if err != nil {
				log.Println("[ERROR] read connection error:", err, string(msg))
				break
			}
			fmt.Printf("\033[1E%v\n", string(msg)) // display in next line
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

			_, err = common.WriteData(tcpConn, text)
			if err != nil {
				log.Fatal("[ERROR] write connetion error:", err)
			}
		}
	}()
	wg.Wait()
}
