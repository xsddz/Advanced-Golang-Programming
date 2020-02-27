package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

var (
	wsAddr    string // param
	roomNo    string // param
	clientNum int    // param

	clientMap     sync.Map            // saved connected client
	userInputChan = make(chan string) // transfer user input text
	recvMsgChan   = make(chan int)    // control only one connected client receive server msg

	isQuit = false
)

func init() {
	flag.StringVar(&wsAddr, "wsAddr", "ws://127.0.0.1:8070/chat?authorization=2W8EAFNSUPY8EY5", "websocket address")
	flag.StringVar(&roomNo, "roomNo", "258EAFA5", "room number")
	flag.IntVar(&clientNum, "clientNum", 20, "client number")
	flag.Parse()

	// add room number to wsAddr
	wsAddr = fmt.Sprintf("%s&room=%s", wsAddr, roomNo)
}

func main() {
	initClient()
	defer cleanClient()

	// send & recv message loop
	go eventLoop()
	recvMsgChan <- 1

	// input text loop
	fmt.Println("\n",
		"| for receive message from server,\n",
		"|     we will use a random choosed connected client;\n",
		"| and you can input text below, we will send it to server\n",
		"|     by a random choosed connected client everytime,\n",
		"| and ^D for quit.\n",
		"| \n",
		"| enjoy yourself.\n",
		"")
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			if err != io.EOF {
				log.Println("[ERROR] read input:", err)
			}
			isQuit = true
			break
		}

		text = strings.TrimRight(text, "\n")
		if len(text) == 0 {
			continue
		}

		userInputChan <- text
	}
}

func makeClientNo(number int) string {
	return fmt.Sprintf("client-%07d", number)
}

func initClient() {
	for i := 0; i < clientNum; i++ {
		clientNo := makeClientNo(i + 1)
		fmt.Printf("\r=======begin initClient: [%d] %s", i+1, clientNo)

		conn, _, _, err := ws.Dial(context.Background(), wsAddr)
		if err != nil {
			log.Println("\n[ERROR] initClient error:", err)
			i-- // retry
			continue
		}

		clientMap.Store(clientNo, conn)
	}
	fmt.Println(", end initClient.=======")
}

func cleanClient() {
	fmt.Printf("\n\n")
	count := 0
	clientMap.Range(func(k, v interface{}) bool {
		clientNo := k.(string)
		conn := v.(net.Conn)
		conn.Close()

		count++
		fmt.Printf("\r=======begin cleanclient: [%d] %s", count, clientNo)

		return true
	})
	fmt.Println(", end cleanclient.=======")
}

func randomClient(callback func(conn net.Conn)) {
	var conn net.Conn
	for {
		// Choose a random connection
		rand.Seed(time.Now().UnixNano())
		v, ok := clientMap.Load(makeClientNo(rand.Intn(clientNum) + 1))
		if !ok {
			continue
		}

		conn = v.(net.Conn)
		break
	}
	callback(conn)
}

func walkClient(callback func(conn net.Conn)) {
	clientMap.Range(func(k, v interface{}) bool {
		conn := v.(net.Conn)
		callback(conn)
		return true
	})
}

func eventLoop() {
	for {
		select {
		// Grab the next message from the userInputChan channel
		case msg := <-userInputChan:
			randomClient(func(conn net.Conn) {
				clientAddr := conn.LocalAddr().String()
				wsutil.WriteClientText(conn, []byte(fmt.Sprintf("%s 说: \"%s\"", clientAddr, msg)))
			})
		// Grab the ball from recvMsgChan channel
		case <-recvMsgChan:
			randomClient(func(conn net.Conn) {
				go func(conn net.Conn) {
					clientAddr := conn.LocalAddr().String()
					for {
						msg, err := wsutil.ReadServerText(conn)
						if err != nil {
							if !isQuit {
								log.Println("[ERROR] wsutil.ReadServerText:", err)
							}
							break
						}

						fmt.Println(clientAddr, "recv msg:\t", string(msg), err)
					}
				}(conn)
			})
		}
	}
}
