package router

import (
	"gin-app/controllers/grpc"
	"gin-app/entity/pb"
	"gin-app/library/app"
)

func GRPCRouter() {
	grpcServer := app.GetGRPCServer()

	// 为 Foobar 服务注册业务实现 将 Foobar 服务绑定到 RPC 服务容器上
	pb.RegisterFoobarServer(grpcServer, &grpc.Foobar{})
}
