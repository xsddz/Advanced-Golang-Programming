package router

import "google.golang.org/grpc"

func SetGRPCRouter(app *grpc.Server) {

	// 为 User 服务注册业务实现 将 User 服务绑定到 RPC 服务容器上
	// user.RegisterUserServer(app, &UserService{})
	// 注册反射服务 这个服务是CLI使用的 跟服务本身没有关系
	// reflection.Register(app)

}
