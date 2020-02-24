package main

import (
	"Advanced-Golang-Programming/tinyim/wsserver/auth"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
)

var (
	clientsMap sync.Map
	broadcast  = make(chan []byte)
)

func main() {
	// Configure websocket route
	http.HandleFunc("/chat", handleConnections)

	// Start deliver for incoming chat messages
	go deliverMessages()

	// Start the server on localhost port 8080
	log.Println("http server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Authorization
	client, ok := auth.Login(r)
	if !ok {
		// handle authorization fail
		fmt.Println("authorization fial.")
		return
	}

	// Upgrade initial GET request to a websocket
	conn, _, _, err := ws.UpgradeHTTP(r, w)
	if err != nil {
		// handle error
		fmt.Println(err)
	}
	// Make sure we close the connection when the function returns
	defer conn.Close()

	// Register new client
	clientsMap.Store(client.UserID, conn)

	// handle client send message
	for {
		msg, err := wsutil.ReadClientText(conn)
		if err != nil {
			// handle error
			fmt.Println(err)
			// unregister client
			clientsMap.Delete(client.UserID)
			return
		}
		fmt.Println("read message:\t", string(msg), err)

		// Send the newly received message to the broadcast channel
		broadcast <- msg
	}
}

func deliverMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast

		// Send it out to every client that is currently connected
		clientsMap.Range(func(key, value interface{}) bool {
			fmt.Println("\tsend message to", key)
			c := value.(net.Conn)
			err := wsutil.WriteServerText(c, msg)
			if err != nil {
				// handle error
				return true
			}

			return true
		})
	}
}
