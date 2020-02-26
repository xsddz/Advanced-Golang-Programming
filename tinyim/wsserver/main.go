package main

import (
	"Advanced-Golang-Programming/tinyim/wsserver/asset"
	"Advanced-Golang-Programming/tinyim/wsserver/auth"
	"Advanced-Golang-Programming/tinyim/wsserver/msg"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

var (
	port string // param

	clientsChannelMap      = make(map[string]map[string]net.Conn)
	clientsChannelMapMutex sync.RWMutex

	broadcast = make(chan msg.Message)
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

func registerClient(client *auth.ClientInfo, conn net.Conn) {
	clientsChannelMapMutex.Lock()
	if _, exist := clientsChannelMap[client.Channel]; !exist {
		clientsChannelMap[client.Channel] = make(map[string]net.Conn)
	}
	clientsChannelMap[client.Channel][client.UserID] = conn
	clientsChannelMapMutex.Unlock()
}

func unregisterClient(client *auth.ClientInfo) {
	clientsChannelMapMutex.Lock()
	conn := clientsChannelMap[client.Channel][client.UserID]
	delete(clientsChannelMap[client.Channel], client.UserID)
	conn.Close()
	fmt.Println("close client:", client.UserID)
	clientsChannelMapMutex.Unlock()
}

func main() {
	// Start deliver for incoming chat messages
	go deliverMessages()

	// Configure websocket client page
	http.Handle("/", http.FileServer(http.Dir("./webclient/")))
	// Configure websocket route
	http.HandleFunc("/chat", handleConnections)
	// Start the server on localhost port
	address := fmt.Sprintf(":%s", port)
	log.Println("http server started on:", address)
	err := http.ListenAndServe(address, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// step1:
	// Authorization
	client, ok := auth.Login(r)
	if !ok {
		// handle authorization fail
		fmt.Println("authorization fial.")
		return
	}

	// step2:
	// Upgrade initial GET request to a websocket
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		// handle error
		fmt.Println(err)
	}

	// step3:
	// Register new client
	// Make sure we close the connection when the function returns
	registerClient(client, conn)
	defer unregisterClient(client)

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
			Client: client,
			Data:   data,
		}
	}
}

func deliverMessages() {
	for {
		// Grab the next message from the broadcast channel
		// Send it out to every client that is currently connected
		chmsg := <-broadcast

		count := 0

		clients := clientsChannelMap[chmsg.Client.Channel]
		for u, c := range clients {
			count++

			err := wsutil.WriteServerText(c, chmsg.Data)
			if err != nil {
				fmt.Println("\tWriteServerText to:", u, "\t[FAIL]")
			}
		}

		fmt.Println("\tmsg:", string(chmsg.Data), ", send to channel:", chmsg.Client.Channel, ", send to client number:", count)
	}
}
