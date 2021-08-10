package server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRCPServer struct {
	routerSetter func(*App)
}

func NewGRPCServer(rs func(*App)) *GRCPServer {
	return &GRCPServer{
		routerSetter: rs,
	}
}

func (s *GRCPServer) Run(app *App, wg *sync.WaitGroup) {
	defer func() {
		fmt.Println("rpc server end.")
		wg.Done()
	}()

	// 初始化 GRPC server
	var opts []grpc.ServerOption
	srv := grpc.NewServer(opts...)

	// 设置路由
	app.ServerGRPC = srv
	s.routerSetter(app)

	// 注册反射服务 这个服务是CLI使用的 跟服务本身没有关系
	reflection.Register(srv)

	// 启动服务
	address := ":8081"
	fmt.Printf("Listening and serving GRPC on %s\n", address)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve GRPC: %v", err)
	}
}
