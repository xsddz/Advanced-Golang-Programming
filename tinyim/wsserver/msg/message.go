package msg

import "Advanced-Golang-Programming/tinyim/wsserver/auth"

// Message Message struct
type Message struct {
	Client *auth.ClientInfo
	Data   []byte
}
