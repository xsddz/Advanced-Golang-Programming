package server

import (
	"fmt"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GRPCConf -
type GRPCConf struct {
	Port string `toml:"port"`
}

// GRPCServer -
type GRPCServer struct {
	conf   *GRPCConf
	router Router
}

// NewGRPCServer -
func NewGRPCServer(conf GRPCConf, router Router) *GRPCServer {
	return &GRPCServer{&conf, router}
}

// Run 实现ServerI接口
func (s *GRPCServer) Run(app *Engine, wg *sync.WaitGroup) {
	defer func() {
		fmt.Println("rpc server end.")
		wg.Done()
	}()

	// GRPC server初始化
	var opts []grpc.ServerOption
	srv := grpc.NewServer(opts...)
	app.SetGRPCServer(srv)

	// GRPC server路由设置
	s.router()

	// 注册反射服务，这个服务是 CLI 使用的，跟服务本身没有关系，可使用 grpcui 测试接口
	// - go get github.com/fullstorydev/grpcui/...
	// - go install github.com/fullstorydev/grpcui/cmd/grpcui
	// - grpcui -plaintext 127.0.0.1:8081
	reflection.Register(srv)

	// 启动服务
	address := fmt.Sprintf(":%v", s.conf.Port)
	fmt.Printf("Listening and serving GRPC on %s\n", address)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		panic(fmt.Errorf("failed to listen: %v", err))
	}
	if err := srv.Serve(lis); err != nil {
		panic(fmt.Errorf("failed to serve GRPC: %v", err))
	}
}
