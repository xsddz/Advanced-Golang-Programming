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
	*grpc.Server

	address string
	routers []Router
}

// NewGRPCServer -
func NewGRPCServer(conf GRPCConf, routers ...Router) *GRPCServer {
	// GRPC server选项
	var opts []grpc.ServerOption

	address := fmt.Sprintf(":%v", conf.Port)
	return &GRPCServer{grpc.NewServer(opts...), address, routers}
}

// RegisterRouter 实现ServerI接口
func (s *GRPCServer) RegisterRouter(routers ...Router) {
	s.routers = append(s.routers, routers...)
}

// Run 实现ServerI接口
func (s *GRPCServer) Run(wg *sync.WaitGroup) {
	defer func() {
		fmt.Println("rpc server end.")
		wg.Done()
	}()

	// GRPC server路由设置
	for _, router := range s.routers {
		router()
	}

	// 注册反射服务，这个服务是 CLI 使用的，跟服务本身没有关系，可使用 grpcui 测试接口
	// - go get github.com/fullstorydev/grpcui/...
	// - go install github.com/fullstorydev/grpcui/cmd/grpcui
	// - grpcui -plaintext 127.0.0.1:8081
	reflection.Register(s)

	// 启动服务
	fmt.Printf("Listening and serving GRPC on %s\n", s.address)
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		panic(fmt.Errorf("failed to listen: %v", err))
	}
	if err := s.Serve(lis); err != nil {
		panic(fmt.Errorf("failed to serve GRPC: %v", err))
	}
}
