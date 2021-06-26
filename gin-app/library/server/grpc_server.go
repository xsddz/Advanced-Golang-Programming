package server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
)

func RCPServerMaker(routerSetter func(*grpc.Server)) ServerHandle {
	return func(wg *sync.WaitGroup) {
		defer func() {
			fmt.Println("rpc server end.")
			wg.Done()
		}()

		// 初始化 RPC app
		var opts []grpc.ServerOption
		app := grpc.NewServer(opts...)

		// 设置路由
		routerSetter(app)

		// 启动服务
		address := ":8081"
		fmt.Printf("Listening and serving GRPC on %s\n", address)
		lis, err := net.Listen("tcp", address)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}
		if err := app.Serve(lis); err != nil {
			log.Fatalf("failed to serve GRPC: %v", err)
		}
	}
}
