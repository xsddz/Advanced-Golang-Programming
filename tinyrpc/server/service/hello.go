package service

import (
	"fmt"
	"net/rpc"
)

func init() {
	fmt.Println("======init hello service")
	rpc.RegisterName(HelloServiceName, new(HelloService))
}

// HelloServiceName HelloServiceName
const HelloServiceName = "HelloService"

// HelloService HelloService
type HelloService struct{}

// Hello Hello
func (hs *HelloService) Hello(request string, reply *string) error {
	*reply = "hello " + request
	return nil
}
