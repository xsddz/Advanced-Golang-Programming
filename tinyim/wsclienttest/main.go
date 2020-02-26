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
	wsaddr    string // param
	roomNo    string // param
	clientNum int    // param

	clientMap     sync.Map            // saved connected client
	userInputChan = make(chan string) // transfer user input text
	recvMsgChan   = make(chan int)    // control only one connected client receive server msg

	isQuit = false
)

func init() {
	flag.StringVar(&wsaddr, "wsaddr", "ws://127.0.0.1:8070/chat?authorization=2W8EAFNSUPY8EY5", "websocket address")
	flag.StringVar(&roomNo, "roomNo", "258EAFA5", "room number")
	flag.IntVar(&clientNum, "clientNum", 20, "client number")
	flag.Parse()

	// add room number to wsaddr
	wsaddr = fmt.Sprintf("%s&room=%s", wsaddr, roomNo)
}

func main() {
	initClient()
	defer cleanClient()

	go eventLoop()
	recvMsgChan <- 1

	fmt.Println("\n",
		"| we will random choose a connection to show message from server,\n",
		"| and you can input some text, we will send it to server with a random choose connection everytime,\n",
		"| and ^D for quit.\n",
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

func initClient() {
	fmt.Println("=======begin initClient:")
	for i := 0; i < clientNum; i++ {
		clientNo := makeClientNo(i + 1)
		fmt.Printf("\r[%d] %s", i+1, clientNo)

		conn, _, _, err := ws.Dial(context.Background(), wsaddr)
		if err != nil {
			log.Println("\n[ERROR] initClient error:", err)
			i-- // retry
			continue
		}

		clientMap.Store(clientNo, conn)
	}
	fmt.Println("\n=======end initClient.")
}

func cleanClient() {
	fmt.Println("\n=======begin cleanclient:")
	count := 0
	clientMap.Range(func(k, v interface{}) bool {
		clientNo := k.(string)
		conn := v.(net.Conn)
		conn.Close()

		count++
		fmt.Printf("\r[%d] %s", count, clientNo)

		return true
	})
	fmt.Println("\n=======end cleanclient.")
}

func makeClientNo(number int) string {
	return fmt.Sprintf("client-%07d", number)
}

func eventLoop() {
	for {
		// Choose a random connection
		rand.Seed(time.Now().UnixNano())
		v, ok := clientMap.Load(makeClientNo(rand.Intn(clientNum) + 1))
		if !ok {
			continue
		}
		conn := v.(net.Conn)
		clientAddr := conn.LocalAddr().String()

		select {
		// Grab the next message from the userInputChan channel
		case msg := <-userInputChan:
			wsutil.WriteClientText(conn, []byte(fmt.Sprintf("%s 说: \"%s\"", clientAddr, msg)))
		// Grab the recvMsgChan from recvMsgChan channel
		case <-recvMsgChan:
			go func(c net.Conn) {
				for {
					msg, err := wsutil.ReadServerText(c)
					if err != nil {
						if !isQuit {
							log.Println("[ERROR] wsutil.ReadServerText:", err)
						}
						break
					}

					fmt.Println(clientAddr, " read message:\t", string(msg), err)
				}
			}(conn)
		}
	}
}
