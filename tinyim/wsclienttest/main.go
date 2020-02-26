package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

var (
	wsaddr    string
	clientNum int
)

func init() {
	flag.StringVar(&wsaddr, "wsaddr", "ws://127.0.0.1:8070/chat?room=258EAFA5&authorization=2W8EAFNSUPY8EY5", "websocket address")
	flag.IntVar(&clientNum, "clientNum", 20, "client number")
	flag.Parse()
}

var broadcast = make(chan string)
var ball = make(chan int)

func main() {
	fmt.Println("====begin initClient: ", clientNum)
	for i := 0; i < clientNum; i++ {
		go initClient()
	}
	fmt.Println("====end initClient: ", clientNum)

	ball <- 1
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal("[ERROR] read input:", err)
		}

		text = strings.TrimRight(text, "\n")
		if len(text) == 0 {
			continue
		}

		broadcast <- text
	}
}

func initClient() {
	conn, _, _, err := ws.Dial(context.Background(), wsaddr)
	if err != nil {
		fmt.Println("=======initClient error:", err)
		return
	}
	defer conn.Close()

	clientAddr := conn.LocalAddr().String()
	// wsutil.WriteClientText(conn, []byte(fmt.Sprintf("login: %s", clientAddr)))

	for {
		select {
		// Grab the next message from the broadcast channel
		case msg := <-broadcast:
			wsutil.WriteClientText(conn, []byte(fmt.Sprintf("%s 说: \"%s\"", clientAddr, msg)))
		// Grab the ball from ball channel
		case <-ball:
			go func() {
				for {
					msg, err := wsutil.ReadServerText(conn)
					if err != nil {
						fmt.Println("=======wsutil.ReadServerText error:", err)
						return
					}

					fmt.Println(clientAddr, " read message:\t", string(msg), err)
				}
			}()
		}
	}
}
