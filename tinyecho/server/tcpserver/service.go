package tcpserver

import (
	"Advanced-Golang-Programming/tinyecho/common"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const (
	// Protocol Protocol
	Protocol = "tcp"
)

var (
	// clientInfoListMap login client map
	clientInfoListMap sync.Map
)

type clientInfo struct {
	name      string
	address   string
	conn      net.Conn
	loginTime time.Time
}

func registerClient(conn net.Conn) *clientInfo {
	address := conn.RemoteAddr().String()
	name := common.GenerateUsername(address)

	// TODO: read login message
	// send welcome message
	welcomeDetail := fmt.Sprintf("Welcome, your name is <%v>.", name)
	n, err := common.WriteData(conn, welcomeDetail)
	fmt.Printf("Send message to <%v>: %v,%v,%v\n", address, welcomeDetail, n, err)

	c := &clientInfo{
		name:      name,
		address:   address,
		conn:      conn,
		loginTime: time.Now(),
	}
	clientInfoListMap.Store(name, c)
	return c
}

func cleanClient(c *clientInfo) {
	// clean from map
	clientInfoListMap.Delete(c.name)

	// close connection
	c.conn.Close()
}

// Run run a tcp server on tcp address
func Run(address string) {
	tcpAddr, err := net.ResolveTCPAddr(Protocol, address)
	if err != nil {
		log.Fatal("[ERROR] address:", Protocol, err)
	}
	listerer, err := net.ListenTCP(Protocol, tcpAddr)
	if err != nil {
		log.Fatal("[ERROR] listen:", err)
	}
	defer listerer.Close()
	fmt.Printf("Server running on %s %s ...\n\n", Protocol, address)

	for {
		// wait for a tcp client connection.
		tcpConn, err := listerer.AcceptTCP()
		if err != nil {
			log.Fatal("[ERROR] receive tcp connection:", err)
		}
		fmt.Printf("Receive tcp client connection: %s\n", tcpConn.RemoteAddr().String())

		// handle the client connection in a new goroutine.
		// the loop then returns to AcceptTCP, so that multiple connections
		// may be served concurrently.
		go func(conn net.Conn) {
			// register client
			c := registerClient(conn)
			defer cleanClient(c)

			// receive and send message loop
			for {
				// receive message from client
				msg, err := common.ReadMessage(c.conn)
				if err != nil {
					log.Printf("[ERROR] read client data: %v,%v,%v\n", err, c.address, string(msg))
					break
				}
				recvDetail := fmt.Sprintf("%v %v", c.name, string(msg))
				fmt.Printf("Receive message: %v\n", recvDetail)

				// send client message to other login client
				clientInfoListMap.Range(func(key, value interface{}) bool {
					sc := value.(*clientInfo)
					if sc.conn != c.conn {
						n, err := common.WriteData(sc.conn, recvDetail)
						fmt.Printf("Send message to <%v>: %v,%v,%v\n", sc.name, recvDetail, n, err)
					}

					return true
				})
			}
		}(tcpConn)
	}
}