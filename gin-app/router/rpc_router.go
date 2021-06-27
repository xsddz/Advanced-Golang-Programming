package router

import (
	"gin-app/controller"
	"gin-app/protos/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func SetGRPCRouter(app *grpc.Server) {

	// 为 Simple 服务注册业务实现 将 Simple 服务绑定到 RPC 服务容器上
	pb.RegisterSimpleServer(app, &controller.SimpleService{})

	// 注册反射服务 这个服务是CLI使用的 跟服务本身没有关系
	reflection.Register(app)

}
