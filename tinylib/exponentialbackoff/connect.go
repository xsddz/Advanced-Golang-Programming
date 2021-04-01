package exponentialbackoff

import (
	"fmt"
	"net"
	"time"
)

// MaxSleep -
const MaxSleep = 20

// ConnectRetry -
func ConnectRetry(network string, address string) (conn net.Conn, err error) {
	for numsec := 1; numsec <= MaxSleep; numsec <<= 1 {
		fmt.Printf("[ConnectRetry]connect to [%s]%s ", network, address)
		if conn, err = net.DialTimeout(network, address, time.Duration(1)); err == nil {
			// connect success, return
			fmt.Printf("success!\n")
			break
		}
		fmt.Printf("faild.\n")

		// delay before trying again
		if numsec <= MaxSleep/2 {
			fmt.Println("[ConnectRetry]sleep second:", numsec)
			time.Sleep(time.Duration(numsec) * time.Second)
		}
	}
	return conn, err
}
