package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	_ "tinyrpc/server/service"
)

func main() {
	listerer, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("Listen TCP error:", err)
	}

	for {
		conn, err := listerer.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}

		// go rpc.ServeConn(conn)
		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
