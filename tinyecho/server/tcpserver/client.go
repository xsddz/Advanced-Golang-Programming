package tcpserver

import (
	"Advanced-Golang-Programming/tinyecho/common"
	"fmt"
	"net"
	"sync"
	"time"
)

var (
	// clientInfoMap 登录客户端集合
	clientInfoMap sync.Map
)

type clientInfo struct {
	name       string    // 客户端别名
	address    string    // 客户端地址信息
	conn       net.Conn  // 客户端链接
	loginTime  time.Time // 客户端登录时间
	expireTime time.Time // 预留字段，到期时间，客户端会定期续约
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
	clientInfoMap.Store(address, c)
	return c
}

func cleanClient(c *clientInfo) {
	// close connection
	c.conn.Close()
	// clean from map
	clientInfoMap.Delete(c.name)
}
