package udpserver

import (
	"Advanced-Golang-Programming/tinyecho/common"
	"bytes"
	"fmt"
	"log"
	"net"
	"sync"
	"time"
)

const (
	// Protocol Protocol
	Protocol = "udp"
)

var (
	// clientInfoListMap login client map
	clientInfoListMap sync.Map
)

type clientInfo struct {
	name      string
	addr      *net.UDPAddr
	loginTime time.Time
}

func registerClient(address string) (c *clientInfo, ok bool) {
	name := common.GenerateUsername(address)

	if v, ok := clientInfoListMap.Load(name); ok {
		c = v.(*clientInfo)
		return c, false
	}

	udpAddr, _ := net.ResolveUDPAddr(Protocol, address)
	c = &clientInfo{
		name:      name,
		addr:      udpAddr,
		loginTime: time.Now(),
	}
	clientInfoListMap.Store(name, c)
	return c, true
}

// Run run a udp server on udp address
func Run(address string) {
	serverAddr, err := net.ResolveUDPAddr(Protocol, address)
	if err != nil {
		log.Fatal("[ERROR] address:", Protocol, err)
	}
	udpConn, err := net.ListenUDP(Protocol, serverAddr)
	if err != nil {
		log.Fatal("[ERROR] listen:", err)
	}
	fmt.Printf("Server running on %s %s ...\n\n", Protocol, address)

	data := make([]byte, 4096)
	for {
		// wait for a udp client connection，and read the data from this udp client
		n, clientAddr, err := udpConn.ReadFromUDP(data)
		if err != nil {
			log.Fatal("[ERROR] receive udp connection:", err)
		}
		fmt.Printf("Receive udp client connection: %s, %s\n", clientAddr.String(), data[:n])

		// handle the client connection and data in a new goroutine.
		// the loop then returns to ReadFromUDP, so that multiple connections
		// may be served concurrently.
		go func(address string, data []byte) {
			// parse client data
			msg, err := common.ReadMessage(bytes.NewReader(data))
			if err != nil {
				log.Printf("[ERROR] parse client data: %v,%v,%v\n", err, address, string(data))
				return
			}

			// when register client success, send login message
			// when register client failed, client already login, no need send login message
			// TODO: we need a keepalive logic to monitor login client in clientInfoListMap
			c, ok := registerClient(address)
			if ok {
				welcomeDetail := fmt.Sprintf("Welcome, your name is <%v>.", c.name)
				n, err := udpConn.WriteToUDP(common.MakeV1ProtoData(welcomeDetail), c.addr)
				fmt.Printf("Send message to <%v>: %v,%v,%v\n", c.addr.String(), welcomeDetail, n, err)
			}

			// display client message detail
			recvDetail := fmt.Sprintf("%v %v", c.name, string(msg))
			fmt.Printf("Receive message: %v\n", recvDetail)

			// send client message to other login client
			clientInfoListMap.Range(func(key, value interface{}) bool {
				sc := value.(*clientInfo)
				if sc.addr != c.addr {
					n, err := udpConn.WriteToUDP(common.MakeV1ProtoData(recvDetail), sc.addr)
					fmt.Printf("Send message to <%v>: %v,%v,%v\n", sc.name, recvDetail, n, err)
				}

				return true
			})
		}(clientAddr.String(), data[:n])
	}
}
