package router

import (
	"yawebapp/Illuminate/app"
	"yawebapp/controllers/grpc"
	"yawebapp/entity/pb"
)

func GRPCRouter() {
	grpcServer := app.GetGRPCServer()

	// 为 Foobar 服务注册业务实现 将 Foobar 服务绑定到 RPC 服务容器上
	pb.RegisterFoobarServer(grpcServer, &grpc.Foobar{})
}
