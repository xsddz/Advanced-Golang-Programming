package router

import (
	"gin-app/controller"
	"gin-app/entity/pb"
	"gin-app/library/server"
)

func GRPCRouter(app *server.App) {
	// 为 Foo 服务注册业务实现 将 Foo 服务绑定到 RPC 服务容器上
	pb.RegisterFooServer(app.ServerGRPC, &controller.FooController{})
}
