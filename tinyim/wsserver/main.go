package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"

	"Advanced-Golang-Programming/tinyim/wsserver/asset"
	"Advanced-Golang-Programming/tinyim/wsserver/auth"
	"Advanced-Golang-Programming/tinyim/wsserver/client"
	"Advanced-Golang-Programming/tinyim/wsserver/msg"
)

var (
	port string // param

	broadcast = make(chan msg.Message, 4096)
)

func init() {
	flag.StringVar(&port, "port", "8070", "port number")
	flag.Parse()

	// 解压静态文件
	dirs := []string{"webclient"}
	for _, dir := range dirs {
		if err := asset.RestoreAssets("./", dir); err != nil {
			break
		}
	}
}

func main() {
	// Start deliver for incoming chat messages
	go deliverMessages()

	mux := http.NewServeMux()
	// Configure websocket client page
	mux.Handle("/", http.FileServer(http.Dir("./webclient/")))
	// Configure websocket route
	mux.HandleFunc("/chat/", handleConnections)

	// Start the server on localhost port
	address := fmt.Sprintf(":%s", port)
	log.Println("http server started on:", address)
	err := http.ListenAndServe(address, mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// step1:
	// Authorization
	clientInfo, ok := auth.Login(r)
	if !ok {
		// handle authorization fail
		log.Println("authorization fail.")
		return
	}

	// step2:
	// Upgrade initial GET request to a websocket
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		// handle error
		log.Println(err)
		return
	}

	// step3:
	// Register new client
	// Make sure we close the connection when the function returns
	client.RegisterClientConn(clientInfo, conn)
	defer client.CleanClientConn(clientInfo)

	// step4:
	// Receive message loop
	for {
		data, err := wsutil.ReadClientText(conn)
		if err != nil {
			// handle error
			// log.Println("ReadClientText error:", err)
			break
		}
		fmt.Println("ReadClientText:", string(data), err)

		// step5:
		// Report the newly received message to the broadcast channel
		broadcast <- msg.Message{
			Channel: clientInfo.Channel,
			Data:    data,
		}
	}
}

func deliverMessages() {
	for {
		select {
		// Every 1 minutes send a keepalive message
		case <-time.After(60 * time.Second):
			client.WalkChannel(func(ch string) {
				broadcast <- msg.Message{
					Channel: ch,
					Data:    []byte("system: " + time.Now().Format("2006/01/02 15:04:05")),
				}
			})

		// Grab the next message from the broadcast channel
		// Send it out to every client that is currently connected
		case chmsg := <-broadcast:
			fmt.Println("deliverMessages:", string(chmsg.Data))
			client.WalkChannelClient(chmsg.Channel, func(userID string, conn net.Conn) {
				err := wsutil.WriteServerText(conn, chmsg.Data)
				if err != nil {
					fmt.Printf("\tWriteServerText: %s,%s [FAIL]\n", chmsg.Channel, userID)
				}
			})
		}
	}
}
