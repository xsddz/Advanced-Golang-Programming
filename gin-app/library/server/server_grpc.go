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
	router Router
}

func NewGRPCServer(r Router) *GRCPServer {
	return &GRCPServer{
		router: r,
	}
}

func (s *GRCPServer) Run(app *Engine, wg *sync.WaitGroup) {
	defer func() {
		fmt.Println("rpc server end.")
		wg.Done()
	}()

	// 初始化 GRPC server
	var opts []grpc.ServerOption
	srv := grpc.NewServer(opts...)
	app.GRPCServer = srv

	// 设置路由
	s.router()

	// 注册反射服务，这个服务是 CLI 使用的，跟服务本身没有关系，可使用 grpcui 测试接口
	// - go get github.com/fullstorydev/grpcui/...
	// - go install github.com/fullstorydev/grpcui/cmd/grpcui
	// - grpcui -plaintext 127.0.0.1:8081
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
