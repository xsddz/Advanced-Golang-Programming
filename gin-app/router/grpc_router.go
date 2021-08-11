package router

import (
	"gin-app/controller"
	"gin-app/entity/pb"
	"gin-app/library/server"
)

func GRPCRouter(app *server.App) {
	// 为 Foobar 服务注册业务实现 将 Foobar 服务绑定到 RPC 服务容器上
	pb.RegisterFoobarServer(app.GRPCServer, &controller.FoobarController{})
}
