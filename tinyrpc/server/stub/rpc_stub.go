package stub

import (
	"fmt"
	"net"
)

func RunRPC() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("[RunRPC] listen err:", ln, err)
	}

	fmt.Println("listen :8080 ...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("[RunRPC] accept err:", conn, err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
}
