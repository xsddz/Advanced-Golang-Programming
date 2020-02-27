package client

import (
	"fmt"
	"net"
	"sync"
)

// UserInfo 客户端信息
type UserInfo struct {
	Channel string
	UserID  string
}

var (
	clientMap      sync.Map
	clientMapMutex sync.RWMutex
)

// RegisterClientConn 注册
func RegisterClientConn(u UserInfo, conn net.Conn) {
	clientMapMutex.Lock()
	v, exist := clientMap.LoadOrStore(u.Channel, map[string]net.Conn{u.UserID: conn})
	if exist {
		cm := v.(map[string]net.Conn)
		cm[u.UserID] = conn
		clientMap.Store(u.Channel, cm)
	}
	clientMapMutex.Unlock()
}

// CleanClientConn 释放
func CleanClientConn(u UserInfo) {
	clientMapMutex.Lock()
	v, ok := clientMap.Load(u.Channel)
	if ok {
		cm := v.(map[string]net.Conn)
		delete(cm, u.UserID)
		if len(cm) > 1 {
			clientMap.Store(u.Channel, cm)
		} else {
			clientMap.Delete(u.Channel)
		}
	}
	clientMapMutex.Unlock()
}

// WalkChannel WalkChannel
func WalkChannel(callback func(string)) {
	clientMap.Range(func(k, v interface{}) bool {
		ch := k.(string)
		callback(ch)
		return true
	})
}

// WalkChannelClient WalkChannelClient
func WalkChannelClient(ch string, callback func(string, net.Conn)) {
	v, ok := clientMap.Load(ch)
	if !ok {
		return
	}

	cm := v.(map[string]net.Conn)
	for userID, conn := range cm {
		callback(userID, conn)
	}

	fmt.Println("\tWalkChannelClient: channel:", ch, ", client number:", len(cm))
}
