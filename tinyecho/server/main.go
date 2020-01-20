package main

import (
	"Advanced-Golang-Programming/tinyecho/common"
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	protocol string
	address  string

	clientInfoMap sync.Map
)

type tcpClientInfo struct {
	name      string
	addresses string
	conn      net.Conn
	loginTime time.Time
}

type udpClientInfo struct {
	name      string
	addr      *net.UDPAddr
	loginTime time.Time
}

func init() {
	flag.StringVar(&protocol, "proto", "tcp", "network protocol")
	flag.StringVar(&address, "addr", "127.0.0.1:8897", "address")
	flag.Parse()
}

func main() {
	switch protocol {
	case "tcp":
		tcpServer()
	case "udp":
		udpServer()
	default:
		fmt.Println("Server protocol error: ", protocol)
	}
}

func udpServer() {
	localUDPAddr, err := net.ResolveUDPAddr(protocol, address)
	if err != nil {
		log.Fatal("Address error:", protocol, err)
	}
	udpConn, err := net.ListenUDP(protocol, localUDPAddr)
	if err != nil {
		log.Fatal("Listen error:", err)
	}
	fmt.Printf("Server running on %s %s ...\n\n", protocol, address)

	data := make([]byte, 4096)
	for {
		// 获取udp连接
		n, clientUDPAddr, err := udpConn.ReadFromUDP(data)
		if err != nil {
			fmt.Printf("error during read: %s", err)
		}
		fmt.Printf("<%s> %s\n", clientUDPAddr, data[:n])

		// 处理udp连接
		go func(address string, r io.Reader) {
			// 注册失败（即，已经注册）时，发送登录信息
			// TODO: 后续需要对注册信息进行keepalive验证
			uci, ok := registerUDPClient(address)
			if ok {
				welcomeDetail := fmt.Sprintf("Welcome, your name is <%v>.", uci.name)
				n, err := udpConn.WriteToUDP(common.MakeV1ProtoData(welcomeDetail), uci.addr)
				fmt.Printf("Send message to <%v>: %v,%v,%v\n", uci.addr.String(), welcomeDetail, n, err)
			}

			// 读取当前连接客户端信息
			msg, err := common.ReadMessage(r)
			recvDetail := fmt.Sprintf("%v... %v", uci.name, string(msg))
			if err != nil {
				fmt.Printf("\nRead connection error: %v,%v\n", err, recvDetail)
			}
			fmt.Printf("\nReceive message: %v\n", recvDetail)

			// 发送给其他客户端
			clientInfoMap.Range(func(key, value interface{}) bool {
				suci := value.(*udpClientInfo)
				if suci.addr != uci.addr {
					n, err := udpConn.WriteToUDP(common.MakeV1ProtoData(recvDetail), suci.addr)
					fmt.Printf("Send message to <%v>: %v,%v,%v\n", suci.name, recvDetail, n, err)
				}

				return true
			})
		}(clientUDPAddr.String(), bytes.NewReader(data[:n]))
	}
}

func registerUDPClient(address string) (ci *udpClientInfo, ok bool) {
	// TODO: generate uniq random name
	hash := md5.Sum([]byte(address))
	name := base64.StdEncoding.EncodeToString(hash[:])
	name = strings.TrimRight(name, "=")

	if v, ok := clientInfoMap.Load(name); ok {
		ci = v.(*udpClientInfo)
		return ci, false
	}

	udpAddr, _ := net.ResolveUDPAddr(protocol, address)
	ci = &udpClientInfo{
		name:      name,
		addr:      udpAddr,
		loginTime: time.Now(),
	}
	clientInfoMap.Store(name, ci)
	return ci, true
}

func tcpServer() {
	tcpAddr, err := net.ResolveTCPAddr(protocol, address)
	if err != nil {
		log.Fatal("Address error:", protocol, err)
	}
	// Listen on <network> port <port> on IP addresses <host> of the local system.
	listerer, err := net.ListenTCP(protocol, tcpAddr)
	if err != nil {
		log.Fatal("Listen error:", err)
	}
	defer listerer.Close()
	fmt.Printf("Server running on %s %s ...\n\n", protocol, address)

	for {
		// Wait for a connection.
		tcpConn, err := listerer.AcceptTCP()
		if err != nil {
			log.Fatal("Accept error:", err)
		}

		// Handle the connection in a new goroutine.
		// The loop then returns to accepting, so that
		// multiple connections may be served concurrently.
		fmt.Printf("Accept connection from <%v>.\n", tcpConn.RemoteAddr().String())
		go func(c net.Conn) {
			ci := registerTCPClient(c)
			defer cleanTCPClient(ci)

			// Send welcome message
			welcomeDetail := fmt.Sprintf("Welcome, your name is <%v>.", ci.name)
			n, err := common.WriteData(ci.conn, welcomeDetail)
			fmt.Printf("Send message to <%v>: %v,%v,%v\n", ci.addresses, welcomeDetail, n, err)

			for {
				// Receive message from client
				msg, err := common.ReadMessage(ci.conn)
				recvDetail := fmt.Sprintf("%v... %v", ci.name, string(msg))
				if err != nil {
					fmt.Printf("\nRead connection error: %v,%v\n", err, recvDetail)
					break
				}
				fmt.Printf("\nReceive message: %v\n", recvDetail)

				// Send message to all clients except this
				clientInfoMap.Range(func(key, value interface{}) bool {
					sci := value.(*tcpClientInfo)
					if sci.conn != ci.conn {
						n, err := common.WriteData(sci.conn, recvDetail)
						fmt.Printf("Send message to <%v>: %v,%v,%v\n", sci.name, recvDetail, n, err)
					}

					return true
				})
			}
		}(tcpConn)
	}
}

func registerTCPClient(c net.Conn) *tcpClientInfo {
	addresses := c.RemoteAddr().String()

	// TODO: generate uniq random name
	hash := md5.Sum([]byte(addresses))
	name := base64.StdEncoding.EncodeToString(hash[:])
	name = strings.TrimRight(name, "=")

	ci := &tcpClientInfo{
		name:      name,
		addresses: addresses,
		conn:      c,
	}

	clientInfoMap.Store(name, ci)

	return ci
}

func cleanTCPClient(ci *tcpClientInfo) {
	// clean from map
	clientInfoMap.Delete(ci.name)

	// close connection
	ci.conn.Close()
}
